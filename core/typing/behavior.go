package typing

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type HumanActions interface {
	DeliberateAtHut()
	Deliberate()
	Act()
	UpdateAgent()
}

type HumanBehavior struct {
	H *Human
}

type ChildBehavior struct {
	C *Child
}

func (hb *HumanBehavior) DeliberateAtHut() {
	/** If he is tired and have a home, he should sleep **/
	if hb.H.Ag.Body.Tiredness > 0 {
		hb.H.Ag.Action = SLEEP
		return
	}
	/** If he is home and not partner he should wait **/
	if hb.H.Ag.Procreate.Partner != nil && !hb.H.Ag.Procreate.isHome {
		hb.H.Ag.Action = SLEEP
		return
	}

	if hb.H.Ag.Procreate.Partner == nil {
		hb.H.Ag.Action = MOVE
	}

	/** If he is home with partner he should procreate **/
	if hb.H.Ag.Procreate.Partner != nil && hb.H.Ag.Procreate.isHome {
		hb.H.Ag.Action = PROCREATE
		return
	}

	/** If he is hungry and have food in home, he should eat **/
	if hb.H.Ag.Body.Hungriness > 80 {
		if slices.Contains(hb.H.Ag.HutInventoryVision, ANIMAL) || slices.Contains(hb.H.Ag.HutInventoryVision, FRUIT) {
			hb.H.Ag.Action = EATFROMHOME
			return
		} else {
			hb.H.Ag.Action = MOVE
		}
	}

	/** If he has stuff in inventory, he should store it **/
	if hb.H.Ag.Inventory.Weight > 0 {
		hb.H.Ag.Action = STOREATHOME
		return
	}

	if hb.H.Ag.Clan != nil && hb.H.Ag.Clan.chief == hb.H.Ag && len(hb.H.Ag.Clan.members) < 15 && len(hb.H.Ag.Clan.members) > 0 && hb.H.Ag.Hut.Ballot.VoteInProgress == false && hb.H.Ag.Looking4Someone == false {
		hb.H.Ag.Action = CREATEVOTENEWMEMBER
		return
	}
	if hb.H.Ag.Clan != nil && hb.H.Ag.Hut.Ballot.VoteInProgress && slices.Contains(hb.H.Ag.Hut.Ballot.VotersID, hb.H.Ag.ID) {
		hb.H.Ag.Action = VOTE
		return
	}
	if hb.H.Ag.Clan != nil && hb.H.Ag.Clan.chief == hb.H.Ag && hb.H.Ag.Hut.Ballot.VoteInProgress && hb.H.Ag.Hut.Ballot.EndTimeVote.Before(time.Now()) {
		hb.H.Ag.Action = GETRESULT
		return
	}

}

func (hb *HumanBehavior) Deliberate() {
	hb.H.Ag.Action = NOOP

	/** Stacked actions **/
	if len(hb.H.Ag.StackAction) > 0 {
		hb.H.Ag.Action = hb.H.Ag.StackAction[0]
		hb.H.Ag.StackAction = hb.H.Ag.StackAction[1:]
		return
	}

	/** Early game actions **/
	if hb.H.Ag.Hut == nil {
		/** if the agent is tired and don't have a home, he should rest in the nature **/
		if hb.H.Ag.Body.Tiredness > 80 {
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, NOOP)
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, NOOP)
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, NOOP)
			return
		}

		/** if he can build a home and don't have one, he should build it **/
		if hb.H.Ag.Inventory.Object[WOOD] >= Needs["hut"][WOOD] && hb.H.Ag.Inventory.Object[ROCK] >= Needs["hut"][ROCK] {
			for _, v := range hb.H.Ag.Board.GetNeighbours(hb.H.Ag.Position) {
				if v == nil {
					continue
				}
				if v.Biome == WATER {
					hb.H.Ag.Action = BUILD
					return
				}
			}
		}
	}

	/** In Hut actions **/
	if hb.H.Ag.Hut != nil && hb.H.Ag.Position.Position == hb.H.Ag.Hut.Position.Position {
		hb.DeliberateAtHut()
		if hb.H.Ag.Action != NOOP {
			return
		}
	}

	/** General actions **/
	if hb.H.Ag.Body.Thirstiness > 80 || hb.H.Ag.Body.Hungriness > 80 {
		if !hb.H.Ag.MovingToTarget {
			hb.H.Ag.Action = MOVE
			return
		}
	}

	if hb.H.Ag.Hut != nil {
		if len(hb.H.Ag.Neighbours) > 0 && hb.H.Ag.Clan == nil {
			hb.H.Ag.Action = CREATECLAN
			return
		}
	}

	if hb.H.Ag.Procreate.Partner != nil && hb.H.Ag.Position != hb.H.Ag.Hut.Position {
		hb.H.Ag.Action = MOVE
		return
	}

	if hb.H.Ag.Clan != nil && hb.H.Ag.Clan.chief == hb.H.Ag && hb.H.Ag.Looking4Someone {
		hb.H.Ag.Action = LOOK4SOMEONE
		return
	}

	if !hb.H.Ag.MovingToTarget {
		hb.H.Ag.Action = MOVE
		return
	}

}

