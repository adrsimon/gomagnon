package hexmap

import (
	"fmt"
	"image/color"
	"math/rand"
)

func NewGame(
	screenWidth, screenHeight int,
	backgroundColor, dirtColor, forestColor, waterColor, caveColor color.RGBA,
	xmax, ymax int,
	hexSize int,
	fruits, animals, rocks, woods int,
) *GameMap {
	return &GameMap{
		Board:           NewBoard(xmax, ymax, hexSize, fruits, animals, rocks, woods),
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		BackgroundColor: backgroundColor,
		DirtColor:       dirtColor,
		ForestColor:     forestColor,
		WaterColor:      waterColor,
		CaveColor:       caveColor,
	}
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

func NewResourceManager(fruits, animals, rocks, woods int) *ResourceManager {
	return &ResourceManager{
		Resources:         make([]ResourceType, 0),
		MaxFruitQuantity:  fruits,
		MaxAnimalQuantity: animals,
		MaxRockQuantity:   rocks,
		MaxWoodQuantity:   woods,
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
		biome := BiomesType(rand.Intn(4))
		biomeHexs := Biome{
			BiomeType: biome,
			Hexs:      make([]*Hexagone, 0),
		}
		biomeHexs.Hexs = append(biomeHexs.Hexs, hex)
		delete(availableHexs, pos)

		neighbours := b.GetNeighbours(hex)
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			key := fmt.Sprintf("%d:%d", neighbour.Position.X, neighbour.Position.Y)
			_, ok := availableHexs[key]
			if try := rand.Intn(100); try > 1 && ok {
				biomeHexs.Hexs = append(biomeHexs.Hexs, neighbour)
				delete(availableHexs, key)
				neighbours = append(neighbours, b.GetNeighbours(neighbour)...)
			}
		}
		b.Biomes = append(b.Biomes, &biomeHexs)
	}
}

func (b *Board) GetHexBiome(hex *Hexagone) *Biome {
	for _, biome := range b.Biomes {
		for _, biomeHex := range biome.Hexs {
			if biomeHex == hex {
				return biome
			}
		}
	}
	return nil
}

func (b *Board) GenerateResources() {
	for _, biome := range b.Biomes {
		var resourceType ResourceType

		hex := biome.Hexs[rand.Intn(len(biome.Hexs))]
		if (b.ResourceManager.MaxFruitQuantity > b.ResourceManager.FruitQuantity) ||
			(b.ResourceManager.MaxAnimalQuantity > b.ResourceManager.AnimalQuantity) ||
			(b.ResourceManager.MaxRockQuantity > b.ResourceManager.RockQuantity) ||
			(b.ResourceManager.MaxWoodQuantity > b.ResourceManager.WoodQuantity) {
			switch biome.BiomeType {
			case DIRT:
				resourceType = ANIMAL
			case FOREST:
				if rand.Intn(2) == 0 {
					resourceType = FRUIT
				} else {
					resourceType = WOOD
				}
			case WATER:
				resourceType = NONE
			case CAVE:
				resourceType = ROCK
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
		} else {
			hex.Resource = NONE
		}
	}
}
