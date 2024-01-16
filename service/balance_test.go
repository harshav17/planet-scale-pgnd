package service

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

type test struct {
	name         string
	expenses     []*planetscale.Expense
	settlements  []*planetscale.Settlement
	participants []*planetscale.ExpenseParticipant
	expected     []*planetscale.Balance
}

func TestBalanceService_GetGroupBalances_TwoParticipants(t *testing.T) {
	var twoExpenseParticipants = []*planetscale.ExpenseParticipant{
		{
			ExpenseID: 1,
			UserID:    "test-user-id",
			Note:      "test expense",
		},
		{
			ExpenseID: 1,
			UserID:    "test-user-id-2",
			Note:      "test expense",
		},
		{
			ExpenseID: 2,
			UserID:    "test-user-id",
			Note:      "test expense",
		},
		{
			ExpenseID: 2,
			UserID:    "test-user-id-2",
			Note:      "test expense",
		},
	}

	var twoExpenses = []*planetscale.Expense{
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
	}

	tests := []test{
		{
			name:         "1 owes 2",
			expenses:     twoExpenses,
			settlements:  []*planetscale.Settlement{},
			participants: twoExpenseParticipants,
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 50,
					BalanceItems: map[string]float64{
						"test-user-id-2": 50,
					},
				},
				{
					UserID: "test-user-id-2",
					Amount: -50,
					BalanceItems: map[string]float64{
						"test-user-id": -50,
					},
				},
			},
		},
		{
			name:     "balance settled",
			expenses: twoExpenses,
			settlements: []*planetscale.Settlement{
				{
					SettlementID: 1,
					GroupID:      1,
					PaidBy:       "test-user-id-2",
					PaidTo:       "test-user-id",
					Amount:       50,
				},
			},
			participants: twoExpenseParticipants,
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 0,
					BalanceItems: map[string]float64{
						"test-user-id-2": 0,
					},
				},
				{
					UserID: "test-user-id-2",
					Amount: 0,
					BalanceItems: map[string]float64{
						"test-user-id": 0,
					},
				},
			},
		},
		{
			name:     "balance settled with multiple settlements",
			expenses: twoExpenses,
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
			participants: twoExpenseParticipants,
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 0,
					BalanceItems: map[string]float64{
						"test-user-id-2": 0,
					},
				},
				{
					UserID: "test-user-id-2",
					Amount: 0,
					BalanceItems: map[string]float64{
						"test-user-id": 0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testGetBalancesHelper(t, test.expenses, test.settlements, test.participants, test.expected)
		})
	}
}

func TestBalanceService_GetGroupBalances_ThreeParticipants(t *testing.T) {
	var threeExpenseParticipants = []*planetscale.ExpenseParticipant{
		{
			ExpenseID: 1,
			UserID:    "test-user-id",
			Note:      "test expense",
		},
		{
			ExpenseID: 1,
			UserID:    "test-user-id-2",
			Note:      "test expense",
		},
		{
			ExpenseID: 1,
			UserID:    "test-user-id-3",
			Note:      "test expense",
		},
		{
			ExpenseID: 2,
			UserID:    "test-user-id-3",
			Note:      "test expense",
		},
		{
			ExpenseID: 2,
			UserID:    "test-user-id-2",
			Note:      "test expense",
		},
	}

	var threeExpenses = []*planetscale.Expense{
		{
			ExpenseID:   1,
			GroupID:     1,
			PaidBy:      "test-user-id",
			Amount:      300,
			SplitTypeID: 1,
		},
		{
			ExpenseID:   2,
			GroupID:     1,
			PaidBy:      "test-user-id-2",
			Amount:      200,
			SplitTypeID: 1,
		},
	}

	tests := []test{
		{
			name:         "1 owes 2 and 3",
			expenses:     threeExpenses,
			settlements:  []*planetscale.Settlement{},
			participants: threeExpenseParticipants,
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 200,
					BalanceItems: map[string]float64{
						"test-user-id-2": 100,
						"test-user-id-3": 100,
					},
				},
				{
					UserID: "test-user-id-2",
					Amount: 0,
					BalanceItems: map[string]float64{
						"test-user-id":   -100,
						"test-user-id-3": 100,
					},
				},
				{
					UserID: "test-user-id-3",
					Amount: -200,
					BalanceItems: map[string]float64{
						"test-user-id":   -100,
						"test-user-id-2": -100,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testGetBalancesHelper(t, test.expenses, test.settlements, test.participants, test.expected)
		})
	}
}

func testGetBalancesHelper(t *testing.T, expenses []*planetscale.Expense, settlements []*planetscale.Settlement, participants []*planetscale.ExpenseParticipant, expected []*planetscale.Balance) {
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

	balanceService.repos.ExpenseParticipant = &db_mock.ExpenseParticipantRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.ExpenseParticipantFilter) ([]*planetscale.ExpenseParticipant, error) {
			return filterParitipantsByExpenseID(participants, filter.ExpenseID), nil
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

func filterParitipantsByExpenseID(participants []*planetscale.ExpenseParticipant, expenseID int64) []*planetscale.ExpenseParticipant {
	var filtered []*planetscale.ExpenseParticipant
	for _, participant := range participants {
		if participant.ExpenseID == expenseID {
			filtered = append(filtered, participant)
		}
	}
	return filtered
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

	// validate balance items
	for userID, expectedBalance := range expectedMap {
		gotBalance, ok := gotMap[userID]
		if !ok {
			return false
		}
		if len(expectedBalance.BalanceItems) != len(gotBalance.BalanceItems) {
			return false
		}
		for paidToUserID, expectedAmount := range expectedBalance.BalanceItems {
			gotAmount, ok := gotBalance.BalanceItems[paidToUserID]
			if !ok {
				return false
			}
			if expectedAmount != gotAmount {
				return false
			}
		}
	}

	return true
}
