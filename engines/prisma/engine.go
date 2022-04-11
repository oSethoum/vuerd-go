package prisma

import (
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Prisma(State types.State) {
	var helper engines.Helper
	nodes := engines.Simplify(State, types.DataTypes, PrismaTypes, helper.CorrectCamel)
	file := Schema(nodes, "sqlite")
	utils.WriteFile(file, "")
}
