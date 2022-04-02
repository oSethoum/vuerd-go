package ent

var EntTypes = map[string]string{
	"int":      "Int",
	"long":     "Int64",
	"float":    "Float",
	"uuid":     "UUID",
	"double":   "Float64",
	"decimal":  "Int",
	"boolean":  "Bool",
	"string":   "String",
	"lob":      "String",
	"date":     "Time",
	"json":     "Json",
	"datetime": "Time",
	"time":     "Time",
}

var GQLTypes = map[string]string{
	"int":      "Int",
	"long":     "Int",
	"float":    "Float",
	"uuid":     "ID",
	"double":   "Float",
	"decimal":  "Int",
	"boolean":  "Bool",
	"string":   "String",
	"lob":      "String",
	"date":     "Time",
	"json":     "Json",
	"datetime": "Time",
	"time":     "Time",
}

type SchemaConfig struct {
	Graphql    bool
	SingleFile bool
}
