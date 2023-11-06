package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type productRepo struct {
	db *DB
}

func NewProductRepo(db *DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Get(tx *sql.Tx, productID int64) (*planetscale.Product, error) {
	query := `
		SELECT
			id
		FROM
			products
		WHERE
			id = ?
	`

	slog.Info("loading all products")
	var product planetscale.Product
	row := tx.QueryRow(query, productID)
	err := row.Scan(&product.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no story found with ID %d", productID)
		}
		return nil, err
	}
	slog.Info("loaded product", "productID", product.ID)

	return &product, nil
}
