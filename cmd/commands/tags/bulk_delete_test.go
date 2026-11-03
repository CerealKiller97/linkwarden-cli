package tags

import (
	"fmt"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestBulkDelete_ParsesIDs(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	first, second, kept := e2etest.Tag(t, client), e2etest.Tag(t, client), e2etest.Tag(t, client)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", fmt.Sprintf("%d,%d", first.ID, second.ID)})
	assert.NoError(cmd.Execute())

	for _, id := range []int{first.ID, second.ID} {
		_, err := client.Tags.Get(t.Context(), id)
		assert.True(linkwarden.IsNotFound(err))
	}
	_, err := client.Tags.Get(t.Context(), kept.ID)
	assert.NoError(err)
}

func TestBulkDelete_InvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
