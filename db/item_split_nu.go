package db

import (
	"database/sql"
	"fmt"

	planetscale "github.com/harshav17/planet_scale"
)

type itemSplitNURepo struct {
	db *DB
}

func NewItemSplitNURepo(db *DB) *itemSplitNURepo {
	return &itemSplitNURepo{
		db: db,
	}
}

func (r *itemSplitNURepo) Get(tx *sql.Tx, itemSplitID int64) (*planetscale.ItemSplitNU, error) {
	query := `
		SELECT
			item_split_id,
			item_id,
			user_id,
			amount
		FROM
			item_splits_nu
		WHERE
			item_split_id = ?
	`

	var itemSplit planetscale.ItemSplitNU
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

func (r *itemSplitNURepo) Create(tx *sql.Tx, itemSplit *planetscale.ItemSplitNU) error {
	query := `
		INSERT INTO
			item_splits_nu
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

func (r *itemSplitNURepo) Update(tx *sql.Tx, itemSplitID int64, itemSplit *planetscale.ItemSplitNUUpdate) (*planetscale.ItemSplitNU, error) {
	query := `
		UPDATE
			item_splits_nu
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

func (r *itemSplitNURepo) Delete(tx *sql.Tx, itemSplitID int64) error {
	query := `
		DELETE FROM
			item_splits_nu
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

func (r *itemSplitNURepo) Find(tx *sql.Tx, filter planetscale.ItemSplitNUFilter) ([]*planetscale.ItemSplitNU, error) {
	where := &findWhereClause{}
	if filter.ItemID != 0 {
		where.Add("item_id", filter.ItemID)
	}

	query := `
		SELECT
			item_split_id,
			item_id,
			user_id,
			amount,
			initials
		FROM
			item_splits_nu
	` + where.ToClause()

	rows, err := tx.Query(query, where.values...)
	if err != nil {
		return nil, err
	}

	var itemSplits []*planetscale.ItemSplitNU
	for rows.Next() {
		var itemSplit planetscale.ItemSplitNU
		err := rows.Scan(&itemSplit.ItemSplitID, &itemSplit.ItemID, &itemSplit.UserID, &itemSplit.Amount, &itemSplit.Initials)
		if err != nil {
			return nil, err
		}
		itemSplits = append(itemSplits, &itemSplit)
	}
	return itemSplits, nil
}
