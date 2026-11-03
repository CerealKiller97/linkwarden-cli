package archives

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestUploadStandalone_CreatesLinkWithArchive(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	url := e2etest.UniqueURL()
	filePath := e2etest.WriteFile(t, "page.html", []byte("<html><body>Saved page</body></html>"))

	cmd := UploadStandalone()
	cmd.SetArgs([]string{"--file", filePath, "--url", url, "--format", "4"})
	assert.NoError(cmd.Execute())

	links, err := client.Dashboard.GetV1(t.Context())
	assert.NoError(err)

	var created *linkwarden.Link
	for i := range links {
		if links[i].URL == url {
			created = &links[i]
		}
	}
	assert.NotNil(created, "a link for the uploaded URL should exist")
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), created.ID) })
}

func TestUploadStandalone_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)
	filePath := e2etest.WriteFile(t, "page.html", []byte("<html></html>"))

	cmd := UploadStandalone()
	cmd.SetArgs([]string{"--file", filePath, "--url", e2etest.UniqueURL(), "--format", "4"})
	assert.Error(cmd.Execute())
}
