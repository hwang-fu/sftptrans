package sftp

import (
	"path/filepath"
	"time"
)

type FileEntry struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	IsDir       bool      `json:"isDir"`
	ModTime     time.Time `json:"modTime"`
	Permissions string    `json:"permissions"`
}

func (c *Client) ListDir(path string) ([]FileEntry, error) {
	files, err := c.sftpClient.ReadDir(path)
	if err != nil {
		return nil, err
	}

	entries := make([]FileEntry, 0, len(files))
	for _, f := range files {
		entries = append(entries, FileEntry{
			Name:        f.Name(),
			Path:        filepath.Join(path, f.Name()),
			Size:        f.Size(),
			IsDir:       f.IsDir(),
			ModTime:     f.ModTime(),
			Permissions: f.Mode().String(),
		})
	}

	return entries, nil
}

func (c *Client) MkDir(path string) error {
	return c.sftpClient.MkdirAll(path)
}

func (c *Client) Rename(oldPath, newPath string) error {
	return c.sftpClient.Rename(oldPath, newPath)
}

func (c *Client) Delete(path string) error {
	info, err := c.sftpClient.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return c.deleteDir(path)
	}
	return c.sftpClient.Remove(path)
}

func (c *Client) deleteDir(path string) error {
	files, err := c.sftpClient.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			if err := c.deleteDir(fullPath); err != nil {
				return err
			}
		} else {
			if err := c.sftpClient.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return c.sftpClient.RemoveDirectory(path)
}

func (c *Client) Stat(path string) (*FileEntry, error) {
	info, err := c.sftpClient.Stat(path)
	if err != nil {
		return nil, err
	}
	return &FileEntry{
		Name:        info.Name(),
		Path:        path,
		Size:        info.Size(),
		IsDir:       info.IsDir(),
		ModTime:     info.ModTime(),
		Permissions: info.Mode().String(),
	}, nil
}
