package moving

import (
	"fmt"
	"github.com/adrsimon/gomagnon/core/typing"
	"math"
)

func distance(from typing.Hexagone, to typing.Hexagone) float64 {
	q1, r1 := from.OodrToAxial()
	q2, r2 := to.OodrToAxial()
	q3, r3 := q1-q2, r1-r2
	return (math.Abs(float64(q3)) + math.Abs(float64(q3+r3)) + math.Abs(float64(r3))) / 2
}

func HauteurNoeud(s string, d map[string]string) int {
	cnt := 1
	for d[s] != "" {
		cnt++
		parent := d[s][0]
		s = parent
	}
	return cnt
}

func Astar(agent typing.Human, goal *typing.Hexagone) (typing.Human, map[string]string) { // goal methode agent
	l := make(PriorityQueue, 50)
	l.Push(Item{agent, distance(*agent.Map[agent.Position], goal)})
	var save map[string]string
	save[agent.Position] = ""
	for l.Len() != 0 {
		var agTemp typing.Human
		a := l.Pop()
		agTemp = a.value
		for _, succ := range agTemp.GetNeighbours(agTemp.Position) {
			val, ok := save[fmt.Sprintf("%d:%d", succ.Position.X, succ.Position.Y)]
			// If the key exists
			if !ok {
				save[fmt.Sprintf("%d:%d", succ.Position.X, succ.Position.Y)] = fmt.Sprintf("%d:%d", agTemp.Position.X, agTemp.Position.Y)
				if succ == goal {
					return succ, save // bon vasy j'ai la flemme ce soir
				}

				g := HauteurNoeud(fmt.Sprintf("%d:%d", v.Position.X, v.Position.Y), save)
				l.Push(v, goalCal-g) // ou plus en fait j'ai pun probleme dans la gestion des goals avec la profondeur du chemin
			}

		}
	}
	bestIfnot95 := findMax(save)
	return bestIfnot95
}
