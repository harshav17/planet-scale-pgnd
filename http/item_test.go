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

func TestHnadleItem_All(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server.Server)

	t.Run("POST /items", func(t *testing.T) {
		t.Run("successful create", func(t *testing.T) {
			userID := "test-user-id"
			userID2 := "test-user-id-2"
			item := &planetscale.Item{
				Name:      "test-item",
				Price:     10.0,
				Quantity:  1,
				ExpenseID: 1,
				Splits: []*planetscale.ItemSplitNU{
					{
						UserID: &userID,
						Amount: 10.0,
					},
					{
						UserID: &userID2,
						Amount: 0.0,
					},
				},
			}
			server.repos.Item = &db_mock.ItemRepo{
				CreateFn: func(tx *sql.Tx, item *planetscale.Item) error {
					item.ItemID = 1
					return nil
				},
			}
			server.repos.ItemSplitNu = &db_mock.ItemSplitNURepo{
				CreateFn: func(tx *sql.Tx, split *planetscale.ItemSplitNU) error {
					return nil
				},
			}

			// generate token
			token := server.buildJWTForTesting(t, "test-user-id")
			body, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/items", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.router.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusCreated {
				t.Fatalf("expected status code %d, got %d", http.StatusCreated, status)
			}
		})
	})
}
