package user

import (
	"context"

	"github.com/QosmuratSamat0/user-service-market/internal/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	UpdateUser(ctx context.Context, u *user.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]*user.User, error)
}
