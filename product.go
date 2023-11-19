package planetscale

import (
	"database/sql"
	"net/http"
)

type (
	Product struct {
		ID int64 `json:"ID"`
	}

	ProductRepo interface {
		Get(tx *sql.Tx, productID int64) (*Product, error)
	}

	ProductController interface {
		HandleGetProduct(w http.ResponseWriter, r *http.Request)
		HandlePostProduct(w http.ResponseWriter, r *http.Request)
	}
)
