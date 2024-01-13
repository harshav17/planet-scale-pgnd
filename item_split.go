package planetscale

import "database/sql"

type (
	ItemSplit struct {
		ItemSplitID int64   `json:"item_split_id"`
		ItemID      int64   `json:"item_id"`
		UserID      string  `json:"user_id"`
		Amount      float64 `json:"amount"`
	}

	ItemSplitRepo interface {
		Get(tx *sql.Tx, itemSplitID int64) (*ItemSplit, error)
		Create(tx *sql.Tx, itemSplit *ItemSplit) error
		Update(tx *sql.Tx, itemSplitID int64, itemSplit *ItemSplitUpdate) (*ItemSplit, error)
		Delete(tx *sql.Tx, itemSplitID int64) error
		Find(tx *sql.Tx, filter ItemSplitFilter) ([]*ItemSplit, error)
	}

	ItemSplitUpdate struct {
		Amount *float64 `json:"amount"`
	}

	ItemSplitFilter struct {
		ItemID int64 `json:"item_id"`
	}
)
