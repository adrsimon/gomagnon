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
		s.zoomFactor = 0.6
		s.cameraX = 0
		s.cameraY = 0
		s.SelectedAgent = ""
		s.Selector.SetSelectedEntry(AgentChoice{id: ""})
	}

	if s.SelectedAgent != "" {
		ag, ok := s.Board.AgentManager.Agents[s.SelectedAgent]
		if !ok {
			s.SelectedAgent = ""
			s.AgentDesc.SetText("Select an agent to see it's statistics")
		} else {
			camX, camY := drawing.GetHexGraphicalCenter(ag.Position.Position.X, ag.Position.Position.Y, s.Board.HexSize)
			s.cameraX = camX - float32(s.ScreenHeight/2)
			s.cameraY = camY - float32(s.ScreenWidth/4)
			s.zoomFactor = 1.5
		}
	}

	screen.Fill(s.backgroundColor)
	drawing.DrawBoard(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor)
	drawing.DrawAgents(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor, s.Debug)

	if !ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.UI.Draw(screen)
	}
}
