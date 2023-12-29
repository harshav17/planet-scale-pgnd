package planetscale

import (
	"database/sql"
	"time"
)

type (
	ExpenseGroup struct {
		ExpenseGroupID int64     `json:"group_id"`
		GroupName      string    `json:"group_name"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		CreateBy       string    `json:"created_by"`
		UpdatedBy      string    `json:"updated_by"`
	}

	ExpenseGroupRepo interface {
		Get(tx *sql.Tx, groupID int64) (*ExpenseGroup, error)
		Create(tx *sql.Tx, group *ExpenseGroup) error
		Update(tx *sql.Tx, groupID int64, update *ExpenseGroupUpdate) (*ExpenseGroup, error)
		Delete(tx *sql.Tx, groupID int64) error
		ListAllForUser(tx *sql.Tx, userID string) ([]*ExpenseGroup, error)
	}

	ExpenseGroupUpdate struct {
		GroupName string `json:"group_name"`
	}
)
