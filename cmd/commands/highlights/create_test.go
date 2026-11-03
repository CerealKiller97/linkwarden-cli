package highlights

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestCreate_UsesLinkIDFlag(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.CreateHighlightRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/highlights", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"linkId":9,"color":"yellow","text":"quoted","startOffset":0,"endOffset":7}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{
		"--link-id", "9",
		"--text", "quoted",
		"--color", "yellow",
		"--start-offset", "0",
		"--end-offset", "7",
	})
	assert.NoError(cmd.Execute())
	assert.Equal(9, gotBody.LinkID)
	assert.Equal("quoted", gotBody.Text)
}

func TestCreate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid highlight"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{
		"--link-id", "9",
		"--text", "quoted",
		"--color", "yellow",
		"--start-offset", "0",
		"--end-offset", "7",
	})
	assert.Error(cmd.Execute())
}
