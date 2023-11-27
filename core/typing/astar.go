package typing

import (
	"container/heap"
	"math"
)

func distance(from Hexagone, to Hexagone) float64 {
	q1, r1 := from.EvenRToAxial()
	q2, r2 := to.EvenRToAxial()
	q3, r3 := q1-q2, r1-r2
	return (math.Abs(float64(q3)) + math.Abs(float64(q3+r3)) + math.Abs(float64(r3))) / 2
}

func createPath(maps map[*Hexagone]*Hexagone, hexagon *Hexagone) []*Hexagone {
	path := make([]*Hexagone, 0)
	path = append(path, hexagon)
	val, ok := maps[hexagon]
	for ok {
		path = append(path, val)
		val, ok = maps[val]
	}
	return path
}

func HauteurNoeud(node *Hexagone, save map[*Hexagone]*Hexagone) int {
	cnt := 1
	for save[node] != nil {
		cnt++
		parent := save[node]
		node = parent
	}
	return cnt
}

func AStar(agent Human, goal *Hexagone) map[*Hexagone]*Hexagone {
	l := make(PriorityQueue, 0)
	heap.Init(&l)
	l.Push(Item{agent, distance(*agent.Position, *goal)})
	save := make(map[*Hexagone]*Hexagone)
	save[agent.Position] = nil

	for l.Len() != 0 {
		a := heap.Pop(&l).(*Item)
		agTemp := a.value

		if agTemp.Position == goal {
			return save
		}

		for _, succ := range agTemp.Board.GetNeighbours(agTemp.Position) {
			if succ.Biome.BiomeType == WATER {
				continue
			}
			_, ok := save[succ]
			if !ok {
				save[succ] = agTemp.Position
				newHum := NewHuman(agent.ID, agent.Type, agent.Race, agent.Body, agent.Stats, succ, agent.Target, agent.MovingToTarget, agent.CurrentPath, agent.Board, agent.ComOut, agent.ComIn, agent.Hut, agent.Inventory, agent.AgentRelation)
				g := HauteurNoeud(succ, save)
				dist := distance(*newHum.Position, *goal)
				l.Push(Item{*newHum, dist + float64(g) + a.priority})
			}
		}
	}

	return save
}
