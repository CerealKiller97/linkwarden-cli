package collections

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestUpdate_AppliesFlags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	renamed := e2etest.UniqueName("renamed")

	cmd := Update()
	cmd.SetArgs([]string{"--id", strconv.Itoa(collection.ID), "--name", renamed, "--public"})
	assert.NoError(cmd.Execute())

	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(renamed, got.Name)
	assert.True(got.IsPublic)
}

func TestUpdate_KeepsUnsetFields(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)

	cmd := Update()
	cmd.SetArgs([]string{"--id", strconv.Itoa(collection.ID), "--color", "#10b981"})
	assert.NoError(cmd.Execute())

	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(collection.Name, got.Name)
	assert.Equal("#10b981", got.Color)
}

func TestUpdate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "999999999", "--name", "Renamed"})
	assert.Error(cmd.Execute())
}
