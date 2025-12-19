package session

import (
	"sync"

	"sftptrans/internal/sftp"
)

type Session struct {
	client      *sftp.Client
	mu          sync.RWMutex
	downloadDir string
}

var (
	current *Session
	once    sync.Once
)

func Initialize(client *sftp.Client, downloadDir string) {
	once.Do(func() {
		current = &Session{
			client:      client,
			downloadDir: downloadDir,
		}
	})
}
