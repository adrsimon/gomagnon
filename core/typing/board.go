package typing

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Cases           map[string]*Hexagone
	XMax            int
	YMax            int
	HexSize         int
	Biomes          []*Biome
	ResourceManager *ResourceManager
	AgentManager    *AgentManager
	Agents          []*Human
}

func NewBoard(xmax, ymax, hexSize, fruits, animals, rocks, woods int) *Board {
	return &Board{
		Cases:           make(map[string]*Hexagone),
		XMax:            xmax,
		YMax:            ymax,
		HexSize:         hexSize,
		Biomes:          make([]*Biome, 0),
		ResourceManager: NewResourceManager(fruits, animals, rocks, woods),
	}
}

func (b *Board) Generate() {
	for x := 0; x < b.XMax; x++ {
		for y := 0; y < b.YMax; y++ {
			b.Cases[fmt.Sprintf("%d:%d", x, y)] = &Hexagone{
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
	if hex.Position.Y%2 == 0 {
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X+1, hex.Position.Y+1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X, hex.Position.Y-1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X+1, hex.Position.Y-1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X-1, hex.Position.Y)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X+1, hex.Position.Y)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X, hex.Position.Y+1)])
	} else {
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X-1, hex.Position.Y)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X, hex.Position.Y-1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X+1, hex.Position.Y)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X-1, hex.Position.Y+1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X-1, hex.Position.Y-1)])
		neighbours = append(neighbours, b.Cases[fmt.Sprintf("%d:%d", hex.Position.X, hex.Position.Y+1)])
	}
	return neighbours
}

func (b *Board) GenerateBiomes() {
	availableHexs := make(map[string]*Hexagone)
	for k, v := range b.Cases {
		availableHexs[k] = v
	}

	for pos, hex := range availableHexs {
		if hex == nil {
			continue
		}
		biomeType := BiomesType(rand.Intn(4))
		biome := Biome{
			BiomeType: biomeType,
			Hexs:      make([]*Hexagone, 0),
		}
		biome.Hexs = append(biome.Hexs, hex)
		hex.Biome = &biome
		delete(availableHexs, pos)

		neighbours := b.GetNeighbours(hex)
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			key := fmt.Sprintf("%d:%d", neighbour.Position.X, neighbour.Position.Y)
			_, ok := availableHexs[key]
			if try := rand.Intn(100); try > 1 && ok {
				biome.Hexs = append(biome.Hexs, neighbour)
				neighbour.Biome = &biome
				delete(availableHexs, key)
				neighbours = append(neighbours, b.GetNeighbours(neighbour)...)
			}
		}
		b.Biomes = append(b.Biomes, &biome)
	}
}

func (b *Board) GenerateResources() {
	for _, biome := range b.Biomes {
		resourceType := NONE

		hex := biome.Hexs[rand.Intn(len(biome.Hexs))]
		switch biome.BiomeType {
		case PLAINS:
			if b.ResourceManager.MaxAnimalQuantity > b.ResourceManager.AnimalQuantity {
				resourceType = ANIMAL
			}
		case FOREST:
			if rand.Intn(2) == 0 && b.ResourceManager.MaxFruitQuantity > b.ResourceManager.FruitQuantity {
				resourceType = FRUIT
			} else if b.ResourceManager.MaxWoodQuantity > b.ResourceManager.WoodQuantity {
				resourceType = WOOD
			}
		case CAVE:
			if b.ResourceManager.MaxRockQuantity > b.ResourceManager.RockQuantity {
				resourceType = ROCK
			}
		}
		hex.Resource = resourceType
		b.ResourceManager.Resources = append(b.ResourceManager.Resources, resourceType)
		switch resourceType {
		case FRUIT:
			b.ResourceManager.FruitQuantity++
		case ANIMAL:
			b.ResourceManager.AnimalQuantity++
		case ROCK:
			b.ResourceManager.RockQuantity++
		case WOOD:
			b.ResourceManager.WoodQuantity++
		}
	}
}

func (b *Board) GenerateHumans() {
	humans := make([]*Human, 10)

	availableHexs := make(map[string]*Hexagone)
	for k, v := range b.Cases {
		availableHexs[k] = v
	}

	for i := range humans {
		for pos := range availableHexs {
			humans[i] = &Human{
				id:          i,
				Position:    pos,
				Type:        rune(rand.Intn(2)), // 0 or 1
				Hungriness:  rand.Intn(101),     // 0 to 100
				Thirstiness: rand.Intn(101),     // 0 to 100
				Age:         rand.Intn(101),     // 0 to 100
				Gender:      rune(rand.Intn(2)), // 0 or 1
				Strength:    rand.Intn(101),     // 0 to 100
				Sociability: rand.Intn(101),     // 0 to 100
			}
			delete(availableHexs, pos)
			break
		}
	}

	b.Agents = humans
	//print agents
	// for _, agent := range b.Agents {
	// 	fmt.Println(*agent)
	// }

}
