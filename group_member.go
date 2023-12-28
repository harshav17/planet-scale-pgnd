package planetscale

import (
	"database/sql"
	"time"
)

type (
	GroupMember struct {
		GroupID  int64     `json:"group_id"`
		UserID   string    `json:"user_id"`
		JoinedAt time.Time `json:"joined_at"`
	}

	GroupMemberRepo interface {
		Get(tx *sql.Tx, groupID int64, userID string) (*GroupMember, error)
		Create(tx *sql.Tx, group *GroupMember) error
		Delete(tx *sql.Tx, groupID int64, userID string) error
	}

	GroupMemberUpdate struct {
		GroupName string `json:"group_name"`
	}
)
