package httputil

import (
	"net/http"

	"diplom/internal/apperror"
	"diplom/internal/auth"
	"diplom/pkg/response"

	"github.com/google/uuid"
)

// UserIDFromContext returns authenticated user id from JWT middleware.
func UserIDFromContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, response.CodeUnauthorized, "authentication required")
		return uuid.Nil, false
	}
	return claims.UserID, true
}

// RequireSelf ensures path user id matches token (for legacy /users/{id} routes).
func RequireSelf(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	authUserID, ok := UserIDFromContext(w, r)
	if !ok {
		return uuid.Nil, false
	}
	pathID, ok := URLParamUUID(w, r, param)
	if !ok {
		return uuid.Nil, false
	}
	if authUserID != pathID {
		HandleError(w, r, apperror.ErrForbidden)
		return uuid.Nil, false
	}
	return authUserID, true
}
