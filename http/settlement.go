package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
)

type settlementController struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewSettlementController(repos *planetscale.RepoProvider, tm planetscale.TransactionManager) *settlementController {
	return &settlementController{
		repos: repos,
		tm:    tm,
	}
}

// HandleGetGroupSettlements handles the GET /groups/{groupID}/settlements endpoint.
func (c *settlementController) HandleGetGroupSettlements(w http.ResponseWriter, r *http.Request) {
	group32, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	groupID := int64(group32)

	var settlements []*planetscale.Settlement
	getSettlementFunc := func(tx *sql.Tx) error {
		settlements, err = c.repos.Settlement.Find(tx, planetscale.SettlementFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getSettlementFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(findSettlementsResponse{
			Settlements: settlements,
			N:           len(settlements),
		}); err != nil {
			Error(w, r, err)
			return
		}
	default:
		Error(w, r, &planetscale.Error{
			Code:    planetscale.ENOTIMPLEMENTED,
			Message: "not implemented",
		})
		return
	}
}

type findSettlementsResponse struct {
	Settlements []*planetscale.Settlement `json:"settlements"`
	N           int                       `json:"n"`
}

func (c *settlementController) HandlePostSettlement(w http.ResponseWriter, r *http.Request) {
	var settlement planetscale.Settlement
	err := ReceiveJson(w, r, &settlement)
	if err != nil {
		Error(w, r, err)
		return
	}

	createSettlementFunc := func(tx *sql.Tx) error {
		err = c.repos.Settlement.Create(tx, &settlement)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), createSettlementFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(settlement); err != nil {
			Error(w, r, err)
			return
		}
	default:
		Error(w, r, &planetscale.Error{
			Code:    planetscale.ENOTIMPLEMENTED,
			Message: "not implemented",
		})
		return
	}
}

func (c *settlementController) HandleGetSettlement(w http.ResponseWriter, r *http.Request) {
	settlement32, err := strconv.Atoi(chi.URLParam(r, "settlementID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	settlementID := int64(settlement32)

	var settlement *planetscale.Settlement
	getSettlementFunc := func(tx *sql.Tx) error {
		settlement, err = c.repos.Settlement.Get(tx, settlementID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getSettlementFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(settlement); err != nil {
			Error(w, r, err)
			return
		}
	default:
		Error(w, r, &planetscale.Error{
			Code:    planetscale.ENOTIMPLEMENTED,
			Message: "not implemented",
		})
		return
	}
}

func (c *settlementController) HandlePatchSettlement(w http.ResponseWriter, r *http.Request) {
	settlement32, err := strconv.Atoi(chi.URLParam(r, "settlementID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	settlementID := int64(settlement32)

	var settlement planetscale.SettlementUpdate
	err = ReceiveJson(w, r, &settlement)
	if err != nil {
		Error(w, r, err)
		return
	}

	updateSettlementFunc := func(tx *sql.Tx) error {
		_, err = c.repos.Settlement.Update(tx, settlementID, &settlement)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), updateSettlementFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(settlement); err != nil {
			Error(w, r, err)
			return
		}
	default:
		Error(w, r, &planetscale.Error{
			Code:    planetscale.ENOTIMPLEMENTED,
			Message: "not implemented",
		})
		return
	}
}

func (c *settlementController) HandleDeleteSettlement(w http.ResponseWriter, r *http.Request) {
	settlement32, err := strconv.Atoi(chi.URLParam(r, "settlementID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	settlementID := int64(settlement32)

	deleteSettlementFunc := func(tx *sql.Tx) error {
		err = c.repos.Settlement.Delete(tx, settlementID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), deleteSettlementFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.WriteHeader(http.StatusNoContent)
	default:
		Error(w, r, &planetscale.Error{
			Code:    planetscale.ENOTIMPLEMENTED,
			Message: "not implemented",
		})
		return
	}
}
