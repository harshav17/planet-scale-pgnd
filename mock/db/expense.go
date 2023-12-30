package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ExpenseRepo struct {
	GetFn    func(tx *sql.Tx, expenseID int64) (*planetscale.Expense, error)
	CreateFn func(tx *sql.Tx, expense *planetscale.Expense) error
	DeleteFn func(tx *sql.Tx, expenseID int64) error
	UpdateFn func(tx *sql.Tx, expenseID int64, expense *planetscale.ExpenseUpdate) (*planetscale.Expense, error)
	FindFn   func(tx *sql.Tx, filter planetscale.ExpenseFilter) ([]*planetscale.Expense, error)
}

func (s ExpenseRepo) Get(tx *sql.Tx, expenseID int64) (*planetscale.Expense, error) {
	return s.GetFn(tx, expenseID)
}

func (s ExpenseRepo) Create(tx *sql.Tx, expense *planetscale.Expense) error {
	return s.CreateFn(tx, expense)
}

func (s ExpenseRepo) Delete(tx *sql.Tx, expenseID int64) error {
	return s.DeleteFn(tx, expenseID)
}

func (s ExpenseRepo) Update(tx *sql.Tx, expenseID int64, expense *planetscale.ExpenseUpdate) (*planetscale.Expense, error) {
	return s.UpdateFn(tx, expenseID, expense)
}

func (s ExpenseRepo) Find(tx *sql.Tx, filter planetscale.ExpenseFilter) ([]*planetscale.Expense, error) {
	return s.FindFn(tx, filter)
}
