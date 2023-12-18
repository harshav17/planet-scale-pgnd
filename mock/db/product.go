package db_mock

import (
	"database/sql"

	planetscale "github.com/harshav17/planet_scale"
)

type ProductRepo struct {
	GetFn    func(tx *sql.Tx, storyID int64) (*planetscale.Product, error)
	GetAllFn func(tx *sql.Tx) ([]*planetscale.Product, error)
	CreateFn func(tx *sql.Tx, story *planetscale.Product) error
}

func (s ProductRepo) Get(tx *sql.Tx, storyID int64) (*planetscale.Product, error) {
	return s.GetFn(tx, storyID)
}

func (s ProductRepo) GetAll(tx *sql.Tx) ([]*planetscale.Product, error) {
	return s.GetAllFn(tx)
}

func (s ProductRepo) Create(tx *sql.Tx, story *planetscale.Product) error {
	return s.CreateFn(tx, story)
}
