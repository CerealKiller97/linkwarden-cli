package collections

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestDelete_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", strconv.Itoa(collection.ID)})
	assert.NoError(cmd.Execute())

	_, err := client.Collections.Get(t.Context(), collection.ID)
	assert.True(linkwarden.IsNotFound(err))
}

func TestDelete_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "999999999"})
	assert.Error(cmd.Execute())
}
