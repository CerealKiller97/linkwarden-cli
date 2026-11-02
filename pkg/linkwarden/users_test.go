package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsersService_Me(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/users/me", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"name":null,"username":"dev","email":null,"image":"","locale":"en","dashboardSections":[],"subscription":{},"hasUnIndexedLinks":true}}`))
	})

	user, err := client.Users.Me(t.Context())
	assert.NoError(err)
	assert.Equal(1, user.ID)
	assert.Equal("dev", user.Username)
	assert.Nil(user.Name)
	assert.True(user.HasUnIndexedLinks)
}

func TestUsersService_List(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/users", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"username":"dev","email":null,"emailVerified":null,"createdAt":"2024-12-01T22:04:16.492Z"}]}`))
	})

	users, err := client.Users.List(t.Context())
	assert.NoError(err)
	assert.Len(users, 1)
	assert.Equal("dev", users[0].Username)
}

func TestUsersService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/users", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":2,"username":"new"}}`))
	})

	user, err := client.Users.Create(t.Context(), CreateUserRequest{Username: "new", Password: "password123"})
	assert.NoError(err)
	assert.Equal(2, user.ID)
	assert.Equal(map[string]any{"username": "new", "password": "password123"}, gotBody)
}

func TestUsersService_Update_OmitsNilFields(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/users/1", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"locale":"de"}}`))
	})

	locale, private := "de", false
	user, err := client.Users.Update(t.Context(), 1, UpdateUserRequest{Locale: &locale, IsPrivate: &private})
	assert.NoError(err)
	assert.Equal("de", user.Locale)
	assert.Equal(map[string]any{"locale": "de", "isPrivate": false}, gotBody)
}

func TestUsersService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/users/4", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"User account and all related data deleted successfully."}`))
	})

	assert.NoError(client.Users.Delete(t.Context(), 4))
}

func TestUsersService_UpdatePreferences(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/users/1/preference", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"theme":"dark","readableFontSize":null}}`))
	})

	prefs, err := client.Users.UpdatePreferences(t.Context(), 1, UserPreferences{Theme: "dark"})
	assert.NoError(err)
	assert.Equal("dark", prefs.Theme)
	assert.Equal(map[string]any{"theme": "dark"}, gotBody)
}
