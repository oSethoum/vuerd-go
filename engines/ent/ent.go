package ent

import (
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Ent(state types.State) {
	var helper = engines.Helper{}
	nodes := engines.Simplify(state, types.DataTypes, EntTypes, helper.Snake)
	files := Schema(nodes, &SchemaConfig{Graphql: true})
	dir := "out"
	utils.WriteFiles(files, &dir)
}
