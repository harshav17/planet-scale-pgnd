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

		members, err := s.repos.GroupMember.Find(tx, planetscale.GroupMemberFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}

		balances = calculateBalances(expenses, settlements, members)
		return nil
	}

	err := s.tm.ExecuteInTx(ctx, getBalancesFunc)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func calculateBalances(expenses []*planetscale.Expense, settlements []*planetscale.Settlement, members []*planetscale.GroupMember) []*planetscale.Balance {
	// compile a list of balance records

	balances := make(map[string]*planetscale.Balance)
	for _, expense := range expenses {
		if _, ok := balances[expense.PaidBy]; !ok {
			balances[expense.PaidBy] = &planetscale.Balance{
				UserID: expense.PaidBy,
			}
		}
		balances[expense.PaidBy].Amount += expense.Amount

		if expense.SplitTypeID == 1 { // TODO load from DB
			// TODO refactor into an extensible split type system
			// split equally among all members
			for _, member := range members {
				if _, ok := balances[member.UserID]; !ok {
					balances[member.UserID] = &planetscale.Balance{
						UserID: member.UserID,
					}
				}
				balances[member.UserID].Amount -= expense.Amount / float64(len(members))
			}
		}
	}
	for _, settlement := range settlements {
		if _, ok := balances[settlement.PaidBy]; !ok {
			balances[settlement.PaidBy] = &planetscale.Balance{
				UserID: settlement.PaidBy,
			}
		}
		balances[settlement.PaidBy].Amount += settlement.Amount
	}
	for _, settlement := range settlements {
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
	return balanceSlice
}
