package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateGroupMember(tb testing.TB, tx *sql.Tx, db *DB, u *planetscale.GroupMember) *planetscale.GroupMember {
	tb.Helper()

	if err := NewGroupMemberRepo(db).Create(tx, u); err != nil {
		tb.Fatal(err)
	}

	return u
}

func TestGroupMemberRepo_All(t *testing.T) {
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

			u := MustCreateUser(t, tx, db.DB, &planetscale.User{
				UserID: "test-user-id",
				Name:   "test user",
				Email:  "",
			})
			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			gm := &planetscale.GroupMember{
				GroupID: g.ExpenseGroupID,
				UserID:  "non-existent-user-id",
			}

			// TODO - return better error messages from the repo
			if err := NewGroupMemberRepo(db.DB).Create(tx, gm); err == nil {
				t.Fatal("expected error, got nil")
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

			gm := &planetscale.GroupMember{
				GroupID: 1,
				UserID:  u.UserID,
			}

			if err := NewGroupMemberRepo(db.DB).Create(tx, gm); err == nil {
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

			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			gm := MustCreateGroupMember(t, tx, db.DB, &planetscale.GroupMember{
				GroupID: g.ExpenseGroupID,
				UserID:  u.UserID,
			})

			if got, err := NewGroupMemberRepo(db.DB).Get(tx, gm.GroupID, gm.UserID); err != nil {
				t.Fatal(err)
			} else if got.GroupID != gm.GroupID {
				t.Fatalf("expected group id to be %d, got %d", gm.GroupID, got.GroupID)
			}
		})

		t.Run("invalid user id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			// TODO - return better error messages from the repo
			if _, err := NewGroupMemberRepo(db.DB).Get(tx, 1, "non-existent-user-id"); err == nil {
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

			g := MustCreateExpenseGroup(t, tx, db.DB, &planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  u.UserID,
			})

			gm := MustCreateGroupMember(t, tx, db.DB, &planetscale.GroupMember{
				GroupID: g.ExpenseGroupID,
				UserID:  u.UserID,
			})

			if err := NewGroupMemberRepo(db.DB).Delete(tx, gm.GroupID, gm.UserID); err != nil {
				t.Fatal(err)
			}
			// get the group member again to make sure it was deleted
			if _, err := NewGroupMemberRepo(db.DB).Get(tx, gm.GroupID, gm.UserID); err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("invalid user id", func(t *testing.T) {
			tx, err := db.db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			// TODO - return better error messages from the repo
			if err := NewGroupMemberRepo(db.DB).Delete(tx, 1, "non-existent-user-id"); err == nil {
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

			gm := MustCreateGroupMember(t, tx, db.DB, &planetscale.GroupMember{
				GroupID: g.ExpenseGroupID,
				UserID:  u.UserID,
			})

			if got, err := NewGroupMemberRepo(db.DB).Find(tx, planetscale.GroupMemberFilter{
				GroupID: gm.GroupID,
			}); err != nil {
				t.Fatal(err)
			} else if len(got) != 1 {
				t.Fatalf("expected length of group members to be 1, got %d", len(got))
			} else if got[0].GroupID != gm.GroupID {
				t.Fatalf("expected group id to be %d, got %d", gm.GroupID, got[0].GroupID)
			} else if got[0].UserID != gm.UserID {
				t.Fatalf("expected user id to be %s, got %s", gm.UserID, got[0].UserID)
			} else if got[0].User == nil {
				t.Fatal("expected user to be populated, got nil")
			} else if got[0].User.UserID != u.UserID {
				t.Fatalf("expected user id to be %s, got %s", u.UserID, got[0].User.UserID)
			} else if got[0].User.Name != u.Name {
				t.Fatalf("expected user name to be %s, got %s", u.Name, got[0].User.Name)
			} else if got[0].User.Email != u.Email {
				t.Fatalf("expected user email to be %s, got %s", u.Email, got[0].User.Email)
			}
		})
	})
}
