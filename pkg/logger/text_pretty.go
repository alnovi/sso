package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type TextPrettyHandler struct {
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

func NewTextPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *TextPrettyHandler {
	return &TextPrettyHandler{
		Handler: slog.NewJSONHandler(out, opts),
		l:       log.New(out, "", 0),
	}
}

func (h *TextPrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.Marshal(fields)
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *TextPrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TextPrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   append(h.attrs, attrs...),
	}
}

func (h *TextPrettyHandler) WithGroup(name string) slog.Handler {
	return &TextPrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
