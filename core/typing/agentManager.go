package typing

import (
	"fmt"
	"slices"
)

type agentToManager struct {
	AgentID string
	Action  string
	Pos     *Hexagone
	commOut chan managerToAgent
}

type managerToAgent struct {
	Valid    bool
	Map      [][]*Hexagone
	Resource ResourceType
}

type AgentManager struct {
	Map             *[][]*Hexagone
	messIn          chan agentToManager
	Agents          map[string]*Human
	ResourceManager *ResourceManager
	Count           int
}

func NewAgentManager(Map [][]*Hexagone, messIn chan agentToManager, agents map[string]*Human, ressourceManager *ResourceManager) *AgentManager {
	return &AgentManager{Map: &Map, messIn: messIn, Agents: agents, ResourceManager: ressourceManager, Count: 0}
}

func (agMan *AgentManager) startResources() {
	for {
		request := <-agMan.messIn
		agMan.executeResources(request)
	}
}

func (agMan *AgentManager) executeResources(request agentToManager) {
	switch request.Action {
	case "get":
		switch (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource {
		case NONE:
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		default:
			res := (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource
			(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource = NONE
			respawnCD := Randomizer.Intn(20) + 10
			agMan.ResourceManager.RespawnCDs = append(agMan.ResourceManager.RespawnCDs, CoolDown{Current: respawnCD, Resource: res})
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: res}
		}
	case "build":
		(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Hut = &Hut{Position: request.Pos, Inventory: make([]ResourceType, 0), Owner: agMan.Agents[request.AgentID]}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "leave-house":
		ag := agMan.Agents[request.AgentID]
		ag.Hut.Owner = nil
		(*agMan.Map)[ag.Hut.Position.Position.X][ag.Hut.Position.Position.Y].Hut.Owner = nil
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "isHome":
		ag := agMan.Agents[request.AgentID]
		if ag.Procreate.Partner != nil && ag.Procreate.Partner.Position.Position == ag.Hut.Position.Position {
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
		} else {
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		}

	case "procreate":
		ag := agMan.Agents[request.AgentID]
		if ag.Procreate.Partner != nil {
			numChildren := Randomizer.Intn(2) + 1
			for i := 0; i < numChildren; i++ {
				newHuman := MakeChild(ag, ag.Procreate.Partner, agMan.Count)
				if newHuman != nil {
					agMan.Count++
					agMan.Agents[newHuman.ID] = newHuman
				}
			}
			ag.Procreate.Partner.Procreate.Partner = nil
			ag.Procreate.Partner.Procreate.Timer = 200
			ag.Procreate.Partner = nil
			ag.Procreate.Timer = 200
		}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "die":
		agent := agMan.Agents[request.AgentID]
		if agent != nil {
			if agent.Clan != nil {
				if agent.Procreate.Partner != nil {
					agent.Procreate.Partner.Procreate.Partner = nil
					agent.Procreate.Partner.Procreate.Timer = 300
					agent.Procreate.Partner = nil
				}

				if agent.Clan.chief.ID == request.AgentID {
					if len(agMan.Agents[request.AgentID].Clan.members) > 0 {
						agent.Clan.chief = agent.Clan.members[0]
						agent.Hut.Owner = agent.Clan.members[0]
						(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = agent.Clan.members[0]
						agent.Clan.members = agent.Clan.members[1:]
					} else {
						agMan.Agents[request.AgentID].Clan = nil
						agent.Hut.Owner = nil
						(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil

					}
				} else {
					// JE PENSE QUIL FAUT LE VIRER DES MEMBRES DU CLANS ICI
				}
			} else {
				if agent.Hut != nil {
					agent.Hut.Owner = nil
					(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil
				}
			}

			delete(agMan.Agents, request.AgentID)
		}
	case "store-at-home":
		ag := agMan.Agents[request.AgentID]
		if ag.Inventory.Weight <= 0 {
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
			return
		}
		for res, val := range ag.Inventory.Object {
			for i := 0; i < val; i++ {
				ag.Hut.Inventory = append(ag.Hut.Inventory, res)
			}
			ag.Inventory.Object[res] = 0
		}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "eat-from-home":
		ag := agMan.Agents[request.AgentID]
		if !slices.Contains(ag.Hut.Inventory, ANIMAL) && !slices.Contains(ag.Hut.Inventory, FRUIT) {
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
			return
		}
		if slices.Contains(ag.Hut.Inventory, ANIMAL) {
			i := slices.Index(ag.Hut.Inventory, ANIMAL)
			if i == -1 {
				request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
				return
			}
			ag.Hut.Inventory = append(ag.Hut.Inventory[:i], ag.Hut.Inventory[i+1:]...)
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: ANIMAL}
		} else {
			i := slices.Index(ag.Hut.Inventory, FRUIT)
			if i == -1 {
				request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
				return
			}
			ag.Hut.Inventory = append(ag.Hut.Inventory[:i], ag.Hut.Inventory[i+1:]...)
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: FRUIT}
		}
	case "VoteNewPerson":
		valid := agMan.Agents[request.AgentID].Hut.StartNewVote(agMan.Agents[request.AgentID], "VoteNewPerson") //(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Hut.StartNewVote(agMan.Agents[request.AgentID], "VoteNewPerson")
		request.commOut <- managerToAgent{Valid: valid, Map: *agMan.Map, Resource: NONE}
	case "VoteYes":
		valid := agMan.Agents[request.AgentID].Hut.Vote(agMan.Agents[request.AgentID], "VoteYes")
		request.commOut <- managerToAgent{Valid: valid, Map: *agMan.Map, Resource: NONE}

	case "VoteNo":
		valid := agMan.Agents[request.AgentID].Hut.Vote(agMan.Agents[request.AgentID], "VoteNo")
		request.commOut <- managerToAgent{Valid: valid, Map: *agMan.Map, Resource: NONE}

	case "GetResult":
		result := agMan.Agents[request.AgentID].Hut.GetResult(agMan.Agents[request.AgentID])
		request.commOut <- managerToAgent{Valid: result, Map: *agMan.Map, Resource: NONE}
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startResources()
}
