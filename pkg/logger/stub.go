package logger

import (
	"context"
	"io"
	"log/slog"
)

type StubHandler struct{}

func NewStubHandler(_ io.Writer, _ *slog.HandlerOptions) *StubHandler {
	return &StubHandler{}
}

func (h *StubHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *StubHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *StubHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *StubHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
