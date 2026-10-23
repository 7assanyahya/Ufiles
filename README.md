<p align="center">
  <img src="assets/icon.svg" alt="Ufiles logo" width="120" height="120">
</p>

<h1 align="center">Ufiles</h1>

<p align="center">
  A lightweight local network file uploader written in Go.<br>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
  <a href="https://github.com/7assanyahya/ufiles/actions/workflows/ci.yml"><img src="https://github.com/7assanyahya/ufiles/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
</p>

Ufiles is a lightweight Go application that lets you upload, list, download, and delete files over your local network through a simple web interface. It binds to your machine’s local IP and serves a small UI from the `static/` directory, storing uploaded files in `static/servfile/`.

Run it on your computer and access it from any device on the same Wi‑Fi/LAN.

## Features
- Upload files via the browser (drag-and-drop style input)
- View recent uploads (last 5)
- Browse all server files
- Download or delete files from the server
- Clean UI with light/dark toggle
- Simple logging to `server.log`
- Path-traversal protection on upload, download, and delete

## Requirements
- Go 1.18+
- macOS, Linux, or Windows
- Devices must be on the same local network

## Quick Start
1. Build and run directly:
   - `go run .`
   - or: `go build -o ufiles && ./ufiles`
2. Watch the console for the URL (e.g. `http://192.168.1.20:8080`).
3. Open that URL from any device on the same network.

By default, the app binds to your machine’s local IP on port 8080 and writes logs to `server.log`.

## Options
| Flag | Default | Description |
| --- | --- | --- |
| `-host` | local network IP | Address to listen on (use `127.0.0.1` for local-only testing) |
| `-port` | `8080` | Port to listen on |
| `-max-size` | `1024` | Maximum upload size in MB |

Example: `go run . -port 9000 -max-size 4096`

## Usage
- Upload: open the root page and choose a file to upload. Files are saved to `static/servfile/`.
- Server files: click “Server Files” or go to `/file.html` to list, download, or delete files.
- Recent uploads: the home page shows the last 5 uploaded files.

## Endpoints
- `GET /` — Serves the web UI from `static/`
- `POST /upload` — Upload a file (form field name: `file`)
- `GET /last-uploaded` — JSON list of the 5 most recent files
- `GET /server-files` — JSON list of all files on the server
- `DELETE /delete?file=<name>` — Delete a file by name
- `GET /download?file=<name>` — Download a file by name

Notes:
- The server uses your local IP (not `localhost`) so other devices on the same LAN can reach it.
- No authentication is implemented. Only run on trusted networks.
- Ensure your firewall allows inbound traffic on port 8080 if needed.

## File Locations
- Web assets: `static/`
- Uploaded files: `static/servfile/`
- Log file: `server.log`

## Project Structure (overview)
```
.
├── main.go
├── assets/
│   └── icon.svg           # Project logo
├── static/
│   ├── index.html         # Upload UI + recent uploads
│   ├── file.html          # Server files list (download/delete)
│   ├── styles.css         # UI styles (light/dark)
│   ├── script.js          # Front-end logic
│   ├── favicon.svg        # Browser tab icon
│   └── servfile/          # Uploaded files live here
└── server.log             # Created at runtime
```

## Troubleshooting
- Can’t reach the URL from another device? Confirm you’re on the same network and that your OS firewall allows inbound connections on port 8080.
- The `static/servfile/` folder is created automatically on startup.
- Port already in use? Start with a different one: `go run . -port 9000`.

## Contributing
Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) and our [Code of Conduct](CODE_OF_CONDUCT.md) before opening an issue or pull request. To report a security problem, see [SECURITY.md](SECURITY.md).

## License
Ufiles is released under the [MIT License](LICENSE).


