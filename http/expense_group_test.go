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
	service_mock "github.com/harshav17/planet_scale/mock/service"
)

func TestHandleExpenseGroups_All(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server.Server)

	t.Run("GET /groups", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			user_id := "test-user-id"
			groups := []*planetscale.ExpenseGroup{
				{
					GroupName: "test group",
					CreateBy:  user_id,
				},
				{
					GroupName: "test group 2",
					CreateBy:  "test-user-id-2",
				},
			}

			filteredGroups := []*planetscale.ExpenseGroup{}
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				ListAllForUserFn: func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
					// filter groups by user id
					for _, group := range groups {
						if group.CreateBy == userID {
							filteredGroups = append(filteredGroups, group)
						}
					}
					return filteredGroups, nil
				},
			}

			token := server.buildJWTForTesting(t, user_id)
			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Fatalf("expected status code %d, got %d", http.StatusOK, status)
			}

			expected := findExpenseGroupsResponse{
				ExpenseGroups: filteredGroups,
				N:             1,
			}
			var got findExpenseGroupsResponse
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, got) {
				t.Fatalf("expected %v, got %v", expected, got)
			}

		})

		t.Run("no items found", func(t *testing.T) {
			groups := []*planetscale.ExpenseGroup{
				{
					GroupName: "test group 2",
					CreateBy:  "test-user-id-2",
				},
			}

			filteredGroups := []*planetscale.ExpenseGroup{}
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				ListAllForUserFn: func(tx *sql.Tx, userID string) ([]*planetscale.ExpenseGroup, error) {
					// filter groups by user id
					for _, group := range groups {
						if group.CreateBy == userID {
							filteredGroups = append(filteredGroups, group)
						}
					}
					return filteredGroups, nil
				},
			}

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			expected := findExpenseGroupsResponse{
				ExpenseGroups: []*planetscale.ExpenseGroup{},
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

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("GET", "/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

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

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

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
				t.Fatalf("expected group id 1, got %d", got.ExpenseGroupID)
			} else if got.GroupName != "test group" {
				t.Fatalf("expected group name test group, got %s", got.GroupName)
			} else if got.CreateBy != "test_user_id" {
				t.Fatalf("expected create by test_user_id, got %s", got.CreateBy)
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

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

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
			userID := "test-user-id"
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       userID,
				ExpenseGroupID: 1,
			}
			updateGroupName := "test group updated"
			expenseGroupUpdate := planetscale.ExpenseGroupUpdate{
				GroupName: updateGroupName,
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				GetFn: func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
				UpdateFn: func(tx *sql.Tx, groupID int64, group *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error) {
					expenseGroup.GroupName = updateGroupName
					return &expenseGroup, nil
				},
			}

			body, err := json.Marshal(expenseGroupUpdate)
			if err != nil {
				t.Fatal(err)
			}

			token := server.buildJWTForTesting(t, userID)
			req, err := http.NewRequest("PATCH", "/groups/1", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Fatalf("expected status code %d, got %d", http.StatusOK, status)
			}

			var got planetscale.ExpenseGroup
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if got.ExpenseGroupID != 1 {
				t.Fatalf("expected group id 1, got %d", got.ExpenseGroupID)
			} else if got.GroupName != updateGroupName {
				t.Fatalf("expected group name %s, got %s", updateGroupName, got.GroupName)
			}
		})

		t.Run("unauthorized", func(t *testing.T) {
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       "test-user-id-1",
				ExpenseGroupID: 1,
			}
			expenseGroupUpdate := planetscale.ExpenseGroupUpdate{
				GroupName: "test group updated",
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				GetFn: func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
				UpdateFn: func(tx *sql.Tx, groupID int64, group *planetscale.ExpenseGroupUpdate) (*planetscale.ExpenseGroup, error) {
					return nil, nil
				},
			}

			body, err := json.Marshal(expenseGroupUpdate)
			if err != nil {
				t.Fatal(err)
			}

			token := server.buildJWTForTesting(t, "test-user-id-2")
			req, err := http.NewRequest("PATCH", "/groups/1", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, status)
			}
		})
	})

	t.Run("DELETE /groups/:id", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			userID := "test-user-id"
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       userID,
				ExpenseGroupID: 1,
			}
			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				GetFn: func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
				DeleteFn: func(tx *sql.Tx, groupID int64) error {
					return nil
				},
			}

			token := server.buildJWTForTesting(t, userID)
			req, err := http.NewRequest("DELETE", "/groups/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNoContent {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}
		})

		t.Run("unauthorized", func(t *testing.T) {
			expenseGroup := planetscale.ExpenseGroup{
				GroupName:      "test group",
				CreateBy:       "test-user-id-1",
				ExpenseGroupID: 1,
			}

			server.repos.ExpenseGroup = &db_mock.ExpenseGroupRepo{
				GetFn: func(tx *sql.Tx, groupID int64) (*planetscale.ExpenseGroup, error) {
					return &expenseGroup, nil
				},
				DeleteFn: func(tx *sql.Tx, groupID int64) error {
					return nil
				},
			}

			token := server.buildJWTForTesting(t, "test-user-id-2")
			req, err := http.NewRequest("DELETE", "/groups/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, status)
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

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("GET", "/groups/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

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

	t.Run("GET /groups/:id/balances", func(t *testing.T) {
		t.Run("successful get", func(t *testing.T) {
			balances := []*planetscale.Balance{
				{
					UserID: "test-user-id",
					Amount: 100,
				},
			}

			server.services.Balance = &service_mock.BalanceService{
				GetGroupBalancesFn: func(groupID int64) ([]*planetscale.Balance, error) {
					return balances, nil
				},
			}

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("GET", "/groups/1/balances", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
			}

			var got []planetscale.Balance
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Errorf("expected 1 balance, got %d", len(got))
			} else if got[0].UserID != "test-user-id" {
				t.Errorf("expected user id test-user-id, got %s", got[0].UserID)
			} else if got[0].Amount != 100 {
				t.Errorf("expected amount 100, got %f", got[0].Amount)
			}
		})
	})
}
