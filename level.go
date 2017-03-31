package logtee

import (
	"strings"
)

type Level int

const (
	UnknownLevel = iota
	BizLevel
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case BizLevel:
		return "biz"
	}
	return "unknown"
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "biz":
		return BizLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn", "warning":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return UnknownLevel
	}
}
