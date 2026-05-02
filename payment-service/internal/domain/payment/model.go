package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusSucceeded  PaymentStatus = "succeeded"
	StatusFailed     PaymentStatus = "failed"
	StatusCanceled   PaymentStatus = "canceled"
	StatusRefunded   PaymentStatus = "refunded"
)

type Payment struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	UserID         uuid.UUID
	Amount         int64
	Currency       string
	Provider       string
	ProviderID     string
	Status         PaymentStatus
	IdempotencyKey string
	Metadata       map[string]interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PaymentRequest struct {
	OrderID        uuid.UUID
	UserID         uuid.UUID
	Amount         int64
	Currency       string
	Provider       string
	IdempotencyKey string
}

type InitPaymentResult struct {
	PaymentID  uuid.UUID
	PaymentURL string
	Status     PaymentStatus
	Existing   bool
}
