package utils

import (
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DateTimeFormat Default datetime format
const DateTimeFormat = "2006-01-02 15:04:05"

// ConfigureDefaultLogger Set default options for Zerolog with log level and pretty print
// Logs always go to stderr so they never mix with command output on stdout.
func ConfigureDefaultLogger(level string, output ...io.Writer) {
	zerologLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("Failed to parse logging level: " + level)
	}

	zerolog.SetGlobalLevel(zerologLevel)
	zerolog.TimeFieldFormat = DateTimeFormat
	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	var w io.Writer = newStderrConsoleWriter()

	if len(output) > 0 {
		w = zerolog.MultiLevelWriter(w, output[0])
	}

	log.Logger = log.Output(w)
}

// newStderrConsoleWriter returns a human-readable log writer on stderr, coloured
// only when stderr is a terminal and NO_COLOR is unset, so redirected logs
// (2> err.log) don't fill up with escape codes.
func newStderrConsoleWriter() zerolog.ConsoleWriter {
	return zerolog.NewConsoleWriter(func(cw *zerolog.ConsoleWriter) {
		cw.Out = os.Stderr
		cw.NoColor = !colorEnabledFor(os.Getenv("NO_COLOR"), isatty.IsTerminal(os.Stderr.Fd()))
	})
}

// New Returns an instance of zerolog.Logger with configured log level and pretty print flag
// Logs always go to stderr so they never mix with command output on stdout.
func New(level string, prettyPrint bool) zerolog.Logger {
	var logger zerolog.Logger

	zerologLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("Failed to parse logging level: " + level)
	}

	var w io.Writer

	if prettyPrint {
		w = newStderrConsoleWriter()
	} else {
		w = os.Stderr
	}

	logger = zerolog.New(w).
		With().
		Timestamp().
		Logger().
		Level(zerologLevel)

	return logger
}
