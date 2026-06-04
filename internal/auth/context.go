package auth

import (
	"context"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type contextKey struct{}

// Claims is stored in request context after JWT validation.
type Claims struct {
	UserID uuid.UUID
	Email  string
	Role   models.Role
}

func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(contextKey{}).(Claims)
	return c, ok
}
