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

func createPath(maps map[*Point2D]*Point2D, hexagon *Hexagone) []*Point2D {
	path := make([]*Point2D, 0)
	path = append(path, hexagon.Position)
	val, ok := maps[hexagon.Position]
	for ok {
		path = append(path, val)
		val, ok = maps[val]
	}
	return path
}

func HauteurNoeud(node *Point2D, save map[*Point2D]*Point2D) int {
	cnt := 1
	for save[node] != nil {
		cnt++
		parent := save[node]
		node = parent
	}
	return cnt
}

func AStar(agent Agent, goal Hexagone) map[*Point2D]*Point2D {
	l := make(PriorityQueue, 0)
	heap.Init(&l)
	l.Push(Item{agent, distance(*agent.Position, goal)})
	save := make(map[*Point2D]*Point2D)
	save[agent.Position.Position] = nil

	for l.Len() != 0 {
		a := heap.Pop(&l).(*Item)
		agTemp := a.value

		if agTemp.Position == &goal {
			return save
		}

		for _, succ := range agTemp.Board.GetNeighbours(*agTemp.Position) {
			if succ.Biome == DEEP_WATER {
				continue
			}
			_, ok := save[succ.Position]
			if !ok {
				save[succ.Position] = agTemp.Position.Position
				newHum := NewHuman(agent.ID, agent.Type, agent.Race, agent.Body, agent.Stats, agent.MapVision, succ, agent.Target, agent.MovingToTarget, agent.CurrentPath, agent.Board, agent.ComOut, agent.ComIn, agent.Hut, agent.Inventory, agent.AgentRelation, agent.Procreate)
				g := HauteurNoeud(succ.Position, save)
				dist := distance(*newHum.Position, goal)
				l.Push(Item{*newHum, dist + float64(g) + a.priority})
			}
		}
	}

	return save
}
