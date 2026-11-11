package tui

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	keyLeft  = tea.KeyPressMsg{Code: tea.KeyLeft}
	keyRight = tea.KeyPressMsg{Code: tea.KeyRight}
	keyEnd   = tea.KeyPressMsg{Code: tea.KeyEnd}
)

func TestSelectMoveAndPick(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a", "b", "c"}), keyDown, keyDown, keyDown, keyUp, keyEnter).(selectModel)

	assert.True(m.done)
	assert.Equal(1, m.choice, "cursor clamps at the last item, then moves back up one")
}

func TestSelectFilterNarrowsAndMapsBackToOriginalIndex(t *testing.T) {
	assert := require.New(t)

	opts := []string{"Go docs", "Rust book", "golang blog"}
	m := press(newSelectModel("Pick", opts), typeText("GO")...).(selectModel)
	assert.Equal([]int{0, 2}, m.matches, "filter is case-insensitive")

	m = press(m, keyDown, keyEnter).(selectModel)
	assert.Equal(2, m.choice)
}

func TestSelectBackspaceWidensFilter(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"ab", "ac"}), typeText("ab")...).(selectModel)
	assert.Len(m.matches, 1)

	m = press(m, keyBack).(selectModel)
	assert.Len(m.matches, 2)
}

func TestSelectEnterWithNoMatchesDoesNothing(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a"}), append(typeText("zzz"), keyEnter)...).(selectModel)
	assert.False(m.done)
	assert.Contains(m.View().Content, "no matches")
}

func TestSelectEscAborts(t *testing.T) {
	assert := require.New(t)

	m := press(newSelectModel("Pick", []string{"a"}), keyEsc).(selectModel)
	assert.True(m.aborted)
}

func TestSelectPagesWithPaginator(t *testing.T) {
	assert := require.New(t)

	opts := make([]string, 25)
	for i := range opts {
		opts[i] = fmt.Sprintf("option %02d", i+1)
	}

	var m tea.Model = newSelectModel("Pick", opts)
	for range 15 {
		m = press(m, keyDown)
	}
	sm := m.(selectModel)

	assert.Equal(15, sm.cursor)
	assert.Equal(1, sm.pages.Page, "moving past the first 10 turns the page")
	view := sm.View().Content
	assert.Contains(view, "page 2/3")
	assert.Contains(view, "option 16")
	assert.NotContains(view, "option 10", "only the current page is shown")

	sm = press(sm, keyRight).(selectModel)
	assert.Equal(2, sm.pages.Page)
	assert.Equal(24, sm.cursor, "→ moves a page, stopping at the last option")

	sm = press(sm, keyLeft, keyLeft).(selectModel)
	assert.Equal(0, sm.pages.Page)
}

func TestSelectWithEmptyOptions(t *testing.T) {
	_, err := SelectWith("Pick:", nil, Options{})
	require.Error(t, err)
}

