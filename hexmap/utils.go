package hexmap

import (
	"fmt"
	"image/color"
	"math/rand"
)

func NewGame(
	screenWidth int,
	screenHeight int,
	backgroundColor color.RGBA,
	dirtColor color.RGBA,
	forestColor color.RGBA,
	waterColor color.RGBA,
	caveColor color.RGBA,
	xmax int,
	ymax int,
	hexSize int,
) *GameMap {
	return &GameMap{
		Board:           NewBoard(xmax, ymax, hexSize),
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		BackgroundColor: backgroundColor,
		DirtColor:       dirtColor,
		ForestColor:     forestColor,
		WaterColor:      waterColor,
		CaveColor:       caveColor,
	}
}

func NewBoard(xmax int, ymax int, hexSize int) *Board {
	return &Board{
		Cases:   make(map[string]*Hexagone),
		XMax:    xmax,
		YMax:    ymax,
		HexSize: hexSize,
		Biomes:  make([]*Biome, 0),
		Agents:  make([]*Human, 0),
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
			if try := rand.Intn(100); try > 80 && ok {
				biomeHexs.Hexs = append(biomeHexs.Hexs, neighbour)
				delete(availableHexs, key)
				neighbours = append(neighbours, b.GetNeighbours(neighbour)...)
			}
		}
		b.Biomes = append(b.Biomes, &biomeHexs)
	}
}

func (b *Board) GenerateHumans() {
	humans := make([]*Human, 10)

	availableHexs := make(map[string]*Hexagone)
	for k, v := range b.Cases {
		availableHexs[k] = v
	}

	for i := range humans {
		// Ensure unique position for each human
		for pos, hex := range availableHexs {
			humans[i] = &Human{
				id:          i,
				Position:    *hex,
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
}
