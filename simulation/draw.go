package simulation

import (
	"github.com/adrsimon/gomagnon/core/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Simulation) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.cameraX += 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.cameraX -= 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.cameraY += 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.cameraY -= 20
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
		ag := s.Board.AgentManager.Agents["ag-0"]
		camX, camY := drawing.GetHexGraphicalCenter(ag.Position.Position.X, ag.Position.Position.Y, s.Board.HexSize)
		s.cameraX = camX - float32(s.ScreenHeight/2)
		s.cameraY = camY - float32(s.ScreenWidth/4)
		s.zoomFactor = 1.5
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.Debug = false
	} else {
		s.Debug = true
	}

	screen.Fill(s.backgroundColor)
	drawing.DrawBoard(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor)
	drawing.DrawAgents(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor, s.Debug)

	s.UI.Draw(screen)
}
