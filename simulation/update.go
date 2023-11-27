package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"sync"
)

func (s *Simulation) Update() error {
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

	s.GameMap.Board.GenerateResources()

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
