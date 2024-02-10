package http

import (
	"database/sql"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	planetscale "github.com/harshav17/planet_scale"
	"github.com/patrickmn/go-cache"
)

type Middleware struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
	c     *cache.Cache
	clerk *clerk.Client
}

func NewMiddleware(repoProvider *planetscale.RepoProvider, tm planetscale.TransactionManager, c *cache.Cache, clerk *clerk.Client) *Middleware {
	return &Middleware{
		repos: repoProvider,
		tm:    tm,
		c:     c,
		clerk: clerk,
	}
}

func (m Middleware) withClerkSession() func(next http.Handler) http.Handler {
	return clerk.WithSessionV2(*m.clerk)
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func (m Middleware) EnsureValidToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sessClaims, ok := clerk.SessionFromContext(ctx)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"Unauthorized"}`))
				return
			}

			getUserFunc := func(tx *sql.Tx) error {
				user, err := m.repos.User.Get(tx, sessClaims.Claims.Subject)
				if err != nil {
					return err
				}

				// set the user in the context
				r = r.WithContext(planetscale.NewContextWithUser(ctx, user))

				return nil
			}

			err := m.tm.ExecuteInTx(r.Context(), getUserFunc)
			if err != nil {
				Error(w, r, err)
				return
			}

			// signed out
			next.ServeHTTP(w, r)
		})
	}
}
