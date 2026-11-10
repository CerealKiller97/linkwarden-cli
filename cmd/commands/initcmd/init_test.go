package initcmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

// fakeLinkwarden serves the three endpoints init uses. Only sessionToken may
// mint tokens, and only apiToken may call /users/me.
func fakeLinkwarden(t *testing.T, sessionToken, apiToken string) *httptest.Server {
	t.Helper()

	respond := func(w http.ResponseWriter, status int, body any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/session", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Username, Password, SessionName string }
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Username != "alice" || req.Password != "hunter2" {
			respond(w, http.StatusUnauthorized, map[string]string{"response": "Invalid credentials."})
			return
		}
		respond(w, http.StatusOK, map[string]any{"response": map[string]string{"token": sessionToken}})
	})
	mux.HandleFunc("POST /api/v1/tokens", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+sessionToken {
			respond(w, http.StatusUnauthorized, map[string]string{"response": "nope"})
			return
		}
		respond(w, http.StatusOK, map[string]any{"response": map[string]any{"secretKey": apiToken}})
	})
	mux.HandleFunc("GET /api/v1/users/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+apiToken {
			respond(w, http.StatusUnauthorized, map[string]string{"response": "nope"})
			return
		}
		respond(w, http.StatusOK, map[string]any{"response": map[string]any{"id": 1, "username": "alice"}})
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

func TestLoginForTokenMintsAccessTokenFromSession(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")

	token, err := loginForToken(context.Background(), srv.URL, "alice", "hunter2")
	assert.NoError(err)
	assert.Equal("api-token", token)
}

func TestLoginForTokenWrongPassword(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")

	_, err := loginForToken(context.Background(), srv.URL, "alice", "wrong")
	assert.ErrorContains(err, "sign-in failed")
}

func TestVerify(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")

	user, err := verify(context.Background(), srv.URL, "api-token")
	assert.NoError(err)
	assert.Equal("alice", user.Username)

	_, err = verify(context.Background(), srv.URL, "bad-token")
	assert.ErrorContains(err, "rejected")
}

func TestAccessTokensURL(t *testing.T) {
	require.Equal(t, "https://links.example.com/settings/access-tokens",
		accessTokensURL("https://links.example.com"))
}

func TestValidateURL(t *testing.T) {
	assert := require.New(t)

	assert.NoError(validateURL("https://cloud.linkwarden.app"))
	assert.NoError(validateURL(" http://localhost:3000/ "))
	assert.Error(validateURL("links.example.com"))
	assert.Error(validateURL("ftp://links.example.com"))
	assert.Error(validateURL(""))
}

func TestCommandFailsWithoutInteractiveTerminal(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	utils.SetNonInteractive(true)
	t.Cleanup(func() { utils.SetNonInteractive(false) })

	cmd := Command(nil)
	cmd.SetArgs(nil)
	require.ErrorContains(t, cmd.Execute(), "interactive terminal")
}

func labelsOf(methods []signInMethod) []string {
	labels := make([]string, len(methods))
	for i, m := range methods {
		labels[i] = m.Label
	}

	return labels
}

func TestSignInMethods(t *testing.T) {
	assert := require.New(t)
	str := func(s string) *string { return &s }

	// Instance couldn't be asked: token and email/password only.
	assert.Equal([]string{"Access Token (Recommended)", "Email / Password"}, labelsOf(signInMethods(nil)))

	// SSO providers are appended after the built-in methods.
	methods := signInMethods(&linkwarden.LoginConfig{
		CredentialsEnabled: str("true"),
		ButtonAuths: []linkwarden.ButtonAuth{
			{Method: "google", Name: "Google"},
			{Method: "apple", Name: "Apple"},
			{Method: "keycloak"}, // no display name: falls back to the method
			{},                   // nothing to show: skipped
		},
	})
	assert.Equal([]string{
		"Access Token (Recommended)", "Email / Password",
		"Sign in with Google", "Sign in with Apple", "Sign in with keycloak",
	}, labelsOf(methods))
	assert.Equal("Google", methods[2].Provider)
	assert.True(methods[1].Password)

	// Email/password hidden when the instance disabled it.
	methods = signInMethods(&linkwarden.LoginConfig{
		CredentialsEnabled: str("false"),
		ButtonAuths:        []linkwarden.ButtonAuth{{Method: "apple", Name: "Apple"}},
	})
	assert.Equal([]string{"Access Token (Recommended)", "Sign in with Apple"}, labelsOf(methods))
}

func TestTokenSteps(t *testing.T) {
	assert := require.New(t)
	page := "https://links.example.com/settings/access-tokens"

	plain := tokenSteps(page, "", false)
	assert.Len(plain, 5) // heading + 4 steps
	assert.Contains(plain[1], page)

	sso := tokenSteps(page, "Apple", false)
	assert.Len(sso, 6) // extra sign-in step
	assert.Contains(sso[0], "Sign in with Apple")
	assert.Contains(sso[2], "Sign in with Apple")
}

func TestBanner(t *testing.T) {
	assert := require.New(t)

	plain := banner(false)
	assert.NotContains(plain, "\x1b[", "no ANSI codes without color")
	assert.Contains(plain, "Welcome to Linkwarden CLI")
	assert.Len(strings.Split(plain, "\n"), len(logo))

	assert.Contains(banner(true), "\x1b[", "styled with color on")
}
