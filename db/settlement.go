package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type settlementRepo struct {
	db *DB
}

func NewSettlementRepo(db *DB) *settlementRepo {
	return &settlementRepo{
		db: db,
	}
}

func (r *settlementRepo) Get(tx *sql.Tx, settlementID int64) (*planetscale.Settlement, error) {
	query := `SELECT settlement_id, group_id, paid_by, paid_to, amount, timestamp FROM settlements WHERE settlement_id = ?`

	var settlement planetscale.Settlement
	row := tx.QueryRow(query, settlementID)
	err := row.Scan(&settlement.SettlementID, &settlement.GroupID, &settlement.PaidBy, &settlement.PaidTo, &settlement.Amount, (*NullTime)(&settlement.Timestamp))
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no settlement found with ID %d", settlementID)
		}
		return nil, err
	}
	slog.Info("loaded settlement", slog.Int64("id", settlement.SettlementID))

	return &settlement, nil
}

func (r *settlementRepo) Create(tx *sql.Tx, settlement *planetscale.Settlement) error {
	query := `INSERT INTO settlements (group_id, paid_by, paid_to, amount) VALUES (?, ?, ?, ?)`

	result, err := tx.Exec(query, settlement.GroupID, settlement.PaidBy, settlement.PaidTo, settlement.Amount)
	if err != nil {
		return err
	}
	settlement_id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	settlement.SettlementID = settlement_id
	slog.Info("created settlement", slog.Int64("id", settlement.SettlementID))

	return nil
}

func (r *settlementRepo) Delete(tx *sql.Tx, settlementID int64) error {
	query := `DELETE FROM settlements WHERE settlement_id = ?`

	result, err := tx.Exec(query, settlementID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no group member found with ID %d", settlementID)
	}
	slog.Info("deleted group member", slog.Int64("id", settlementID))

	return nil
}

func (r *settlementRepo) Update(tx *sql.Tx, settlementID int64, settlement *planetscale.SettlementUpdate) (*planetscale.Settlement, error) {
	query := `UPDATE settlements SET group_id = ? WHERE settlement_id = ?`

	result, err := tx.Exec(query, settlement.GroupID, settlementID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no group member found with ID %d", settlementID)
	}
	slog.Info("updated group member", slog.Int64("id", settlementID))

	return r.Get(tx, settlementID)
}

func (r *settlementRepo) Find(tx *sql.Tx, filter planetscale.SettlementFilter) ([]*planetscale.Settlement, error) {
	where := &findWhereClause{}
	if filter.GroupID != 0 {
		where.Add("group_id", filter.GroupID)
	}

	query := `
		SELECT
			settlement_id,
			group_id,
			paid_by,
			paid_to,
			Amount,
			timestamp
		FROM settlements
		` + where.ToClause()

	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var settlements []*planetscale.Settlement
	for rows.Next() {
		var settlement planetscale.Settlement
		err := rows.Scan(&settlement.SettlementID, &settlement.GroupID, &settlement.PaidBy, &settlement.PaidTo, &settlement.Amount, (*NullTime)(&settlement.Timestamp))
		if err != nil {
			return nil, err
		}
		settlements = append(settlements, &settlement)
	}

	return settlements, nil
}
