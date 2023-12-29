package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpenseGroup(tb testing.TB, tx *sql.Tx, db *DB, u *planetscale.ExpenseGroup) *planetscale.ExpenseGroup {
	tb.Helper()

	if err := NewExpenseGroupRepo(db).Create(tx, u); err != nil {
		tb.Fatal(err)
	}

	return u
}

func TestExpenseGroupRepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("Create Tests", func(t *testing.T) {
		t.Run("invalid user id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			eg := &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  "non-existent-user-id",
			}

			if err := NewExpenseGroupRepo(db.DB).Create(tx, eg); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})

	t.Run("Get Tests", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   "test user",
				Email:  "",
			})

			groupName := "test group"
			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: groupName,
				CreateBy:  u.UserID,
			})

			if got, err := NewExpenseGroupRepo(db.DB).Get(tx, eg.ExpenseGroupID); err != nil {
				t.Fatal(err)
			} else if got.GroupName != groupName {
				t.Fatalf("expected title to be %s, got %s", groupName, got.GroupName)
			}
		})
	})

	t.Run("Update Tests", func(t *testing.T) {
		t.Run("successful update", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   "test user",
				Email:  "",
			})

			groupName := "test group"
			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				ExpenseGroupID: 1,
				GroupName:      groupName,
				CreateBy:       u.UserID,
			})

			update := &planetscale.ExpenseGroupUpdate{
				GroupName: "updated group name",
			}
			if got, err := NewExpenseGroupRepo(db.DB).Update(tx, eg.ExpenseGroupID, update); err != nil {
				t.Fatal(err)
			} else if got.GroupName != update.GroupName {
				t.Fatalf("expected title to be %s, got %s", update.GroupName, got.GroupName)
			}
		})
	})

	t.Run("Delete Tests", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   "test user",
				Email:  "",
			})

			groupName := "test group"
			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				ExpenseGroupID: 1,
				GroupName:      groupName,
				CreateBy:       u.UserID,
			})

			if err := NewExpenseGroupRepo(db.DB).Delete(tx, eg.ExpenseGroupID); err != nil {
				t.Fatal(err)
			}

			// Verify that the group was deleted
			if _, err := NewExpenseGroupRepo(db.DB).Get(tx, eg.ExpenseGroupID); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
}
