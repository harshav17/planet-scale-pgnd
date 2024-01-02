package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

func TestHandleExpenseGroups_All(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server.Server)

	t.Run("GET /groups", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			groups := []*planetscale.ExpenseGroup{
				{
					GroupName: "test group",
					CreateBy:  "test-user-id",
				},
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				ListAllForUserFn: func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
					return groups, nil
				},
			}

			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			expected := findExpenseGroupsResponse{
				ExpenseGroups: groups,
				N:             1,
			}
			var got findExpenseGroupsResponse
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected %v, got %v", expected, got)
			}

		})

		t.Run("no items found", func(t *testing.T) {
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				ListAllForUserFn: func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
					return nil, nil
				},
			}

			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			expected := findExpenseGroupsResponse{
				ExpenseGroups: nil,
				N:             0,
			}
			var got findExpenseGroupsResponse
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})

		t.Run("error", func(t *testing.T) {
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				ListAllForUserFn: func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
					return nil, &planetscale.Error{
						Code:    planetscale.ENOTFOUND,
						Message: "not found",
					}
				},
			}

			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
			} else if rr.Body.String() != "{\"error\":\"not found\"}\n" {
				t.Errorf("expected error message, got %s", rr.Body.String())
			}
		})
	})

	t.Run("POST /groups", func(t *testing.T) {
		t.Run("successful post", func(t *testing.T) {
			expenseGroup := planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  "test-user-id",
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				CreateFn: func(tx *sql.Tx, eg *planetscale.ExpenseGroup) error {
					eg.ExpenseGroupID = 1
					return nil
				},
			}
			server.repos.GroupMember = &db_mock.GroupMemberRepo{
				CreateFn: func(tx *sql.Tx, gm *planetscale.GroupMember) error {
					return nil
				},
			}
			// TODO how do you test a real service / complex operations involved?

			body, err := json.Marshal(expenseGroup)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			var got planetscale.ExpenseGroup
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseGroupID != 1 {
				t.Errorf("expected group id 1, got %d", got.ExpenseGroupID)
			}
		})

		t.Run("error", func(t *testing.T) {
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				CreateFn: func(tx *sql.Tx, eg *planetscale.ExpenseGroup) error {
					return &planetscale.Error{
						Code:    planetscale.ECONFLICT,
						Message: "Another group with the same name already exists",
					}
				},
			}

			expenseGroup := planetscale.ExpenseGroup{
				GroupName: "test group",
				CreateBy:  "test-user-id",
			}
			body, err := json.Marshal(expenseGroup)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusConflict {
				t.Errorf("expected status code %d, got %d", http.StatusConflict, status)
			}
		})
	})

	t.Run("PATCH /groups/:id", func(t *testing.T) {
		t.Run("successful patch", func(t *testing.T) {
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       "test-user-id",
				ExpenseGroupID: 1,
			}
			expenseGroupUpdate := planetscale.ExpenseGroupUpdate{
				GroupName: "test group",
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				UpdateFn: func(tx *sql.Tx, groupID int64, group *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
			}

			body, err := json.Marshal(expenseGroupUpdate)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("PATCH", "/groups/1", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			var got planetscale.ExpenseGroup
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseGroupID != 1 {
				t.Errorf("expected group id 1, got %d", got.ExpenseGroupID)
			}
		})
	})

	t.Run("DELETE /groups/:id", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				DeleteFn: func(tx *sql.Tx, groupID int64) error {
					return nil
				},
			}

			req, err := http.NewRequest("DELETE", "/groups/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNoContent {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}
		})
	})

	t.Run("GET /groups/:id", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       "test-user-id",
				ExpenseGroupID: 1,
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				GetFn: func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
			}

			req, err := http.NewRequest("GET", "/groups/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			var got planetscale.ExpenseGroup
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseGroupID != 1 {
				t.Errorf("expected group id 1, got %d", got.ExpenseGroupID)
			}
		})
	})
}
