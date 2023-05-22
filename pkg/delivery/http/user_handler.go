package route

import (
	"encoding/json"
	"net/http"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/application"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/delivery/helpers"
	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"

	"go.uber.org/zap"
)

type UserHandler struct {
	Logger      *zap.SugaredLogger
	AuthService *application.AuthService
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	token, err := h.AuthService.SignUp(ctx, user.Username, user.Password)
	if err == model.ErrUserExist {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	token, err := h.AuthService.LogIn(ctx, user.Username, user.Password)
	if err != nil {
		if err == model.ErrInvalidCredentials {
			http.Error(w, `{"message": "invalid username or password"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
