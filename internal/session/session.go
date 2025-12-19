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
