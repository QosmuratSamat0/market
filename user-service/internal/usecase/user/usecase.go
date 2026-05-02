package user

import (
	"context"
	"fmt"

	"github.com/QosmuratSamat0/user-service-market/internal/domain/user"
	"github.com/QosmuratSamat0/user-service-market/internal/lib/errs"
	"github.com/QosmuratSamat0/user-service-market/internal/lib/passwordUtils"
)

type UseCase struct {
	userRepo UserRepository
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type CreateInternalUserInput struct {
	Name         string
	Email        string
	PasswordHash string
	Role         string
}

type UpdateUserInput struct {
	Name  string
	Email string
	Role  string
}

func NewUseCase(userRepo UserRepository) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

func (uc *UseCase) CreateUser(ctx context.Context, input CreateUserInput) (*user.User, error) {
	existing, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, errs.EmailAlreadyExists
	}

	passwordHash, err := passwordUtils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &user.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         user.RoleUser,
	}

	if err := uc.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UseCase) CreateInternalUser(ctx context.Context, input CreateInternalUserInput) (*user.User, error) {
	existing, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, errs.EmailAlreadyExists
	}

	if input.PasswordHash == "" {
		return nil, errs.ErrInvalidInput
	}

	role := user.Role(input.Role)
	if role == "" {
		role = user.RoleUser
	}

	u := &user.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         role,
	}

	if err := uc.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UseCase) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return uc.userRepo.GetUserByID(ctx, id)
}

func (uc *UseCase) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return uc.userRepo.GetUserByEmail(ctx, email)
}

func (uc *UseCase) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*user.User, error) {
	u, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		u.Name = input.Name
	}
	if input.Email != "" {
		u.Email = input.Email
	}
	if input.Role != "" {
		u.Role = user.Role(input.Role)
	}

	if err := uc.userRepo.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.DeleteUser(ctx, id)
}

func (uc *UseCase) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	return uc.userRepo.GetAllUsers(ctx)
}
