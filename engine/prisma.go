package engine

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"vuerd/types"
)

type Prisma struct {
	Provider string
	Nodes    []PNode
}

type PNode struct {
	Name    string
	Fields  []PField
	Options string
}

type PField struct {
	Name    string
	Type    string
	Options string
	DType   string
	Default string
}

func PrismaEngine(Nodes []types.Node) types.File {
	var schema types.File
	prisma := PSimplify(Nodes)
	prisma.Provider = "sqlite"
	file, err := os.ReadFile("templates/prisma.go.tmpl")

	if err != nil {
		log.Println(err.Error())
	}

	t := template.Must(template.New("prisma").Parse(string(file)))

	b, _ := json.Marshal(prisma)

	fmt.Println(string(b))

	err = t.Execute(os.Stdout, prisma)

	if err != nil {
		log.Println(err)
	}

	return schema
}

func PSimplify(nodes []types.Node) Prisma {
	var h Helper
	var prisma Prisma
	for _, node := range nodes {
		var opts []string
		var pfks []string
		var pnode PNode
		pnode.Name = h.Pascal(h.Singular(node.Name))
		for _, f := range node.Fields {
			if !f.Fk {
				var field PField
				field.DType = f.Type
				field.Default = f.Default
				field.Name = h.Camel(f.Name)
				field.Options = Options(f)
				if f.Nullable {
					field.Type = f.Type + "?"
				} else {
					field.Type = f.Type
				}
				if f.Pfk {
					pfks = append(pfks, field.Name)
				}
				pnode.Fields = append(pnode.Fields, field)
			}
		}

		if len(pfks) > 0 {
			opts = append(opts, "@@id(["+strings.Join(pfks, ", "+"])"))
		}

		var flds []PField
		for _, edge := range node.Edges {
			if edge.Direction == "Out" {
				switch edge.Type {
				case "0..1", "1..1":
					{
						flds = append(flds, PField{
							Name: h.Camel(h.Singular(edge.Name)),
							Type: h.Pascal(h.Singular(edge.Name)) + "?",
						})
					}
				case "1..N", "0..N":
					{
						flds = append(flds, PField{
							Name: h.Camel(h.Plural(edge.Name)),
							Type: h.Pascal(h.Singular(edge.Name)) + "[]",
						})
					}

				}

			}

			if edge.Direction == "In" {
				switch edge.Type {
				case "1..N", "1..1":
					{
						flds = append(flds, PField{
							Name: strings.TrimSuffix(edge.Field.Name, "Id"),
							Type: h.Pascal(h.Singular(edge.Name)),
						})
						flds = append(flds, PField{
							Name: edge.Field.Name,
							Type: edge.Field.Type,
						})
					}
				case "0..N", "0..1":
					{
						flds = append(flds, PField{
							Name:    strings.TrimSuffix(edge.Field.Name, "Id"),
							Type:    h.Pascal(h.Singular(edge.Name)) + "?",
							Options: fmt.Sprintf("@relation(fields: [%s], references: [%s])", edge.Field.Name, edge.Reference.Name),
						})

						flds = append(flds, PField{
							Name: edge.Field.Name,
							Type: edge.Field.Type + "?",
						})
					}

				}
			}
		}

		pnode.Fields = append(pnode.Fields, flds...)
		pnode.Options = strings.Join(opts, "\n")
		prisma.Nodes = append(prisma.Nodes, pnode)
	}
	return prisma
}

func Options(f types.Field) string {
	var buffer []string

	if f.Pk {
		buffer = append(buffer, `@id`)
	}

	if f.AutoIncrement {
		buffer = append(buffer, "@default(autoincrement())")
	}

	if f.Name == "created_at" {
		buffer = append(buffer, "@default(now())")
	}

	if f.Name == "updated_at" {
		buffer = append(buffer, "@updatedAt")
	}

	if f.Unique {
		buffer = append(buffer, "@unique()")
	}

	return strings.Join(buffer, " ")
}
