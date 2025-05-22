package logger

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
	"time"
	"unicode"
)

var emptyAttr = slog.Attr{}

type PrettyHandler struct {
	w          io.Writer
	mu         *sync.Mutex
	opts       *Options
	bufferPool *bufferPool

	attrsPrefix string
	groupPrefix string
	groups      []string
}

func NewPrettyHandler(w io.Writer, opts *Options) slog.Handler {
	if opts == nil {
		opts = new(Options)
	}

	h := &PrettyHandler{
		w:          w,
		mu:         new(sync.Mutex),
		opts:       opts,
		bufferPool: newBufferPool(),
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	if h.opts.TimeFormatter == nil {
		h.opts.TimeFormatter = DefaultTimeFormatter(time.DateTime)
	}
	if h.opts.LevelFormatter == nil {
		h.opts.LevelFormatter = DefaultLevelFormatter(!h.opts.DisableColor)
	}
	if h.opts.SourceFormatter == nil {
		h.opts.SourceFormatter = DefaultSourceFormatter(!h.opts.DisableColor)
	}
	return h
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, record slog.Record) error {
	rep := h.opts.ReplaceAttr
	buf := h.bufferPool.Acquire()
	defer h.bufferPool.Free(buf)

	// Time
	h.appendTime(buf, rep, record)

	// Level
	if rep == nil {
		h.opts.LevelFormatter(buf, record.Level)
	} else if a := rep(nil, slog.Any(slog.LevelKey, record.Level)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
	}
	buf.AppendByte(' ')

	// Source
	h.appendSource(buf, rep, record)

	// Message
	if rep == nil {
		buf.AppendString(record.Message)
	} else if a := rep(nil, slog.String(slog.MessageKey, record.Message)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
	}
	buf.AppendByte(' ')

	// handler attributes
	if len(h.attrsPrefix) > 0 {
		buf.AppendString(h.attrsPrefix)
	}

	// Write attributes
	record.Attrs(func(attr slog.Attr) bool {
		if rep != nil {
			attr = rep(h.groups, attr)
		}
		h.appendAttr(buf, attr, h.groupPrefix)
		return true
	})

	if buf.Len() == 0 {
		return nil
	}
	buf.Replace(buf.Len()-1, '\n') // Replace the last space with a newline

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := buf.WriteTo(h.w)
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()

	buf := h.bufferPool.Acquire()
	defer h.bufferPool.Free(buf)

	for _, attr := range attrs {
		if h.opts.ReplaceAttr != nil {
			attr = h.opts.ReplaceAttr(h.groups, attr)
		}
		h.appendAttr(buf, attr, h.groupPrefix)
	}
	h2.attrsPrefix += buf.String()
	return h2
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *PrettyHandler) clone() *PrettyHandler {
	return &PrettyHandler{
		w:           h.w,
		mu:          h.mu,
		opts:        h.opts,
		bufferPool:  h.bufferPool,
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix,
		groups:      h.groups,
	}
}

func (h *PrettyHandler) appendTime(buf *Buffer, rep ReplaceAttrFunc, record slog.Record) {
	if !record.Time.IsZero() {
		val := record.Time.Round(0)
		if rep == nil {
			h.opts.TimeFormatter(buf, val)
		} else if a := rep(nil, slog.Time(slog.TimeKey, val)); a.Key != "" {
			if a.Value.Kind() == slog.KindTime {
				h.opts.TimeFormatter(buf, a.Value.Time())
			} else {
				h.appendValue(buf, a.Value, false)
			}
		}
		buf.AppendByte(' ')
	}
}

func (h *PrettyHandler) appendSource(buf *Buffer, rep ReplaceAttrFunc, record slog.Record) {
	if h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := &slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			if rep == nil {
				h.opts.SourceFormatter(buf, src)
			} else if a := rep(nil, slog.Any(slog.SourceKey, src)); a.Key != "" {
				h.appendValue(buf, a.Value, false)
			}
			buf.AppendByte(' ')
		}
	}
}

func (h *PrettyHandler) appendAttr(buf *Buffer, attr slog.Attr, groupsPrefix string) {
	if attr.Equal(emptyAttr) {
		return
	}
	attr.Value = attr.Value.Resolve()

	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix)
		}
		return
	}

	h.appendKey(buf, attr.Key, groupsPrefix)
	h.appendValue(buf, attr.Value, true)
	buf.AppendByte(' ')
}

func (h *PrettyHandler) appendKey(buf *Buffer, key, groups string) {
	if !h.opts.DisableColor {
		buf.AppendString(ansiFaint)
		defer buf.AppendString(ansiReset)
	}
	appendString(buf, groups+key, true)
	buf.AppendByte('=')
}

// nolint: cyclop
func (h *PrettyHandler) appendValue(buf *Buffer, v slog.Value, quote bool) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote)
	case slog.KindInt64:
		buf.AppendInt(v.Int64())
	case slog.KindUint64:
		buf.AppendUint(v.Uint64())
	case slog.KindFloat64:
		buf.AppendFloat64(v.Float64())
	case slog.KindBool:
		buf.AppendBool(v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote)
	case slog.KindTime:
		appendString(buf, v.Time().String(), quote)
	case slog.KindAny, slog.KindLogValuer:
		if tm, ok := v.Any().(encoding.TextMarshaler); ok {
			b, err := tm.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(b), quote)
			return
		}

		appendString(buf, fmt.Sprint(v.Any()), quote)
	case slog.KindGroup:
		// Nothing to do
	}
}

func appendString(buf *Buffer, s string, quote bool) {
	if quote && needsQuoting(s) {
		buf.AppendQuote(s)
		return
	}
	buf.AppendString(s)
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '"' || r == '=' || !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}
