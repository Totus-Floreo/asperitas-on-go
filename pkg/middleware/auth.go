package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"

	"github.com/gorilla/mux"
)

func Auth(jwtService model.IJWTService) mux.MiddlewareFunc {
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
				w.WriteHeader(http.StatusUnauthorized)
				resp, errMarshal := json.Marshal(map[string]interface{}{"message": "no auth"})
				if errMarshal != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Print("auth: marsh error")
					return
				}

				_, errWrite := w.Write(resp)
				if errWrite != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Print("auth: body write error")
					return
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), "author", author)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
