package typing

type Human struct {
	id          int
	Position    Hexagone
	Type        rune // can be cromagnon or neandertal
	Hungriness  int  // 0 to 100
	Thirstiness int  // 0 to 100
	Age         int
	Gender      rune
	Strength    int // 0 to 100
	Sociability int // 0 to 100
}
