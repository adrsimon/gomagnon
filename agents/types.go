package agents

import "github.com/adrsimon/gomagnon/hexmap"

type Human struct {
	id          int
	Position    hexmap.Hexagone
	Type        rune
	Hungriness  int
	Thirstiness int
	Age         int
	Gender      rune
	Strength    int
	Sociability int
}
