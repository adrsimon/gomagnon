package typing

import (
	"fmt"
)

type HumanStats struct {
	Strength    int
	Sociability int
}

type HumanBody struct {
	Age         int
	Gender      rune
	Hungriness  int
	Thirstiness int
}

type agentComm struct {
	AgentID string
	Action  string
	comm    chan agentComm
}

type Clan struct {
	members []string
	chief   string
}

type Human struct {
	id    string
	Type  rune
	Body  HumanBody
	Stats HumanStats

	Position       *Hexagone
	Target         *Hexagone
	MovingToTarget bool
	CurrentPath    []*Hexagone
	Board          *Board

	ComOut agentToManager
	ComIn  managerToAgent

	AgentRelation map[string]string
	AgentCommOut  agentComm
	AgentCommIn   chan agentComm
	Clan          *Clan
}

func NewHuman(id string, Type rune, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent) *Human {
	return &Human{id: id, Type: Type, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn, AgentRelation: make(map[string]string)}
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
)

func (h *Human) EvaluateOneHex(hex *Hexagone) float64 {
	var score = 0.0

	if hex == nil {
		return score
	}

	switch hex.Resource {
	case ANIMAL:
		score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 0.01
	case FRUIT:
		score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.01
	case ROCK:
		score += 0.5
	case WOOD:
		score += 0.5
	}

	if hex.Biome.BiomeType == WATER {
		score = (float64(h.Body.Thirstiness) / 100) * WaterValueMultiplier
	}

	return score
}

func (h *Human) GetNeighborsWithin5() []*Hexagone {
	neighbours := h.Board.GetNeighbours(h.Position)
	visited := make(map[*Hexagone]bool)
	for i := 0; i < 4; i++ {
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			if _, ok := visited[neighbour]; !ok {
				visited[neighbour] = true
				neighbours = append(neighbours, h.Board.GetNeighbours(neighbour)...)
			}
		}
	}

	return neighbours
}

func (h *Human) BestNeighbor(surroundingHexagons []*Hexagone) (*Hexagone, *Human) {
	best := 0.
	indexBest := -1
	listFriends := make([]*Human, 0)
	for i, v := range surroundingHexagons {
		score := h.EvaluateOneHex(v)
		if score > best {
			indexBest = i
			best = score
		}
		if len(listFriends) < 1 {
			for _, f := range v.Agents {
				_, ok := h.AgentRelation[f.id]
				if ok {
					listFriends = append(listFriends, f)
				}
			}
		}
	}

	if indexBest != -1 {
		return surroundingHexagons[indexBest], nil
	}

	valid := false
	randHex := &Hexagone{}
	for !valid {
		if h.Clan == nil {
			if len(listFriends) > 0 {
				return nil, listFriends[0]
			}
		} else {
			randHex = surroundingHexagons[r.Intn(len(surroundingHexagons))]
			if h.Board.isValidHex(randHex) && randHex.Biome.BiomeType != WATER {
				valid = true
			}
		}
	}

	return randHex, nil
}

func (h *Human) UpdateAgent() {
	select {
	case interruption := <-h.AgentCommIn:
		fmt.Printf("%s wants to %s", interruption.AgentID, interruption.Action)
	default:
		if len(h.Position.Agents) > 1 {
			fmt.Println("a")
			for _, v := range h.Position.Agents {
				fmt.Println("b")
				if v != h {
					fmt.Println("c")
					_, ok := h.AgentRelation[v.id]
					if !ok {
						fmt.Println("D")
						if r.Intn(2) >= 1 {
							h.AgentRelation[v.id] = "Friend"
							fmt.Println("friend")
						} else {
							h.AgentRelation[v.id] = "Ennemy"
							fmt.Println("ennemy")
						}
					}
				}
			}
		}
		if !h.MovingToTarget {
			surroundingHexagons := h.GetNeighborsWithin5()
			targetHexagon, possibleFriend := h.BestNeighbor(surroundingHexagons)
			if targetHexagon != nil {
				res := AStar(*h, targetHexagon)
				path := createPath(res, targetHexagon)
				h.CurrentPath = path
				h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-2]
				h.Target = targetHexagon
				h.MovingToTarget = true
			} else {
				possibleFriend.AgentCommIn <- agentComm{AgentID: h.id, Action: "inviteClan", comm: h.AgentCommIn}
			}

		}

		if h.MovingToTarget && len(h.CurrentPath) > 0 {
			nextHexagon := h.CurrentPath[len(h.CurrentPath)-1]
			h.MoveToHexagon(h.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-1]
		}

		if h.Target.Position == h.Position.Position {
			h.MovingToTarget = false
			h.Target = nil
			if h.Position.Resource != NONE {
				h.ComOut = agentToManager{AgentID: h.id, Action: "get", Pos: h.Position, commOut: make(chan managerToAgent)}
				h.Board.AgentManager.messIn <- h.ComOut
				h.ComIn = <-h.ComOut.commOut
				if h.ComIn.Valid {
					h.UpdateStateBasedOnResource(h.Position)
				}
			}
		}
	}
}

func createPath(maps map[*Hexagone]*Hexagone, hexagon *Hexagone) []*Hexagone {
	path := make([]*Hexagone, 0)
	path = append(path, hexagon)
	val, ok := maps[hexagon]
	for ok {
		path = append(path, val)
		val, ok = maps[val]
	}
	return path
}

func (h *Human) UpdateStateBasedOnResource(hex *Hexagone) {
	if hex.Resource == ANIMAL {
		h.Body.Hungriness = max(0, h.Body.Hungriness-r.Intn(20))
	}
	if hex.Resource == FRUIT {
		h.Body.Hungriness = max(0, h.Body.Hungriness-r.Intn(10))
	}
	if hex.Biome.BiomeType == WATER {
		h.Body.Thirstiness = max(0, h.Body.Thirstiness-r.Intn(30))
	}
}

func (h *Human) MoveToHexagon(hex *Hexagone) {
	h.Position = hex
}

func (h *Human) Start() {

}
