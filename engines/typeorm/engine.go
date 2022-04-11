package typeorm

import (
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Engine(state types.State) {
	h := engines.Helper{}
	nodes := engines.Simplify(state, types.DataTypes, TypeOrmTypes, h.CorrectCamel)
	files := Schema(nodes)
	files = append(files, GraphQL(nodes)...)
	utils.WriteFiles(files, "src/graphql")
}
