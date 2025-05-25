//go:build windows
// +build windows

package pipe

import (
	"github.com/GregoireW/oidc-helper/internal/logutil"
	"github.com/Microsoft/go-winio"
	"os"
)

func RunDaemon(pipeName string) {
	pipePath := `\\.\pipe\` + pipeName
	ln, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		logutil.Logf(logutil.LogError, "Windows pipe listen error: %v", err)
		os.Exit(1)
	}
	defer ln.Close()
	ts := NewTokenStore()
	logutil.Logf(logutil.LogInfo, "Windows daemon listening on %s", pipePath)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go ts.HandleConnection(conn)
	}
}

func ConnectToDaemonWithProvider(pipeName string, provider string) (string, error) {
	pipePath := `\\.\pipe\` + pipeName
	conn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.Write([]byte("GET " + provider))
	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	return string(buf[:n]), nil
}

func StoreToken(pipeName *string, setCmd string) {
	pipePath := `\\.\pipe\` + *pipeName
	conn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		logutil.Logf(logutil.LogError, "Failed to connect to daemon for storing token: %v", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte(setCmd))
	buf := make([]byte, 4096)
	conn.Read(buf) // ignore response
}
