package typing

import (
	"github.com/aquilax/go-perlin"
)

type Board struct {
	Cases           [][]*Hexagone
	XMax            int
	YMax            int
	HexSize         float32
	ResourceManager *ResourceManager
	AgentManager    *AgentManager
}

func NewBoard(xmax, ymax int, hexSize float32, maxResources map[ResourceType]int) *Board {
	cases := make([][]*Hexagone, 0)
	for x := 0; x < xmax; x++ {
		cases = append(cases, make([]*Hexagone, ymax))
	}
	agents := make([]*Agent, 0)
	resMan := NewResourceManager(maxResources)
	return &Board{
		Cases:           cases,
		XMax:            xmax,
		YMax:            ymax,
		HexSize:         hexSize,
		ResourceManager: resMan,
		AgentManager:    NewAgentManager(cases, make(chan agentToManager, 100), agents, resMan),
	}
}

func (b *Board) Generate() {
	for x := 0; x < b.XMax; x++ {
		for y := 0; y < b.YMax; y++ {
			b.Cases[x][y] = &Hexagone{
				Position: &Point2D{
					X: x,
					Y: y,
				},
			}
		}
	}
}

func (b *Board) GetNeighbours(hex *Hexagone) []*Hexagone {
	neighbours := make([]*Hexagone, 0)

	addIfExist := func(x, y int) {
		if x >= 0 && x < b.XMax && y >= 0 && y < b.YMax {
			neighbours = append(neighbours, b.Cases[x][y])
		}
	}

	if hex.Position.Y%2 == 0 {
		addIfExist(hex.Position.X+1, hex.Position.Y+1)
		addIfExist(hex.Position.X, hex.Position.Y-1)
		addIfExist(hex.Position.X+1, hex.Position.Y-1)
		addIfExist(hex.Position.X-1, hex.Position.Y)
		addIfExist(hex.Position.X+1, hex.Position.Y)
		addIfExist(hex.Position.X, hex.Position.Y+1)
	} else {
		addIfExist(hex.Position.X-1, hex.Position.Y)
		addIfExist(hex.Position.X, hex.Position.Y-1)
		addIfExist(hex.Position.X+1, hex.Position.Y)
		addIfExist(hex.Position.X-1, hex.Position.Y+1)
		addIfExist(hex.Position.X-1, hex.Position.Y-1)
		addIfExist(hex.Position.X, hex.Position.Y+1)
	}
	return neighbours
}

func (b *Board) GenerateBiomes() {
	elevation := perlin.NewPerlin(1, 2.7, 3, Seed)
	moisture := perlin.NewPerlin(0.8, 2, 5, Seed)

	for i, line := range b.Cases {
		for j := range line {
			hex := b.Cases[i][j]
			if hex == nil {
				continue
			}

			var biomeType BiomeType

			x := float64(i) / float64(b.XMax)
			y := float64(j) / float64(b.YMax)

			elevationValue := elevation.Noise2D(x, y)
			moistureValue := moisture.Noise2D(x, y)

			switch {
			case elevationValue > 0.3:
				biomeType = CAVE
			case elevationValue < -0.4:
				biomeType = WATER
			default:
				if moistureValue > 0 {
					biomeType = FOREST
				} else {
					biomeType = PLAINS
				}
			}

			hex.Biome = biomeType
		}
	}
}

func (b *Board) GenerateResources() {
	for i := 0; i < int(NUM_RESOURCE_TYPES); i++ {
		res := ResourceType(i)
		for b.ResourceManager.CurrentQuantities[res] < b.ResourceManager.maxQuantities[res] {
			hex := b.Cases[Randomizer.Intn(b.XMax)][Randomizer.Intn(b.YMax)]
			if hex.Resource != NONE {
				continue
			}
			if (res == ANIMAL && hex.Biome != PLAINS) || b.CountResourcesAround(hex, ANIMAL, FRUIT, 5) > 2 {
				continue
			} else if res == FRUIT && hex.Biome != FOREST || b.CountResourcesAround(hex, ANIMAL, FRUIT, 5) > 2 {
				continue
			} else if res == WOOD && hex.Biome != FOREST {
				continue
			} else if res == ROCK && hex.Biome != CAVE {
				continue
			} else if res == MAMMOTH && hex.Biome != PLAINS {
				continue
			}
			hex.Resource = res
			b.ResourceManager.CurrentQuantities[res]++
		}
	}
}

func (b *Board) CountResourcesAround(hex *Hexagone, resType1, resType2 ResourceType, acuity int) int {
	neighbours := b.GetNeighbours(hex)
	visited := make(map[*Hexagone]bool)
	count := 0

	for i := 0; i < acuity; i++ {
		newNeighbours := []*Hexagone{}
		for _, neighbour := range neighbours {
			if neighbour == nil || visited[neighbour] {
				continue
			}
			visited[neighbour] = true

			if neighbour.Resource == resType1 || neighbour.Resource == resType2 {
				count++
			}

			newNeighbours = append(newNeighbours, b.GetNeighbours(neighbour)...)
		}
		neighbours = newNeighbours
	}
	return count
}

func (b *Board) isValidHex(hex *Hexagone) bool {
	if hex == nil {
		return false
	}

	if hex.Position.X < 0 || hex.Position.X >= b.XMax || hex.Position.Y < 0 || hex.Position.Y >= b.YMax {
		return false
	}

	return true
}
