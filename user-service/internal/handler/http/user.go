package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/QosmuratSamat0/user-service-market/internal/domain/user"
	"github.com/QosmuratSamat0/user-service-market/internal/handler/http/middleware"
	resp "github.com/QosmuratSamat0/user-service-market/internal/lib/api/response"
	"github.com/QosmuratSamat0/user-service-market/internal/lib/errs"
	userUseCase "github.com/QosmuratSamat0/user-service-market/internal/usecase/user"
)

type UserHandler struct {
	userUC *userUseCase.UseCase
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InternalCreateUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type InternalUserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

func NewUserHandler(userUC *userUseCase.UseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

// GetMe returns the currently authenticated user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetMe"
	log := slog.With("operation", op)

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))
		return
	}

	u, err := h.userUC.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toUserResponse(u))
}

// GetUser returns a user by ID (admin/manager only)
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetUser"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("user id is required"))
		return
	}

	u, err := h.userUC.GetUserByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toUserResponse(u))
}

// GetAllUsers returns all users (admin/manager only)
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetAllUsers"
	log := slog.With("operation", op)

	users, err := h.userUC.GetAllUsers(r.Context())
	if err != nil {
		log.Error("failed to get users", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
		return
	}

	var result []UserResponse
	for _, u := range users {
		result = append(result, toUserResponse(u))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}

// CreateUser creates a new user (admin only)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.CreateUser"
	log := slog.With("operation", op)

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("name, email and password are required"))
		return
	}

	u, err := h.userUC.CreateUser(r.Context(), userUseCase.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Error("failed to create user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, toUserResponse(u))
}

// UpdateUser updates user info
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUser"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("user id is required"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	u, err := h.userUC.UpdateUser(r.Context(), id, userUseCase.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		log.Error("failed to update user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toUserResponse(u))
}

// DeleteUser deletes a user by ID (admin only)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.DeleteUser"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("user id is required"))
		return
	}

	if err := h.userUC.DeleteUser(r.Context(), id); err != nil {
		log.Error("failed to delete user", "error", err)
		handleUserError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetInternalUserByEmail(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetInternalUserByEmail"
	log := slog.With("operation", op)

	email := r.URL.Query().Get("email")
	if email == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("email is required"))
		return
	}

	u, err := h.userUC.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Error("failed to get user by email", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toInternalUserResponse(u))
}

func (h *UserHandler) GetInternalUserByID(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetInternalUserByID"
	log := slog.With("operation", op)

	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("user id is required"))
		return
	}

	u, err := h.userUC.GetUserByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toInternalUserResponse(u))
}

func (h *UserHandler) CreateInternalUser(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.CreateInternalUser"
	log := slog.With("operation", op)

	var req InternalCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	if req.Name == "" || req.Email == "" || req.PasswordHash == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("name, email and password_hash are required"))
		return
	}

	u, err := h.userUC.CreateInternalUser(r.Context(), userUseCase.CreateInternalUserInput{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		Role:         req.Role,
	})
	if err != nil {
		log.Error("failed to create internal user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, toInternalUserResponse(u))
}

// UpdateMe updates the currently authenticated user's profile
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateMe"
	log := slog.With("operation", op)

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	// Users cannot change their own role
	u, err := h.userUC.UpdateUser(r.Context(), userID, userUseCase.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		log.Error("failed to update user", "error", err)
		handleUserError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, toUserResponse(u))
}

func handleUserError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.UserNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp.Error("user not found"))

	case errors.Is(err, errs.EmailAlreadyExists):
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, resp.Error("email already exists"))

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

func toUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toInternalUserResponse(u *user.User) InternalUserResponse {
	return InternalUserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
	}
}

func (h *UserHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, []map[string]string{
		{"from": "system", "msg": "Welcome to the Market Place Chat!"},
		{"from": "admin", "msg": "How can I help you today?"},
	})
}

func (h *UserHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "sent", "msg": req.Message})
}
