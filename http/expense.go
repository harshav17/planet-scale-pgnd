package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
)

type expenseController struct {
	repos    *planetscale.RepoProvider
	services *planetscale.ServiceProvider
	tm       planetscale.TransactionManager
}

func NewExpenseController(repos *planetscale.RepoProvider, services *planetscale.ServiceProvider, tm planetscale.TransactionManager) *expenseController {
	return &expenseController{
		repos:    repos,
		services: services,
		tm:       tm,
	}
}

func (c *expenseController) HandleGetGroupExpenses(w http.ResponseWriter, r *http.Request) {
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

	var expenses []*planetscale.Expense
	getExpenseFunc := func(tx *sql.Tx) error {
		// check if user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, groupID, user.UserID)
		if err != nil {
			return err
		}

		expenses, err = c.repos.Expense.Find(tx, planetscale.ExpenseFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getExpenseFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(findExpensesResponse{
			Expenses: expenses,
			N:        len(expenses),
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

type findExpensesResponse struct {
	Expenses []*planetscale.Expense `json:"expenses"`
	N        int                    `json:"n"`
}

func (c *expenseController) HandleGetExpense(w http.ResponseWriter, r *http.Request) {
	expense32, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	expenseID := int64(expense32)

	var expense *planetscale.Expense
	getExpenseFunc := func(tx *sql.Tx) error {
		expense, err = c.repos.Expense.Get(tx, expenseID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getExpenseFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expense); err != nil {
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

func (c *expenseController) HandlePostExpense(w http.ResponseWriter, r *http.Request) {
	var expense planetscale.Expense
	err := ReceiveJson(w, r, &expense)
	if err != nil {
		Error(w, r, err)
		return
	}

	err = c.services.Expense.CreateExpense(r.Context(), &expense)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expense); err != nil {
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

func (c *expenseController) HandleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	expense32, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	expenseID := int64(expense32)

	deleteExpenseFunc := func(tx *sql.Tx) error {
		err = c.repos.Expense.Delete(tx, expenseID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), deleteExpenseFunc)
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

func (c *expenseController) HandlePatchExpense(w http.ResponseWriter, r *http.Request) {
	expense32, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	expenseID := int64(expense32)

	var expenseUpdate planetscale.ExpenseUpdate
	err = ReceiveJson(w, r, &expenseUpdate)
	if err != nil {
		Error(w, r, err)
		return
	}

	var expense *planetscale.Expense
	patchExpenseFunc := func(tx *sql.Tx) error {
		expense, err = c.repos.Expense.Update(tx, expenseID, &expenseUpdate)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), patchExpenseFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expense); err != nil {
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
