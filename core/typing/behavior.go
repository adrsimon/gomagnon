package typing

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

type HumanActions interface {
	DeliberateAtHut()
	Deliberate()
	Act()
	GetNeighboursWithinAcuity() []*Hexagone
}

type HumanBehavior struct {
	H *Agent
}

type ChildBehavior struct {
	C *Agent
}

func (hb *HumanBehavior) GetNeighboursWithinAcuity() []*Hexagone {
	neighbours := hb.H.Board.GetNeighbours(hb.H.Position)
	visited := make(map[*Hexagone]bool)
	for i := 1; i < hb.H.Stats.Acuity; i++ {
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			if _, ok := visited[neighbour]; !ok {
				visited[neighbour] = true
				neighbours = append(neighbours, hb.H.Board.GetNeighbours(neighbour)...)
			}
		}
	}

	return neighbours
}

func (hb *HumanBehavior) DeliberateAtHut() {
	/** If he is tired and have a home, he should sleep **/
	if hb.H.Body.Tiredness > 0 {
		hb.H.Action = SLEEP
		return
	}
	/** If he is home and not partner he should wait **/
	if hb.H.Procreate.Partner != nil && !hb.H.Procreate.isHome {
		hb.H.Action = SLEEP
		return
	}

	if hb.H.Procreate.Partner == nil {
		hb.H.Action = MOVE
	}

	/** If he is home with partner he should procreate **/
	if hb.H.Procreate.Partner != nil && hb.H.Procreate.isHome {
		hb.H.Action = PROCREATE
		return
	}

	/** If he is hungry and have food in home, he should eat **/
	if hb.H.Body.Hungriness > 80 {
		if slices.Contains(hb.H.HutInventoryVision, ANIMAL) || slices.Contains(hb.H.HutInventoryVision, FRUIT) {
			hb.H.Action = EATFROMHOME
			return
		} else {
			hb.H.Action = MOVE
		}
	}

	/** If he has stuff in inventory, he should store it **/
	if hb.H.Inventory.Weight > 0 {
		hb.H.Action = STOREATHOME
		return
	}

	if hb.H.Clan != nil && hb.H.Clan.chief == hb.H && len(hb.H.Clan.members) < 15 && len(hb.H.Clan.members) > 0 && !hb.H.Hut.Ballot.VoteInProgress && !hb.H.Looking4Someone {
		hb.H.Action = CREATEVOTENEWMEMBER
		return
	}
	if hb.H.Clan != nil && hb.H.Hut.Ballot.VoteInProgress && slices.Contains(hb.H.Hut.Ballot.VotersID, hb.H.ID) {
		hb.H.Action = VOTE
		return
	}
	if hb.H.Clan != nil && hb.H.Clan.chief == hb.H && hb.H.Hut.Ballot.VoteInProgress && hb.H.Hut.Ballot.EndTimeVote.Before(time.Now()) {
		hb.H.Action = GETRESULT
		return
	}

}

func (hb *HumanBehavior) Deliberate() {
	hb.H.Action = NOOP

	/** Stacked actions **/
	if len(hb.H.StackAction) > 0 {
		hb.H.Action = hb.H.StackAction[0]
		hb.H.StackAction = hb.H.StackAction[1:]
		return
	}

	/** Early game actions **/
	if hb.H.Hut == nil {
		/** if the agent is tired and don't have a home, he should rest in the nature **/
		if hb.H.Body.Tiredness > 80 {
			hb.H.StackAction = append(hb.H.StackAction, NOOP)
			hb.H.StackAction = append(hb.H.StackAction, NOOP)
			hb.H.StackAction = append(hb.H.StackAction, NOOP)
			return
		}

		/** if he can build a home and don't have one, he should build it **/
		if hb.H.Inventory.Object[WOOD] >= Needs["hut"][WOOD] && hb.H.Inventory.Object[ROCK] >= Needs["hut"][ROCK] {
			for _, v := range hb.H.Board.GetNeighbours(hb.H.Position) {
				if v == nil {
					continue
				}
				if v.Biome == WATER {
					hb.H.Action = BUILD
					return
				}
			}
		}
	}

	if hb.H.Opponent != nil && hb.H.Body.Thirstiness < 80 && hb.H.Body.Hungriness < 80 && hb.H.Body.Tiredness < 80 && hb.H.Fightcooldown == 0 {
		hb.H.Action = FIGHT
		hb.H.Fightcooldown = 100
		return
	}

	/** In Hut actions **/
	if hb.H.Hut != nil && hb.H.Position.Position == hb.H.Hut.Position.Position {
		hb.DeliberateAtHut()
		if hb.H.Action != NOOP {
			return
		}
	}

	/** General actions **/
	if hb.H.Body.Thirstiness > 80 || hb.H.Body.Hungriness > 80 {
		if !hb.H.MovingToTarget {
			hb.H.Action = MOVE
			return
		}
	}

	if hb.H.Hut != nil {
		if len(hb.H.Neighbours) > 0 && hb.H.Clan == nil {
			hb.H.Action = CREATECLAN
			return
		}
	}

	if hb.H.Procreate.Partner != nil && hb.H.Position != hb.H.Hut.Position {
		hb.H.Action = MOVE
		return
	}

	if hb.H.Clan != nil && hb.H.Clan.chief == hb.H && hb.H.Looking4Someone {
		hb.H.Action = LOOK4SOMEONE
		return
	}

	if !hb.H.MovingToTarget {
		hb.H.Action = MOVE
		return
	}

}

