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
}

func NewHuman(id string, Type rune, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent) *Human {
	return &Human{id: id, Type: Type, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn}
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

func (h *Human) BestNeighbor(surroundingHexagons []*Hexagone) *Hexagone {
	best := 0.
	indexBest := -1
	for i, v := range surroundingHexagons {
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
		if h.Board.isValidHex(randHex) {
			valid = true
		}
	}

	return randHex
}

func (h *Human) UpdateAgent() {
	if !h.MovingToTarget {
		fmt.Println("Looking for target")
		surroundingHexagons := h.GetNeighborsWithin5()
		targetHexagon := h.BestNeighbor(surroundingHexagons)

		h.CurrentPath = AStar(*h, targetHexagon)
		h.Target = targetHexagon
		h.MovingToTarget = true
		fmt.Println("New target:", targetHexagon.ToString())
	}

	if h.MovingToTarget && len(h.CurrentPath) > 0 {
		nextHexagon := h.CurrentPath[0]
		h.MoveToHexagon(nextHexagon)
		h.CurrentPath = h.CurrentPath[1:]
		h.UpdateStateBasedOnResource(nextHexagon)
		if nextHexagon == h.Target {
			h.MovingToTarget = false
		}
	}

	if h.Target.Position == h.Position.Position {
		fmt.Println("Reached target")
		h.MovingToTarget = false
		h.Target = nil
		h.Board.Cases[h.Position.ToString()].Resource = NONE
	}
}

func (h *Human) UpdateStateBasedOnResource(hex *Hexagone) {
	if hex.Resource == ANIMAL {
		// TODO: send request to take resource and if yes:
		h.Body.Hungriness = max(0, h.Body.Hungriness-r.Intn(20))
	}
	if hex.Resource == FRUIT {
		// TODO: send request to take resource and if yes:
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
