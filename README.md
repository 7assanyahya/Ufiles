# Ufiles — Local Network File Uploader

Ufiles is a lightweight Go application that lets you upload, list, download, and delete files over your local network through a simple web interface. It binds to your machine’s local IP and serves a small UI from the `static/` directory, storing uploaded files in `static/servfile/`.

Run it on your computer and access it from any device on the same Wi‑Fi/LAN.

## Features
- Upload files via the browser (drag-and-drop style input)
- View recent uploads (last 5)
- Browse all server files
- Download or delete files from the server
- Clean UI with light/dark toggle
- Simple logging to `server.log`

## Requirements
- Go 1.18+
- macOS, Linux, or Windows
- Devices must be on the same local network

## Quick Start
1. Build and run directly:
   - `go run .`
   - or: `go build -o ufiles && ./ufiles`
2. Watch the console for the URL (e.g. `http://192.XXX.X.XX:8080`).
3. Open that URL from any device on the same network.

By default, the app binds to your machine’s local IP on port 8080 and writes logs to `server.log`.

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
├── static/
│   ├── index.html         # Upload UI + recent uploads
│   ├── file.html          # Server files list (download/delete)
│   ├── styles.css         # UI styles (light/dark)
│   ├── script.js          # Front-end logic
│   └── servfile/          # Uploaded files live here
└── server.log             # Created at runtime
```

## Troubleshooting
- Can’t reach the URL from another device? Confirm you’re on the same network and that your OS firewall allows inbound connections on port 8080.
- If `static/servfile/` is missing, create it before uploading.


