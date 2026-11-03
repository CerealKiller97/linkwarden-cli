package migrations

import (
	"fmt"
	"testing"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestImport_ImportsBookmarksFile(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	folder := e2etest.UniqueName("imported")
	url := e2etest.UniqueURL()
	filePath := e2etest.WriteFile(t, "bookmarks.html", fmt.Appendf(nil, `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>%s</H3>
    <DL><p>
        <DT><A HREF="%s">Imported link</A>
    </DL><p>
</DL><p>
`, folder, url))

	cmd := Import()
	cmd.SetArgs([]string{"--file", filePath, "--format", "1"})
	assert.NoError(cmd.Execute())

	// The import runs in the background; wait for its folder to show up as a collection.
	var imported *linkwarden.Collection
	assert.Eventually(func() bool {
		collections, err := client.Collections.List(t.Context())
		if err != nil {
			return false
		}
		for i := range collections {
			if collections[i].Name == folder {
				imported = &collections[i]
				return true
			}
		}

		return false
	}, time.Minute, time.Second)
	t.Cleanup(func() { _ = client.Collections.Delete(t.Context(), imported.ID) })

	links, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: imported.ID})
	assert.NoError(err)
	assert.Len(links, 1)
	assert.Equal(url, links[0].URL)
}

func TestImport_MissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Import()
	cmd.SetArgs([]string{"--file", "/nonexistent/file.json", "--format", "0"})
	assert.Error(cmd.Execute())
}
