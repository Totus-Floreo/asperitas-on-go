package application

import (
	"context"
	"errors"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	userStorage  model.IUserStorage
	tokenStorage model.ITokenStorage
	jwtService   model.IJWTService
}

func NewAuthService(userStorage model.IUserStorage, tokenStorage model.ITokenStorage, jwtService model.IJWTService) *AuthService {
	return &AuthService{
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		jwtService:   jwtService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, username string, password string) (string, error) {
	if _, err := s.userStorage.GetUser(ctx, username); err == nil {
		return "", model.ErrUserExist
	}

	user := &model.User{
		Username: username,
		Password: password,
	}
	if err := s.userStorage.AddUser(ctx, user); err != nil {
		return "", err
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}

	if err := s.tokenStorage.SetToken(ctx, user.ID, token); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) LogIn(ctx context.Context, username, password string) (string, error) {
	user, err := s.userStorage.GetUser(ctx, username)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}
	if user.Password != password {
		return "", model.ErrInvalidCredentials
	}

	token, err := s.tokenStorage.GetToken(ctx, user.ID)
	if errors.Is(err, redis.Nil) {
		token, err := s.jwtService.GenerateToken(user)
		if err != nil {
			return "", err
		}
		if err := s.tokenStorage.SetToken(ctx, user.ID, token); err != nil {
			return "", err
		}
		return token, nil
	} else if err != nil {
		return "", err
	}
	return token, nil
}
