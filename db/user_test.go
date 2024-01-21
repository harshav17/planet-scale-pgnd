package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateUser(tb testing.TB, tx *sql.Tx, db *DB, u *planetscale.User) *planetscale.User {
	tb.Helper()

	if err := NewUserRepo(db).Create(tx, u); err != nil {
		tb.Fatal(err)
	}

	return u
}

func TestUserRepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("Get Tests", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			email := "test@user.com"
			name := "test user"
			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   name,
				Email:  email,
			})
			if got, err := NewUserRepo(db.DB).Get(tx, u.UserID); err != nil {
				t.Fatal(err)
			} else if got.Email != email {
				t.Fatalf("expected title to be %s, got %s", email, got.Email)
			} else if got.Name != name {
				t.Fatalf("expected title to be %s, got %s", name, got.Name)
			}
		})
	})

	t.Run("Create Tests", func(t *testing.T) {
		t.Run("successful create", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			email := "test@user.com"
			name := "test user"
			u := &planetscale.User{
				UserID: "test-user-id",
				Name:   name,
				Email:  email,
			}
			if err := NewUserRepo(db.DB).Create(tx, u); err != nil {
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
	})

	t.Run("Upsert Tests", func(t *testing.T) {
		t.Run("successful upsert", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			email := "test@user.com"
			name := "test user"
			u := &planetscale.User{
				UserID: "test-user-id",
				Name:   name,
				Email:  email,
			}
			if err := NewUserRepo(db.DB).Upsert(tx, u); err != nil {
				t.Fatal(err)
			}
			if got, err := NewUserRepo(db.DB).Get(tx, u.UserID); err != nil {
				t.Fatal(err)
			} else if got.Email != email {
				t.Fatalf("expected title to be %s, got %s", email, got.Email)
			}

			// update email
			email = "test2@user.com"
			u.Email = email
			if err := NewUserRepo(db.DB).Upsert(tx, u); err != nil {
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
	})
}
