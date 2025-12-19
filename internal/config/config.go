// Package config
package config

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

type Settings struct {
	User         string
	Host         string
	Port         int
	IdentityFile string
	Password     string
	ListenAddr   string
	DownloadDir  string
	LogLevel     string
}

func ParseSettings() (*Settings, error) {
	settings := &Settings{}

	flag.StringVar(&settings.User, "u", "", "SFTP host username (required)")
	flag.StringVar(&settings.Host, "h", "", "SFTP host address (required)")
	flag.IntVar(&settings.Port, "p", 22, "SFTP port (optional, default 22)")
	flag.StringVar(&settings.IdentityFile, "i", "", "SFTP host identity file path")
	flag.StringVar(&settings.Password, "pass", "", "SFTP host password")
	flag.StringVar(&settings.ListenAddr, "listen", ":8080", "HTTP listen address")
	flag.StringVar(&settings.DownloadDir, "outdir", "", "Default download folder (optional, default $HOME/tmp")
	flag.StringVar(&settings.LogLevel, "log", "info", "Log level: debug|info|warn|error (optional, default info)")

	flag.Parse()

	if settings.Host == "" {
		return nil, errors.New("host is required")
	}
	if settings.User == "" {
		return nil, errors.New("username is required")
	}
	if settings.Password == "" && settings.IdentityFile == "" {
		return nil, errors.New("either a password or an identity file path is required")
	}
	if settings.Password != "" && settings.IdentityFile != "" {
		return nil, errors.New("either a password OR an identity file path is required, not both")
	}

	if settings.DownloadDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "/tmp"
		}
		settings.DownloadDir = filepath.Join(home, "temporary")
	}

	if err := os.MkdirAll(settings.DownloadDir, 0o755); err != nil {
		return nil, err
	}

	return settings, nil
}
