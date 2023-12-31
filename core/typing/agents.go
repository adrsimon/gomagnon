package typing

// check if everyone make childs, add age + childs behaviour,

import (
	"fmt"
	"math"
	"strconv"
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
	ASK4PROCREATE
	PROCREATE
	FIGHT
)

func (h *Agent) actionToStr() (action string) {
	switch h.Action {
	case NOOP:
		action = "NOOP"
	case MOVE:
		action = "MOVE"
	case GET:
		action = "GET"
	case BUILD:
		action = "BUILD"
	case SLEEP:
		action = "SLEEP"
	case STOREATHOME:
		action = "STOREATHOME"
	case EATFROMHOME:
		action = "EATFROMHOME"
	case CREATECLAN:
		action = "CREATECLAN"
	case CREATEVOTENEWMEMBER:
		action = "CREATEVOTENEWMEMBER"
	case VOTE:
		action = "VOTE"
	case GETRESULT:
		action = "GETRESULT"
	case LOOK4SOMEONE:
		action = "LOOK4SOMEONE"
	case PROCREATE:
		action = "PROCREATE"
	case ASK4PROCREATE:
		action = "ASK4PROCREATE"
	case FIGHT:
		action = "FIGHT"
	}
	return
}

type Race int

const (
	NEANDERTHAL Race = iota
	SAPIENS
)

type Inventory struct {
	Object map[ResourceType]int
	Weight float64
}

type Procreate struct {
	Partner *Agent
	Valide  bool
	Timer   int
	IsHome  bool
}

type Agent struct {
	ID       string
	Type     rune
	Race     Race
	Body     HumanBody
	Stats    HumanStats
	DeathAge float64

	Board     *Board
	MapVision [][]Hexagone

	Inventory Inventory

	Position       *Hexagone
	Target         *Hexagone
	MovingToTarget bool
	CurrentPath    []*Point2D

	Hut                *Hut
	HutInventoryVision []ResourceType

	ComOut agentToManager
	ComIn  managerToAgent

	Action          Action
	StackAction     StackAction
	Looking4Someone bool

	Neighbours    []*Agent
	AgentRelation map[string]string
	AgentCommIn   chan AgentComm
	Clan          *Clan
	Procreate     Procreate
	Terminated    bool

	String string

	Opponent      *Agent
	Fightcooldown int

	Behavior HumanActions
}

func (h *Agent) PerformAction() bool {
	randomNumber := Randomizer.Intn(101)
	return randomNumber <= h.Stats.Sociability+100
}

type AgentComm struct {
	Agent   *Agent
	Action  string
	commOut chan AgentComm
}

type Clan struct {
	ID      string
	members []*Agent
	chief   *Agent
}

const (
	AnimalFoodValueMultiplier = 5.0
	FruitFoodValueMultiplier  = 3.0
	WaterValueMultiplier      = 4.0
	DistanceMultiplier        = 0.5
)

func NewHuman(id string, Type rune, Race Race, body HumanBody, stats HumanStats, mapVision [][]Hexagone, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Point2D, board *Board, comOut agentToManager, comIn managerToAgent, hut *Hut, inventory Inventory, agentRelation map[string]string, procreate Procreate) *Agent {
	return &Agent{ID: id, Type: Type, Race: Race, Body: body, Stats: stats, MapVision: mapVision, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn, Hut: hut, Inventory: inventory, AgentRelation: agentRelation, Procreate: procreate}
}

