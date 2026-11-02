package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkerService_Stats(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/worker", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"","data":{"link":{"pending":3,"done":10,"failed":1},"search":{"pending":2,"done":11}}}`))
	})

	stats, err := client.Worker.Stats(t.Context())
	assert.NoError(err)
	assert.Equal(3, stats.Link.Pending)
	assert.Equal(1, stats.Link.Failed)
	assert.Equal(11, stats.Search.Done)
}

func TestWorkerService_ResetPreservation(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/worker/preservation", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Started."}`))
	})

	msg, err := client.Worker.ResetPreservation(t.Context(), PreservationResetBroken)
	assert.NoError(err)
	assert.Equal("Started.", msg)
	assert.Equal("allBroken", gotBody["action"])
}
