package main

import (
	"encoding/json"
	"os"
	"vuerd/engine"
	"vuerd/models"
)

func main() {
	var state models.State
	buffer, err := os.ReadFile("./db/db.vuerd.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(buffer, &state)

	nodes := engine.Simplify(state)

	//	jsonNodes, _ := json.Marshal(nodes)

	//os.WriteFile("data.json", jsonNodes, os.ModeAppend)

	engine.PrismaEngine(nodes)

}

func Write() {

}
