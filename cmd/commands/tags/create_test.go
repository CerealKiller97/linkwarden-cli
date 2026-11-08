package tags

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateSendsLabelsAsTags(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		Tags []struct {
			Label string `json:"label"`
		} `json:"tags"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":2,"name":"Reading"}]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{"--label", "Reading", "--label", "Tech"})
	assert.NoError(cmd.Execute())
	assert.Len(gotBody.Tags, 2)
	assert.Equal("Reading", gotBody.Tags[0].Label)
}

func TestCreatePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{"--label", "Reading"})
	assert.Error(cmd.Execute())
}
