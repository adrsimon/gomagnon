package typing

import (
	"fmt"
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

func NewAgentManager(Map [][]*Hexagone, messIn chan agentToManager, agents map[string]*Human, ressourceManager *ResourceManager, count int) *AgentManager {
	return &AgentManager{Map: &Map, messIn: messIn, Agents: agents, ResourceManager: ressourceManager, Count: count}
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
			respawnCD := Randomizer.Intn(20)
			agMan.ResourceManager.RespawnCDs = append(agMan.ResourceManager.RespawnCDs, CoolDown{Current: respawnCD, Resource: res})
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: res}
		}
	case "build":
		(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Hut = &Hut{Position: request.Pos, Inventory: make(map[ResourceType]int), Owner: agMan.Agents[request.AgentID]}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "leave-house":
		ag := agMan.Agents[request.AgentID]
		ag.Hut.Owner = nil
		(*agMan.Map)[ag.Hut.Position.Position.X][ag.Hut.Position.Position.Y].Hut.Owner = nil
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "procreate":
		ag := agMan.Agents[request.AgentID]
		if ag.Procreate.Partner != nil {
			newHuman := MakeChild(ag, ag.Procreate.Partner, agMan.Count)
			if newHuman != nil {
				agMan.Count++
				agMan.Agents[newHuman.ID] = newHuman
			}
		}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	case "die":
		agent := agMan.Agents[request.AgentID]
		if agent.Clan != nil {
			for i, ag := range agent.Clan.members {
				if ag.ID == request.AgentID {
					agent.Clan.members = append(agent.Clan.members[:i], agent.Clan.members[i+1:]...)
					break
				}
			}

			if agent.Clan.chief.ID == request.AgentID {
				if len(agMan.Agents[request.AgentID].Clan.members) > 0 {
					agent.Clan.chief = agent.Clan.members[0]
					agent.Hut.Owner = agent.Clan.members[0]
					(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = agent.Clan.members[0]
				} else {
					agMan.Agents[request.AgentID].Clan = nil
					agent.Hut.Owner = nil
					(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil
				}
			}
		} else {
			if agent.Hut != nil {
				agent.Hut.Owner = nil
				(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil
			}
		}

		delete(agMan.Agents, request.AgentID)
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startResources()
}
