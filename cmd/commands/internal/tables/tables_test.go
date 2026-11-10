package tables

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestLinksOneRowPerLink(t *testing.T) {
	links := []linkwarden.Link{
		{ID: 1, Name: "Bubble Tea", URL: "https://github.com/charmbracelet/bubbletea", Tags: []linkwarden.Tag{{Name: "go"}, {Name: "tui"}}},
		{ID: 2, Name: "Lip Gloss"},
	}

	require.Equal(t, 2, Links(nil, "Links", links).Len())
}

func TestTagNames(t *testing.T) {
	assert := require.New(t)

	assert.Equal("go, tui", tagNames([]linkwarden.Tag{{Name: "go"}, {Name: "tui"}}))
	assert.Empty(tagNames(nil), "no tags is blank, which the table shows as a dash")
}

func TestHighlightsOneRowPerHighlight(t *testing.T) {
	require.Equal(t, 1, Highlights(nil, "Highlights", []linkwarden.Highlight{{ID: 1, LinkID: 9, Text: "quote"}}).Len())
}
