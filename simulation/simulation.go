package simulation

import (
	"fmt"
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
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

	Debug  bool
	Paused bool
	TPS    int

	SelectedAgent string
	SavedLen      int
	Selector      *widget.List
	AgentDesc     *widget.TextArea

	UI *ebitenui.UI
}

func NewSimulation() Simulation {
	simu := Simulation{}
	ressourcesMap := map[typing.ResourceType]int{
		typing.FRUIT:   30,
		typing.ANIMAL:  30,
		typing.ROCK:    30,
		typing.WOOD:    30,
		typing.MAMMOTH: 4,
	}

	simu.Board = typing.NewBoard(46, 41, 40, ressourcesMap)

	simu.cameraX = 0
	simu.cameraY = 0
	simu.zoomFactor = 0.6

	simu.Board.Generate()
	simu.Board.GenerateBiomes()
	simu.Board.GenerateResources()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	simu.Board.AgentManager.Start()

	for i := 0; i < 20; i++ {
		x, y := -1, -1
		for x == -1 && y == -1 {
			x = typing.Randomizer.Intn(simu.Board.XMax)
			y = typing.Randomizer.Intn(simu.Board.YMax)
			if simu.Board.Cases[x][y].Biome == typing.WATER {
				x, y = -1, -1
			}
		}

		ag := &typing.Agent{
			ID:   fmt.Sprintf("ag-%d", i),
			Type: []rune{'M', 'F'}[typing.Randomizer.Intn(2)],
			Race: typing.Race(typing.Randomizer.Intn(2)),
			Body: typing.HumanBody{
				Thirstiness: 50,
				Hungriness:  50,
				Age:         float64(25),
			},
			Stats: typing.HumanStats{
				Strength:    50,
				Sociability: 50,
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
			Procreate:      typing.Procreate{Partner: nil, Timer: 100},
		}
		simu.Board.AgentManager.Agents = append(simu.Board.AgentManager.Agents, ag)
		ag.Behavior = &typing.HumanBehavior{H: ag}
		simu.Board.AgentManager.Count++
	}

	simu.Debug = false
	simu.Paused = false
	simu.SelectedAgent = ""
	simu.TPS = 20

	return simu
}
