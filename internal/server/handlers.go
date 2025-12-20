package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"sftptrans/internal/localfs"
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

func handleRemoteMkdir(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	sess := session.Current()
	if err := sess.Client().MkDir(req.Path); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_MKDIR_ERROR")
		return
	}
	writeSuccess(w, nil)
}

func handleRemoteRename(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	sess := session.Current()
	if err := sess.Client().Rename(req.OldPath, req.NewPath); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_RENAME_ERROR")
		return
	}
	writeSuccess(w, nil)
}

func handleRemoteDelete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, http.StatusBadRequest, "Path is required", "INVALID_REQUEST")
		return
	}

	sess := session.Current()
	if err := sess.Client().Delete(path); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_DELETE_ERROR")
		return
	}
	writeSuccess(w, nil)
}

func handleRemoteDownload(w http.ResponseWriter, r *http.Request) {
	remotePath := r.URL.Query().Get("path")
	if remotePath == "" {
		writeError(w, http.StatusBadRequest, "Path is required", "INVALID_REQUEST")
		return
	}

	sess := session.Current()
	filename := filepath.Base(remotePath)
	localPath := filepath.Join(sess.DownloadDir(), filename)

	if err := sess.Client().Download(remotePath, localPath); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_DOWNLOAD_ERROR")
		return
	}

	writeSuccess(w, map[string]string{"localPath": localPath})
}

func handleRemoteUpload(w http.ResponseWriter, r *http.Request) {
	remotePath := r.URL.Query().Get("path")
	if remotePath == "" {
		writeError(w, http.StatusBadRequest, "Remote path is required", "INVALID_REQUEST")
		return
	}

	// Parse multipart form (max 1GB)
	if err := r.ParseMultipartForm(1 << 30); err != nil {
		writeError(w, http.StatusBadRequest, "Failed to parse form", "INVALID_REQUEST")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "No file provided", "INVALID_REQUEST")
		return
	}
	defer file.Close()

	// Save to temp file first
	tmpFile, err := os.CreateTemp("", "sftptrans-upload-*")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "LOCAL_ERROR")
		return
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.ReadFrom(file); err != nil {
		tmpFile.Close()
		writeError(w, http.StatusInternalServerError, err.Error(), "LOCAL_ERROR")
		return
	}
	tmpFile.Close()

	// Upload to remote
	fullRemotePath := filepath.Join(remotePath, header.Filename)
	sess := session.Current()
	if err := sess.Client().Upload(tmpPath, fullRemotePath); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "SFTP_UPLOAD_ERROR")
		return
	}

	writeSuccess(w, map[string]string{"remotePath": fullRemotePath})
}

func handleLocalList(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = localfs.GetHomeDir()
	}

	entries, err := localfs.ListDir(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), "LOCAL_LIST_ERROR")
		return
	}
	writeSuccess(w, entries)
}

// Status and control

func handleStatus(w http.ResponseWriter, r *http.Request) {
	sess := session.Current()
	writeSuccess(w, map[string]any{
		"connected":   true,
		"connection":  sess.Client().ConnectionInfo(),
		"downloadDir": sess.DownloadDir(),
	})
}

var shutdownChan chan struct{}

func SetShutdownChan(ch chan struct{}) {
	shutdownChan = ch
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	writeSuccess(w, map[string]string{"message": "Shutting down..."})
	if shutdownChan != nil {
		close(shutdownChan)
	}
}
