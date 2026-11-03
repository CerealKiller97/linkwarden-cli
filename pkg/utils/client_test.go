package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// This package can't import internal/e2etest (which imports it), so it reads the test
// instance from the environment itself.
func TestNewLinkwardenClient_UsesConfiguredURLAndToken(t *testing.T) {
	assert := require.New(t)

	baseURL, token := os.Getenv("LINKWARDEN_URL"), os.Getenv("LINKWARDEN_TOKEN")
	if baseURL == "" || token == "" {
		if os.Getenv("CI") != "" {
			t.Fatal("LINKWARDEN_URL and LINKWARDEN_TOKEN must be set in CI")
		}
		t.Skip("LINKWARDEN_URL / LINKWARDEN_TOKEN not set; start an instance with scripts/dev-linkwarden.sh up")
	}

	SetVars(token, baseURL)
	_, err := NewLinkwardenClient().Collections.List(t.Context())
	assert.NoError(err)

	SetVars("invalid-token", baseURL)
	_, err = NewLinkwardenClient().Collections.List(t.Context())
	assert.Error(err)
}
