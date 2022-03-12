package prisma

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Schema(nodes []types.Node, provider string) types.File {
	var buffer []string
	var helper engines.Helper

	buffer = append(buffer,
		"database db {",
		"\turl = env(\"DATABASE_URL\")",
		fmt.Sprintf("\tprovider = \"%s\"", provider),
		"}",
		"",
		"generator client {",
		"\tprovider = \"prisma-client-js\"",
		"}",
		"",
	)

	for _, node := range nodes {

		buffer = append(buffer,
			fmt.Sprintf("model %s {", helper.Pascal(helper.Singular(node.Name))),
		)

		var pfks []string
		for _, field := range node.Fields {
			options := []string{}

			if field.Pfk {
				pfks = append(pfks, helper.Camel(field.Name))
			}

			if field.Nullable {
				field.Type += "?"
			}

			if field.Pk {
				options = append(options, "@id")
			}

			if field.AutoIncrement {
				options = append(options, "@default(@autoincrement())")
			} else if field.Name == "createdAt" {
				options = append(options, "@default(now())")
			} else {
				if len(field.Default) > 0 {
					if field.Type == "String" {
						options = append(options, fmt.Sprintf(`@default("%s")`, field.Default))
					} else {
						options = append(options, fmt.Sprintf(`@default(%s)`, field.Default))
					}
				}
			}

			if field.Unique {
				options = append(options, "@unique")
			}

			if field.Name == "updatedAt" {
				options = append(options, "@updatedAt")
			}

			if !field.Fk {
				buffer = append(buffer, fmt.Sprintf("\t%s %s %s",
					helper.Camel(field.Name),
					field.Type,
					strings.Join(options, " ")))
			}
		}

		for _, edge := range node.Edges {
			if edge.Direction == "Out" {
				switch edge.Type {
				case "0..1", "1..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s",
							helper.Camel(helper.Plural(edge.Name)),
							helper.Pascal(helper.Singular(edge.Name))+"?",
						))
					}
				case "0..N", "1..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s",
							helper.Camel(helper.Plural(edge.Name)),
							helper.Pascal(helper.Singular(edge.Name))+"[]",
						))
					}
				}
			}

			if edge.Direction == "In" {
				switch edge.Type {
				case "1..1", "1..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s @relation(fields:[%s], references:[%s])",
							strings.TrimSuffix(edge.Field.Name, "Id"),
							helper.Pascal(helper.Singular(edge.Name)),
							edge.Field.Name,
							edge.Reference.Name,
						))

						buffer = append(buffer, fmt.Sprintf("\t%s %s",
							edge.Field.Name,
							edge.Field.Type,
						))
					}
				case "0..1", "0..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s @relation(fields:[%s], references:[%s])",
							strings.TrimSuffix(edge.Field.Name, "Id"),
							helper.Pascal(helper.Singular(edge.Name)),
							edge.Field.Name,
							edge.Reference.Name,
						))

						buffer = append(buffer, fmt.Sprintf("\t%s %s",
							edge.Field.Name,
							edge.Field.Type+"?",
						))
					}
				}
			}
		}

		if len(pfks) > 0 {
			buffer = append(buffer, fmt.Sprintf("\t@@id([%s])", strings.Join(pfks, ", ")))
		}
		buffer = append(buffer, "}", "")
	}

	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "prisma/schema.prisma",
	}
}
