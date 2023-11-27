package typing

type Board struct {
	Cases           [][]*Hexagone
	XMax            int
	YMax            int
	HexSize         float32
	Biomes          []*Biome
	ResourceManager *ResourceManager
	AgentManager    *AgentManager
}

func NewBoard(xmax, ymax int, hexSize float32, fruits, animals, rocks, woods int) *Board {
	cases := make([][]*Hexagone, 0)
	for x := 0; x < xmax; x++ {
		cases = append(cases, make([]*Hexagone, ymax))
	}
	agents := make(map[string]*Human)
	resMan := NewResourceManager(fruits, animals, rocks, woods)
	return &Board{
		Cases:           cases,
		XMax:            xmax,
		YMax:            ymax,
		HexSize:         hexSize,
		Biomes:          make([]*Biome, 0),
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
	availableHexs := make([][]*Hexagone, b.XMax)
	for i := range availableHexs {
		availableHexs[i] = make([]*Hexagone, b.YMax)
		for j := range availableHexs[i] {
			availableHexs[i][j] = b.Cases[i][j]
		}
	}

	for i := range availableHexs {
		for j := range availableHexs[i] {
			hex := availableHexs[i][j]
			if hex == nil {
				continue
			}

			biomeType := BiomesType(r.Intn(4))
			biome := Biome{
				BiomeType: biomeType,
				Hexs:      make([]*Hexagone, 0),
			}
			biome.Hexs = append(biome.Hexs, hex)
			hex.Biome = &biome
			availableHexs[i][j] = nil

			neighbours := b.GetNeighbours(hex)
			for _, neighbour := range neighbours {
				if neighbour == nil {
					continue
				}
				neighbourHex := availableHexs[neighbour.Position.X][neighbour.Position.Y]
				if try := r.Intn(100); try > 1 && neighbourHex != nil && neighbourHex.Biome == nil {
					biome.Hexs = append(biome.Hexs, neighbour)
					neighbour.Biome = &biome
					availableHexs[neighbour.Position.X][neighbour.Position.Y] = nil
					neighbours = append(neighbours, b.GetNeighbours(neighbour)...)
				}
			}
			b.Biomes = append(b.Biomes, &biome)
		}
	}
}

func (b *Board) GenerateResources() {
	for b.ResourceManager.FruitQuantity < b.ResourceManager.MaxFruitQuantity {
		hex := b.Cases[r.Intn(b.XMax)][r.Intn(b.YMax)]
		if hex.Biome.BiomeType == FOREST {
			hex.Resource = FRUIT
			b.ResourceManager.FruitQuantity++
		}
	}

	for b.ResourceManager.AnimalQuantity < b.ResourceManager.MaxAnimalQuantity {
		hex := b.Cases[r.Intn(b.XMax)][r.Intn(b.YMax)]
		if hex.Biome.BiomeType == PLAINS {
			hex.Resource = ANIMAL
			b.ResourceManager.AnimalQuantity++
		}
	}

	for b.ResourceManager.RockQuantity < b.ResourceManager.MaxRockQuantity {
		hex := b.Cases[r.Intn(b.XMax)][r.Intn(b.YMax)]
		if hex.Biome.BiomeType == CAVE {
			hex.Resource = ROCK
			b.ResourceManager.RockQuantity++
		}
	}

	for b.ResourceManager.WoodQuantity < b.ResourceManager.MaxWoodQuantity {
		hex := b.Cases[r.Intn(b.XMax)][r.Intn(b.YMax)]
		if hex.Biome.BiomeType == FOREST {
			hex.Resource = WOOD
			b.ResourceManager.WoodQuantity++
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
