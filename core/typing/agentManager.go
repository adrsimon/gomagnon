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
	Agents          map[string]*Agent
	ResourceManager *ResourceManager
	Count           int
}

func NewAgentManager(Map [][]*Hexagone, messIn chan agentToManager, agents map[string]*Agent, ressourceManager *ResourceManager) *AgentManager {
	return &AgentManager{Map: &Map, messIn: messIn, Agents: agents, ResourceManager: ressourceManager, Count: 0}
}

func (agMan *AgentManager) startResources() {
	for {
		request := <-agMan.messIn
		agMan.executeResources(request)
	}
}

func MakeChild(parent1 *Agent, parent2 *Agent, count int) *Agent {
	var failChance int
	var newHuman *Agent
	newHuman = nil
	if parent1.Race == NEANDERTHAL {
		failChance = Randomizer.Intn(2)
	} else {
		failChance = Randomizer.Intn(1)
	}
	if failChance == 0 {
		newHuman = &Agent{
			ID:   fmt.Sprintf("ag-%d", count),
			Type: []rune{'M', 'F'}[Randomizer.Intn(2)],
			Race: parent1.Race,
			Body: HumanBody{
				Thirstiness: 50,
				Hungriness:  50,
				Age:         0,
			},
			Stats: HumanStats{
				Strength:    (parent1.Stats.Strength + parent2.Stats.Strength) / 2,
				Sociability: (parent1.Stats.Sociability + parent2.Stats.Sociability) / 2,
				Acuity:      (parent1.Stats.Acuity + parent2.Stats.Acuity) / 2,
			},
			Position:       parent1.Position,
			Target:         nil,
			MovingToTarget: false,
			CurrentPath:    nil,
			Hut:            parent1.Hut,
			Board:          parent1.Board,
			Inventory:      Inventory{Weight: 0, Object: make(map[ResourceType]int)},
			AgentRelation:  make(map[string]string),
			AgentCommIn:    make(chan AgentComm),
			Clan:           parent1.Clan,
			Procreate:      Procreate{Partner: nil, Timer: 200},
		}
		newHuman.Behavior = &ChildBehavior{C: newHuman}
	}
	return newHuman
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
		fmt.Println("\033[33mNew hut built at\033[0m", request.Pos.Position.X, request.Pos.Position.Y, "\033[33mby\033[0m", request.AgentID)
	case "leave-house":
		ag := agMan.Agents[request.AgentID]
		(*agMan.Map)[ag.Hut.Position.Position.X][ag.Hut.Position.Position.Y].Hut.Owner = nil
		request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
		fmt.Println("\033[33mAgent\033[0m", request.AgentID, "\033[33mleft his house and joined clan\033[0m", ag.Clan.ID)
	case "isHome":
		ag := agMan.Agents[request.AgentID]
		if ag != nil {
			if ag.Procreate.Partner != nil && ag.Procreate.Partner.Position.Position == ag.Hut.Position.Position {
				request.commOut <- managerToAgent{Valid: true, Map: *agMan.Map, Resource: NONE}
			} else {
				request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
			}
		} else {
			request.commOut <- managerToAgent{Valid: false, Map: *agMan.Map, Resource: NONE}
		}
	case "procreate":
		ag := agMan.Agents[request.AgentID]
		if len(ag.Clan.members) > 15 {
			fmt.Println("\033[35mAgent\033[0m", request.AgentID, "\033[35mtried to procreate but his clan\033[0m", ag.Clan.ID, "\033[35mwas too big\033[0m")
			if ag.Procreate.Partner != nil {
				ag.Procreate.Partner.Procreate.Partner = nil
				ag.Procreate.Partner.Procreate.Timer = 100
			}
			ag.Procreate.Partner = nil
			ag.Procreate.Timer = 100
			return
		}
		if ag.Procreate.Partner != nil {
			numChildren := Randomizer.Intn(2) + 1
			for i := 0; i < numChildren; i++ {
				newHuman := MakeChild(ag, ag.Procreate.Partner, agMan.Count)
				if newHuman != nil {
					agMan.Count++
					agMan.Agents[newHuman.ID] = newHuman
					ag.Clan.members = append(ag.Clan.members, newHuman)
					fmt.Println("\033[32mNew human\033[0m", newHuman.ID, "\033[32mborn with race\033[0m", ag.Race, "\033[32mfrom:\033[0m", ag.ID, ag.Procreate.Partner.ID, "\033[32m- There are now\033[0m", len(agMan.Agents), "\033[32magents,\033[0m", len(ag.Clan.members), "\033[32mmembers in the clan\033[0m", ag.Clan.ID)
				}
			}
			ag.Procreate.Partner.Procreate.Partner = nil
			ag.Procreate.Partner.Procreate.Timer = 100
			ag.Procreate.Partner = nil
			ag.Procreate.Timer = 100
		}
	case "die":
		agent, ok := agMan.Agents[request.AgentID]
		if !ok {
			fmt.Println("\033[31mAgent\033[0m", request.AgentID, "\033[31mis supposed to die but he was already dead\033[0m")
			return
		}
		if agent.Clan != nil {
			if len(agent.Clan.members) <= 0 {
				fmt.Println("\033[31mClan\033[0m", agent.Clan.ID, "\033[31m has no more members.\033[0m")
				agent.Clan = nil
				agent.Hut.Owner = nil
				(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil
			} else if agent.Clan.chief.ID == agent.ID {
				newChief := agent.Clan.members[Randomizer.Intn(len(agent.Clan.members))]
				if agent.Hut.Owner.ID == agent.ID {
					agent.Hut.Owner = newChief
					(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = newChief
				}
				agent.Clan.chief = newChief
				for i, v := range agent.Clan.members {
					if v.ID == newChief.ID {
						agent.Clan.members = append(agent.Clan.members[:i], agent.Clan.members[i+1:]...)
					}
				}
				fmt.Println("\033[35mClan\033[0m", agent.Clan.ID, "\033[35mchief died, new chief is\033[0m", newChief.ID)
			} else {
				for i, v := range agent.Clan.members {
					if v.ID == agent.ID {
						agent.Clan.members = append(agent.Clan.members[:i], agent.Clan.members[i+1:]...)
					}
				}
			}
		} else {
			if agent.Hut != nil {
				agent.Hut.Owner = nil
				(*agMan.Map)[agent.Hut.Position.Position.X][agent.Hut.Position.Position.Y].Hut.Owner = nil
			}
		}

		delete(agMan.Agents, request.AgentID)
		fmt.Println("\033[31mAgent\033[0m", agent.ID, "\033[31mdied, there are\033[0m", len(agMan.Agents), "\033[31magents left\033[0m.")
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
		fmt.Println("\033[33mNew vote started by\033[0m", request.AgentID, "\033[33min clan\033[0m", agMan.Agents[request.AgentID].Clan.ID)
	case "VoteYes":
		valid := agMan.Agents[request.AgentID].Hut.Vote(agMan.Agents[request.AgentID], "VoteYes")
		request.commOut <- managerToAgent{Valid: valid, Map: *agMan.Map, Resource: NONE}
	case "VoteNo":
		valid := agMan.Agents[request.AgentID].Hut.Vote(agMan.Agents[request.AgentID], "VoteNo")
		request.commOut <- managerToAgent{Valid: valid, Map: *agMan.Map, Resource: NONE}
	case "GetResult":
		result := agMan.Agents[request.AgentID].Hut.GetResult(agMan.Agents[request.AgentID])
		request.commOut <- managerToAgent{Valid: result, Map: *agMan.Map, Resource: NONE}
		if result {
			fmt.Println("\033[33mNew agent admitted in clan\033[0m", agMan.Agents[request.AgentID].Clan.ID, "\033[33m. Looking for an agent to include in the clan\033[0m")
		} else {
			fmt.Println("\033[35mVote rejected in clan\033[0m", agMan.Agents[request.AgentID].Clan.ID)
		}
	}
}

func (agMan *AgentManager) Start() {
	fmt.Println("Starting agent manager")
	go agMan.startResources()
}
