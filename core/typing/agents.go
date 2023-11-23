package typing

type HumanStats struct {
	Strength    int
	Sociability int
}

type HumanBody struct {
	Age         int
	Gender      rune
	Hungriness  int
	Thirstiness int
}

type Human struct {
	ID    string
	Type  rune
	Body  HumanBody
	Stats HumanStats

	Inventory map[ResourceType]int

	Position       *Hexagone
	Target         *Hexagone
	MovingToTarget bool
	CurrentPath    []*Hexagone
	Board          *Board

	ComOut agentToManager
	ComIn  managerToAgent
}

func NewHuman(id string, Type rune, body HumanBody, stats HumanStats, position *Hexagone, target *Hexagone, movingToTarget bool, currentPath []*Hexagone, board *Board, comOut agentToManager, comIn managerToAgent) *Human {
	return &Human{ID: id, Type: Type, Body: body, Stats: stats, Position: position, Target: target, MovingToTarget: movingToTarget, CurrentPath: currentPath, Board: board, ComOut: comOut, ComIn: comIn}
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
)

func (h *Human) EvaluateOneHex(hex *Hexagone) float64 {
	var score = 0.0

	if hex == nil {
		return score
	}

	switch hex.Resource {
	case ANIMAL:
		score += (float64(h.Body.Hungriness)/100)*AnimalFoodValueMultiplier + 0.01
	case FRUIT:
		score += (float64(h.Body.Hungriness)/100)*FruitFoodValueMultiplier + 0.01
	case ROCK:
		score += 0.5
	case WOOD:
		score += 0.5
	}

	return score
}

func (h *Human) GetNeighborsWithin5() []*Hexagone {
	neighbours := h.Board.GetNeighbours(h.Position)
	visited := make(map[*Hexagone]bool)
	for i := 0; i < 4; i++ {
		for _, neighbour := range neighbours {
			if neighbour == nil {
				continue
			}
			if _, ok := visited[neighbour]; !ok {
				visited[neighbour] = true
				neighbours = append(neighbours, h.Board.GetNeighbours(neighbour)...)
			}
		}
	}

	return neighbours
}

func (h *Human) BestNeighbor(surroundingHexagons []*Hexagone) *Hexagone {
	best := 0.
	indexBest := -1
	for i, v := range surroundingHexagons {
		score := h.EvaluateOneHex(v)
		if score > best {
			indexBest = i
			best = score
		}
	}

	if indexBest != -1 {
		return surroundingHexagons[indexBest]
	}

	valid := false
	randHex := &Hexagone{}
	for !valid {
		randHex = surroundingHexagons[r.Intn(len(surroundingHexagons))]
		if h.Board.isValidHex(randHex) && randHex.Biome.BiomeType != WATER {
			valid = true
		}
	}

	return randHex
}

func (h *Human) UpdateAgent() {
	if !h.MovingToTarget {
		surroundingHexagons := h.GetNeighborsWithin5()
		targetHexagon := h.BestNeighbor(surroundingHexagons)

		res := AStar(*h, targetHexagon)
		path := createPath(res, targetHexagon)
		h.CurrentPath = path
		h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-2]
		h.Target = targetHexagon
		h.MovingToTarget = true
	}

	if h.MovingToTarget && len(h.CurrentPath) > 0 {
		nextHexagon := h.CurrentPath[len(h.CurrentPath)-1]
		h.MoveToHexagon(h.Board.Cases[nextHexagon.Position.X][nextHexagon.Position.Y])
		h.CurrentPath = h.CurrentPath[:len(h.CurrentPath)-1]
	}

	if h.Target.Position == h.Position.Position {
		h.MovingToTarget = false
		h.Target = nil
		if h.Position.Resource != NONE {
			h.ComOut = agentToManager{AgentID: h.ID, Action: "get", Pos: h.Position, commOut: make(chan managerToAgent)}
			h.Board.AgentManager.messIn <- h.ComOut
			h.ComIn = <-h.ComOut.commOut
			if h.ComIn.Valid {
				h.Update(h.Position.Resource)
			}
		}
	}
}

func (h *Human) Update(resource ResourceType) {
	switch resource {
	case ANIMAL:
		h.Body.Hungriness -= 10 * AnimalFoodValueMultiplier
	case FRUIT:
		h.Body.Hungriness -= 10 * FruitFoodValueMultiplier
	case ROCK, WOOD:
		h.Inventory[resource] += 1
	}

	neighbours := h.Board.GetNeighbours(h.Position)
	for _, neighbour := range neighbours {
		if neighbour == nil {
			continue
		}
		if neighbour.Biome.BiomeType == WATER {
			h.Body.Thirstiness -= 10
		}
	}
}

func createPath(maps map[*Hexagone]*Hexagone, hexagon *Hexagone) []*Hexagone {
	path := make([]*Hexagone, 0)
	path = append(path, hexagon)
	val, ok := maps[hexagon]
	for ok {
		path = append(path, val)
		val, ok = maps[val]
	}
	return path
}

func (h *Human) MoveToHexagon(hex *Hexagone) {
	h.Position = hex
	h.Body.Hungriness += 1
	h.Body.Thirstiness += 2
}

func (h *Human) Start() {

}
