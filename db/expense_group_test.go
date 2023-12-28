package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpenseGroup(tb testing.TB, ctx context.Context, db *DB, u *planetscale.ExpenseGroup) (*planetscale.ExpenseGroup, context.Context) {
	tb.Helper()

	createExpenseGroupFunc := func(tx *sql.Tx) error {
		if err := NewExpenseGroupRepo(db).Create(tx, u); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createExpenseGroupFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return u, ctx
}

func TestExpenseGroupRepo_CreateExpenseGroup(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("invalid user id", func(t *testing.T) {
		eg := &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  "non-existent-user-id",
		}

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewExpenseGroupRepo(db.DB).Create(tx, eg); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestExpenseGroupRepo_GetExpenseGroup(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		groupName := "test group"
		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: groupName,
			CreateBy:  u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := NewExpenseGroupRepo(db.DB).Get(tx, eg.ExpenseGroupID); err != nil {
			t.Fatal(err)
		} else if got.GroupName != groupName {
			t.Fatalf("expected title to be %s, got %s", groupName, got.GroupName)
		}
	})
}

func TestExpenseGroupRepo_UpdateExpenseGroup(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		groupName := "test group"
		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			ExpenseGroupID: 1,
			GroupName:      groupName,
			CreateBy:       u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		update := &planetscale.ExpenseGroupUpdate{
			GroupName: "updated group name",
		}
		if got, err := NewExpenseGroupRepo(db.DB).Update(tx, eg.ExpenseGroupID, update); err != nil {
			t.Fatal(err)
		} else if got.GroupName != update.GroupName {
			t.Fatalf("expected title to be %s, got %s", update.GroupName, got.GroupName)
		}
	})
}

func TestExpenseGroupRepo_DeleteExpenseGroup(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		groupName := "test group"
		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			ExpenseGroupID: 1,
			GroupName:      groupName,
			CreateBy:       u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewExpenseGroupRepo(db.DB).Delete(tx, eg.ExpenseGroupID); err != nil {
			t.Fatal(err)
		}

		// Verify that the group was deleted
		if _, err := NewExpenseGroupRepo(db.DB).Get(tx, eg.ExpenseGroupID); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
