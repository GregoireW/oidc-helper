package pipe

// Platform-specific pipe server/client stubs will be implemented in pipe_unix.go and pipe_windows.go

// Add shared interfaces or types here if needed

import (
	"net"
	"strings"
	"sync"
	"time"
)

// TokenStore holds tokens and their expiry times, protected by a mutex for concurrency.
type TokenStore struct {
	tokens map[string]string
	expiry map[string]time.Time
	mu     sync.RWMutex
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

// HandleConnection processes a single connection for GET/SET commands.
func (ts *TokenStore) HandleConnection(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	n, _ := c.Read(buf)
	cmd := string(buf[:n])
	if strings.HasPrefix(cmd, "GET ") {
		provider := strings.TrimSpace(cmd[4:])
		ts.mu.RLock()
		tok, ok := ts.tokens[provider]
		exp, expok := ts.expiry[provider]
		ts.mu.RUnlock()
		if ok && expok && tok != "" && time.Now().Before(exp) {
			c.Write([]byte(tok))
		} else {
			c.Write([]byte(""))
		}
	} else if strings.HasPrefix(cmd, "SET ") {
		// SET provider|token|expiry
		parts := strings.SplitN(cmd[4:], "|", 3)
		if len(parts) == 3 {
			provider := parts[0]
			ts.mu.Lock()
			ts.tokens[provider] = parts[1]
			e, _ := time.Parse(time.RFC3339, parts[2])
			ts.expiry[provider] = e
			ts.mu.Unlock()
			c.Write([]byte("OK"))
		} else {
			c.Write([]byte("ERR"))
		}
	} else {
		c.Write([]byte("ERR"))
	}
}
