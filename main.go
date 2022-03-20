package main

import (
	"os"
	"path"
	"vuerd/engines/ent"
	"vuerd/types"
	"vuerd/utils"
)

func main() {

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var state types.State
	utils.ReadJSON(&state, path.Join(cwd, "db/db.vuerd.json"))
	ent.Ent(state)
}
