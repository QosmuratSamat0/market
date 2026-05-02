package order

import (
	"context"

	domain "github.com/QosmuratSamat/order-service/internal/domain/order"
	orderService "github.com/QosmuratSamat/order-service/internal/service/order"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o *domain.Order) error
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	GetOrdersByUser(ctx context.Context, userID string) ([]*domain.Order, error)
	GetOrdersBySeller(ctx context.Context, sellerID string) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status domain.Status) error
	DeleteOrder(ctx context.Context, id string) error
}

type OrderService interface {
	PrepareOrder(ctx context.Context, inputItems []orderService.CreateItemInput) (*orderService.PreparedOrder, error)
}
