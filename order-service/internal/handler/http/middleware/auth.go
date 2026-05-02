package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	resp "github.com/QosmuratSamat/order-service/internal/lib/api/response"
	"github.com/QosmuratSamat/order-service/internal/lib/tokens"
	"github.com/go-chi/chi/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func Auth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("missing authorization header"))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid authorization header format"))
				return
			}

			claims, err := tokens.ParseJWT(parts[1], secret)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid or expired token"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(RoleKey).(string); ok {
		return v
	}
	return ""
}

func RequireRole(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetRole(r.Context())
			if role != requiredRole {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("access denied: insufficient permissions"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserOnly(paramName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenUserID := GetUserID(r.Context())
			role := GetRole(r.Context())

			if role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			urlUserID := chi.URLParam(r, paramName)
			if tokenUserID != urlUserID {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("access denied: you can only access your own data"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
