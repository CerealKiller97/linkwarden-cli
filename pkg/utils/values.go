package utils

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
