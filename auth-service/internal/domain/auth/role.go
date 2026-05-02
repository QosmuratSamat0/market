package domain

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
	RoleManager Role = "manager"
)
