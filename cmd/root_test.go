package cmd

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestResolveCredentials(t *testing.T) {
	saved := utils.Config{URL: "https://saved.example.com", Token: "saved-token"}

	cases := map[string]struct {
		envToken, envURL string
		cfg              utils.Config
		wantToken        string
		wantURL          string
		wantErr          bool
	}{
		"env token, no env URL uses default": {
			envToken: "env-token", cfg: saved, wantToken: "env-token", wantURL: utils.DefaultInstanceURL,
		},
		"env token and env URL": {
			envToken: "env-token", envURL: "https://env.example.com/", cfg: saved,
			wantToken: "env-token", wantURL: "https://env.example.com",
		},
		"config only": {
			cfg: saved, wantToken: "saved-token", wantURL: "https://saved.example.com",
		},
		"config token, env URL overrides": {
			envURL: "http://localhost:3000", cfg: saved, wantToken: "saved-token", wantURL: "http://localhost:3000",
		},
		"config token without URL uses default": {
			cfg: utils.Config{Token: "saved-token"}, wantToken: "saved-token", wantURL: utils.DefaultInstanceURL,
		},
		"nothing configured": {wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)

			token, url, err := resolveCredentials(tc.envToken, tc.envURL, tc.cfg)
			if tc.wantErr {
				assert.ErrorIs(err, errNotLoggedIn)
				return
			}

			assert.NoError(err)
			assert.Equal(tc.wantToken, token)
			assert.Equal(tc.wantURL, url)
		})
	}
}

func TestInitSkipsCredentialCheck(t *testing.T) {
	assert := require.New(t)

	t.Setenv("LINKWARDEN_TOKEN", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	root := Root()
	root.SetArgs([]string{"init", "--non-interactive"})

	// With no credentials anywhere, any other command fails with errNotLoggedIn;
	// init must get past PersistentPreRunE and fail on its first prompt instead.
	err := root.Execute()
	assert.Error(err)
	assert.NotErrorIs(err, errNotLoggedIn)
	assert.ErrorContains(err, "interactive terminal")
}

func TestOtherCommandsRequireCredentials(t *testing.T) {
	t.Setenv("LINKWARDEN_TOKEN", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	root := Root()
	root.SetArgs([]string{"tags", "list"})

	require.ErrorIs(t, root.Execute(), errNotLoggedIn)
}
