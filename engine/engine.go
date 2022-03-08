package engine

import (
	"strings"
	"vuerd/models"
	"vuerd/types"
)

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

func Simplify(state models.State) []types.Node {
	nodes := make([]types.Node, 0)

	for _, table := range state.TableState.Tables {
		node := types.Node{}
		node.Name = table.Name

		for _, column := range table.Columns {
			node.Fields = append(node.Fields, types.Field{
				Name:          column.Name,
				Comment:       column.Comment,
				Type:          column.DataType,
				Default:       column.Default,
				Pk:            column.Ui.Pk,
				Fk:            column.Ui.Fk,
				Pfk:           column.Ui.Pfk,
				Unique:        column.Option.Unique,
				AutoIncrement: column.Option.AutoIncrement,
				Sensitive:     strings.Contains(column.Comment, "-s"),
			})
		}

		for _, relationship := range state.RelationshipState.Relationships {
			if table.Id == relationship.Start.TableId {
				for _, endTable := range state.TableState.Tables {
					if endTable.Id == relationship.End.TableId {
						node.Edges = append(node.Edges, types.Edge{
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
