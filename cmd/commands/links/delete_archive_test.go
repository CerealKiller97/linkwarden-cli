package links

import (
	"fmt"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDeleteArchive_ParsesIDs(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	first := e2etest.Link(t, client, collection.ID)
	second := e2etest.Link(t, client, collection.ID)

	cmd := DeleteArchive()
	cmd.SetArgs([]string{"--ids", fmt.Sprintf("%d, %d", first.ID, second.ID)})
	assert.NoError(cmd.Execute())
}

func TestDeleteArchive_InvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := DeleteArchive()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
