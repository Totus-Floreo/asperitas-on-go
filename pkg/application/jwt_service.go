package application

import (
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"

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
	t := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
		"iat": t.Unix(),
		"exp": t.Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*model.Author, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

	if payload, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, model.ErrInvalidToken
	} else {
		user := payload["user"].(map[string]interface{})
		author := &model.Author{
			ID:       user["id"].(string),
			Username: user["username"].(string),
		}
		return author, nil
	}
}
