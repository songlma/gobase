package sqlz

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func Strings(querys ...string) string {
	return strings.Join(querys, " ")
}

var INSERT_TPL = `INSERT:
insert into XXXXX
( %s) values
(:%s)
`
var SELECT_TPL = `SELECT *:
select %s
from XXXXX
where id=0
`
var UPDATE_TPL = `UPDATE ALL:
update XXXXX set %s
where id=0
`

func TemplateSql(model interface{}) string {
	ref := reflect.ValueOf(model)
	refType := ref.Type()
	var cols []string
	var sets []string
	for i := 0; i < refType.NumField(); i++ {
		col := refType.Field(i).Tag.Get("db")
		cols = append(cols, col)
		if col == "id" || col == "ctime" || col == "utime" {
			continue
		}
		sets = append(sets, fmt.Sprintf("%s=:%s", col, col))
	}
	return fmt.Sprintf("\n生成时间:%s\n", time.Now()) + fmt.Sprintf(INSERT_TPL, strings.Join(cols, ", "), strings.Join(cols, ",:")) +
		fmt.Sprintf(SELECT_TPL, strings.Join(cols, ", ")) +
		fmt.Sprintf(UPDATE_TPL, strings.Join(sets, ", "))
}
