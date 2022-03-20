package ent

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func countNonKeyFields(node types.Node) int {
	count := 0
	for _, field := range node.Fields {
		if !field.Pk && !field.Fk && !field.Pfk {
			count++
		}
	}
	return count
}

func Schema(nodes []types.Node, config *SchemaConfig) []types.File {
	files := []types.File{}
	schemaBuffer := []string{
		"package schema",
		"",
		"import (",
		"\t\"time\"",
		"\t\"entgo/ent\"",
		"\t\"entgo/ent/dialect/entsql\"",
		"\t\"entgo/ent/schema/field\"",
	}

	if config.Graphql {
		schemaBuffer = append(schemaBuffer, "\t\"entgo.io/contrib/entgql\"")
	}
	schemaBuffer = append(schemaBuffer, "", ")")
	helper := engines.Helper{}

	for _, node := range nodes {
		buffer := []string{}
		pascal := helper.Pascal(helper.Singular(node.Name))

		buffer = append(buffer,
			"package schema",
			"",
			"import (",
			"\t\"time\"",
			"\t\"entgo/ent\"",
			"\t\"entgo/ent/dialect/entsql\"",
			"\t\"entgo/ent/schema/field\"",
			"\t\"entgo/ent/dialect/entsql/edge\"",
			")",
		)

		if config.Graphql {
			buffer = append(buffer, "\t\"entgo.io/contrib/entgql\"")
		}

		if len(node.Edges) > 0 {
			buffer = append(buffer, "\t\"entgo/ent/dialect/entsql/edge\"")
		}
		buffer = append(buffer, ")")

		bodyBuffer := []string{}
		// schema
		bodyBuffer = append(bodyBuffer,
			fmt.Sprintf("// %s Schema", pascal),
			fmt.Sprintf("type %s struct {", pascal),
			"\tent.Schema",
			"}",
			"",
		)

		// annotation
		bodyBuffer = append(bodyBuffer,
			fmt.Sprintf("// %s Annotations", pascal),
			fmt.Sprintf("func (%s) Annotations() []schema.Anotation {", pascal),
			"\treturn []schema.Annotation {",
			fmt.Sprintf("\t\tentsql.Annotation {Table: \"%s\"},", node.Name),
			"\t}",
			"}",
			"",
		)

		if countNonKeyFields(node) > 0 {
			bodyBuffer = append(bodyBuffer,
				fmt.Sprintf("// %s Fields", pascal),
				fmt.Sprintf("func (%s) Fields() []ent.Field {", pascal),
			)
			bodyBuffer = append(bodyBuffer, "\treturn []ent.Field{")

			// fields
			for _, field := range node.Fields {
				options := []string{}

				if field.Sensitive {
					options = append(options, "Sensitive()")
				}

				if field.Unique {
					options = append(options, "Unique()")
				}

				if field.Nullable {
					options = append(options, "Optional()", "Nillable()")
				} else {
					options = append(options, "NotEmpty()")
				}

				if strings.HasPrefix(strings.ToLower(field.Type), "enum") {
					enums := strings.Split(strings.Split(strings.Split(field.Type, "(")[1], ")")[0], ",")
					if config.Graphql {
						namedValues := []string{}
						for _, enum := range enums {
							namedValues = append(namedValues,
								fmt.Sprintf("\t\t\"%s\"", helper.Pascal(enum)),
								", ",
								"\t"+strings.ToUpper(enum),
								",",
								"\n\t\t\t",
							)
						}
						options = append(options, fmt.Sprintf("NamedValues(\n\t\t%s)", strings.Join(namedValues, "")))
					} else {
						options = append(options, fmt.Sprintf("Values(\"%s\")", strings.Join(enums, "\", \"")))
					}
				}

				if len(field.Default) > 0 {
					if field.Type == "UUID" {
						options = append(options, "Default(uuid.New)")
					} else if field.Type == "String" || field.Type == "Enum" {
						options = append(options, fmt.Sprintf("Default(\"%s\")", field.Default))
					} else {
						options = append(options, fmt.Sprintf("Default(%s)", field.Default))
					}
				}

				if field.Pk && field.Type == "UUID" {
					bodyBuffer = append(bodyBuffer, fmt.Sprintf("\t\tfield.%s(\"id\", uuid.UUID{}).%s", field.Type, strings.Join(options, ".\n\t")))
				}

				if config.Graphql {
					options = append(options, fmt.Sprintf("Annotation(entql.OrderField(\"%s\")),", strings.ToUpper(field.Name)))
				}

				if !field.Pk && !field.Fk && !field.Pfk {
					if field.Type == "UUID" {
						bodyBuffer = append(bodyBuffer, fmt.Sprintf("\t\tfield.%s(\"%s\", uuid.UUID{}).%s", field.Type, field.Name, strings.Join(options, ".\n\t\t\t")))
					} else if field.Type == "Json" {
						// bodyBuffer = append(bodyBuffer, fmt.Sprintf("field.%s(\"%s\").%s", field.Type, field.Name, strings.Join(options, ".\n\t")))
					} else if strings.HasPrefix(field.Type, "Enum") {
						bodyBuffer = append(bodyBuffer, fmt.Sprintf("\t\tfield.Enum(\"%s\", uuid.UUID{}).%s", field.Name, strings.Join(options, ".\n\t\t\t")))
					} else {
						bodyBuffer = append(bodyBuffer, fmt.Sprintf("\t\tfield.%s(\"%s\").%s", field.Type, field.Name, strings.Join(options, ".\n\t\t\t")))
					}
				}

				if field.Name == "created_at" {
					options = append(options, "Default(time.Now)", "Immutable()")
				}

				if field.Name == "updated_at" {
					options = append(options, "Default(time.Now)", "UpdateDefault(time.Now)")
				}

			}

			bodyBuffer = append(bodyBuffer, "\t}")
		}
		bodyBuffer = append(bodyBuffer, "}", "")

		// edges
		if len(node.Edges) > 0 {

			bodyBuffer = append(bodyBuffer, fmt.Sprintf("// %s Edges", pascal))
			bodyBuffer = append(bodyBuffer, fmt.Sprintf("func (%s) Edges() []ent.Edge {", pascal))
			bodyBuffer = append(bodyBuffer, "\treturn []ent.Edge{")

			gqlOptions := ""
			if config.Graphql {
				gqlOptions = ".Annotations(entgql.Bind())"
			}

			for _, edge := range node.Edges {
				if edge.Direction == "Out" {
					switch edge.Type {
					case "0..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type)%s,", edge.Name, helper.Singular(helper.Pascal(edge.Name)), gqlOptions),
							)
						}
					case "1..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type)%s.Required(),", edge.Name, helper.Singular(helper.Pascal(edge.Name)), gqlOptions),
							)
						}
					case "0..1", "1..1":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type)%s.Unique(),", helper.Singular(edge.Name), helper.Singular(helper.Pascal(edge.Name)), gqlOptions),
							)
						}
					}
				}
				if edge.Direction == "In" {
					switch edge.Type {
					case "0..N", "1..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type)%s.Ref(\"%s\").Unique(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									gqlOptions,
									node.Name,
								))

						}
					case "1..1":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type)%s.Ref(\"%s\").Unique().Required(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									gqlOptions,
									helper.Singular(node.Name),
								),
							)
						}

					case "0..1":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type)%s.Ref(\"%s\").Unique(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									gqlOptions,
									helper.Singular(node.Name),
								),
							)
						}
					}
				}
			}

			bodyBuffer = append(bodyBuffer, "\t}", "}", "")
		}

		if config.SingleFile {
			schemaBuffer = append(schemaBuffer, bodyBuffer...)
		}

		bodyBuffer = append(buffer, bodyBuffer...)
		files = append(files, types.File{
			Buffer: strings.Join(bodyBuffer, "\n"),
			Path:   fmt.Sprintf("ent/schema/%s.go", helper.Singular(helper.Snake(node.Name))),
		})

	}

	if config.SingleFile {
		files = []types.File{
			{
				Buffer: strings.Join(schemaBuffer, "\n"),
				Path:   "ent/schema/schema.go",
			},
		}
	}

	return files
}

