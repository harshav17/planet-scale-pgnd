package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	planetscale "github.com/harshav17/planet_scale"
)

type itemController struct {
	repos    *planetscale.RepoProvider
	services *planetscale.ServiceProvider
	tm       planetscale.TransactionManager
}

func NewItemController(repos *planetscale.RepoProvider, services *planetscale.ServiceProvider, tm planetscale.TransactionManager) *itemController {
	return &itemController{
		repos:    repos,
		services: services,
		tm:       tm,
	}
}

func (c *itemController) HandlePostItem(w http.ResponseWriter, r *http.Request) {
	var item planetscale.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		Error(w, r, err)
		return
	}

	createItemFunc := func(tx *sql.Tx) error {
		err := c.repos.Item.Create(tx, &item)
		if err != nil {
			return err
		}

		for _, split := range item.Splits {
			split.ItemID = item.ItemID
			err = c.repos.ItemSplitNu.Create(tx, split)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), createItemFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
