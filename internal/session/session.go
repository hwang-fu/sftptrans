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

func Current() *Session {
	return current
}

func (s *Session) Client() *sftp.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.client
}
