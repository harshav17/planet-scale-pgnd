package http

import (
	"database/sql"
	"fmt"
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

	c.tm.ExecuteInTx(r.Context(), getProductFunc)

	w.Write([]byte(fmt.Sprintf("title:%v", product)))
}
