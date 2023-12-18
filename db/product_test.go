package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateProduct(tb testing.TB, ctx context.Context, db *DB, p *planetscale.Product) (*planetscale.Product, context.Context) {
	tb.Helper()

	createProductFunc := func(tx *sql.Tx) error {
		if err := NewProductRepo(db).Create(tx, p); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createProductFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return p, ctx
}

func TestProductRepo_GetProduct(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		p, ctx := MustCreateProduct(t, ctx, db.DB, &planetscale.Product{
			Name:  "test product",
			Price: 100,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := NewProductRepo(db.DB).Get(tx, p.ID); err != nil {
			t.Fatal(err)
		} else if got.Price != 100 {
			t.Fatalf("expected title to be %d, got %f", 100, got.Price)
		}
	})
}
