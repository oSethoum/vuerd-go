package main

import (
	"vuerd/engines/prisma"
	"vuerd/types"
	"vuerd/utils"
)

func main() {
	var state types.State
	utils.ReadJSON(&state, "db/db.vuerd.json")
	prisma.Prisma(state)
}
