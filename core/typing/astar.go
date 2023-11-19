package typing

import (
	"container/heap"
	"math"
)

func distance(from Hexagone, to Hexagone) float64 {
	q1, r1 := from.OddRToAxial()
	q2, r2 := to.OddRToAxial()
	q3, r3 := q1-q2, r1-r2
	return (math.Abs(float64(q3)) + math.Abs(float64(q3+r3)) + math.Abs(float64(r3))) / 2
}

func HauteurNoeud(node string, save map[string]string) int {
	cnt := 1
	for save[node] != "" {
		cnt++
		parent := save[node]
		node = parent
	}
	return cnt
}

func AStar(human Human, goal *Hexagone) []*Hexagone {
	l := make(PriorityQueue, 0)
	heap.Init(&l)
	l.Push(&Item{value: human, priority: 0})
	visited := make(map[*Hexagone]bool)
	dist := make(map[*Hexagone]float64)
	dist[human.Position] = 0
	d := make(map[*Hexagone]*Hexagone)
	d[human.Position] = nil
	for len(l) > 0 {
		current := heap.Pop(&l).(*Item).value
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
func Astarnew(agent Human, goal *Hexagone) map[string]string { // goal methode agent
	l := make(PriorityQueue, 0)
	heap.Init(&l)
	l.Push(Item{agent, distance(*agent.Position, *goal)})
	save := make(map[string]string, 0)
	save[agent.Position.ToString()] = ""
	for l.Len() != 0 {
		var agTemp Human
		a := heap.Pop(&l).(*Item)
		agTemp = a.value
		for _, succ := range agTemp.Board.GetNeighbours(agTemp.Position) {
			_, ok := save[succ.ToString()]
			// If the key exists
			if !ok {
				save[succ.ToString()] = agTemp.Position.ToString()
				newHum := NewHuman(agent.id, agent.Type, agent.Body, agent.Stats, succ, agent.Target, agent.MovingToTarget, agent.CurrentPath, agent.Board, agent.ComOut, agent.ComIn)
				if succ.ToString() == goal.ToString() {
					return save
				}
				g := HauteurNoeud(succ.ToString(), save)
				dist := distance(*newHum.Position, *goal)
				l.Push(Item{*newHum, dist + float64(g)}) // ou plus en fait j'ai pun probleme dans la gestion des goals avec la profondeur du chemin
			}

		}
	}
	return save
}
