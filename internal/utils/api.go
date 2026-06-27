package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var errorResponse string = `{"error": "%s"}`

func WriteJsonResponse[T any](w http.ResponseWriter, payload T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	http.Error(w, fmt.Sprintf(errorResponse, message), code)
}
