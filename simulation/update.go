package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
	"sync"
)

func (s *Simulation) Update() error {
	/**
	 * UI UPDATE
	 */
	s.UI.Update()

	if s.Paused {
		return nil
	}

	s.Agents = s.Board.AgentManager.Agents
	m := makeAgentList(s)
	if len(m) != s.SavedLen {
		s.Selector.SetEntries(m)
		s.Selector.SetSelectedEntry(AgentChoice{id: s.SelectedAgent})
		s.SavedLen = len(m)
	}

	ebiten.SetTPS(s.TPS)

	/**
	 * ENVIRONMENT UPDATE
	 */
	for _, line := range s.Board.Cases {
		for _, hex := range line {
			hex.Agents = nil
		}
	}
	for _, agent := range s.Board.AgentManager.Agents {
		if agent.Position != nil {
			s.Board.Cases[agent.Position.Position.X][agent.Position.Position.Y].Agents = append(s.Board.Cases[agent.Position.Position.X][agent.Position.Position.Y].Agents, agent)
		}
	}

	for i := 0; i < len(s.Board.AgentManager.ResourceManager.RespawnCDs); i++ {
		res := s.Board.AgentManager.ResourceManager.RespawnCDs[i]
		res.Current--
		if res.Current == 0 {
			s.Board.AgentManager.ResourceManager.CurrentQuantities[res.Resource]--
			s.Board.AgentManager.ResourceManager.RespawnCDs = append(s.Board.AgentManager.ResourceManager.RespawnCDs[:i], s.Board.AgentManager.ResourceManager.RespawnCDs[i+1:]...)
			i--
		} else {
			s.Board.AgentManager.ResourceManager.RespawnCDs[i] = res
		}
	}

	s.Board.GenerateResources()

	/**
	 * AGENTS UPDATE
	 */
	var wg sync.WaitGroup
	for _, agent := range s.Board.AgentManager.Agents {
		wg.Add(1)
		go func(a *typing.Agent) {
			defer wg.Done()
			a.UpdateAgent()
		}(agent)
	}
	wg.Wait()
	return nil
}
