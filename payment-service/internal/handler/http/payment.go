package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/QosmuratSamat0/payment-service/internal/domain/payment"
	"github.com/QosmuratSamat0/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

type InitPaymentRequest struct {
	OrderID        string `json:"order_id"`
	UserID         string `json:"user_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	Provider       string `json:"provider"`
	IdempotencyKey string `json:"idempotency_key"`
}

type InitPaymentResponse struct {
	PaymentID  string                `json:"payment_id"`
	PaymentURL string                `json:"payment_url,omitempty"`
	Status     payment.PaymentStatus `json:"status"`
	Existing   bool                  `json:"existing"`
}

type PaymentResponse struct {
	ID             string                `json:"id"`
	OrderID        string                `json:"order_id"`
	UserID         string                `json:"user_id"`
	Amount         int64                 `json:"amount"`
	Currency       string                `json:"currency"`
	Provider       string                `json:"provider"`
	ProviderID     string                `json:"provider_id,omitempty"`
	Status         payment.PaymentStatus `json:"status"`
	IdempotencyKey string                `json:"idempotency_key,omitempty"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at,omitempty"`
}

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) InitPayment(w http.ResponseWriter, r *http.Request) {
	var req InitPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		slog.Error("failed to parse order_id", "order_id", req.OrderID, "error", err)
		writeError(w, http.StatusBadRequest, "invalid order_id")
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		slog.Error("failed to parse user_id", "user_id", req.UserID, "error", err)
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	result, err := h.paymentService.InitPayment(r.Context(), payment.PaymentRequest{
		OrderID:        orderID,
		UserID:         userID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Provider:       req.Provider,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, InitPaymentResponse{
		PaymentID:  result.PaymentID.String(),
		PaymentURL: result.PaymentURL,
		Status:     result.Status,
		Existing:   result.Existing,
	})
}

func (h *PaymentHandler) ProcessWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read webhook payload")
		return
	}

	if err := h.paymentService.ProcessWebhook(r.Context(), chi.URLParam(r, "provider"), payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (h *PaymentHandler) GetUserPayments(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)

	payments, err := h.paymentService.GetUserPayments(r.Context(), chi.URLParam(r, "userID"), limit, offset)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := make([]PaymentResponse, 0, len(payments))
	for _, p := range payments {
		result = append(result, toPaymentResponse(p))
	}
	writeJSON(w, http.StatusOK, result)
}

func toPaymentResponse(p payment.Payment) PaymentResponse {
	resp := PaymentResponse{
		ID:             p.ID.String(),
		OrderID:        p.OrderID.String(),
		UserID:         p.UserID.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		Provider:       p.Provider,
		ProviderID:     p.ProviderID,
		Status:         p.Status,
		IdempotencyKey: p.IdempotencyKey,
	}
	if !p.CreatedAt.IsZero() {
		resp.CreatedAt = p.CreatedAt.Format(time.RFC3339)
	}
	if !p.UpdatedAt.IsZero() {
		resp.UpdatedAt = p.UpdatedAt.Format(time.RFC3339)
	}
	return resp
}

func queryInt(r *http.Request, name string, fallback int) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Status: "Error", Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
