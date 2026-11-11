package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/stretchr/testify/require"
)

// answerSelect is a Prompter whose Select picks the option named answer (or
// returns err). pkg/utils/prompttest can't be used here: it imports this package.
type answerSelect struct {
	Prompter
	answer string
	err    error
	asked  *[]string
}

func (a answerSelect) Select(label string, options []string) (int, error) {
	*a.asked = append(*a.asked, label)
	if a.err != nil {
		return 0, a.err
	}
	for i, o := range options {
		if o == a.answer {
			return i, nil
		}
	}

	return 0, errors.New("no such option: " + a.answer)
}

func useSelectAnswer(t *testing.T, answer string, err error) *[]string {
	t.Helper()

	var asked []string
	t.Cleanup(SetPrompter(answerSelect{answer: answer, err: err, asked: &asked}))

	return &asked
}

func TestConfirmAsksAndHonorsTheAnswer(t *testing.T) {
	for answer, want := range map[string]bool{"Yes": true, "No": false} {
		t.Run(answer, func(t *testing.T) {
			assert := require.New(t)
			asked := useSelectAnswer(t, answer, nil)

			ok, err := Confirm(newConfirmCommand(t), "Delete link 3?")

			assert.NoError(err)
			assert.Equal(want, ok)
			assert.Equal([]string{"Delete link 3?"}, *asked)
		})
	}
}

func TestConfirmAbortedPromptIsAnError(t *testing.T) {
	assert := require.New(t)
	useSelectAnswer(t, "", tui.ErrAborted)

	ok, err := Confirm(newConfirmCommand(t), "Delete?")

	assert.ErrorIs(err, tui.ErrAborted, "Esc must not read as Yes or a silent No")
	assert.False(ok)
}

func TestConfirmWithQuestionErrorAsksNothing(t *testing.T) {
	assert := require.New(t)
	asked := useSelectAnswer(t, "Yes", nil)

	boom := errors.New("fetching collection")
	ok, err := ConfirmWith(newConfirmCommand(t), func() (string, error) { return "", boom })

	assert.ErrorIs(err, boom)
	assert.False(ok)
	assert.Empty(*asked)
}

func TestNonInteractiveFlagWinsOverAvailablePrompter(t *testing.T) {
	assert := require.New(t)
	asked := useSelectAnswer(t, "Yes", nil)

	SetNonInteractive(true)
	t.Cleanup(func() { SetNonInteractive(false) })

	ok, err := Confirm(newConfirmCommand(t), "Delete?")
	assert.ErrorIs(err, ErrConfirmationRequired)
	assert.False(ok)

	_, err = Password("Access token:", nil)
	assert.ErrorContains(err, "interactive terminal")
	assert.Empty(*asked)
}

func TestSaveConfigTightensExistingLoosePermissions(t *testing.T) {
	assert := require.New(t)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	path, err := ConfigFile()
	assert.NoError(err)
	assert.NoError(os.MkdirAll(filepath.Dir(path), 0o755))
	assert.NoError(os.WriteFile(path, []byte(`{"url":"https://old.example.com","token":"a-much-longer-old-token"}`), 0o644))

	_, err = SaveConfig(Config{URL: "https://new.example.com", Token: "new"})
	assert.NoError(err)

	info, err := os.Stat(path)
	assert.NoError(err)
	assert.Equal(os.FileMode(0o600), info.Mode().Perm(), "a token must not stay world-readable")

	got, err := LoadConfig()
	assert.NoError(err, "the shorter new config replaces the old one entirely")
	assert.Equal(Config{URL: "https://new.example.com", Token: "new"}, got)
}

func TestSaveConfigCreatesOwnerOnlyDirectory(t *testing.T) {
	assert := require.New(t)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	path, err := SaveConfig(Config{Token: "secret"})
	assert.NoError(err)

	info, err := os.Stat(filepath.Dir(path))
	assert.NoError(err)
	assert.Equal(os.FileMode(0o700), info.Mode().Perm())
}

func TestAppendTokenToShellFileCreatesOwnerOnlyFile(t *testing.T) {
	assert := require.New(t)
	t.Setenv("SHELL", "/bin/zsh")
	setHome(t, t.TempDir())

	assert.NoError(AppendTokenToShellFile("secret"))

	info, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".zshenv"))
	assert.NoError(err)
	assert.Equal(os.FileMode(0o600), info.Mode().Perm())
}

func TestAppendTokenToShellFileKeepsEarlierLines(t *testing.T) {
	assert := require.New(t)
	t.Setenv("SHELL", "/bin/bash")
	setHome(t, t.TempDir())

	assert.NoError(AppendTokenToShellFile("first"))
	assert.NoError(AppendTokenToShellFile("second"))

	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".bashrc"))
	assert.NoError(err)
	assert.Equal("\n# Managed by Linkwarden CLI\nexport LINKWARDEN_TOKEN=\"first\"\n"+
		"\n# Managed by Linkwarden CLI\nexport LINKWARDEN_TOKEN=\"second\"\n", string(data),
		"appends, never rewrites; the last export wins when the shell starts")
}
