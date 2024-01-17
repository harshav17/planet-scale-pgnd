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
	items        []*planetscale.Item
	itemSplits   []*planetscale.ItemSplit
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
			items:        []*planetscale.Item{},
			itemSplits:   []*planetscale.ItemSplit{},
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
			testGetBalancesHelper(t, test.expenses, test.settlements, test.participants, test.items, test.itemSplits, test.expected)
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
			items:        []*planetscale.Item{},
			itemSplits:   []*planetscale.ItemSplit{},
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
			testGetBalancesHelper(t, test.expenses, test.settlements, test.participants, test.items, test.itemSplits, test.expected)
		})
	}
}

func TestBalanceService_GetGroupBalances_ItemizedSplit(t *testing.T) {
	var items = []*planetscale.Item{
		{
			ItemID:    1,
			Name:      "test item",
			Price:     100,
			Quantity:  1,
			ExpenseID: 1,
		},
		{
			ItemID:    2,
			Name:      "test item 2",
			Price:     200,
			Quantity:  2,
			ExpenseID: 1,
		},
		{
			ItemID:    3,
			Name:      "test item 3",
			Price:     300,
			Quantity:  3,
			ExpenseID: 1,
		},
	}

	itemSplits := []*planetscale.ItemSplit{
		{
			ItemSplitID: 1,
			ItemID:      1,
			UserID:      "test-user-id",
		},
		{
			ItemSplitID: 2,
			ItemID:      2,
			UserID:      "test-user-id",
		},
		{
			ItemSplitID: 3,
			ItemID:      2,
			UserID:      "test-user-id-2",
		},
		{
			ItemSplitID: 4,
			ItemID:      3,
			UserID:      "test-user-id",
		},
		{
			ItemSplitID: 5,
			ItemID:      3,
			UserID:      "test-user-id-2",
		},
		{
			ItemSplitID: 6,
			ItemID:      3,
			UserID:      "test-user-id-3",
		},
	}

	var expense = []*planetscale.Expense{
		{
			ExpenseID:   1,
			GroupID:     1,
			PaidBy:      "test-user-id",
			Amount:      600,
			SplitTypeID: 3,
		},
	}

	tests := []test{
		{
			name:         "1 owes 2 and 3",
			expenses:     expense,
			settlements:  []*planetscale.Settlement{},
			participants: []*planetscale.ExpenseParticipant{},
			items:        items,
			itemSplits:   itemSplits,
			expected: []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 300,
					BalanceItems: map[string]float64{
						"test-user-id-2": 200,
						"test-user-id-3": 100,
					},
				},
				{
					UserID: "test-user-id-2",
					Amount: -200,
					BalanceItems: map[string]float64{
						"test-user-id": -200,
					},
				},
				{
					UserID: "test-user-id-3",
					Amount: -100,
					BalanceItems: map[string]float64{
						"test-user-id": -100,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testGetBalancesHelper(t, test.expenses, test.settlements, test.participants, test.items, test.itemSplits, test.expected)
		})
	}
}

func testGetBalancesHelper(
	t *testing.T,
	expenses []*planetscale.Expense,
	settlements []*planetscale.Settlement,
	participants []*planetscale.ExpenseParticipant,
	items []*planetscale.Item,
	itemSplits []*planetscale.ItemSplit,
	expected []*planetscale.Balance) {
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

	balanceService.repos.Item = &db_mock.ItemRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.ItemFilter) ([]*planetscale.Item, error) {
			return items, nil
		},
	}

	balanceService.repos.ItemSplit = &db_mock.ItemSplitRepo{
		FindFn: func(tx *sql.Tx, filter planetscale.ItemSplitFilter) ([]*planetscale.ItemSplit, error) {
			return filterItemSplitsByItemID(itemSplits, filter.ItemID), nil
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

func filterItemSplitsByItemID(itemSplits []*planetscale.ItemSplit, itemID int64) []*planetscale.ItemSplit {
	var filtered []*planetscale.ItemSplit
	for _, itemSplit := range itemSplits {
		if itemSplit.ItemID == itemID {
			filtered = append(filtered, itemSplit)
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
