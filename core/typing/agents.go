package typing

// check if everyone make childs, add age + childs behaviour,

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

type HumanStats struct {
	Strength    int
	Sociability int
	Acuity      int
}

type HumanBody struct {
	Age    float64
	Gender rune

	Hungriness  float64
	Thirstiness float64

	Tiredness float64
}

type Action int

type StackAction []Action

const (
	NOOP Action = iota
	MOVE
	GET
	BUILD
	SLEEP
	STOREATHOME
	EATFROMHOME
	CREATECLAN
	CREATEVOTENEWMEMBER
	VOTE
	GETRESULT
	LOOK4SOMEONE
	PROCREATE
)

type Race int

const (
	NEANDERTHAL Race = iota
	SAPIENS
)

const (
	MaxWeightInv = 10.0
	WeightRock   = 2.0
	WeightWood   = 1.0
	WeightAnimal = 0.5
	WeightFruit  = 0.1
)

type Inventory struct {
	Object map[ResourceType]int
	Weight float64
}

type Procreate struct {
	Partner *Human
	Timer   int
	isHome  bool
}

type Human struct {
	ID    string
	Type  rune
	Race  Race
	Body  HumanBody
	Stats HumanStats

	Inventory Inventory

	Position       *Hexagone
	Target         *Hexagone
	MovingToTarget bool
	CurrentPath    []*Hexagone
	Board          *Board

	Hut                *Hut
	HutInventoryVision []ResourceType

	ComOut agentToManager
	ComIn  managerToAgent

	Action          Action
	StackAction     StackAction
	Looking4Someone bool

	Neighbours    []*Human
	AgentRelation map[string]string
	AgentCommIn   chan AgentComm
	Clan          *Clan
	Procreate     Procreate
	Terminated    bool
}

func (h *Human) PerformAction() bool {
	randomNumber := Randomizer.Intn(101)
	return randomNumber <= h.Stats.Sociability
}

type AgentComm struct {
	Agent   *Human
	Action  string
	commOut chan AgentComm
}

type Clan struct {
	ID      string
	members []*Human
	chief   *Human
}

const (
	AnimalFoodValueMultiplier = 5.0
	FruitFoodValueMultiplier  = 3.0
	WaterValueMultiplier      = 2.0
	DistanceMultiplier        = 0.2
)

func NewHuman(id string, Type rune, Race Race, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent, hut *Hut, inventory Inventory, agentRelation map[string]string) *Human {
	return &Human{ID: id, Type: Type, Race: Race, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn, Hut: hut, Inventory: inventory, AgentRelation: agentRelation}
}

func (h *Human) BestMatchHuman() *Human {
	if len(h.Neighbours) == 0 {
		return nil
	}

	bestMatch := h.Neighbours[0]
	highestScore := calculateScore(h, bestMatch)

	for _, neighbour := range h.Neighbours[1:] {
		score := calculateScore(h, neighbour)
		if score > highestScore {
			bestMatch = neighbour
			highestScore = score
		}
	}

	return bestMatch
}

func calculateScore(h, n *Human) float64 {
	var score float64
	score += float64(n.Stats.Sociability / 100)
	score += float64(n.Stats.Strength / 100)
	if n.Type != h.Type && h.Clan != nil && len(h.Clan.members) < 4 {
		score += 2
	}
	if n.Race == h.Race {
		score += 1
	}
	return score
}

