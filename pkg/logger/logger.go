package logger

import "log/slog"

func New(format, level string) *slog.Logger {
	config := Config{
		Format: format,
		Level:  level,
	}

	options := slog.HandlerOptions{Level: config.level()}
	handler := config.handler(options)

	return slog.New(handler)
}
