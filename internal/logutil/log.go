package logutil

import (
	"fmt"
	"os"
	"strings"
)

const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
)

var LogLevel = LogWarn

func ParseLogLevel(s string) int {
	switch strings.ToLower(s) {
	case "debug":
		return LogDebug
	case "info":
		return LogInfo
	case "warn":
		return LogWarn
	case "error":
		return LogError
	default:
		return LogWarn
	}
}

func Logf(level int, format string, a ...interface{}) {
	if level < LogLevel {
		return
	}
	prefix := "[INFO]"
	switch level {
	case LogDebug:
		prefix = "[DEBUG]"
	case LogInfo:
		prefix = "[INFO]"
	case LogWarn:
		prefix = "[WARN]"
	case LogError:
		prefix = "[ERROR]"
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s "+format+"\n", append([]interface{}{prefix}, a...)...)
}
