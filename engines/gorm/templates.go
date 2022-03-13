package gorm

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

var idType string

func Schema(nodes []types.Node) types.File {

	helper := engines.Helper{}
	buffer := []string{}

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
		"package model",
		"",
		"import (",
		"\t\"gorm.io/gorm\"",
		"\t\"time\"",
		")",
		"",
		"type Model struct {",
		"\tID " + idType + " `gorm:\"primaryKey\" json:\"id\"`",
		"\tCreatedAt time.Time `json:\"createdAt\"`",
		"\tUpdatedAt time.Time `json:\"updatedAt\"`",
		"\tDeletedAt gorm.DeletedAt `gorm:\"index\" json:\"deletedAt\"`",
		"}",
		"",
	}

	buffer = append(base, buffer...)

	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "models/models.go",
	}
}

func DB(nodes []types.Node, driver, mod string) types.File {
	buffer := []string{}
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
		models = swap(models, indexOf(models, e.first), indexOf(models, e.last))
	}

	buffer = append(buffer,
		"package db",
		"",
		"import (",
		fmt.Sprintf("\t\"%s/models\"", mod),
		"\t\"log\"",
		"\t\"os\"",
		"",
		"\t\"gorm.io/driver/sqlite\"",
		"\t\"gorm.io/gorm\"",
		"\t\"gorm.io/gorm/logger\"",
		")",
		"",
		"var DB *gorm.DB",
		"",
		"func Init() error {",
		fmt.Sprintf("\tdns := \"%s\"", "dev.sqlite"),
		fmt.Sprintf("\tdb, err := gorm.Open(%s.Open(dns), &gorm.Config{", driver),
		"\t\tLogger: logger.New(log.New(os.Stdout, \"\\r\\n\", log.LstdFlags), logger.Config{",
		"\t\t\tLogLevel: logger.Info,",
		"\t\t\tColorful: true,",
		"\t\t}),",
		"\t})",
		"\tif err != nil {",
		"\t\treturn err",
		"\t}",
		"\tdb.AutoMigrate(",
	)

	for _, model := range models {
		buffer = append(buffer, fmt.Sprintf("\t\t&models.%s{},", model))
	}

	buffer = append(buffer,
		"\t)",
		"",
		"\tDB = db",
		"\treturn nil",
		"}",
	)

	return types.File{
		Buffer: strings.Join(buffer, "\n"),
		Path:   "db/db.go",
	}
}

func swap(models []string, first, last int) []string {
	tmp := models[first]
	models[first] = models[last]
	models[last] = tmp
	return models
}

func indexOf(models []string, model string) int {
	for i := 0; i < len(models); i++ {
		if models[i] == model {
			return i
		}
	}
	return -1
}

func Services(nodes []types.Node, mod string) []types.File {
	files := []types.File{}
	helper := engines.Helper{}
	serviceBuffer := []string{}
	serviceBuffer = append(serviceBuffer,
		"package services",
		"",
		"import (",
		fmt.Sprintf("\"%s/db\"", mod),
		"\"gorm.io/gorm\"",
		")",
		"",
		"var DB *gorm.DB",
		"",
		"func Init() {",
		"\tDB = db.DB",
		"}",
	)

	files = append(files, types.File{
		Buffer: strings.Join(serviceBuffer, "\n"),
		Path:   "services/services.go",
	})

	for _, node := range nodes {
		buffer := []string{
			"package services",
			"",
			"import (",
			fmt.Sprintf("\t\"%s/models\"", mod),
			"\t\"errors\"",
			")",
			"",
		}
		camel := helper.Camel(helper.Singular(node.Name))
		camels := helper.Camel(helper.Plural(node.Name))
		pascal := helper.Pascal(helper.Singular(node.Name))
		pascals := helper.Pascal(helper.Plural(node.Name))
		buffer = append(buffer,
			fmt.Sprintf("func Find%s(%s *models.%s) (*models.%s, error) {", pascal, camel, pascal, pascal),
			fmt.Sprintf("\tresult := DB.First(%s, %s.ID)", camel, camel),
			fmt.Sprintf("\treturn %s, result.Error", camel),
			"}",
			"",
			fmt.Sprintf("func Find%s() (*[]*models.%s, error) {", pascals, pascal),
			fmt.Sprintf("\t%s := new([]*models.%s)", camels, pascal),
			fmt.Sprintf("\tresult := DB.Find(%s)", camels),
			fmt.Sprintf("\treturn %s, result.Error", camels),
			"}",
			"",
			fmt.Sprintf("func Create%s(%s *models.%s) (*models.%s, error) {", pascal, camel, pascal, pascal),
			fmt.Sprintf("\tresult := DB.Create(%s)", camel),
			fmt.Sprintf("\treturn %s, result.Error", camel),
			"}",
			"",
			fmt.Sprintf("func Create%s(%s *[]*models.%s) (*[]*models.%s, error) {", pascals, camels, pascal, pascal),
			fmt.Sprintf("\tresult := DB.Create(%s)", camels),
			fmt.Sprintf("\treturn %s, result.Error", camels),
			"}",
			"",
			fmt.Sprintf("func Update%s(%s *models.%s) (*models.%s, error) {", pascal, camel, pascal, pascal),
			fmt.Sprintf("\tresult := DB.First(%s.ID)", camel),
			"\tif result.Error != nil {",
			"\t\treturn nil, result.Error",
			"\t}",
			fmt.Sprintf("\tresult = DB.Save(%s)", camel),
			fmt.Sprintf("\treturn %s, result.Error", camel),
			"}",
			"",
			fmt.Sprintf("func Update%s(%s *[]*models.%s) (*[]*models.%s, error) {", pascals, camels, pascal, pascal),
			fmt.Sprintf("\tIDs := make([]%s, 0)", idType),
			fmt.Sprintf("\tfor _, %s := range *%s {", camel, camels),
			fmt.Sprintf("\t	IDs = append(IDs, %s.ID)", camel),
			"\t}",
			fmt.Sprintf("\tres := new([]*models.%s)", pascal),
			"\tDB.Find(res, IDs)",
			fmt.Sprintf("\tif len(*res) != len(*%s) {", camels),
			"\t	return nil, errors.New(\"error in IDs\")",
			"\t}",
			fmt.Sprintf("\tresult := DB.Save(%s)", camels),
			fmt.Sprintf("\treturn %s, result.Error", camels),
			"}",
			"",
			fmt.Sprintf("func Delete%s(%s *models.%s) (*models.%s, error) {", pascal, camel, pascal, pascal),
			fmt.Sprintf("\tresult := DB.Delete(%s)", camel),
			fmt.Sprintf("\treturn %s, result.Error", camel),
			"}",
			"",
			fmt.Sprintf("func Delete%s(%s *[]*models.%s) (*[]*models.%s, error) {", pascals, camels, pascal, pascal),
			fmt.Sprintf("\tresult := DB.Delete(%s)", camels),
			fmt.Sprintf("\treturn %s, result.Error", camels),
			"}",
		)
		files = append(files, types.File{
			Buffer: strings.Join(buffer, "\n"),
			Path:   fmt.Sprintf("services/%s.go", camels),
		})
	}

	return files
}
