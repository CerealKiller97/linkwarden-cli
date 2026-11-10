// Package initcmd implements `linkwarden init`, the interactive first-run setup.
// (It can't be called "init": that name is reserved in Go.)
package initcmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

// SkipAuthAnnotation marks a command that must run without credentials, so the
// root command doesn't demand LINKWARDEN_TOKEN before it.
const SkipAuthAnnotation = "linkwarden/skip-auth"

// cliTokenName names the access token minted for the Email / Password option and
// the session used to create it.
const cliTokenName = "linkwarden-cli"

// noExpiry is the Tokens.Create expiry value for a token that never expires.
const noExpiry = 4

// signInMethod is one entry in the "How would you like to set up your account?" menu.
type signInMethod struct {
	Label string
	// Provider is the SSO provider's display name (e.g. "Google"); empty for the
	// access token and email/password methods.
	Provider string
	Password bool
}

var (
	methodToken    = signInMethod{Label: "Access Token (Recommended)"}
	methodPassword = signInMethod{Label: "Email / Password", Password: true}
)

// signInMethods builds the menu from the instance's login config: the access
// token option always, email/password unless the instance disabled it, and one
// "Sign in with <provider>" entry per SSO button (Google, Apple, ...). Without a
// config (the instance couldn't be asked), it offers the first two.
func signInMethods(cfg *linkwarden.LoginConfig) []signInMethod {
	methods := []signInMethod{methodToken}

	if cfg == nil {
		return append(methods, methodPassword)
	}

	if cfg.CredentialsEnabled == nil || !strings.EqualFold(*cfg.CredentialsEnabled, "false") {
		methods = append(methods, methodPassword)
	}

	for _, b := range cfg.ButtonAuths {
		name := strings.TrimSpace(b.Name)
		if name == "" {
			name = b.Method
		}
		if name == "" {
			continue
		}
		methods = append(methods, signInMethod{Label: "Sign in with " + name, Provider: name})
	}

	return methods
}

// openBrowser opens a URL in the user's browser; a variable so tests don't.
var openBrowser = utils.OpenBrowser

// newClient builds an SDK client; a variable so tests can point it elsewhere.
var newClient = func(baseURL, token string) *linkwarden.Client {
	return linkwarden.New(baseURL, token)
}

func Command(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Set up the CLI with your Linkwarden account",
		Long: "Walks you through signing in to your Linkwarden instance, with an access token " +
			"(recommended), your email and password, or a provider such as Google or Apple, " +
			"and saves the result to the CLI's config file.",
		Example: `  # Sign in (run it again any time to switch accounts or instances)
  linkwarden init`,
		Args:        cobra.NoArgs,
		Annotations: map[string]string{SkipAuthAnnotation: "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(app, cmd, args)
		},
	}
}

func run(app *app.App, cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	existing, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("reading existing config: %w", err)
	}

	fmt.Println()
	fmt.Println(banner(utils.ColorEnabled()))
	fmt.Println()

	if existing.Token != "" {
		_, answer, err := utils.Select(app,
			fmt.Sprintf("You're already set up for %s. Overwrite it?", existing.URL), []string{"No", "Yes"})
		if err != nil {
			return err
		}
		if answer != "Yes" {
			fmt.Println("Keeping your existing configuration.")
			return nil
		}
	}

	defaultURL := existing.URL
	if defaultURL == "" {
		defaultURL = utils.DefaultInstanceURL
	}

	rawURL, err := utils.InputDefault(app, "Linkwarden URL:", defaultURL, validateURL)
	if err != nil {
		return err
	}
	instanceURL := utils.NormalizeInstanceURL(rawURL)

	// Best effort: an older or unreachable instance just gets the default menu.
	loginCfg, err := newClient(instanceURL, "").Logins.Get(ctx)
	if err != nil {
		loginCfg = nil
	}

	methods := signInMethods(loginCfg)
	labels := make([]string, len(methods))
	for i, m := range methods {
		labels[i] = m.Label
	}

	idx, _, err := utils.Select(app, "How would you like to set up your account?", labels)
	if err != nil {
		return err
	}

	var token string
	switch method := methods[idx]; {
	case method.Password:
		token, err = passwordFlow(ctx, app, instanceURL)
	default:
		token, err = tokenFlow(app, instanceURL, method.Provider)
	}
	if err != nil {
		return err
	}

	user, err := verify(ctx, instanceURL, token)
	if err != nil {
		return err
	}

	path, err := utils.SaveConfig(utils.Config{URL: instanceURL, Token: token})
	if err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s Signed in as %s. Saved to %s\n",
		color(doneStyle, "✔"), user.Username, path)
	fmt.Println("You're all set. Try " + color(linkStyle, "linkwarden links list") + ".")

	return nil
}

