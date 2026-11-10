package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
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

func TestExitCode(t *testing.T) {
	cases := map[string]struct {
		err  error
		want int
	}{
		"success":         {err: nil, want: 0},
		"failure":         {err: errors.New("boom"), want: 1},
		"aborted":         {err: tui.ErrAborted, want: exitAborted},
		"wrapped aborted": {err: fmt.Errorf("resolve collection ID: %w", tui.ErrAborted), want: exitAborted},
		"interrupted":     {err: fmt.Errorf("list links: %w", context.Canceled), want: exitAborted},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.want, exitCode(tc.err))
		})
	}
}

// Every command needs a Short, or --help lists it with a blank description.
func TestEveryCommandHasShort(t *testing.T) {
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if !c.Hidden && c.Name() != "help" && c.Name() != "completion" {
			require.NotEmpty(t, c.Short, "%q has no Short", c.CommandPath())
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}

	walk(Root())
}

// Every command you can run needs an Example: they show up in --help and docs/cli/.
func TestEveryCommandHasExample(t *testing.T) {
	for _, c := range ownCommands(Root()) {
		if c.Runnable() {
			require.NotEmpty(t, c.Example, "%q has no Example", c.CommandPath())
		}
	}
}

// Every "linkwarden …" line in an Example must name a real command and only use
// flags it has, with values they accept, so examples can't drift from the flags.
func TestExamplesParse(t *testing.T) {
	for _, c := range ownCommands(Root()) {
		for _, line := range exampleCommands(c.Example) {
			args, err := splitArgs(line)
			require.NoError(t, err, "%q: %s", c.CommandPath(), line)

			root := Root()
			found, rest, err := root.Find(args[1:])
			require.NoError(t, err, "%q: %s", c.CommandPath(), line)
			require.NotEqual(t, root, found, "%q: %s names no subcommand", c.CommandPath(), line)
			require.NoError(t, found.ParseFlags(rest), "%q: %s", c.CommandPath(), line)
		}
	}
}

// ownCommands lists every command under root except Cobra's help and completion.
func ownCommands(root *cobra.Command) []*cobra.Command {
	var all []*cobra.Command
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Hidden || c.Name() == "help" || c.Name() == "completion" {
			return
		}
		all = append(all, c)
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)

	return all
}

// exampleCommands returns the lines of an Example that run linkwarden, with
// backslash continuations joined and anything after a pipe dropped.
func exampleCommands(example string) []string {
	var cmds []string
	var cur string
	for _, line := range strings.Split(example, "\n") {
		line = strings.TrimSpace(line)
		if cur != "" {
			line = cur + " " + line
			cur = ""
		}
		if strings.HasSuffix(line, `\`) {
			cur = strings.TrimSuffix(line, `\`)
			continue
		}
		if !strings.HasPrefix(line, "linkwarden ") {
			continue
		}
		line, _, _ = strings.Cut(line, " | ")
		cmds = append(cmds, line)
	}

	return cmds
}

// splitArgs splits a command line into words like a shell would for the simple
// quoting examples use: single and double quotes, no escapes or expansion.
func splitArgs(line string) ([]string, error) {
	var args []string
	var word strings.Builder
	inWord := false
	var quote rune
	for _, r := range line {
		switch {
		case quote != 0 && r == quote:
			quote = 0
		case quote != 0:
			word.WriteRune(r)
		case r == '"' || r == '\'':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			if inWord {
				args = append(args, word.String())
				word.Reset()
				inWord = false
			}
		default:
			word.WriteRune(r)
			inWord = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unclosed %c quote", quote)
	}
	if inWord {
		args = append(args, word.String())
	}

	return args, nil
}
