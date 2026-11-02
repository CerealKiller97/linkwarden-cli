package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	var gotAuth string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/session", r.URL.Path)
		gotAuth = r.Header.Get("Authorization")
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"token":"session-abc"}}`))
	})
	client.token = ""

	token, err := client.Session.Create(t.Context(), CreateSessionRequest{Username: "dev", Password: "devpassword"})
	assert.NoError(err)
	assert.Equal("session-abc", token)
	assert.Empty(gotAuth)
	assert.Equal(map[string]any{"username": "dev", "password": "devpassword"}, gotBody)
}

func TestSessionService_Create_EmailNotVerified(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"response":"Please verify your email address before logging in.","code":"EMAIL_NOT_VERIFIED"}`))
	})

	_, err := client.Session.Create(t.Context(), CreateSessionRequest{Username: "dev", Password: "x"})
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal("Please verify your email address before logging in.", apiErr.Message)
}

func TestAuthService(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		var body map[string]any
		switch r.URL.Path {
		case "/api/v1/auth/forgot-password":
			assert.NoError(json.NewDecoder(r.Body).Decode(&body))
			assert.Equal("a@example.com", body["email"])
		case "/api/v1/auth/reset-password":
			assert.NoError(json.NewDecoder(r.Body).Decode(&body))
			assert.Equal("tok", body["token"])
			assert.Equal("newpassword", body["password"])
		case "/api/v1/auth/verify-email":
			assert.Equal("tok", r.URL.Query().Get("token"))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"ok"}`))
	})

	msg, err := client.Auth.ForgotPassword(t.Context(), "a@example.com")
	assert.NoError(err)
	assert.Equal("ok", msg)

	_, err = client.Auth.ResetPassword(t.Context(), "tok", "newpassword")
	assert.NoError(err)

	_, err = client.Auth.VerifyEmail(t.Context(), "tok")
	assert.NoError(err)
}

func TestLoginsService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/logins", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"credentialsEnabled":"true","emailEnabled":null,"registrationDisabled":"false","buttonAuths":[{"method":"github","name":"GitHub"}]}`))
	})

	cfg, err := client.Logins.Get(t.Context())
	assert.NoError(err)
	assert.Equal("true", *cfg.CredentialsEnabled)
	assert.Nil(cfg.EmailEnabled)
	assert.Equal("GitHub", cfg.ButtonAuths[0].Name)
}
