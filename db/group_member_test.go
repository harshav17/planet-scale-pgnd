package db

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
)

func MustCreateGroupMember(tb testing.TB, ctx context.Context, db *DB, u *planetscale.GroupMember) (*planetscale.GroupMember, context.Context) {
	tb.Helper()

	createGroupMemberFunc := func(tx *sql.Tx) error {
		if err := NewGroupMemberRepo(db).Create(tx, u); err != nil {
			return err
		}
		return nil
	}

	tm := NewTransactionManager(db)
	err := tm.ExecuteInTx(ctx, createGroupMemberFunc)
	if err != nil {
		tb.Fatal(err)
	}
	return u, ctx
}

func TestGroupMemberRepo_CreateGroupMember(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("invalid user id", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})
		g, ctx := MustCreateExpenseGroup(t, ctx, db.DB, &planetscale.ExpenseGroup{
			GroupName: "test group",
			CreateBy:  u.UserID,
		})

		gm := &planetscale.GroupMember{
			GroupID: g.ExpenseGroupID,
			UserID:  "non-existent-user-id",
		}

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		// TODO - return better error messages from the repo
		if err := NewGroupMemberRepo(db.DB).Create(tx, gm); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid group id", func(t *testing.T) {
		u, ctx := MustCreateUser(t, ctx, db.DB, &planetscale.User{
			UserID: "test-user-id",
			Name:   "test user",
			Email:  "",
		})

		gm := &planetscale.GroupMember{
			GroupID: 1,
			UserID:  u.UserID,
		}

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := NewGroupMemberRepo(db.DB).Create(tx, gm); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGroupMemberRepo_GetGroupMember(t *testing.T) {
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

		gm, ctx := MustCreateGroupMember(t, ctx, db.DB, &planetscale.GroupMember{
			GroupID: g.ExpenseGroupID,
			UserID:  u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
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
		// TODO - return better error messages from the repo
		if _, err := NewGroupMemberRepo(db.DB).Get(tx, 1, "non-existent-user-id"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGroupMemberRepo_DeleteGroupMember(t *testing.T) {
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

		gm, ctx := MustCreateGroupMember(t, ctx, db.DB, &planetscale.GroupMember{
			GroupID: g.ExpenseGroupID,
			UserID:  u.UserID,
		})

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
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
		// TODO - return better error messages from the repo
		if err := NewGroupMemberRepo(db.DB).Delete(tx, 1, "non-existent-user-id"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
