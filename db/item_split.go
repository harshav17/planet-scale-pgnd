package db

import (
	"database/sql"
	"fmt"

	planetscale "github.com/harshav17/planet_scale"
)

type itemSplitRepo struct {
	db *DB
}

func NewItemSplitRepo(db *DB) *itemSplitRepo {
	return &itemSplitRepo{
		db: db,
	}
}

func (r *itemSplitRepo) Get(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplit, error) {
	query := `
		SELECT
			item_split_id,
			item_id,
			user_id,
			amount
		FROM
			item_splits
		WHERE
			item_split_id = ?
	`

	var itemSplit planetscale.ItemSplit
	row := tx.QueryRow(query, itemSplitID)
	err := row.Scan(&itemSplit.ItemSplitID, &itemSplit.ItemID, &itemSplit.UserID, &itemSplit.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no item split found with ID %d", itemSplitID)
		}
		return nil, err
	}
	return &itemSplit, nil
}

func (r *itemSplitRepo) Create(tx *sql.Tx, itemSplit *planetscale.ItemSplit) error {
	query := `
		INSERT INTO
			item_splits
			(item_id, user_id, amount)
		VALUES
			(?, ?, ?)
	`

	result, err := tx.Exec(query, itemSplit.ItemID, itemSplit.UserID, itemSplit.Amount)
	if err != nil {
		return err
	}
	itemSplitID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	itemSplit.ItemSplitID = itemSplitID
	return nil
}

func (r *itemSplitRepo) Update(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitUpdate) (*planetscale.ItemSplit, error) {
	query := `
		UPDATE
			item_splits
		SET
			amount = ?
		WHERE
			item_split_id = ?
	`

	result, err := tx.Exec(query, itemSplit.Amount, itemSplitID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no item split found with ID %d", itemSplitID)
	}
	return r.Get(tx, itemSplitID)
}

func (r *itemSplitRepo) Delete(tx *sql.Tx, itemSplitID int64) error {
	query := `
		DELETE FROM
			item_splits
		WHERE
			item_split_id = ?
	`

	result, err := tx.Exec(query, itemSplitID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no item split found with ID %d", itemSplitID)
	}
	return nil
}

func (r *itemSplitRepo) Find(tx *sql.Tx, filter planetscale.ItemSplitFilter) ([]*planetscale.ItemSplit, error) {
	where := &findWhereClause{}
	if filter.ItemID != 0 {
		where.Add("item_id", filter.ItemID)
	}

	query := `
		SELECT
			item_split_id,
			item_id,
			user_id,
			amount
		FROM
			item_splits
	` + where.ToClause()

	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var itemSplits []*planetscale.ItemSplit
	for rows.Next() {
		var itemSplit planetscale.ItemSplit
		if err := rows.Scan(&itemSplit.ItemSplitID, &itemSplit.ItemID, &itemSplit.UserID, &itemSplit.Amount); err != nil {
			return nil, err
		}
		itemSplits = append(itemSplits, &itemSplit)
	}
	return itemSplits, nil
}
