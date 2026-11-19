package utils

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// OpenBrowser opens rawURL in the user's default browser without waiting for it.
// Only absolute http(s) URLs are accepted, so the value handed to the opener can
// never be read as a command-line flag or a local file/app to launch.
func OpenBrowser(rawURL string) error {
	target, err := browserURL(rawURL)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	// Fixed opener binaries and no shell; target is a validated http(s) URL.
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target) // #nosec G204 -- validated http(s) URL
	default:
		cmd = exec.Command("xdg-open", target) // #nosec G204 -- validated http(s) URL
	}

	return cmd.Start()
}

// browserURL parses rawURL and returns it re-encoded, or an error unless it's an
// absolute http or https URL with a host.
func browserURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL %q: %w", rawURL, err)
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("refusing to open %q: only http(s) URLs are supported", rawURL)
	}

	return u.String(), nil
}
