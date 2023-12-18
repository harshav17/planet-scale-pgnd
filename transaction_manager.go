package planetscale

import (
	"context"
	"database/sql"
)

type TransactionManager interface {
	ExecuteInTx(ctx context.Context, fn func(*sql.Tx) error) error
}
