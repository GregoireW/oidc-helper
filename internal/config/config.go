package config

import (
	"fmt"
	"github.com/GregoireW/oidc-helper/internal/logutil"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
)

type Provider struct {
	OIDCUrl  string `yaml:"oidc_url"`
	ClientID string `yaml:"client_id"`
}

type Config struct {
	Providers map[string]Provider `yaml:"providers"`
	Default   string              `yaml:"default"`
}

// Example YAML structure:
// default: "provider1"
// providers:
//   provider1:
//     oidc_url: "https://example.com"
//     client_id: "abc123"
//   provider2:
//     oidc_url: "https://another.com"
//     client_id: "def456"

func configPath() (string, error) {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			paths = append(paths, appData+"/oidc-helper/config.yaml")
		}
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		if homeDir != "" {
			paths = append(paths, homeDir+"/Library/Application Support/oidc-helper/config.yaml")
		}
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig != "" {
			paths = append(paths, xdgConfig+"/oidc-helper/config.yaml")
		}
	case "linux":
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		homeDir, _ := os.UserHomeDir()
		if xdgConfig != "" {
			paths = append(paths, xdgConfig+"/oidc-helper/config.yaml")
		}
		if homeDir != "" {
			paths = append(paths, homeDir+"/.config/oidc-helper/config.yaml")
		}
	}

	// Fallback: executable directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := execPath
		if idx := len(execPath) - len("/oidc-helper"); idx > 0 {
			execDir = execPath[:idx]
		}
		paths = append(paths, execDir+"/config.yaml")
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("config.yaml not found in standard config locations or executable directory")
}

func LoadConfig() (*Config, error) {
	configPath, err := configPath()
	if err != nil {
		logutil.Logf(logutil.LogError, "Error finding config: %v", err)
		os.Exit(1)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
