package planetscale

import (
	"database/sql"
	"net/http"
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
		Upsert(tx *sql.Tx, user *User) error
	}

	UserController interface {
		HandlePutUser(w http.ResponseWriter, r *http.Request)
	}
)
