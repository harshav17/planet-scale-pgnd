package db

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	db *DB
}

func NewTransactionManager(db *DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm TransactionManager) ExecuteInTx(ctx context.Context, fn func(*sql.Tx) error) error {
	// create a trascation from ctx and execute fn
	tx, err := tm.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			// Ignore commit errors. The tx has already been committed by RELEASE.
			_ = tx.Commit()
		} else {
			// We always need to execute a Rollback() so sql.DB releases the
			// connection.
			_ = tx.Rollback()
		}
	}()

	// execute fn
	err = fn(tx)
	if err != nil {
		return err
	}

	return nil
}
