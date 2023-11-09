package gui

import (
	_map "github.com/adrsimon/gomagnon/hexmap"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func NewSimulation() Simulation {
	simu := Simulation{}
	simu.gameMap = _map.NewGame(
		ScreenWidth,
		ScreenHeight,
		colornames.Black,
		colornames.Yellow,
		colornames.Darkgreen,
		colornames.Blue,
		colornames.Darkslategrey,
		30,
		30,
		40,
	)
	simu.gameMap.Board.Generate()
	simu.gameMap.Board.GenerateBiomes()

	simu.ScreenWidth = ScreenWidth
	simu.ScreenHeight = ScreenHeight

	return simu
}

func (s *Simulation) Update() error {
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(s.gameMap.BackgroundColor)

	for _, biome := range s.gameMap.Board.Biomes {
		color := s.gameMap.BackgroundColor
		switch biome.BiomeType {
		case _map.DIRT:
			color = s.gameMap.DirtColor
		case _map.FOREST:
			color = s.gameMap.ForestColor
		case _map.WATER:
			color = s.gameMap.WaterColor
		case _map.CAVE:
			color = s.gameMap.CaveColor
		}

		for _, hex := range biome.Hexs {
			hexSize := float32(s.gameMap.Board.HexSize)
			x := float32(hex.Position.X)
			y := float32(hex.Position.Y)

			var offsetX, offsetY float32
			offsetY = 0.75 * hexSize
			offsetX = 0

			if int(y)%2 == 0 {
				offsetX = hexSize / 2
				s.gameMap.DrawHex(screen, x*hexSize+offsetX, y*offsetY, color, hexSize)
			} else {
				s.gameMap.DrawHex(screen, x*hexSize+offsetX, y*offsetY, color, hexSize)
			}
		}
	}

	/** ORIGINAL VERSION OF GENERATION FUNCTION -- TO KEEP IF BIOMES ARE DELETED
	for i := 0; i < s.gameMap.Board.XMax; i++ {
		for j := 0; j < s.gameMap.Board.YMax; j++ {
			hex := s.gameMap.Board.Cases[fmt.Sprintf("%d:%d", i, j)]
			hexSize := float32(s.gameMap.Board.HexSize)
			x := float32(hex.Position.X)
			y := float32(hex.Position.Y)

			var offsetX, offsetY float32
			offsetY = 0.75 * hexSize
			offsetX = 0

			if j%2 == 0 {
				offsetX = hexSize/2
				s.gameMap.DrawHex(screen, x*hexSize+offsetX, y*offsetY+100, colornames.White, hexSize)
			} else {
				s.gameMap.DrawHex(screen, x*hexSize+offsetX, y*offsetY+100, colornames.White, hexSize)
			}

			/** TEXT DEBUG TO DISPLAY HEXAGONE POSITION -- TO REMOVE LATER
			middleX := x*hexSize + offsetX
			middleY := y*offsetY

			textX := middleX
			textY := middleY

			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d:%d", i, j), int(textX), int(textY))
			 *\/
		}
	}
	*/
}

func (s *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
