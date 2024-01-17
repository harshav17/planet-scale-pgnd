package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ItemRepo struct {
	GetFn    func(tx *sql.Tx, itemID int64) (*planetscale.Item, error)
	CreateFn func(tx *sql.Tx, item *planetscale.Item) error
	DeleteFn func(tx *sql.Tx, itemID int64) error
	UpdateFn func(tx *sql.Tx, itemID int64, item *planetscale.ItemUpdate) (*planetscale.Item, error)
	FindFn   func(tx *sql.Tx, filter planetscale.ItemFilter) ([]*planetscale.Item, error)
}

func (s ItemRepo) Get(tx *sql.Tx, itemID int64) (*planetscale.Item, error) {
	return s.GetFn(tx, itemID)
}

func (s ItemRepo) Create(tx *sql.Tx, item *planetscale.Item) error {
	return s.CreateFn(tx, item)
}

func (s ItemRepo) Delete(tx *sql.Tx, itemID int64) error {
	return s.DeleteFn(tx, itemID)
}

func (s ItemRepo) Update(tx *sql.Tx, itemID int64, item *planetscale.ItemUpdate) (*planetscale.Item, error) {
	return s.UpdateFn(tx, itemID, item)
}

func (s ItemRepo) Find(tx *sql.Tx, filter planetscale.ItemFilter) ([]*planetscale.Item, error) {
	return s.FindFn(tx, filter)
}
