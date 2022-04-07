package typeorm

var TypeOrmTypes = map[string]string{
	"int":      "number",
	"long":     "number",
	"float":    "float",
	"double":   "number",
	"decimal":  "number",
	"boolean":  "boolean",
	"string":   "string",
	"lob":      "string",
	"json":     "simple-json",
	"date":     "Date",
	"datetime": "Date",
	"time":     "Date",
}
