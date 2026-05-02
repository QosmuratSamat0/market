package http

import (
	"github.com/go-chi/chi/v5"

	userMiddleware "github.com/QosmuratSamat0/user-service-market/internal/handler/http/middleware"
	userUseCase "github.com/QosmuratSamat0/user-service-market/internal/usecase/user"
)

type UserModule struct {
	userHandler *UserHandler
	jwtSecret   string
}

func NewUserModule(userUC *userUseCase.UseCase, jwtSecret string) *UserModule {
	return &UserModule{
		userHandler: NewUserHandler(userUC),
		jwtSecret:   jwtSecret,
	}
}

func (m *UserModule) RegisterRoutes(r chi.Router) {
	r.Route("/internal/users", func(r chi.Router) {
		r.Get("/by-email", m.userHandler.GetInternalUserByEmail)
		r.Get("/{id}", m.userHandler.GetInternalUserByID)
		r.Post("/", m.userHandler.CreateInternalUser)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(userMiddleware.Auth(m.jwtSecret))

		// Authenticated user routes
		r.Get("/me", m.userHandler.GetMe)
		r.Put("/me", m.userHandler.UpdateMe)

		// Admin routes
		r.Get("/", m.userHandler.GetAllUsers)
		r.Post("/", m.userHandler.CreateUser)
		r.Get("/{id}", m.userHandler.GetUser)
		r.Put("/{id}", m.userHandler.UpdateUser)
		r.Delete("/{id}", m.userHandler.DeleteUser)
	})

	r.Route("/chat", func(r chi.Router) {
		r.Use(userMiddleware.Auth(m.jwtSecret))
		r.Get("/history", m.userHandler.GetChatHistory)
		r.Post("/send", m.userHandler.SendMessage)
	})
}
