package ent

import (
	"os"
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Ent(state types.State, pkg string) {
	var helper = engines.Helper{}
	nodes := engines.Simplify(state, types.DataTypes, EntTypes, helper.Snake)
	files := Schema(nodes, &SchemaConfig{Graphql: true, SingleFile: true})
	nodes = engines.Simplify(state, types.DataTypes, GQLTypes, helper.Snake)
	files = append(files, GQL(nodes, pkg)...)
	cwd, _ := os.Getwd()
	utils.WriteFiles(files, &cwd)
}
