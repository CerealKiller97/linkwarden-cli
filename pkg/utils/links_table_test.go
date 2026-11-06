package utils

import (
	"bytes"
	"testing"

	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestNewLinksTable_TagsColumn(t *testing.T) {
	assert := require.New(t)

	var withoutTags bytes.Buffer
	NewLinksTable("Links", []linkwarden.Link{{ID: 1, Name: "a", URL: "https://a"}}).RenderTo(&withoutTags, 120, false)
	assert.NotContains(withoutTags.String(), "TAGS")

	var withTags bytes.Buffer
	NewLinksTable("Links", []linkwarden.Link{
		{ID: 1, Name: "a", URL: "https://a"},
		{ID: 2, Name: "b", URL: "https://b", Tags: []linkwarden.Tag{{Name: "go"}, {Name: "cli"}}},
	}).RenderTo(&withTags, 120, false)
	assert.Contains(withTags.String(), "TAGS")
	assert.Contains(withTags.String(), "go, cli")
}
