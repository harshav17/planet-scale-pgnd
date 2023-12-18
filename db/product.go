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
			id,
			name,
			price
		FROM
			products
		WHERE
			id = ?
	`

	var product planetscale.Product
	row := tx.QueryRow(query, productID)
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no story found with ID %d", productID)
		}
		return nil, err
	}
	slog.Info("loaded product", slog.Int64("id", product.ID))

	return &product, nil
}

func (r *productRepo) Create(tx *sql.Tx, product *planetscale.Product) error {
	query := `
		INSERT INTO
			products
			(name, price)
		VALUES
			(?, ?)
	`

	result, err := tx.Exec(query, product.Name, product.Price)
	if err != nil {
		return err
	}
	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = productID
	slog.Info("created product", slog.Int64("id", product.ID))

	return nil
}
