package main

import (
	"vuerd/engines/gorm"
	"vuerd/types"
	"vuerd/utils"
)

func main() {
	var state types.State
	utils.ReadJSON(&state, "db/db.vuerd.json")
	gorm.Gorm(state)
}
