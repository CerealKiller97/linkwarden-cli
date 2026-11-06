package tui

import (
	"bytes"
	"sync"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

func TestProgressModel_CountsAndCaps(t *testing.T) {
	assert := require.New(t)

	m := newProgressModel("Work", 2)
	for range 3 {
		next, _ := m.Update(progressIncrMsg{})
		m = next.(progressModel)
	}

	assert.Equal(2, m.current, "never counts past total")
	assert.InDelta(1.0, m.percent(), 0.0001)
}

func TestProgress_ConcurrentIncrementsThenStop(t *testing.T) {
	var out bytes.Buffer
	// A bytes.Buffer has no size to detect, so give the renderer one.
	p := startProgress("Work", 50, &out, tea.WithWindowSize(80, 24))

	var wg sync.WaitGroup
	for range 50 {
		wg.Go(p.Increment)
	}
	wg.Wait()
	p.Stop()
	p.Stop() // idempotent

	require.Contains(t, out.String(), "50/50")
}
