package utils

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JsonError(w http.ResponseWriter, msg string, status int) {
	JsonResponse(w, map[string]string{"error": msg}, status)
}
