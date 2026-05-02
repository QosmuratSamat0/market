package provider

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/QosmuratSamat0/payment-service/internal/domain/payment"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) CreateTransaction(ctx context.Context, p domain.Payment) (CreateTransactionResponse, error) {
	// В моке мы просто генерируем фиктивную ссылку
	// В реальности здесь был бы запрос к API (Stripe, PayPal, Kaspi и т.д.)
	return CreateTransactionResponse{
		PaymentURL: fmt.Sprintf("https://mock-payment-gateway.com/pay/%s", p.ID),
		ExternalID: "mock_ext_" + p.ID.String(),
	}, nil
}

func (m *MockProvider) ParseWebhook(payload []byte) (WebhookResult, error) {
	var data struct {
		InternalID string               `json:"internal_id"`
		ExternalID string               `json:"external_id"`
		OrderID    string               `json:"order_id"`
		Status     domain.PaymentStatus `json:"status"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return WebhookResult{}, err
	}

	return WebhookResult{
		InternalID: data.InternalID,
		OrderID:    data.OrderID,
		Status:     data.Status,
		ExternalID: data.ExternalID,
	}, nil
}
