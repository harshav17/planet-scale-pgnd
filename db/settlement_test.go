package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateSettlement(tb testing.TB, tx *sql.Tx, db *DB, s *planetscale.Settlement) *planetscale.Settlement {
	tb.Helper()

	if err := NewSettlementRepo(db).Create(tx, s); err != nil {
		tb.Fatal(err)
	}

	return s
}

func TestSettlementRepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

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

			u2 := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id-2",
				Name:   "test user 2",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			s := &planetscale.Settlement{
				GroupID: g.ExpenseGroupID,
				PaidBy:  u.UserID,
				PaidTo:  u2.UserID,
				Amount:  100,
			}

			if err := NewSettlementRepo(db.DB).Create(tx, s); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("invalid group id", func(t *testing.T) {
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
			s := &planetscale.Settlement{
				GroupID: 100,
				PaidBy:  u.UserID,
				PaidTo:  u.UserID,
				Amount:  100,
			}

			// TODO - return better error messages from the repo
			if err := NewSettlementRepo(db.DB).Create(tx, s); err == nil {
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
			u2 := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id-2",
				Name:   "test user 2",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			s := MustCreateSettlement(t, tx, db.DB, &planetscale.Settlement{
				GroupID: g.ExpenseGroupID,
				PaidBy:  u.UserID,
				PaidTo:  u2.UserID,
				Amount:  100,
			})

			if got, err := NewSettlementRepo(db.DB).Get(tx, s.SettlementID); err != nil {
				t.Fatal(err)
			} else if got.Amount != 100 {
				t.Fatalf("expected title to be %d, got %f", 100, got.Amount)
			}
		})

		t.Run("invalid settlement id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback() // TODO - pull this out into a helper

			// TODO - return better error messages from the repo
			if _, err := NewSettlementRepo(db.DB).Get(tx, 1); err == nil {
				t.Fatal("expected error, got nil")
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
			u2 := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id-2",
				Name:   "test user 2",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			s := MustCreateSettlement(t, tx, db.DB, &planetscale.Settlement{
				GroupID: g.ExpenseGroupID,
				PaidBy:  u.UserID,
				PaidTo:  u2.UserID,
				Amount:  100,
			})

			if err := NewSettlementRepo(db.DB).Delete(tx, s.SettlementID); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("invalid settlement id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			// TODO - return better error messages from the repo
			if err := NewSettlementRepo(db.DB).Delete(tx, 100); err == nil {
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
			u2 := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id-2",
				Name:   "test user 2",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			g2 := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group 2",
				CreateBy:  u.UserID,
			})
			s := MustCreateSettlement(t, tx, db.DB, &planetscale.Settlement{
				GroupID: g.ExpenseGroupID,
				PaidBy:  u.UserID,
				PaidTo:  u2.UserID,
				Amount:  100,
			})

			su := &planetscale.SettlementUpdate{
				GroupID: &g2.ExpenseGroupID,
			}

			if _, err := NewSettlementRepo(db.DB).Update(tx, s.SettlementID, su); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("invalid settlement id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			gid := int64(100)
			su := &planetscale.SettlementUpdate{
				GroupID: &gid,
			}

			// TODO - return better error messages from the repo
			if _, err := NewSettlementRepo(db.DB).Update(tx, 100, su); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})

	t.Run("find tests", func(t *testing.T) {
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
			u2 := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id-2",
				Name:   "test user 2",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})
			s := MustCreateSettlement(t, tx, db.DB, &planetscale.Settlement{
				GroupID: g.ExpenseGroupID,
				PaidBy:  u.UserID,
				PaidTo:  u2.UserID,
				Amount:  100,
			})

			f := planetscale.SettlementFilter{
				GroupID: g.ExpenseGroupID,
			}

			if got, err := NewSettlementRepo(db.DB).Find(tx, f); err != nil {
				t.Fatal(err)
			} else if len(got) != 1 {
				t.Fatalf("expected 1 settlement, got %d", len(got))
			} else if got[0].Amount != s.Amount {
				t.Fatalf("expected title to be %f, got %f", s.Amount, got[0].Amount)
			}
		})
	})
}
