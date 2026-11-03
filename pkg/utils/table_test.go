package utils

import (
	"bytes"
	"strings"
	"testing"

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

	assert.False(shouldBrowse(5, 40), "fits on screen -> just print")
	assert.True(shouldBrowse(100, 40), "taller than the screen -> viewer")
	assert.False(shouldBrowse(100, 0), "unknown height -> just print")
}

func TestTable_RendersStaticWhenNotATerminal(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	tbl := NewTable("Example")
	tbl.out = &out
	tbl.Header("ID", "Name")
	tbl.Row(1, "first")
	tbl.Row(2, nil)

	tbl.Render()

	assert.Equal(2, tbl.Len())
	got := out.String()
	assert.True(strings.HasPrefix(got, "Example\n"))
	assert.Contains(got, "╭")
	assert.Contains(got, "first")
	assert.Contains(got, "<nil>")
	assert.NotContains(got, "\x1b[", "test output isn't a terminal, so no ANSI")
}

func TestStartProgress_NoOpWithoutTerminal(t *testing.T) {
	p := StartProgress("x", 3)
	p.Increment()
	p.Stop()

	require.IsType(t, noProgress{}, p)
}
