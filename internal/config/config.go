// Package config
package config

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	User           string
	Host           string
	Port           int
	Password       string
	KeyPath        string
	ListenAddress  string
	DownloadFolder string
	LogLevel       string
}

func Parse() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.User, "user", "", "SFTP host username (required)")
	flag.StringVar(&cfg.Host, "host", "", "SFTP host address (required)")
	flag.IntVar(&cfg.Port, "port", 22, "SFTP port (optional, default 22)")
	flag.StringVar(&cfg.Password, "password", "", "SFTP host password")
	flag.StringVar(&cfg.KeyPath, "keypath", "", "SFTP host key file path")
	flag.StringVar(&cfg.ListenAddress, "listen", ":8080", "HTTP listen address")
	flag.StringVar(&cfg.DownloadFolder, "downloadto", "", "Default download folder (optional, default $HOME/tmp")
	flag.StringVar(&cfg.LogLevel, "log", "info", "Log level: debug|info|warn|error (optional, default info)")

	flag.Parse()

	if cfg.Host == "" {
		return nil, errors.New("host is required")
	}
	if cfg.User == "" {
		return nil, errors.New("user is required")
	}
	if cfg.Password == "" && cfg.KeyPath == "" {
		return nil, errors.New("either a password or a key file path is required")
	}

	if cfg.DownloadFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "/tmp"
		}
		cfg.DownloadFolder = filepath.Join(home, "temporary")
	}

	if err := os.MkdirAll(cfg.DownloadFolder, 0o755); err != nil {
		return nil, err
	}

	return cfg, nil
}
