package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sftptrans/internal/config"
	"sftptrans/internal/server"
	"sftptrans/internal/session"
	"sftptrans/internal/sftp"
)

func main() {
	settings, err := config.ParseSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUsage: sftptrans -host <host> -user <user> [-password <pass> | -key <keyfile>]\n")
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	switch settings.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	slog.Info("Connecting to SFTP server", "host", settings.Host, "port", settings.Port, "user", settings.User)
	client, err := sftp.NewClient(settings.Host, settings.Port, settings.User, settings.Password, settings.IdentityFile)
	if err != nil {
		slog.Error("Failed to connect to SFTP server", "error", err)
		os.Exit(1)
	}
	defer client.Close()
	slog.Info("Connected successfully", "connection", client.ConnectionInfo())

	session.Initialize(client, settings.DownloadDir)

	shutdownChan := make(chan struct{})
	server.SetShutdownChan(shutdownChan)

	srv := server.NewServer(settings.ListenAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting HTTP server", "address", settings.ListenAddr)
		fmt.Printf("\n  sftptrans is running!\n")
		fmt.Printf("  Open http://localhost%s in your browser\n\n", settings.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	select {
	case <-sigChan:
		slog.Info("Received shutdown signal")
	case <-shutdownChan:
		slog.Info("Shutdown requested via API")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
	}

	session.Current().Close()
	slog.Info("Goodbye!")
}
