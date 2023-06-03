package model

import "context"

type IAuthService interface {
	LogIn(context.Context, string, string) (string, error)
	SignUp(context.Context, string, string) (string, error)
}
