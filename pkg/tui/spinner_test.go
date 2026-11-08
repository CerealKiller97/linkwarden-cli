package tui

import (
	"bytes"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

func TestSpinnerModel_HiddenUntilShownAndAfterStop(t *testing.T) {
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

func TestSpinner_StopIsIdempotent(t *testing.T) {
	var out bytes.Buffer
	s := startSpinner("Work", 0, &out, tea.WithWindowSize(80, 24))

	time.Sleep(50 * time.Millisecond)
	s.Stop()
	s.Stop()

	require.Contains(t, out.String(), "Work")
}

func TestSpinner_StopBeforeDelayDrawsNothing(t *testing.T) {
	var out bytes.Buffer
	s := startSpinner("Work", time.Hour, &out, tea.WithWindowSize(80, 24))
	s.Stop()

	require.NotContains(t, out.String(), "Work")
}
