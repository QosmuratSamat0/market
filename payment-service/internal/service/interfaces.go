package service

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/QosmuratSamat0/payment-service/internal/domain/payment"
)

type PaymentRepository interface {
	Create(ctx context.Context, p domain.Payment) error
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, externalID string) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]domain.Payment, error)
}

type MessageBroker interface {
	PublishPaymentSuccess(ctx context.Context, orderID string) error
}

type OrderServiceClient interface {
	GetOrderAmount(ctx context.Context, orderID string) (int64, error)
}

type UserServiceClient interface {
	GetUserEmail(ctx context.Context, userID string) (string, error)
}
