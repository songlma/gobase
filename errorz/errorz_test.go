package errorz

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	e := New(1000, "msg err")
	t.Log(e)                    //{code=1000,msg=msg err}
	t.Log(e.Error())            //code:1000,msg err-wrapErr:[]-errorz_test.go:12
	t.Log(fmt.Sprintf("%v", e)) //{code=1000,msg=msg err}
	t.Log(fmt.Sprintf("%+v", e))
}

func TestWrap(t *testing.T) {
	e := Wrap(1000, "MsgErr", sql.ErrNoRows)
	e = Wrap(1001, "BatchErr", e)
	t.Log(e)         //{code=1000,msg=MSG ERR,alert=,wrapErr=sql: no rows in result set}
	t.Log(e.Error()) //MSG ERR-errorz.go:76
	t.Log(fmt.Sprintf("%v", e))
	t.Log(fmt.Sprintf("%+v", e))
}

func TestWrap2(t *testing.T) {
	e := Wrap(1000, "sql err", sql.ErrNoRows)
	t.Log(e) //{code=1000,msg=sql err,alert=,wrapErr=sql: no rows in result set}
	t.Log(e.Error())
	e = Wrap(1200, "mse err", e)
	t.Log(e) //{code=1200,msg=mse err,alert=,wrapErr={code=1000,msg=sql err,alert=,wrapErr=sql: no rows in result set}}
	t.Log(e.Error())
}

func Test_Unwrap(t *testing.T) {
	goErr := sql.ErrNoRows
	e := Wrap(1000, "sql err", goErr)
	t.Log(e)
	t.Log(errors.Unwrap(e))
	if reflect.TypeOf(errors.Unwrap(e)).String() != "*errors.errorString" {
		t.Error(reflect.TypeOf(errors.Unwrap(e)).String())
	}
}

func TestErrorStatus(t *testing.T) {
	e := New(1000, "message")
	var xbErr Error
	if !errors.As(e, &xbErr) {
		t.Error(e)
		return
	}
}

func TestIs(t *testing.T) {
	goErr := sql.ErrNoRows
	e := Wrap(1000, "sql err", goErr)

	if !errors.Is(errors.Unwrap(e), goErr) {
		t.Error(errors.Unwrap(e))
		t.Error(goErr)
	}
	if !errors.Is(e, goErr) {
		t.Error(e)
		t.Error(goErr)
	}

}

func TestAs(t *testing.T) {
	//err := New(1000, "message err")
	//var xbErr Error
	//if !errors.As(err, &xbErr) {
	//	t.Log(err)
	//}
	sTime := time.Now().UnixNano()
	time.Sleep(time.Second)
	ts := float64(time.Now().UnixNano()-sTime) / 1000000
	t.Log(ts)
}

func TestStack(t *testing.T) {
	err := Call()
	t.Log(fmt.Sprintf("%+v", err))
}

func TestAlert(t *testing.T) {
	e := New(1000, "msg err", "alert msg")
	t.Log(e)
	t.Log(e.Error())
	if e.Alert() != "alert msg" {
		t.Error("alert msg error")
	}

	e2 := Wrap(1001, "wrap err", e, "wrap alert")
	t.Log(e2)
	if e2.Alert() != "wrap alert" {
		t.Error("wrap alert error")
	}

	e3 := NewAlertError(1002, "new alert error", "alert msg 3")
	if e3.Alert() != "alert msg 3" {
		t.Error("new alert error msg error")
	}

	e4 := WrapWithAlert(e, "wrap with alert")
	if e4.Alert() != "wrap with alert" {
		t.Error("wrap with alert error")
	}
}

func Call() error {
	return Wrap(1200, "messageErr", SelectSql())
}

func SelectSql() error {
	return sql.ErrNoRows
}
