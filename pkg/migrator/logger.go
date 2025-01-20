package migrator

import (
	"fmt"
	"log/slog"
	"strings"
)

type GooseLogger struct {
	log *slog.Logger
}

func NewGooseLogger(log *slog.Logger) *GooseLogger {
	return &GooseLogger{log: log}
}

func (l *GooseLogger) Fatalf(format string, v ...interface{}) {
	msg := l.clearMsg(fmt.Sprintf(format, v...))
	l.log.Error(msg)
}

func (l *GooseLogger) Printf(format string, v ...interface{}) {
	msg := l.clearMsg(fmt.Sprintf(format, v...))
	l.log.Info(msg)
}

func (l *GooseLogger) clearMsg(msg string) string {
	msg = strings.TrimPrefix(msg, "goose: ")
	msg = strings.TrimSuffix(msg, "\n")
	return msg
}