// TestSelectWithRunsProgram drives a real Bubble Tea program end-to-end through
// a fake terminal: "down, enter" picks the second option.
func TestSelectWithRunsProgram(t *testing.T) {
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

// options returns n options labelled prefix 1…n.
func options(prefix string, from int, n int) []SelectOption {
	opts := make([]SelectOption, n)
	for i := range opts {
		opts[i] = SelectOption{Label: fmt.Sprintf("%s %d", prefix, from+i)}
	}

	return opts
}

// batches returns a LoadMore that hands out each batch in turn, reporting
// more until the last one, and counts its calls.
func batches(calls *int, bs ...[]SelectOption) LoadMore {
	return func() ([]SelectOption, bool, error) {
		*calls++
		if *calls > len(bs) {
			return nil, false, nil
		}

		return bs[*calls-1], *calls < len(bs), nil
	}
}

// settle runs cmd, feeds its message back into m, and repeats with whatever
// command that returns, like the Bubble Tea runtime would for loads.
func settle(m tea.Model, cmd tea.Cmd) tea.Model {
	for cmd != nil {
		m, cmd = m.Update(cmd())
	}

	return m
}

func TestSelectPagedLoadsNextBatchOnLastPage(t *testing.T) {
	assert := require.New(t)

	calls := 0
	m := newPagedSelectModel("Pick", options("link", 1, 20), batches(&calls, options("link", 21, 20), options("link", 41, 5)))

	assert.Nil(m.Init(), "two pages loaded and on the first: nothing to fetch yet")

	next, cmd := m.Update(keyRight)
	assert.NotNil(cmd, "reaching the last loaded page fetches the next batch")
	assert.True(next.(selectModel).loading)
	assert.Contains(next.(selectModel).View().Content, "loading more…")

	sm := settle(next, cmd).(selectModel)
	assert.Equal(1, calls)
	assert.Len(sm.options, 40)
	assert.True(sm.more)
	assert.Equal(10, sm.cursor, "the cursor stays where it was")
	assert.Contains(sm.View().Content, "page 2/4+", "+ says more can be loaded")

	next, cmd = press(sm, keyRight).(selectModel).Update(keyRight)
	assert.NotNil(cmd, "page 4 of 4 is the last loaded page")
	sm = settle(next, cmd).(selectModel)
	assert.Equal(2, calls)
	assert.Len(sm.options, 45)
	assert.False(sm.more, "the last batch said there's no more")

	_, cmd = sm.Update(keyRight)
	assert.Nil(cmd, "nothing more to load")

	sm = press(sm, keyEnd, keyEnter).(selectModel)
	assert.Equal(44, sm.choice, "the index counts across every batch")
}

func TestSelectPagedInitLoadsWhenFirstBatchFitsOnePage(t *testing.T) {
	assert := require.New(t)

	calls := 0
	m := newPagedSelectModel("Pick", options("link", 1, 3), batches(&calls, options("link", 4, 3)))

	sm := settle(m, m.Init()).(selectModel)
	assert.Equal(1, calls)
	assert.Len(sm.options, 6)
}

func TestSelectPagedFilterKeepsLoadingUntilMatchesFillAPage(t *testing.T) {
	assert := require.New(t)

	calls := 0
	later := append(options("link", 21, 19), SelectOption{Label: "needle in batch two"})
	m := newPagedSelectModel("Pick", options("link", 1, 20), batches(&calls, later, options("link", 41, 20)))

	var tm tea.Model = m
	var cmd tea.Cmd
	for _, k := range typeText("needle") {
		tm, cmd = tm.Update(k)
		tm = settle(tm, cmd)
	}
	sm := tm.(selectModel)

	assert.Equal(2, calls, "with no matches yet it loads every batch")
	assert.Equal([]int{39}, sm.matches)
	assert.Contains(sm.View().Content, "needle in batch two")
}

func TestSelectPagedShowsLoadError(t *testing.T) {
	assert := require.New(t)

	m := newPagedSelectModel("Pick", options("link", 1, 3), func() ([]SelectOption, bool, error) {
		return nil, false, errors.New("boom")
	})

	sm := settle(m, m.Init()).(selectModel)
	assert.False(sm.more, "no retrying after an error")
	assert.Contains(sm.View().Content, "couldn't load more: boom")
	assert.Len(sm.options, 3, "what was loaded can still be picked")
}

func TestSelectRowShowsDetail(t *testing.T) {
	assert := require.New(t)

	m := newPagedSelectModel("Pick", []SelectOption{{Label: "Effective Go", Detail: "go.dev/doc/effective_go"}}, nil)

	assert.Contains(m.View().Content, "Effective Go")
	assert.Contains(m.View().Content, "go.dev/doc/effective_go")

	m = press(m, typeText("effective_go")...).(selectModel)
	assert.Len(m.matches, 1, "the detail is searchable too")
}

func TestSelectRowFitsWidth(t *testing.T) {
	assert := require.New(t)

	long := SelectOption{Label: strings.Repeat("Long title ", 10), Detail: "example.com/a/long/path"}
	m := newPagedSelectModel("Pick", []SelectOption{long}, nil)
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	m = sized.(selectModel)

	row := m.row(long, false)
	assert.LessOrEqual(lipgloss.Width(row), 60)
	assert.Contains(row, "example.com/a/long/path", "the label is cut, not the URL")
	assert.Contains(row, "…")

	sized, _ = m.Update(tea.WindowSizeMsg{Width: 30, Height: 20})
	m = sized.(selectModel)
	row = m.row(long, false)
	assert.LessOrEqual(lipgloss.Width(row), 30)
	assert.NotContains(row, "example.com", "too narrow for both: the URL goes")
}
