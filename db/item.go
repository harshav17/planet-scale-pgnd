package db

import (
	"database/sql"
	"fmt"

	planetscale "github.com/harshav17/planet_scale"
)

type itemRepo struct {
	db *DB
}

func NewItemRepo(db *DB) *itemRepo {
	return &itemRepo{
		db: db,
	}
}

func (r *itemRepo) Get(tx *sql.Tx, itemID int64) (*planetscale.Item, error) {
	query := `
		SELECT
			item_id,
			name,
			price,
			quantity,
			expense_id
		FROM
			items
		WHERE
			item_id = ?
	`

	var item planetscale.Item
	row := tx.QueryRow(query, itemID)
	err := row.Scan(&item.ItemID, &item.Name, &item.Price, &item.Quantity, &item.ExpenseID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no item found with ID %d", itemID)
		}
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) Create(tx *sql.Tx, item *planetscale.Item) error {
	query := `
		INSERT INTO
			items
			(name, price, quantity, expense_id)
		VALUES
			(?, ?, ?, ?)
	`

	result, err := tx.Exec(query, item.Name, item.Price, item.Quantity, item.ExpenseID)
	if err != nil {
		return err
	}
	itemID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ItemID = itemID
	return nil
}

func (r *itemRepo) Update(tx *sql.Tx, itemID int64, update *planetscale.ItemUpdate) (*planetscale.Item, error) {
	query := `
		UPDATE
			items
		SET
			name = ?,
			price = ?,
			quantity = ?
		WHERE
			item_id = ?
	`

	result, err := tx.Exec(query, update.Name, update.Price, update.Quantity, itemID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no item found with ID %d", itemID)
	}
	return r.Get(tx, itemID)
}

func (r *itemRepo) Delete(tx *sql.Tx, itemID int64) error {
	query := `DELETE FROM items WHERE item_id = ?`

	result, err := tx.Exec(query, itemID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no item found with ID %d", itemID)
	}
	return nil
}

func (r *itemRepo) Find(tx *sql.Tx, filter planetscale.ItemFilter) ([]*planetscale.Item, error) {
	where := &findWhereClause{}

	if filter.ExpenseID != 0 {
		where.Add("expense_id", filter.ExpenseID)
	}

	query := `
		SELECT
			item_id,
			name,
			price,
			quantity,
			expense_id
		FROM
			items
	` + where.ToClause()

	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var items []*planetscale.Item
	for rows.Next() {
		var item planetscale.Item
		err := rows.Scan(&item.ItemID, &item.Name, &item.Price, &item.Quantity, &item.ExpenseID)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
