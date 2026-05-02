package order

import (
	"context"

	domain "github.com/QosmuratSamat/order-service/internal/domain/order"
	"github.com/QosmuratSamat/order-service/internal/lib/errs"
)

type Service struct {
	productClient ProductClient
}

type CreateItemInput struct {
	ProductID string
	Quantity  int
}

type PreparedOrder struct {
	Items []domain.OrderItem
	Total float64
}

func NewService(productClient ProductClient) *Service {
	return &Service{productClient: productClient}
}

func (s *Service) PrepareOrder(ctx context.Context, inputItems []CreateItemInput) (*PreparedOrder, error) {
	if len(inputItems) == 0 {
		return nil, errs.ErrInvalidInput
	}

	items := make([]domain.OrderItem, 0, len(inputItems))
	var total float64

	for _, input := range inputItems {
		if input.ProductID == "" || input.Quantity <= 0 {
			return nil, errs.ErrInvalidInput
		}

		product, err := s.productClient.GetProductByID(ctx, input.ProductID)
		if err != nil {
			return nil, err
		}
		if product.Stock < input.Quantity {
			return nil, errs.OutOfStock
		}

		items = append(items, domain.OrderItem{
			ProductID: product.ID,
			SellerID:  product.SellerID,
			Quantity:  input.Quantity,
			Price:     product.Price,
		})
		total += product.Price * float64(input.Quantity)
	}

	return &PreparedOrder{Items: items, Total: total}, nil
}
