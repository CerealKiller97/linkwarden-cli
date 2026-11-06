package config

import (
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderConfigTable(cfg *linkwarden.PublicConfig) {
	t := utils.NewTable("Instance Configuration")
	t.Header("Key", "Value")

	t.Row("DISABLE_REGISTRATION", boolValue(cfg.DisableRegistration))
	t.Row("ADMIN", intValue(cfg.Admin))
	t.Row("RSS_POLLING_INTERVAL_MINUTES", intValue(cfg.RSSPollingIntervalMinutes))
	t.Row("EMAIL_PROVIDER", boolValue(cfg.EmailProvider))
	t.Row("MAX_FILE_BUFFER", intValue(cfg.MaxFileBuffer))
	t.Row("USER_CONTENT_DOMAIN", stringValue(cfg.UserContentDomain))
	t.Row("AI_ENABLED", boolValue(cfg.AIEnabled))
	t.Row("INSTANCE_VERSION", stringValue(cfg.InstanceVersion))

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
