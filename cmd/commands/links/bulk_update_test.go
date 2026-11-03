package links

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestBulkUpdate_AppliesFlags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	from := e2etest.Collection(t, client)
	to := e2etest.Collection(t, client)
	first := e2etest.Link(t, client, from.ID)
	second := e2etest.Link(t, client, from.ID)
	tagA, tagB := e2etest.UniqueName("a"), e2etest.UniqueName("b")
	e2etest.CleanupTagsNamed(t, client, tagA, tagB)

	cmd := BulkUpdate()
	cmd.SetArgs([]string{
		"--ids", fmt.Sprintf("%d,%d", first.ID, second.ID),
		"--collection-id", strconv.Itoa(to.ID),
		"--tag", tagA,
		"--tag", tagB,
		"--remove-previous-tags",
	})
	assert.NoError(cmd.Execute())

	moved, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: to.ID})
	assert.NoError(err)
	assert.ElementsMatch([]int{first.ID, second.ID}, e2etest.LinkIDs(moved))
	for _, l := range moved {
		var names []string
		for _, tag := range l.Tags {
			names = append(names, tag.Name)
		}
		assert.ElementsMatch([]string{tagA, tagB}, names)
	}
}

func TestBulkUpdate_InvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := BulkUpdate()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
