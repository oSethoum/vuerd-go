package iris

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Main(mod string) types.File {
	buffer := []string{}

	buffer = append(buffer,
		"package main",
		"",
		"import (",
		fmt.Sprintf("\t\"%s/db\"", mod),
		fmt.Sprintf("\t\"%s/routes\"", mod),
		fmt.Sprintf("\t\"%s/services\"", mod),
		"",
		"\t\"github.com/kataras/iris/v12\"",
		")",
		"",
		"func main () {",
		"\tapp := iris.Default()",
		"\tdb.Init()",
		"\tservices.Init()",
		"\troutes.Init(app)",
		"\tapp.Listen(\":3000\")",
		"}",
	)
	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "main.go",
	}
}

func Routes(nodes []types.Node) types.File {
	helper := engines.Helper{}
	buffer := []string{}
	buffer = append(buffer,
		"package routes",
		"",
		"func Init(app *iris.Application) {",
		"\t app.Party(\"/api\")",
		"",
	)
	for _, node := range nodes {
		kebabs := helper.Kebab(helper.Plural(node.Name))
		camels := helper.Camel(helper.Plural(node.Name))
		pascal := helper.Pascal(helper.Singular(node.Name))
		pascals := helper.Pascal(helper.Plural(node.Name))

		buffer = append(buffer,
			fmt.Sprintf("%s := api.Party(\"/%s\")", camels, kebabs),
			fmt.Sprintf("%s.Get(\"/{id}\", handlers.Get%s)", camels, pascal),
			fmt.Sprintf("%s.Get(\"/\", handlers.Get%s)", camels, pascals),
			fmt.Sprintf("%s.Post(\"/\", handlers.Post%s)", camels, pascal),
			fmt.Sprintf("%s.Post(\"/bulk\", handlers.Post%s)", camels, pascals),
			fmt.Sprintf("%s.Put(\"/\", handlers.Put%s)", camels, pascal),
			fmt.Sprintf("%s.Put(\"/bulk\", handlers.Put%s)", camels, pascals),
			fmt.Sprintf("%s.Delete(\"/\", handlers.Delete%s)", camels, pascal),
			fmt.Sprintf("%s.Delete(\"/bulk\", handlers.Delete%s)", camels, pascals),
		)
	}

	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "routes/routes.go",
	}
}

func Handlers(nodes []types.Node) {
	// handlersBuffer := []string{}

	// for _, node := range nodes {
	// 	buffer := []string{}
	// }
}
