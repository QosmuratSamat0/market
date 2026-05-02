package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	domain "github.com/QosmuratSamat/order-service/internal/domain/order"
	"github.com/QosmuratSamat/order-service/internal/handler/http/middleware"
	resp "github.com/QosmuratSamat/order-service/internal/lib/api/response"
	"github.com/QosmuratSamat/order-service/internal/lib/errs"
	orderUseCase "github.com/QosmuratSamat/order-service/internal/usecase/order"
)

type OrderHandler struct {
	orderUC *orderUseCase.UseCase
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateOrderStatusRequest struct {
	Status domain.Status `json:"status"`
}

type OrderResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Status    domain.Status       `json:"status"`
	Total     float64             `json:"total"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	SellerID  string  `json:"seller_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func NewOrderHandler(orderUC *orderUseCase.UseCase) *OrderHandler {
	return &OrderHandler{orderUC: orderUC}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.CreateOrder"
	log := slog.With("operation", op)

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	input := orderUseCase.CreateOrderInput{
		Items: make([]orderUseCase.CreateOrderItemInput, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		input.Items = append(input.Items, orderUseCase.CreateOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	o, err := h.orderUC.CreateOrder(r.Context(), middleware.GetUserID(r.Context()), input)
	if err != nil {
		log.Error("failed to create order", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, toOrderResponse(o))
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.GetOrder"
	log := slog.With("operation", op)

	o, err := h.orderUC.GetOrder(
		r.Context(),
		chi.URLParam(r, "id"),
		middleware.GetUserID(r.Context()),
		middleware.GetRole(r.Context()),
	)
	if err != nil {
		log.Error("failed to get order", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toOrderResponse(o))
}

func (h *OrderHandler) GetOrderInternal(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.GetOrderInternal"
	log := slog.With("operation", op)

	// Bypass auth for internal requests by providing admin role
	o, err := h.orderUC.GetOrder(
		r.Context(),
		chi.URLParam(r, "id"),
		"",      // no specific user
		"admin", // admin role allows access to any order
	)
	if err != nil {
		log.Error("failed to get order internal", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toOrderResponse(o))
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.GetMyOrders"
	log := slog.With("operation", op)

	orders, err := h.orderUC.GetMyOrders(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		log.Error("failed to get my orders", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toOrderListResponse(orders))
}

func (h *OrderHandler) GetSellerOrders(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.GetSellerOrders"
	log := slog.With("operation", op)

	orders, err := h.orderUC.GetSellerOrders(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		log.Error("failed to get seller orders", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toOrderListResponse(orders))
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.UpdateOrderStatus"
	log := slog.With("operation", op)

	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	o, err := h.orderUC.UpdateOrderStatus(
		r.Context(),
		chi.URLParam(r, "id"),
		req.Status,
		middleware.GetUserID(r.Context()),
		middleware.GetRole(r.Context()),
	)
	if err != nil {
		log.Error("failed to update order status", "error", err)
		handleOrderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toOrderResponse(o))
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	const op = "OrderHandler.DeleteOrder"
	log := slog.With("operation", op)

	err := h.orderUC.DeleteOrder(
		r.Context(),
		chi.URLParam(r, "id"),
		middleware.GetUserID(r.Context()),
		middleware.GetRole(r.Context()),
	)
	if err != nil {
		log.Error("failed to delete order", "error", err)
		handleOrderError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleOrderError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.OrderNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp.Error("order not found"))
	case errors.Is(err, errs.ProductNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp.Error("product not found"))
	case errors.Is(err, errs.ErrInvalidInput):
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid input"))
	case errors.Is(err, errs.OutOfStock):
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, resp.Error("product out of stock"))
	case errors.Is(err, errs.ErrUnauthorized):
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))
	case errors.Is(err, errs.ErrForbidden):
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, resp.Error("forbidden"))
	default:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
	}
}

func toOrderResponse(o *domain.Order) OrderResponse {
	return OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Status:    o.Status,
		Total:     o.Total,
		Items:     toOrderItemListResponse(o.Items),
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toOrderListResponse(orders []*domain.Order) []OrderResponse {
	result := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toOrderResponse(o))
	}
	return result
}

func toOrderItemListResponse(items []domain.OrderItem) []OrderItemResponse {
	result := make([]OrderItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			SellerID:  item.SellerID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	return result
}
