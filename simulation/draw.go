package simulation

import (
	"github.com/adrsimon/gomagnon/core/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(s.GameMap.BackgroundColor)
	drawing.DrawBoard(screen, s.GameMap.Board)
	drawing.DrawAgents(screen, s.GameMap.Board)
}
