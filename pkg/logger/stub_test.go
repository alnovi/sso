package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	h   = NewStubHandler(os.Stdout, &slog.HandlerOptions{})
	ctx = context.Background()
	rec = slog.Record{Time: time.Now(), Message: "test message", Level: slog.LevelDebug}
)

func TestStubHandler_Handle(t *testing.T) {
	assert.Nil(t, h.Handle(ctx, rec))
}

func TestStubHandler_WithAttrs(t *testing.T) {
	attr := []slog.Attr{slog.String("msg", "text")}
	assert.Equal(t, h.WithAttrs(attr), h)
}

func TestStubHandler_WithGroup(t *testing.T) {
	g := "group name"
	assert.Equal(t, h.WithGroup(g), h)
}

func TestStubHandler_Enabled(t *testing.T) {
	assert.False(t, h.Enabled(ctx, slog.LevelDebug))
}
