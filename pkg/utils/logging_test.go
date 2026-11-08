package utils

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

// captureStd redirects stdout and stderr while fn runs and returns what each received.
func captureStd(t *testing.T, fn func()) (string, string) {
	t.Helper()

	origOut, origErr := os.Stdout, os.Stderr
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = outW, errW
	defer func() {
		os.Stdout, os.Stderr = origOut, origErr
	}()

	fn()

	_ = outW.Close()
	_ = errW.Close()
	outBytes, _ := io.ReadAll(outR)
	errBytes, _ := io.ReadAll(errR)

	return string(outBytes), string(errBytes)
}

func TestConfigureDefaultLoggerWritesToStderr(t *testing.T) {
	origLogger := log.Logger
	defer func() {
		log.Logger = origLogger
	}()

	stdout, stderr := captureStd(t, func() {
		ConfigureDefaultLogger("info")
		log.Error().Msg("boom")
	})

	if stdout != "" {
		t.Errorf("expected nothing on stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "boom") {
		t.Errorf("expected log line on stderr, got %q", stderr)
	}
}

func TestNewWritesToStderr(t *testing.T) {
	for _, pretty := range []bool{true, false} {
		stdout, stderr := captureStd(t, func() {
			logger := New("info", pretty)
			logger.Error().Msg("boom")
		})

		if stdout != "" {
			t.Errorf("pretty=%v: expected nothing on stdout, got %q", pretty, stdout)
		}
		if !strings.Contains(stderr, "boom") {
			t.Errorf("pretty=%v: expected log line on stderr, got %q", pretty, stderr)
		}
	}
}

func TestLogsHaveNoColorWhenStderrIsNotATerminal(t *testing.T) {
	_, stderr := captureStd(t, func() {
		logger := New("info", true)
		logger.Error().Msg("boom")
	})

	if strings.Contains(stderr, "\x1b[") {
		t.Errorf("expected no ANSI escapes on a non-terminal stderr, got %q", stderr)
	}
}