func (h *Human) EvaluateOneHex(hex *Hexagone) float64 {
	var score = 0.0
	threshold := 85.0

	dist := distance(*h.Position, *hex)
	score -= dist * DistanceMultiplier

	if hex.Biome == WATER {
		return math.Inf(-1)
	}
	if hex == nil {
		return score
	}

	switch hex.Resource {
	case ANIMAL:
		if h.Race == NEANDERTHAL {
			score += AnimalFoodValueMultiplier + 0.5
			score += AnimalFoodValueMultiplier + 0.5
		}
		if h.Race == SAPIENS {
			score += AnimalFoodValueMultiplier + 1.0
			score += AnimalFoodValueMultiplier + 1.0
		}
		if h.Body.Hungriness > threshold {
			score += 3
		}
	case FRUIT:
		if h.Race == NEANDERTHAL {
			score += FruitFoodValueMultiplier + 0.01
		}
		if h.Race == SAPIENS {
			score += FruitFoodValueMultiplier + 0.5
		}
		if h.Body.Hungriness > threshold {
			score += 3
		}
	case ROCK:
		if h.Hut == nil && h.Inventory.Object[ROCK] < Needs["hut"][ROCK] && h.Inventory.Weight <= MaxWeightInv-WeightRock {
			score += 3
		} else if (h.Hut != nil || h.Inventory.Object[ROCK] > Needs["hut"][ROCK]) && h.Inventory.Weight <= MaxWeightInv-WeightRock {
			score += 0.5
		} else {
			score -= 1
		}
	case WOOD:
		if h.Hut == nil && h.Inventory.Object[WOOD] < Needs["hut"][WOOD] && h.Inventory.Weight <= MaxWeightInv-WeightWood {
			score += 3
		} else if (h.Hut != nil || h.Inventory.Object[WOOD] > Needs["hut"][WOOD]) && h.Inventory.Weight <= MaxWeightInv-WeightWood {
			score += 0.5
		} else {
			score -= 1
		}
	}

	for _, nb := range h.Board.GetNeighbours(hex) {
		if nb.Biome == WATER && h.Body.Thirstiness > threshold {
			score += (float64(h.Body.Thirstiness)/100)*WaterValueMultiplier + 0.5
			break
		}
	}

	return score
}

func (h *Human) GetNeighboursWithinAcuity() []*Hexagone {
	neighbours := h.Board.GetNeighbours(h.Position)
	visited := make(map[*Hexagone]bool)
	for i := 1; i < h.Stats.Acuity; i++ {
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
		score := h.EvaluateOneHex(v)
		if score > best {
			indexBest = i
			best = score
		}
	}

	if indexBest != -1 && surroundingHexagons[indexBest] != h.Position {
		return surroundingHexagons[indexBest]
	}

	valid := false
	randHex := &Hexagone{}
	for !valid {
		randHex = surroundingHexagons[Randomizer.Intn(len(surroundingHexagons))]
		if h.Board.isValidHex(randHex) && randHex.Biome != WATER {
			valid = true
		}
	}

	return randHex
}

func (h *Human) MoveToHexagon(hex *Hexagone) {
	h.Position = hex
	h.Body.Hungriness += 0.5
	h.Body.Thirstiness += 1
	h.Body.Tiredness += 0.5
}

func (h *Human) UpdateState(resource ResourceType) {
	switch resource {
	case ANIMAL:
		if h.Body.Hungriness > 85 || h.Hut == nil || h.Inventory.Weight >= MaxWeightInv-3*WeightAnimal {
			h.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			break
		} else {
			h.Inventory.Object[resource] += 3
			h.Inventory.Weight += 3 * WeightAnimal
		}
	case FRUIT:
		if h.Body.Hungriness > 85 || h.Hut == nil || h.Inventory.Weight >= MaxWeightInv-WeightFruit {
			h.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			break
		} else {
			h.Inventory.Object[resource] += 1
			h.Inventory.Weight += WeightFruit
		}
	case ROCK:
		h.Inventory.Object[resource] += 1
		h.Inventory.Weight += WeightRock
	case WOOD:
		h.Inventory.Object[resource] += 1
		h.Inventory.Weight += WeightWood
	}

	neighbours := h.Board.GetNeighbours(h.Position)
	for _, neighbour := range neighbours {
		if neighbour == nil {
			continue
		}
		if neighbour.Biome == WATER {
			h.Body.Thirstiness = 0
		}
	}
}

