package planetscale

import "database/sql"

type (
	ExpenseParticipant struct {
		ExpenseID       int64   `json:"expense_id"`
		UserID          string  `json:"user_id"`
		AmountOwed      float64 `json:"amount_owed"`
		SharePercentage float64 `json:"share_percentage"`
		SplitMethod     string  `json:"split_method"`
		Note            string  `json:"note"`
	}

	ExpenseParticipantRepo interface {
		Get(tx *sql.Tx, expenseID int64, userID string) (*ExpenseParticipant, error)
		Create(tx *sql.Tx, expense *ExpenseParticipant) error
		Delete(tx *sql.Tx, expenseID int64, userID string) error
		Update(tx *sql.Tx, expenseID int64, userID string, expense *ExpenseParticipantUpdate) (*ExpenseParticipant, error)
		Find(tx *sql.Tx, filter ExpenseParticipantFilter) ([]*ExpenseParticipant, error)
	}

	ExpenseParticipantUpdate struct {
		AmountOwed      *float64 `json:"amount_owed"`
		SharePercentage *float64 `json:"share_percentage"`
		SplitMethod     *string  `json:"split_method"`
		Note            *string  `json:"note"`
	}

	ExpenseParticipantFilter struct {
		ExpenseID int64
	}
)
