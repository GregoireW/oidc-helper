# OIDC Helper – Technical Documentation

## Overview
OIDC Helper is a command-line tool designed to simplify the OpenID Connect (OIDC) authentication process for CLI applications. It automates browser-based authentication, captures the authorization code, and retrieves the access token, displaying it securely in the terminal.

## Architecture
- **Language:** Go
- **Main Components:**
  - `main.go`: Entry point, handles CLI arguments, config loading, and authentication flow.
  - `logutil/`: Custom logging utilities for consistent log formatting and level control.
  - `config.yaml`: Stores OIDC provider and client configuration.
  - `pipe_unix.go` / `pipe_windows.go`: Platform-specific helpers for inter-process communication.

## Dependencies
- [coreos/go-oidc](https://github.com/coreos/go-oidc): OIDC protocol implementation.
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2): OAuth2 client library.
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3): YAML parsing for configuration.

## Setup
1. **Clone the repository** and navigate to the project directory.
2. **Configure OIDC:**
   - Edit `config.yaml` with your OIDC provider, client ID, client secret, and redirect URI.
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Run the application:**
   ```sh
   go run main.go -log=debug
   ```

## Authentication Flow
1. Loads configuration from `config.yaml`.
2. Initializes OIDC and OAuth2 clients.
3. Starts a local HTTP server to capture the authorization code.
4. Opens the default browser to the OIDC provider's login page.
5. Receives the authorization code and exchanges it for an access token.
6. Prints the access token to the terminal (no extra formatting).

## Logging
- All output except the access token is formatted as logs (e.g., `[WARN] ...`).
- Log level is set via the `-log` CLI argument.

## Security Considerations
- **Client credentials** and **access tokens** must be kept confidential.
- The access token is only printed to stdout, with no additional text.
- Do not share `config.yaml` or access tokens.

## Extensibility & Improvements
- Advanced CLI argument parsing (e.g., with cobra).
- Secure token storage (e.g., OS keyring integration).

## License
Refer to the repository for license information.

