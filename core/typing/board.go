package typing

import (
	"github.com/aquilax/go-perlin"
	"math"
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

func (b *Board) GenerateContinentBiomes() {
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
			case elevationValue < -0.7:
				biomeType = DEEP_WATER
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
			b.ResourceManager.FreeSpots[biomeType] = append(b.ResourceManager.FreeSpots[biomeType], Point2D{X: i, Y: j})
		}
	}
}

func (b *Board) GenerateIslandBiomes() {
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

			lerp := func(a, b, t float64) float64 {
				return a + t*(b-a)
			}

			d := func(x, y float64) float64 {
				return math.Sqrt(math.Pow(x-0.5, 2) + math.Pow(y-0.5, 2))
			}

			elevationValue = math.Abs(lerp(elevationValue, d(x, y), 0.75))

			switch {
			case elevationValue > 0.4:
				biomeType = DEEP_WATER
			case elevationValue > 0.3:
				biomeType = WATER
			case elevationValue < 0.05:
				biomeType = CAVE
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
			biome := ResourceToBiome[res][Randomizer.Intn(len(ResourceToBiome[res]))]
			if len(b.ResourceManager.FreeSpots[biome]) == 0 {
				break
			}
			spot := Randomizer.Intn(len(b.ResourceManager.FreeSpots[biome]))
			hexPos := b.ResourceManager.FreeSpots[biome][spot]
			hex := b.Cases[hexPos.X][hexPos.Y]
			hex.Resource = res
			b.ResourceManager.CurrentQuantities[res]++
			b.ResourceManager.FreeSpots[biome] = append(b.ResourceManager.FreeSpots[biome][:spot], b.ResourceManager.FreeSpots[biome][spot+1:]...)
		}
	}
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
