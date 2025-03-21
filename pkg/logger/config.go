package logger

import (
	"log/slog"
	"os"
)

const (
	FormatJson       = "json"
	FormatText       = "text"
	FormatJsonPretty = "json-pretty"
	FormatTextPretty = "text-pretty"
	FormatDiscard    = "discard"

	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Config struct {
	Format string
	Level  string
}

func (c Config) level() slog.Level {
	switch c.Level {
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

func (c Config) handler(opts slog.HandlerOptions) slog.Handler {
	switch c.Format {
	case FormatJson:
		return slog.NewJSONHandler(os.Stdout, &opts)
	case FormatText:
		return slog.NewTextHandler(os.Stdout, &opts)
	case FormatJsonPretty:
		return NewJsonPrettyHandler(os.Stdout, &opts)
	case FormatTextPretty:
		return NewTextPrettyHandler(os.Stdout, &opts)
	case FormatDiscard:
		return slog.DiscardHandler
	default:
		return slog.NewJSONHandler(os.Stdout, &opts)
	}
}
