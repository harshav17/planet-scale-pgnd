package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ItemSplitNURepo struct {
	GetFn    func(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplitNU, error)
	CreateFn func(tx *sql.Tx, itemSplit *planetscale.ItemSplitNU) error
	UpdateFn func(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitNUUpdate) (*planetscale.ItemSplitNU, error)
	DeleteFn func(tx *sql.Tx, itemSplitID int64) error
	FindFn   func(tx *sql.Tx, filter planetscale.ItemSplitNUFilter) ([]*planetscale.ItemSplitNU, error)
}

func (s ItemSplitNURepo) Get(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplitNU, error) {
	return s.GetFn(tx, itemSplitID)
}

func (s ItemSplitNURepo) Create(tx *sql.Tx, itemSplit *planetscale.ItemSplitNU) error {
	return s.CreateFn(tx, itemSplit)
}

func (s ItemSplitNURepo) Update(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitNUUpdate) (*planetscale.ItemSplitNU, error) {
	return s.UpdateFn(tx, itemSplitID, itemSplit)
}

func (s ItemSplitNURepo) Delete(tx *sql.Tx, itemSplitID int64) error {
	return s.DeleteFn(tx, itemSplitID)
}

func (s ItemSplitNURepo) Find(tx *sql.Tx, filter planetscale.ItemSplitNUFilter) ([]*planetscale.ItemSplitNU, error) {
	return s.FindFn(tx, filter)
}
