package authmw

import (
	"net/http"
	"strings"

	"diplom/internal/apperror"
	"diplom/internal/auth"
	"diplom/pkg/httputil"
	"diplom/pkg/response"
)

// Bearer validates JWT and stores claims in request context.
func Bearer(jwt *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httputil.WriteError(w, http.StatusUnauthorized, response.CodeUnauthorized, "missing Authorization header")
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				httputil.WriteError(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid Authorization header")
				return
			}
			claims, err := jwt.Parse(parts[1])
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid or expired token")
				return
			}
			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

// Optional helper for handlers.
func MustClaims(w http.ResponseWriter, r *http.Request) (auth.Claims, bool) {
	c, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		httputil.HandleError(w, r, apperror.ErrUnauthorized)
		return auth.Claims{}, false
	}
	return c, true
}
