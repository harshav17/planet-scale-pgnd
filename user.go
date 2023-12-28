package planetscale

import (
	"database/sql"
	"time"
)

type (
	User struct {
		UserID    string    `json:"user_id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	UserRepo interface {
		Get(tx *sql.Tx, userID string) (*User, error)
		Create(tx *sql.Tx, user *User) error
	}
)
