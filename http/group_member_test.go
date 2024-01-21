package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

func TestHandleGroupMembers_All(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server.Server)

	t.Run("GET /groups/1/members", func(t *testing.T) {
		t.Run("successful find", func(t *testing.T) {
			user_id := "test-user-id"
			server.repos.GroupMember = &db_mock.GroupMemberRepo{
				FindFn: func(tx *sql.Tx, filter planetscale.GroupMemberFilter) ([]*planetscale.GroupMember, error) {
					return []*planetscale.GroupMember{
						{
							GroupID: 1,
							UserID:  user_id,
						},
					}, nil
				},
			}

			// generate token
			token := server.buildJWTForTesting(t, user_id)
			req, err := http.NewRequest("GET", "/groups/1/members", nil)
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

			var got findGroupMembersResponse
			err = json.Unmarshal(rr.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.GroupMembers) != 1 && got.N != 1 {
				t.Fatalf("expected 1 group member, got %d", len(got.GroupMembers))
			}
			if got.GroupMembers[0].GroupID != 1 {
				t.Fatalf("expected group id 1, got %d", got.GroupMembers[0].GroupID)
			}
		})
	})

	t.Run("POST /groups/1/members", func(t *testing.T) {
		t.Run("successful post", func(t *testing.T) {
			server.repos.GroupMember = &db_mock.GroupMemberRepo{
				CreateFn: func(tx *sql.Tx, groupMember *planetscale.GroupMember) error {
					return nil
				},
			}

			groupMember := planetscale.GroupMember{
				UserID: "test-user-id",
			}
			body, err := json.Marshal(groupMember)
			if err != nil {
				t.Fatal(err)
			}

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("POST", "/groups/1/members", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
			}
		})
	})

	t.Run("DELETE /groups/1/members/user-id", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			server.repos.GroupMember = &db_mock.GroupMemberRepo{
				DeleteFn: func(tx *sql.Tx, groupID int64, userID string) error {
					return nil
				},
			}

			token := server.buildJWTForTesting(t, "test_user_id")
			req, err := http.NewRequest("DELETE", "/groups/1/members/test-user-id", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNoContent {
				t.Errorf("expected status code %d, got %d", http.StatusNoContent, status)
			}
		})
	})
}
