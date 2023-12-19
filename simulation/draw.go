package simulation

import (
	"github.com/adrsimon/gomagnon/core/drawing"
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Simulation) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.cameraX += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.cameraX -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.cameraY += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.cameraY -= 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyComma) {
		s.zoomFactor += 0.005
	}
	if ebiten.IsKeyPressed(ebiten.KeyPeriod) {
		s.zoomFactor -= 0.005
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.zoomFactor = 1
		s.cameraX = 0
		s.cameraY = 0
		s.SelectedAgent = ""
		s.AgentDesc.SetText("Select an agent to see it's statistics")
		s.Selector.SetSelectedEntry(AgentChoice{id: ""})
	}

	if s.SelectedAgent != "" {
		_, ag := s.Board.AgentManager.GetAgent(s.SelectedAgent)
		if ag == nil {
			s.SelectedAgent = ""
			s.AgentDesc.SetText("Select an agent to see it's statistics")
		} else {
			camX, camY := drawing.GetHexGraphicalCenter(ag.Position.Position.X, ag.Position.Position.Y, s.Board.HexSize)
			s.cameraX = camX - float32(s.ScreenHeight/2)
			s.cameraY = camY - float32(s.ScreenWidth/4)
			s.zoomFactor = 1.5
			s.AgentDesc.SetText(ag.String)
		}
	}

	screen.Fill(s.backgroundColor)
	if s.SelectedAgent != "" {
		_, ag := s.Board.AgentManager.GetAgent(s.SelectedAgent)
		if ag != nil {
			drawing.DrawBoard(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor, ag)
			drawing.DrawAgents(screen, []*typing.Agent{ag}, s.cameraX, s.cameraY, s.zoomFactor, s.Board.HexSize, s.Debug)
		}
	} else {
		drawing.DrawBoard(screen, s.Board, s.cameraX, s.cameraY, s.zoomFactor, nil)
		drawing.DrawAgents(screen, s.Board.AgentManager.Agents, s.cameraX, s.cameraY, s.zoomFactor, s.Board.HexSize, s.Debug)
	}

	if !ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.UI.Draw(screen)
	}
}
