package controller

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("API Gateway is running!"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}
