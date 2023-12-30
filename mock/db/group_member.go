package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type GroupMemberRepo struct {
	GetFn    func(tx *sql.Tx, groupID int64, userID string) (*planetscale.GroupMember, error)
	CreateFn func(tx *sql.Tx, group *planetscale.GroupMember) error
	DeleteFn func(tx *sql.Tx, groupID int64, userID string) error
	FindFn   func(tx *sql.Tx, filter planetscale.GroupMemberFilter) ([]*planetscale.GroupMember, error)
}

func (s GroupMemberRepo) Get(tx *sql.Tx, groupID int64, userID string) (*planetscale.GroupMember, error) {
	return s.GetFn(tx, groupID, userID)
}

func (s GroupMemberRepo) Create(tx *sql.Tx, group *planetscale.GroupMember) error {
	return s.CreateFn(tx, group)
}

func (s GroupMemberRepo) Delete(tx *sql.Tx, groupID int64, userID string) error {
	return s.DeleteFn(tx, groupID, userID)
}

func (s GroupMemberRepo) Find(tx *sql.Tx, filter planetscale.GroupMemberFilter) ([]*planetscale.GroupMember, error) {
	return s.FindFn(tx, filter)
}
