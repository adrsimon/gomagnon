package typing

type agentToManager struct {
	Agent   Agent
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
	Map          map[string]*Hexagone
	messIn       chan agentToManager
	stackRequest []agentToManager
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
func (agMan *AgentManager) removeAgent(hexagone string, a *Agent) {
	ind := -1
	for i, v := range agMan.Map[hexagone].Agents {
		if v == a {
			ind = i
		}
	}
	if ind != -1 {
		agMan.Map[hexagone].Agents = append(agMan.Map[hexagone].Agents[:ind], agMan.Map[hexagone].Agents[ind+1:]...)
	}
}

func (agMan *AgentManager) execute(request agentToManager) {
	switch request.Action {
	case "walk":
		agMan.Map[request.Pos].Agents = append(agMan.Map[request.Pos].Agents, &request.Agent)
		agMan.removeAgent(request.Agent.Position, &request.Agent)
		request.commOut <- managerToAgent{Valid: true, Map: agMan.Map}
	case "get":
		switch agMan.Map[request.Pos].Resource {
		case 1:
			request.commOut <- managerToAgent{Valid: false, Map: agMan.Map}
		default:
			request.commOut <- managerToAgent{Valid: true, Map: agMan.Map, resource: agMan.Map[request.Pos].Resource}
			agMan.Map[request.Pos].Resource = NONE
		}
	}
}

// mouvement action ou pas
func (agMan *AgentManager) Start() {
	go agMan.startListening()
	go agMan.startAnswering()
}
