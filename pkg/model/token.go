package model

import "context"

type ITokenStorage interface {
	GetToken(context.Context, string) (string, error)
	CreateToken(context.Context, string, string) error
	DeleteToken(context.Context, string) error
}

type IJWTService interface {
	GenerateToken(*User) (string, error)
	VerifyToken(string) (*Author, error)
}
