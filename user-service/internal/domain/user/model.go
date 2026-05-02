package user

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
	RoleManager Role = "manager"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