func GQL(nodes []types.Node) []types.File {
	files := []types.File{}
	helper := engines.Helper{}
	files = append(files, types.File{
		Buffer: MutationInput,
		Path:   "ent/template/gql_mutation_input.go.tmpl",
	})

	schemaBuffer := []string{}
	schemaBuffer = append(schemaBuffer,
		"interface Node {",
		"\tid: ID!",
		"}",
		"",
		"scalar Time",
		"",
		"scalar Cursor",
		"",
		"type PageInfo {",
		"\thasNextPage: Boolean!",
		"\thasPreviousPage: Boolean!",
		"\tstartCursor: Cursor",
		"\tendCursor: Cursor",
		"}",
		"",
		"enum OrderDirection {",
		"\tASC",
		"\tDESC",
		"}",
		"",
	)

	queryBuffer := []string{"type Query {"}
	mutationBuffer := []string{"type Mutation {"}
	for _, n := range nodes {
		pascal := helper.Pascal(helper.Singular(n.Name))
		camels := helper.Camel(helper.Plural(n.Name))
		buffer := []string{}
		orderInputsBuffer := []string{}
		orderInputsBuffer = append(orderInputsBuffer,
			fmt.Sprintf("input %sOrder {", pascal),
			"\tdirection: OrderDirection!",
			fmt.Sprintf("\tfield: %sOrderField", pascal),
			"}",
			"",
		)

		queryBuffer = append(queryBuffer, fmt.Sprintf("\t%s(input: %sQueryInput!):%sConnection!", camels, pascal, pascal))
		mutationBuffer = append(mutationBuffer, fmt.Sprintf("\tcreate%s(input:Create%sInput!):%s!", pascal, pascal, pascal))
		mutationBuffer = append(mutationBuffer, fmt.Sprintf("\tupdate%s(id:ID!,input:Update%sInput!):%s!", pascal, pascal, pascal))
		mutationBuffer = append(mutationBuffer, fmt.Sprintf("\tdelete%s(id:ID!):%s!", pascal, pascal))

		orderFieldEnumBuffer := []string{fmt.Sprintf("enum %sOrderField {", pascal)}
		createInputBuffer := []string{fmt.Sprintf("input Create%sInput {", pascal)}
		updateInputBuffer := []string{fmt.Sprintf("input Update%sInput {", pascal)}

		buffer = append(buffer, fmt.Sprintf("type %s implements Node {", pascal))
		buffer = append(buffer, "\tid: ID!")

		for _, f := range n.Fields {
			if !f.Nullable {
				f.Type += "!"
			}

			if !f.Pk && !f.Fk && !f.Pfk {
				if f.Type != "Json" {
					orderFieldEnumBuffer = append(orderFieldEnumBuffer, "\t"+strings.ToUpper(helper.Snake(f.Name)))
				}
				createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%s: %s", helper.Camel(f.Name), f.Type))
				updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\t%s: %s", helper.Camel(f.Name), f.Type))
				buffer = append(buffer, fmt.Sprintf("\t%s: %s", helper.Camel(f.Name), f.Type))
			}
		}

		for _, e := range n.Edges {
			if e.Direction == "Out" {
				switch e.Type {
				case "0..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: [%s!]", helper.Camel(helper.Plural(e.Name)), helper.Pascal(helper.Singular(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sIDs: [ID!]", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tadd%sIDs: [ID!]!", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tremove%sIDs: [ID!]!", helper.Camel(helper.Singular(e.Name))))
					}
				case "1..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: [%s!]!", helper.Camel(helper.Plural(e.Name)), helper.Pascal(helper.Singular(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sIDs: [ID!]!", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tadd%sIDs: [ID!]!", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tremove%sIDs: [ID!]!", helper.Camel(helper.Singular(e.Name))))

					}
				case "1..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: %s!", helper.Singular(helper.Camel(e.Name)), helper.Singular(helper.Pascal(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sID: ID!", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tClear%s: Boolean", helper.Pascal(helper.Singular(e.Name))))
					}
				case "0..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: %s", helper.Singular(helper.Camel(e.Name)), helper.Singular(helper.Pascal(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tClear%s: Boolean", helper.Pascal(helper.Singular(e.Name))))
					}

				}

			} else {
				switch e.Type {
				case "0..N", "0..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: %s", helper.Singular(helper.Camel(e.Name)), helper.Singular(helper.Pascal(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tclear%s: Boolean", helper.Pascal(helper.Singular(e.Name))))
					}
				case "1..N", "1..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s: %s!", helper.Singular(helper.Camel(e.Name)), helper.Singular(helper.Pascal(e.Name))))
						createInputBuffer = append(createInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\t%sID: ID", helper.Camel(helper.Singular(e.Name))))
						updateInputBuffer = append(updateInputBuffer, fmt.Sprintf("\tclear%s: Boolean", helper.Pascal(helper.Singular(e.Name))))
					}
				}

			}
		}
		createInputBuffer = append(createInputBuffer, "}", "")
		updateInputBuffer = append(updateInputBuffer, "}", "")
		orderFieldEnumBuffer = append(orderFieldEnumBuffer, "}", "")

		buffer = append(buffer, "}", "")
		buffer = append(buffer,
			fmt.Sprintf("type %sEdge {", pascal),
			fmt.Sprintf("\tnode: %s", pascal),
			"\tcursor: Cursor!",
			"}",
			"",
			fmt.Sprintf("type %sConnection {", pascal),
			"\ttotalCount: Int!",
			"\tpageInfo: PageInfo!",
			fmt.Sprintf("\tedges: [%sEdge]", pascal),
			"}",
			"",
			fmt.Sprintf("input %sQueryInput {", pascal),
			"\tafter: Cursor",
			"\tfirst: Int",
			"\tbefore: Cursor",
			"\tlast: Int",
			fmt.Sprintf("\torderBy: %sOrder", pascal),
			fmt.Sprintf("\twhere: %sWhereInput", pascal),
			"}",
			"",
		)
		buffer = append(buffer, orderFieldEnumBuffer...)
		buffer = append(buffer, createInputBuffer...)
		buffer = append(buffer, updateInputBuffer...)
		schemaBuffer = append(schemaBuffer, buffer...)
	}

	queryBuffer = append(queryBuffer, "}", "")
	mutationBuffer = append(mutationBuffer, "}", "")

	schemaBuffer = append(schemaBuffer, queryBuffer...)
	schemaBuffer = append(schemaBuffer, mutationBuffer...)

	// handlers
	gqlhandlers := []string{
		"package handlers",
		"",
		"import (",
		"\t\"net/http\"",
		"\t\"time\"",
		"\t\"todo/ent\"",
		"\t\"todo/graph\"",
		"",
		"\t\"entgo.io/contrib/entgql\"",
		"\t\"github.com/99designs/gqlgen/graphql/handler\"",
		"\t\"github.com/99designs/gqlgen/graphql/handler/extension\"",
		"\t\"github.com/99designs/gqlgen/graphql/handler/transport\"",
		"\t\"github.com/99designs/gqlgen/graphql/playground\"",
		"\t\"github.com/gorilla/websocket\"",
		"\t\"github.com/labstack/echo/v4\"",
		")",
		"",
		"func PlaygroundHandler() echo.HandlerFunc {",
		"\th := playground.Handler(\"GraphQL\", \"/query\")",
		"",
		"\treturn func(c echo.Context) error {",
		"\t\th.ServeHTTP(c.Response(), c.Request())",
		"\t\treturn nil",
		"\t}",
		"}",
		"",
		"func PlaygroundWsHandler() echo.HandlerFunc {",
		"\th := playground.Handler(\"GraphQL WS\", \"/subscription\")",
		"\treturn func(c echo.Context) error {",
		"\t\th.ServeHTTP(c.Response(), c.Request())",
		"\t\treturn nil",
		"\t}",
		"}",
		"",
		"func GraphqlHandlers(client *ent.Client) echo.HandlerFunc {",
		"",
		"\th := handler.NewDefaultServer(graph.NewSchema(client))",
		"\th.Use(entgql.Transactioner{TxOpener: client})",
		"\th.AddTransport(transport.POST{})",
		"\th.AddTransport(&transport.Websocket{",
		"\t\tKeepAlivePingInterval: 10 * time.Second,",
		"\t\tUpgrader: websocket.Upgrader{",
		"\t\t\tCheckOrigin: func(r *http.Request) bool {",
		"\t\t\t\treturn true",
		"\t\t\t},",
		"\t\t},",
		"\t})",
		"\th.Use(extension.Introspection{})",
		"\treturn func(c echo.Context) error {",
		"\t\th.ServeHTTP(c.Response(), c.Request())",
		"\t\treturn nil",
		"\t}",
		"}",
		"",
		"func GraphqlHandler(client *ent.Client) echo.HandlerFunc {",
		"\th := handler.NewDefaultServer(graph.NewSchema(client))",
		"\th.Use(entgql.Transactioner{TxOpener: client})",
		"",
		"\treturn func(c echo.Context) error {",
		"\t\th.ServeHTTP(c.Response(), c.Request())",
		"\t\treturn nil",
		"\t}",
		"}",
		"",
		"func GraphqlWsHandler(client *ent.Client) echo.HandlerFunc {",
		"\th := handler.New(graph.NewSchema(client))",
		"\th.AddTransport(transport.POST{})",
		"\th.AddTransport(&transport.Websocket{",
		"\t\tKeepAlivePingInterval: 10 * time.Second,",
		"\t\tUpgrader: websocket.Upgrader{",
		"\t\t\tCheckOrigin: func(r *http.Request) bool {",
		"\t\t\t\treturn true",
		"\t\t\t},",
		"\t\t},",
		"\t})",
		"\th.Use(extension.Introspection{})",
		"\treturn func(c echo.Context) error {",
		"\t\th.ServeHTTP(c.Response(), c.Request())",
		"\t\treturn nil",
		"\t}",
		"}",
		"",
	}

	entcBuffer := []string{
		"//go:build ignore",
		"// +build ignore",
		"",
		"package main",
		"",
		"import (",
		"\t\"log\"",
		"",
		"\t\"entgo.io/contrib/entgql\"",
		"\t\"entgo.io/ent/entc\"",
		"\t\"entgo.io/ent/entc/gen\"",
		")",
		"",
		"func main() {",
		"\tex, err := entgql.NewExtension(",
		"\t\tentgql.WithWhereFilters(true),",
		"\t\tentgql.WithConfigPath(\"../gqlgen.yml\"),",
		"\t\tentgql.WithSchemaPath(\"../graph/ent.graphqls\"),",
		"\t)",
		"\tif err != nil {",
		"\t\tlog.Fatalf(\"creating entgql extension: %v\", err)",
		"\t}",
		"\topts := []entc.Option{",
		"\t\tentc.Extensions(ex),",
		"\t\tentc.TemplateDir(\"./template\"),",
		"\t}",
		"\t",
		"\tif err := entc.Generate(\"./schema\", &gen.Config{}, opts...); err != nil {",
		"\t\tlog.Fatalf(\"running ent codegen: %v\", err)",
		"\t}",
		"}",
	}

	generateBuffer := []string{
		"package ent",
		"",
		"//go:generate go run entc.go && gqlgen ../",
	}

	files = append(files,
		types.File{
			Buffer: strings.Join(schemaBuffer, "\n"),
			Path:   "graph/schema.graphqls",
		},
		types.File{
			Buffer: strings.Join(entcBuffer, "\n"),
			Path:   "ent/entc.go",
		},
		types.File{
			Buffer: strings.Join(generateBuffer, "\n"),
			Path:   "ent/generate.go",
		},
		types.File{
			Buffer: strings.Join(gqlhandlers, "\n"),
			Path:   "handlers/gql.go",
		},
	)

	return files
}
