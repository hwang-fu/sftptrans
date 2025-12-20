package server

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

var staticFiles embed.FS

func NewServer(listenAddr string) *http.Server {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("", handleRemoteList)

	// Static files (Angular SPA)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		slog.Error("Failed to load static files", "error", err)
	}

	return &http.Server{}
}

// spaHandler serves static files and falls back to index.html for SPA routing
type spaHandler struct {
	staticFS http.Handler
}
