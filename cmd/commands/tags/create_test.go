package tags

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestCreate_CreatesLabelsAsTags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	reading, tech := e2etest.UniqueName("reading"), e2etest.UniqueName("tech")
	e2etest.CleanupTagsNamed(t, client, reading, tech)

	cmd := Create()
	cmd.SetArgs([]string{"--label", reading, "--label", tech})
	assert.NoError(cmd.Execute())

	for _, name := range []string{reading, tech} {
		tags, err := client.Tags.List(t.Context(), linkwarden.ListTagsOptions{Search: name})
		assert.NoError(err)
		assert.Len(tags, 1)
		assert.Equal(name, tags[0].Name)
	}
}

func TestCreate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Create()
	cmd.SetArgs([]string{"--label", "Reading"})
	assert.Error(cmd.Execute())
}
