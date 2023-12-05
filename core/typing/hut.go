package typing

var Needs = map[string]map[ResourceType]int{
	"hut": map[ResourceType]int{
		WOOD: 1,
		ROCK: 1,
	},
}

type Hut struct {
	Position  *Hexagone
	Inventory map[ResourceType]int
	Owner     *Human
}
