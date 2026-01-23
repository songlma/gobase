package sqlz

import (
	"context"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var TextCtx = context.Background()

var testdb *DB

var testSqlite3DDL = `
create table finance
(
    id          integer not null
        constraint finance_pk
            primary key autoincrement,
    time        datetime default 0 not null,
    code        text    not null,
    open_price  real     default 0 not null,
    close_price real     default 0 not null,
    high_price  read     default 0 not null,
    low_price   read     default 0 not null
);
`
var testMysqlDDL = `
CREATE TABLE finance (
  id int unsigned NOT NULL AUTO_INCREMENT,
  code varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  time datetime NOT NULL,
  open_price float NOT NULL,
  high_price float NOT NULL,
  close_price float NOT NULL,
  low_price float NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1746 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
`
var testPgDDL = `
create table finance
(
    id          serial
        constraint finance_pk
            primary key,
    code        varchar(16)      not null,
    open_price  double precision not null,
    close_price double precision not null,
    high_price  double precision not null,
    low_price   double precision not null,
    time        timestamp with time zone
);
`

func init() {
	var err error
	//testdb, err = Open(TextCtx, "postgres","postgres://root:root@localhost/finance?sslmode=disable")
	testdb, err = Open(TextCtx, "sqlite3", "hackb.sqlite3?parseTime=true")
	//testdb, err = Open(TextCtx, "mysql", "root:password@tcp(127.0.0.1:3306)/hackb?parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		log.Fatal(err)
	}
	testdb.sqlxDB.SetMaxOpenConns(1)
}
func getConn(t *testing.T, ctx context.Context) *Conn {
	conn, err := testdb.Conn(TextCtx)
	if err != nil {
		t.Errorf("%+v", err)
	}
	return conn
}
