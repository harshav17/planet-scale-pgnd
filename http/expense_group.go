package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
)

type expenseGroupController struct {
	repos    *planetscale.RepoProvider
	services *planetscale.ServiceProvider
	tm       planetscale.TransactionManager
}

func NewExpenseGroupController(repos *planetscale.RepoProvider, services *planetscale.ServiceProvider, tm planetscale.TransactionManager) *expenseGroupController {
	return &expenseGroupController{
		repos:    repos,
		services: services,
		tm:       tm,
	}
}

func (c *expenseGroupController) HandleGetExpenseGroups(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	var expenseGroups []*planetscale.ExpenseGroup
	var err error
	getExpenseGroupFunc := func(tx *sql.Tx) error {
		expenseGroups, err = c.repos.ExpenseGroup.ListAllForUser(tx, user.UserID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getExpenseGroupFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(findExpenseGroupsResponse{
		ExpenseGroups: expenseGroups,
		N:             len(expenseGroups),
	}); err != nil {
		Error(w, r, err)
		return
	}
}

type findExpenseGroupsResponse struct {
	ExpenseGroups []*planetscale.ExpenseGroup `json:"expenseGroups"`
	N             int                         `json:"n"`
}

func (c *expenseGroupController) HandlePostExpenseGroup(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	var expenseGroup planetscale.ExpenseGroup
	err := ReceiveJson(w, r, &expenseGroup)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Set the user ID of the expense group to the user ID of the user who created it.
	expenseGroup.CreateBy = user.UserID
	expenseGroup.UpdatedBy = user.UserID

	createExpenseGroupFunc := func(tx *sql.Tx) error {
		err = c.repos.ExpenseGroup.Create(tx, &expenseGroup)
		if err != nil {
			return err
		}

		err = c.repos.GroupMember.Create(tx, &planetscale.GroupMember{
			GroupID: expenseGroup.ExpenseGroupID,
			UserID:  expenseGroup.CreateBy,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), createExpenseGroupFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expenseGroup); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *expenseGroupController) HandlePatchExpenseGroup(w http.ResponseWriter, r *http.Request) {
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

	var update planetscale.ExpenseGroupUpdate
	err = ReceiveJson(w, r, &update)
	if err != nil {
		Error(w, r, err)
		return
	}

	var expenseGroup *planetscale.ExpenseGroup
	patchExpenseGroupFunc := func(tx *sql.Tx) error {
		expenseGroup, err = c.repos.ExpenseGroup.Get(tx, groupID)
		if err != nil {
			return err
		}
		if expenseGroup.CreateBy != user.UserID {
			// TODO: should the users of the group be able to update the group?
			return planetscale.Errorf(planetscale.EUNAUTHORIZED, "user %s is not authorized to update expense group %d", user.UserID, groupID)
		}

		expenseGroup, err = c.repos.ExpenseGroup.Update(tx, groupID, &update)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), patchExpenseGroupFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expenseGroup); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *expenseGroupController) HandleDeleteExpenseGroup(w http.ResponseWriter, r *http.Request) {
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

	deleteExpenseGroupFunc := func(tx *sql.Tx) error {
		expenseGroup, err := c.repos.ExpenseGroup.Get(tx, groupID)
		if err != nil {
			return err
		}
		if expenseGroup.CreateBy != user.UserID {
			return planetscale.Errorf(planetscale.EUNAUTHORIZED, "user %s is not authorized to update expense group %d", user.UserID, groupID)
		}

		err = c.repos.ExpenseGroup.Delete(tx, groupID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), deleteExpenseGroupFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *expenseGroupController) HandleGetExpenseGroup(w http.ResponseWriter, r *http.Request) {
	group32, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	groupID := int64(group32)

	var expenseGroup *planetscale.ExpenseGroup
	getExpenseGroupFunc := func(tx *sql.Tx) error {
		expenseGroup, err = c.repos.ExpenseGroup.Get(tx, groupID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getExpenseGroupFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expenseGroup); err != nil {
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

func (c *expenseGroupController) HandleGetGroupBalances(w http.ResponseWriter, r *http.Request) {
	group32, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	groupID := int64(group32)

	balances, err := c.services.Balance.GetGroupBalances(r.Context(), groupID)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(balances); err != nil {
		Error(w, r, err)
		return
	}
}
