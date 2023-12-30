package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ExpenseGroupRepo struct {
	ListAllForUserFn func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error)
	GetFn            func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error)
	CreateFn         func(tx *sql.Tx, group *planetscale.ExpenseGroup) error
	UpdateFn         func(tx *sql.Tx, groupID int64, group *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error)
	DeleteFn         func(tx *sql.Tx, groupID int64) error
}

func (s ExpenseGroupRepo) ListAllForUser(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
	return s.ListAllForUserFn(tx, userID)
}

func (s ExpenseGroupRepo) Get(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
	return s.GetFn(tx, groupID)
}

func (s ExpenseGroupRepo) Create(tx *sql.Tx, group *planetscale.ExpenseGroup) error {
	return s.CreateFn(tx, group)
}

func (s ExpenseGroupRepo) Update(tx *sql.Tx, groupID int64, group *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error) {
	return s.UpdateFn(tx, groupID, group)
}

func (s ExpenseGroupRepo) Delete(tx *sql.Tx, groupID int64) error {
	return s.DeleteFn(tx, groupID)
}
