package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

var sample = Table{
	Title:   "Links",
	Headers: []string{"ID", "Name"},
	Rows:    [][]string{{"1", "Bubble Tea"}, {"2", "Lip Gloss"}},
}

func TestRenderStatic_Plain(t *testing.T) {
	assert := require.New(t)

	out := RenderStatic(sample, false)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	assert.Equal("Links", lines[0])
	assert.Len(lines, TableChromeHeight+len(sample.Rows))
	assert.Contains(out, "Bubble Tea")
	assert.NotContains(out, "\x1b[")
}

func TestRenderStatic_Color(t *testing.T) {
	require.Contains(t, RenderStatic(sample, true), "\x1b[")
}

func TestFitColumns_ShrinksToWidth(t *testing.T) {
	assert := require.New(t)

	wide := Table{Headers: []string{"ID", "URL"}, Rows: [][]string{{"1", strings.Repeat("x", 200)}}}
	cols := fitColumns(wide, 80)

	total := 0
	for _, c := range cols {
		total += c.Width + 2
	}
	assert.LessOrEqual(total, 80)
	assert.Equal(2, cols[0].Width, "narrow columns keep their natural width")
}

func TestBrowse_ScrollAndQuit(t *testing.T) {
	assert := require.New(t)

	var m tea.Model = newBrowseModel(sample, 80, 24)
	assert.Contains(m.View().Content, "Bubble Tea", "rows are visible, not just the header")

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Contains(m.View().Content, "2/2")
	assert.True(m.View().AltScreen)

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	assert.NotNil(cmd)
	assert.IsType(tea.QuitMsg{}, cmd())
}
