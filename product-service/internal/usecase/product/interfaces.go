package product

import (
	"context"

	"github.com/QosmuratSamat0/product-service/internal/domain/product"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *product.Product) error
	GetProductByID(ctx context.Context, id string) (*product.Product, error)
	GetAllProducts(ctx context.Context) ([]*product.Product, error)
	GetProductsByCategory(ctx context.Context, categoryID string) ([]*product.Product, error)
	GetProductsBySeller(ctx context.Context, sellerID string) ([]*product.Product, error)
	UpdateProduct(ctx context.Context, p *product.Product) error
	DeleteProduct(ctx context.Context, id string) error

	CreateCategory(ctx context.Context, c *product.Category) error
	GetAllCategories(ctx context.Context) ([]*product.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*product.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}
