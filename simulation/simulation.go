package simulation

import (
	"fmt"
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

	cameraX, cameraY float32
	zoomFactor       float32
}

func NewSimulation() Simulation {
	simu := Simulation{}
	simu.GameMap = typing.NewGame(
		ScreenWidth, ScreenHeight,
		colornames.Black,
		28, 25,
		40,
		10, 10, 10, 10,
	)

	simu.cameraX = 0
	simu.cameraY = 0
	simu.zoomFactor = 1

	simu.GameMap.Board.Generate()
	simu.GameMap.Board.GenerateBiomes()
	simu.GameMap.Board.GenerateResources()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	simu.GameMap.Board.AgentManager.Start()

	for i := 0; i < 50; i++ {
		simu.GameMap.Board.AgentManager.Agents[fmt.Sprintf("ag-%d", i)] = &typing.Human{
			Type:           0,
			Body:           typing.HumanBody{},
			Stats:          typing.HumanStats{},
			Position:       simu.GameMap.Board.Cases[1][1],
			Target:         nil,
			MovingToTarget: false,
			CurrentPath:    nil,
			Board:          simu.GameMap.Board,
			AgentRelation:  make(map[string]string),
		}
	}

	return simu
}
