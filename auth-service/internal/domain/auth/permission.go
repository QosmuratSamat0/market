package domain

type Permission string

const (
	PermRead   Permission = "read"
	PermWrite  Permission = "write"
	PermUpdate Permission = "update"
	PermDelete Permission = "delete"
)
