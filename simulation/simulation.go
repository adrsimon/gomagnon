package simulation

import (
	"fmt"
	"image/color"

	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/adrsimon/gomagnon/settings"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
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

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	resourcesMap := map[typing.ResourceType]int{
		typing.FRUIT:  settings.Setting.World.Resources.MaxFruits,
		typing.ANIMAL: settings.Setting.World.Resources.MaxAnimals,
		typing.ROCK:   settings.Setting.World.Resources.MaxRocks,
		typing.WOOD:   settings.Setting.World.Resources.MaxWoods,
	}

	hexSize := float32(simu.ScreenWidth) / float32(settings.Setting.World.Size.X-1)
	simu.Board = typing.NewBoard(
		settings.Setting.World.Size.X,
		settings.Setting.World.Size.Y,
		hexSize,
		resourcesMap,
	)

	simu.cameraX = 0
	simu.cameraY = 0
	simu.zoomFactor = 1

	simu.Board.Generate()
	if settings.Setting.World.Type == "island" {
		simu.Board.GenerateIslandBiomes()
	} else if settings.Setting.World.Type == "continent" {
		simu.Board.GenerateContinentBiomes()
	}
	simu.Board.GenerateResources()

	simu.Board.AgentManager.Start()

	for i := 0; i < settings.Setting.Agents.InitialNumber; i++ {
		hex := simu.Board.Spawnables[typing.Randomizer.Intn(len(simu.Board.Spawnables))]

		mapVision := make([][]typing.Hexagone, simu.Board.XMax)
		for i := range mapVision {
			mapVision[i] = make([]typing.Hexagone, simu.Board.YMax)
			for j := range mapVision[i] {
				mapVision[i][j] = typing.Hexagone{Position: &typing.Point2D{X: -1, Y: -1}}
			}
		}

		if i < settings.Setting.Agents.InitialNumber/2 {
			ag := &typing.Agent{
				ID:   fmt.Sprintf("ag-%d", i),
				Type: []rune{'M', 'F'}[typing.Randomizer.Intn(2)],
				Race: typing.Race(0),
				Body: typing.HumanBody{
					Thirstiness: float64(int(typing.Randomizer.NormFloat64()*10 + 50)),
					Hungriness:  float64(int(typing.Randomizer.NormFloat64()*10 + 50)),
					Age:         float64(typing.Randomizer.Intn(10) + 10),
				},
				Stats: typing.HumanStats{
					Strength:    int(typing.Randomizer.NormFloat64()*10 + 60),
					Sociability: int(typing.Randomizer.NormFloat64()*10 + 40),
					Acuity:      typing.Randomizer.Intn(2) + 4,
				},
				MapVision:      mapVision,
				Position:       simu.Board.Cases[hex.X][hex.Y],
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
				Opponent:       nil,
				Fightcooldown:  50 + typing.Randomizer.Intn(200),
				DeathAge:       typing.Randomizer.NormFloat64()*6 + 37,
			}

			simu.Board.AgentManager.Agents = append(simu.Board.AgentManager.Agents, ag)
			ag.Behavior = &typing.HumanBehavior{H: ag}
			simu.Board.AgentManager.Count++
		} else {
			ag := &typing.Agent{
				ID:   fmt.Sprintf("ag-%d", i),
				Type: []rune{'M', 'F'}[typing.Randomizer.Intn(2)],
				Race: typing.Race(1),
				Body: typing.HumanBody{
					Thirstiness: float64(int(typing.Randomizer.NormFloat64()*10 + 50)),
					Hungriness:  float64(int(typing.Randomizer.NormFloat64()*10 + 50)),
					Age:         float64(typing.Randomizer.Intn(10) + 10),
				},
				Stats: typing.HumanStats{
					Strength:    int(typing.Randomizer.NormFloat64()*10 + 50),
					Sociability: int(typing.Randomizer.NormFloat64()*10 + 70),
					Acuity:      typing.Randomizer.Intn(3) + 4,
				},
				MapVision:      mapVision,
				Position:       simu.Board.Cases[hex.X][hex.Y],
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
				Opponent:       nil,
				Fightcooldown:  50 + typing.Randomizer.Intn(200),
				DeathAge:       typing.Randomizer.NormFloat64()*6 + 37,
			}

			simu.Board.AgentManager.Agents = append(simu.Board.AgentManager.Agents, ag)
			ag.Behavior = &typing.HumanBehavior{H: ag}
			simu.Board.AgentManager.Count++
		}
	}

	simu.Debug = false
	simu.Paused = false
	simu.SelectedAgent = ""
	simu.TPS = 20

	return simu
}
