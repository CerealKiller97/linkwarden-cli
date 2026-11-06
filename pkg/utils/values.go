package utils

import "strings"

var (
	LinkwardenToken string
	InstanceURL     string
	NonInteractive  bool
)

func SetVars(linkwardenToken string, instanceURL string) {
	LinkwardenToken = linkwardenToken
	InstanceURL = instanceURL
}

func GetLinkwardenToken() string {
	return LinkwardenToken
}

func GetInstanceURL() string {
	return InstanceURL
}

// SetNonInteractive forces (or releases) non-interactive mode, overriding the
// default TTY auto-detection used by ResolveID. Set from the --non-interactive
// root flag.
func SetNonInteractive(nonInteractive bool) {
	NonInteractive = nonInteractive
}

func GetNonInteractive() bool {
	return NonInteractive
}

// DefaultInstanceURL is used when no instance URL is configured.
const DefaultInstanceURL = "https://cloud.linkwarden.app"

// NormalizeInstanceURL drops surrounding whitespace and trailing slashes (so
// "https://links.example.com/" works) and falls back to DefaultInstanceURL when
// nothing is left.
func NormalizeInstanceURL(raw string) string {
	u := strings.TrimRight(strings.TrimSpace(raw), "/")
	if u == "" {
		return DefaultInstanceURL
	}

	return u
}
