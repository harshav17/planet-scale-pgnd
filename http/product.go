package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
)

type productController struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewProductController(repos *planetscale.RepoProvider, tm planetscale.TransactionManager) *productController {
	return &productController{
		repos: repos,
		tm:    tm,
	}
}

func (c *productController) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	product32, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	productID := int64(product32)

	var product *planetscale.Product
	getProductFunc := func(tx *sql.Tx) error {
		product, err = c.repos.Product.Get(tx, productID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getProductFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(product); err != nil {
			Error(w, r, err)
			return
		}
	default:
		// TODO load template using go embed
		var data interface{}
		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			Error(w, r, err)
			return
		}
	}
}

func (c *productController) HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	var product planetscale.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		Error(w, r, err)
		return
	}

	createProductFunc := func(tx *sql.Tx) error {
		err := c.repos.Product.Create(tx, &product)
		if err != nil {
			return err
		}
		return nil
	}

	err := c.tm.ExecuteInTx(r.Context(), createProductFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		Error(w, r, err)
		return
	}
}
