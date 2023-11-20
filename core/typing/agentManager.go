package typing

import "fmt"

type agentToManager struct {
	AgentID string
	Action  string
	Pos     string
	commOut chan managerToAgent
}

type managerToAgent struct {
	Valid    bool
	Map      map[string]*Hexagone
	Resource ResourceType
}

type AgentManager struct {
	Map          *map[string]*Hexagone
	messIn       chan agentToManager
	stackRequest []agentToManager
	Agents       map[string]*Human
}

func NewAgentManager(Map *map[string]*Hexagone, messIn chan agentToManager, stackRequest []agentToManager, agents map[string]*Human) *AgentManager {
	return &AgentManager{Map: Map, messIn: messIn, stackRequest: stackRequest, Agents: agents}
}

func (agMan *AgentManager) startListening() {
	for {
		request := <-agMan.messIn
		agMan.stackRequest = append(agMan.stackRequest, request)
	}
}

func (agMan *AgentManager) startAnswering() {
	for {
		if len(agMan.stackRequest) > 0 {
			request := agMan.stackRequest[0]
			agMan.execute(request)
			agMan.stackRequest = agMan.stackRequest[1:]
		}
	}
}

func (agMan *AgentManager) execute(request agentToManager) {
	switch request.Action {
	case "get":
		switch (*agMan.Map)[request.Pos].Resource {
		case NONE:
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		default:
			res := (*agMan.Map)[request.Pos].Resource
			(*agMan.Map)[request.Pos].Resource = NONE
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: res}
		}
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startListening()
	go agMan.startAnswering()
}
