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
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	group32, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	groupID := int64(group32)

	var settlements []*planetscale.Settlement
	getSettlementFunc := func(tx *sql.Tx) error {
		// validate user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, groupID, user.UserID)
		if err != nil {
			return err
		}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(findSettlementsResponse{
		Settlements: settlements,
		N:           len(settlements),
	}); err != nil {
		Error(w, r, err)
		return
	}
}

type findSettlementsResponse struct {
	Settlements []*planetscale.Settlement `json:"settlements"`
	N           int                       `json:"n"`
}

func (c *settlementController) HandlePostSettlement(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	var settlement planetscale.Settlement
	err := ReceiveJson(w, r, &settlement)
	if err != nil {
		Error(w, r, err)
		return
	}

	createSettlementFunc := func(tx *sql.Tx) error {
		// validate user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, settlement.GroupID, user.UserID)
		if err != nil {
			// rewrap error to add more context
			return planetscale.Errorf(planetscale.ENOTFOUND, "you are not a member of this group")
		}

		// validate paid by user is context user
		if settlement.PaidBy != user.UserID {
			return planetscale.Errorf(planetscale.EINVALID, "you cannot create a settlement for another user")
		}

		// validate paid to user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, settlement.GroupID, settlement.PaidTo)
		if err != nil {
			return planetscale.Errorf(planetscale.ENOTFOUND, "paid to user is not a member of this group")
		}

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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settlement); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *settlementController) HandleGetSettlement(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

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

		// validate user is a member of the group
		_, err := c.repos.GroupMember.Get(tx, settlement.GroupID, user.UserID)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settlement); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *settlementController) HandlePatchSettlement(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

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
		// validate user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, *settlement.GroupID, user.UserID)
		if err != nil {
			return err
		}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settlement); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *settlementController) HandleDeleteSettlement(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	settlement32, err := strconv.Atoi(chi.URLParam(r, "settlementID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	settlementID := int64(settlement32)

	deleteSettlementFunc := func(tx *sql.Tx) error {
		// validate user is a member of the group
		settlement, err := c.repos.Settlement.Get(tx, settlementID)
		if err != nil {
			return err
		}

		_, err = c.repos.GroupMember.Get(tx, settlement.GroupID, user.UserID)
		if err != nil {
			return err
		}

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

	w.WriteHeader(http.StatusNoContent)
}
