// Package sftp
package sftp

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sshClient  *ssh.Client  // underlying SSH connection
	sftpClient *sftp.Client // SFTP session over SSH
	host       string       // stored for ConnectionInfo()
	user       string       // stored for ConnectionInfo()
}

func NewClient(host string, port int, user string, password string, identityFile string) (*Client, error) {
	var authMethods []ssh.AuthMethod

	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if identityFile != "" {
		key, err := os.ReadFile(identityFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}
		// parses a PEM-encoded private key
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	return &Client{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		host:       host,
		user:       user,
	}, nil
}

func (c *Client) Close() error {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	if c.sshClient != nil {
		c.sshClient.Close()
	}
	return nil
}

// SFTP xposes the underlying `*sftp.Client` so you can do file operations
func (c *Client) SFTP() *sftp.Client {
	return c.sftpClient
}

func (c *Client) ConnectionInfo() string {
	return fmt.Sprintf("%s@%s", c.user, c.host)
}
