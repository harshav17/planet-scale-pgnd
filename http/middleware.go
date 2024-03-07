package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

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

func extractToken(r *http.Request) string {
	headerToken := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimPrefix(headerToken, "Bearer ")
}

func isJWT(token string) bool {
	return strings.Contains(token, ".")
}

func setUserContext(r *http.Request, userID string) {
	user := &planetscale.User{
		UserID: userID,
	}
	*r = *r.WithContext(planetscale.NewContextWithUser(r.Context(), user))
}

func fetchUserInfo(token string) (*UserInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://clerk.skwabbl.com/oauth/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func httpError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (m *Middleware) OpaqueTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if isJWT(token) {
			next.ServeHTTP(w, r)
			return
		}

		userInfo, err := fetchUserInfo(token)
		if err != nil {
			httpError(w, http.StatusUnauthorized, err.Error())
			return
		}

		setUserContext(r, userInfo.UserId)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if !isJWT(token) {
			next.ServeHTTP(w, r)
			return
		}

		clerkMiddleware := clerk.WithSessionV2(*m.clerk)
		clerkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !m.setUserFromClerkSession(w, r) {
				return // Unauthorized error is handled in setUserFromClerkSession
			}
			next.ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})
}

func (m *Middleware) setUserFromClerkSession(w http.ResponseWriter, r *http.Request) bool {
	ctx := r.Context()
	sessClaims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		httpError(w, http.StatusUnauthorized, `{"message":"Unauthorized"}`)
		return false
	}

	setUserContext(r, sessClaims.Subject)
	return true
}

type UserInfo struct {
	UserId string `json:"user_id"`
}
