package typing

import (
	"fmt"
	"math/rand"
)

type Human struct {
	id             string
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
	ComOut         agentToManager
	Comin          managerToAgent
	Map            map[string]*Hexagone
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
)

func (h *Human) EvaluateOneHex(hex *Hexagone) float64 {
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

func (h *Human) GetNeighbours(coord string) []*Hexagone {
	neighbours := make([]*Hexagone, 0)
	pos := h.Map[coord].Position
	if pos.Y%2 == 0 {
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X+1, pos.Y+1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X, pos.Y-1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X+1, pos.Y-1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X-1, pos.Y)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X+1, pos.Y)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X, pos.Y+1)])
	} else {
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X-1, pos.Y)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X, pos.Y-1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X+1, pos.Y)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X-1, pos.Y+1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X-1, pos.Y-1)])
		neighbours = append(neighbours, h.Map[fmt.Sprintf("%d:%d", pos.X, pos.Y+1)]) // SI JAMAIS EXISTE PAS ????
	}
	return neighbours
}

func (h *Human) GetNeighborsWithin5() []*Hexagone {
	currentLevel := []*Hexagone{h.Map[h.Position]}
	visited := make(map[string]bool)

	for i := 0; i < 5; i++ {
		nextLevel := make([]*Hexagone, 0)

		for _, currentHex := range currentLevel {
			var neighbors []*Hexagone
			neighbors = h.GetNeighbours(currentHex)

			for _, neighbor := range neighbors {
				if !visited[fmt.Sprintf("%d:%d", neighbor.Position.X, neighbor.Position.Y)] {
					visited[fmt.Sprintf("%d:%d", neighbor.Position.X, neighbor.Position.Y)] = true
					nextLevel = append(nextLevel, neighbor)
				}
			}
		}

		// Concaténation des niveaux pour passer au niveau suivant
		currentLevel = append(currentLevel, nextLevel...)
	}

	// Retourne tous les voisins à 5 cases de distance ou moins
	return currentLevel
}

func (h *Human) BestNeighbor(surroundingHexagons []*Hexagone) *Hexagone {
	best := 0
	for i, v := range surroundingHexagons {
		score := h.EvaluateOneHex(v)
		if score > best {
			best = i
		}
	}
	return surroundingHexagons[best]
}

// Adapt to quentin A* + add concurency handling + send request to take resource + add updating resources when taken
func (h *Human) UpdateAgent() {
	if !h.MovingToTarget {
		surroundingHexagons := h.GetNeighborsWithin5()
		targetHexagon := h.BestNeighbor(surroundingHexagons)

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

func (h *Human) Start() {

}
