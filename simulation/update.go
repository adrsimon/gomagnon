package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"sync"
)

func (s *Simulation) Update() error {
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
