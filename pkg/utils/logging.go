package utils

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// DateTimeFormat Default datetime format
const DateTimeFormat = "2006-01-02 15:04:05"

// ConfigureDefaultLogger sets the default slog logger to a colorized tint handler on
// standard output with the given level. If output is given, records are also written
// to it as JSON.
func ConfigureDefaultLogger(level string, output ...io.Writer) {
	var handler slog.Handler = newConsoleHandler(os.Stdout, parseLevel(level))

	if len(output) > 0 {
		handler = slog.NewMultiHandler(handler, newJSONHandler(output[0], parseLevel(level)))
	}

	slog.SetDefault(slog.New(handler))
}

// New Returns an instance of *slog.Logger with configured log level and pretty print flag
// If pretty print is false, it will write JSON to Standard Output
func New(level string, prettyPrint bool) *slog.Logger {
	if prettyPrint {
		return slog.New(newConsoleHandler(os.Stdout, parseLevel(level)))
	}

	return slog.New(newJSONHandler(os.Stdout, parseLevel(level)))
}

func newConsoleHandler(w *os.File, level slog.Level) slog.Handler {
	return tint.NewTextHandler(w, &tint.Options{
		Level:      level,
		TimeFormat: DateTimeFormat,
		NoColor:    !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()),
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if _, ok := a.Value.Any().(error); ok && a.Value.Kind() == slog.KindAny {
				return tint.Attr(9, a) // red, like zerolog's console error field
			}
			return a
		},
	})
}

func newJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(slog.TimeKey, a.Value.Time().UTC().Format(DateTimeFormat))
			}
			return a
		},
	})
}

func parseLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		panic("Failed to parse logging level: " + level)
	}

	return l
}
