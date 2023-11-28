package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"sync"
)

func (s *Simulation) Update() error {
	/**
	 * ENVIRONMENT UPDATE
	 */
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

	for i := 0; i < len(s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs); i++ {
		res := s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs[i]
		res.Current--
		if res.Current == 0 {
			s.GameMap.Board.AgentManager.RessourceManager.CurrentQuantities[res.Resource]--
			s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs = append(s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs[:i], s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs[i+1:]...)
			i--
		} else {
			s.GameMap.Board.AgentManager.RessourceManager.RespawnCDs[i] = res
		}
	}

	s.GameMap.Board.GenerateResources()

	/**
	 * AGENTS UPDATE
	 */
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
