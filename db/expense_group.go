package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type expenseGroupRepo struct {
	db *DB
}

func NewExpenseGroupRepo(db *DB) *expenseGroupRepo {
	return &expenseGroupRepo{
		db: db,
	}
}

func (r *expenseGroupRepo) Get(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
	query := `SELECT group_id, group_name, created_at, created_by FROM expense_groups WHERE group_id = ?`

	var group planetscale.ExpenseGroup
	row := tx.QueryRow(query, groupID)
	err := row.Scan(&group.ExpenseGroupID, &group.GroupName, (*NullTime)(&group.CreatedAt), &group.CreateBy)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no expense group found with ID %d", groupID)
		}
		return nil, err
	}
	slog.Info("loaded expense group", slog.Int64("id", group.ExpenseGroupID))

	return &group, nil
}

func (r *expenseGroupRepo) Create(tx *sql.Tx, group *planetscale.ExpenseGroup) error {
	query := `INSERT INTO expense_groups (group_name, created_by) VALUES (?, ?)`

	result, err := tx.Exec(query, group.GroupName, group.CreateBy)
	if err != nil {
		return err
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	group.ExpenseGroupID = groupID
	slog.Info("created expense group", slog.Int64("id", group.ExpenseGroupID))

	return nil
}

func (r *expenseGroupRepo) Update(tx *sql.Tx, groupID int64, update *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error) {
	query := `UPDATE expense_groups SET group_name = ? WHERE group_id = ?`

	result, err := tx.Exec(query, update.GroupName, groupID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no group member found with ID %d", groupID)
	}
	slog.Info("updated expense group", slog.Int64("id", groupID))

	return r.Get(tx, groupID)
}

func (r *expenseGroupRepo) Delete(tx *sql.Tx, groupID int64) error {
	query := `DELETE FROM expense_groups WHERE group_id = ?`

	result, err := tx.Exec(query, groupID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no group member found with ID %d", groupID)
	}
	slog.Info("deleted expense group", slog.Int64("id", groupID))

	return nil
}
