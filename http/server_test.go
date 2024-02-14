package http

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
	"github.com/patrickmn/go-cache"
	svix "github.com/svix/svix-webhooks/go"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type TestServer struct {
	*Server
	repos    *planetscale.RepoProvider
	services *planetscale.ServiceProvider
	jwk      *jose.JSONWebKey
}

func MustOpenServer(tb testing.TB) TestServer {
	tb.Helper()

	os.Setenv("AUTH0_DOMAIN", "http://localhost:8080")
	os.Setenv("AUTH0_AUDIENCE", "http://localhost:8080")

	tm := db_mock.TransactionManager{}
	tm.ExecuteInTxFn = func(ctx context.Context, fn func(*sql.Tx) error) error {
		return fn(nil)
	}

	repos := planetscale.RepoProvider{}
	repos.User = db_mock.UserRepo{
		GetFn: func(tx *sql.Tx, userID string) (*planetscale.User, error) {
			return &planetscale.User{
				UserID: userID,
			}, nil
		},
		UpsertFn: func(tx *sql.Tx, user *planetscale.User) error {
			return nil
		},
	}

	services := planetscale.ServiceProvider{}

	controllers := planetscale.ControllerProvider{}
	controllers.Product = NewProductController(&repos, &tm)
	controllers.ExpenseGroup = NewExpenseGroupController(&repos, &services, &tm)
	controllers.GroupMember = NewGroupMemberController(&repos, &tm)
	controllers.Expense = NewExpenseController(&repos, &services, &tm)
	controllers.Settlement = NewSettlementController(&repos, &tm)
	controllers.SplitType = NewSplitTypeController(&repos, &tm)
	controllers.User = NewUserController(&repos, &tm, &svix.Webhook{})
	controllers.Item = NewItemController(&repos, &services, &tm)

	c := cache.New(5*time.Minute, 10*time.Minute)
	client, _ := clerk.NewClient("test", clerk.WithBaseURL("http://localhost:8080"))
	middleware := NewMiddleware(&repos, &tm, c, &client)

	server := NewServer(&controllers, middleware)

	// handle JWT cycles
	jwk := generateJWK(tb)
	server.router.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		wk := WellKnownEndpoints{JWKSURI: "http://localhost:8080/.well-known/jwks.json"}
		err := json.NewEncoder(w).Encode(wk)
		if err != nil {
			tb.Fatal(err)
		}
	})
	server.router.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{jwk.Public()},
		})
		if err != nil {
			tb.Fatal(err)
		}
	})

	// Begin running test server.
	if err := server.Open(); err != nil {
		tb.Fatal(err)
	}

	return TestServer{
		Server:   server,
		repos:    &repos,
		services: &services,
		jwk:      jwk,
	}
}

// MustCloseServer is a test helper function for shutting down the server.
// Fail on error.
func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}

func (s TestServer) buildJWTForTesting(t testing.TB, subject string) string {
	t.Helper()

	key := jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(s.jwk.Algorithm),
		Key:       s.jwk,
	}
	signer, err := jose.NewSigner(key, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		t.Fatalf("could not build signer: %s", err.Error())
	}

	claims := jwt.Claims{
		Issuer:   "https://clerk.com:8080/",
		Audience: []string{"http://localhost:8080"},
		Subject:  subject,
	}

	token, err := jwt.Signed(signer).
		Claims(claims).
		CompactSerialize()

	if err != nil {
		t.Fatalf("could not build token: %s", err.Error())
	}

	return token
}

func generateJWK(t testing.TB) *jose.JSONWebKey {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("failed to generate private key")
	}

	return &jose.JSONWebKey{
		Key:       privateKey,
		KeyID:     "kid",
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
}

type WellKnownEndpoints struct {
	JWKSURI string `json:"jwks_uri"`
}
