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

type ResourceManager struct {
	maxQuantities     map[ResourceType]int
	currentQuantities map[ResourceType]int
}

func NewResourceManager(maxs map[ResourceType]int) *ResourceManager {
	return &ResourceManager{
		maxQuantities:     maxs,
		currentQuantities: make(map[ResourceType]int),
	}
}