func (h *Agent) EvaluateOneHex(hex *Hexagone) float64 {
	var score = 1.0
	threshold := 85.0

	dist := distance(*h.Position, *hex)
	score -= dist * DistanceMultiplier

	if hex.Biome == DEEP_WATER || hex.Biome == WATER {
		return math.Inf(-1)
	}
	if hex == nil {
		return score
	}

	switch hex.Resource {
	case ANIMAL:
		if h.Race == NEANDERTHAL {
			score += AnimalFoodValueMultiplier
		}
		if h.Race == SAPIENS {
			score += AnimalFoodValueMultiplier
		}
		if h.Body.Hungriness > threshold {
			score += 3
		}
	case FRUIT:
		if h.Race == NEANDERTHAL {
			score = math.Inf(-1)
		}
		if h.Race == SAPIENS {
			score += FruitFoodValueMultiplier + 0.5
		}
		if h.Body.Hungriness > threshold {
			score += 3
		}
	case ROCK:
		if h.Hut == nil && h.Inventory.Object[ROCK] < Needs["hut"][ROCK] && h.Inventory.Weight <= MaxWeightInv-WeightRock {
			score += 10
		} else if h.Hut == nil && h.Inventory.Object[ROCK] >= Needs["hut"][ROCK] {
			score = math.Inf(-1)
		}
	case WOOD:
		if h.Hut == nil && h.Inventory.Object[WOOD] < Needs["hut"][WOOD] && h.Inventory.Weight <= MaxWeightInv-WeightWood {
			score += 10
		} else if h.Hut == nil && h.Inventory.Object[WOOD] >= Needs["hut"][WOOD] {
			score = math.Inf(-1)
		}
	}

	for _, nb := range h.Board.GetNeighbours(*hex) {
		if nb.Biome == WATER && h.Hut == nil && h.Inventory.Object[WOOD] >= Needs["hut"][WOOD] && h.Inventory.Object[ROCK] >= Needs["hut"][ROCK] {
			score += 100
			break
		}
		if nb.Biome == WATER && h.Body.Thirstiness > threshold {
			score += WaterValueMultiplier + 0.5
			break
		}
	}

	return score
}

func (h *Agent) BestMatchHuman() *Agent {
	if len(h.Neighbours) == 0 {
		return nil
	}

	bestMatch := h.Neighbours[0]
	highestScore := h.calculateScore(bestMatch)

	for _, neighbour := range h.Neighbours[1:] {
		score := h.calculateScore(neighbour)
		if score > highestScore {
			bestMatch = neighbour
			highestScore = score
		}
	}

	return bestMatch
}

func (h *Agent) calculateScore(n *Agent) float64 {
	var score float64
	score += float64(n.Stats.Sociability / 10)
	score += float64(n.Stats.Strength / 10)
	if n.Type != h.Type {
		score += 5
	}
	if h.Clan != nil && len(h.Clan.members) < 4 {
		score += 2
	}
	if n.Race == h.Race {
		score += 4
	}
	return score
}

func (h *Agent) BestMove(surroundingHexagons []*Hexagone) *Hexagone {
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

	if indexBest != -1 && surroundingHexagons[indexBest].Position.X != h.Position.Position.X && surroundingHexagons[indexBest].Position.Y != h.Position.Position.Y {
		return surroundingHexagons[indexBest]
	}

	valid := false
	randHex := &Hexagone{}
	for !valid {
		n := len(surroundingHexagons)
		randHex = surroundingHexagons[Randomizer.Intn(n)]
		if h.Board.isValidHex(randHex) && randHex.Biome != DEEP_WATER && randHex.Biome != WATER {
			valid = true
		}
	}

	return randHex
}

func (h *Agent) MoveToHexagon(hex *Hexagone) {
	h.Position = hex
	h.Body.Hungriness += 0.2
	h.Body.Thirstiness += 0.4
	h.Body.Tiredness += 0.2
}

