package db

import (
	"database/sql"
	"fmt"

	planetscale "github.com/harshav17/planet_scale"
)

type splitTypeRepo struct {
	db *DB
}

func NewSplitTypeRepo(db *DB) *splitTypeRepo {
	return &splitTypeRepo{
		db: db,
	}
}

func (r *splitTypeRepo) Get(tx *sql.Tx, splitTypeID int64) (*planetscale.SplitType, error) {
	q := "SELECT * FROM split_types WHERE split_type_id = ?"

	var splitType planetscale.SplitType
	err := tx.QueryRow(q, splitTypeID).Scan(
		&splitType.SplitTypeID,
		&splitType.TypeName,
		&splitType.Description,
		(*NullTime)(&splitType.CreatedAt),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no split type found with ID %d", splitTypeID)
		}
		return nil, err
	}

	return &splitType, nil
}

func (r *splitTypeRepo) GetAll(tx *sql.Tx) ([]*planetscale.SplitType, error) {
	q := "SELECT * FROM split_types"

	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splitTypes []*planetscale.SplitType
	for rows.Next() {
		var splitType planetscale.SplitType
		err := rows.Scan(
			&splitType.SplitTypeID,
			&splitType.TypeName,
			&splitType.Description,
			(*NullTime)(&splitType.CreatedAt),
		)
		if err != nil {
			return nil, err
		}
		splitTypes = append(splitTypes, &splitType)
	}

	return splitTypes, nil
}
