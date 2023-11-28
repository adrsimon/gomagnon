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
	Map              *[][]*Hexagone
	messIn           chan agentToManager
	Agents           map[string]*Human
	RessourceManager *ResourceManager
}

func NewAgentManager(Map [][]*Hexagone, messIn chan agentToManager, agents map[string]*Human, ressourceManager *ResourceManager) *AgentManager {
	return &AgentManager{Map: &Map, messIn: messIn, Agents: agents, RessourceManager: ressourceManager}
}

func (agMan *AgentManager) startRessources() {
	for {
		request := <-agMan.messIn
		agMan.executeRessources(request)
	}
}

func (agMan *AgentManager) executeRessources(request agentToManager) {
	switch request.Action {
	case "get":
		switch (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource {
		case NONE:
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		default:
			res := (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource
			(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource = NONE
			agMan.RessourceManager.currentQuantities[res]--
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: res}
		}
	case "build":
		(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Hut = &Hut{Position: request.Pos, Inventory: make(map[ResourceType]int)}
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startRessources()
}
