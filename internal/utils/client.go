package utils

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// NewLinkwardenClient builds a Linkwarden API client from the instance URL and
// token configured via SetVars.
//
// Legacy, removed once every command reads the App: use the Client of the App passed to the command.
func NewLinkwardenClient() *linkwarden.Client {
	return app.New(app.WithCredentials(GetLinkwardenToken(), GetInstanceURL())).Client()
}