func (h *Agent) UpdateState(resource ResourceType) {
	switch resource {
	case ANIMAL:
		if h.Body.Hungriness > 85 || h.Hut == nil || h.Inventory.Weight >= MaxWeightInv-3*WeightAnimal {
			h.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			if h.Body.Hungriness < 0 {
				h.Body.Hungriness = 0
			}
			break
		} else {
			h.Inventory.Object[resource] += 3
			h.Inventory.Weight += 3 * WeightAnimal
		}
	case FRUIT:
		if h.Body.Hungriness > 85 || h.Hut == nil || h.Inventory.Weight >= MaxWeightInv-WeightFruit {
			h.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			if h.Body.Hungriness < 0 {
				h.Body.Hungriness = 0
			}
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

	neighbours := h.Board.GetNeighbours(*h.Position)
	for _, neighbour := range neighbours {
		if neighbour == nil {
			continue
		}
		if neighbour.Biome == WATER {
			h.Body.Thirstiness = 0
		}
	}
}

func (h *Agent) Perceive() {
	listHumans := make([]*Agent, 0)
	cases := make([]*Hexagone, 0)
	cases = append(cases, h.Position)
	cases = append(cases, h.Board.GetNeighbours(*h.Position)...)
	for _, v := range cases {
		for _, p := range v.Agents {
			if p != h && p.Body.Age >= 10 {
				_, ok := h.AgentRelation[p.ID]
				listHumans = append(listHumans, p)
				if !ok {
					var relation string
					h.Opponent = nil
					if h.Clan != nil && p.Clan == h.Clan {
						if Randomizer.Intn(100) < 10 {
							relation = "Enemy"
							if h.Opponent == nil {
								h.Opponent = p
							}
						} else {
							relation = "Friend"
						}
					} else {
						if Randomizer.Intn(100) < 50 {
							relation = "Enemy"
							if h.Opponent == nil {
								h.Opponent = p
							}
						} else {
							relation = "Friend"
						}
					}

					h.AgentRelation[p.ID] = relation
				}
			}
		}
	}
	h.Neighbours = listHumans
	if h.Procreate.Partner != nil {
		for _, neighbour := range h.Neighbours {
			if neighbour == h.Procreate.Partner && neighbour.Procreate.Partner == h && neighbour.Position.Position == h.Position.Position {
				h.Procreate.IsHome = true
				break
			}
		}
	}

	if h.Hut != nil && h.Procreate.Partner == nil && !h.Procreate.Valide && h.Procreate.Timer <= 0 && h.Clan != nil && h.PerformAction() && len(h.Clan.members) < 15 {
		for _, neighbour := range h.Neighbours {
			if neighbour.Clan == h.Clan && neighbour.Procreate.Partner == nil && neighbour.Hut == h.Hut && neighbour.Body.Age > 10 && h.Type != neighbour.Type && h.Procreate.Timer <= 100 {
				h.Procreate.Partner = neighbour
				break
			}
		}
	}

	if h.Hut != nil && h.Position.Position == h.Hut.Position.Position {
		h.HutInventoryVision = h.Hut.Inventory
	}

	for _, v := range h.Behavior.GetNeighboursWithinAcuity() {
		if v == nil {
			continue
		}
		h.MapVision[v.Position.X][v.Position.Y] = *v
	}
}

func (h *Agent) AnswerAgents(res AgentComm) {
	switch res.Action {
	case "CREATECLAN":
		score := h.calculateScore(res.Agent)
		if h.Clan != nil || score < 12 {
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
	case "PROCREATE":
		if math.Abs(float64(res.Agent.Stats.Sociability-h.Stats.Sociability)) < 50 && h.Procreate.Timer <= 0 {
			res.commOut <- AgentComm{Agent: h, Action: "ACCEPTPROCREATE", commOut: h.AgentCommIn}
			h.Procreate.Valide = true
			h.Procreate.Partner = res.Agent
		} else {
			res.commOut <- AgentComm{Agent: h, Action: "REFUSEPROCREATE", commOut: h.AgentCommIn}
		}
	case "FIGHT":
		h.Opponent = res.Agent
		SociabilityOpp := h.Opponent.Stats.Sociability
		SociabilityAg := h.Stats.Sociability
		if 1.25*float64(SociabilityOpp) > float64(SociabilityAg) {
			res.commOut <- AgentComm{Agent: h, Action: "OKFIGHT", commOut: h.AgentCommIn}
			res2 := <-h.AgentCommIn
			if res2.Action == "YOUWIN" {
				h.ComOut = agentToManager{AgentID: h.ID, Action: "transfer-inventory", Pos: h.Position, commOut: make(chan managerToAgent)}
				h.Board.AgentManager.messIn <- h.ComOut
				h.Opponent.AgentCommIn <- AgentComm{Agent: h, Action: "LOOTED", commOut: h.AgentCommIn}
			} else {
				res3 := <-h.AgentCommIn
				if res3.Action == "LOOTED" {
					h.ComOut = agentToManager{AgentID: h.ID, Action: "die", Pos: h.Position, commOut: make(chan managerToAgent)}
					h.Board.AgentManager.messIn <- h.ComOut
				}
			}
		} else {
			res.commOut <- AgentComm{Agent: h, Action: "NOFIGHT", commOut: h.AgentCommIn}
			h.Opponent = nil
		}
	}
}

func (h *Agent) IsDead() bool {
	return h.Body.Hungriness >= 100 || h.Body.Thirstiness >= 100 || h.Body.Tiredness >= 100 || h.Body.Age >= h.DeathAge
}

func (h *Agent) CloseUpdate() {
	if h.IsDead() {
		h.ComOut = agentToManager{AgentID: h.ID, Action: "die", Pos: h.Position, commOut: make(chan managerToAgent)}
		h.Board.AgentManager.messIn <- h.ComOut
	} else {
		h.UpdateState(NONE)
		h.Body.Age += 0.01
		h.Body.Hungriness += 0.1
		h.Body.Thirstiness += 0.2
		h.Body.Tiredness += 0.1
		h.Procreate.Timer -= 1
		h.Fightcooldown -= 1
	}
	if int(h.Body.Age) == 10 {
		h.Behavior = &HumanBehavior{H: h}
		if Randomizer.Intn(200) < 1 && h.Clan != nil && len(h.Clan.members) > 5 {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "leave-clan", Pos: h.Position, commOut: make(chan managerToAgent)}
			h.Board.AgentManager.messIn <- h.ComOut
			h.ComIn = <-h.ComOut.commOut
		}
	}
	if h.Body.Tiredness < 0 {
		h.Body.Tiredness = 0
	}

}

func (h *Agent) UpdateAgent() {
	h.Terminated = false
	h.Perceive()
	h.Behavior.Deliberate()
	h.Behavior.Act()
	select {
	case res := <-h.AgentCommIn:
		h.AnswerAgents(res)
		h.Terminated = true
	case <-time.After(5 * time.Millisecond):
		h.Terminated = true
	}
	h.CloseUpdate()
	h.String = h.ToString()
}

func (h *Agent) ToString() string {
	race := "Neanderthal"
	if h.Race == SAPIENS {
		race = "Sapiens"
	}

	str := h.ID + " - " + race + h.Position.ToString() + "\n\n"
	str += "--- Body ---" + "\n"
	str += "Age : " + fmt.Sprintf("%d", int(h.Body.Age)) + "\n"
	str += "Hungriness : " + fmt.Sprintf("%d", int(h.Body.Hungriness)) + "\n"
	str += "Thirstiness : " + fmt.Sprintf("%d", int(h.Body.Thirstiness)) + "\n"
	str += "Tiredness : " + fmt.Sprintf("%d", int(h.Body.Tiredness)) + "\n\n"
	str += "--- Hut and Clan ---\n"
	if h.Hut != nil {
		str += "Hut pos : " + strconv.Itoa(h.Hut.Position.Position.X) + " " + strconv.Itoa(h.Hut.Position.Position.Y) + "\n"
	} else {
		str += "No hut" + "\n"
	}
	if h.Clan != nil {
		str += "Clan ID : " + h.Clan.ID + "\n"
		str += "Chief : " + h.Clan.chief.ID + "\n"
		str += "Members : " + strconv.Itoa(len(h.Clan.members)) + "\n\n"
	} else {
		str += "No clan" + "\n\n"
	}
	str += "--- Inventory ---" + "\n"
	str += "Fruits : " + strconv.Itoa(h.Inventory.Object[FRUIT]) + "\n"
	str += "Animals : " + strconv.Itoa(h.Inventory.Object[ANIMAL]) + "\n"
	str += "Woods : " + strconv.Itoa(h.Inventory.Object[WOOD]) + "\n"
	str += "Rocks : " + strconv.Itoa(h.Inventory.Object[ROCK]) + "\n\n"

	str += "Doing : " + h.actionToStr()
	return str
}
