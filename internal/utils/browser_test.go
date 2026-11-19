package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowserURL(t *testing.T) {
	assert := require.New(t)

	got, err := browserURL("https://links.example.com/settings/access-tokens")
	assert.NoError(err)
	assert.Equal("https://links.example.com/settings/access-tokens", got)

	_, err = browserURL("http://localhost:3000")
	assert.NoError(err)

	for _, bad := range []string{"", "-a Calculator", "file:///etc/passwd", "javascript:alert(1)", "https://", "links.example.com"} {
		_, err := browserURL(bad)
		assert.Error(err, "input %q", bad)
	}
}