func (hb *HumanBehavior) Act() {
	switch hb.H.Ag.Action {
	case NOOP:
		hb.H.Ag.Body.Tiredness -= 1
	case MOVE:
		if !hb.H.Ag.MovingToTarget {
			var targetHexagon *Hexagone

			if hb.H.Ag.Hut != nil {
				if hb.H.Ag.Body.Tiredness > 80 || hb.H.Ag.Procreate.Partner != nil {
					targetHexagon = hb.H.Ag.Hut.Position
				} else if hb.H.Ag.Body.Hungriness > 80 && (slices.Contains(hb.H.Ag.HutInventoryVision, ANIMAL) || slices.Contains(hb.H.Ag.HutInventoryVision, FRUIT)) {
					targetHexagon = hb.H.Ag.Hut.Position
				}
			}

			if targetHexagon == nil {
				surroundingHexagons := hb.H.Ag.GetNeighboursWithinAcuity()
				targetHexagon = hb.H.Ag.BestNeighbor(surroundingHexagons)
			}

			res := AStar(*hb.H.Ag, targetHexagon)
			hb.H.Ag.CurrentPath = createPath(res, targetHexagon)
			if len(hb.H.Ag.CurrentPath) < 2 {
				hb.H.Ag.CurrentPath = nil
				break
			}
			hb.H.Ag.CurrentPath = hb.H.Ag.CurrentPath[:len(hb.H.Ag.CurrentPath)-2]
			hb.H.Ag.Target = targetHexagon
			hb.H.Ag.MovingToTarget = true
		}

		if hb.H.Ag.MovingToTarget && len(hb.H.Ag.CurrentPath) > 0 {
			nextHexagon := hb.H.Ag.CurrentPath[len(hb.H.Ag.CurrentPath)-1]
			hb.H.Ag.MoveToHexagon(hb.H.Ag.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			hb.H.Ag.CurrentPath = hb.H.Ag.CurrentPath[:len(hb.H.Ag.CurrentPath)-1]
		}

		/** Next move stacking **/
		if hb.H.Ag.MovingToTarget && len(hb.H.Ag.CurrentPath) > 0 {
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, MOVE)
		}

		if hb.H.Ag.Position.Position == hb.H.Ag.Target.Position {
			if hb.H.Ag.Target.Resource != NONE {
				hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, GET)
			}
			hb.H.Ag.Target = nil
			hb.H.Ag.MovingToTarget = false
		}
	case GET:
		if hb.H.Ag.Position.Resource != NONE {
			hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "get", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
			hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
			hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
			if hb.H.Ag.ComIn.Valid {
				hb.H.Ag.UpdateState(hb.H.Ag.ComIn.Resource)
			}
		}
	case BUILD:
		hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "build", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
			hb.H.Ag.Hut = &Hut{Position: hb.H.Ag.Position, Inventory: make([]ResourceType, 0), Owner: hb.H.Ag}
			hb.H.Ag.Inventory.Object[WOOD] -= Needs["hut"][WOOD]
			hb.H.Ag.Inventory.Object[ROCK] -= Needs["hut"][ROCK]
			hb.H.Ag.Inventory.Weight -= WeightWood * float64(Needs["hut"][WOOD])
			hb.H.Ag.Inventory.Weight -= WeightRock * float64(Needs["hut"][ROCK])
		}
	case SLEEP:
		if hb.H.Ag.Body.Tiredness > 0 {
			hb.H.Ag.Body.Tiredness -= 3
			// hb.H.Ag.Body.Hungriness += 0.5
			// hb.H.Ag.Body.Thirstiness += 0.5
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, SLEEP)
		}
	case STOREATHOME:
		hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "store-at-home", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
			hb.H.Ag.Inventory.Weight = 0
		}
	case EATFROMHOME:
		hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "eat-from-home", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
			if hb.H.Ag.ComIn.Resource == ANIMAL {
				hb.H.Ag.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			} else {
				hb.H.Ag.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			}
		}
	case CREATECLAN:
		var bestH *Agent
		if len(hb.H.Ag.Neighbours) > 1 {
			//TO DEVELOPP bestH=find bestMatchHuman(humans)
			bestH = hb.H.Ag.Neighbours[0] // waiting function
		} else if len(hb.H.Ag.Neighbours) == 1 {
			bestH = hb.H.Ag.Neighbours[0]
		} else {
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, MOVE)
			break
		}
		if !bestH.Terminated {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: hb.H.Ag, Action: "CREATECLAN", commOut: hb.H.Ag.AgentCommIn}:
				select {
				case res := <-hb.H.Ag.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						clanID := fmt.Sprintf("clan-%s", strings.Split(hb.H.Ag.ID, "-")[1])
						clan := &Clan{ID: clanID, members: []*Agent{bestH}, chief: hb.H.Ag}
						hb.H.Ag.Clan = clan
						bestH.AgentCommIn <- AgentComm{Agent: hb.H.Ag, Action: "INVITECLAN", commOut: hb.H.Ag.AgentCommIn}
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case CREATEVOTENEWMEMBER:
		hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "VoteNewPerson", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
			hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "VoteYes", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
			hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
			hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
			if hb.H.Ag.ComIn.Valid {
			}
		}
	case VOTE:
		if Randomizer.Intn(2) >= 1 {
			hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "VoteYes", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		} else {
			hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "VoteNo", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
		}
	case GETRESULT:
		hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "GetResult", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
		hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		hb.H.Ag.ComIn = <-hb.H.Ag.ComOut.commOut
		if hb.H.Ag.ComIn.Valid {
			hb.H.Ag.Looking4Someone = true
		} else {
			hb.H.Ag.Looking4Someone = false
		}
	case LOOK4SOMEONE:
		var bestH *Agent
		if len(hb.H.Ag.Neighbours) > 1 {
			//TO DEVELOPP bestH=find bestMatchHuman(humans)
			bestH = hb.H.Ag.Neighbours[0] // waiting function
		} else if len(hb.H.Ag.Neighbours) == 1 {
			bestH = hb.H.Ag.Neighbours[0]
		} else {
			hb.H.Ag.StackAction = append(hb.H.Ag.StackAction, MOVE)
			break
		}
		if bestH.Terminated == false {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: hb.H.Ag, Action: "INVITECLAN", commOut: hb.H.Ag.AgentCommIn}:
				select {
				case res := <-hb.H.Ag.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						hb.H.Ag.Looking4Someone = false
					} else {
						hb.H.Ag.Action = MOVE
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case PROCREATE:
		if hb.H.Ag.Type == 'F' {
			hb.H.Ag.ComOut = agentToManager{AgentID: hb.H.Ag.ID, Action: "procreate", Pos: hb.H.Ag.Position, commOut: make(chan managerToAgent)}
			hb.H.Ag.Board.AgentManager.messIn <- hb.H.Ag.ComOut
		}
	default:
		fmt.Println("Should not be here")
	}

}

func (hb *HumanBehavior) UpdateAgent() {
	hb.H.Ag.Terminated = false
	hb.H.Ag.Perceive()
	hb.Deliberate()
	hb.Act()
	select {
	case res := <-hb.H.Ag.AgentCommIn:
		hb.H.Ag.Terminated = true
		hb.H.Ag.AnswerAgents(res)
	default:
		hb.H.Ag.Terminated = true
	}
	hb.H.Ag.CloseUpdate()
}
