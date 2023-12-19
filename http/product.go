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
		var data interface{}
		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			Error(w, r, err)
			return
		}
	}
}

func (c *productController) HandleProductAdd(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	err := templates.ExecuteTemplate(w, "product-add.html", data)
	if err != nil {
		Error(w, r, err)
		return
	}
}

func (c *productController) HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	var product *planetscale.Product
	var err error
	switch r.Header.Get("Accept") {
	case "application/json":
		product, err = handlePostProductJSON(r)
		if err != nil {
			Error(w, r, err)
			return
		}
	default:
		product, err = handlePostProductHTML(r)
		if err != nil {
			Error(w, r, err)
			return
		}
	}

	createProductFunc := func(tx *sql.Tx) error {
		err := c.repos.Product.Create(tx, product)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), createProductFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(product); err != nil {
			Error(w, r, err)
			return
		}
	default:
		products := []*planetscale.Product{product}
		err := templates.ExecuteTemplate(w, "products.html", products)
		if err != nil {
			Error(w, r, err)
			return
		}
	}
}

func handlePostProductHTML(r *http.Request) (*planetscale.Product, error) {
	product := &planetscale.Product{}
	r.ParseForm()
	product.Name = r.FormValue("name")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		return nil, err
	}
	product.Price = price

	return product, nil
}

func handlePostProductJSON(r *http.Request) (*planetscale.Product, error) {
	var product *planetscale.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}
