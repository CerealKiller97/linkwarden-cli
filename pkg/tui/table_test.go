package tui

import (
	"fmt"
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
	view := m.View()
	assert.Contains(view.Content, "Links (2)")
	assert.Contains(view.Content, "Bubble Tea", "rows are visible, not just the header")
	assert.Contains(view.Content, "q quit")
	assert.False(view.AltScreen, "renders inline, like the Bubble Tea table example")

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(1, m.(browseModel).table.Cursor())

	m, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	assert.NotNil(cmd)
	assert.IsType(tea.QuitMsg{}, cmd())
	assert.NotContains(m.View().Content, "q quit", "help line is dropped once you quit")
	assert.Contains(m.View().Content, "Bubble Tea", "the table itself stays on screen")
}

func TestVisibleRows(t *testing.T) {
	assert := require.New(t)

	assert.Equal(17, visibleRows(100, 24), "fills the screen minus chrome")
	assert.Equal(2+1, visibleRows(2, 24), "never below 3")
	assert.Equal(10, visibleRows(10, 24), "never more rows than there are")
}

func TestBrowse_ShowsAsManyRowsAsFit(t *testing.T) {
	assert := require.New(t)

	big := Table{Title: "T", Headers: []string{"N"}}
	for i := range 50 {
		big.Rows = append(big.Rows, []string{fmt.Sprintf("row-%02d", i)})
	}

	view := newBrowseModel(big, 80, 24).View().Content
	lines := strings.Split(strings.TrimRight(view, "\n"), "\n")

	assert.Contains(view, "row-16")
	assert.NotContains(view, "row-17", "exactly visibleRows (17) rows for a 24-row terminal")
	assert.Len(lines, 24-1, "fills the terminal, leaving one line for the prompt")
}

func TestNewBrowseModel_WidthFitsContent(t *testing.T) {
	m := newBrowseModel(sample, 200, 24)
	// "ID"(2) + "Bubble Tea"(10), each with 2 cells of padding.
	require.Equal(t, 2+2+10+2, tableWidth(m.table.Columns()))
}

func TestRenderStatic_TonesColorCells(t *testing.T) {
	assert := require.New(t)

	tbl := Table{
		Title:   "Results",
		Headers: []string{"Success", "Failed"},
		Rows:    [][]string{{"9", "1"}},
		Tone: func(_, col int) Tone {
			if col == 1 {
				return ToneDanger
			}
			return ToneSuccess
		},
	}

	colored := RenderStatic(tbl, true)
	assert.Contains(colored, toneStyles[ToneSuccess].Render("9")[:10], "success cell uses the success color")
	assert.NotContains(RenderStatic(tbl, false), "\x1b[", "tones disappear without color")
}
