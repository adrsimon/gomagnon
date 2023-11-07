package gui

import _map "github.com/adrsimon/gomagnon/hexmap"

const (
	ScreenWidth  = 1080
	ScreenHeight = 720
)

type Simulation struct {
	gameMap *_map.GameMap

	ScreenWidth  int
	ScreenHeight int
}
