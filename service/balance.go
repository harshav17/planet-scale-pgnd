package service

import (
	"context"
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type balanceService struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewBalanceService(repoProvider *planetscale.RepoProvider, tm planetscale.TransactionManager) *balanceService {
	return &balanceService{
		repos: repoProvider,
		tm:    tm,
	}
}

func (s *balanceService) GetGroupBalances(ctx context.Context, groupID int64) ([]*planetscale.Balance, error) {
	var balances []*planetscale.Balance
	getBalancesFunc := func(tx *sql.Tx) error {
		expenses, err := s.repos.Expense.Find(tx, planetscale.ExpenseFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}

		settlements, err := s.repos.Settlement.Find(tx, planetscale.SettlementFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}

		balances, err = s.calculateBalances(tx, expenses, settlements)
		if err != nil {
			return err
		}

		return nil
	}

	err := s.tm.ExecuteInTx(ctx, getBalancesFunc)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (s *balanceService) calculateBalances(tx *sql.Tx, expenses []*planetscale.Expense, settlements []*planetscale.Settlement) ([]*planetscale.Balance, error) {
	// compile a list of balance records
	balances := make(map[string]*planetscale.Balance)
	for _, expense := range expenses {
		if _, ok := balances[expense.PaidBy]; !ok {
			balances[expense.PaidBy] = &planetscale.Balance{
				UserID:       expense.PaidBy,
				BalanceItems: map[string]float64{},
			}
		}
		balances[expense.PaidBy].Amount += expense.Amount

		if expense.SplitTypeID == 1 { // TODO load from DB
			s.handleEqualSplitType(tx, balances, expense)
		} else if expense.SplitTypeID == 3 {
			s.handleItemizedSplitType(tx, balances, expense)
		}
	}
	for _, settlement := range settlements {
		if _, ok := balances[settlement.PaidBy]; !ok {
			balances[settlement.PaidBy] = &planetscale.Balance{
				UserID: settlement.PaidBy,
			}
		}
		balances[settlement.PaidBy].Amount += settlement.Amount

		// update paidBy user's balance items
		if _, ok := balances[settlement.PaidBy].BalanceItems[settlement.PaidTo]; !ok {
			balances[settlement.PaidBy].BalanceItems[settlement.PaidTo] = 0
		}
		balances[settlement.PaidBy].BalanceItems[settlement.PaidTo] += settlement.Amount
		balances[settlement.PaidTo].BalanceItems[settlement.PaidBy] -= settlement.Amount

		if _, ok := balances[settlement.PaidTo]; !ok {
			balances[settlement.PaidTo] = &planetscale.Balance{
				UserID: settlement.PaidTo,
			}
		}
		balances[settlement.PaidTo].Amount -= settlement.Amount
	}

	// convert map to slice
	var balanceSlice []*planetscale.Balance
	for _, balance := range balances {
		// truncate to 2 decimal places
		balance.Amount = float64(int64(balance.Amount*100)) / 100
		balanceSlice = append(balanceSlice, balance)
	}
	return balanceSlice, nil
}

func (s *balanceService) handleEqualSplitType(tx *sql.Tx, balances map[string]*planetscale.Balance, expense *planetscale.Expense) error {
	participants, err := s.repos.ExpenseParticipant.Find(tx, planetscale.ExpenseParticipantFilter{
		ExpenseID: expense.ExpenseID,
	})
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if _, ok := balances[participant.UserID]; !ok {
			balances[participant.UserID] = &planetscale.Balance{
				UserID:       participant.UserID,
				BalanceItems: map[string]float64{},
			}
		}
		balances[participant.UserID].Amount -= expense.Amount / float64(len(participants))

		if participant.UserID == expense.PaidBy {
			continue
		}
		balances[participant.UserID].BalanceItems[expense.PaidBy] -= expense.Amount / float64(len(participants))
		balances[expense.PaidBy].BalanceItems[participant.UserID] += expense.Amount / float64(len(participants))
	}
	return nil
}

func (s *balanceService) handleItemizedSplitType(tx *sql.Tx, balances map[string]*planetscale.Balance, expense *planetscale.Expense) error {
	return nil
}
