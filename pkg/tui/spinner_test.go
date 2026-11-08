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

	m := newSpinnerModel("Fetching links", time.Hour, 0)
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
	s := startSpinner("Work", 0, 0, &out, tea.WithWindowSize(80, 24))

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
	s := startSpinner("Work", time.Hour, time.Hour, &out, tea.WithWindowSize(80, 24))
	s.Stop()

	final, ok := s.final.(spinnerModel)
	require.True(t, ok)
	require.False(t, final.shown, "stopped before the delay passed")
	require.NotContains(t, out.String(), "Work")
}

func TestSpinnerStaysVisibleForMinimumTime(t *testing.T) {
	assert := require.New(t)

	m := newSpinnerModel("Fetching links", 0, time.Hour)
	next, _ := m.Update(spinnerShowMsg{})
	m = next.(spinnerModel)

	next, cmd := m.Update(spinnerStopMsg{})
	m = next.(spinnerModel)
	assert.False(m.stopped, "a just-shown spinner isn't stopped yet")
	assert.Contains(m.View().Content, "Fetching links", "the label stays readable")
	assert.NotNil(cmd, "the stop is retried once the minimum time has passed")

	m.shownAt = time.Now().Add(-2 * time.Hour)
	next, _ = m.Update(spinnerStopMsg{})
	m = next.(spinnerModel)
	assert.True(m.stopped, "stops once the minimum time has passed")
	assert.Empty(m.View().Content)
}

func TestSpinnerStopWaitsForMinimumTime(t *testing.T) {
	var out bytes.Buffer
	const minVisible = 50 * time.Millisecond
	s := startSpinner("Work", 0, minVisible, &out, tea.WithWindowSize(80, 24))

	s.program.Send(spinnerShowMsg{})
	start := time.Now()
	s.Stop()

	require.GreaterOrEqual(t, time.Since(start), minVisible-10*time.Millisecond)
	final, ok := s.final.(spinnerModel)
	require.True(t, ok)
	require.True(t, final.stopped)
}
