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

	debug bool
}

func NewSimulation() Simulation {
	simu := Simulation{}
	simu.GameMap = typing.NewGame(
		ScreenWidth, ScreenHeight,
		colornames.Black,
		90, 80,
		40,
		map[typing.ResourceType]int{
			typing.FRUIT:  100,
			typing.ANIMAL: 100,
			typing.ROCK:   100,
			typing.WOOD:   100,
		},
	)

	simu.cameraX = 0
	simu.cameraY = 0
	simu.zoomFactor = 0.3

	simu.GameMap.Board.Generate()
	simu.GameMap.Board.GenerateBiomes()
	simu.GameMap.Board.GenerateResources()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	simu.GameMap.Board.AgentManager.Start()

	for i := 0; i < 100; i++ {
		x, y := -1, -1
		for x == -1 && y == -1 {
			x = typing.Randomizer.Intn(simu.GameMap.Board.XMax)
			y = typing.Randomizer.Intn(simu.GameMap.Board.YMax)
			if simu.GameMap.Board.Cases[x][y].Biome.BiomeType == typing.WATER {
				x, y = -1, -1
			}
		}

		simu.GameMap.Board.AgentManager.Agents[fmt.Sprintf("ag-%d", i)] = &typing.Human{
			ID:   fmt.Sprintf("ag-%d", i),
			Race: typing.Race(typing.Randomizer.Intn(2)),
			Body: typing.HumanBody{
				Thirstiness: 50,
				Hungriness:  50,
			},
			Stats: typing.HumanStats{
				Strength:    10,
				Sociability: 10,
				Acuity:      typing.Randomizer.Intn(2) + 4,
			},
			Position:       simu.GameMap.Board.Cases[x][y],
			Target:         nil,
			MovingToTarget: false,
			CurrentPath:    nil,
			Hut:            nil,
			Board:          simu.GameMap.Board,
			Inventory:      typing.Inventory{Weight: 0, Object: make(map[typing.ResourceType]int)},
			AgentRelation:  make(map[string]string),
			AgentCommIn:    make(chan typing.AgentComm),
			Clan:           nil,
		}
	}

	simu.debug = false

	return simu
}
