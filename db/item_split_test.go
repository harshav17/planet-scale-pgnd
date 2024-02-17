package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateItemSplit(tb testing.TB, tx *sql.Tx, db *DB, i *planetscale.ItemSplit) *planetscale.ItemSplit {
	tb.Helper()

	if err := NewItemSplitRepo(db).Create(tx, i); err != nil {
		tb.Fatal(err)
	}

	return i
}

func TestItemSplitRepo_All(t *testing.T) {
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

			if _, err := NewItemSplitRepo(db.DB).Get(tx, 0); err == nil {
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
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     100,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})
			is := MustCreateItemSplit(t, tx, db.DB, &planetscale.ItemSplit{
				ItemID: i.ItemID,
				UserID: u.UserID,
				Amount: 100,
			})

			is2, err := NewItemSplitRepo(db.DB).Get(tx, is.ItemSplitID)
			if err != nil {
				t.Fatal(err)
			}

			if is2.ItemSplitID != is.ItemSplitID {
				t.Fatalf("expected item split id %d, got %d", is.ItemSplitID, is2.ItemSplitID)
			} else if is2.ItemID != is.ItemID {
				t.Fatalf("expected item id %d, got %d", is.ItemID, is2.ItemID)
			} else if is2.UserID != is.UserID {
				t.Fatalf("expected user id %s, got %s", is.UserID, is2.UserID)
			} else if is2.Amount != is.Amount {
				t.Fatalf("expected amount %f, got %f", is.Amount, is2.Amount)
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

			if _, err := NewItemSplitRepo(db.DB).Update(tx, 0, &planetscale.ItemSplitUpdate{}); err == nil {
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
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     100,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})
			is := MustCreateItemSplit(t, tx, db.DB, &planetscale.ItemSplit{
				ItemID: i.ItemID,
				UserID: u.UserID,
				Amount: 100,
			})

			amount := 200.0
			is2, err := NewItemSplitRepo(db.DB).Update(tx, is.ItemSplitID, &planetscale.ItemSplitUpdate{
				Amount: &amount,
			})
			if err != nil {
				t.Fatal(err)
			}

			if is2.ItemSplitID != is.ItemSplitID {
				t.Fatalf("expected item split id %d, got %d", is.ItemSplitID, is2.ItemSplitID)
			} else if is2.ItemID != is.ItemID {
				t.Fatalf("expected item id %d, got %d", is.ItemID, is2.ItemID)
			} else if is2.UserID != is.UserID {
				t.Fatalf("expected user id %s, got %s", is.UserID, is2.UserID)
			} else if is2.Amount != amount {
				t.Fatalf("expected amount %f, got %f", amount, is2.Amount)
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

			if err := NewItemSplitRepo(db.DB).Delete(tx, 0); err == nil {
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
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     100,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})
			is := MustCreateItemSplit(t, tx, db.DB, &planetscale.ItemSplit{
				ItemID: i.ItemID,
				UserID: u.UserID,
				Amount: 100,
			})

			if err := NewItemSplitRepo(db.DB).Delete(tx, is.ItemSplitID); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("List Tests", func(t *testing.T) {
		t.Run("invalid item id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			filter := planetscale.ItemSplitFilter{
				ItemID: 0,
			}
			if items, err := NewItemSplitRepo(db.DB).Find(tx, filter); err != nil {
				t.Fatal(err)
			} else if len(items) != 0 {
				t.Fatalf("expected 0 items, got %d", len(items))
			}
		})

		t.Run("successful list", func(t *testing.T) {
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
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     100,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})
			is := MustCreateItemSplit(t, tx, db.DB, &planetscale.ItemSplit{
				ItemID: i.ItemID,
				UserID: u.UserID,
				Amount: 100,
			})

			filter := planetscale.ItemSplitFilter{
				ItemID: i.ItemID,
			}
			itemSplits, err := NewItemSplitRepo(db.DB).Find(tx, filter)
			if err != nil {
				t.Fatal(err)
			}

			if len(itemSplits) != 1 {
				t.Fatalf("expected 1 item split, got %d", len(itemSplits))
			} else if itemSplits[0].ItemSplitID != is.ItemSplitID {
				t.Fatalf("expected item split id %d, got %d", is.ItemSplitID, itemSplits[0].ItemSplitID)
			} else if itemSplits[0].ItemID != is.ItemID {
				t.Fatalf("expected item id %d, got %d", is.ItemID, itemSplits[0].ItemID)
			} else if itemSplits[0].UserID != is.UserID {
				t.Fatalf("expected user id %s, got %s", is.UserID, itemSplits[0].UserID)
			} else if itemSplits[0].Amount != is.Amount {
				t.Fatalf("expected amount %f, got %f", is.Amount, itemSplits[0].Amount)
			}
		})
	})
}
