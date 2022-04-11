package typeorm

import (
	"fmt"
	"strings"
	"vuerd/engines"
	"vuerd/types"
)

func Schema(nodes []types.Node) []types.File {
	files := make([]types.File, 0)

	// base entity
	baseEntityBuffer := make([]string, 0)
	baseEntityBuffer = append(baseEntityBuffer, "import { PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn, DeleteDateColumn } from 'typeorm';\n")
	baseEntityBuffer = append(baseEntityBuffer, "export class BaseEntity {\n\t@PrimaryGeneratedColumn()\n\tid: number;\n\n\t@CreateDateColumn()\n\tcreatedAt: Date;\n\n\t@UpdateDateColumn()\n\tupdatedAt: Date;\n\n\t@DeleteDateColumn()\n\tdeletedAt: Date;\n}")

	files = append(files, types.File{
		Path:   "bases/base.entity.ts",
		Buffer: strings.Join(baseEntityBuffer, "\n"),
	})

	for _, n := range nodes {
		h := engines.Helper{}
		pascal := h.Pascal(h.Singular(n.Name))
		camel := h.Camel(h.Singular(n.Name))
		camels := h.Camel(h.Plural(n.Name))

		importsMap := map[string][]string{}

		importBuffer := make([]string, 0)

		importsMap["../bases/base.entity"] = safeAppend(importsMap["../bases/base.entity"], "BaseEntity")
		importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "Entity")
		importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "PrimaryGeneratedColumn")
		importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "Column")

		entityBuffer := make([]string, 0)
		entityBuffer = append(entityBuffer, fmt.Sprintf("@Entity()\nexport class %s extends BaseEntity {", pascal))

		// fields
		for _, f := range n.Fields {
			// column options
			columnOptionsBuffer := make([]string, 0)
			if f.Nullable {
				columnOptionsBuffer = append(columnOptionsBuffer, "nullable: true")
			}

			if f.Unique {
				columnOptionsBuffer = append(columnOptionsBuffer, "unique: true")
			}

			if f.Default != "" {
				if f.Type == "string" {
					columnOptionsBuffer = append(columnOptionsBuffer, fmt.Sprintf("default: \"%s\"", f.Default))
				} else {
					columnOptionsBuffer = append(columnOptionsBuffer, fmt.Sprintf("default: %s", f.Default))
				}
			}

			columnoptions := ""

			if len(columnOptionsBuffer) > 0 {
				columnoptions = "{ " + strings.Join(columnOptionsBuffer, ", ") + " }"
			}

			// push column to buffer
			if !f.Pk && !f.Pfk && !f.Fk {
				entityBuffer = append(entityBuffer, fmt.Sprintf("\t@Column(%s)\n\t%s: %s;\n", columnoptions, f.Name, f.Type))
			}
		}

		// relationshipts
		for _, e := range n.Edges {

			if e.Direction == "In" {
				switch e.Type {
				case "0..1", "1..1":
					{
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@OneToOne(()=> %s)", h.Pascal(h.Singular(e.Name))))
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@JoinColumn()\n\t%s: %s;\n", h.Camel(e.Name), h.Pascal(h.Singular(e.Name))))

						object := fmt.Sprintf("../%s/%s.entity", h.Camel(h.Plural(e.Name)), h.Camel(h.Singular(e.Name)))
						importsMap[object] = safeAppend(importsMap[object], h.Pascal(h.Singular(e.Name)))

						importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "OneToOne")
						importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "JoinColumn")
					}

				case "0..N", "1..N":
					{
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t@ManyToOne(()=> %s, (%s)=> %s.%s)", h.Pascal(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), h.Camel(h.Singular(e.Name)), camels))
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t%s: %s\n", h.Camel(h.Singular(e.Name)), h.Pascal(h.Singular(e.Name))))

						object := fmt.Sprintf("../%s/%s.entity", h.Camel(h.Plural(e.Name)), h.Camel(h.Singular(e.Name)))
						importsMap[object] = safeAppend(importsMap[object], h.Pascal(h.Singular(e.Name)))

						importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "ManyToOne")
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
						entityBuffer = append(entityBuffer, fmt.Sprintf("\t%s: %s[]\n", h.Camel(h.Plural(e.Name)), h.Pascal(h.Singular(e.Name))))

						object := fmt.Sprintf("../%s/%s.entity", h.Camel(h.Plural(e.Name)), h.Camel(h.Singular(e.Name)))
						importsMap[object] = safeAppend(importsMap[object], h.Pascal(h.Singular(e.Name)))

						importsMap["typeorm"] = safeAppend(importsMap["typeorm"], "OneToMany")
					}
				}
			}

		}

		entityBuffer = append(entityBuffer, "}")
		for k, v := range importsMap {
			importBuffer = append(importBuffer, fmt.Sprintf("import { %s } from '%s';", strings.Join(v, ", "), k))
		}

		importBuffer = append(importBuffer, "")
		entityBuffer = append(importBuffer, strings.Join(entityBuffer, "\n"))

		// append the file
		files = append(files, types.File{
			Path:   fmt.Sprintf("%s/%s.entity.ts", camels, camel),
			Buffer: strings.Join(entityBuffer, "\n"),
		})
	}
	return files
}

