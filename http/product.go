package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
	"github.com/harshav17/planet_scale/db"
)

type productController struct {
	repos *planetscale.RepoProvider
	tm    *db.TransactionManager
}

func NewProductController(repos *planetscale.RepoProvider, tm *db.TransactionManager) *productController {
	return &productController{
		repos: repos,
		tm:    tm,
	}
}

func (c *productController) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	product32, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productID := int64(product32)

	slog.Info("Hello from inside the story controller")

	var product *planetscale.Product
	getProductFunc := func(tx *sql.Tx) error {
		product, err = c.repos.Product.Get(tx, productID)
		if err != nil {
			slog.Error(err.Error())
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getProductFunc)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(product); err != nil {
			slog.Error(err.Error())
			return
		}
		w.Write([]byte(fmt.Sprintf("title:%v", product)))
	default:
		var tmplHtml = "index.html"
		tmpl, err := template.New(tmplHtml).ParseFiles(tmplHtml)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = tmpl.Execute(w, product)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func (c *productController) HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	var product planetscale.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var productID int64
	createProductFunc := func(tx *sql.Tx) error {
		err := c.repos.Product.Create(tx, &product)
		if err != nil {
			slog.Error(err.Error())
		}
		return nil
	}

	err := c.tm.ExecuteInTx(r.Context(), createProductFunc)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/products/%d", productID))
	w.WriteHeader(http.StatusCreated)
}
