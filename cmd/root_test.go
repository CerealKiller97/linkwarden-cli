package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstanceURL(t *testing.T) {
	cases := map[string]string{
		"":                            defaultInstanceURL,
		"   ":                         defaultInstanceURL,
		"/":                           defaultInstanceURL,
		"https://links.example.com":   "https://links.example.com",
		"https://links.example.com/":  "https://links.example.com",
		"https://links.example.com//": "https://links.example.com",
		" http://localhost:3100/ ":    "http://localhost:3100",
		"https://example.com/lw/":     "https://example.com/lw",
	}

	for in, want := range cases {
		require.Equal(t, want, instanceURL(in), "input %q", in)
	}
}
