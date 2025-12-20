# sftptrans

A lightweight, web-based SFTP client with a Go backend and Angular frontend.
It runs as a single binary and provides a classic FTP browser experience through your web browser.

## Features

- **Single Binary Deployment** — Go backend serves the Angular SPA and REST API together
- **Dual Authentication** — Supports both password and SSH private key authentication
- **Full File Operations** — Upload, download, rename, delete, and create directories
- **Dual-Pane Browser** — Navigate both local and remote filesystems side by side
- **Concurrent Transfers** — Leverages Go's goroutines for parallel file operations
- **Secure by Design** — No credential storage; all session data is held in memory only

## Building

```bash
# Build the Angular frontend and the Go binary (embeds the frontend)
make
```

## Usage

```
sftptrans [options]

Required:
  -h string       SFTP host address
  -u string       SFTP username

Authentication (one required):
  -pass string       SFTP password
  -key string        Path to SSH private key

Optional:
  -port int          SFTP port (default: 22)
  -listen string     HTTP listen address (default: :8080)
```

### Examples (After Building)

```bash
# Password authentication
./sftptrans -h 192.168.1.100 -u admin -pass secret

# SSH key authentication
./sftptrans -h example.com -u deploy -key ~/.ssh/id_ed25519
```

Once running, open `http://localhost:8080` in your browser.

## Architecture

```
┌─────────────────────────────────────────┐
│           Go Backend (single binary)    │
│  ┌─────────────────────────────────┐    │
│  │   Static File Server (Angular)  │    │
│  └─────────────────────────────────┘    │
│  ┌─────────────────────────────────┐    │      ┌──────────────┐
│  │       REST API Handler          │◄───┼─────►│ Remote SFTP  │
│  └─────────────────────────────────┘    │      │    Server    │
│  ┌─────────────────────────────────┐    │      └──────────────┘
│  │   SFTP Client (golang.org/x)    │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
            ▲
            │ :8080
      User's Browser
```


## Tech Stack

**Backend:** Go (stdlib + golang.org/x/crypto/ssh + github.com/pkg/sftp)

**Frontend:** Angular with minimal CSS styling

## License

MIT
