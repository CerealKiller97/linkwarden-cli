package utils

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/require"
)

func TestColorEnabledFor(t *testing.T) {
	assert := require.New(t)

	assert.True(colorEnabledFor("", true), "no NO_COLOR, real terminal -> color on")
	assert.False(colorEnabledFor("", false), "no NO_COLOR, not a terminal -> color off")
	assert.False(colorEnabledFor("1", true), "NO_COLOR set -> color off even on a terminal")
	assert.False(colorEnabledFor("1", false), "NO_COLOR set, not a terminal -> color off")
}

func TestTableStyleFor_ColorDisabled(t *testing.T) {
	assert := require.New(t)

	style := tableStyleFor(false)
	assert.Equal(table.StyleRounded.Box, style.Box, "borders stay rounded regardless of color")
	assert.Empty(style.Color.Header, "no header color applied when color is disabled")
}

func TestTableStyleFor_ColorEnabled(t *testing.T) {
	assert := require.New(t)

	style := tableStyleFor(true)
	assert.Equal(table.StyleRounded.Box, style.Box, "borders stay rounded")
	assert.NotEmpty(style.Color.Header, "header gets a color when color is enabled")
}

func TestNewTable_ReturnsConfiguredWriter(t *testing.T) {
	assert := require.New(t)

	tbl := NewTable("Example")
	assert.NotNil(tbl)

	tbl.AppendHeader(table.Row{"A", "B"})
	tbl.AppendRow(table.Row{1, 2})
	assert.Equal(1, tbl.Length())
}
