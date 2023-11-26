package simulation

import (
	"github.com/adrsimon/gomagnon/core/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Simulation) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.cameraX += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.cameraX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.cameraY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.cameraY -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyComma) {
		s.zoomFactor += 0.005
	}
	if ebiten.IsKeyPressed(ebiten.KeyPeriod) {
		s.zoomFactor -= 0.005
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.cameraX = 0
		s.cameraY = 0
		s.zoomFactor = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		ag := s.GameMap.Board.AgentManager.Agents["ag-0"]
		camX, camY := drawing.GetHexGraphicalCenter(ag.Position.Position.X, ag.Position.Position.Y, s.GameMap.Board.HexSize)
		s.cameraX = camX - float32(s.ScreenHeight/2)
		s.cameraY = camY - float32(s.ScreenWidth/4)
		s.zoomFactor = 1.5
	}

	screen.Fill(s.GameMap.BackgroundColor)
	drawing.DrawBoard(screen, s.GameMap.Board, s.cameraX, s.cameraY, s.zoomFactor)
	drawing.DrawAgents(screen, s.GameMap.Board, s.cameraX, s.cameraY, s.zoomFactor)

	widget := NewWidget(0, 0, 50, 50, s.GameMap.Board.AgentManager.Agents["ag-0"])
	widget.Draw(screen)
}
