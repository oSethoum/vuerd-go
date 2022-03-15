package main

import (
	"vuerd/engines/ent"
	"vuerd/types"
	"vuerd/utils"
)

func main() {
	var state types.State
	utils.ReadJSON(&state, "db/db.vuerd.json")
	ent.Ent(state)
}
