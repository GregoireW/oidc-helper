# oidc-helper

A simple CLI tool to fetch OIDC (OpenID Connect) access tokens from any compatible provider. Ideal for developers and operators who need to quickly obtain tokens for testing or automation.

---

## Features
- Loads OIDC provider and client info from a YAML config file
- Discovers OIDC endpoints automatically
- Guides you through the browser-based authentication flow
- Retrieves and prints your access token

---

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/GregoireW/oidc-helper.git
   cd oidc-helper
   ```
2. **Build:**
   ```sh
   go build -o oidc-helper
   ./oidc-helper
   ```

---

## Configuration

The tool looks for a `config.yaml` file in the following standard locations depending on your operating system:

- **Linux:**
  - `$XDG_CONFIG_HOME/oidc-helper/config.yaml` (if `XDG_CONFIG_HOME` is set)
  - `$HOME/.config/oidc-helper/config.yaml`
- **macOS:**
  - `$HOME/Library/Application Support/oidc-helper/config.yaml`
  - `$XDG_CONFIG_HOME/oidc-helper/config.yaml` (if `XDG_CONFIG_HOME` is set)
- **Windows:**
  - `%APPDATA%\oidc-helper\config.yaml`
- **Fallback:**
  - The directory containing the `oidc-helper` executable

Create a `config.yaml` file in one of these locations with the following content:

```yaml
default: "provider1"
providers:
  provider1:
    oidc_url: "https://your-oidc-provider.com"
    client_id: "your-client-id"
  provider2:
    oidc_url: "https://another-provider.com"
    client_id: "another-client-id"
```

- `default`: The name of the default provider to use if none is specified.
- `providers`: A map of provider names to their OIDC configuration.
- `oidc_url`: The base URL of your OIDC provider (e.g., Auth0, Google, Okta)
- `client_id`: The client ID registered with your provider

---

## Usage

1. Fill in `config.yaml` with your provider details.
2. Run the tool:
   ```sh
   ./oidc-helper [options]
   ```
   
   **Options:**
   - `--daemon`   Run as daemon to hold token in memory
   - `--pipe`     Named pipe/socket name for daemon communication (default: "oidc-helper-pipe")
   - `--provider` OIDC provider name (overrides default in config)
   - `--log`      Set log level: `debug`, `info`, `warn`, `error` (default: `warn`)
   - `--help`, `-h`   Display help and exit

3. Follow the on-screen instructions:
   - Open the provided URL in your browser
   - Log in and authorize
   - Paste the returned authorization code into the CLI
4. The tool will display your access token.

---

## Example

```sh
$ ./oidc-helper
eyJhbGciOi... (the access token)
```

---

## Daemon & Pipe System

The oidc-helper can run as a background daemon to securely store and serve OIDC tokens to local clients via a pipe system. This is useful for sharing tokens between processes without exposing them on disk.

### How it works
- **Unix:** Uses a Unix domain socket at `/tmp/<pipeName>.sock`.
- **Windows:** Uses a named pipe at `\\.\pipe\<pipeName>`.
- The daemon supports two commands (provider name required):
  - `GET <provider>`: Returns the current valid token for the given provider (if available and not expired).
  - `SET <provider>|<token>|<expiry>`: Stores a new token and its expiry time (RFC3339 format) for the given provider.

### Example Usage
1. **Start the daemon:**
   - The daemon will listen for connections from local clients.
2. **Store a token:**
   - Send `SET provider1|<token>|<expiry>` to the pipe to store a token for provider1.
3. **Fetch a token:**
   - Send `GET provider1` to the pipe to retrieve the current token for provider1.

This system allows multiple processes to share OIDC tokens for multiple providers securely in memory, without writing them to disk.

---

## Troubleshooting
- Ensure your `config.yaml` is correct and matches your OIDC provider settings
- Your client must be registered for the "authorization code" flow
- If you see errors, check your network connection and provider status

---

## Dependencies
- [coreos/go-oidc](https://github.com/coreos/go-oidc)
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)

---
```
