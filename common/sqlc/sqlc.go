package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goweb/common/db"
	"goweb/common/logx"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
)

//SqlBuilder.Result() type
type SqlResultFn func() (string, []any, error)

type SqlConn struct {
	*sqlx.DB
}

//TODO session, db connect pool impl
func NewSqlConnect() *SqlConn {
	return &SqlConn{
		DB: db.GetDb(),
	}
}

func (m *SqlConn) ExecContextWithLog(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := m.ExecContext(ctx, query, args...)
	logx.Loggerx.Debugf("STMT: %v\nPARAMS: %v\n", query, args)
	return result, err
}

func (m *SqlConn) BeginTxxWithLog(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := m.BeginTxx(ctx, &sql.TxOptions{Isolation: 0, ReadOnly: false})
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

func (m *SqlConn) Trans(ctx context.Context, fn func(tx *Tx) error) (err error) {
	tx, err := m.BeginTxxWithLog(ctx, &sql.TxOptions{Isolation: 0, ReadOnly: false})
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			txErr := tx.Rollback()
			panicResult := fmt.Sprintln(p)
			stackBytes := debug.Stack()
			err = errors.New(txErr.Error() + panicResult + string(stackBytes))
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)

}

type Tx struct {
	*sqlx.Tx
}

func (tx *Tx) PrepareContextWithLog(ctx context.Context, query string) (*Stmt, error) {
	logx.Loggerx.Infof("PREPARE: %v", query)
	stmt, err := tx.PrepareContext(ctx, query)
	return &Stmt{
		Stmt: stmt,
	}, err
}

type Stmt struct {
	*sql.Stmt
}

func (stmt *Stmt) ExecContextWithLog(ctx context.Context, args ...any) (sql.Result, error) {
	logx.Loggerx.Infof("PARAMS: %v", args)
	return stmt.ExecContext(ctx, args...)
}
