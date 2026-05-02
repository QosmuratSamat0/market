package token

import (
	"context"
	"time"

	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
	token "github.com/QosmuratSamat0/auth-service/internal/lib/tokens"
)

type Service struct {
	tokenRepo RefreshTokenRepository
	secret    string
}

func NewService(
	tokenRepo RefreshTokenRepository,
	secret string,
) *Service {
	return &Service{
		tokenRepo: tokenRepo,
		secret:    secret,
	}
}

func (s *Service) SaveRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error {
	return s.tokenRepo.SaveRefreshToken(ctx, token, userID, expiresAt)
}

func (s *Service) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	return s.tokenRepo.GetRefreshToken(ctx, token)
}

func (s *Service) DeleteRefreshToken(ctx context.Context, token string) (string, error) {
	return s.tokenRepo.DeleteRefreshToken(ctx, token)
}

func (s *Service) GenerateAccessToken(ctx context.Context, userID string, role string) (string, error) {
	accessToken, err := token.GenerateJWT(
		s.secret,
		userID,
		role,
		15*time.Minute,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *Service) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = s.tokenRepo.SaveRefreshToken(
		ctx,
		newRefreshToken,
		userID,
		time.Now().Add(30*24*time.Hour),
	)
	if err != nil {
		return "", err
	}
	return newRefreshToken, nil
}
