package service

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

func TestExpenseService_CreateExpense(t *testing.T) {
	groupID := int64(1)
	t.Run("successful create", func(t *testing.T) {
		participant := &planetscale.ExpenseParticipant{
			UserID: "test-user-id",
			Note:   "test expense",
		}
		participant2 := &planetscale.ExpenseParticipant{
			UserID: "test-user-id-2",
			Note:   "test expense 2",
		}

		expense := &planetscale.Expense{
			GroupID:     &groupID,
			PaidBy:      "test-user-id",
			Amount:      100,
			SplitTypeID: 1,
			Participants: []*planetscale.ExpenseParticipant{
				participant,
				participant2,
			},
		}

		repoProvider := &planetscale.RepoProvider{}
		tm := db_mock.TransactionManager{}
		tm.ExecuteInTxFn = func(ctx context.Context, fn func(*sql.Tx) error) error {
			return fn(nil)
		}
		expenseService := NewExpenseService(repoProvider, tm)

		expenseService.repos.Expense = &db_mock.ExpenseRepo{
			CreateFn: func(tx *sql.Tx, expense *planetscale.Expense) error {
				expense.ExpenseID = 1
				return nil
			},
		}

		expenseService.repos.ExpenseParticipant = &db_mock.ExpenseParticipantRepo{
			CreateFn: func(tx *sql.Tx, expenseParticipant *planetscale.ExpenseParticipant) error {
				return nil
			},
		}

		err := expenseService.CreateExpense(context.Background(), expense)
		if err != nil {
			t.Fatal(err)
		}

		if expense.ExpenseID != 1 {
			t.Fatalf("expected expense id to be 1, got %d", expense.ExpenseID)
		} else if expense.Participants[0].ExpenseID != 1 {
			t.Fatalf("expected expense participant expense id to be 1, got %d", expense.Participants[0].ExpenseID)
		} else if expense.Participants[1].ExpenseID != 1 {
			t.Fatalf("expected expense participant expense id to be 1, got %d", expense.Participants[1].ExpenseID)
		}
	})

	t.Run("create equal split type", func(t *testing.T) {
		expense := &planetscale.Expense{
			GroupID:     &groupID,
			PaidBy:      "test-user-id",
			Amount:      100,
			SplitTypeID: 1,
		}

		repoProvider := &planetscale.RepoProvider{}
		tm := db_mock.TransactionManager{}
		tm.ExecuteInTxFn = func(ctx context.Context, fn func(*sql.Tx) error) error {
			return fn(nil)
		}
		expenseService := NewExpenseService(repoProvider, tm)

		expenseService.repos.Expense = &db_mock.ExpenseRepo{
			CreateFn: func(tx *sql.Tx, expense *planetscale.Expense) error {
				expense.ExpenseID = 1
				return nil
			},
		}
		expenseService.repos.GroupMember = &db_mock.GroupMemberRepo{
			FindFn: func(tx *sql.Tx, filter planetscale.GroupMemberFilter) ([]*planetscale.GroupMember, error) {
				return []*planetscale.GroupMember{
					{
						GroupID: 1,
						UserID:  "test-user-id",
					},
					{
						GroupID: 1,
						UserID:  "test-user-id-2",
					},
				}, nil
			},
		}
		expenseService.repos.ExpenseParticipant = &db_mock.ExpenseParticipantRepo{
			CreateFn: func(tx *sql.Tx, expenseParticipant *planetscale.ExpenseParticipant) error {
				return nil
			},
		}

		err := expenseService.CreateExpense(context.Background(), expense)
		if err != nil {
			t.Fatal(err)
		}

		if expense.ExpenseID != 1 {
			t.Fatalf("expected expense id to be 1, got %d", expense.ExpenseID)
		} else if len(expense.Participants) != 2 {
			t.Fatalf("expected 1 expense participant, got %d", len(expense.Participants))
		} else if expense.Participants[0].UserID != "test-user-id" {
			t.Fatalf("expected user id to be test-user-id, got %s", expense.Participants[0].UserID)
		} else if expense.Participants[1].UserID != "test-user-id-2" {
			t.Fatalf("expected user id to be test-user-id-2, got %s", expense.Participants[1].UserID)
		}
	})
}
