package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type userRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Get(tx *sql.Tx, userID string) (*planetscale.User, error) {
	query := `SELECT user_id, email, name, created_at FROM users WHERE user_id = ?`

	var user planetscale.User
	row := tx.QueryRow(query, userID)
	err := row.Scan(&user.UserID, &user.Email, &user.Name, (*NullTime)(&user.CreatedAt))
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, err
	}
	slog.Info("loaded user", slog.String("id", user.UserID))

	return &user, nil
}

func (r *userRepo) Create(tx *sql.Tx, user *planetscale.User) error {
	query := `INSERT INTO users (user_id, email, name) VALUES (?, ?, ?)`

	result, err := tx.Exec(query, user.UserID, user.Email, user.Name)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	slog.Info("created user", slog.String("id", user.UserID))

	return nil
}
