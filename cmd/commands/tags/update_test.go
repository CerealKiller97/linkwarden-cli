package tags

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestUpdate_UsesIDAndNameFlags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	tag := e2etest.Tag(t, client)
	renamed := e2etest.UniqueName("renamed")

	cmd := Update()
	cmd.SetArgs([]string{"--id", strconv.Itoa(tag.ID), "--name", renamed})
	assert.NoError(cmd.Execute())

	got, err := client.Tags.Get(t.Context(), tag.ID)
	assert.NoError(err)
	assert.Equal(renamed, got.Name)
}

func TestUpdate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "999999999", "--name", "Renamed"})
	assert.Error(cmd.Execute())
}
