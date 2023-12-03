package simulation

import (
	"fmt"
	"image/color"

	"github.com/adrsimon/gomagnon/core/typing"
)

const (
	ScreenWidth  = 1080
	ScreenHeight = 720
)

type Simulation struct {
	Board *typing.Board

	ScreenWidth     int
	ScreenHeight    int
	backgroundColor color.RGBA

	cameraX, cameraY float32
	zoomFactor       float32

	debug bool
}

func NewSimulation() Simulation {
	simu := Simulation{}
	ressourcesMap := map[typing.ResourceType]int{
		typing.FRUIT:  20,
		typing.ANIMAL: 20,
		typing.ROCK:   20,
		typing.WOOD:   20,
	}

	simu.Board = typing.NewBoard(45, 40, 40, ressourcesMap)

	simu.cameraX = 0
	simu.cameraY = 0
	simu.zoomFactor = 0.6

	simu.Board.Generate()
	simu.Board.GenerateBiomes()
	simu.Board.GenerateResources()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	simu.Board.AgentManager.Start()

	for i := 0; i < 3; i++ {
		x, y := -1, -1
		for x == -1 && y == -1 {
			x = typing.Randomizer.Intn(simu.Board.XMax)
			y = typing.Randomizer.Intn(simu.Board.YMax)
			if simu.Board.Cases[x][y].Biome == typing.WATER {
				x, y = -1, -1
			}
		}
		simu.Board.AgentManager.Agents[fmt.Sprintf("ag-%d", simu.Board.AgentManager.Count)] = &typing.Human{
			ID:   fmt.Sprintf("ag-%d", i),
			Race: typing.Race(typing.Randomizer.Intn(2)),
			Body: typing.HumanBody{
				Thirstiness: 50,
				Hungriness:  50,
				Age:         float64(25),
			},
			Stats: typing.HumanStats{
				Strength:    10,
				Sociability: 10,
				Acuity:      typing.Randomizer.Intn(2) + 4,
			},
			Position:       simu.Board.Cases[x][y],
			Target:         nil,
			MovingToTarget: false,
			CurrentPath:    nil,
			Hut:            nil,
			Board:          simu.Board,
			Inventory:      typing.Inventory{Weight: 0, Object: make(map[typing.ResourceType]int)},
			AgentRelation:  make(map[string]string),
			AgentCommIn:    make(chan typing.AgentComm),
			Clan:           nil,
			Procreate:      typing.Procreate{Partner: nil, Timer: 100, Potential: true},
		}
		simu.Board.AgentManager.Count++
	}

	simu.debug = false

	return simu
}
