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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(findExpensesResponse{
		Expenses: expenses,
		N:        len(expenses),
	}); err != nil {
		Error(w, r, err)
		return
	}
}

type findExpensesResponse struct {
	Expenses []*planetscale.Expense `json:"expenses"`
	N        int                    `json:"n"`
}

func (c *expenseController) HandleGetExpense(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

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

		// check if user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, *expense.GroupID, user.UserID)
		if err != nil {
			return err
		}

		// get expense participants
		expense.Participants, err = c.repos.ExpenseParticipant.Find(tx, planetscale.ExpenseParticipantFilter{
			ExpenseID: expenseID,
		})

		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getExpenseFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expense); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *expenseController) HandlePostExpense(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	var expense planetscale.Expense
	err := ReceiveJson(w, r, &expense)
	if err != nil {
		Error(w, r, err)
		return
	}

	// add user id to expense
	expense.CreatedBy = user.UserID
	expense.UpdatedBy = user.UserID
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
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

	expense32, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	if err != nil {
		Error(w, r, err)
		return
	}
	expenseID := int64(expense32)

	deleteExpenseFunc := func(tx *sql.Tx) error {
		// check if the user is a member of the group
		expense, err := c.repos.Expense.Get(tx, expenseID)
		if err != nil {
			return err
		}
		_, err = c.repos.GroupMember.Get(tx, *expense.GroupID, user.UserID)
		if err != nil {
			return err
		}

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

	w.WriteHeader(http.StatusNoContent)
}

func (c *expenseController) HandlePatchExpense(w http.ResponseWriter, r *http.Request) {
	user, found := planetscale.UserFromContext(r.Context())
	if !found {
		Error(w, r, planetscale.Errorf(planetscale.ENOTFOUND, "user context not set"))
		return
	}

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
		// check if the user is a member of the group
		foundExp, err := c.repos.Expense.Get(tx, expenseID)
		if err != nil {
			return err
		}
		_, err = c.repos.GroupMember.Get(tx, *foundExp.GroupID, user.UserID)
		if err != nil {
			return err
		}

		expense, err = c.repos.Expense.Update(tx, expenseID, &expenseUpdate)
		if err != nil {
			return err
		}

		// update expense participants
		if expenseUpdate.Participants != nil {
			// TODO conver all the queries in this block to Batch queries
			// find existing participants
			existingParticipants, err := c.repos.ExpenseParticipant.Find(tx, planetscale.ExpenseParticipantFilter{
				ExpenseID: expenseID,
			})
			if err != nil {
				return err
			}

			// delete participants that are not in the update
			for _, existingParticipant := range existingParticipants {
				found := false
				for _, participant := range expenseUpdate.Participants {
					if existingParticipant.UserID == participant.UserID {
						found = true
						break
					}
				}
				if !found {
					err := c.repos.ExpenseParticipant.Delete(tx, expenseID, existingParticipant.UserID)
					if err != nil {
						return err
					}
				}
			}

			// upsert participants
			for _, participant := range expenseUpdate.Participants {
				err := c.repos.ExpenseParticipant.Upsert(tx, participant)
				if err != nil {
					return err
				}
			}

			// get updated participants
			expense.Participants, err = c.repos.ExpenseParticipant.Find(tx, planetscale.ExpenseParticipantFilter{
				ExpenseID: expenseID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), patchExpenseFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expense); err != nil {
		Error(w, r, err)
		return
	}
}
