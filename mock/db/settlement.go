package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type SettlementRepo struct {
	GetFn    func(tx *sql.Tx, settlementID int64) (*planetscale.Settlement, error)
	CreateFn func(tx *sql.Tx, settlement *planetscale.Settlement) error
	DeleteFn func(tx *sql.Tx, settlementID int64) error
	UpdateFn func(tx *sql.Tx, settlementID int64, settlement *planetscale.SettlementUpdate) (*planetscale.Settlement, error)
	FindFn   func(tx *sql.Tx, filter planetscale.SettlementFilter) ([]*planetscale.Settlement, error)
}

func (s SettlementRepo) Get(tx *sql.Tx, settlementID int64) (*planetscale.Settlement, error) {
	return s.GetFn(tx, settlementID)
}

func (s SettlementRepo) Create(tx *sql.Tx, settlement *planetscale.Settlement) error {
	return s.CreateFn(tx, settlement)
}

func (s SettlementRepo) Delete(tx *sql.Tx, settlementID int64) error {
	return s.DeleteFn(tx, settlementID)
}

func (s SettlementRepo) Update(tx *sql.Tx, settlementID int64, settlement *planetscale.SettlementUpdate) (*planetscale.Settlement, error) {
	return s.UpdateFn(tx, settlementID, settlement)
}

func (s SettlementRepo) Find(tx *sql.Tx, filter planetscale.SettlementFilter) ([]*planetscale.Settlement, error) {
	return s.FindFn(tx, filter)
}
