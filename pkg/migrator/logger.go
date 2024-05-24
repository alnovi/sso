package migrator

import (
	"fmt"
	"log/slog"
	"strings"
)

type Logger struct {
	log *slog.Logger
}

func NewLogger(log *slog.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	msg := l.clearMsg(fmt.Sprintf(format, v...))
	l.log.Error(msg)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	msg := l.clearMsg(fmt.Sprintf(format, v...))
	l.log.Info(msg)
}

func (l *Logger) clearMsg(msg string) string {
	msg = strings.TrimPrefix(msg, "goose: ")
	msg = strings.TrimSuffix(msg, "\n")
	return msg
}
