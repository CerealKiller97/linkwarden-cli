package config

import (
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderConfigTable(cfg *linkwarden.PublicConfig) {
	t := utils.NewListTable("Instance Configuration",
		utils.Column{Header: "Key"},
		utils.Column{Header: "Value", Flex: true},
	)

	t.AppendRow("DISABLE_REGISTRATION", boolValue(cfg.DisableRegistration))
	t.AppendRow("ADMIN", intValue(cfg.Admin))
	t.AppendRow("RSS_POLLING_INTERVAL_MINUTES", intValue(cfg.RSSPollingIntervalMinutes))
	t.AppendRow("EMAIL_PROVIDER", boolValue(cfg.EmailProvider))
	t.AppendRow("MAX_FILE_BUFFER", intValue(cfg.MaxFileBuffer))
	t.AppendRow("USER_CONTENT_DOMAIN", stringValue(cfg.UserContentDomain))
	t.AppendRow("AI_ENABLED", boolValue(cfg.AIEnabled))
	t.AppendRow("INSTANCE_VERSION", stringValue(cfg.InstanceVersion))

	t.Render()
}

func boolValue(v *bool) string {
	if v == nil {
		return "N/A"
	}

	return strconv.FormatBool(*v)
}

func intValue(v *int) string {
	if v == nil {
		return "N/A"
	}

	return strconv.Itoa(*v)
}

func stringValue(v *string) string {
	if v == nil {
		return "N/A"
	}

	return *v
}
