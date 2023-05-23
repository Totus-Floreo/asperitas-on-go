package helpers

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
