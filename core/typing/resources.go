package typing

type ResourceType int

const (
	NONE ResourceType = iota
	FRUIT
	ANIMAL
	ROCK
	WOOD
	NUM_RESOURCE_TYPES
)

const (
	MaxWeightInv = 10.0
	WeightFruit  = 0.1
	WeightAnimal = 0.5
	WeightRock   = 2.0
	WeightWood   = 1.0
)

var Weights = map[ResourceType]float64{
	FRUIT:  WeightFruit,
	ANIMAL: WeightAnimal,
	ROCK:   WeightRock,
	WOOD:   WeightWood,
}

var ResourceToBiome = map[ResourceType][]BiomeType{
	FRUIT:  {FOREST},
	ANIMAL: {PLAINS},
	ROCK:   {CAVE},
	WOOD:   {FOREST},
}

type CoolDown struct {
	Current  int
	Resource ResourceType
}

type ResourceManager struct {
	maxQuantities     map[ResourceType]int
	CurrentQuantities map[ResourceType]int
	RespawnCDs        []CoolDown
	FreeSpots         map[BiomeType][]Point2D
}

func NewResourceManager(maxs map[ResourceType]int) *ResourceManager {
	return &ResourceManager{
		maxQuantities:     maxs,
		CurrentQuantities: make(map[ResourceType]int),
		RespawnCDs:        make([]CoolDown, 0),
		FreeSpots:         make(map[BiomeType][]Point2D),
	}
}
