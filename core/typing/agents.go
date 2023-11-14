package typing

import (
	"fmt"
	"math/rand"
)

type Human struct {
	id             int
	Position       string
	Type           rune // can be cromagnon or neandertal
	Hungriness     int  // 0 to 100
	Thirstiness    int  // 0 to 100
	Age            int
	Gender         rune
	Strength       int // 0 to 100
	Sociability    int // 0 to 100
	CurrentPath    []*Hexagone
	MovingToTarget bool
	Target         *Hexagone
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
)

func (h *Human) EvaluateSurroundings(hex *Hexagone) float64 {
	var score = 0.0

	switch hex.Resource {
	case ANIMAL:
		score = (float64(h.Hungriness) / 100) * AnimalFoodValueMultiplier
	case FRUIT:
		score = (float64(h.Hungriness) / 100) * FruitFoodValueMultiplier
	}

	if hex.Biome.BiomeType == WATER {
		score = (float64(h.Thirstiness) / 100) * WaterValueMultiplier
	}

	return score
}

// Adapt to quentin A* + add concurency handling + send request to take resource + add updating resources when taken
func (h *Human) UpdateAgent() {
	if !h.MovingToTarget {
		surroundingHexagons := h.GetSurroundingHexagons()
		targetHexagon := h.ChooseTargetHexagon(surroundingHexagons)

		if targetHexagon != nil {
			h.CurrentPath = AStarPathfinding(h.CurrentHexagon(), targetHexagon)
			h.Target = targetHexagon
			h.MovingToTarget = true
		}
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
}

func (h *Human) UpdateStateBasedOnResource(hex *Hexagone) {
	if hex.Resource == ANIMAL {
		// TODO: send request to take resource and if yes:
		h.Hungriness = max(0, h.Hungriness-rand.Intn(20))
	}
	if hex.Resource == FRUIT {
		// TODO: send request to take resource and if yes:
		h.Hungriness = max(0, h.Hungriness-rand.Intn(10))
	}
	if hex.Biome.BiomeType == WATER {
		h.Thirstiness = max(0, h.Thirstiness-rand.Intn(30))
	}
}

func (h *Human) MoveToHexagon(hex *Hexagone) {
	h.Position = fmt.Sprintf("%d:%d", hex.Position.X, hex.Position.Y)
}
