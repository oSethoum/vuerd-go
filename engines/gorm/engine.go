package gorm

import (
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Gorm(state types.State) {
	var helper engines.Helper
	var files []types.File
	nodes := engines.Simplify(state, types.DataTypes, GormTypes, helper.Pascal)
	files = append(files, Schema(nodes))
	files = append(files, DB(nodes, "sqlite", "api"))
	files = append(files, Service(nodes, "api")...)
	dir := "out"
	utils.WriteFiles(files, &dir)
}
