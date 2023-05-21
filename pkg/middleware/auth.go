package middleware

import (
	"context"
	"encoding/json"
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
				http.Redirect(w, r, "/", http.StatusUnauthorized)
				return
			}
			//  i don't shure about this block, but rn i have no idea how to do this
			dbtoken, err := tokenStorage.GetToken(r.Context(), author.ID)
			if err != nil || dbtoken != token {
				resp, errMarshal := json.Marshal(map[string]interface{}{"message": "bad token"})
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
				http.Redirect(w, r, "/", http.StatusUnauthorized)
				return
			}
			// end of conversation
			ctx := context.WithValue(r.Context(), AuthorContextKey, author)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