func (h *Human) Perceive() {
	listHumans := make([]*Human, 0)
	cases := make([]*Hexagone, 0)
	cases = append(cases, h.Position)
	cases = append(cases, h.Board.GetNeighbours(h.Position)...)
	for _, v := range cases {
		for _, p := range v.Agents {
			if p != h {
				_, ok := h.AgentRelation[p.ID]
				listHumans = append(listHumans, p)
				if !ok {
					if Randomizer.Intn(2) >= 1 {
						h.AgentRelation[p.ID] = "Friend"
					} else {
						h.AgentRelation[p.ID] = "Enemy"
					}
				}
			}
		}
	}
	h.Neighbours = listHumans
	if h.Hut != nil && h.Procreate.Partner == nil && h.Procreate.Timer <= 0 && h.Clan != nil && h.PerformAction() {
		for _, neighbour := range h.Neighbours {
			if neighbour.Clan == h.Clan && neighbour.Procreate.Partner == nil && neighbour.Hut == h.Hut && neighbour.Body.Age > 15 /*&& h.Type != neighbour.Type */ {
				h.Procreate.Partner = neighbour
				neighbour.Procreate.Partner = h
				break
			}
		}

	} else if h.Hut != nil && h.Procreate.Partner != nil && h.Position.Position == h.Hut.Position.Position {
		h.ComOut = agentToManager{AgentID: h.ID, Action: "isHome", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		h.Procreate.isHome = h.ComIn.Valid
	}

	if h.Hut != nil && h.Position.Position == h.Hut.Position.Position {
		h.HutInventoryVision = h.Hut.Inventory
	}
}

func (h *Human) DeliberateAtHut() {
	/** If he is tired and have a home, he should sleep **/
	if h.Body.Tiredness > 0 {
		h.Action = SLEEP
		return
	}
	/** If he is home and not partner he should wait **/
	if h.Procreate.Partner != nil && !h.Procreate.isHome {
		h.Action = SLEEP
		return
	}

	if h.Procreate.Partner == nil {
		h.Action = MOVE
	}

	/** If he is home with partner he should procreate **/
	if h.Procreate.Partner != nil && h.Procreate.isHome {
		h.Action = PROCREATE
		return
	}

	/** If he is hungry and have food in home, he should eat **/
	if h.Body.Hungriness > 80 {
		if slices.Contains(h.HutInventoryVision, ANIMAL) || slices.Contains(h.HutInventoryVision, FRUIT) {
			h.Action = EATFROMHOME
			return
		} else {
			h.Action = MOVE
		}
	}

	/** If he has stuff in inventory, he should store it **/
	if h.Inventory.Weight > 0 {
		h.Action = STOREATHOME
		return
	}

	if h.Clan != nil && h.Clan.chief == h && len(h.Clan.members) < 15 && len(h.Clan.members) > 0 && !h.Hut.Ballot.VoteInProgress && !h.Looking4Someone {
		h.Action = CREATEVOTENEWMEMBER
		return
	}
	if h.Clan != nil && h.Hut.Ballot.VoteInProgress && slices.Contains(h.Hut.Ballot.VotersID, h.ID) {
		h.Action = VOTE
		return
	}
	if h.Clan != nil && h.Clan.chief == h && h.Hut.Ballot.VoteInProgress && h.Hut.Ballot.EndTimeVote.Before(time.Now()) {
		h.Action = GETRESULT
		return
	}

}

func (h *Human) Deliberate() {
	h.Action = NOOP

	/** Stacked actions **/
	if len(h.StackAction) > 0 {
		h.Action = h.StackAction[0]
		h.StackAction = h.StackAction[1:]
		return
	}

	/** Early game actions **/
	if h.Hut == nil {
		/** if the agent is tired and don't have a home, he should rest in the nature **/
		if h.Body.Tiredness > 80 {
			h.StackAction = append(h.StackAction, NOOP)
			h.StackAction = append(h.StackAction, NOOP)
			h.StackAction = append(h.StackAction, NOOP)
			return
		}

		/** if he can build a home and don't have one, he should build it **/
		if h.Inventory.Object[WOOD] >= Needs["hut"][WOOD] && h.Inventory.Object[ROCK] >= Needs["hut"][ROCK] {
			for _, v := range h.Board.GetNeighbours(h.Position) {
				if v == nil {
					continue
				}
				if v.Biome == WATER {
					h.Action = BUILD
					return
				}
			}
		}
	}

	/** In Hut actions **/
	if h.Hut != nil && h.Position.Position == h.Hut.Position.Position {
		h.DeliberateAtHut()
		if h.Action != NOOP {
			return
		}
	}

	/** General actions **/
	if h.Body.Thirstiness > 80 || h.Body.Hungriness > 80 {
		if !h.MovingToTarget {
			h.Action = MOVE
			return
		}
	}

	if h.Hut != nil {
		if len(h.Neighbours) > 0 && h.Clan == nil {
			h.Action = CREATECLAN
			return
		}
	}

	if h.Procreate.Partner != nil && h.Position != h.Hut.Position {
		h.Action = MOVE
		return
	}

	if h.Clan != nil && h.Clan.chief == h && h.Looking4Someone {
		h.Action = LOOK4SOMEONE
		return
	}

	if !h.MovingToTarget {
		h.Action = MOVE
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

			if h.Hut != nil {
				if h.Body.Tiredness > 80 || h.Procreate.Partner != nil {
					targetHexagon = h.Hut.Position
				} else if h.Body.Hungriness > 80 && (slices.Contains(h.HutInventoryVision, ANIMAL) || slices.Contains(h.HutInventoryVision, FRUIT)) {
					targetHexagon = h.Hut.Position
				}
			}

			if targetHexagon == nil {
				surroundingHexagons := h.GetNeighboursWithinAcuity()
				targetHexagon = h.BestNeighbor(surroundingHexagons)
			}

			res := AStar(*h, targetHexagon)
			h.CurrentPath = createPath(res, targetHexagon)
			if len(h.CurrentPath) < 2 {
				h.CurrentPath = nil
				break
			}
			h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-2]
			h.Target = targetHexagon
			h.MovingToTarget = true
		}

		if h.MovingToTarget && len(h.CurrentPath) > 0 {
			nextHexagon := h.CurrentPath[len(h.CurrentPath)-1]
			h.MoveToHexagon(h.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-1]
		}

		/** Next move stacking **/
		if h.MovingToTarget && len(h.CurrentPath) > 0 {
			h.StackAction = append(h.StackAction, MOVE)
		}

		if h.Position.Position == h.Target.Position {
			if h.Target.Resource != NONE {
				h.StackAction = append(h.StackAction, GET)
			}
			h.Target = nil
			h.MovingToTarget = false
		}
	case GET:
		if h.Position.Resource != NONE {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "get", Pos: h.Position, commOut: make(chan managerToAgent)}
			h.Board.AgentManager.messIn <- h.ComOut
			h.ComIn = <-h.ComOut.commOut
			if h.ComIn.Valid {
				h.UpdateState(h.ComIn.Resource)
			}
		}
	case BUILD:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "build", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			h.Hut = &Hut{Position: h.Position, Inventory: make([]ResourceType, 0), Owner: h}
			h.Inventory.Object[WOOD] -= Needs["hut"][WOOD]
			h.Inventory.Object[ROCK] -= Needs["hut"][ROCK]
			h.Inventory.Weight -= WeightWood * float64(Needs["hut"][WOOD])
			h.Inventory.Weight -= WeightRock * float64(Needs["hut"][ROCK])
		}
	case SLEEP:
		if h.Body.Tiredness > 0 {
			h.Body.Tiredness -= 3
			// h.Body.Hungriness += 0.5
			// h.Body.Thirstiness += 0.5
			h.StackAction = append(h.StackAction, SLEEP)
		}
	case STOREATHOME:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "store-at-home", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			h.Inventory.Weight = 0
		}
	case EATFROMHOME:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "eat-from-home", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			if h.ComIn.Resource == ANIMAL {
				h.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			} else {
				h.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			}
		}
	case CREATECLAN:
		var bestH *Human
		if len(h.Neighbours) > 1 {
			bestH = h.BestMatchHuman()
		} else if len(h.Neighbours) == 1 {
			bestH = h.Neighbours[0]
		} else {
			h.StackAction = append(h.StackAction, MOVE)
			break
		}
		if !bestH.Terminated {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: h, Action: "CREATECLAN", commOut: h.AgentCommIn}:
				select {
				case res := <-h.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						clanID := fmt.Sprintf("clan-%s", strings.Split(h.ID, "-")[1])
						clan := &Clan{ID: clanID, members: []*Human{bestH}, chief: h}
						h.Clan = clan
						bestH.AgentCommIn <- AgentComm{Agent: h, Action: "INVITECLAN", commOut: h.AgentCommIn}
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case CREATEVOTENEWMEMBER:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "VoteNewPerson", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "VoteYes", Pos: h.Position, commOut: make(chan managerToAgent)}
			h.Board.AgentManager.messIn <- h.ComOut
			h.ComIn = <-h.ComOut.commOut
		}
	case VOTE:
		if h.PerformAction() {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "VoteYes", Pos: h.Position, commOut: make(chan managerToAgent)}
		} else {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "VoteNo", Pos: h.Position, commOut: make(chan managerToAgent)}
		}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut

	case GETRESULT:
		h.ComOut = agentToManager{AgentID: h.ID, Action: "GetResult", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
		h.ComIn = <-h.ComOut.commOut
		if h.ComIn.Valid {
			h.Looking4Someone = true
		} else {
			h.Looking4Someone = false
		}
	case LOOK4SOMEONE:
		var bestH *Human
		if len(h.Neighbours) > 1 {
			//TO DEVELOPP bestH=find bestMatchHuman(humans)
			bestH = h.Neighbours[0] // waiting function
		} else if len(h.Neighbours) == 1 {
			bestH = h.Neighbours[0]
		} else {
			h.StackAction = append(h.StackAction, MOVE)
			break
		}
		if !bestH.Terminated {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: h, Action: "INVITECLAN", commOut: h.AgentCommIn}:
				select {
				case res := <-h.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						h.Looking4Someone = false
					} else {
						h.Action = MOVE
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case PROCREATE:
		if h.Type == 'F' {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "procreate", Pos: h.Position, commOut: make(chan managerToAgent)}
			h.Board.AgentManager.messIn <- h.ComOut
		}
	default:
		fmt.Println("Should not be here")
	}

}

