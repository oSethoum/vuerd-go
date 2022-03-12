package gorm

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Schema(nodes []types.Node) types.File {

	helper := engines.Helper{}
	buffer := []string{}
	var idType string

	for _, node := range nodes {

		buffer = append(buffer, fmt.Sprintf("type %s struct {", helper.Pascal(helper.Singular(node.Name))))
		buffer = append(buffer, "\tModel")
		for _, field := range node.Fields {
			var options []string

			if idType == "" && field.Pk {
				idType = field.Type
			}

			if field.Unique {
				options = append(options, "unique")
			}

			if !field.Nullable {
				options = append(options, "not null")
			}

			if len(field.Default) > 0 {
				options = append(options, fmt.Sprintf("default:%s", field.Default))
			}

			if !field.Pk && !field.Fk && field.Name != "CreatedAt" && field.Name != "UpdatedAt" {
				if len(options) > 0 {
					buffer = append(buffer, fmt.Sprintf("\t%s %s `gorm:\"%s\" json:\"%s\"`", field.Name, field.Type, strings.Join(options, ","), helper.CorrectCamel(field.Name)))
				} else {
					buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`", field.Name, field.Type, helper.CorrectCamel(field.Name)))
				}
			}
		}

		for _, edge := range node.Edges {
			if edge.Direction == "Out" {
				switch edge.Type {
				case "0..N", "1..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s []%s `json:\"%s\"`",
							helper.Pascal(edge.Name),
							helper.Singular(helper.Pascal(edge.Name)),
							helper.Plural(helper.CorrectCamel(edge.Name)),
						))
					}

				case "0..1", "1..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`",
							helper.Singular(helper.Pascal(edge.Name)),
							helper.Singular(helper.Pascal(edge.Name)),
							helper.Singular(helper.CorrectCamel(edge.Name)),
						))
					}
				}

			}

			if edge.Direction == "In" {
				switch edge.Type {
				case "0..N", "1..N":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`",
							helper.Pascal(helper.Singular(edge.Field.Name)),
							edge.Field.Type,
							helper.CorrectCamel(edge.Field.Name),
						))
					}
				case "0..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`",
							helper.Pascal(edge.Field.Name),
							edge.Field.Type,
							helper.CorrectCamel(edge.Field.Name),
						))

						buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`",
							strings.TrimSuffix(helper.Pascal(edge.Field.Name), "ID"),
							helper.Pascal(helper.Singular(edge.Name)),
							strings.TrimSuffix(helper.Camel(edge.Field.Name), "ID"),
						))
					}
				case "1..1":
					{
						buffer = append(buffer, fmt.Sprintf("\t%s %s `json:\"%s\"`",
							helper.Pascal(edge.Field.Name),
							edge.Field.Type,
							helper.CorrectCamel(edge.Field.Name),
						))
					}
				}
			}
			buffer = append(buffer, "}", "")
		}
	}

	base := []string{
		"type Model struct {",
		"\tID " + idType + " `gorm:\"primaryKey\" json:\"id\"`",
		"\tCreatedAt time.Time `json:\"createdAt\"`",
		"\tUpdatedAt time.Time `json:\"updatedAt\"`",
		"\tDeletedAt gorm.DeletedAt `gorm:\"index\" json:\"deletedAt\"`",
		"",
	}

	buffer = append(base, buffer...)

	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "models/models.go",
	}
}

func Migration(nodes []types.Node) {
	type edge struct {
		first string
		last  string
	}
	helper := engines.Helper{}
	models := []string{}
	edges := []edge{}
	for _, node := range nodes {
		models = append(models, helper.Pascal(helper.Singular(node.Name)))
		for _, e := range node.Edges {
			if e.Direction == "Out" {
				edges = append(edges, edge{first: helper.Pascal(helper.Singular(node.Name)), last: helper.Pascal(helper.Singular(node.Name))})
			}
		}
	}

	for _, e := range edges {
		swap(&models, indexOf(models, e.first), indexOf(models, e.last))
	}

	fmt.Printf("%v", models)
}

func swap(models *[]string, first, last int) {

}

func indexOf(models []string, model string) int {
	for i := 0; i < len(models); i++ {
		if models[i] == model {
			return i
		}
	}
	return -1
}
