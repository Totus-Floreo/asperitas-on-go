package model

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	User Author `json:"user"`
	jwt.RegisteredClaims
}

func NewTokenClaims(now time.Time, id string, username string) TokenClaims {
	return TokenClaims{
		Author{
			ID:       id,
			Username: username,
		},
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)),
		},
	}
}

type ITokenStorage interface {
	GetToken(context.Context, string) (string, error)
	SetToken(context.Context, string, string) error
	DeleteToken(context.Context, string) error
}
