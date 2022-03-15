package ent

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Schema(nodes []types.Node, config *SchemaConfig) []types.File {
	files := []types.File{}
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
		)

		if config.Graphql {
			buffer = append(buffer, "\t\"entgo.io/contrib/entgql\"")
		}

		if len(node.Edges) > 0 {
			buffer = append(buffer, "\t\"entgo/ent/dialect/entsql/edge\"")
		}
		buffer = append(buffer, ")")

		// schema
		buffer = append(buffer,
			fmt.Sprintf("// %s Schema", pascal),
			fmt.Sprintf("type %s struct {", pascal),
			"\tent.Schema",
			"}",
			"",
		)

		// annotation
		buffer = append(buffer,
			fmt.Sprintf("// %s Annotations", pascal),
			fmt.Sprintf("func (%s) Annotations() []schema.Anotation {", pascal),
			"\treturn []schema.Annotation {",
			fmt.Sprintf("\t\tentsql.Annotation {Table: \"%s\"},", node.Name),
			"\t}",
			"}",
			"",
		)

		if countNonKeyFields(node) > 0 {
			buffer = append(buffer,
				fmt.Sprintf("// %s Fields", pascal),
				fmt.Sprintf("func (%s) Fields() []ent.Field{", pascal),
			)
			buffer = append(buffer, "\treturn []ent.Field{")

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
				}

				if strings.HasPrefix(strings.ToLower(field.Type), "enum") {
					enums := strings.Split(strings.Split(strings.Split(field.Type, "(")[1], ")")[0], ",")
					if config.Graphql {
						namedValues := []string{}
						for _, enum := range enums {
							namedValues = append(namedValues,
								fmt.Sprintf("\"%s\"", helper.Pascal(enum)),
								fmt.Sprintf("\"%s\"", strings.ToUpper(helper.Pascal(enum))),
								"\n",
							)
						}
						options = append(options, fmt.Sprintf("NamedValues(%s)", strings.Join(namedValues, "")))
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

				if field.Name == "created_at" {
					options = append(options, "Default(time.Now)", "Immutable()")
				}

				if field.Name == "updated_at" {
					options = append(options, "Default(time.Now)", "UpdateDefault(time.Now)")
				}

				if field.Pk && field.Type == "UUID" {
					buffer = append(buffer, fmt.Sprintf("\t\tfield.%s(\"id\", uuid.UUID{}).%s", field.Type, strings.Join(options, ".\n\t")))
				}

				if config.Graphql {
					options = append(options, fmt.Sprintf("Annotation(entql.OrderField(\"%s\")),", strings.ToUpper(field.Name)))
				}

				if !field.Pk && !field.Fk && !field.Pfk {
					if field.Type == "UUID" {
						buffer = append(buffer, fmt.Sprintf("\t\tfield.%s(\"%s\", uuid.UUID{}).%s", field.Type, field.Name, strings.Join(options, ".\n\t\t\t")))
					} else if field.Type == "Json" {
						// buffer = append(buffer, fmt.Sprintf("field.%s(\"%s\").%s", field.Type, field.Name, strings.Join(options, ".\n\t")))
					} else {
						buffer = append(buffer, fmt.Sprintf("\t\tfield.%s(\"%s\").%s", field.Type, field.Name, strings.Join(options, ".\n\t\t\t")))
					}
				}

			}

			buffer = append(buffer, "\t}")
		}
		buffer = append(buffer, "}", "")

		// edges
		if len(node.Edges) > 0 {

			buffer = append(buffer, fmt.Sprintf("// %s Edges", pascal))
			buffer = append(buffer, fmt.Sprintf("func (%s) Edges() []ent.Edge {", pascal))
			buffer = append(buffer, "\treturn []ent.Edge{")

			for _, edge := range node.Edges {
				if edge.Direction == "Out" {
					switch edge.Type {
					case "0..N":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type),", edge.Name, helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					case "1..N":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type).Required(),", edge.Name, helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					case "0..1", "1..1":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type).Unique(),", helper.Singular(edge.Name), helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					}
				}
				if edge.Direction == "In" {
					switch edge.Type {
					case "0..N", "1..N":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type).Ref(\"%s\").Unique(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									node.Name,
								))

						}
					case "1..1":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type).Ref(\"%s\").Unique().Required(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									helper.Singular(node.Name),
								),
							)
						}

					case "0..1":
						{
							buffer = append(buffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type).Ref(\"%s\").Unique(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									helper.Singular(node.Name),
								),
							)
						}
					}
				}
			}

			buffer = append(buffer, "\t}", "}", "")
		}

		files = append(files, types.File{
			Buffer: strings.Join(buffer, "\n"),
			Path:   fmt.Sprintf("ent/schema/%s.go", helper.Singular(helper.Snake(node.Name))),
		})

		// if config.Graphql {
		// 	files = append(files, types.File{
		// 		Buffer: strings.Join(gqlBuffer, "\n"),
		// 		Path:   path.Join(config.GraphqlFolder, fmt.Sprintf("/%s.graphqls", helper.Singular(helper.Snake(node.Name)))),
		// 	})
		// }
	}

	return files
}

func GQL(nodes []types.Node) {

}

func countNonKeyFields(node types.Node) int {
	count := 0
	for _, field := range node.Fields {
		if !field.Pk && !field.Fk && !field.Pfk {
			count++
		}
	}
	return count
}
