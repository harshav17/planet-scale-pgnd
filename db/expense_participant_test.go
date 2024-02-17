package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateExpenseParticipant(tb testing.TB, tx *sql.Tx, db *DB, u *planetscale.ExpenseParticipant) *planetscale.ExpenseParticipant {
	tb.Helper()

	if err := NewExpenseParticipantRepo(db).Create(tx, u); err != nil {
		tb.Fatal(err)
	}

	return u
}

func TestExpenseParticipantRepo_All(t *testing.T) {
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

			ep := &planetscale.ExpenseParticipant{
				ExpenseID:       1,
				UserID:          "non-existent-user-id",
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			}

			if err := NewExpenseParticipantRepo(db.DB).Create(tx, ep); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})

	t.Run("Upsert Tests", func(t *testing.T) {
		t.Run("successful insert and update", func(t *testing.T) {
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

			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     &g.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})

			ep := &planetscale.ExpenseParticipant{
				ExpenseID:       e.ExpenseID,
				UserID:          u.UserID,
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			}

			if err := NewExpenseParticipantRepo(db.DB).Upsert(tx, ep); err != nil {
				t.Fatal(err)
			}

			// get expense and check if expense id and created at is set
			got, err := NewExpenseParticipantRepo(db.DB).Get(tx, ep.ExpenseID, ep.UserID)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseID != ep.ExpenseID {
				t.Fatalf("expected expense id %d, got %d", ep.ExpenseID, got.ExpenseID)
			} else if got.UserID != ep.UserID {
				t.Fatalf("expected user id %s, got %s", ep.UserID, got.UserID)
			} else if got.AmountOwed != ep.AmountOwed {
				t.Fatalf("expected amount owed %f, got %f", ep.AmountOwed, got.AmountOwed)
			} else if got.SharePercentage != ep.SharePercentage {
				t.Fatalf("expected share percentage %f, got %f", ep.SharePercentage, got.SharePercentage)
			} else if got.Note != ep.Note {
				t.Fatalf("expected note %s, got %s", ep.Note, got.Note)
			}

			ep.AmountOwed = 200
			if err := NewExpenseParticipantRepo(db.DB).Upsert(tx, ep); err != nil {
				t.Fatal(err)
			}

			// get expense and check if amount is updated
			got, err = NewExpenseParticipantRepo(db.DB).Get(tx, ep.ExpenseID, ep.UserID)
			if err != nil {
				t.Fatal(err)
			}
			if got.AmountOwed != ep.AmountOwed {
				t.Fatalf("expected amount owed %f, got %f", ep.AmountOwed, got.AmountOwed)
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
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     &g.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})
			ep := MustCreateExpenseParticipant(t, tx, db.DB, &planetscale.ExpenseParticipant{
				ExpenseID:       e.ExpenseID,
				UserID:          u.UserID,
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			})

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

			epRepo := NewExpenseParticipantRepo(db.DB)
			if _, err := epRepo.Get(tx, 1, u.UserID); err == nil {
				t.Fatal("expected error, got nil")
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
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     &g.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})
			ep := MustCreateExpenseParticipant(t, tx, db.DB, &planetscale.ExpenseParticipant{
				ExpenseID:       e.ExpenseID,
				UserID:          u.UserID,
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			})

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
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     &g.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})
			ep := MustCreateExpenseParticipant(t, tx, db.DB, &planetscale.ExpenseParticipant{
				ExpenseID:       e.ExpenseID,
				UserID:          u.UserID,
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			})

			if err := NewExpenseParticipantRepo(db.DB).Delete(tx, ep.ExpenseID, ep.UserID); err != nil {
				t.Fatal(err)
			}

			if _, err := NewExpenseParticipantRepo(db.DB).Get(tx, ep.ExpenseID, ep.UserID); err == nil {
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
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				GroupID:     &g.ExpenseGroupID,
				PaidBy:      u.UserID,
				SplitTypeID: 1,
				Amount:      100,
				Description: "test expense",
				Timestamp:   time.Now(),
				CreatedBy:   u.UserID,
				UpdatedBy:   u.UserID,
			})
			ep := MustCreateExpenseParticipant(t, tx, db.DB, &planetscale.ExpenseParticipant{
				ExpenseID:       e.ExpenseID,
				UserID:          u.UserID,
				AmountOwed:      100,
				SharePercentage: 100,
				Note:            "test expense",
			})

			epRepo := NewExpenseParticipantRepo(db.DB)
			got, err := epRepo.Find(tx, planetscale.ExpenseParticipantFilter{ExpenseID: e.ExpenseID})
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("expected 1 expense participant, got %d", len(got))
			}
			if got[0].ExpenseID != ep.ExpenseID {
				t.Fatalf("expected expense id %d, got %d", ep.ExpenseID, got[0].ExpenseID)
			}
			if got[0].UserID != ep.UserID {
				t.Fatalf("expected user id %s, got %s", ep.UserID, got[0].UserID)
			}
		})
	})
}
