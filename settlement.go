package planetscale

import (
	"database/sql"
	"net/http"
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
		Update(tx *sql.Tx, settlementID int64, settlement *SettlementUpdate) (*Settlement, error)
		Find(tx *sql.Tx, filter SettlementFilter) ([]*Settlement, error)
	}

	// what kind of fields can be updated?
	SettlementUpdate struct {
		GroupID *int64 `json:"group_id"`
	}

	SettlementFilter struct {
		GroupID int64
	}

	SettlementController interface {
		HandleGetSettlement(w http.ResponseWriter, r *http.Request)
		HandlePostSettlement(w http.ResponseWriter, r *http.Request)
		HandleDeleteSettlement(w http.ResponseWriter, r *http.Request)
		HandlePatchSettlement(w http.ResponseWriter, r *http.Request)
		HandleGetGroupSettlements(w http.ResponseWriter, r *http.Request)
	}
)
