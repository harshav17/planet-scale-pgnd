package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

func TestHandleGetProduct(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server.Server)

	t.Run("successful get", func(t *testing.T) {
		server.repos.Product = &db_mock.ProductRepo{
			GetFn: func(tx *sql.Tx, storyID int64) (*planetscale.Product, error) {
				return &planetscale.Product{
					ID:    storyID,
					Name:  "test product",
					Price: 100,
				}, nil
			},
		}

		req, err := http.NewRequest("GET", "/products/1", nil)
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

		expected := planetscale.Product{
			ID:    1,
			Name:  "test product",
			Price: 100,
		}
		var got planetscale.Product
		err = json.Unmarshal(rr.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