// tokenFlow explains how to create an access token, opens the right settings
// page in the browser and asks the user to paste the token.
//
// provider, when set (e.g. "Google"), is the SSO provider the user signs in with
// first. Linkwarden's Google/Apple/... sign-in is an OAuth redirect that only
// ends in a browser session cookie, so the CLI can't receive a token from it
// directly; instead the user signs in in the browser and creates a token there.
func tokenFlow(app *app.App, instanceURL, provider string) (string, error) {
	pageURL := accessTokensURL(instanceURL)

	fmt.Println()
	for _, line := range tokenSteps(pageURL, provider, utils.ColorEnabled()) {
		fmt.Println(line)
	}
	fmt.Println()

	if err := openBrowser(pageURL); err != nil {
		fmt.Println("Couldn't open a browser; open the URL above manually.")
	}

	token, err := utils.Password(app, "Access token:", required("access token"))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(token), nil
}

// tokenSteps is the numbered how-to tokenFlow prints, with an extra sign-in
// step for an SSO provider.
func tokenSteps(pageURL, provider string, colorOn bool) []string {
	step := stepStyle
	link := linkStyle
	c := func(style lipgloss.Style, s string) string {
		if !colorOn {
			return s
		}
		return style.Render(s)
	}

	var steps []string
	if provider == "" {
		steps = append(steps, fmt.Sprintf("Open %s (opening it in your browser now)", c(link, pageURL)))
	} else {
		steps = append(steps,
			fmt.Sprintf("Your browser is opening %s", c(link, pageURL)),
			fmt.Sprintf("Click %s and finish signing in. If you land on the dashboard, "+
				"go to Settings → Access Tokens", c(link, "Sign in with "+provider)),
		)
	}
	steps = append(steps,
		fmt.Sprintf("Click %s", c(link, "New Access Token")),
		fmt.Sprintf("Name it (e.g. %q), choose an expiry and click %s", cliTokenName, c(link, "Create Access Token")),
		"Copy the token and paste it below. It's only shown once.",
	)

	heading := "Create a Personal Access Token in Linkwarden:"
	if provider != "" {
		heading = fmt.Sprintf("Sign in with %s in your browser, then create an access token for the CLI:", provider)
	}

	lines := []string{c(step, heading)}
	for i, s := range steps {
		lines = append(lines, fmt.Sprintf("%s %s", c(step, fmt.Sprintf("%d.", i+1)), s))
	}

	return lines
}

// passwordFlow asks for email/username and password, then exchanges them for a
// long-lived access token.
func passwordFlow(ctx context.Context, app *app.App, instanceURL string) (string, error) {
	username, err := utils.Input(app, "Email or username:", required("email or username"))
	if err != nil {
		return "", err
	}

	password, err := utils.Password(app, "Password:", required("password"))
	if err != nil {
		return "", err
	}

	return loginForToken(ctx, instanceURL, strings.TrimSpace(username), password)
}

// loginForToken signs in with a username and password, then uses that session to
// mint a non-expiring access token for the CLI, so the password is never stored.
func loginForToken(ctx context.Context, instanceURL, username, password string) (string, error) {
	session, err := newClient(instanceURL, "").Session.Create(ctx, linkwarden.CreateSessionRequest{
		Username:    username,
		Password:    password,
		SessionName: cliTokenName,
	})
	if err != nil {
		if hasStatus(err, http.StatusUnauthorized) {
			return "", errors.New("sign-in failed: check your email/username and password " +
				"(Linkwarden also refuses accounts whose email isn't verified yet)")
		}

		return "", fmt.Errorf("sign-in failed: %w", err)
	}

	created, err := newClient(instanceURL, session).Tokens.Create(ctx, cliTokenName, noExpiry)
	if err != nil {
		return "", fmt.Errorf("creating an access token for the CLI: %w", err)
	}

	return created.SecretKey, nil
}

// verify checks token against the instance and returns the signed-in user.
func verify(ctx context.Context, instanceURL, token string) (*linkwarden.User, error) {
	user, err := newClient(instanceURL, token).Users.Me(ctx)
	if err != nil {
		if hasStatus(err, http.StatusUnauthorized) {
			return nil, errors.New("that token was rejected by " + instanceURL + ": check it and run init again")
		}

		return nil, fmt.Errorf("couldn't reach %s: %w", instanceURL, err)
	}

	return user, nil
}

// accessTokensURL is the Linkwarden settings page where access tokens are created.
func accessTokensURL(instanceURL string) string {
	return instanceURL + "/settings/access-tokens"
}

func validateURL(s string) error {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("enter a full URL, e.g. https://links.example.com")
	}

	return nil
}

func required(name string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s is required", name)
		}

		return nil
	}
}

func hasStatus(err error, status int) bool {
	var apiErr *linkwarden.APIError

	return errors.As(err, &apiErr) && apiErr.StatusCode == status
}

// Styles share pkg/tui's Linkwarden palette with the rest of the CLI.
var (
	stepStyle = lipgloss.NewStyle().Bold(true).Foreground(tui.Sky)
	linkStyle = lipgloss.NewStyle().Foreground(tui.Ocean)
	doneStyle = lipgloss.NewStyle().Bold(true).Foreground(tui.Ocean)
)

func color(style lipgloss.Style, s string) string {
	if !utils.ColorEnabled() {
		return s
	}

	return style.Render(s)
}
