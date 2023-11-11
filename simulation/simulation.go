package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"golang.org/x/image/colornames"
)

const (
	ScreenWidth  = 1080
	ScreenHeight = 720
)

type Simulation struct {
	GameMap *typing.GameMap

	ScreenWidth  int
	ScreenHeight int
}

func NewSimulation() Simulation {
	simu := Simulation{}
	simu.GameMap = typing.NewGame(
		ScreenWidth, ScreenHeight,
		colornames.Black,
		27, 23,
		40,
		10, 10, 10, 10,
	)
	simu.GameMap.Board.Generate()
	simu.GameMap.Board.GenerateBiomes()
	simu.GameMap.Board.GenerateResources()
	simu.GameMap.Board.GenerateHumans()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	return simu
}
