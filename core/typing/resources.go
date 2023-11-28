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

type CoolDown struct {
	Current  int
	Resource ResourceType
}

type ResourceManager struct {
	maxQuantities     map[ResourceType]int
	CurrentQuantities map[ResourceType]int
	RespawnCDs        []CoolDown
}

func NewResourceManager(maxs map[ResourceType]int) *ResourceManager {
	return &ResourceManager{
		maxQuantities:     maxs,
		CurrentQuantities: make(map[ResourceType]int),
		RespawnCDs:        make([]CoolDown, 0),
	}
}
