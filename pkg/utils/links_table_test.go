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
	tbl := NewLinksTable("Links", []linkwarden.Link{{ID: 1, Name: "a", URL: "https://a"}})
	tbl.out = &withoutTags
	tbl.renderStatic(120, false)
	assert.NotContains(withoutTags.String(), "Tags")

	var withTags bytes.Buffer
	tbl = NewLinksTable("Links", []linkwarden.Link{
		{ID: 1, Name: "a", URL: "https://a"},
		{ID: 2, Name: "b", URL: "https://b", Tags: []linkwarden.Tag{{Name: "go"}, {Name: "cli"}}},
	})
	tbl.out = &withTags
	tbl.renderStatic(120, false)
	assert.Contains(withTags.String(), "Tags")
	assert.Contains(withTags.String(), "go, cli")
}
