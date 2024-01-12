package db

import (
	"context"
	"testing"
)

func TestSplitTypeRepo_All(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	defer MustCloseDB(t, db)
	ctx := context.Background()

	t.Run("Get Tests", func(t *testing.T) {
		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()

		if got, err := NewSplitTypeRepo(db.DB).Get(tx, 1); err != nil {
			t.Fatal(err)
		} else if got.SplitTypeID != 1 {
			t.Fatalf("expected title to be %d, got %d", 1, got.SplitTypeID)
		} else if got.TypeName != "Equal" {
			t.Fatalf("expected title to be %s, got %s", "Equal", got.TypeName)
		}
	})

	t.Run("GetAll Tests", func(t *testing.T) {
		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()

		if got, err := NewSplitTypeRepo(db.DB).GetAll(tx); err != nil {
			t.Fatal(err)
		} else if len(got) != 5 {
			t.Fatalf("expected length to be %d, got %d", 3, len(got))
		} else if got[0].TypeName != "Equal" {
			t.Fatalf("expected title to be %s, got %s", "Equal", got[0].TypeName)
		} else if got[1].TypeName != "Unequal" {
			t.Fatalf("expected title to be %s, got %s", "Unequal", got[1].TypeName)
		} else if got[2].TypeName != "ItemBased" {
			t.Fatalf("expected title to be %s, got %s", "ItemBased", got[2].TypeName)
		} else if got[3].TypeName != "ShareBased" {
			t.Fatalf("expected title to be %s, got %s", "ShareBased", got[3].TypeName)
		} else if got[4].TypeName != "PercentageBased" {
			t.Fatalf("expected title to be %s, got %s", "PercentageBased", got[4].TypeName)
		}
	})
}
