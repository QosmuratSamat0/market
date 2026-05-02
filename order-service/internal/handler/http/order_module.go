package http

import (
	"github.com/go-chi/chi/v5"

	orderMiddleware "github.com/QosmuratSamat/order-service/internal/handler/http/middleware"
	orderUseCase "github.com/QosmuratSamat/order-service/internal/usecase/order"
)

type OrderModule struct {
	orderHandler *OrderHandler
	jwtSecret    string
}

func NewOrderModule(orderUC *orderUseCase.UseCase, jwtSecret string) *OrderModule {
	return &OrderModule{
		orderHandler: NewOrderHandler(orderUC),
		jwtSecret:    jwtSecret,
	}
}

func (m *OrderModule) RegisterRoutes(r chi.Router) {
	authMiddleware := orderMiddleware.Auth(m.jwtSecret)

	r.Route("/orders", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", m.orderHandler.CreateOrder)
		r.Get("/my", m.orderHandler.GetMyOrders)
		r.Get("/seller", m.orderHandler.GetSellerOrders)
		r.Get("/{id}", m.orderHandler.GetOrder)
		r.Patch("/{id}/status", m.orderHandler.UpdateOrderStatus)
		r.Delete("/{id}", m.orderHandler.DeleteOrder)
	})

	// Internal route for service-to-service communication (no auth)
	r.Get("/internal/orders/{id}", m.orderHandler.GetOrderInternal)
}
