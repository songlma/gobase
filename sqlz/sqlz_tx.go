package sqlz

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type TxConn struct {
	sqlxTx   *sqlx.Tx
	ctx      context.Context
	bindType int
	execCnt  int
	checkCnt int
}

type Result struct {
	parent sql.Result
	txConn *TxConn
}

func (this *Conn) Begin(opts *sql.TxOptions) (*TxConn, error) {
	tx, err := this.sqlxDB.BeginTxx(this.ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TxConn{
		sqlxTx:   tx,
		ctx:      this.ctx,
		bindType: this.bindType,
		execCnt:  0,
		checkCnt: 0,
	}, nil
}

func (this Result) LastInsertId() (int64, error) {
	this.txConn.checkCnt++
	return this.parent.LastInsertId()
}
func (this Result) RowsAffected() (int64, error) {
	this.txConn.checkCnt++
	return this.parent.RowsAffected()
}

func (this *TxConn) Rollback() error {
	return this.sqlxTx.Rollback()
}
func (this *TxConn) Commit() error {
	if this.execCnt != this.checkCnt {
		log.Printf("execCnt: %d, checkCnt: %d not eq", this.execCnt, this.checkCnt)
	}
	return this.sqlxTx.Commit()
}

func (this *TxConn) QueryOne(dest interface{}, query string, args ...interface{}) error {
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.sqlxTx.GetContext(this.ctx, dest, query, args...)
}
func (this *TxConn) Query(dest interface{}, query string, args ...interface{}) error {
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.sqlxTx.SelectContext(this.ctx, dest, query, args...)
}
func (this *TxConn) QueryWithIn(dest interface{}, query string, args ...interface{}) error {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.sqlxTx.SelectContext(this.ctx, dest, query, args...)
}

func (this *TxConn) NamedQueryOne(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.sqlxTx.GetContext(this.ctx, dest, query, args...)
}

func (this *TxConn) NamedQuery(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.Query(dest, query, args...)
}
func (this *TxConn) NamedQueryWithIn(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	return this.Query(dest, query, args...)
}

func (this *TxConn) NamedInsert(query string, arg interface{}) (sql.Result, error) {
	this.execCnt++
	var re Result
	sqlRe, err := this.sqlxTx.NamedExecContext(this.ctx, query, arg)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
func (this *TxConn) Update(query string, args ...interface{}) (sql.Result, error) {
	this.execCnt++
	var re Result
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	sqlRe, err := this.sqlxTx.ExecContext(this.ctx, query, args...)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
func (this *TxConn) UpdateWithIn(query string, args ...interface{}) (sql.Result, error) {
	var re Result
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return re, err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	this.execCnt++
	sqlRe, err := this.sqlxTx.ExecContext(this.ctx, query, args...)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
func (this *TxConn) NamedUpdateByStruct(query string, arg interface{}) (sql.Result, error) {
	this.execCnt++
	var re Result
	sqlRe, err := this.sqlxTx.NamedExecContext(this.ctx, query, arg)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
func (this *TxConn) NamedUpdate(query string, args map[string]interface{}) (sql.Result, error) {
	this.execCnt++
	var re Result
	sqlRe, err := this.sqlxTx.NamedExecContext(this.ctx, query, args)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
func (this *TxConn) NamedUpdateWithIn(query string, arg map[string]interface{}) (sql.Result, error) {
	var re Result
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return re, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return re, err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxTx.Rebind(query)
	}
	this.execCnt++
	sqlRe, err := this.sqlxTx.ExecContext(this.ctx, query, args...)
	if err != nil {
		return re, err
	}
	re.parent = sqlRe
	re.txConn = this
	return re, err
}
