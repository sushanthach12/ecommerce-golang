package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/sushanthach12/ecom-go/internal/httpx"
)

type contextKey string

const userIDKey contextKey = "userID"
const userEmailKey contextKey = "userEmail"

// RequireAuth validates the Bearer JWT on the request and attaches the
// authenticated user's ID and email to the request context.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			httpx.HandleError(w, &httpx.UnauthorizedError{
				Message: "missing authorization header",
			})
			return
		}

		tokenStr, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenStr == "" {
			httpx.HandleError(w, &httpx.UnauthorizedError{
				Message: "malformed authorization header",
			})
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			httpx.HandleError(w, &httpx.UnauthorizedError{
				Message: "invalid or expired token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, userEmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts the authenticated user's ID set by RequireAuth.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

// UserEmailFromContext extracts the authenticated user's email set by RequireAuth.
func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}
