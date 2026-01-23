package sqlz

import (
	"testing"
	"time"
)

func TestTxConn_Rollback(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	var time1 time.Time
	err = tx.QueryOne(&time1, "select time from finance where id=?", 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	re, err := tx.Update("update finance set time=?,open_price=open_price+1 where id =?", time.Now(), 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 1 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Rollback()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	conn2 := getConn(t, TextCtx)
	var time2 time.Time
	err = conn2.QueryOne(&time2, "select time from finance where id=?", 2)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if time1.Unix() != time2.Unix() {
		t.Error("rollback fail")
		return
	}
}

func TestTxConn_Commit(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	var time1 time.Time
	err = tx.QueryOne(&time1, "select time from finance where id=?", 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	re, err := tx.Update("update finance set time=?,open_price=open_price+1 where id =?", time.Now(), 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 1 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	conn2 := getConn(t, TextCtx)
	var time2 time.Time
	err = conn2.QueryOne(&time2, "select time from finance where id=?", 2)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	if time1.Unix() > time2.Unix() {
		t.Error("rollback fail")
		return
	}
}
func TestTxConn_QueryOne(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	var id int
	err = tx.QueryOne(&id, "select id from finance where id=?", 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if id != 2 {
		tx.Rollback()
		t.Error(id)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_Query(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	var ids []int
	err = tx.Query(&ids, "select id from finance where id>? limit ?", 3, 2)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		tx.Rollback()
		t.Error(ids)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}
func TestTxConn_QueryWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	var ids []int
	err = tx.QueryWithIn(&ids, "select id from finance where id in ( ? )", []int{1, 2})
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		tx.Rollback()
		t.Error(ids)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_NamedQueryOne(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	args := map[string]interface{}{
		"id": "3",
	}
	var id int
	err = tx.NamedQueryOne(&id, "select id from finance where id =:id", args)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if id != 3 {
		tx.Rollback()
		t.Error(id)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_NamedQuery(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
	}
	args := map[string]interface{}{
		"max_id": "4",
		"min_id": "1",
	}
	var ids []int
	err = tx.NamedQuery(&ids, "select id from finance where id >:min_id and id<:max_id", args)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		tx.Rollback()
		t.Error(ids)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_NamedQueryWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	args := map[string]interface{}{
		"ids":    []int{2, 4},
		"min_id": 1,
	}
	var ids []int
	err = tx.NamedQueryWithIn(&ids, "select id from finance where id > :min_id and id in (:ids)", args)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if len(ids) != 2 {
		tx.Rollback()
		t.Error(ids)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_Update(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	re, err := tx.Update("update finance set time=?,open_price=open_price+1 where id = ?", time.Now(), 1)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 1 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_NamedUpdateWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
	}
	args := map[string]interface{}{
		"time": time.Now(),
		"ids":  []int{1, 2, 3},
	}
	re, err := tx.NamedUpdateWithIn("update finance set time= :time ,open_price=open_price+1 where id in ( :ids )", args)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 3 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_NamedInsert(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
	}
	arg := map[string]interface{}{
		"time":        time.Now(),
		"open_price":  123,
		"close_price": 123,
		"high_price":  123,
		"low_price":   123,
		"code":        "1231",
	}
	re, err := tx.NamedInsert("insert into finance (time,code,open_price,close_price,high_price,low_price) values (:time,:code,:open_price,:close_price,:high_price,:low_price);", arg)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	li, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
	if li == 0 {
		tx.Rollback()
		t.Error(li)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTxConn_NamedUpdate1(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
	}
	args := map[string]interface{}{
		"time": time.Now(),
		"id":   1,
	}
	re, err := tx.NamedUpdate("update finance set time=:time,open_price=open_price+1 where id =:id", args)
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 1 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestTxConn_UpdateWithIn(t *testing.T) {
	conn := getConn(t, TextCtx)
	tx, err := conn.Begin(nil)
	if err != nil {
		t.Error(err)
		return
	}
	re, err := tx.UpdateWithIn("update finance set time=?,open_price=open_price+1 where id in (?)", time.Now(), []int{1, 2, 3})
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	af, err := re.RowsAffected()
	if err != nil {
		tx.Rollback()
		t.Errorf("%+v", err)
		return
	}
	if af != 3 {
		tx.Rollback()
		t.Error(af)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
}
