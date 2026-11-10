package utils

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestColorEnabledFor(t *testing.T) {
	assert := require.New(t)

	assert.True(colorEnabledFor("", true), "no NO_COLOR, real terminal -> color on")
	assert.False(colorEnabledFor("", false), "no NO_COLOR, not a terminal -> color off")
	assert.False(colorEnabledFor("1", true), "NO_COLOR set -> color off even on a terminal")
	assert.False(colorEnabledFor("1", false), "NO_COLOR set, not a terminal -> color off")
}

func TestShouldBrowse(t *testing.T) {
	assert := require.New(t)

	assert.False(shouldBrowse(5, 5, 40), "fits on screen -> just print")
	assert.True(shouldBrowse(100, 5, 40), "taller than the screen -> viewer")
	assert.False(shouldBrowse(100, 5, 0), "unknown height -> just print")
	assert.False(shouldBrowse(34, 5, 40), "exactly fits, with a line for the prompt")
	assert.True(shouldBrowse(34, 6, 40), "a caption line tips it over")
}

func TestTableRendersStaticWhenNotATerminal(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	tbl := NewTable(t.Context(), "Example")
	tbl.out = &out
	tbl.Header("ID", "Name")
	tbl.Row(1, "first")
	tbl.Row(2, nil)

	tbl.Render()

	assert.Equal(2, tbl.Len())
	got := out.String()
	assert.True(strings.HasPrefix(got, "Example (2)\n"), "list tables show their row count")
	assert.Contains(got, "╭", "rounded border, same as the viewer")
	assert.Contains(got, "first")
	assert.NotContains(got, "<nil>", "nil shows as EmptyCell")
	assert.NotContains(got, "\x1b[", "test output isn't a terminal, so no ANSI")
}

func TestNewSummaryHasNoCount(t *testing.T) {
	var out bytes.Buffer
	tbl := NewSummary(t.Context(), "Totals")
	tbl.out = &out
	tbl.Header("Total")
	tbl.Row(3)
	tbl.Render()

	require.True(t, strings.HasPrefix(out.String(), "Totals\n"))
}

func TestFormatCell(t *testing.T) {
	assert := require.New(t)

	when := time.Date(2026, 11, 6, 18, 30, 0, 0, time.Local)
	n := 7
	var nilInt *int
	word := "hi"

	assert.Equal(EmptyCell, FormatCell(nil))
	assert.Equal(EmptyCell, FormatCell(""))
	assert.Equal(EmptyCell, FormatCell("   "))
	assert.Equal("hi", FormatCell(word))
	assert.Equal("Yes", FormatCell(true))
	assert.Equal("No", FormatCell(false))
	assert.Equal("2026-11-06 18:30", FormatCell(when))
	assert.Equal(EmptyCell, FormatCell(time.Time{}))
	assert.Equal("7", FormatCell(&n))
	assert.Equal("hi", FormatCell(&word))
	assert.Equal(EmptyCell, FormatCell(nilInt))
	assert.Equal("0", FormatCell(0), "zero numbers are real values, not blanks")
}

func TestStartProgressNoOpWithoutTerminal(t *testing.T) {
	p := StartProgress("x", 3)
	p.Increment()
	p.Stop()

	require.IsType(t, noProgress{}, p)
}

func TestColorSwatchHexWithoutColor(t *testing.T) {
	// Tests don't run on a terminal, so color is off: the hex code stays.
	assert := require.New(t)

	assert.Equal("#00bba7", ColorSwatch("#00bba7"))
	assert.Equal(EmptyCell, ColorSwatch(""))
}

func TestLinkTarget(t *testing.T) {
	assert := require.New(t)

	assert.Equal("https://go.dev/doc", linkTarget("https://go.dev/doc"))
	assert.Equal("http://example.com", linkTarget(" http://example.com "))
	assert.Empty(linkTarget(EmptyCell), "the empty-cell dash isn't a link")
	assert.Empty(linkTarget("go.dev/doc"), "no scheme")
	assert.Empty(linkTarget("javascript:alert(1)"), "only http and https")
	assert.Empty(linkTarget("file:///etc/passwd"))
}

func TestLinkColumnLinksOnlyThatColumn(t *testing.T) {
	assert := require.New(t)

	tbl := NewTable(t.Context(), "Links")
	tbl.Header("ID", "URL")
	tbl.Row(1, "https://go.dev")
	tbl.Row(2, "")
	tbl.LinkColumn(1)

	assert.Equal("https://go.dev", tbl.data.Link(0, 1))
	assert.Empty(tbl.data.Link(0, 0), "other columns aren't linked")
	assert.Empty(tbl.data.Link(1, 1), "a blank URL shows as a dash and isn't linked")
	assert.Empty(tbl.data.Link(5, 1), "out of range")
}
