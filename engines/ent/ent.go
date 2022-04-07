package ent

import (
	"os"
	"os/exec"
	"vuerd/engines"
	"vuerd/types"
	"vuerd/utils"
)

func Ent(state types.State, pkg string) {
	var helper = engines.Helper{}
	nodes := engines.Simplify(state, types.DataTypes, EntTypes, helper.Snake)

	files := Schema(nodes, &SchemaConfig{Graphql: true, SingleFile: true}, pkg)

	nodes = engines.Simplify(state, types.DataTypes, GQLTypes, helper.Snake)

	files = append(files, GQL(nodes, pkg, "#")...)

	cwd, _ := os.Getwd()
	// create resolvers
	files = append(files, Resolvers(nodes, pkg)...)

	utils.WriteFiles(files, &cwd)

	err := exec.Command("go", "mod", "init").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("go", "get", "github.com/99designs/gqlgen").Run()
	if err != nil {
		panic(err)
	}

}
