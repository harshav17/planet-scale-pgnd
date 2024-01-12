package planetscale

import (
	"database/sql"
	"time"
)

type (
	SplitType struct {
		SplitTypeID int64     `json:"split_type_id"`
		TypeName    string    `json:"type_name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
	}

	SplitTypeRepo interface {
		Get(tx *sql.Tx, splitTypeID int64) (*SplitType, error)
		GetAll(tx *sql.Tx) ([]*SplitType, error)
	}
)
