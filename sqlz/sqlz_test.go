package sqlz

import (
	"testing"
)

func TestOpen(t *testing.T) {
	mydb1, err := Open(TextCtx, "sqlite3", "hackb.sqlite3")
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	err = mydb1.sqlxDB.Ping()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestPGOpen(t *testing.T) {
	mydb1, err := Open(TextCtx, "postgres", "")
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	defer mydb1.sqlxDB.Close()
	err = mydb1.sqlxDB.Ping()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestConn(t *testing.T) {
	mydb1, err := Open(TextCtx, "sqlite3", "hackb.sqlite3")
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	defer mydb1.sqlxDB.Close()
	conn, err := mydb1.Conn(TextCtx)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	err = conn.sqlxDB.PingContext(TextCtx)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn(t *testing.T) {
	mydb1, err := Open(TextCtx, "sqlite3", "hackb.sqlite3")
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	defer mydb1.sqlxDB.Close()
	conn, err := mydb1.Conn(TextCtx)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	var i int
	err = tx.sqlxTx.Get(&i, "select 1", nil)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestDB_Close(t *testing.T) {
	mydb1, err := Open(TextCtx, "sqlite3", "hackb.sqlite3")
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	err = mydb1.Close()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}
