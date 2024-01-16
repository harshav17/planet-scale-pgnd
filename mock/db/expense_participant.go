package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ExpenseParticipantRepo struct {
	GetFn    func(tx *sql.Tx, expenseID int64, userID string) (*planetscale.ExpenseParticipant, error)
	CreateFn func(tx *sql.Tx, expense *planetscale.ExpenseParticipant) error
	DeleteFn func(tx *sql.Tx, expenseID int64, userID string) error
	UpdateFn func(tx *sql.Tx, expenseID int64, userID string, expense *planetscale.ExpenseParticipantUpdate) (*planetscale.ExpenseParticipant, error)
	FindFn   func(tx *sql.Tx, filter planetscale.ExpenseParticipantFilter) ([]*planetscale.ExpenseParticipant, error)
}

func (s ExpenseParticipantRepo) Get(tx *sql.Tx, expenseID int64, userID string) (*planetscale.ExpenseParticipant, error) {
	return s.GetFn(tx, expenseID, userID)
}

func (s ExpenseParticipantRepo) Create(tx *sql.Tx, expense *planetscale.ExpenseParticipant) error {
	return s.CreateFn(tx, expense)
}

func (s ExpenseParticipantRepo) Delete(tx *sql.Tx, expenseID int64, userID string) error {
	return s.DeleteFn(tx, expenseID, userID)
}

func (s ExpenseParticipantRepo) Update(tx *sql.Tx, expenseID int64, userID string, expense *planetscale.ExpenseParticipantUpdate) (*planetscale.ExpenseParticipant, error) {
	return s.UpdateFn(tx, expenseID, userID, expense)
}

func (s ExpenseParticipantRepo) Find(tx *sql.Tx, filter planetscale.ExpenseParticipantFilter) ([]*planetscale.ExpenseParticipant, error) {
	return s.FindFn(tx, filter)
}
