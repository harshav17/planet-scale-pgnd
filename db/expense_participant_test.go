package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpenseParticipant(tb testing.TB, ctx context.Context, db *DB, u *planetscale.ExpenseParticipant) (*planetscale.ExpenseParticipant, context.Context) {
	tb.Helper()

	createExpenseParticipantFunc := func(tx *sql.Tx) error {
		if err := NewExpenseParticipantRepo(db).Create(tx, u); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createExpenseParticipantFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return u, ctx
}

func TestExpenseParticipantRepo_CreateExpenseParticipant(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("invalid user id", func(t *testing.T) {
		ep := &planetscale.ExpenseParticipant{
			ExpenseID:       1,
			UserID:          "non-existent-user-id",
			AmountOwed:      100,
			SharePercentage: 100,
			SplitMethod:     "EQUAL",
			Note:            "test expense",
		}

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewExpenseParticipantRepo(db.DB).Create(tx, ep); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestExpenseParticipantRepo_GetExpenseParticipant(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})
		g, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})
		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     g.ExpenseGroupID,
			PaidBy:      u.UserID,
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Now(),
			CreatedBy:   u.UserID,
			UpdatedBy:   u.UserID,
		})
		ep, ctx := MustCreateExpenseParticipant(t, ctx, db.DB, &planetscale.ExpenseParticipant{
			ExpenseID:       e.ExpenseID,
			UserID:          u.UserID,
			AmountOwed:      100,
			SharePercentage: 100,
			SplitMethod:     "EQUAL",
			Note:            "test expense",
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		epRepo := NewExpenseParticipantRepo(db.DB)
		ep2, err := epRepo.Get(tx, ep.ExpenseID, ep.UserID)
		if err != nil {
			t.Fatal(err)
		}
		if ep2.ExpenseID != ep.ExpenseID {
			t.Fatalf("expected expense id %d, got %d", ep.ExpenseID, ep2.ExpenseID)
		}
		if ep2.UserID != ep.UserID {
			t.Fatalf("expected user id %s, got %s", ep.UserID, ep2.UserID)
		}
	})

	t.Run("invalid expense id", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		epRepo := NewExpenseParticipantRepo(db.DB)
		if _, err := epRepo.Get(tx, 1, u.UserID); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestExpenseParticipantRepo_UpdateExpenseParticipant(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})
		g, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})
		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     g.ExpenseGroupID,
			PaidBy:      u.UserID,
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Now(),
			CreatedBy:   u.UserID,
			UpdatedBy:   u.UserID,
		})
		ep, ctx := MustCreateExpenseParticipant(t, ctx, db.DB, &planetscale.ExpenseParticipant{
			ExpenseID:       e.ExpenseID,
			UserID:          u.UserID,
			AmountOwed:      100,
			SharePercentage: 100,
			SplitMethod:     "EQUAL",
			Note:            "test expense",
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		amountOwed := 200.0
		sharePercentage := 50.0
		splitMethod := "PERCENTAGE"
		note := "updated note"
		update := &planetscale.ExpenseParticipantUpdate{
			AmountOwed:      &amountOwed,
			SharePercentage: &sharePercentage,
			SplitMethod:     &splitMethod,
			Note:            &note,
		}
		if got, err := NewExpenseParticipantRepo(db.DB).Update(tx, ep.ExpenseID, ep.UserID, update); err != nil {
			t.Fatal(err)
		} else if got.ExpenseID != ep.ExpenseID {
			t.Fatalf("expected expense id %d, got %d", ep.ExpenseID, got.ExpenseID)
		} else if got.UserID != ep.UserID {
			t.Fatalf("expected user id %s, got %s", ep.UserID, got.UserID)
		} else if got.AmountOwed != amountOwed {
			t.Fatalf("expected amount owed %f, got %f", amountOwed, got.AmountOwed)
		} else if got.SharePercentage != sharePercentage {
			t.Fatalf("expected share percentage %f, got %f", sharePercentage, got.SharePercentage)
		} else if got.SplitMethod != splitMethod {
			t.Fatalf("expected split method %s, got %s", splitMethod, got.SplitMethod)
		}
	})
}

func TestExpenseParticipantRepo_DeleteExpenseParticipant(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})
		g, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})
		e, ctx := MustCreateExpense(t, ctx, db.DB, &planetscale.Expense{
			GroupID:     g.ExpenseGroupID,
			PaidBy:      u.UserID,
			Amount:      100,
			Description: "test expense",
			Timestamp:   time.Now(),
			CreatedBy:   u.UserID,
			UpdatedBy:   u.UserID,
		})
		ep, ctx := MustCreateExpenseParticipant(t, ctx, db.DB, &planetscale.ExpenseParticipant{
			ExpenseID:       e.ExpenseID,
			UserID:          u.UserID,
			AmountOwed:      100,
			SharePercentage: 100,
			SplitMethod:     "EQUAL",
			Note:            "test expense",
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewExpenseParticipantRepo(db.DB).Delete(tx, ep.ExpenseID, ep.UserID); err != nil {
			t.Fatal(err)
		}

		if _, err := NewExpenseParticipantRepo(db.DB).Get(tx, ep.ExpenseID, ep.UserID); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
