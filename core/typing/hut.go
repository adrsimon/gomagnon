package typing

var Needs = map[string]map[ResourceType]int{
	"hut": map[ResourceType]int{
		WOOD: 3,
		ROCK: 3,
	},
}

type Hut struct {
	Position  *Hexagone
	Inventory map[ResourceType]int
	Owner     *Human
}
