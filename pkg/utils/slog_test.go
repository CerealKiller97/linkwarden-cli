package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewSlogLogger_ForwardsToZerolog(t *testing.T) {
	assert := require.New(t)

	var buf bytes.Buffer
	logger := NewSlogLogger(zerolog.New(&buf).Level(zerolog.InfoLevel))

	logger.With("svc", "sdk").WithGroup("req").Error("boom",
		slog.Any("error", errors.New("closed")),
		slog.Int("status", 500),
		slog.Group("hdr", slog.String("ct", "json")),
	)

	var got map[string]any
	assert.NoError(json.Unmarshal(buf.Bytes(), &got))
	assert.Equal("error", got["level"])
	assert.Equal("boom", got["message"])
	assert.Equal("sdk", got["svc"])
	assert.Equal("closed", got["req.error"])
	assert.InDelta(500, got["req.status"], 0)
	assert.Equal("json", got["req.hdr.ct"])
}

func TestNewSlogLogger_RespectsZerologLevel(t *testing.T) {
	assert := require.New(t)

	var buf bytes.Buffer
	logger := NewSlogLogger(zerolog.New(&buf).Level(zerolog.WarnLevel))

	assert.False(logger.Enabled(t.Context(), slog.LevelInfo))
	logger.Info("hidden")
	assert.Empty(buf.String())

	logger.Warn("shown")
	assert.Contains(buf.String(), "shown")
}
