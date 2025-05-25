//go:build !windows
// +build !windows

package pipe

import (
	"github.com/GregoireW/oidc-helper/internal/logutil"
	"net"
	"os"
)

func RunDaemon(pipeName string) {
	sockPath := "/tmp/" + pipeName + ".sock"
	os.Remove(sockPath)
	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		logutil.Logf(logutil.LogError, "Unix socket listen error: %v", err)
		os.Exit(1)
	}
	defer ln.Close()
	defer os.Remove(sockPath)
	ts := NewTokenStore()
	logutil.Logf(logutil.LogInfo, "Unix daemon listening on %s", sockPath)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go ts.HandleConnection(conn)
	}
}

func ConnectToDaemonWithProvider(pipeName string, provider string) (string, error) {
	sockPath := "/tmp/" + pipeName + ".sock"
	conn, err := net.Dial("unix", sockPath)
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
	sockPath := "/tmp/" + *pipeName + ".sock"
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		logutil.Logf(logutil.LogError, "Failed to connect to daemon for storing token: %v", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte(setCmd))
	buf := make([]byte, 4096)
	conn.Read(buf) // ignore response
}
