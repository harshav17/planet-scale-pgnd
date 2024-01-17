package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ItemSplitRepo struct {
	GetFn    func(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplit, error)
	CreateFn func(tx *sql.Tx, itemSplit *planetscale.ItemSplit) error
	UpdateFn func(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitUpdate) (*planetscale.ItemSplit, error)
	DeleteFn func(tx *sql.Tx, itemSplitID int64) error
	FindFn   func(tx *sql.Tx, filter planetscale.ItemSplitFilter) ([]*planetscale.ItemSplit, error)
}

func (s ItemSplitRepo) Get(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplit, error) {
	return s.GetFn(tx, itemSplitID)
}

func (s ItemSplitRepo) Create(tx *sql.Tx, itemSplit *planetscale.ItemSplit) error {
	return s.CreateFn(tx, itemSplit)
}

func (s ItemSplitRepo) Update(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitUpdate) (*planetscale.ItemSplit, error) {
	return s.UpdateFn(tx, itemSplitID, itemSplit)
}

func (s ItemSplitRepo) Delete(tx *sql.Tx, itemSplitID int64) error {
	return s.DeleteFn(tx, itemSplitID)
}

func (s ItemSplitRepo) Find(tx *sql.Tx, filter planetscale.ItemSplitFilter) ([]*planetscale.ItemSplit, error) {
	return s.FindFn(tx, filter)
}
