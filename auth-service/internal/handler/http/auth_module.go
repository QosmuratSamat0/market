package http

import (
	"github.com/go-chi/chi/v5"
	authUseCase "github.com/QosmuratSamat0/auth-service/internal/usecase/auth"
)

type AuthModule struct {
	authHandler *AuthHandler
}

func NewAuthModule(authUC *authUseCase.UseCase) *AuthModule {
	return &AuthModule{
		authHandler: NewAuthHandler(authUC),
	}
}

func (m *AuthModule) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", m.authHandler.Register)
		r.Post("/login", m.authHandler.Login)
		r.Post("/refresh", m.authHandler.Refresh)
		r.Post("/logout", m.authHandler.Logout)
	})
}
