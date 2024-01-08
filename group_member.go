package planetscale

import (
	"database/sql"
	"net/http"
	"time"
)

type (
	GroupMember struct {
		GroupID  int64     `json:"group_id"`
		UserID   string    `json:"user_id"`
		JoinedAt time.Time `json:"joined_at"`

		User *User `json:"user"`
	}

	GroupMemberRepo interface {
		Get(tx *sql.Tx, groupID int64, userID string) (*GroupMember, error)
		Create(tx *sql.Tx, group *GroupMember) error
		Delete(tx *sql.Tx, groupID int64, userID string) error
		Find(tx *sql.Tx, filter GroupMemberFilter) ([]*GroupMember, error)
	}

	GroupMemberUpdate struct {
		GroupName string `json:"group_name"`
	}

	GroupMemberFilter struct {
		GroupID int64
	}

	GroupMemberController interface {
		HandleGetGroupMembers(w http.ResponseWriter, r *http.Request)
		HandlePostGroupMember(w http.ResponseWriter, r *http.Request)
		HandleDeleteGroupMember(w http.ResponseWriter, r *http.Request)
	}
)
