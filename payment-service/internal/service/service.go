package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/QosmuratSamat0/payment-service/internal/domain/payment"
	"github.com/QosmuratSamat0/payment-service/internal/provider"
	"github.com/google/uuid"
)

type PaymentService struct {
	repo        PaymentRepository
	providers   map[string]provider.PaymentProvider
	broker      MessageBroker
	orderClient OrderServiceClient
	userClient  UserServiceClient
}

func NewPaymentService(
	repo PaymentRepository,
	broker MessageBroker,
	providers map[string]provider.PaymentProvider,
	orderClient OrderServiceClient,
	userClient UserServiceClient,
) *PaymentService {
	return &PaymentService{
		repo:        repo,
		broker:      broker,
		providers:   providers,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (s *PaymentService) InitPayment(ctx context.Context, req domain.PaymentRequest) (*domain.InitPaymentResult, error) {
	if req.OrderID == uuid.Nil {
		return nil, errors.New("order_id is required")
	}
	if req.UserID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if req.Currency == "" {
		return nil, errors.New("currency is required")
	}
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}
	if req.IdempotencyKey == "" {
		return nil, errors.New("idempotency_key is required")
	}

	orderAmount, err := s.orderClient.GetOrderAmount(ctx, req.OrderID.String())
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки заказа: %w", err)
	}

	if orderAmount != req.Amount {
		return nil, fmt.Errorf("сумма платежа (%d) не совпадает с суммой заказа (%d)", req.Amount, orderAmount)
	}

	if _, err := s.userClient.GetUserEmail(ctx, req.UserID.String()); err != nil {
		return nil, fmt.Errorf("ошибка проверки пользователя: %w", err)
	}

	existing, err := s.repo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки idempotency key: %w", err)
	}
	if existing != nil {
		return &domain.InitPaymentResult{
			PaymentID: existing.ID,
			Status:    existing.Status,
			Existing:  true,
		}, nil
	}

	p, ok := s.providers[req.Provider]
	if !ok {
		return nil, fmt.Errorf("провайдер %s не поддерживается", req.Provider)
	}

	newPayment := domain.Payment{
		ID:             uuid.New(),
		OrderID:        req.OrderID,
		UserID:         req.UserID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Provider:       req.Provider,
		Status:         domain.StatusPending,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.repo.Create(ctx, newPayment); err != nil {
		return nil, err
	}

	res, err := p.CreateTransaction(ctx, newPayment)
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, newPayment.ID, domain.StatusFailed, "")
		return nil, err
	}

	if res.ExternalID != "" {
		if err := s.repo.UpdateStatus(ctx, newPayment.ID, domain.StatusPending, res.ExternalID); err != nil {
			return nil, err
		}
	}

	// Simulation: if provider is mock, automatically succeed after 5 seconds
	if req.Provider == "mock" {
		go func(paymentID uuid.UUID, orderID uuid.UUID) {
			time.Sleep(5 * time.Second)
			// Use Background context because the request context might be cancelled
			_ = s.repo.UpdateStatus(context.Background(), paymentID, domain.StatusSucceeded, "mock_ext_"+paymentID.String())
			_ = s.broker.PublishPaymentSuccess(context.Background(), orderID.String())
		}(newPayment.ID, newPayment.OrderID)
	}

	return &domain.InitPaymentResult{
		PaymentID:  newPayment.ID,
		PaymentURL: res.PaymentURL,
		Status:     domain.StatusPending,
		Existing:   false,
	}, nil
}

func (s *PaymentService) GetUserPayments(ctx context.Context, userID string, limit, offset int) ([]domain.Payment, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user_id")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

func (s *PaymentService) ProcessWebhook(ctx context.Context, providerName string, payload []byte) error {
	p, ok := s.providers[providerName]
	if !ok {
		return fmt.Errorf("провайдер %s не поддерживается", providerName)
	}

	result, err := p.ParseWebhook(payload)
	if err != nil {
		return err
	}

	internalID, err := uuid.Parse(result.InternalID)
	if err != nil {
		return fmt.Errorf("invalid internal payment id: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, internalID, result.Status, result.ExternalID); err != nil {
		return err
	}

	if result.Status == domain.StatusSucceeded {
		return s.broker.PublishPaymentSuccess(ctx, result.OrderID)
	}

	return nil
}
