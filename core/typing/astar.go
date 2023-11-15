package typing

import (
	"math"
)

func distance(from Hexagone, to Hexagone) float64 {
	q1, r1 := from.OddRToAxial()
	q2, r2 := to.OddRToAxial()
	q3, r3 := q1-q2, r1-r2
	return (math.Abs(float64(q3)) + math.Abs(float64(q3+r3)) + math.Abs(float64(r3))) / 2
}

func AStar(human Human, goal *Hexagone) []*Hexagone {
	l := make(PriorityQueue, 0)
	l.Push(&Item{value: human, priority: 0})
	visited := make(map[*Hexagone]bool)
	dist := make(map[*Hexagone]float64)
	dist[human.Position] = 0
	d := make(map[*Hexagone]*Hexagone)
	d[human.Position] = nil
	for len(l) > 0 {
		current := l.Pop().value
		visited[current.Position] = true
		if current.Position == goal {
			break
		}
		for _, v := range current.GetNeighborsWithin5() {
			if !visited[v] && v != nil {
				if dist[v] == 0 || dist[current.Position]+distance(*v, *current.Board.Cases[current.Position.ToString()])+current.EvaluateOneHex(v) < dist[v] {
					dist[v] = dist[current.Position] + distance(*v, *current.Board.Cases[current.Position.ToString()]) + current.EvaluateOneHex(v)
					d[v] = current.Position
					l.Push(&Item{value: current, priority: dist[v] + distance(*v, *goal)})
				}
			}
		}
	}
	path := make([]*Hexagone, 0)
	path = append(path, goal)
	for _, v := range l {
		path = append(path, v.value.Position)
	}
	return path
}
