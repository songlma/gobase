package sqlz

import (
	"testing"
	"time"
)

type testModel struct {
	Id    int       `json:"-" db:"id"`
	Name  string    `json:"-" db:"name"`
	age   uint      `json:"-" db:"age"`
	ctime time.Time `json:"-" db:"ctime"`
}

// ( id, name, age, ctime) values
// (:id,:name,:age,:ctime)
func TestTemplateSql(t *testing.T) {
	t.Log(TemplateSql(testModel{}))
}

func TestStrings(t *testing.T) {
	sqlStr := Strings("insert into test_model",
		"( id, name, age, ctime) values",
		"(:id,:name,:age,:ctime)")

	if sqlStr != "insert into test_model ( id, name, age, ctime) values (:id,:name,:age,:ctime)" {
		t.Error(sqlStr)
	}
}
