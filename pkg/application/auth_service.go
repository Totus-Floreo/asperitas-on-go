package application

import (
	"log"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
)

type AuthService struct {
	userStorage model.IUserStorage
	jwtService  model.IJWTService
}

func NewAuthService(userStorage model.IUserStorage, jwtService model.IJWTService) *AuthService {
	return &AuthService{
		userStorage: userStorage,
		jwtService:  jwtService,
	}
}

func (s *AuthService) SignUp(username string, password string) (string, error) {
	if _, err := s.userStorage.GetUser(username); err == nil {
		return "", model.ErrUserExist
	}

	user := &model.User{
		Username: username,
		Password: password,
	}
	if err := s.userStorage.AddUser(user); err != nil {
		return "", err
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) LogIn(username, password string) (string, error) {
	user, err := s.userStorage.GetUser(username)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}
	log.Println(user.Password)
	log.Println(password)
	if user.Password != password {
		return "", model.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}