func (h *Human) AnswerAgents(res AgentComm) {
	switch res.Action {
	case "CREATECLAN":
		if h.Clan != nil {
			res.commOut <- AgentComm{Agent: h, Action: "REFUSECLAN", commOut: h.AgentCommIn}
		} else {
			res.commOut <- AgentComm{Agent: h, Action: "ACCEPTCLAN", commOut: h.AgentCommIn}
			res2 := <-h.AgentCommIn
			if res2.Action == "INVITECLAN" {
				h.Clan = res2.Agent.Clan
				if h.Hut != nil && h.Hut.Owner != nil {
					h.ComOut = agentToManager{AgentID: h.ID, Action: "leave-house", Pos: h.Position, commOut: make(chan managerToAgent)}
					h.Board.AgentManager.messIn <- h.ComOut
					h.ComIn = <-h.ComOut.commOut
					if h.ComIn.Valid {
						h.Hut = nil
					}
				}
				h.Hut = res2.Agent.Hut
			}
		}
	case "INVITECLAN":
		if h.Clan != nil {
			res.commOut <- AgentComm{Agent: h, Action: "REFUSECLAN", commOut: h.AgentCommIn}
		} else {
			res.commOut <- AgentComm{Agent: h, Action: "ACCEPTCLAN", commOut: h.AgentCommIn}
			h.Clan = res.Agent.Clan
			if h.Hut != nil && h.Hut.Owner != nil {
				h.ComOut = agentToManager{AgentID: h.ID, Action: "leave-house", Pos: h.Position, commOut: make(chan managerToAgent)}
				h.Board.AgentManager.messIn <- h.ComOut
				h.ComIn = <-h.ComOut.commOut
				if h.ComIn.Valid {
					h.Hut = nil
				}
			}
			h.Hut = res.Agent.Hut
		}
	}
}

func (h *Human) IsDead() bool {
	return h.Body.Hungriness >= 100 || h.Body.Thirstiness >= 100 || h.Body.Tiredness >= 100
}

func (h *Human) CloseUpdate() {
	if h.IsDead() {
		h.ComOut = agentToManager{AgentID: h.ID, Action: "die", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
	} else {
		h.UpdateState(NONE)
		h.Body.Age += 0.05
		h.Procreate.Timer -= 1
		h.Body.Hungriness += 0.2
		h.Body.Thirstiness += 0.4
		h.Body.Tiredness += 0.4
	}
}

func (h *Human) UpdateAgent() {
	h.Terminated = false
	h.Perceive()
	h.Deliberate()
	h.Act()
	select {
	case res := <-h.AgentCommIn:
		h.Terminated = true
		h.AnswerAgents(res)
	default:
		h.Terminated = true
	}
	h.CloseUpdate()
}
