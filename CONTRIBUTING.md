# Contributing to Ufiles

Thanks for your interest in improving Ufiles! Contributions of all kinds are welcome: bug reports, feature ideas, documentation, and code.

## Reporting bugs and requesting features
- Search [existing issues](../../issues) first to avoid duplicates.
- Open a new issue using the matching template and include steps to reproduce, expected vs. actual behavior, your OS, and your Go version.
- For security problems, **do not open a public issue** — see [SECURITY.md](SECURITY.md).

## Development setup
1. Install Go 1.18 or newer.
2. Fork and clone the repository.
3. Run the server locally:
   ```sh
   go run . -host 127.0.0.1
   ```
4. Open `http://127.0.0.1:8080` in your browser.

## Before opening a pull request
- Format and vet your code:
  ```sh
  gofmt -l .     # should print nothing
  go vet ./...
  go build ./...
  ```
- Keep changes focused; one feature or fix per pull request.
- Update `README.md` if you change behavior, flags, or endpoints.
- Never commit uploaded files, logs, or personal data.

## Commit messages
Use short, descriptive messages in the imperative mood, e.g. `Add upload size limit`.

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
