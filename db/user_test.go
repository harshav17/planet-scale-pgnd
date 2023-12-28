package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateUser(tb testing.TB, ctx context.Context, db *DB, u *planetscale.User) (*planetscale.User, context.Context) {
	tb.Helper()

	createUserFunc := func(tx *sql.Tx) error {
		if err := NewUserRepo(db).Create(tx, u); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createUserFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return u, ctx
}

func TestUserRepo_GetUser(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		email := "test@user.com"
		name := "test user"
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   name,
			Email:  email,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := NewUserRepo(db.DB).Get(tx, u.UserID); err != nil {
			t.Fatal(err)
		} else if got.Email != email {
			t.Fatalf("expected title to be %s, got %s", email, got.Email)
		} else if got.Name != name {
			t.Fatalf("expected title to be %s, got %s", name, got.Name)
		}
	})
}
