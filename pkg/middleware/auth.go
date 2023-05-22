package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
	"github.com/gorilla/mux"
)

type AuthContextKey string

const (
	AuthorContextKey AuthContextKey = "author"
)

func Auth(jwtService model.IJWTService, tokenStorage model.ITokenStorage) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Print("auth: no token")
				return
			}

			token := authHeader[len("Bearer "):]
			author, err := jwtService.VerifyToken(token)
			if err != nil {
				http.Error(w, `{"message": "no auth"}`, http.StatusUnauthorized)
				return
			}

			dbtoken, err := tokenStorage.GetToken(r.Context(), author.ID)
			if err != nil || dbtoken != token {
				http.Error(w, `{"message": "bad token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthorContextKey, author)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
