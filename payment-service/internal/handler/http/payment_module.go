package http

import (
	"github.com/QosmuratSamat0/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type PaymentModule struct {
	paymentHandler *PaymentHandler
}

func NewPaymentModule(paymentService *service.PaymentService) *PaymentModule {
	return &PaymentModule{
		paymentHandler: NewPaymentHandler(paymentService),
	}
}

func (m *PaymentModule) RegisterRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Post("/init", m.paymentHandler.InitPayment)
		r.Post("/webhooks/{provider}", m.paymentHandler.ProcessWebhook)
		r.Get("/users/{userID}", m.paymentHandler.GetUserPayments)
	})
}
