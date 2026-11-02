package utils

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

// NewSlogLogger returns a *slog.Logger that writes through the given zerolog logger,
// so libraries that log via slog (like the Linkwarden SDK) honour the CLI's zerolog
// level and output configuration.
func NewSlogLogger(logger zerolog.Logger) *slog.Logger {
	return slog.New(&zerologHandler{logger: logger})
}

// zerologHandler is a slog.Handler backed by a zerolog.Logger.
type zerologHandler struct {
	logger zerolog.Logger
	attrs  []slog.Attr
	group  string
}

func (h *zerologHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.WithLevel(zerologLevel(level)).Enabled()
}

func (h *zerologHandler) Handle(_ context.Context, r slog.Record) error {
	e := h.logger.WithLevel(zerologLevel(r.Level))
	if !e.Enabled() {
		return nil
	}

	for _, a := range h.attrs {
		addAttr(e, "", a)
	}
	r.Attrs(func(a slog.Attr) bool {
		addAttr(e, h.group, a)
		return true
	})

	e.Msg(r.Message)

	return nil
}

func (h *zerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	next.attrs = append(next.attrs, h.attrs...)
	for _, a := range attrs {
		a.Key = joinKey(h.group, a.Key)
		next.attrs = append(next.attrs, a)
	}

	return &next
}

func (h *zerologHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := *h
	next.group = joinKey(h.group, name)

	return &next
}

func addAttr(e *zerolog.Event, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}

	key := joinKey(prefix, a.Key)

	switch a.Value.Kind() {
	case slog.KindGroup:
		for _, ga := range a.Value.Group() {
			addAttr(e, key, ga)
		}
	case slog.KindString:
		e.Str(key, a.Value.String())
	case slog.KindInt64:
		e.Int64(key, a.Value.Int64())
	case slog.KindUint64:
		e.Uint64(key, a.Value.Uint64())
	case slog.KindFloat64:
		e.Float64(key, a.Value.Float64())
	case slog.KindBool:
		e.Bool(key, a.Value.Bool())
	case slog.KindDuration:
		e.Dur(key, a.Value.Duration())
	case slog.KindTime:
		e.Time(key, a.Value.Time())
	default:
		if err, ok := a.Value.Any().(error); ok {
			e.AnErr(key, err)
			return
		}
		e.Interface(key, a.Value.Any())
	}
}

func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	if key == "" {
		return prefix
	}

	return prefix + "." + key
}

func zerologLevel(level slog.Level) zerolog.Level {
	switch {
	case level >= slog.LevelError:
		return zerolog.ErrorLevel
	case level >= slog.LevelWarn:
		return zerolog.WarnLevel
	case level >= slog.LevelInfo:
		return zerolog.InfoLevel
	case level >= slog.LevelDebug:
		return zerolog.DebugLevel
	default:
		return zerolog.TraceLevel
	}
}
