package main

import (
	"encoding/json"
	"log"
	"os"
	"vuerd/engine"
	"vuerd/models"
	"vuerd/utils"
)

func _main() {
	println(utils.SnakeToCamel("user_id"))
}

func main() {
	var state models.State
	buffer, err := os.ReadFile("db/db.vuerd.json")
	logErr(err)
	json.Unmarshal(buffer, &state)

	t := map[string]string{}
	m := map[string]string{}

	buffer, err = os.ReadFile("data.types.json")
	logErr(err)
	json.Unmarshal(buffer, &t)

	buffer, err = os.ReadFile("prisma.types.json")
	logErr(err)
	json.Unmarshal(buffer, &m)

	nodes := engine.Simplify(state, t, m, utils.SnakeToCamel)
	jsonNodes, _ := json.Marshal(nodes)

	os.WriteFile("data.json", jsonNodes, os.ModeAppend)

	engine.PrismaEngine(nodes)

}

func logErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
