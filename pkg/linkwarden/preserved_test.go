package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreservedService_Token(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/preserved/token", r.URL.Path)
		assert.Equal("7", r.URL.Query().Get("linkId"))
		assert.Equal("4", r.URL.Query().Get("format"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"url":"https://usercontent.example.com/api/v1/preserved/view?token=abc"}}`))
	})

	signed, err := client.Preserved.Token(t.Context(), 7, ArchiveFormatMonolithHTML)
	assert.NoError(err)
	assert.Equal("https://usercontent.example.com/api/v1/preserved/view?token=abc", signed.URL)
}

func TestPreservedService_View(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/preserved/view", r.URL.Path)
		assert.Equal("abc", r.URL.Query().Get("token"))
		assert.Equal("1", r.URL.Query().Get("download"))
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="7.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF-1.7"))
	})

	file, err := client.Preserved.View(t.Context(), "abc", true)
	assert.NoError(err)
	assert.Equal("application/pdf", file.ContentType)
	assert.Equal("7.pdf", file.Filename)
	assert.Equal([]byte("%PDF-1.7"), file.Data)
}

func TestAvatarService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/avatar/2", r.URL.Path)
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{0xff, 0xd8})
	})

	avatar, err := client.Avatar.Get(t.Context(), 2)
	assert.NoError(err)
	assert.Equal("image/jpeg", avatar.ContentType)
	assert.Equal([]byte{0xff, 0xd8}, avatar.Data)
}
