package gorm

var GormTypes = map[string]string{
	"int":      "int",
	"long":     "int64",
	"float":    "float64",
	"double":   "float64",
	"decimal":  "int",
	"boolean":  "bool",
	"string":   "string",
	"lob":      "string",
	"json":     "string",
	"date":     "time.Time",
	"datetime": "time.Time",
	"time":     "time.Time",
}
