package utils

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
)

// NewLinkwardenClient builds a Linkwarden API client from the instance URL and
// token configured via SetVars (LINKWARDEN_URL / LINKWARDEN_TOKEN).
func NewLinkwardenClient() *linkwarden.Client {
	return linkwarden.New(GetInstanceURL(), GetLinkwardenToken(), linkwarden.WithLogger(slogBridge()))
}

// slogBridge adapts the global zerolog logger to the *slog.Logger the SDK
// expects. zerolog.Logger is an io.Writer, so slog's text output is forwarded
// through it; slog's own timestamp is dropped because zerolog adds one.
func slogBridge() *slog.Logger {
	return slog.New(slog.NewTextHandler(log.Logger, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}
