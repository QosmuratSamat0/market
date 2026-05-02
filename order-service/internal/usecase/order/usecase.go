package order

import (
	"context"

	domain "github.com/QosmuratSamat/order-service/internal/domain/order"
	"github.com/QosmuratSamat/order-service/internal/lib/errs"
	orderService "github.com/QosmuratSamat/order-service/internal/service/order"
)

type UseCase struct {
	orderRepo    OrderRepository
	orderService OrderService
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
}

func NewUseCase(orderRepo OrderRepository, orderService OrderService) *UseCase {
	return &UseCase{
		orderRepo:    orderRepo,
		orderService: orderService,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, userID string, input CreateOrderInput) (*domain.Order, error) {
	if userID == "" {
		return nil, errs.ErrUnauthorized
	}

	serviceItems := make([]orderService.CreateItemInput, 0, len(input.Items))
	for _, item := range input.Items {
		serviceItems = append(serviceItems, orderService.CreateItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	prepared, err := uc.orderService.PrepareOrder(ctx, serviceItems)
	if err != nil {
		return nil, err
	}

	o := &domain.Order{
		UserID: userID,
		Status: domain.StatusPending,
		Total:  prepared.Total,
		Items:  prepared.Items,
	}
	if err := uc.orderRepo.CreateOrder(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (uc *UseCase) GetOrder(ctx context.Context, id string, userID string, role string) (*domain.Order, error) {
	o, err := uc.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canAccessOrder(o, userID, role) {
		return nil, errs.ErrForbidden
	}
	return o, nil
}

func (uc *UseCase) GetMyOrders(ctx context.Context, userID string) ([]*domain.Order, error) {
	if userID == "" {
		return nil, errs.ErrUnauthorized
	}
	return uc.orderRepo.GetOrdersByUser(ctx, userID)
}

func (uc *UseCase) GetSellerOrders(ctx context.Context, sellerID string) ([]*domain.Order, error) {
	if sellerID == "" {
		return nil, errs.ErrUnauthorized
	}
	return uc.orderRepo.GetOrdersBySeller(ctx, sellerID)
}

func (uc *UseCase) UpdateOrderStatus(ctx context.Context, id string, status domain.Status, userID string, role string) (*domain.Order, error) {
	if !isValidStatus(status) {
		return nil, errs.ErrInvalidInput
	}

	o, err := uc.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == "admin" || isSellerOrder(o, userID) || (o.UserID == userID && status == domain.StatusCancelled) {
		if err := uc.orderRepo.UpdateOrderStatus(ctx, id, status); err != nil {
			return nil, err
		}
		o.Status = status
		return o, nil
	}

	return nil, errs.ErrForbidden
}

func (uc *UseCase) DeleteOrder(ctx context.Context, id string, userID string, role string) error {
	o, err := uc.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return err
	}
	if role != "admin" && o.UserID != userID {
		return errs.ErrForbidden
	}
	return uc.orderRepo.DeleteOrder(ctx, id)
}

func canAccessOrder(o *domain.Order, userID string, role string) bool {
	return role == "admin" || o.UserID == userID || isSellerOrder(o, userID)
}

func isSellerOrder(o *domain.Order, sellerID string) bool {
	for _, item := range o.Items {
		if item.SellerID == sellerID {
			return true
		}
	}
	return false
}

func isValidStatus(status domain.Status) bool {
	switch status {
	case domain.StatusPending, domain.StatusPaid, domain.StatusShipped, domain.StatusDelivered, domain.StatusCancelled:
		return true
	default:
		return false
	}
}
