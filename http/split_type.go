package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	planetscale "github.com/harshav17/planet_scale"
)

type splitTypeController struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewSplitTypeController(repos *planetscale.RepoProvider, tm planetscale.TransactionManager) *splitTypeController {
	return &splitTypeController{
		repos: repos,
		tm:    tm,
	}
}

func (c *splitTypeController) HandleGetAllSplitTypes(w http.ResponseWriter, r *http.Request) {
	var splitTypes []*planetscale.SplitType
	var err error
	getSplitTypesFunc := func(tx *sql.Tx) error {
		splitTypes, err = c.repos.SplitType.GetAll(tx)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getSplitTypesFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(splitTypes); err != nil {
		Error(w, r, err)
		return
	}
}
