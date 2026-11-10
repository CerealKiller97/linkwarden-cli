package utils

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// NewLinkwardenClient builds a Linkwarden API client from the instance URL and
// token configured via SetVars.
//
// Legacy, removed once every command reads the App: use app.From(cmd.Context()).Client().
func NewLinkwardenClient() *linkwarden.Client {
	return app.New(app.WithCredentials(GetLinkwardenToken(), GetInstanceURL())).Client()
}
