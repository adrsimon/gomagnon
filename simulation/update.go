package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
	"sync"
)

func (s *Simulation) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.Key1) {
		ebiten.SetTPS(1)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		ebiten.SetTPS(10)
	}

	for _, line := range s.GameMap.Board.Cases {
		for _, hex := range line {
			hex.Agents = nil
		}
	}
	for _, agent := range s.GameMap.Board.AgentManager.Agents {
		if agent.Position != nil {
			s.GameMap.Board.Cases[agent.Position.Position.X][agent.Position.Position.Y].Agents = append(s.GameMap.Board.Cases[agent.Position.Position.X][agent.Position.Position.Y].Agents, agent)
		}
	}

	var wg sync.WaitGroup
	for _, agent := range s.GameMap.Board.AgentManager.Agents {
		wg.Add(1)
		go func(a *typing.Human) {
			defer wg.Done()
			a.UpdateAgent()
		}(agent)
	}
	wg.Wait()
	return nil
}
