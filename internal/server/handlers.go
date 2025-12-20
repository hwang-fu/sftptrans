package server

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func writeJSON(writer http.ResponseWriter, status int, response APIResponse) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(response)
}

func writeError(writer http.ResponseWriter, status int, message string, code string) {
	writeJSON(writer, status, APIResponse{
		Success: false,
		Error:   &APIError{Message: message, Code: code},
	})
}
