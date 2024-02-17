package planetscale

import "database/sql"

type (
	ItemSplitNU struct {
		ItemSplitID int64   `json:"item_split_id"`
		ItemID      int64   `json:"item_id"`
		UserID      *string `json:"user_id"`
		Amount      float64 `json:"amount"`
		Initials    *string `json:"initials"`
	}

	ItemSplitNURepo interface {
		Get(tx *sql.Tx, itemSplitID int64) (*ItemSplitNU, error)
		Create(tx *sql.Tx, itemSplit *ItemSplitNU) error
		Update(tx *sql.Tx, itemSplitID int64, itemSplit *ItemSplitNUUpdate) (*ItemSplitNU, error)
		Delete(tx *sql.Tx, itemSplitID int64) error
		Find(tx *sql.Tx, filter ItemSplitNUFilter) ([]*ItemSplitNU, error)
	}

	ItemSplitNUUpdate struct {
		Amount *float64 `json:"amount"`
	}

	ItemSplitNUFilter struct {
		ItemID int64 `json:"item_id"`
	}
)
