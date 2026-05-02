package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/QosmuratSamat0/product-service/internal/domain/product"
	"github.com/QosmuratSamat0/product-service/internal/handler/http/middleware"
	resp "github.com/QosmuratSamat0/product-service/internal/lib/api/response"
	"github.com/QosmuratSamat0/product-service/internal/lib/errs"
	productUseCase "github.com/QosmuratSamat0/product-service/internal/usecase/product"
)

type ProductHandler struct {
	productUC *productUseCase.UseCase
}

// ---------- Request / Response ----------

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	SellerID    string  `json:"seller_id"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func NewProductHandler(productUC *productUseCase.UseCase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

// ---------- Product Handlers ----------

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.GetAllProducts"
	log := slog.With("operation", op)

	products, err := h.productUC.GetAllProducts(r.Context())
	if err != nil {
		log.Error("failed to get products", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toProductListResponse(products))
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.GetProduct"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	p, err := h.productUC.GetProductByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get product", "error", err)
		handleProductError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toProductResponse(p))
}

func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.GetProductsByCategory"
	log := slog.With("operation", op)

	categoryID := chi.URLParam(r, "categoryID")
	products, err := h.productUC.GetProductsByCategory(r.Context(), categoryID)
	if err != nil {
		log.Error("failed to get products by category", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toProductListResponse(products))
}

func (h *ProductHandler) GetMyProducts(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.GetMyProducts"
	log := slog.With("operation", op)

	sellerID := middleware.GetUserID(r.Context())
	products, err := h.productUC.GetMyProducts(r.Context(), sellerID)
	if err != nil {
		log.Error("failed to get my products", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toProductListResponse(products))
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.CreateProduct"
	log := slog.With("operation", op)

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	sellerID := middleware.GetUserID(r.Context())

	p, err := h.productUC.CreateProduct(r.Context(), sellerID, productUseCase.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		Stock:       req.Stock,
	})
	if err != nil {
		log.Error("failed to create product", "error", err)
		handleProductError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, toProductResponse(p))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.UpdateProduct"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	sellerID := middleware.GetUserID(r.Context())

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	p, err := h.productUC.UpdateProduct(r.Context(), id, sellerID, productUseCase.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		Stock:       req.Stock,
	})
	if err != nil {
		log.Error("failed to update product", "error", err)
		handleProductError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toProductResponse(p))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.DeleteProduct"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	sellerID := middleware.GetUserID(r.Context())
	role := middleware.GetRole(r.Context())

	if err := h.productUC.DeleteProduct(r.Context(), id, sellerID, role); err != nil {
		log.Error("failed to delete product", "error", err)
		handleProductError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---------- Category Handlers ----------

func (h *ProductHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productUC.GetAllCategories(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toCategoryListResponse(categories))
}

func (h *ProductHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	c, err := h.productUC.CreateCategory(r.Context(), req.Name)
	if err != nil {
		handleProductError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, toCategoryResponse(c))
}

func (h *ProductHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.productUC.DeleteCategory(r.Context(), id); err != nil {
		handleProductError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Helpers ----------

func handleProductError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.ProductNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp.Error("product not found"))
	case errors.Is(err, errs.CategoryNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp.Error("category not found"))
	case errors.Is(err, errs.ErrInvalidInput):
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid input"))
	case errors.Is(err, errs.Forbidden):
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, resp.Error("forbidden"))
	default:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
	}
}

func toProductResponse(p *product.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
		SellerID:    p.SellerID,
		ImageURL:    p.ImageURL,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toProductListResponse(products []*product.Product) []ProductResponse {
	result := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		result = append(result, toProductResponse(p))
	}
	return result
}

func toCategoryResponse(c *product.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toCategoryListResponse(categories []*product.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(c))
	}
	return result
}
