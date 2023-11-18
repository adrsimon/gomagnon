package typing

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
	Map          map[string]*Hexagone
	messIn       chan agentToManager
	stackRequest []agentToManager
	Agents       map[string]*Human
}

func NewAgentManager(Map map[string]*Hexagone, messIn chan agentToManager, stackRequest []agentToManager, agents map[string]*Human) *AgentManager {
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
func (agMan *AgentManager) removeAgent(hexagone string, a *Human) {
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
		agMan.Map[request.Pos].Agents = append(agMan.Map[request.Pos].Agents, agMan.Agents[request.AgentID])
		agMan.removeAgent(request.Pos, agMan.Agents[request.AgentID])
		request.commOut <- managerToAgent{Valid: true, Map: agMan.Map, Resource: NONE}
	case "get":
		switch agMan.Map[request.Pos].Resource {
		case 1:
			request.commOut <- managerToAgent{Valid: false, Map: agMan.Map, Resource: NONE}
		default:
			res := agMan.Map[request.Pos].Resource
			agMan.Map[request.Pos].Resource = NONE
			request.commOut <- managerToAgent{Valid: true, Map: agMan.Map, Resource: res}
		}
	}
}

func (agMan *AgentManager) Start() {
	go agMan.startListening()
	go agMan.startAnswering()
}
