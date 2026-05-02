package http

import (
	"github.com/go-chi/chi/v5"

	productMiddleware "github.com/QosmuratSamat0/product-service/internal/handler/http/middleware"
	productUseCase "github.com/QosmuratSamat0/product-service/internal/usecase/product"
)

type ProductModule struct {
	productHandler *ProductHandler
	jwtSecret      string
}

func NewProductModule(productUC *productUseCase.UseCase, jwtSecret string) *ProductModule {
	return &ProductModule{
		productHandler: NewProductHandler(productUC),
		jwtSecret:      jwtSecret,
	}
}

func (m *ProductModule) RegisterRoutes(r chi.Router) {
	authMiddleware := productMiddleware.Auth(m.jwtSecret)

	// Public routes — no auth required
	r.Get("/products", m.productHandler.GetAllProducts)
	r.Get("/products/", m.productHandler.GetAllProducts)
	r.Get("/products/{id}", m.productHandler.GetProduct)
	r.Get("/categories", m.productHandler.GetAllCategories)
	r.Get("/categories/", m.productHandler.GetAllCategories)
	r.Get("/categories/{categoryID}/products", m.productHandler.GetProductsByCategory)

	// Protected routes — auth required
	r.With(authMiddleware).Get("/products/my", m.productHandler.GetMyProducts)
	r.With(authMiddleware).Post("/products", m.productHandler.CreateProduct)
	r.With(authMiddleware).Post("/products/", m.productHandler.CreateProduct)
	r.With(authMiddleware).Put("/products/{id}", m.productHandler.UpdateProduct)
	r.With(authMiddleware).Delete("/products/{id}", m.productHandler.DeleteProduct)

	r.With(authMiddleware).Post("/categories", m.productHandler.CreateCategory)
	r.With(authMiddleware).Post("/categories/", m.productHandler.CreateCategory)
	r.With(authMiddleware).Delete("/categories/{id}", m.productHandler.DeleteCategory)
}
