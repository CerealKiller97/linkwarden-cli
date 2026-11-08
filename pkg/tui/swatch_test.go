package tui

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/require"
)

func TestSwatch(t *testing.T) {
	assert := require.New(t)

	dot, ok := Swatch("#00bba7")
	assert.True(ok)
	assert.Contains(dot, "●")
	assert.Contains(dot, "\x1b[", "the dot is drawn in the color")
	assert.Equal(1, lipgloss.Width(dot))

	_, ok = Swatch("#0ea")
	assert.True(ok, "short #rgb form works too")

	for _, bad := range []string{"", "00bba7", "#00bb", "red", "#gggggg"} {
		_, ok := Swatch(bad)
		assert.False(ok, bad)
	}
}
