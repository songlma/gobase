package ddl2struct

import (
	"testing"
)

var ddl = `create table finance1
(
    id          integer constraint finance_pk primary key autoincrement,
    stime       datetime  not null,
    code        text    not null,
    DECIMAL_c  DECIMAL(10,5)     not null,
    close_price VARYING CHARACTER(255)     not null,
    VARCHAR_c  VARCHAR(255)  ,
    bigint_c  bigint not null  ,
    blob_c  blob ,
    timestamp_c  timestamp ,
    UNSIGNED_c   UNSIGNED BIG INT     not null
);`

func TestSqlite2Struct(t *testing.T) {
	re := Sqlite2Struct(ddl)
	t.Logf("%s", re)
}
