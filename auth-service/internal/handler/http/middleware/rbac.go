package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
	resp "github.com/QosmuratSamat0/auth-service/internal/lib/api/response"
)

func RequirePermission(permission domain.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetRole(r.Context())
			if role == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("unauthorized"))
				return
			}

			if !domain.HasPermission(domain.Role(role), permission) {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
