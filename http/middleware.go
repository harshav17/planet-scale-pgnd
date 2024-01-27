package http

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	planetscale "github.com/harshav17/planet_scale"
	"github.com/patrickmn/go-cache"
)

type Middleware struct {
	repos *planetscale.RepoProvider
	tm    planetscale.TransactionManager
	c     *cache.Cache
}

func NewMiddleware(repoProvider *planetscale.RepoProvider, tm planetscale.TransactionManager, c *cache.Cache) *Middleware {
	return &Middleware{
		repos: repoProvider,
		tm:    tm,
		c:     c,
	}
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func (m Middleware) EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse(os.Getenv("AUTH0_DOMAIN") + "/") // TODO make sure prod code works
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &planetscale.CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func (m Middleware) SetUserContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			auth0ID, customClaims := planetscale.Auth0IDFromContext(ctx)

			r = r.WithContext(planetscale.NewContextWithUser(ctx, &planetscale.User{
				UserID: auth0ID,
				Email:  customClaims.Email,
				Name:   customClaims.Name,
			}))

			// check if the user data is latest
			// TODO use redis instead eventually
			cacheUser, found := m.c.Get(auth0ID)
			if found {
				retClaims := cacheUser.(*planetscale.CustomClaims)
				if retClaims.Email == customClaims.Email && retClaims.Name == customClaims.Name {
					next.ServeHTTP(w, r)
					return
				}
			}

			// if not found OR if the data is not latest, then update the cache and DB
			// Update user to the latest information from Auth0.
			upsertUserFunc := func(tx *sql.Tx) error {
				err := m.repos.User.Upsert(tx, &planetscale.User{
					UserID: auth0ID,
					Email:  customClaims.Email,
					Name:   customClaims.Name,
				})
				if err != nil {
					return err
				}
				return nil
			}

			err := m.tm.ExecuteInTx(r.Context(), upsertUserFunc)
			if err != nil {
				Error(w, r, err)
				return
			}

			m.c.Set(auth0ID, customClaims, cache.DefaultExpiration)

			next.ServeHTTP(w, r)
		})
	}
}
