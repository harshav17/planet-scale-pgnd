package planetscale

import (
	"database/sql"
	"time"
)

type (
	Settlement struct {
		SettlementID int64     `json:"settlement_id"`
		GroupID      int64     `json:"group_id"`
		PaidBy       string    `json:"paid_by"`
		PaidTo       string    `json:"paid_to"`
		Amount       float64   `json:"amount"`
		Timestamp    time.Time `json:"timestamp"`
	}

	SettlementRepo interface {
		Get(tx *sql.Tx, settlementID int64) (*Settlement, error)
		Create(tx *sql.Tx, settlement *Settlement) error
		Delete(tx *sql.Tx, settlementID int64) error
		Update(tx *sql.Tx, settlementID int64, settlement *SettlementUpdate) error
	}

	// what kind of fields can be updated?
	SettlementUpdate struct {
		GroupID *int64 `json:"group_id"`
	}
)
