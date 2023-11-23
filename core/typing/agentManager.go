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
	Map          *[][]*Hexagone
	messIn       chan agentToManager
	stackRequest []agentToManager
	Agents       map[string]*Human
	signal       chan struct{}
}

func NewAgentManager(Map [][]*Hexagone, messIn chan agentToManager, stackRequest []agentToManager, agents map[string]*Human) *AgentManager {
	return &AgentManager{Map: &Map, messIn: messIn, stackRequest: stackRequest, Agents: agents, signal: make(chan struct{})}
}

func (agMan *AgentManager) startListening() {
	for {
		request := <-agMan.messIn
		fmt.Println("Request received: ", request)
		agMan.stackRequest = append(agMan.stackRequest, request)
		agMan.signal <- struct{}{}
	}
}

func (agMan *AgentManager) startAnswering() {
	for {
		select {
		case <-agMan.signal:
			request := agMan.stackRequest[0]
			fmt.Println("Request to execute: ", request)
			agMan.execute(request)
			agMan.stackRequest = agMan.stackRequest[1:]
		}
	}
}

func (agMan *AgentManager) execute(request agentToManager) {
	switch request.Action {
	case "get":
		switch (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource {
		case NONE:
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		default:
			res := (*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource
			(*agMan.Map)[request.Pos.Position.X][request.Pos.Position.Y].Resource = NONE
			request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: res}
		}
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startListening()
	go agMan.startAnswering()
}
