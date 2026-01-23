package sqlz

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Conn struct {
	sqlxDB   *sqlx.DB
	ctx      context.Context
	bindType int
}

func (this *Conn) QueryOne(dest interface{}, query string, args ...interface{}) error {
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.GetContext(this.ctx, dest, query, args...)
}

func (this *Conn) Query(dest interface{}, query string, args ...interface{}) error {
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.SelectContext(this.ctx, dest, query, args...)
}
func (this *Conn) QueryWithIn(dest interface{}, query string, args ...interface{}) error {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.SelectContext(this.ctx, dest, query, args...)
}

func (this *Conn) NameQueryOne(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.GetContext(this.ctx, dest, query, args...)
}

func (this *Conn) NamedQuery(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.SelectContext(this.ctx, dest, query, args...)
}
func (this *Conn) NamedQueryWithIn(dest interface{}, query string, arg interface{}) error {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	if this.bindType != sqlx.QUESTION {
		query = this.sqlxDB.Rebind(query)
	}
	return this.sqlxDB.SelectContext(this.ctx, dest, query, args...)
}
