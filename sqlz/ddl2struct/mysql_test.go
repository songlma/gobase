package ddl2struct

import "testing"

var mysqlDDL = `
CREATE TABLE finance (
	id bigint unsigned NOT NULL AUTO_INCREMENT,
	code varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
	stime datetime NOT NULL,
	open_price float NOT NULL,
	high_price float NOT NULL,
	close_price float NOT NULL,
	low_price float NOT NULL,
	null_c int DEFAULT NULL,
	status int DEFAULT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1746 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
`

func TestMysql2Struct(t *testing.T) {
	re := Mysql2Struct(mysqlDDL)
	t.Logf("%s", re)
}
