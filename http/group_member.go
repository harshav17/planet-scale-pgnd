package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
)

type groupMemberController struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
}

func NewGroupMemberController(repos *planetscale.RepoProvider, tm planetscale.TransactionManager) *groupMemberController {
	return &groupMemberController{repos: repos, tm: tm}
}

func (c *groupMemberController) HandleGetGroupMembers(w http.ResponseWriter, r *http.Request) {
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

	var groupMembers []*planetscale.GroupMember
	getGroupMemberFunc := func(tx *sql.Tx) error {
		// check if user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, groupID, user.UserID)
		if err != nil {
			return err
		}

		groupMembers, err = c.repos.GroupMember.Find(tx, planetscale.GroupMemberFilter{
			GroupID: groupID,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), getGroupMemberFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(findGroupMembersResponse{
			GroupMembers: groupMembers,
			N:            len(groupMembers),
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

type findGroupMembersResponse struct {
	GroupMembers []*planetscale.GroupMember `json:"group_members"`
	N            int                        `json:"n"`
}

func (c *groupMemberController) HandlePostGroupMember(w http.ResponseWriter, r *http.Request) {
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

	var groupMember planetscale.GroupMember
	groupMember.GroupID = int64(group32)

	if err := json.NewDecoder(r.Body).Decode(&groupMember); err != nil {
		Error(w, r, err)
		return
	}

	createGroupMemberFunc := func(tx *sql.Tx) error {
		// check if user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, groupMember.GroupID, user.UserID)
		if err != nil {
			return err
		}

		err = c.repos.GroupMember.Create(tx, &groupMember)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), createGroupMemberFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groupMember); err != nil {
		Error(w, r, err)
		return
	}
}

func (c *groupMemberController) HandleDeleteGroupMember(w http.ResponseWriter, r *http.Request) {
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
	userID := chi.URLParam(r, "userID")

	deleteGroupMemberFunc := func(tx *sql.Tx) error {
		// check if user is a member of the group
		_, err = c.repos.GroupMember.Get(tx, groupID, user.UserID)
		if err != nil {
			return err
		}

		err = c.repos.GroupMember.Delete(tx, groupID, userID)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.tm.ExecuteInTx(r.Context(), deleteGroupMemberFunc)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
