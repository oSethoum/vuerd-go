package engines

import (
	"strings"
	"vuerd/types"
	"vuerd/utils"

	"github.com/codemodus/kace"
	"github.com/gertd/go-pluralize"
)

type Helper struct{}

var Pluralize = pluralize.NewClient()

func (Helper) Plural(word string) string {
	return Pluralize.Plural(word)
}

// func (Helper) MultiPlural(word string) string {
// 	words := strings.Split(word, "_")
// 	for i := 0; i < len(words); i++ {
// 		words[i] = Pluralize.Plural(words[i])
// 	}
// 	return strings.Join(words, "_")
// }

// func (Helper) MultiSingular(word string) string {
// 	words := strings.Split(word, "_")
// 	for i := 0; i < len(words); i++ {
// 		words[i] = Pluralize.Singular(words[i])
// 	}
// 	return strings.Join(words, "_")
// }

func (Helper) Singular(word string) string {
	return Pluralize.Singular(word)
}

func (Helper) Pascal(word string) string {
	return kace.Pascal(word)
}

func (Helper) Snake(word string) string {
	return kace.Snake(word)
}

func (Helper) Camel(word string) string {
	return kace.Camel(word)
}
func (Helper) CorrectCamel(word string) string {
	return utils.CorrectCase(kace.Camel(word))
}

func (Helper) Kebab(word string) string {
	return kace.Kebab(word)
}

func parseRelationshipType(relationshiptType string) string {
	switch relationshiptType {
	case "ZeroN":
		return "0..N"
	case "OneN":
		return "1..N"
	case "OneOnly":
		return "1..1"
	case "ZeroOne":
		return "0..1"
	}
	return "WTF"
}

func Simplify(state types.State, t map[string]string, m map[string]string, c func(string) string) []types.Node {
	nodes := make([]types.Node, 0)

	for _, table := range state.TableState.Tables {
		node := types.Node{}
		node.Name = table.Name
		node.ID = table.Id
		node.Comment = table.Comment

		for _, column := range table.Columns {
			node.Fields = append(node.Fields, types.Field{
				ID:            column.Id,
				Name:          c(column.Name),
				Comment:       column.Comment,
				Type:          getType(column.DataType, t, m),
				Default:       column.Default,
				Pk:            column.Ui.Pk,
				Fk:            column.Ui.Fk,
				Pfk:           column.Ui.Pfk,
				Unique:        column.Option.Unique,
				Nullable:      !column.Option.NotNull,
				AutoIncrement: column.Option.AutoIncrement,
				Sensitive:     strings.Contains(column.Comment, "-s"),
			})
		}

		for _, relationship := range state.RelationshipState.Relationships {
			if table.Id == relationship.Start.TableId {
				for _, endTable := range state.TableState.Tables {
					if endTable.Id == relationship.End.TableId {
						node.Edges = append(node.Edges, types.Edge{
							ID:        relationship.Id,
							Field:     findColumnById(relationship.Start.ColumnIds[0], table.Columns, t, m, c),
							Reference: findColumnById(relationship.End.ColumnIds[0], endTable.Columns, t, m, c),
							Name:      endTable.Name,
							Type:      parseRelationshipType(relationship.RelationshipType),
							Direction: "Out",
						})
					}
				}

			}

			if table.Id == relationship.End.TableId {
				for _, startTable := range state.TableState.Tables {
					if startTable.Id == relationship.Start.TableId {
						node.Edges = append(node.Edges, types.Edge{
							ID:        relationship.Id,
							Field:     findColumnById(relationship.End.ColumnIds[0], table.Columns, t, m, c),
							Reference: findColumnById(relationship.Start.ColumnIds[0], startTable.Columns, t, m, c),
							Name:      startTable.Name,
							Type:      parseRelationshipType(relationship.RelationshipType),
							Direction: "In",
						})
					}
				}
			}
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func findColumnById(id string, columns []types.Column, t, m map[string]string, c func(string) string) types.Field {
	for _, column := range columns {
		if id == column.Id {
			return types.Field{
				ID:            column.Id,
				Name:          c(column.Name),
				Comment:       column.Comment,
				Type:          getType(column.DataType, t, m),
				Default:       column.Default,
				Pk:            column.Ui.Pk,
				Fk:            column.Ui.Fk,
				Pfk:           column.Ui.Pfk,
				Unique:        column.Option.Unique,
				Nullable:      !column.Option.NotNull,
				AutoIncrement: column.Option.AutoIncrement,
				Sensitive:     strings.Contains(column.Comment, "-s"),
			}
		}
	}

	return types.Field{}
}

func getType(ct string, t, m map[string]string) string {
	if m[t[strings.ToLower(ct)]] == "" {
		return ct
	}
	return m[t[strings.ToLower(ct)]]
}
