package auth

import (
	"context"
	"fmt"

	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
	"github.com/QosmuratSamat0/auth-service/internal/lib/errs"
	"github.com/QosmuratSamat0/auth-service/internal/lib/passwordUtils"
)

type UseCase struct {
	tokenService TokenService
	userClient   UserClient
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

func NewUseCase(tokenService TokenService, userClient UserClient) *UseCase {
	return &UseCase{
		tokenService: tokenService,
		userClient:   userClient,
	}
}

func (s *UseCase) Register(ctx context.Context, input RegisterInput) error {
	existing, _ := s.userClient.GetUserByEmail(ctx, input.Email)
	if existing != nil {
		return errs.EmailAlreadyExists
	}

	hashedPassword, err := passwordUtils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
	}

	return s.userClient.CreateUser(ctx, user)
}

func (s *UseCase) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := s.userClient.GetUserByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, errs.ErrInvalidCredentials
	}

	if !passwordUtils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UseCase) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.tokenService.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return errs.ErrInvalidRefreshToken
	}
	return nil
}

func (s *UseCase) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	userID, err := s.tokenService.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errs.ErrInvalidRefreshToken
	}

	user, err := s.userClient.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errs.UserNotFound
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
