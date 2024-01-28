package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type expenseParticipantRepo struct {
	db *DB
}

func NewExpenseParticipantRepo(db *DB) *expenseParticipantRepo {
	return &expenseParticipantRepo{
		db: db,
	}
}

func (r *expenseParticipantRepo) Get(tx *sql.Tx, expenseID int64, userID string) (*planetscale.ExpenseParticipant, error) {
	query := `
		SELECT
			expense_id,
			user_id,
			amount_owed,
			share_percentage,
			note
		FROM expense_participants
		WHERE expense_id = ? AND user_id = ?
	`

	var participant planetscale.ExpenseParticipant
	row := tx.QueryRow(query, expenseID, userID)
	err := row.Scan(&participant.ExpenseID, &participant.UserID, &participant.AmountOwed, &participant.SharePercentage, &participant.Note)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no expense participant found with expenseID %d and userID %s", expenseID, userID)
		}
		return nil, err
	}
	slog.Info("loaded expense participant", slog.Int64("id", expenseID), slog.String("user_id", userID))

	return &participant, nil
}

func (r *expenseParticipantRepo) Create(tx *sql.Tx, participant *planetscale.ExpenseParticipant) error {
	query := `
		INSERT INTO expense_participants (
			expense_id,
			user_id,
			amount_owed,
			share_percentage,
			note
		) VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query, participant.ExpenseID, participant.UserID, participant.AmountOwed, participant.SharePercentage, participant.Note)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	slog.Info("created expense participant", slog.Int64("id", participant.ExpenseID), slog.String("user_id", participant.UserID))

	return nil
}

func (r *expenseParticipantRepo) Upsert(tx *sql.Tx, participant *planetscale.ExpenseParticipant) error {
	query := `
		INSERT INTO expense_participants (
			expense_id,
			user_id,
			amount_owed,
			share_percentage,
			note
		) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE amount_owed = ?, share_percentage = ?, note = ?
	`

	result, err := tx.Exec(query, participant.ExpenseID, participant.UserID, participant.AmountOwed, participant.SharePercentage, participant.Note, participant.AmountOwed, participant.SharePercentage, participant.Note)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	slog.Info("upserted expense participant", slog.Int64("id", participant.ExpenseID), slog.String("user_id", participant.UserID))

	return nil
}

func (r *expenseParticipantRepo) Delete(tx *sql.Tx, expenseID int64, userID string) error {
	query := `
		DELETE FROM expense_participants
		WHERE expense_id = ? AND user_id = ?
	`

	result, err := tx.Exec(query, expenseID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no expense participant found with expenseID %d and userID %s", expenseID, userID)
	}
	slog.Info("deleted expense participant", slog.Int64("id", expenseID), slog.String("user_id", userID))

	return nil
}

func (r *expenseParticipantRepo) Update(tx *sql.Tx, expenseID int64, userID string, update *planetscale.ExpenseParticipantUpdate) (*planetscale.ExpenseParticipant, error) {
	participant, err := r.Get(tx, expenseID, userID)
	if err != nil {
		return nil, err
	}

	if update.AmountOwed != nil {
		participant.AmountOwed = *update.AmountOwed
	}
	if update.SharePercentage != nil {
		participant.SharePercentage = *update.SharePercentage
	}
	if update.Note != nil {
		participant.Note = *update.Note
	}

	query := `
		UPDATE expense_participants
		SET amount_owed = ?, share_percentage = ?, note = ?
		WHERE expense_id = ? AND user_id = ?
	`

	result, err := tx.Exec(query, participant.AmountOwed, participant.SharePercentage, participant.Note, expenseID, userID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no expense participant found with expenseID %d and userID %s", expenseID, userID)
	}
	slog.Info("updated expense participant", slog.Int64("id", expenseID), slog.String("user_id", userID))

	return r.Get(tx, expenseID, userID)
}

func (r *expenseParticipantRepo) Find(tx *sql.Tx, filter planetscale.ExpenseParticipantFilter) ([]*planetscale.ExpenseParticipant, error) {
	where := &findWhereClause{}
	if filter.ExpenseID != 0 {
		where.Add("expense_id", filter.ExpenseID)
	}

	query := `
		SELECT
			expense_id,
			user_id,
			amount_owed,
			share_percentage,
			note
		FROM expense_participants
	` + where.ToClause()

	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var participants []*planetscale.ExpenseParticipant
	for rows.Next() {
		var participant planetscale.ExpenseParticipant
		err := rows.Scan(&participant.ExpenseID, &participant.UserID, &participant.AmountOwed, &participant.SharePercentage, &participant.Note)
		if err != nil {
			return nil, err
		}
		participants = append(participants, &participant)
	}
	slog.Info("found expense participants", slog.Int("count", len(participants)))

	return participants, nil
}
