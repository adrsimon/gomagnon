package moving

/*
import (
	"github.com/adrsimon/gomagnon/core/typing"

)

func astar(agent Agent, env []typing.Hexagone) Agent { // goal methode agent
	l := make(PriorityQueue, 0)
	l.Push(Item{agent, agent.goal(env)})
	var save map[Agent]Agent
	save[agent]:=nil
	depth:=0
	for l.Len()!=0{
		agTemp := l.Pop()
		for _,v := range agTemp.successor(env){
			g:= hauteurRecherhce(v,save)
			if g >=5{//seuil de recherche
				continue
			}
			val, ok := save[v]
			// If the key exists
			if ok {
				continue
			}
			goalCal:=v.goal(env)
			if goalCal>0.95{//seuil de décision
				return v
			}
			l.Push(v,goalCal-g) // ou plus en fait j'ai pun probleme dans la gestion des goals avec la profondeur du chemin
		}
	}
	bestIfnot95 :=findMax(save)
	return bestIfnot95
}

*/
