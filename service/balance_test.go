package service

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

type test struct {
	name        string
	expenses    []*planetscale.Expense
	settlements []*planetscale.Settlement
	members     []*planetscale.GroupMember
	expected    []*planetscale.Balance
}

func TestBalanceService_GetGroupBalances(t *testing.T) {
	tests := []test{
		{
			name: "1 owes 2",
			expenses: []*planetscale.Expense{
				{
					ExpenseID:   1,
					GroupID:     1,
					PaidBy:      "test-user-id",
					Amount:      200,
					SplitTypeID: 1,
				},
				{
					ExpenseID:   2,
					GroupID:     1,
					PaidBy:      "test-user-id-2",
					Amount:      100,
					SplitTypeID: 1,
				},
			},
			settlements: []*planetscale.Settlement{},
			members: []*planetscale.GroupMember{
				{
					GroupID: 1,
					UserID:  "test-user-id",
				},
				{
					GroupID: 1,
					UserID:  "test-user-id-2",
				},
			},
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 50,
				},
				{
					UserID: "test-user-id-2",
					Amount: -50,
				},
			},
		},
		{
			name: "balance settled",
			expenses: []*planetscale.Expense{
				{
					ExpenseID:   1,
					GroupID:     1,
					PaidBy:      "test-user-id",
					Amount:      200,
					SplitTypeID: 1,
				},
				{
					ExpenseID:   2,
					GroupID:     1,
					PaidBy:      "test-user-id-2",
					Amount:      100,
					SplitTypeID: 1,
				},
			},
			settlements: []*planetscale.Settlement{
				{
					SettlementID: 1,
					GroupID:      1,
					PaidBy:       "test-user-id-2",
					PaidTo:       "test-user-id",
					Amount:       50,
				},
			},
			members: []*planetscale.GroupMember{
				{
					GroupID: 1,
					UserID:  "test-user-id",
				},
				{
					GroupID: 1,
					UserID:  "test-user-id-2",
				},
			},
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 0,
				},
				{
					UserID: "test-user-id-2",
					Amount: 0,
				},
			},
		},
		{
			name: "balance settled with multiple settlements",
			expenses: []*planetscale.Expense{
				{
					ExpenseID:   1,
					GroupID:     1,
					PaidBy:      "test-user-id",
					Amount:      200,
					SplitTypeID: 1,
				},
				{
					ExpenseID:   2,
					GroupID:     1,
					PaidBy:      "test-user-id-2",
					Amount:      100,
					SplitTypeID: 1,
				},
			},
			settlements: []*planetscale.Settlement{
				{
					SettlementID: 1,
					GroupID:      1,
					PaidBy:       "test-user-id-2",
					PaidTo:       "test-user-id",
					Amount:       25,
				},
				{
					SettlementID: 2,
					GroupID:      1,
					PaidBy:       "test-user-id-2",
					PaidTo:       "test-user-id",
					Amount:       25,
				},
			},
			members: []*planetscale.GroupMember{
				{
					GroupID: 1,
					UserID:  "test-user-id",
				},
				{
					GroupID: 1,
					UserID:  "test-user-id-2",
				},
			},
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 0,
				},
				{
					UserID: "test-user-id-2",
					Amount: 0,
				},
			},
		},
	}

	// TODO add tons more tests

	for _, test := range tests {
		testGetBalancesHelper(t, test.expenses, test.settlements, test.members, test.expected)
	}
}

func testGetBalancesHelper(t *testing.T, expenses []*planetscale.Expense, settlements []*planetscale.Settlement, members []*planetscale.GroupMember, expected []*planetscale.Balance) {
	tm := db_mock.TransactionManager{}
	tm.ExecuteInTxFn = func(ctx context.Context, fn func(*sql.Tx) error) error {
		return fn(nil)
	}
	repos := planetscale.RepoProvider{}
	balanceService := NewBalanceService(&repos, tm)

	balanceService.repos.Expense = &db_mock.ExpenseRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.ExpenseFilter) ([]*planetscale.Expense, error) {
			return expenses, nil
		},
	}

	balanceService.repos.Settlement = &db_mock.SettlementRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.SettlementFilter) ([]*planetscale.Settlement, error) {
			return settlements, nil
		},
	}

	balanceService.repos.GroupMember = &db_mock.GroupMemberRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.GroupMemberFilter) ([]*planetscale.GroupMember, error) {
			return members, nil
		},
	}

	balances, err := balanceService.GetGroupBalances(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if !compareBalances(expected, balances) {
		t.Errorf("Expected and actual balances do not match.")
		t.Error("Expected:")
		for _, b := range expected {
			t.Errorf("%+v", b)
		}
		t.Error("Actual:")
		for _, b := range balances {
			t.Errorf("%+v", b)
		}
	}
}

func compareBalances(expected []*planetscale.Balance, got []*planetscale.Balance) bool {
	expectedMap := make(map[string]*planetscale.Balance)
	for _, balance := range expected {
		expectedMap[balance.UserID] = balance
	}

	gotMap := make(map[string]*planetscale.Balance)
	for _, balance := range got {
		gotMap[balance.UserID] = balance
	}

	if len(expectedMap) != len(gotMap) {
		return false
	}

	for userID, expectedBalance := range expectedMap {
		gotBalance, ok := gotMap[userID]
		if !ok {
			return false
		}
		if expectedBalance.Amount != gotBalance.Amount {
			return false
		}
	}

	return true
}
