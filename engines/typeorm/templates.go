package typeorm

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Schema(nodes []types.Node) []types.File {
	files := make([]types.File, 0)
	for _, n := range nodes {
		h := engines.Helper{}
		pascal := h.Pascal(h.Singular(n.Name))
		camel := h.Camel(h.Singular(n.Name))
		camels := h.Camel(h.Plural(n.Name))

		importBuffer := make([]string, 0)
		importBuffer = append(importBuffer, "import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn, UpdateDateColumn, DeleteDateColumn")

		entityBuffer := make([]string, 0)
		entityBuffer = append(entityBuffer, fmt.Sprintf("@Entity()\nexport class %s {", pascal))
		entityBuffer = append(entityBuffer, "\t@PrimaryGeneratedColumn()\n\tid: number;"+"\n")

		// fields
		for _, f := range n.Fields {
			if !f.Pk && !f.Pfk && !f.Fk {
				entityBuffer = append(entityBuffer, fmt.Sprintf("\t@Column()\n\t%s: %s;\n", f.Name, f.Type))
			}
		}

		// relationshipts
		for _, e := range n.Edges {

			if e.Direction == "In" {
				switch e.Type {
				case "0..1", "1..1":
					{
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@OneTnOne(()=> %s)", pascal))
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@JoinColumn()\n\t%s: %s;\n", h.Camel(e.Name), pascal))

						importBuffer = append(importBuffer, " ,OneToOne, JoinColumn")
					}

				case "0..N", "1..N":
					{
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@ManyToOne(()=> %s, (%s)=> %s.%s)", h.Pascal(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), camels))
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t%s: %s", h.Camel(h.Singular(e.Name)), h.Pascal(h.Singular(e.Name))))

						importBuffer = append(importBuffer, " ,ManyToOne")
					}
				}
			}

			if e.Direction == "Out" {
				switch e.Type {
				case "0..N", "1..N":
					{
						entityBuffer = append(
							entityBuffer,
							fmt.Sprintf("\t@OneToMany(()=> %s, (%s)=> %s.%s)", h.Pascal(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), camel),
						)
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t%s: %s[]", h.Camel(h.Plural(e.Name)), h.Pascal(h.Singular(e.Name))))

						importBuffer = append(importBuffer, ", OneToMany")
					}
				}
			}

		}

		importBuffer = append(importBuffer, " } from 'typeorm';\n")

		entityBuffer = append(entityBuffer, "}")
		entityBuffer = append(importBuffer, entityBuffer...)

		// create and append the file
		files = append(files, types.File{
			Path:   fmt.Sprintf("graphql/%s/%s.entity.ts", h.Camel(n.Name), h.Camel(h.Singular(n.Name))),
			Buffer: strings.Join(entityBuffer, "\n"),
		})
	}
	return files
}

func Dtos(nodes []types.Node) []types.File {
	return nil
}