func GraphQL(nodes []types.Node) []types.File {
	files := make([]types.File, 0)

	indexExportBuffer := make([]string, 0)

	for _, n := range nodes {
		h := engines.Helper{}
		pascal := h.Pascal(h.Singular(n.Name))
		camel := h.Camel(h.Singular(n.Name))
		camels := h.Camel(h.Plural(n.Name))
		dtoImportsMap := map[string][]string{}

		indexExportBuffer = append(indexExportBuffer, fmt.Sprintf("export * from './%s/%s.module'", camels, camel))

		dtoImportsMap["@nestjs/graphql"] = safeAppend(dtoImportsMap["@nestjs/graphql"], "ObjectType")

		dtoBuffer := make([]string, 0)
		dtoBuffer = append(dtoBuffer, fmt.Sprintf("@ObjectType('%s')\nexport class %sDTO {", pascal, pascal))
		dtoBuffer = append(dtoBuffer, "\n\t@FilterableField()\n\tid: number;\n\n\t@FilterableField()\n\tcreatedAt: Date;\n\n\t@FilterableField()\n\tupdatedAt: Date;\n\n\t@FilterableField()\n\tdeletedAt: Date;\n")

		// fields
		for _, f := range n.Fields {
			field := ""
			options := ""
			optionsBuffer := make([]string, 0)

			if !f.Sensitive {
				field = "@FilterableField"
				dtoImportsMap["@nestjs-query/query-graphql"] = safeAppend(dtoImportsMap["@nestjs-query/query-graphql"], "FilterableField")
			} else {
				field = "@Field"
				dtoImportsMap["@nestjs/graphql"] = safeAppend(dtoImportsMap["@nestjs/graphql"], "Field")
			}

			if f.Nullable {
				optionsBuffer = append(optionsBuffer, "nullable: true")
			}

			if strings.Trim(f.Default, " ") != "" {
				if f.Type == "string" {
					optionsBuffer = append(optionsBuffer, fmt.Sprintf("defaultValue: \"%s\"", f.Default))
				} else {
					optionsBuffer = append(optionsBuffer, "defaultValue:"+f.Default)
				}
			}

			if len(optionsBuffer) > 0 {
				options = "{ " + strings.Join(optionsBuffer, ", ") + " }"
			}

			// push column to buffer
			if !f.Pk && !f.Pfk && !f.Fk {
				dtoBuffer = append(dtoBuffer, fmt.Sprintf("\t%s(%s)\n\t%s: %s;\n", field, options, f.Name, f.Type))
			}
		}

		relationshipBuffer := make([]string, 0)

		// relationshipts
		for _, e := range n.Edges {
			if e.Direction == "In" {
				relationshipBuffer = append(relationshipBuffer, fmt.Sprintf("@Relation('%s',() => %sDTO)", h.Camel(h.Singular(e.Name)), h.Pascal(h.Singular(e.Name))))
				dtoImportsMap["@nestjs-query/query-graphql"] = safeAppend(dtoImportsMap["@nestjs-query/query-graphql"], "Relation")

				object := fmt.Sprintf("../%s/%s.dto", h.Camel(h.Plural(e.Name)), h.Camel(h.Singular(e.Name)))
				dtoImportsMap[object] = safeAppend(dtoImportsMap[object], h.Pascal(h.Singular(e.Name))+"DTO")
			}

			if e.Direction == "Out" {
				switch e.Type {
				case "0..1", "1..1":
					relationshipBuffer = append(relationshipBuffer, fmt.Sprintf("@Relation('%s',() => %sDTO)", h.Camel(h.Singular(e.Name)), h.Pascal(h.Singular(e.Name))))
					dtoImportsMap["@nestjs-query/query-graphql"] = safeAppend(dtoImportsMap["@nestjs-query/query-graphql"], "Relation")

				case "0..N", "1..N":
					relationshipBuffer = append(relationshipBuffer, fmt.Sprintf("@FilterableOffsetConnection('%s', () => %sDTO, { nullable: true })", h.Camel(h.Plural(e.Name)), h.Pascal(h.Singular(e.Name))))
					dtoImportsMap["@nestjs-query/query-graphql"] = safeAppend(dtoImportsMap["@nestjs-query/query-graphql"], "FilterableOffsetConnection")
				}

				object := fmt.Sprintf("../%s/%s.dto", h.Camel(h.Plural(e.Name)), h.Camel(h.Singular(e.Name)))
				dtoImportsMap[object] = safeAppend(dtoImportsMap[object], h.Pascal(h.Singular(e.Name))+"DTO")
			}
		}

		dtoBuffer = append(dtoBuffer, "}")
		importBuffer := make([]string, 0)

		for k, v := range dtoImportsMap {
			importBuffer = append(importBuffer, fmt.Sprintf("import { %s } from '%s';", strings.Join(v, ", "), k))
		}

		importBuffer = append(importBuffer, "")

		dtoBuffer = append(importBuffer, append(relationshipBuffer, strings.Join(dtoBuffer, "\n"))...)

		moduleBuffer := make([]string, 0)
		moduleBuffer = append(moduleBuffer, "import { Module } from '@nestjs/common';")
		moduleBuffer = append(moduleBuffer, "import { NestjsQueryGraphQLModule, PagingStrategies } from '@nestjs-query/query-graphql';")
		moduleBuffer = append(moduleBuffer, "import { NestjsQueryTypeOrmModule } from '@nestjs-query/query-typeorm';")
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("import { %s } from './%s.entity';", pascal, camel))
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("import { %sDTO } from './%s.dto';\n", pascal, camel))
		moduleBuffer = append(moduleBuffer, "@Module({")
		moduleBuffer = append(moduleBuffer, "\timports:[")
		moduleBuffer = append(moduleBuffer, "\t\tNestjsQueryGraphQLModule.forFeature({")
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("\t\t\timports: [NestjsQueryTypeOrmModule.forFeature([%s])],", pascal))
		moduleBuffer = append(moduleBuffer, "\t\t\tresolvers:[")
		moduleBuffer = append(moduleBuffer, "\t\t\t\t{")
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("\t\t\t\t\tDTOClass: %sDTO,", pascal))
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("\t\t\t\t\tEntityClass: %s,", pascal))
		moduleBuffer = append(moduleBuffer, "\t\t\t\t\tenableTotalCount: true,")
		moduleBuffer = append(moduleBuffer, "\t\t\t\t\tpagingStrategy: PagingStrategies.OFFSET,")
		moduleBuffer = append(moduleBuffer, "\t\t\t\t},")
		moduleBuffer = append(moduleBuffer, "\t\t\t],")
		moduleBuffer = append(moduleBuffer, "\t\t}),")
		moduleBuffer = append(moduleBuffer, "\t],")
		moduleBuffer = append(moduleBuffer, "})")
		moduleBuffer = append(moduleBuffer, fmt.Sprintf("export class %sModule {}", pascal))

		files = append(files, types.File{
			Path:   fmt.Sprintf("%s/%s.module.ts", h.Camel(n.Name), h.Camel(h.Singular(n.Name))),
			Buffer: strings.Join(moduleBuffer, "\n"),
		})

		files = append(files, types.File{
			Path:   fmt.Sprintf("%s/%s.dto.ts", h.Camel(n.Name), h.Camel(h.Singular(n.Name))),
			Buffer: strings.Join(dtoBuffer, "\n"),
		})
	}

	files = append(files, types.File{
		Path:   "./index.ts",
		Buffer: strings.Join(indexExportBuffer, "\n"),
	})
	return files
}

func safeAppend(s []string, v string) []string {
	for _, a := range s {
		if a == v {
			return s
		}
	}
	return append(s, v)
}
