package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type UserRepo struct {
	GetFn    func(tx *sql.Tx, userID string) (*planetscale.User, error)
	CreateFn func(tx *sql.Tx, user *planetscale.User) error
	UpsertFn func(tx *sql.Tx, user *planetscale.User) error
}

func (s UserRepo) Get(tx *sql.Tx, userID string) (*planetscale.User, error) {
	return s.GetFn(tx, userID)
}

func (s UserRepo) Create(tx *sql.Tx, user *planetscale.User) error {
	return s.CreateFn(tx, user)
}

func (s UserRepo) Upsert(tx *sql.Tx, user *planetscale.User) error {
	return s.UpsertFn(tx, user)
}
