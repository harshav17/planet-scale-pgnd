package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type expenseRepo struct {
	db *DB
}

func NewExpenseRepo(db *DB) *expenseRepo {
	return &expenseRepo{db}
}

func (r *expenseRepo) Get(tx *sql.Tx, expenseID int64) (*planetscale.Expense, error) {
	query := `
		SELECT 
			expense_id, 
			group_id, 
			paid_by, 
			amount, 
			description, 
			timestamp, 
			created_at, 
			updated_at, 
			created_by, 
			updated_by 
		FROM 
			expenses 
		WHERE 
			expense_id = ?`

	var expense planetscale.Expense
	row := tx.QueryRow(query, expenseID)
	err := row.Scan(&expense.ExpenseID, &expense.GroupID, &expense.PaidBy, &expense.Amount, &expense.Description, (*NullTime)(&expense.Timestamp), (*NullTime)(&expense.CreatedAt), (*NullTime)(&expense.UpdatedAt), &expense.CreatedBy, &expense.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no expense group found with ID %d", expenseID)
		}
		return nil, err
	}
	slog.Info("loaded expense", slog.Int64("id", expense.ExpenseID))

	return &expense, nil
}

func (r *expenseRepo) Create(tx *sql.Tx, expense *planetscale.Expense) error {
	query := `INSERT INTO expenses (group_id, paid_by, amount, description, timestamp, created_by, updated_by, split_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, expense.GroupID, expense.PaidBy, expense.Amount, expense.Description, (*NullTime)(&expense.Timestamp), expense.CreatedBy, expense.CreatedBy, expense.SplitTypeID)
	if err != nil {
		return err
	}
	expenseID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	expense.ExpenseID = expenseID
	slog.Info("created expense", slog.Int64("id", expense.ExpenseID))

	return nil
}

func (r *expenseRepo) Upsert(tx *sql.Tx, expense *planetscale.Expense) error {
	query := `
		INSERT INTO 
			expenses (
				group_id, 
				paid_by, 
				amount, 
				description, 
				timestamp, 
				created_by, 
				updated_by, 
				split_type_id
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
			ON DUPLICATE KEY UPDATE 
				group_id = ?, 
				paid_by = ?, 
				amount = ?, 
				description = ?, 
				timestamp = ?, 
				updated_by = ?, 
				split_type_id = ?`

	result, err := tx.Exec(
		query,
		expense.GroupID,
		expense.PaidBy,
		expense.Amount,
		expense.Description,
		(*NullTime)(&expense.Timestamp),
		expense.CreatedBy,
		expense.CreatedBy,
		expense.SplitTypeID,
		expense.GroupID,
		expense.PaidBy,
		expense.Amount,
		expense.Description,
		(*NullTime)(&expense.Timestamp),
		expense.CreatedBy,
		expense.SplitTypeID,
	)
	if err != nil {
		return err
	}
	expenseID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	expense.ExpenseID = expenseID
	slog.Info("upserted expense", slog.Int64("id", expense.ExpenseID))

	return nil
}

func (r *expenseRepo) Update(tx *sql.Tx, expenseID int64, update *planetscale.ExpenseUpdate) (*planetscale.Expense, error) {
	expense, err := r.Get(tx, expenseID)
	if err != nil {
		return nil, err
	}

	if update.GroupID != nil {
		expense.GroupID = *update.GroupID
	}
	if update.PaidBy != nil {
		expense.PaidBy = *update.PaidBy
	}
	if update.Amount != nil {
		expense.Amount = *update.Amount
	}
	if update.Description != nil {
		expense.Description = *update.Description
	}
	if update.Timestamp != nil {
		expense.Timestamp = *update.Timestamp
	}
	if update.UpdatedBy != nil {
		expense.UpdatedBy = *update.UpdatedBy
	}

	query := `UPDATE expenses SET group_id = ?, paid_by = ?, amount = ?, description = ?, timestamp = ?, updated_by = ? WHERE expense_id = ?`

	result, err := tx.Exec(query, expense.GroupID, expense.PaidBy, expense.Amount, expense.Description, expense.Timestamp, expense.UpdatedBy, expenseID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no expense found with ID %d", expenseID)
	}
	slog.Info("updated expense", slog.Int64("id", expenseID))

	return r.Get(tx, expenseID)
}

func (r *expenseRepo) Delete(tx *sql.Tx, expenseID int64) error {
	query := `DELETE FROM expenses WHERE expense_id = ?`

	result, err := tx.Exec(query, expenseID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no expense found with ID %d", expenseID)
	}
	slog.Info("deleted expense", slog.Int64("id", expenseID))

	return nil
}

func (r *expenseRepo) Find(tx *sql.Tx, filter planetscale.ExpenseFilter) ([]*planetscale.Expense, error) {
	where := &findWhereClause{}
	if filter.GroupID != 0 {
		where.Add("group_id", filter.GroupID)
	}

	query := `
		SELECT
			e.expense_id,
			e.group_id,
			e.paid_by,
			e.amount,
			e.description,
			e.timestamp,
			e.created_at,
			e.updated_at,
			e.created_by,
			e.updated_by,
			e.split_type_id,
			u.name
		FROM expenses e JOIN users u ON e.paid_by = u.user_id
		` + where.ToClause()
	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var expenses []*planetscale.Expense
	for rows.Next() {
		var expense planetscale.Expense
		var user planetscale.User
		err := rows.Scan(&expense.ExpenseID, &expense.GroupID, &expense.PaidBy, &expense.Amount, &expense.Description, (*NullTime)(&expense.Timestamp), (*NullTime)(&expense.CreatedAt), (*NullTime)(&expense.UpdatedAt), &expense.CreatedBy, &expense.UpdatedBy, &expense.SplitTypeID, &user.Name)
		if err != nil {
			return nil, err
		}
		expense.PaidByUser = &user

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}
