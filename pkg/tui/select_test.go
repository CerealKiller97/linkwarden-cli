package tui

import (
	"bytes"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

func press(m tea.Model, keys ...tea.KeyPressMsg) tea.Model {
	for _, k := range keys {
		m, _ = m.Update(k)
	}

	return m
}

func typeText(s string) []tea.KeyPressMsg {
	keys := make([]tea.KeyPressMsg, 0, len(s))
	for _, r := range s {
		keys = append(keys, tea.KeyPressMsg{Code: r, Text: string(r)})
	}

	return keys
}

var (
	keyEnter = tea.KeyPressMsg{Code: tea.KeyEnter}
	keyDown  = tea.KeyPressMsg{Code: tea.KeyDown}
	keyUp    = tea.KeyPressMsg{Code: tea.KeyUp}
	keyEsc   = tea.KeyPressMsg{Code: tea.KeyEscape}
	keyBack  = tea.KeyPressMsg{Code: tea.KeyBackspace}
)

func TestSelect_MoveAndPick(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a", "b", "c"}), keyDown, keyDown, keyDown, keyUp, keyEnter).(selectModel)

	assert.True(m.done)
	assert.Equal(1, m.choice, "cursor clamps at the last item, then moves back up one")
}

func TestSelect_FilterNarrowsAndMapsBackToOriginalIndex(t *testing.T) {
	assert := require.New(t)

	opts := []string{"Go docs", "Rust book", "golang blog"}
	m := press(newSelectModel("Pick", opts), typeText("GO")...).(selectModel)
	assert.Equal([]int{0, 2}, m.matches, "filter is case-insensitive")

	m = press(m, keyDown, keyEnter).(selectModel)
	assert.Equal(2, m.choice)
}

func TestSelect_BackspaceWidensFilter(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"ab", "ac"}), typeText("ab")...).(selectModel)
	assert.Len(m.matches, 1)

	m = press(m, keyBack).(selectModel)
	assert.Len(m.matches, 2)
}

func TestSelect_EnterWithNoMatchesDoesNothing(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a"}), append(typeText("zzz"), keyEnter)...).(selectModel)
	assert.False(m.done)
	assert.Contains(m.View().Content, "no matches")
}

func TestSelect_Esc_Aborts(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a"}), keyEsc).(selectModel)
	assert.True(m.aborted)
}

func TestSelect_ScrollsPastPageSize(t *testing.T) {
	assert := require.New(t)

	opts := make([]string, 25)
	for i := range opts {
		opts[i] = strings.Repeat("x", i+1)
	}

	var m tea.Model = newSelectModel("Pick", opts)
	for range 15 {
		m = press(m, keyDown)
	}
	sm := m.(selectModel)

	assert.Equal(15, sm.cursor)
	assert.Equal(6, sm.offset, "cursor stays on the last visible row")
	assert.Contains(sm.View().Content, "16/25")
}

func TestSelectWith_EmptyOptions(t *testing.T) {
	_, err := SelectWith("Pick:", nil, Options{})
	require.Error(t, err)
}

// TestSelectWith_RunsProgram drives a real Bubble Tea program end-to-end through
// a fake terminal: "down, enter" picks the second option.
func TestSelectWith_RunsProgram(t *testing.T) {
	assert := require.New(t)

	in := strings.NewReader("\x1b[B\r")
	var out bytes.Buffer

	done := make(chan struct{})
	var (
		idx int
		err error
	)
	go func() {
		idx, err = SelectWith("Pick", []string{"a", "b", "c"}, Options{Input: in, Output: &out})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("select program did not exit")
	}

	assert.NoError(err)
	assert.Equal(1, idx)
}
