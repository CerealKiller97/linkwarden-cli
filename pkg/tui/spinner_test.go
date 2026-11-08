package tui

import (
	"bytes"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

func TestSpinnerModelHiddenUntilShownAndAfterStop(t *testing.T) {
	assert := require.New(t)

	m := newSpinnerModel("Fetching links", time.Hour)
	assert.Empty(m.View().Content, "nothing drawn before the delay passes")

	next, _ := m.Update(spinnerShowMsg{})
	m = next.(spinnerModel)
	assert.Contains(m.View().Content, "Fetching links")

	next, cmd := m.Update(spinnerStopMsg{})
	m = next.(spinnerModel)
	assert.Empty(m.View().Content, "stopping erases the spinner")
	assert.NotNil(cmd, "stopping quits the program")
}

func TestSpinnerStopIsIdempotent(t *testing.T) {
	var out bytes.Buffer
	s := startSpinner("Work", 0, &out, tea.WithWindowSize(80, 24))

	// Wait until the spinner is visible rather than sleeping a fixed time, which
	// is too short under -race on a slow CI runner.
	s.program.Send(spinnerShowMsg{})
	s.Stop()
	s.Stop()

	final, ok := s.final.(spinnerModel)
	require.True(t, ok)
	require.True(t, final.shown)
	require.True(t, final.stopped)
	require.Empty(t, final.View().Content, "the spinner is erased once stopped")
}

func TestSpinnerStopBeforeDelayDrawsNothing(t *testing.T) {
	var out bytes.Buffer
	s := startSpinner("Work", time.Hour, &out, tea.WithWindowSize(80, 24))
	s.Stop()

	final, ok := s.final.(spinnerModel)
	require.True(t, ok)
	require.False(t, final.shown, "stopped before the delay passed")
	require.NotContains(t, out.String(), "Work")
}
