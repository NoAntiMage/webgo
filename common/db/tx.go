package db

import (
	"github.com/jmoiron/sqlx"
)

// NOT USED, impl in "model" package
type TxFunc func(*sqlx.DB) error

func Trans(db *sqlx.DB, fns ...TxFunc) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, fn := range fns {
		err = fn(db)
		if err != nil {
			return
		}
	}
	return nil
}
