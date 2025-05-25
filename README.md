# oidc-helper

A simple command-line tool to help you authenticate with OpenID Connect (OIDC) providers and manage tokens easily.

---

## 🚀 Quick Start

1. **Download the Latest Release:**
   - Visit the [Releases page](https://github.com/GregoireW/oidc-helper/releases) and download the latest version for your operating system.
   - Extract the archive and place the `oidc-helper` executable somewhere in your `PATH` (e.g., `/usr/local/bin` on Linux/macOS).

2. **Configure the Tool:**
   - Create a configuration file named `config.yaml` (see below for details).

3. **Run the Tool:**
   ```sh
   ./oidc-helper [options]
   ```

---

## ⚙️ Configuration

The tool looks for a `config.yaml` file in these locations (in order):

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

**Example `config.yaml`:**

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
- `oidc_url`: The base URL of your OIDC provider (e.g., Auth0, Google, Okta).
- `client_id`: The client ID registered with your provider.

---

## 📝 Usage

1. Fill in `config.yaml` with your provider details.
2. Run the tool:
   ```sh
   ./oidc-helper [options]
   ```

**Common options:**
- `--provider <name>`: Use a specific provider from your config file (overrides default).
- `--daemon`: Run as a daemon to hold the token in memory (internal use only).
- `--pipe <name>`: Named pipe/socket name for daemon communication (default: oidc-helper-pipe).
- `--log <level>`: Set log level (`debug`, `info`, `warn`, `error`; default: `warn`).
- `--help`, `-h`: Show all available options and usage information.

For help and available options, run:
```sh
./oidc-helper --help
```

---

## 💡 Examples

- **Authenticate with the default provider and print the access token:**
  ```sh
  ./oidc-helper
  ```

- **Authenticate with a specific provider:**
  ```sh
  ./oidc-helper --provider provider2
  ```

- **Show all available options:**
  ```sh
  ./oidc-helper --help
  ```

---

## 📚 Documentation

- See [TECHNICAL.md](documentations/TECHNICAL.md) for advanced usage and development details.
- For issues or feature requests, visit the [GitHub Issues page](https://github.com/GregoireW/oidc-helper/issues).

---

## 🛡️ License

This project is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
