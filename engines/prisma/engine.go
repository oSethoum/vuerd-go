package prisma

import (
	"vuerd/engines"
	"vuerd/types"
)

func Prisma(State types.State) {
	var helper engines.Helper
	nodes := engines.Simplify(State, types.DataTypes, PrismaTypes, helper.Camel)
	file := Schema(nodes, "sqlite")
	println(file.Buffer)
}
