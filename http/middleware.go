package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

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
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerToken := strings.TrimSpace(r.Header.Get("Authorization"))
			token := strings.TrimPrefix(headerToken, "Bearer ")
			if strings.Contains(token, ".") {
				// For non-JWTs, use Clerk's middleware to handle session validation
				clerkMiddleware := clerk.WithSessionV2(*m.clerk)
				clerkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					sessClaims, ok := clerk.SessionFromContext(ctx)
					if !ok {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(`{"message":"Unauthorized"}`))
						return
					}

					r = r.WithContext(planetscale.NewContextWithUser(ctx, &planetscale.User{
						UserID: sessClaims.Subject,
					}))

					next.ServeHTTP(w, r)
				})).ServeHTTP(w, r)
				return
			} else {
				// Opaque token, call https://clerk.skwabbl.com/oauth/userinfo to get the user info
				// and set the user in the context
				client := &http.Client{}
				req, err := http.NewRequest("GET", "https://clerk.skwabbl.com/oauth/userinfo", nil)
				if err != nil {
					Error(w, r, err)
					return
				}

				req.Header.Set("Authorization", "Bearer "+token)
				resp, err := client.Do(req)
				if err != nil {
					Error(w, r, err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					Error(w, r, err)
					return
				}

				var userInfo UserInfo
				if err := json.Unmarshal(body, &userInfo); err != nil {
					Error(w, r, err)
					return
				}

				// set the user in the context
				user := &planetscale.User{
					UserID: userInfo.UserId,
				}
				r = r.WithContext(planetscale.NewContextWithUser(r.Context(), user))
			}

			// Proceed with next middleware or handler if JWT handling does not return
			next.ServeHTTP(w, r)
		})
	}
}

type UserInfo struct {
	UserId string `json:"user_id"`
}
