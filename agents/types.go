package agents

import "github.com/adrsimon/gomagnon/hexmap"

type Human struct {
	id          int
	Position    hexmap.Hexagone
	Type        rune // can me cromagnon or neandertal
	Hungriness  int  // 0 to 100
	Thirstiness int  // 0 to 100
	Age         int
	Gender      rune
	Strength    int // 0 to 100
	Sociability int // 0 to 100
}
