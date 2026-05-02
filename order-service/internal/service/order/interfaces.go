package order

import (
	"context"

	clientProduct "github.com/QosmuratSamat/order-service/internal/client/product"
)

type ProductClient interface {
	GetProductByID(ctx context.Context, id string) (*clientProduct.Product, error)
}
