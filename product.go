package planetscale

import (
	"database/sql"
	"net/http"
)

type (
	Product struct {
		ID    int64   `json:"ID"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	ProductRepo interface {
		Get(tx *sql.Tx, productID int64) (*Product, error)
		Create(tx *sql.Tx, product *Product) error
	}

	ProductController interface {
		HandleGetProduct(w http.ResponseWriter, r *http.Request)
		HandlePostProduct(w http.ResponseWriter, r *http.Request)
		HandleProductAdd(w http.ResponseWriter, r *http.Request)
	}
)
