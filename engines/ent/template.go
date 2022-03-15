package ent

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

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
								", ",
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

				if field.Name == "created_at" {
					options = append(options, "Default(time.Now)", "Immutable()")
				}

				if field.Name == "updated_at" {
					options = append(options, "Default(time.Now)", "UpdateDefault(time.Now)")
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

			}

			bodyBuffer = append(bodyBuffer, "\t}")
		}
		bodyBuffer = append(bodyBuffer, "}", "")

		// edges
		if len(node.Edges) > 0 {

			bodyBuffer = append(bodyBuffer, fmt.Sprintf("// %s Edges", pascal))
			bodyBuffer = append(bodyBuffer, fmt.Sprintf("func (%s) Edges() []ent.Edge {", pascal))
			bodyBuffer = append(bodyBuffer, "\treturn []ent.Edge{")

			for _, edge := range node.Edges {
				if edge.Direction == "Out" {
					switch edge.Type {
					case "0..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type),", edge.Name, helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					case "1..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type).Required(),", edge.Name, helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					case "0..1", "1..1":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.To(\"%s\", %s.Type).Unique(),", helper.Singular(edge.Name), helper.Singular(helper.Pascal(edge.Name))),
							)
						}
					}
				}
				if edge.Direction == "In" {
					switch edge.Type {
					case "0..N", "1..N":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type).Ref(\"%s\").Unique(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									node.Name,
								))

						}
					case "1..1":
						{
							bodyBuffer = append(bodyBuffer,
								fmt.Sprintf("\t\tedge.From(\"%s\",%s.Type).Ref(\"%s\").Unique().Required(),",
									helper.Singular(edge.Name),
									helper.Singular(helper.Pascal(edge.Name)),
									helper.Singular(node.Name),
								),
							)
						}

					case "0..1":
						{
							bodyBuffer = append(bodyBuffer,
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
				Path:   "ent/schema.go",
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
		Path:   "ent/template/mutation_input.tmpl",
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

	for _, node := range nodes {
		buffer := []string{}
		pascal := helper.Pascal(helper.Singular(node.Name))

		buffer = append(buffer,
			fmt.Sprintf("type %s implements Node {", pascal))

	}

	return files
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
