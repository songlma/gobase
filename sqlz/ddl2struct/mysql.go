package ddl2struct

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/xwb1989/sqlparser"
)

var MysqlSqlTypeMap = map[string]string{
	"tinyint":            "int8",
	"smallint":           "int16",
	"mediumint":          "int32",
	"int":                "int32",
	"integer":            "int32",
	"bigint":             "int64",
	"tinyint unsigned":   "uint8",
	"smallint unsigned":  "uint16",
	"mediumint unsigned": "uint32",
	"int unsigned":       "uint32",
	"integer unsigned":   "uint32",
	"bigint unsigned":    "uint64",
	"bit":                "byte",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "[]byte",
	"tinyblob":           "[]byte",
	"mediumblob":         "[]byte",
	"longblob":           "[]byte",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "[]byte",
	"varbinary":          "[]byte",
}

func Mysql2Struct(ddl string) []byte {
	stmt, err := sqlparser.Parse(ddl)
	if err != nil {
		log.Fatalln("ddl illegle")
	}
	ddlStmt, ok := stmt.(*sqlparser.DDL)
	if !ok {
		log.Fatalln("not create ddl")
	}
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("type %s struct {\n", snakeCaseToCamel(ddlStmt.NewName.Name.String())))
	for _, c := range (*ddlStmt).TableSpec.Columns {
		column := *c
		var typeKey = column.Type.Type
		switch column.Type.Type {
		case "int", "integer", "tinyint", "smallint", "mediumint", "bigint":
			if column.Type.Unsigned {
				typeKey += " unsigned"
			}
		}
		goType, ok := MysqlSqlTypeMap[typeKey]
		if ok != true {
			log.Fatalln("no support type", typeKey)
		}
		if column.Type.NotNull == false && column.Type.Autoincrement == false {
			goType = GoNullableType[goType]
		}
		buff.WriteString(fmt.Sprintf("%s\t%s\t`db:\"%s\"`\n", snakeCaseToCamel(strings.ToLower(column.Name.String())), goType, column.Name.String()))
	}
	buff.WriteString(fmt.Sprintf("}"))
	return gofmt(buff.Bytes())
}
