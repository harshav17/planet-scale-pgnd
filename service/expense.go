package service

import (
	"context"
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type expenseService struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewExpenseService(repoProvider *planetscale.RepoProvider, tm planetscale.TransactionManager) *expenseService {
	return &expenseService{
		repos: repoProvider,
		tm:    tm,
	}
}

func (s *expenseService) CreateExpense(ctx context.Context, expense *planetscale.Expense) error {
	createExpenseFunc := func(tx *sql.Tx) error {
		err := s.repos.Expense.Create(tx, expense)
		if err != nil {
			return err
		}

		if expense.SplitTypeID == 1 {
			if len(expense.Participants) == 0 {
				members, err := s.repos.GroupMember.Find(tx, planetscale.GroupMemberFilter{
					GroupID: *expense.GroupID,
				})
				if err != nil {
					return err
				}

				for _, member := range members {
					participant := &planetscale.ExpenseParticipant{
						UserID: member.UserID,
					}
					participant.ExpenseID = expense.ExpenseID
					err := s.repos.ExpenseParticipant.Create(tx, participant)
					if err != nil {
						return err
					}
					expense.Participants = append(expense.Participants, participant)
				}
			} else {
				for _, participant := range expense.Participants {
					participant.ExpenseID = expense.ExpenseID
					err := s.repos.ExpenseParticipant.Create(tx, participant)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	err := s.tm.ExecuteInTx(ctx, createExpenseFunc)
	if err != nil {
		return err
	}

	return nil
}
