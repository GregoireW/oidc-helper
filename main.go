package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/GregoireW/oidc-helper/internal/config"
	daemoninternal "github.com/GregoireW/oidc-helper/internal/daemon"
	"github.com/GregoireW/oidc-helper/internal/logutil"
	"github.com/GregoireW/oidc-helper/internal/pipe"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var Version = "dev"

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	_ = exec.Command(cmd, args...).Start()
}

func main() {
	versionFlag := flag.Bool("version", false, "show version and exit")
	daemon := flag.Bool("daemon", false, "run as daemon to hold token in memory")
	pipeName := flag.String("pipe", "oidc-helper-pipe", "named pipe/socket name for daemon communication")
	help := flag.Bool("help", false, "display help and exit")
	h := flag.Bool("h", false, "display help and exit (shorthand)")
	logFlag := flag.String("log", "warn", "set log level: debug, info, warn, error")
	providerName := flag.String("provider", "", "OIDC provider name (overrides default)")
	listProviders := flag.Bool("list-providers", false, "list available OIDC providers and exit")

	logutil.LogLevel = logutil.ParseLogLevel(*logFlag)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "oidc-helper version: %s\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *listProviders {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available providers:")
		for name := range cfg.Providers {
			if name == cfg.Default {
				fmt.Println("-", name, "*")
			} else {
				fmt.Println("-", name)
			}
		}
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *help || *h {
		flag.Usage()
		os.Exit(0)
	}

	// Perform authentication as before
	cfg, err := config.LoadConfig()
	if err != nil {
		logutil.Logf(logutil.LogError, "Error loading config: %v", err)
		os.Exit(1)
	}

	// Select provider
	selectedProvider := cfg.Default
	if *providerName != "" {
		selectedProvider = *providerName
	}
	prov, ok := cfg.Providers[selectedProvider]
	if !ok {
		logutil.Logf(logutil.LogError, "Provider '%s' not found in config", selectedProvider)
		os.Exit(1)
	}

	if *daemon {
		pipe.RunDaemon(*pipeName)
		return
	}

	// Try to get token from daemon for the selected provider
	tokenStr, err := pipe.ConnectToDaemonWithProvider(*pipeName, selectedProvider)
	if err == nil && tokenStr != "" {
		fmt.Println(tokenStr)
		return
	}

	// If daemon not running, launch it
	if err != nil {
		err = daemoninternal.LaunchDaemon(*pipeName)
		if err != nil {
			logutil.Logf(logutil.LogError, "Failed to launch daemon: %v", err)
			os.Exit(1)
		}

		// Wait a moment for daemon to start
		time.Sleep(500 * time.Millisecond)
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, prov.OIDCUrl)
	if err != nil {
		logutil.Logf(logutil.LogError, "OIDC provider error: %v", err)
		os.Exit(1)
	}

	oauth2Cfg := oauth2.Config{
		ClientID:    prov.ClientID,
		Endpoint:    provider.Endpoint(),
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: "http://localhost:36547/callback",
	}
	authURL := oauth2Cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	logutil.Logf(logutil.LogDebug, "Open this link in your browser:\n%s", authURL)
	openBrowser(authURL)

	type codeResult struct {
		code string
		err  error
	}
	resultCh := make(chan codeResult)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			resultCh <- codeResult{"", fmt.Errorf("OAuth error: %s", errMsg)}
			logutil.Logf(logutil.LogError, "OAuth error: %s", errMsg)
			fmt.Fprintf(w, "OAuth error: %s", errMsg)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			resultCh <- codeResult{"", fmt.Errorf("No code in request")}
			logutil.Logf(logutil.LogError, "No code in request")
			fmt.Fprintf(w, "No code in request")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html><html><body><p>Authentication successful. You may close this window.</p><script>window.close();</script></body></html>`)
		resultCh <- codeResult{code, nil}
	})

	server := &http.Server{Addr: ":36547"}
	go server.ListenAndServe()
	res := <-resultCh
	_ = server.Close()
	if res.err != nil {
		logutil.Logf(logutil.LogError, "Error receiving code: %v", res.err)
		os.Exit(1)
	}

	token, err := oauth2Cfg.Exchange(ctx, res.code)
	if err != nil {
		logutil.Logf(logutil.LogError, "Error exchanging code/token: %v", err)
		os.Exit(1)
	}

	// Send token to daemon for storage
	expiry := token.Expiry.Format(time.RFC3339)
	setCmd := "SET " + selectedProvider + "|" + token.AccessToken + "|" + expiry
	pipe.StoreToken(pipeName, setCmd)

	fmt.Println(token.AccessToken)
}
