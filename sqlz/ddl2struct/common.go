package ddl2struct

import (
	"go/format"
	"log"
	"strings"
)

var GoNullableType = map[string]string{
	"int8":      "sql.NullInt16",
	"int16":     "sql.NullInt16",
	"int":       "sql.NullInt32",
	"int32":     "sql.NullInt32",
	"int64":     "sql.NullInt64",
	"uint8":     "sql.NullInt16",
	"uint16":    "sql.NullInt32",
	"uint":      "sql.NullInt64",
	"uint32":    "sql.NullInt64",
	"uint64":    "sql.NullInt64",
	"bool":      "sql.NullBool",
	"byte":      "sql.NullByte",
	"[]byte":    "[]byte",
	"string":    "sql.NullString",
	"time.Time": "sql.NullTime",
	"float64":   "sql.NullFloat64",
	"float32":   "sql.NullFloat64",
}

func toFirstLower(table string) string {
	first := strings.ToLower(table[0:1])
	return first + table[1:]
}
func snakeCaseToCamel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

func gofmt(source []byte) []byte {
	result, err := format.Source(source)
	if err != nil {
		log.Fatalln(err)
	}
	return result
}
