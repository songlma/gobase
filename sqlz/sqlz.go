package sqlz

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	sqlxDB   *sqlx.DB
	ctx      context.Context
	dsn      string
	bindType int
}

func Open(ctx context.Context, driver, dsn string) (*DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		sqlxDB:   db,
		ctx:      ctx,
		dsn:      dsn,
		bindType: sqlx.BindType(db.DriverName()),
	}, nil
}

func (this *DB) Conn(ctx context.Context) (*Conn, error) {
	return &Conn{
		sqlxDB:   this.sqlxDB,
		ctx:      ctx,
		bindType: this.bindType,
	}, nil
}
func (this *DB) Close() error {
	return this.sqlxDB.Close()
}
