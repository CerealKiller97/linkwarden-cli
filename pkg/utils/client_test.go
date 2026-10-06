package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLinkwardenClient_UsesConfiguredURLAndToken(t *testing.T) {
	assert := require.New(t)

	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[]}`))
	}))
	defer server.Close()

	SetVars("test-token", server.URL)

	client := NewLinkwardenClient()
	_, err := client.Collections.List(t.Context())
	assert.NoError(err)
	assert.Equal("Bearer test-token", gotAuth)
}
