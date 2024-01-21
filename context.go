package planetscale

import (
	"context"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// contextKey represents an internal key for adding context fields.
// This is considered best practice as it prevents other packages from
// interfering with our context keys.
type contextKey int

// List of context keys.
// These are used to store request-scoped information.
const (
	userContextKey = contextKey(iota + 1)
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// HasScope checks whether our claims have a specific scope.
func (c *CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

// Auth0IDFromContext is a helper function that returns the ID of the current
// logged in user.
func Auth0IDFromContext(ctx context.Context) (string, *CustomClaims) {
	claims := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	ret := claims.RegisteredClaims.Subject
	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return ret, nil
	}
	return ret, customClaims
}

// NewContextWithUser returns a new context with the given user.
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, exist := ctx.Value(userContextKey).(*User)
	return user, exist
}
