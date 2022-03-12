package gorm

import (
	"vuerd/engines"
	"vuerd/types"
)

func Gorm(state types.State) {
	var helper engines.Helper
	nodes := engines.Simplify(state, types.DataTypes, GormTypes, helper.Pascal)
	file := Schema(nodes)
	println(file.Buffer)
	Migration(nodes)
}
