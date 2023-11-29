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
	return &Board{
		Cases:           cases,
		XMax:            xmax,
		YMax:            ymax,
		HexSize:         hexSize,
		Biomes:          make([]*Biome, 0),
		ResourceManager: NewResourceManager(fruits, animals, rocks, woods),
		AgentManager:    NewAgentManager(cases, make(chan agentToManager, 100), agents, 0),
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

			biomeType := BiomesType(Randomizer.Intn(4))
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
				if try := Randomizer.Intn(100); try > 1 && neighbourHex != nil && neighbourHex.Biome == nil {
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
	for _, biome := range b.Biomes {
		resourceType := NONE
		hex := biome.Hexs[Randomizer.Intn(len(biome.Hexs))]
		switch biome.BiomeType {
		case PLAINS:
			if b.ResourceManager.MaxAnimalQuantity > b.ResourceManager.AnimalQuantity {
				resourceType = ANIMAL
			}
		case FOREST:
			if Randomizer.Intn(2) == 0 && b.ResourceManager.MaxFruitQuantity > b.ResourceManager.FruitQuantity {
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

func (b *Board) isValidHex(hex *Hexagone) bool {
	if hex == nil {
		return false
	}

	if hex.Position.X < 0 || hex.Position.X >= b.XMax || hex.Position.Y < 0 || hex.Position.Y >= b.YMax {
		return false
	}

	return true
}
