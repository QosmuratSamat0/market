package provider

import (
	"context"

	domain "github.com/QosmuratSamat0/payment-service/internal/domain/payment"
)

type CreateTransactionResponse struct {
	PaymentURL string
	ExternalID string
}

type WebhookResult struct {
	InternalID string
	ExternalID string
	OrderID    string
	Status     domain.PaymentStatus
}

type PaymentProvider interface {
	CreateTransaction(ctx context.Context, p domain.Payment) (CreateTransactionResponse, error)
	ParseWebhook(payload []byte) (WebhookResult, error)
}
