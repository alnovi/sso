package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStubHandler_Handle(t *testing.T) {
	h := NewStubHandler(os.Stdout, &slog.HandlerOptions{})
	rec := slog.Record{Time: time.Now(), Message: "test message", Level: slog.LevelDebug}
	assert.Nil(t, h.Handle(context.Background(), rec))
}

func TestStubHandler_WithAttrs(t *testing.T) {
	h := NewStubHandler(os.Stdout, &slog.HandlerOptions{})
	attr := []slog.Attr{slog.String("msg", "text")}
	assert.Equal(t, h.WithAttrs(attr), h)
}

func TestStubHandler_WithGroup(t *testing.T) {
	g := "group name"
	h := NewStubHandler(os.Stdout, &slog.HandlerOptions{})
	assert.Equal(t, h.WithGroup(g), h)
}

func TestStubHandler_Enabled(t *testing.T) {
	h := NewStubHandler(os.Stdout, &slog.HandlerOptions{})
	assert.False(t, h.Enabled(context.Background(), slog.LevelDebug))
}
