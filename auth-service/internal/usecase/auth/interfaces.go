package auth

import (
	"context"
	"time"

	domain "github.com/QosmuratSamat0/auth-service/internal/domain/auth"
)

type TokenService interface {
	SaveRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) (string, error)
	GenerateAccessToken(ctx context.Context, userID string, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID string) (string, error)
}


type UserClient interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         domain.Role
}
