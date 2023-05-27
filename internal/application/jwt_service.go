package application

import (
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}

func (s *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := model.NewTokenClaims(time.Now(), user.ID, user.Username)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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

	if payload, ok := token.Claims.(*model.TokenClaims); !ok {
		return nil, model.ErrInvalidToken
	} else {
		author := payload.User
		return &author, nil
	}
}
