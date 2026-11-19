package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeInstanceURL(t *testing.T) {
	cases := map[string]string{
		"":                            DefaultInstanceURL,
		"   ":                         DefaultInstanceURL,
		"/":                           DefaultInstanceURL,
		"https://links.example.com":   "https://links.example.com",
		"https://links.example.com/":  "https://links.example.com",
		"https://links.example.com//": "https://links.example.com",
		" http://localhost:3100/ ":    "http://localhost:3100",
		"https://example.com/lw/":     "https://example.com/lw",
	}

	for in, want := range cases {
		require.Equal(t, want, NormalizeInstanceURL(in), "input %q", in)
	}
}
