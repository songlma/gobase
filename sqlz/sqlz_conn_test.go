package sqlz

import (
	"database/sql"
	"testing"
	"time"
)

func TestConn_QueryOne(t *testing.T) {
	conn := getConn(t, TextCtx)
	var id int
	err := conn.QueryOne(&id, "select id from finance where id=?", 2)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if id != 2 {
		t.Error(id)
		return
	}
}

func TestConn_Query(t *testing.T) {
	conn := getConn(t, TextCtx)
	var ids []int
	err := conn.Query(&ids, "select id from finance where id > ? limit 2", 0)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		t.Error(ids)
		return
	}
}

func TestConn_QueryWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	var ids []int
	err := conn.QueryWithIn(&ids, "select id from finance where id in (?)", []int{1, 2})
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		t.Error(ids)
		return
	}
}

func TestConn_NameQueryOne(t *testing.T) {
	conn := getConn(t, TextCtx)
	args := map[string]interface{}{
		"id": "3",
	}
	var id int
	err := conn.NameQueryOne(&id, "select id from finance where id =:id", args)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if id != 3 {
		t.Error(id)
		return
	}
}

func TestConn_NamedQuery(t *testing.T) {
	conn := getConn(t, TextCtx)
	args := map[string]interface{}{
		"max_id": "4",
		"min_id": "1",
	}
	var ids []int
	err := conn.NamedQuery(&ids, "select id from finance where id >:min_id and id<:max_id", args)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		t.Error(ids)
		return
	}
}

func TestConn_NameQueryWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	args := map[string]interface{}{
		"ids":    []int{2, 4},
		"min_id": 1,
	}
	var ids []int
	err := conn.NamedQueryWithIn(&ids, "select id from finance where id >:min_id and id in (:ids)", args)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		t.Error(ids)
		return
	}
}

func TestConn_Like(t *testing.T) {
	conn := getConn(t, TextCtx)
	var i int
	err := conn.QueryOne(&i, "select count(*) from finance where code like ?", "%ab%")
	if err != nil {
		t.Error(err)
		return
	}
	if i != 3 {
		t.Error(i)
		return
	}
}

func TestConn_Null(t *testing.T) {
	conn := getConn(t, TextCtx)
	var i int = 1
	err := conn.QueryOne(&i, "select count(*) from finance where code is null")
	if err != nil {
		t.Error(err)
		return
	}
	if i != 0 {
		t.Error(i)
		return
	}
}
func TestConn_Time(t *testing.T) {
	conn := getConn(t, TextCtx)
	now := time.Now()
	tx, err := conn.Begin(nil)
	var old time.Time
	err = tx.QueryOne(&old, Strings("select time from finance where id=?"), 1)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	re, err := tx.NamedUpdate("update finance set time=:time,open_price=open_price+1 where id=:id", map[string]interface{}{
		"id":   1,
		"time": now,
	})
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	if af != 1 {
		tx.Rollback()
		t.Error(err)
		return
	}
	var qt time.Time
	err = tx.QueryOne(&qt, Strings("select time from finance where id=?"), 1)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	if now.Unix() > qt.Unix() {
		t.Error(old.Unix(), now.Unix(), qt.Unix())
	}
	tx.Commit()
}

func TestSelectNull(t *testing.T) {
	conn := getConn(t, TextCtx)
	var i sql.NullInt64
	err := conn.QueryOne(&i, "select null_c from finance where id=1")
	if err != nil {
		t.Error(err)
		return
	}
	if !i.Valid {
		t.Error(i)
		return
	}
	if i.Int64 != 1 {
		t.Error(i)
		return
	}
}
