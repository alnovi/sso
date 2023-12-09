package logger

import (
	"log/slog"
	"os"
)

const (
	FormatJson = "json"
	FormatText = "text"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func NewLogger(format, level string) *slog.Logger {
	options := slog.HandlerOptions{Level: getSlogLevel(level)}
	handler := getSlogHandler(format, &options)

	return slog.New(handler)
}

func getSlogLevel(level string) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func getSlogHandler(format string, opts *slog.HandlerOptions) slog.Handler {
	switch format {
	case FormatJson:
		return slog.NewJSONHandler(os.Stdout, opts)
	case FormatText:
		return slog.NewTextHandler(os.Stdout, opts)
	default:
		return slog.NewJSONHandler(os.Stdout, opts)
	}
}
