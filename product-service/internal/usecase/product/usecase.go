package product

import (
	"context"
	"fmt"

	"github.com/QosmuratSamat0/product-service/internal/domain/product"
	"github.com/QosmuratSamat0/product-service/internal/lib/errs"
)

type UseCase struct {
	productRepo ProductRepository
}

type CreateProductInput struct {
	Name        string
	Description string
	Price       float64
	CategoryID  string
	ImageURL    string
	Stock       int
}

type UpdateProductInput struct {
	Name        string
	Description string
	Price       float64
	CategoryID  string
	ImageURL    string
	Stock       int
}

func NewUseCase(productRepo ProductRepository) *UseCase {
	return &UseCase{productRepo: productRepo}
}

// ---------- Products ----------

func (uc *UseCase) CreateProduct(ctx context.Context, sellerID string, input CreateProductInput) (*product.Product, error) {
	if input.Name == "" || input.Price <= 0 {
		return nil, errs.ErrInvalidInput
	}

	p := &product.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		SellerID:    sellerID,
		ImageURL:    input.ImageURL,
		Stock:       input.Stock,
	}

	if err := uc.productRepo.CreateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return p, nil
}

func (uc *UseCase) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	return uc.productRepo.GetProductByID(ctx, id)
}

func (uc *UseCase) GetAllProducts(ctx context.Context) ([]*product.Product, error) {
	return uc.productRepo.GetAllProducts(ctx)
}

func (uc *UseCase) GetProductsByCategory(ctx context.Context, categoryID string) ([]*product.Product, error) {
	return uc.productRepo.GetProductsByCategory(ctx, categoryID)
}

func (uc *UseCase) GetMyProducts(ctx context.Context, sellerID string) ([]*product.Product, error) {
	return uc.productRepo.GetProductsBySeller(ctx, sellerID)
}

func (uc *UseCase) UpdateProduct(ctx context.Context, id string, sellerID string, input UpdateProductInput) (*product.Product, error) {
	p, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only the seller who owns the product can update it
	if p.SellerID != sellerID {
		return nil, errs.Forbidden
	}

	if input.Name != "" {
		p.Name = input.Name
	}
	if input.Description != "" {
		p.Description = input.Description
	}
	if input.Price > 0 {
		p.Price = input.Price
	}
	if input.CategoryID != "" {
		p.CategoryID = input.CategoryID
	}
	if input.ImageURL != "" {
		p.ImageURL = input.ImageURL
	}
	if input.Stock >= 0 {
		p.Stock = input.Stock
	}

	if err := uc.productRepo.UpdateProduct(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *UseCase) DeleteProduct(ctx context.Context, id string, sellerID string, role string) error {
	p, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	// Admins can delete any product; sellers can only delete their own
	if role != "admin" && p.SellerID != sellerID {
		return errs.Forbidden
	}

	return uc.productRepo.DeleteProduct(ctx, id)
}

// ---------- Categories ----------

func (uc *UseCase) CreateCategory(ctx context.Context, name string) (*product.Category, error) {
	if name == "" {
		return nil, errs.ErrInvalidInput
	}
	c := &product.Category{Name: name}
	if err := uc.productRepo.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *UseCase) GetAllCategories(ctx context.Context) ([]*product.Category, error) {
	return uc.productRepo.GetAllCategories(ctx)
}

func (uc *UseCase) GetCategoryByID(ctx context.Context, id string) (*product.Category, error) {
	return uc.productRepo.GetCategoryByID(ctx, id)
}

func (uc *UseCase) DeleteCategory(ctx context.Context, id string) error {
	return uc.productRepo.DeleteCategory(ctx, id)
}
