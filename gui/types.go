package gui

import _map "github.com/adrsimon/gomagnon/hexmap"

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

type Simulation struct {
	gameMap *_map.GameMap

	ScreenWidth  int
	ScreenHeight int
}