func (hb *HumanBehavior) Act() {
	switch hb.H.Action {
	case NOOP:
		hb.H.Body.Tiredness -= 1
	case MOVE:
		if !hb.H.MovingToTarget {
			var targetHexagon *Hexagone

			if hb.H.Hut != nil {
				if hb.H.Body.Tiredness > 80 || hb.H.Procreate.Partner != nil {
					targetHexagon = hb.H.Hut.Position
				} else if hb.H.Body.Hungriness > 80 && (slices.Contains(hb.H.HutInventoryVision, ANIMAL) || slices.Contains(hb.H.HutInventoryVision, FRUIT)) {
					targetHexagon = hb.H.Hut.Position
				}
			}

			if hb.H.Opponent != nil {
				targetHexagon = hb.H.Opponent.Position
			}

			if targetHexagon == nil {
				surroundingHexagons := hb.GetNeighboursWithinAcuity()
				targetHexagon = hb.H.BestNeighbor(surroundingHexagons)
			}

			res := AStar(*hb.H, targetHexagon)
			hb.H.CurrentPath = createPath(res, targetHexagon)
			if len(hb.H.CurrentPath) < 2 {
				hb.H.CurrentPath = nil
				break
			}
			hb.H.CurrentPath = hb.H.CurrentPath[:len(hb.H.CurrentPath)-2]
			hb.H.Target = targetHexagon
			hb.H.MovingToTarget = true
		}

		if hb.H.MovingToTarget && len(hb.H.CurrentPath) > 0 {
			nextHexagon := hb.H.CurrentPath[len(hb.H.CurrentPath)-1]
			hb.H.MoveToHexagon(hb.H.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			hb.H.CurrentPath = hb.H.CurrentPath[:len(hb.H.CurrentPath)-1]
		}

		/** Next move stacking **/
		if hb.H.MovingToTarget && len(hb.H.CurrentPath) > 0 {
			hb.H.StackAction = append(hb.H.StackAction, MOVE)
		}

		if hb.H.Position.Position == hb.H.Target.Position {
			if hb.H.Target.Resource != NONE {
				hb.H.StackAction = append(hb.H.StackAction, GET)
			}
			hb.H.Target = nil
			hb.H.MovingToTarget = false
		}
	case GET:
		if hb.H.Position.Resource != NONE {
			hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "get", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
			hb.H.Board.AgentManager.messIn <- hb.H.ComOut
			hb.H.ComIn = <-hb.H.ComOut.commOut
			if hb.H.ComIn.Valid {
				hb.H.UpdateState(hb.H.ComIn.Resource)
			}
		}
	case BUILD:
		hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "build", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut
		if hb.H.ComIn.Valid {
			hb.H.Hut = &Hut{Position: hb.H.Position, Inventory: make([]ResourceType, 0), Owner: hb.H}
			hb.H.Inventory.Object[WOOD] -= Needs["hut"][WOOD]
			hb.H.Inventory.Object[ROCK] -= Needs["hut"][ROCK]
			hb.H.Inventory.Weight -= WeightWood * float64(Needs["hut"][WOOD])
			hb.H.Inventory.Weight -= WeightRock * float64(Needs["hut"][ROCK])
		}
	case SLEEP:
		if hb.H.Body.Tiredness > 0 {
			hb.H.Body.Tiredness -= 3
			// hb.H.Body.Hungriness += 0.5
			// hb.H.Body.Thirstiness += 0.5
			hb.H.StackAction = append(hb.H.StackAction, SLEEP)
		}
	case STOREATHOME:
		hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "store-at-home", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut
		if hb.H.ComIn.Valid {
			hb.H.Inventory.Weight = 0
		}
	case EATFROMHOME:
		hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "eat-from-home", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut
		if hb.H.ComIn.Valid {
			if hb.H.ComIn.Resource == ANIMAL {
				hb.H.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			} else {
				hb.H.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			}
		}
	case CREATECLAN:
		var bestH *Agent
		if len(hb.H.Neighbours) > 1 {
			bestH = hb.H.BestMatchHuman()
		} else if len(hb.H.Neighbours) == 1 {
			bestH = hb.H.Neighbours[0]
		} else {
			hb.H.StackAction = append(hb.H.StackAction, MOVE)
			break
		}
		if !bestH.Terminated {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: hb.H, Action: "CREATECLAN", commOut: hb.H.AgentCommIn}:
				select {
				case res := <-hb.H.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						clanID := fmt.Sprintf("clan-%s", strings.Split(hb.H.ID, "-")[1])
						clan := &Clan{ID: clanID, members: []*Agent{bestH}, chief: hb.H}
						hb.H.Clan = clan
						bestH.AgentCommIn <- AgentComm{Agent: hb.H, Action: "INVITECLAN", commOut: hb.H.AgentCommIn}
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case CREATEVOTENEWMEMBER:
		hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "VoteNewPerson", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut
		if hb.H.ComIn.Valid {
			hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "VoteYes", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
			hb.H.Board.AgentManager.messIn <- hb.H.ComOut
			hb.H.ComIn = <-hb.H.ComOut.commOut
		}
	case VOTE:
		if Randomizer.Intn(2) >= 1 {
			hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "VoteYes", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		} else {
			hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "VoteNo", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut

	case GETRESULT:
		hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "GetResult", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
		hb.H.Board.AgentManager.messIn <- hb.H.ComOut
		hb.H.ComIn = <-hb.H.ComOut.commOut
		if hb.H.ComIn.Valid {
			hb.H.Looking4Someone = true
		} else {
			hb.H.Looking4Someone = false
		}
	case LOOK4SOMEONE:
		var bestH *Agent
		if len(hb.H.Neighbours) > 1 {
			bestH = hb.H.BestMatchHuman()
		} else if len(hb.H.Neighbours) == 1 {
			bestH = hb.H.Neighbours[0]
		} else {
			hb.H.StackAction = append(hb.H.StackAction, MOVE)
			break
		}
		if !bestH.Terminated {
			select {
			case bestH.AgentCommIn <- AgentComm{Agent: hb.H, Action: "INVITECLAN", commOut: hb.H.AgentCommIn}:
				select {
				case res := <-hb.H.AgentCommIn:
					if res.Action == "ACCEPTCLAN" {
						hb.H.Looking4Someone = false
					} else {
						hb.H.Action = MOVE
					}
				case <-time.After(20 * time.Millisecond):
				}
			case <-time.After(20 * time.Millisecond):

			}
		}
	case PROCREATE:
		if hb.H.Type == 'F' {
			hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "procreate", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
			hb.H.Board.AgentManager.messIn <- hb.H.ComOut
			hb.H.ComIn = <-hb.H.ComOut.commOut
			if !hb.H.ComIn.Valid {
				hb.H.StackAction = append(hb.H.StackAction, MOVE)
			}
		}
	case FIGHT:
		if hb.H.Opponent != nil {
			if !hb.H.Opponent.Terminated {
				select {
				case hb.H.Opponent.AgentCommIn <- AgentComm{Agent: hb.H, Action: "FIGHT", commOut: hb.H.AgentCommIn}:
					select {
					case res := <-hb.H.AgentCommIn:
						if res.Action == "OKFIGHT" {
							StrengthOpp := hb.H.Opponent.Stats.Strength
							StrengthAg := hb.H.Stats.Strength
							ThrowRandom := Randomizer.Float64()
							DifForce := math.Abs(float64(StrengthAg) - float64(StrengthOpp))
							ProbaVictoireAgt1 := 1 / (1 + math.Pow(2, DifForce/10))

							if ThrowRandom < ProbaVictoireAgt1 {
								fmt.Println("\033[96mAgent\033[0m", hb.H.ID, "\033[96mfought and wins against\033[0m", hb.H.Opponent.ID)
								hb.H.Opponent.AgentCommIn <- AgentComm{Agent: hb.H, Action: "YOULOSE", commOut: hb.H.AgentCommIn}
								hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "transfer-inventory", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
								hb.H.Board.AgentManager.messIn <- hb.H.ComOut
								hb.H.Opponent.AgentCommIn <- AgentComm{Agent: hb.H, Action: "LOOTED", commOut: hb.H.AgentCommIn}
							} else {
								fmt.Println("\033[96mAgent\033[0m", hb.H.Opponent.ID, "\033[96mfought and wins against\033[0m", hb.H.ID)
								hb.H.Opponent.AgentCommIn <- AgentComm{Agent: hb.H, Action: "YOUWIN", commOut: hb.H.AgentCommIn}
								res2 := <-hb.H.AgentCommIn
								if res2.Action == "LOOTED" {
									hb.H.ComOut = agentToManager{AgentID: hb.H.ID, Action: "die", Pos: hb.H.Position, commOut: make(chan managerToAgent)}
									hb.H.Board.AgentManager.messIn <- hb.H.ComOut
								}
							}
						} else {
							hb.H.Opponent = nil
							hb.H.Fightcooldown = 300
						}
					case <-time.After(100 * time.Millisecond):
					}
				case <-time.After(100 * time.Millisecond):

				}
			}
		}
		// // declenchement bagare : comparer leur force
		// SociabilityOpp := hb.H.Opponent.Stats.Sociability
		// SociabilityAg := hb.H.Stats.Sociability
		// if float64(SociabilityOpp) > 1.25*float64(SociabilityAg) {
		// 	//l'opposant gagne et l'agent h fuit
		// 	hb.H.StackAction = append(hb.H.StackAction, MOVE)
		// 	fmt.Println("Agent", hb.H.Opponent.ID, " fuit de", hb.H.ID)
		// 	fmt.Println("Position agent1", hb.H.Position.Position, "Position agent2", hb.H.Opponent.Position.Position)
		// 	hb.H.IsInFight = false
		// } else {

		// }
		//4 looting

	default:
		fmt.Println("Should not be here")
	}

}

func (hb *ChildBehavior) GetNeighboursWithinAcuity() []*Hexagone {
	neighbours := hb.C.Board.GetNeighbours(hb.C.Position)
	visited := make(map[*Hexagone]bool)
	for i := 1; i < hb.C.Stats.Acuity; i++ {
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			if _, ok := visited[neighbour]; !ok {
				visited[neighbour] = true
				for _, neigbour2 := range hb.C.Board.GetNeighbours(neighbour) {
					if distance(*hb.C.Hut.Position, *neigbour2) <= 5 {
						neighbours = append(neighbours, neigbour2)
					}
				}
			}
		}
	}

	return neighbours
}

func (hb *ChildBehavior) Deliberate() {
	hb.C.Action = NOOP

	/** Stacked actions **/
	if len(hb.C.StackAction) > 0 {
		hb.C.Action = hb.C.StackAction[0]
		hb.C.StackAction = hb.C.StackAction[1:]
		return
	}

	/** In Hut actions **/
	if hb.C.Position.Position == hb.C.Hut.Position.Position {
		hb.DeliberateAtHut()
		if hb.C.Action != NOOP {
			return
		}
	}

	if hb.C.Body.Age > 5 {
		if hb.C.Body.Thirstiness > 80 || hb.C.Body.Hungriness > 80 {
			if !hb.C.MovingToTarget {
				hb.C.Action = MOVE
				return
			}
		}
		if !hb.C.MovingToTarget {
			hb.C.Action = MOVE
			return
		}
	}

}

func (hb *ChildBehavior) DeliberateAtHut() {
	/** If he is tired and have a home, he should sleep **/
	if hb.C.Body.Tiredness > 0 {
		hb.C.Action = SLEEP
		return
	}

	/** If he is hungry and have food in home, he should eat **/
	if hb.C.Body.Hungriness > 80 {
		if slices.Contains(hb.C.HutInventoryVision, ANIMAL) || slices.Contains(hb.C.HutInventoryVision, FRUIT) {
			hb.C.Action = EATFROMHOME
			return
		} else if hb.C.Body.Age > 5 {
			hb.C.Action = MOVE
		}
	}

	/** If he has stuff in inventory, he should store it **/
	if hb.C.Inventory.Weight > 0 {
		hb.C.Action = STOREATHOME
		return
	}
}

func (hb *ChildBehavior) Act() {
	switch hb.C.Action {
	case NOOP:
		hb.C.Body.Tiredness -= 1
	case MOVE:
		if !hb.C.MovingToTarget {
			var targetHexagon *Hexagone

			if hb.C.Hut != nil {
				if hb.C.Body.Tiredness > 80 || hb.C.Procreate.Partner != nil {
					targetHexagon = hb.C.Hut.Position
				} else if hb.C.Body.Hungriness > 80 && (slices.Contains(hb.C.HutInventoryVision, ANIMAL) || slices.Contains(hb.C.HutInventoryVision, FRUIT)) {
					targetHexagon = hb.C.Hut.Position
				}
			}

			if targetHexagon == nil {
				surroundingHexagons := hb.GetNeighboursWithinAcuity()
				targetHexagon = hb.C.BestNeighbor(surroundingHexagons)
			}

			res := AStar(*hb.C, targetHexagon)
			hb.C.CurrentPath = createPath(res, targetHexagon)
			if len(hb.C.CurrentPath) < 2 {
				hb.C.CurrentPath = nil
				break
			}
			hb.C.CurrentPath = hb.C.CurrentPath[:len(hb.C.CurrentPath)-2]
			hb.C.Target = targetHexagon
			hb.C.MovingToTarget = true
		}

		if hb.C.MovingToTarget && len(hb.C.CurrentPath) > 0 {
			nextHexagon := hb.C.CurrentPath[len(hb.C.CurrentPath)-1]
			hb.C.MoveToHexagon(hb.C.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
			hb.C.CurrentPath = hb.C.CurrentPath[:len(hb.C.CurrentPath)-1]
		}

		/** Next move stacking **/
		if hb.C.MovingToTarget && len(hb.C.CurrentPath) > 0 {
			hb.C.StackAction = append(hb.C.StackAction, MOVE)
		}

		if hb.C.Position.Position == hb.C.Target.Position {
			if hb.C.Target.Resource != NONE {
				hb.C.StackAction = append(hb.C.StackAction, GET)
			}
			hb.C.Target = nil
			hb.C.MovingToTarget = false
		}
	case GET:
		if hb.C.Position.Resource != NONE {
			hb.C.ComOut = agentToManager{AgentID: hb.C.ID, Action: "get", Pos: hb.C.Position, commOut: make(chan managerToAgent)}
			hb.C.Board.AgentManager.messIn <- hb.C.ComOut
			hb.C.ComIn = <-hb.C.ComOut.commOut
			if hb.C.ComIn.Valid {
				hb.C.UpdateState(hb.C.ComIn.Resource)
			}
		}
	case SLEEP:
		if hb.C.Body.Tiredness > 0 {
			hb.C.Body.Tiredness -= 3
			// hb.H.Body.Hungriness += 0.5
			// hb.H.Body.Thirstiness += 0.5
			hb.C.StackAction = append(hb.C.StackAction, SLEEP)
		}
	case STOREATHOME:
		hb.C.ComOut = agentToManager{AgentID: hb.C.ID, Action: "store-at-home", Pos: hb.C.Position, commOut: make(chan managerToAgent)}
		hb.C.Board.AgentManager.messIn <- hb.C.ComOut
		hb.C.ComIn = <-hb.C.ComOut.commOut
		if hb.C.ComIn.Valid {
			hb.C.Inventory.Weight = 0
		}
	case EATFROMHOME:
		hb.C.ComOut = agentToManager{AgentID: hb.C.ID, Action: "eat-from-home", Pos: hb.C.Position, commOut: make(chan managerToAgent)}
		hb.C.Board.AgentManager.messIn <- hb.C.ComOut
		hb.C.ComIn = <-hb.C.ComOut.commOut
		if hb.C.ComIn.Valid {
			if hb.C.ComIn.Resource == ANIMAL {
				hb.C.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
			} else {
				hb.C.Body.Hungriness -= 10 * FruitFoodValueMultiplier
			}
		}
	default:
		fmt.Println("Should not be here")
	}

}
