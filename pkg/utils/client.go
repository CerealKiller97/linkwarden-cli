package utils

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// NewLinkwardenClient builds a Linkwarden API client from the instance URL and
// token configured via SetVars (LINKWARDEN_URL / LINKWARDEN_TOKEN).
func NewLinkwardenClient() *linkwarden.Client {
	return linkwarden.New(GetInstanceURL(), GetLinkwardenToken(), linkwarden.WithLogger(slog.Default()))
}
