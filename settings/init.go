package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

var Setting Settings

type Resources struct {
	MaxAnimals int `json:"maxAnimals"`
	MaxFruits  int `json:"maxFruits"`
	MaxWoods   int `json:"maxWoods"`
	MaxRocks   int `json:"maxRocks"`
}

type Size struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type World struct {
	Seed      int64     `json:"seed"`
	Resources Resources `json:"resources"`
	Size      Size      `json:"size"`
}

type Agents struct {
	InitialNumber int `json:"initialNumber"`
}

type Settings struct {
	Agents Agents `json:"agents"`
	World  World  `json:"world"`
}

func init() {
	settingsData, err := os.ReadFile("settings/settings.json")
	if err != nil {
		fmt.Println("Error reading settings.json:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(settingsData, &Setting)
	if err != nil {
		fmt.Println("Error unmarshalling settings.json:", err)
		os.Exit(1)
	}
	fmt.Println(Setting)
}
