package utils

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/rs/zerolog/log"
)

// NewLinkwardenClient builds a Linkwarden API client from the instance URL and
// token configured via SetVars (LINKWARDEN_URL / LINKWARDEN_TOKEN).
func NewLinkwardenClient() *linkwarden.Client {
	return linkwarden.New(GetInstanceURL(), GetLinkwardenToken(), linkwarden.WithLogger(log.Logger))
}
