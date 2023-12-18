package db_mock

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	ExecuteInTxFn func(ctx context.Context, fn func(*sql.Tx) error) error
}

func (t TransactionManager) ExecuteInTx(ctx context.Context, fn func(*sql.Tx) error) error {
	return t.ExecuteInTxFn(ctx, fn)
}
