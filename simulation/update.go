package simulation

import (
	"sync"

	"github.com/adrsimon/gomagnon/core/typing"
)

func (s *Simulation) Update() error {
	var wg sync.WaitGroup
	for _, agent := range s.GameMap.Board.Agents {
		wg.Add(1)
		go func(a *typing.Human) {
			defer wg.Done()
			a.UpdateAgent()
			// handle errors
		}(agent)
	}
	wg.Wait()
	return nil
}

// printAgents functon just has to recover every Simulation.GameMap.Board.Agents.Position (string format) and print them on the right hexagone
