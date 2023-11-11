package typing

type ResourceType int

const (
	NONE ResourceType = iota
	FRUIT
	ANIMAL
	ROCK
	WOOD
)

type ResourceManager struct {
	Resources         []ResourceType
	MaxFruitQuantity  int
	MaxAnimalQuantity int
	MaxRockQuantity   int
	MaxWoodQuantity   int
	FruitQuantity     int
	AnimalQuantity    int
	RockQuantity      int
	WoodQuantity      int
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
