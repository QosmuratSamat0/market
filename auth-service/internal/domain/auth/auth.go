package domain

import (
	"time"
)

type RefreshToken struct {
	ExpiresAt time.Time
	UserID    string
	Token     string
}
