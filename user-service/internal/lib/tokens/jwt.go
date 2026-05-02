package tokens

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string   `json:"user-id"`
	Role   string   `json:"role"`
	Perms  []string `json:"perms"`
	jwt.RegisteredClaims
}

func ParseJWT(tokenStr, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
