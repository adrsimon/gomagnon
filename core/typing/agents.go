package typing

import "fmt"

type HumanStats struct {
	Strength    int
	Sociability int
}

type HumanBody struct {
	Age    int
	Gender rune

	Hungriness  int
	Thirstiness int

	Tiredness float64
	Sleeping  bool
}

type Race string

const (
	Neandertal Race = "Neandertal"
	Sapiens    Race = "Sapiens"

type Action int

const (
	NOOP Action = iota
	MOVE
	GET
	BUILD
	SLEEP
)

type Human struct {
	ID    string
	Type  rune
	Race  Race
	Body  HumanBody
	Stats HumanStats

	Inventory map[ResourceType]int

	Position       *Hexagone
	Target         *Hexagone
	MovingToTarget bool
	CurrentPath    []*Hexagone
	Board          *Board

	Hut *Hut

	ComOut agentToManager
	ComIn  managerToAgent
  
  Action Action
}

func NewHuman(id string, Type rune, Race Race, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent, hut *Hut, inventory map[ResourceType]int) *Human {
	return &Human{ID: id, Type: Type, Race: Race, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn, Hut: hut, Inventory: inventory}
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
	DistanceMultiplier        = 1
)

func NewHuman(id string, Type rune, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent, hut *Hut, inventory map[ResourceType]int) *Human {
	return &Human{ID: id, Type: Type, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn, Hut: hut, Inventory: inventory}
}

func (h *Human) EvaluateOneHex(hex *Hexagone) float64 {
	var score = 0.0
	threshold := 85

	distance := distance(*h.Position, *hex)
	score -= distance * DistanceMultiplier

	if hex.Biome.BiomeType == WATER {
		return -1
	}
	if hex == nil {
		return score
	}

	if h.Hut == nil {
		switch hex.Resource {
		case ANIMAL:
			if h.Race == Neandertal {
				score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 0.5
			}
			if h.Race == Sapiens {
				score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 1.0
			}
			if h.Body.Hungriness > threshold {
				score += 1
			}
		case FRUIT:
			if h.Race == Neandertal {
				score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.01
			}
			if h.Race == Sapiens {
				score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.3
			}
			if h.Body.Hungriness > threshold {
				score += 1
			}
		case ROCK:
			score += 3
		case WOOD:
			score += 3
		}
	} else {
		switch hex.Resource {
		case ANIMAL:
			if h.Race == Neandertal {
				score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 0.5
			}
			if h.Race == Sapiens {
				score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 1.0
			}
			if h.Body.Hungriness > threshold {
				score += 1
			}
		case FRUIT:
			if h.Race == Neandertal {
				score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.01
			}
			if h.Race == Sapiens {
				score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.3
			}
			if h.Body.Hungriness > threshold {
				score += 1
			}
		case ROCK:
			score += 0.5
		case WOOD:
			score += 0.5
		}
	}

	for _, nb := range h.Board.GetNeighbours(hex) {
		if nb.Biome.BiomeType == WATER && h.Body.Thirstiness > threshold {
			score += (float64(h.Body.Thirstiness)/100)*WaterValueMultiplier + 0.5
			break
		}
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

func (h *Human) BestNeighbor(surroundingHexagons []*Hexagone) *Hexagone {
	if h.Body.Tiredness > 70 && h.Hut != nil {
		return h.Hut.Position
	}

	best := 0.
	indexBest := -1
	for i, v := range surroundingHexagons {
		if v.Position == h.Position.Position {
			continue
		}
		score := h.EvaluateOneHex(v)
		if score > best {
			indexBest = i
			best = score
		}
	}

	if indexBest != -1 {
		return surroundingHexagons[indexBest]
	}

	valid := false
	randHex := &Hexagone{}
	for !valid {
		randHex = surroundingHexagons[r.Intn(len(surroundingHexagons))]
		if h.Board.isValidHex(randHex) && randHex.Biome.BiomeType != WATER {
			valid = true
		}
	}

	return randHex
}

func (h *Human) MoveToHexagon(hex *Hexagone) {
	h.Position = hex
	h.Body.Hungriness += 1
	h.Body.Thirstiness += 2
	h.Body.Tiredness += 0.5
}

func (h *Human) UpdateState(resource ResourceType) {
	switch resource {
	case ANIMAL:
		h.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
	case FRUIT:
		h.Body.Hungriness -= 10 * FruitFoodValueMultiplier
	case ROCK, WOOD:
		h.Inventory[resource] += 1
	}

	neighbours := h.Board.GetNeighbours(h.Position)
	for _, neighbour := range neighbours {
		if neighbour == nil {
			continue
		}
		if neighbour.Biome.BiomeType == WATER {
			h.Body.Thirstiness -= 10
		}
	}
}

func (h *Human) Perceive() {}

func (h *Human) Deliberate() {
	h.Action = NOOP
	if h.Hut == nil && h.Inventory[WOOD] >= Needs["hut"][WOOD] && h.Inventory[ROCK] >= Needs["hut"][ROCK] {
		h.Action = BUILD
		return
	}

	if h.Hut != nil && h.Position.Position == h.Hut.Position.Position && h.Body.Tiredness > 0 {
		h.Action = SLEEP
		return
	}

	if !h.MovingToTarget {
		h.Action = MOVE
		return
	}
	if h.MovingToTarget && len(h.CurrentPath) > 0 {
		h.Action = MOVE
		return
	}
	if h.Target.Position == h.Position.Position {
		h.Action = GET
		return
	}
}

func (h *Human) Act() {
	switch h.Action {
	case NOOP:
		h.Body.Tiredness -= 1
	case MOVE:
		if !h.MovingToTarget {
			var targetHexagon *Hexagone

			if h.Body.Tiredness > 70 && h.Hut != nil {
				targetHexagon = h.Hut.Position
			} else {
				surroundingHexagons := h.GetNeighborsWithin5()
				targetHexagon = h.BestNeighbor(surroundingHexagons)
			}

			res := AStar(*h, targetHexagon)
			h.CurrentPath = createPath(res, targetHexagon)
			h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-2]
			h.Target = targetHexagon
			h.MovingToTarget = true
		}

		if h.MovingToTarget && len(h.CurrentPath) > 0 {
			nextHexagon := h.CurrentPath[len(h.CurrentPath)-1]
			h.MoveToHexagon(h.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			h.UpdateState(h.Position.Resource)
			h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-1]
		}
	case GET:
		if h.Target.Position == h.Position.Position {
			h.MovingToTarget = false
			h.Target = nil
			if h.Position.Resource != NONE {
				h.ComOut = agentToManager{AgentID: h.ID, Action: "get", Pos: h.Position, commOut: make(chan managerToAgent)}
				h.Board.AgentManager.messIn <- h.ComOut
				h.ComIn = <-h.ComOut.commOut
				if h.ComIn.Valid {
					h.UpdateState(h.ComIn.Resource)
				}
			}
		}
	case BUILD:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "build", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			h.Hut = &Hut{Position: h.Position, Inventory: make(map[ResourceType]int)}
			h.Inventory[WOOD] -= Needs["hut"][WOOD]
			h.Inventory[ROCK] -= Needs["hut"][ROCK]
		}
	case SLEEP:
		if h.Body.Sleeping && h.Body.Tiredness > 0 {
			h.Body.Tiredness -= 5
			h.Body.Hungriness += 1
			h.Body.Thirstiness += 1
		} else if h.Body.Sleeping && h.Body.Tiredness <= 0 {
			h.Body.Sleeping = false
		} else if !h.Body.Sleeping {
			h.Body.Sleeping = true
		}
	default:
		fmt.Println("Should not be here")
	}
}

func (h *Human) UpdateAgent() {
	h.Perceive()
	h.Deliberate()
	h.Act()
	h.UpdateState(NONE)
}
