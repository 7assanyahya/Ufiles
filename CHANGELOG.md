# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project uses [Semantic Versioning](https://semver.org/).

## [1.0.0] - 2026-10-23

First stable release. 🎉

### Features
- Upload files from any device on your local network through the browser
- Recent uploads list and a server files page with one-tap download and delete
- Light and dark mode that follows your system preference
- `-host`, `-port`, `-max-size` and `-version` command-line flags
- Upload folder is created automatically on startup
- Prebuilt downloads for Windows, macOS and Linux (amd64 and arm64)

### Security
- Upload, download and delete are restricted to the upload folder (path-traversal fix)
- File names are displayed as plain text in the file list (XSS fix)
- Configurable maximum upload size (default 1 GB)

### Project
- MIT license, contributing guide, code of conduct and security policy
- Issue and pull request templates, CI workflow, project icon and screenshots

[1.0.0]: https://github.com/7assanyahya/ufiles/releases/tag/v1.0.0
