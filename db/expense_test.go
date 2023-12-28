package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpense(tb testing.TB, ctx context.Context, db *DB, u *planetscale.Expense) (*planetscale.Expense, context.Context) {
	tb.Helper()

	createExpenseFunc := func(tx *sql.Tx) error {
		if err := NewExpenseRepo(db).Create(tx, u); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createExpenseFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return u, ctx
}

func TestExpenseRepo_CreateExpense(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("invalid user id", func(t *testing.T) {
		e := &planetscale.Expense{
			GroupID:     1,
			PaidBy:      "non-existent-user-id",
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Time{},
			CreatedBy:   "test-user-id",
			UpdatedBy:   "test-user-id",
		}

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		// TODO add better error
		if err := NewExpenseRepo(db.DB).Create(tx, e); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestExpenseRepo_GetExpense(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})

		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     eg.ExpenseGroupID,
			PaidBy:      u.UserID,
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Now(),
			CreatedBy:   u.UserID,
			UpdatedBy:   u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		got, err := NewExpenseRepo(db.DB).Get(tx, e.ExpenseID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ExpenseID != e.ExpenseID {
			t.Fatalf("expected expense id %d, got %d", e.ExpenseID, got.ExpenseID)
		}
	})
}

func TestExpenseRepo_UpdateExpense(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})

		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     eg.ExpenseGroupID,
			PaidBy:      u.UserID,
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

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := NewExpenseRepo(db.DB).Update(tx, e.ExpenseID, update); err != nil {
			t.Fatal(err)
		} else if got.Amount != *update.Amount {
			t.Fatalf("expected amount to be %f, got %f", *update.Amount, got.Amount)
		}
	})
}

func TestExpenseRepo_DeleteExpense(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		eg, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})

		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     eg.ExpenseGroupID,
			PaidBy:      u.UserID,
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Now(),
			CreatedBy:   u.UserID,
			UpdatedBy:   u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewExpenseRepo(db.DB).Delete(tx, e.ExpenseID); err != nil {
			t.Fatal(err)
		}

		_, err = NewExpenseRepo(db.DB).Get(tx, e.ExpenseID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
