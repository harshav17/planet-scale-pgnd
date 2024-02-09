package http

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	planetscale "github.com/harshav17/planet_scale"
	svix "github.com/svix/svix-webhooks/go"
)

type userController struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
	wh    *svix.Webhook
}

func NewUserController(repos *planetscale.RepoProvider, tm planetscale.TransactionManager, wh *svix.Webhook) *userController {
	return &userController{
		repos: repos,
		tm:    tm,
		wh:    wh,
	}
}

func (c *userController) HandlePutUser(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		Error(w, r, err)
		return
	}

	err = c.wh.Verify(payload, r.Header)
	if err != nil {
		Error(w, r, err)
		return
	}

	var response planetscale.ClerkPayload[planetscale.ClerkUserPayload]
	err = json.Unmarshal(payload, &response)
	if err != nil {
		Error(w, r, err)
		return
	}

	// convert the clerk payload to a user
	user := planetscale.User{
		UserID: response.Data.Id,
		Email:  response.Data.EmailAddresses[0].EmailAddress,
		Name:   response.Data.FirstName + " " + response.Data.LastName,
	}

	upsertUserFunc := func(tx *sql.Tx) error {
		err := c.repos.User.Upsert(tx, &user)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), upsertUserFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
