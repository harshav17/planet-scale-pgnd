package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateItemSplitNu(tb testing.TB, tx *sql.Tx, db *DB, itemSplit *planetscale.ItemSplitNU) *planetscale.ItemSplitNU {
	tb.Helper()

	if err := NewItemSplitNURepo(db).Create(tx, itemSplit); err != nil {
		tb.Fatalf("failed to create item split: %v", err)
	}
	return itemSplit
}

func TestItemSplitNURepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("Get Tests", func(t *testing.T) {
		t.Run("invalid item split id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			if _, err := NewItemSplitNURepo(db.DB).Get(tx, 0); err == nil {
				t.Fatal("expected error, got nil")
			}
		})

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
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				PaidBy:      u.UserID,
				Description: "test expense",
				Amount:      10.0,
				CreatedAt:   time.Now(),
				CreatedBy:   u.UserID,
				SplitTypeID: 3,
			})
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			testInitials := "TU"
			is := MustCreateItemSplitNu(t, tx, db.DB, &planetscale.ItemSplitNU{
				ItemID:   i.ItemID,
				Initials: &testInitials,
				Amount:   10.0,
			})

			got, err := NewItemSplitNURepo(db.DB).Get(tx, is.ItemSplitID)
			if err != nil {
				t.Fatalf("failed to get item split: %v", err)
			} else if got == nil {
				t.Fatal("expected item split, got nil")
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

			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   "test user",
				Email:  "",
			})
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				PaidBy:      u.UserID,
				Description: "test expense",
				Amount:      10.0,
				CreatedAt:   time.Now(),
				CreatedBy:   u.UserID,
				SplitTypeID: 3,
			})
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			testInitials := "TU"
			is := &planetscale.ItemSplitNU{
				ItemID:   i.ItemID,
				Initials: &testInitials,
				Amount:   10.0,
			}

			if err := NewItemSplitNURepo(db.DB).Create(tx, is); err != nil {
				t.Fatalf("failed to create item split: %v", err)
			}
			if is.ItemSplitID == 0 {
				t.Fatal("expected item split id to be set, got 0")
			}
		})
	})

	t.Run("Update Tests", func(t *testing.T) {
		t.Run("invalid item split id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			if _, err := NewItemSplitNURepo(db.DB).Update(tx, 0, &planetscale.ItemSplitNUUpdate{}); err == nil {
				t.Fatal("expected error, got nil")
			}
		})

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
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				PaidBy:      u.UserID,
				Description: "test expense",
				Amount:      10.0,
				CreatedAt:   time.Now(),
				CreatedBy:   u.UserID,
				SplitTypeID: 3,
			})
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			testInitials := "TU"
			is := MustCreateItemSplitNu(t, tx, db.DB, &planetscale.ItemSplitNU{
				ItemID:   i.ItemID,
				Initials: &testInitials,
				Amount:   10.0,
			})

			updateAmount := 20.0
			update := &planetscale.ItemSplitNUUpdate{
				Amount: &updateAmount,
			}

			got, err := NewItemSplitNURepo(db.DB).Update(tx, is.ItemSplitID, update)
			if err != nil {
				t.Fatalf("failed to update item split: %v", err)
			} else if got == nil {
				t.Fatal("expected item split, got nil")
			} else if got.Amount != *update.Amount {
				t.Fatalf("expected amount %f, got %f", *update.Amount, got.Amount)
			}
		})
	})

	t.Run("Delete Tests", func(t *testing.T) {
		t.Run("invalid item split id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			if err := NewItemSplitNURepo(db.DB).Delete(tx, 0); err == nil {
				t.Fatal("expected error, got nil")
			}
		})

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
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				PaidBy:      u.UserID,
				Description: "test expense",
				Amount:      10.0,
				CreatedAt:   time.Now(),
				CreatedBy:   u.UserID,
				SplitTypeID: 3,
			})
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			testInitials := "TU"
			is := MustCreateItemSplitNu(t, tx, db.DB, &planetscale.ItemSplitNU{
				ItemID:   i.ItemID,
				Initials: &testInitials,
				Amount:   10.0,
			})

			if err := NewItemSplitNURepo(db.DB).Delete(tx, is.ItemSplitID); err != nil {
				t.Fatalf("failed to delete item split: %v", err)
			}
		})
	})

	t.Run("Find Tests", func(t *testing.T) {
		t.Run("no item splits found", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			found, err := NewItemSplitNURepo(db.DB).Find(tx, planetscale.ItemSplitNUFilter{})
			if err != nil {
				t.Fatalf("failed to find item splits: %v", err)
			} else if len(found) != 0 {
				t.Fatalf("expected 0 item splits, got %d", len(found))
			}
		})

		t.Run("invalid item id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			found, err := NewItemSplitNURepo(db.DB).Find(tx, planetscale.ItemSplitNUFilter{ItemID: 8})
			if err != nil {
				t.Fatalf("failed to find item splits: %v", err)
			} else if len(found) != 0 {
				t.Fatalf("expected 0 item splits, got %d", len(found))
			}
		})

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
			e := MustCreateExpense(t, tx, db.DB, &planetscale.Expense{
				PaidBy:      u.UserID,
				Description: "test expense",
				Amount:      10.0,
				CreatedAt:   time.Now(),
				CreatedBy:   u.UserID,
				SplitTypeID: 3,
			})
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			testInitials := "TU"
			_ = MustCreateItemSplitNu(t, tx, db.DB, &planetscale.ItemSplitNU{
				ItemID:   i.ItemID,
				Initials: &testInitials,
				Amount:   10.0,
			})

			if found, err := NewItemSplitNURepo(db.DB).Find(tx, planetscale.ItemSplitNUFilter{ItemID: i.ItemID}); err != nil {
				t.Fatalf("failed to find item split: %v", err)
			} else if len(found) != 1 {
				t.Fatalf("expected 1 item split, got %d", len(found))
			} else if found[0].ItemID != i.ItemID {
				t.Fatalf("expected item id %d, got %d", i.ItemID, found[0].ItemID)
			}
		})
	})
}
