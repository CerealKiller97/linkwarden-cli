package initcmd

import (
	"errors"
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils/prompttest"
	"github.com/stretchr/testify/require"
)

// isolateInit gives init an empty config dir and records the URLs it would
// open in a browser instead of opening them. browserErr, if set, is what
// opening returns.
func isolateInit(t *testing.T, browserErr error) *[]string {
	t.Helper()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")

	var opened []string
	prev := openBrowser
	openBrowser = func(u string) error {
		opened = append(opened, u)
		return browserErr
	}
	t.Cleanup(func() { openBrowser = prev })

	return &opened
}

func runInit(t *testing.T) error {
	t.Helper()

	cmd := Command(nil)
	cmd.SetArgs(nil)

	return cmd.Execute()
}

func savedConfig(t *testing.T) utils.Config {
	t.Helper()

	cfg, err := utils.LoadConfig()
	require.NoError(t, err)

	return cfg
}

func TestRunTokenFlowSavesVerifiedToken(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	opened := isolateInit(t, nil)

	prompttest.Use(t,
		prompttest.Answer{Label: "Linkwarden URL", Value: srv.URL + "/"},
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Access Token (Recommended)"},
		prompttest.Answer{Label: "Access token", Value: "  api-token \n"},
	)

	assert.NoError(runInit(t))

	assert.Equal(utils.Config{URL: srv.URL, Token: "api-token"}, savedConfig(t),
		"URL without its trailing slash, token trimmed")
	assert.Equal([]string{srv.URL + "/settings/access-tokens"}, *opened)

	path, err := utils.ConfigFile()
	assert.NoError(err)
	info, err := os.Stat(path)
	assert.NoError(err)
	assert.Equal(os.FileMode(0o600), info.Mode().Perm(), "the config file holds a secret")
}

func TestRunPasswordFlowSavesMintedTokenNotPassword(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	opened := isolateInit(t, nil)

	prompttest.Use(t,
		prompttest.Answer{Label: "Linkwarden URL", Value: srv.URL},
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Email / Password"},
		prompttest.Answer{Label: "Email or username", Value: " alice "},
		prompttest.Answer{Label: "Password", Value: "hunter2"},
	)

	assert.NoError(runInit(t))

	assert.Equal(utils.Config{URL: srv.URL, Token: "api-token"}, savedConfig(t))
	assert.Empty(*opened, "the password flow needs no browser")

	path, err := utils.ConfigFile()
	assert.NoError(err)
	data, err := os.ReadFile(path)
	assert.NoError(err)
	assert.NotContains(string(data), "hunter2", "the password is never stored")
	assert.NotContains(string(data), "session-jwt", "neither is the sign-in session")
}

func TestRunRejectedTokenSavesNothing(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	isolateInit(t, nil)

	prompttest.Use(t,
		prompttest.Answer{Label: "Linkwarden URL", Value: srv.URL},
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Access Token (Recommended)"},
		prompttest.Answer{Label: "Access token", Value: "not-the-token"},
	)

	err := runInit(t)
	assert.ErrorContains(err, "rejected")

	path, err := utils.ConfigFile()
	assert.NoError(err)
	assert.NoFileExists(path)
}

func TestRunWrongPasswordSavesNothing(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	isolateInit(t, nil)

	prompttest.Use(t,
		prompttest.Answer{Label: "Linkwarden URL", Value: srv.URL},
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Email / Password"},
		prompttest.Answer{Label: "Email or username", Value: "alice"},
		prompttest.Answer{Label: "Password", Value: "wrong"},
	)

	assert.ErrorContains(runInit(t), "sign-in failed")

	path, err := utils.ConfigFile()
	assert.NoError(err)
	assert.NoFileExists(path)
}

func TestRunKeepsExistingConfigUnlessOverwriteConfirmed(t *testing.T) {
	assert := require.New(t)
	isolateInit(t, nil)

	existing := utils.Config{URL: "https://old.example.com", Token: "old-token"}
	_, err := utils.SaveConfig(existing)
	assert.NoError(err)

	// Answering No must end init: any further prompt would fail the test.
	prompttest.Use(t, prompttest.Answer{Label: "already set up for https://old.example.com", Value: "No"})

	assert.NoError(runInit(t))
	assert.Equal(existing, savedConfig(t))
}

func TestRunOverwriteOffersSavedURLAsDefault(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	isolateInit(t, nil)

	_, err := utils.SaveConfig(utils.Config{URL: srv.URL, Token: "old-token"})
	assert.NoError(err)

	prompttest.Use(t,
		prompttest.Answer{Label: "Overwrite it?", Value: "Yes"},
		prompttest.Answer{Label: "Linkwarden URL"}, // Enter: keep the saved URL
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Access Token (Recommended)"},
		prompttest.Answer{Label: "Access token", Value: "api-token"},
	)

	assert.NoError(runInit(t))
	assert.Equal(utils.Config{URL: srv.URL, Token: "api-token"}, savedConfig(t))
}

func TestRunContinuesWhenBrowserCannotOpen(t *testing.T) {
	assert := require.New(t)
	srv := fakeLinkwarden(t, "session-jwt", "api-token")
	isolateInit(t, errors.New("no display"))

	prompttest.Use(t,
		prompttest.Answer{Label: "Linkwarden URL", Value: srv.URL},
		prompttest.Answer{Label: "How would you like to set up your account?", Value: "Access Token (Recommended)"},
		prompttest.Answer{Label: "Access token", Value: "api-token"},
	)

	assert.NoError(runInit(t))
	assert.Equal("api-token", savedConfig(t).Token)
}
