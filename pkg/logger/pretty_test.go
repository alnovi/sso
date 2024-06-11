package logger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPrettyHandler(t *testing.T) {
	assert.NotNil(t, NewPrettyHandler(os.Stdout, &slog.HandlerOptions{}))
}

func TestPrettyHandler_Handle(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string
		rec  slog.Record
		exp  string
	}{
		{
			name: "Debug",
			rec:  slog.NewRecord(now, slog.LevelDebug, "message", 0),
			exp:  fmt.Sprintf("[%s] DEBUG: message \n", now.Format("2006-01-02 15:04:05 -0700")),
		},
		{
			name: "Info",
			rec:  slog.NewRecord(now, slog.LevelInfo, "message", 0),
			exp:  fmt.Sprintf("[%s] INFO: message \n", now.Format("2006-01-02 15:04:05 -0700")),
		},
		{
			name: "Warn",
			rec:  slog.NewRecord(now, slog.LevelWarn, "message", 0),
			exp:  fmt.Sprintf("[%s] WARN: message \n", now.Format("2006-01-02 15:04:05 -0700")),
		},
		{
			name: "Error",
			rec:  slog.NewRecord(now, slog.LevelError, "message", 0),
			exp:  fmt.Sprintf("[%s] ERROR: message \n", now.Format("2006-01-02 15:04:05 -0700")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			h := NewPrettyHandler(b, &slog.HandlerOptions{})
			err := h.Handle(context.Background(), tc.rec)
			assert.NoError(t, err)
			assert.Equal(t, tc.exp, b.String())
		})
	}
}

func TestPrettyHandler_WithAttrs(t *testing.T) {
	var h slog.Handler
	now := time.Now()
	b := new(bytes.Buffer)
	h = NewPrettyHandler(b, &slog.HandlerOptions{})
	h = h.WithAttrs([]slog.Attr{slog.String("key", "value")})
	rec := slog.NewRecord(now, slog.LevelDebug, "message", 0)
	err := h.Handle(context.Background(), rec)

	assert.NoError(t, err)

	expect := fmt.Sprintf("[%s] DEBUG: message {\"key\":\"value\"}\n", now.Format("2006-01-02 15:04:05 -0700"))
	actual := b.String()

	assert.Equal(t, expect, actual)
}

func TestPrettyHandler_WithGroup(t *testing.T) {
	var h slog.Handler
	now := time.Now()
	b := new(bytes.Buffer)
	h = NewPrettyHandler(b, &slog.HandlerOptions{})
	h = h.WithGroup("group")
	rec := slog.NewRecord(now, slog.LevelDebug, "message", 0)
	err := h.Handle(context.Background(), rec)

	assert.NoError(t, err)

	expect := fmt.Sprintf("[%s] DEBUG: message \n", now.Format("2006-01-02 15:04:05 -0700"))
	actual := b.String()

	assert.Equal(t, expect, actual)
}
