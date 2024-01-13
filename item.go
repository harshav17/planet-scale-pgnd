package planetscale

import "database/sql"

type (
	Item struct {
		ItemID    int64   `json:"item_id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Quantity  int64   `json:"quantity"`
		ExpenseID int64   `json:"expense_id"`
	}

	ItemRepo interface {
		Get(tx *sql.Tx, itemID int64) (*Item, error)
		Create(tx *sql.Tx, item *Item) error
		Update(tx *sql.Tx, itemID int64, item *ItemUpdate) (*Item, error)
		Delete(tx *sql.Tx, itemID int64) error
	}

	ItemUpdate struct {
		Name     *string  `json:"name"`
		Price    *float64 `json:"price"`
		Quantity *int64   `json:"quantity"`
	}
)
