package ddl2struct

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/lempiy/Sqlite3CreateTableParser/parser"
)

var SqliteSqlTypeMap = map[string]string{
	//"int":       "int64",
	"char":      "string",
	"text":      "string",
	"clob":      "string",
	"blob":      "[]byte",
	"real":      "float64",
	"double":    "float64",
	"float":     "float64",
	"numeric":   "float64",
	"decimal":   "float64",
	"boolean":   "bool",
	"datetime":  "time.time",
	"timestamp": "time.time",
	"date":      "time.time",
}

func Sqlite2Struct(ddl string) []byte {
	table, errCode := parser.ParseTable(ddl, 0)
	if errCode != parser.ERROR_NONE {
		log.Fatalln("ddl illegle")
	}
	// do stuff with received data
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("type %s struct {\n", table.Name))
	for _, column := range table.Columns {
		columnType := strings.ToLower(column.Type)
		var goType = "string"
		if strings.Contains(columnType, "int") {
			if strings.Contains(columnType, "big") {
				goType = "int64"
			} else {
				goType = "int32"
			}
			if strings.Contains(columnType, "unsigned") {
				goType = "u" + goType
			}
		} else {
			for k, v := range SqliteSqlTypeMap {
				if strings.Contains(columnType, k) {
					goType = v
					break
				}
			}
		}
		if column.IsNotnull == false && column.IsPrimaryKey == false && column.IsAutoincrement == false {
			goType = GoNullableType[goType]
		}
		buff.WriteString(fmt.Sprintf("%s\t%s\t`db:\"%s\"`\n", snakeCaseToCamel(strings.ToLower(column.Name)), goType, column.Name))
	}
	buff.WriteString("}")
	return gofmt(buff.Bytes())
}
