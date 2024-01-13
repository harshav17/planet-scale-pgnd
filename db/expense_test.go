package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpense(tb testing.TB, tx *sql.Tx, db *DB, u *planetscale.Expense) *planetscale.Expense {
	tb.Helper()

	if err := NewExpenseRepo(db).Create(tx, u); err != nil {
		tb.Fatal(err)
	}
	return u
}

func TestExpenseRepo_All(t *testing.T) {
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

			e := &planetscale.Expense{
				GroupID:     1,
				PaidBy:      "non-existent-user-id",
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   "test-user-id",
				UpdatedBy:   "test-user-id",
			}
			// TODO add better error
			if err := NewExpenseRepo(db.DB).Create(tx, e); err == nil {
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

			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     eg.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})

			got, err := NewExpenseRepo(db.DB).Get(tx, e.ExpenseID)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseID != e.ExpenseID {
				t.Fatalf("expected expense id %d, got %d", e.ExpenseID, got.ExpenseID)
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

			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     eg.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})

			newAmount := 200.0
			update := &planetscale.ExpenseUpdate{
				Amount: &newAmount,
			}

			if got, err := NewExpenseRepo(db.DB).Update(tx, e.ExpenseID, update); err != nil {
				t.Fatal(err)
			} else if got.Amount != *update.Amount {
				t.Fatalf("expected amount to be %f, got %f", *update.Amount, got.Amount)
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

			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     eg.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})

			if err := NewExpenseRepo(db.DB).Delete(tx, e.ExpenseID); err != nil {
				t.Fatal(err)
			}

			_, err = NewExpenseRepo(db.DB).Get(tx, e.ExpenseID)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})

	t.Run("Find Tests", func(t *testing.T) {
		t.Run("successful find", func(t *testing.T) {
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

			eg := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     eg.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})

			filter := planetscale.ExpenseFilter{
				GroupID: eg.ExpenseGroupID,
			}

			if got, err := NewExpenseRepo(db.DB).Find(tx, filter); err != nil {
				t.Fatal(err)
			} else if len(got) != 1 {
				t.Fatalf("expected 1 expense, got %d", len(got))
			} else if got[0].ExpenseID != e.ExpenseID {
				t.Fatalf("expected expense id %d, got %d", e.ExpenseID, got[0].ExpenseID)
			}
		})
	})
}
