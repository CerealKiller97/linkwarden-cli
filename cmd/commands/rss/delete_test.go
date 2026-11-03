package rss

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDelete_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	sub := e2etest.Subscription(t, client, e2etest.Collection(t, client).ID)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", strconv.Itoa(sub.ID)})
	assert.NoError(cmd.Execute())

	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)
	for _, s := range subs {
		assert.NotEqual(sub.ID, s.ID, "deleted subscription should not be listed")
	}
}

func TestDelete_NotFoundErrors(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "999999999"})
	assert.Error(cmd.Execute())
}

func TestDelete_InvalidIDErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "not-a-number"})
	assert.Error(cmd.Execute())
}
