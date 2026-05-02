package token

import (
	"context"
	"time"

	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
)

type RefreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) (string, error)
}
