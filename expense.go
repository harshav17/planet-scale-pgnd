package planetscale

import (
	"database/sql"
	"time"
)

type (
	Expense struct {
		ExpenseID   int64     `json:"expense_id"`
		GroupID     int64     `json:"group_id"`
		PaidBy      string    `json:"paid_by"`
		Amount      float64   `json:"amount"`
		Description string    `json:"description"`
		Timestamp   time.Time `json:"timestamp"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedBy   string    `json:"created_by"`
		UpdatedBy   string    `json:"updated_by"`
	}

	ExpenseRepo interface {
		Get(tx *sql.Tx, expenseID int64) (*Expense, error)
		Create(tx *sql.Tx, expense *Expense) error
		Delete(tx *sql.Tx, expenseID int64) error
		Update(tx *sql.Tx, expenseID int64, expense *ExpenseUpdate) error
		Find(tx *sql.Tx, filter ExpenseFilter) ([]*Expense, error)
	}

	ExpenseFilter struct {
		GroupID int64
	}

	ExpenseUpdate struct {
		GroupID     *int64     `json:"group_id"`
		PaidBy      *string    `json:"paid_by"`
		Amount      *float64   `json:"amount"`
		Description *string    `json:"description"`
		Timestamp   *time.Time `json:"timestamp"`
		UpdatedBy   *string    `json:"updated_by"`
	}
)
