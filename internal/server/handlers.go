package server

import (
	"encoding/json"
	"net/http"

	"sftptrans/internal/session"
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

func writeSuccess(writer http.ResponseWriter, data any) {
	writeJSON(writer, http.StatusOK, APIResponse{Success: true, Data: data})
}

func handleRemoteList(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	sess := session.Current()
	entries, err := sess.Client().ListDir(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_LIST_ERROR")
		return
	}
	writeSuccess(w, entries)
}
