package application

import (
	"github.com/Totus-Floreo/asperitas-on-go/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey      string
	method         jwt.SigningMethod
	timeController model.ITimeController
}

func NewJWTService(secretKey string, method jwt.SigningMethod, timeController model.ITimeController) *JWTService {
	return &JWTService{
		secretKey:      secretKey,
		method:         method,
		timeController: timeController,
	}
}

func (s *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := model.NewTokenClaims(s.timeController.Now(), user.ID, user.Username)
	token := jwt.NewWithClaims(s.method, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*model.Author, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrInvalidSignMethod
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, model.ErrInvalidToken
	}

	if payload, ok := token.Claims.(*model.TokenClaims); ok {
		author := payload.User
		return &author, nil
	}

	return nil, model.ErrInvalidToken
}
