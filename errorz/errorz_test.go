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
	e := New(1000, "msg err", WithAlert("错误提示"))
	t.Log(e)                     // [code=1000,msg=msg err,alert=错误提示,ln=errorz_test.go:13]
	t.Log(e.Error())             // [code=1000,msg=msg err,alert=错误提示,ln=errorz_test.go:13]
	t.Log(fmt.Sprintf("%v", e))  //[code=1000,msg=msg err,alert=错误提示,ln=errorz_test.go:13]
	t.Log(fmt.Sprintf("%+v", e)) //[code=1000,msg=msg err,alert=错误提示,ln=errorz_test.go:13]( gobase/errorz.TestNew:errorz/errorz_test.go:13  ->  testing.tRunner:testing/testing.go:1937  ->  runtime.goexit:runtime/asm_arm64.s:1269 )
}

func TestWrap(t *testing.T) {
	e := Wrap(sql.ErrNoRows, 1000, "MsgErr")
	e = Wrap(e, 1001, "BatchErr")
	e = Wrap(e, 1002, "UpdateErr")
	t.Log(e)
	t.Log(e.Error())
	t.Log(fmt.Sprintf("%v", e))
	t.Log(fmt.Sprintf("%+v", e))

}

func TestWrap2(t *testing.T) {
	e := Wrap(sql.ErrNoRows, 1000, "sql err")
	t.Log(e) //{code=1000,msg=sql err,alert=,wrapErr=sql: no rows in result set}
	t.Log(e.Error())
	e = Wrap(e, 1200, "mse err")
	t.Log(e) //{code=1200,msg=mse err,alert=,wrapErr={code=1000,msg=sql err,alert=,wrapErr=sql: no rows in result set}}
	t.Log(e.Error())
}

func Test_Unwrap(t *testing.T) {
	goErr := sql.ErrNoRows
	e := Wrap(goErr, 1000, "sql err")
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
	e := Wrap(goErr, 1000, "sql err")

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

func Call() error {
	return Wrap(SelectSql(), 1200, "messageErr")
}

func SelectSql() error {
	return sql.ErrNoRows
}
