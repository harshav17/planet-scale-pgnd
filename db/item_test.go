package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateItem(tb testing.TB, tx *sql.Tx, db *DB, i *planetscale.Item) *planetscale.Item {
	tb.Helper()

	if err := NewItemRepo(db).Create(tx, i); err != nil {
		tb.Fatal(err)
	}

	return i
}

func TestItemRepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("Get Tests", func(t *testing.T) {
		t.Run("invalid item id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			if _, err := NewItemRepo(db.DB).Get(tx, 0); err == nil {
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

			got, err := NewItemRepo(db.DB).Get(tx, i.ItemID)
			if err != nil {
				t.Fatal(err)
			}
			if got.ItemID != i.ItemID {
				t.Errorf("expected item id %d, got %d", i.ItemID, got.ItemID)
			} else if got.Name != i.Name {
				t.Errorf("expected item name %s, got %s", i.Name, got.Name)
			} else if got.Price != i.Price {
				t.Errorf("expected item price %f, got %f", i.Price, got.Price)
			} else if got.Quantity != i.Quantity {
				t.Errorf("expected item quantity %d, got %d", i.Quantity, got.Quantity)
			} else if got.ExpenseID != i.ExpenseID {
				t.Errorf("expected item expense id %d, got %d", i.ExpenseID, got.ExpenseID)
			}
		})
	})

	t.Run("Update Tests", func(t *testing.T) {
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

		updateName := "updated item"
		updatePrice := float64(200)
		updateQuantity := int64(2)
		update := &planetscale.ItemUpdate{
			Name:     &updateName,
			Price:    &updatePrice,
			Quantity: &updateQuantity,
		}
		got, err := NewItemRepo(db.DB).Update(tx, i.ItemID, update)
		if err != nil {
			t.Fatal(err)
		}
		if got.ItemID != i.ItemID {
			t.Errorf("expected item id %d, got %d", i.ItemID, got.ItemID)
		} else if got.Name != updateName {
			t.Errorf("expected item name %s, got %s", updateName, got.Name)
		} else if got.Price != updatePrice {
			t.Errorf("expected item price %f, got %f", updatePrice, got.Price)
		} else if got.Quantity != updateQuantity {
			t.Errorf("expected item quantity %d, got %d", update.Quantity, got.Quantity)
		} else if got.ExpenseID != i.ExpenseID {
			t.Errorf("expected item expense %d, got %d", i.ExpenseID, got.ExpenseID)
		}
	})

	t.Run("Delete Tests", func(t *testing.T) {
		t.Run("invalid item id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			if err := NewItemRepo(db.DB).Delete(tx, 0); err == nil {
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

			if err := NewItemRepo(db.DB).Delete(tx, i.ItemID); err != nil {
				t.Fatal(err)
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
			i := MustCreateItem(t, tx, db.DB, &planetscale.Item{
				Name:      "test item",
				Price:     100,
				Quantity:  1,
				ExpenseID: e.ExpenseID,
			})

			items, err := NewItemRepo(db.DB).Find(tx, planetscale.ItemFilter{
				ExpenseID: e.ExpenseID,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 {
				t.Fatalf("expected 1 item, got %d", len(items))
			}
			if items[0].ItemID != i.ItemID {
				t.Errorf("expected item id %d, got %d", i.ItemID, items[0].ItemID)
			} else if items[0].Name != i.Name {
				t.Errorf("expected item name %s, got %s", i.Name, items[0].Name)
			} else if items[0].Price != i.Price {
				t.Errorf("expected item price %f, got %f", i.Price, items[0].Price)
			} else if items[0].Quantity != i.Quantity {
				t.Errorf("expected item quantity %d, got %d", i.Quantity, items[0].Quantity)
			} else if items[0].ExpenseID != i.ExpenseID {
				t.Errorf("expected item expense id %d, got %d", i.ExpenseID, items[0].ExpenseID)
			}
		})
	})
}
