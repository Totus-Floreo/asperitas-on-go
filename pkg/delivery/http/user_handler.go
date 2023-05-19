package route

import (
	"encoding/json"
	"net/http"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/application"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"

	"go.uber.org/zap"
)

type UserHandler struct {
	Logger      *zap.SugaredLogger
	AuthService *application.AuthService
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.AuthService.SignUp(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.LogIn(user.Username, user.Password)
	if err != nil {
		if err == model.ErrInvalidCredentials {
			http.Error(w, `{"message": "invalid username or password"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
