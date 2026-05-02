package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "github.com/QosmuratSamat0/auth-service/internal/lib/api/response"
	"github.com/QosmuratSamat0/auth-service/internal/lib/errs"
	authUseCase "github.com/QosmuratSamat0/auth-service/internal/usecase/auth"
)

type AuthHandler struct {
	authUC   *authUseCase.UseCase
	validate *validator.Validate
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

const refreshTokenCookieName = "refresh_token"

func NewAuthHandler(authUC *authUseCase.UseCase) *AuthHandler {
	return &AuthHandler{
		authUC:   authUC,
		validate: validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   input  body      RegisterRequest  true  "Registration info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Register"
	log := slog.With("operation", op)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		log.Error("validation failed", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("all fields required"))
		return
	}

	err := h.authUC.Register(r.Context(), authUseCase.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Error("failed to register user", "error", err)
		handleAuthError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp.OK())
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   input  body      LoginRequest  true  "Login credentials"
// @Success 200 {object} AccessTokenResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Login"
	log := slog.With("operation", op)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		log.Error("validation failed", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid email format"))
		return
	}

	authResp, err := h.authUC.Login(r.Context(), authUseCase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Error("failed to login user", "error", err)
		handleAuthError(w, r, err)
		return
	}

	setRefreshTokenCookie(w, r, authResp.RefreshToken)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, AccessTokenResponse{AccessToken: authResp.AccessToken})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Produce  json
// @Param   Cookie header string true "refresh_token=<token>"
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Logout"
	log := slog.With("operation", op)

	refreshTokenCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		log.Error("failed to read refresh token cookie", "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))
		return
	}

	err = h.authUC.Logout(r.Context(), refreshTokenCookie.Value)
	if err != nil {
		log.Error("failed to logout user", "error", err)
		handleAuthError(w, r, err)
		return
	}

	clearRefreshTokenCookie(w, r)

	w.WriteHeader(http.StatusNoContent)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token cookie
// @Tags auth
// @Produce  json
// @Param   Cookie header string true "refresh_token=<token>"
// @Success 200 {object} AccessTokenResponse
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Refresh"
	log := slog.With("operation", op)

	refreshTokenCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		log.Error("failed to read refresh token cookie", "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))
		return
	}

	authResp, err := h.authUC.Refresh(r.Context(), refreshTokenCookie.Value)
	if err != nil {
		log.Error("failed to refresh token", "error", err)
		handleAuthError(w, r, err)
		return
	}

	setRefreshTokenCookie(w, r, authResp.RefreshToken)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, AccessTokenResponse{AccessToken: authResp.AccessToken})
}

func setRefreshTokenCookie(w http.ResponseWriter, r *http.Request, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

func clearRefreshTokenCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func handleAuthError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errs.ErrInvalidInput):
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid input"))

	case errors.Is(err, errs.EmailAlreadyExists):
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, resp.Error("email already exists"))

	case errors.Is(err, errs.UserNotFound):
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp.Error("user not found"))

	case errors.Is(err, errs.ErrUnauthorized),
		errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrInvalidRefreshToken):
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("unauthorized"))

	case errors.Is(err, errs.ErrTooManyLoginAttempt):
		render.Status(r, http.StatusTooManyRequests)
		render.JSON(w, r, resp.Error("too many login attempts, please try again later"))

	default:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("internal error"))
	}
}