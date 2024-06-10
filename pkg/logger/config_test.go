package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Level(t *testing.T) {
	testCases := []struct {
		name   string
		cfg    Config
		expect slog.Level
	}{
		{
			name:   "Default",
			cfg:    Config{},
			expect: slog.LevelError,
		},
		{
			name:   "Debug",
			cfg:    Config{Level: LevelDebug},
			expect: slog.LevelDebug,
		},
		{
			name:   "Info",
			cfg:    Config{Level: LevelInfo},
			expect: slog.LevelInfo,
		},
		{
			name:   "Warning",
			cfg:    Config{Level: LevelWarn},
			expect: slog.LevelWarn,
		},
		{
			name:   "Error",
			cfg:    Config{Level: LevelError},
			expect: slog.LevelError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.cfg.level(), tc.expect)
		})
	}
}

func TestConfig_Handler(t *testing.T) {
	testCases := []struct {
		name   string
		cfg    Config
		expect slog.Handler
	}{
		{
			name:   "Default",
			cfg:    Config{},
			expect: &slog.JSONHandler{},
		},
		{
			name:   "Json",
			cfg:    Config{Format: FormatJson},
			expect: &slog.JSONHandler{},
		},
		{
			name:   "Text",
			cfg:    Config{Format: FormatText},
			expect: &slog.TextHandler{},
		},
		{
			name:   "Pretty",
			cfg:    Config{Format: FormatPretty},
			expect: &PrettyHandler{},
		},
		{
			name:   "Stub",
			cfg:    Config{Format: FormatStub},
			expect: &StubHandler{},
		},
	}

	opts := slog.HandlerOptions{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.IsType(t, tc.cfg.handler(opts), tc.expect)
		})
	}
}
