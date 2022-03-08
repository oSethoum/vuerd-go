package engine

import (
	"os"
	"strings"
	"text/template"
	"vuerd/types"
)

type Helper struct {
}

func (Helper) Upper(word string) string {
	return strings.ToUpper(word)
}

func (Helper) Pascal(word string) string {
	return strings.ToUpper(word)
}

type Prisma struct {
	Helper
	Driver string
	Nodes  []types.Node
}

func PrismaEngine(Nodes []types.Node) types.File {
	var schema types.File
	var prisma Prisma
	prisma.Driver = "sqlite"

	t := template.Must(template.New("prisma").Parse(`{{.Upper .Driver}}`))

	err := t.Execute(os.Stdout, prisma)

	if err != nil {
		panic(err)
	}

	return schema
}
