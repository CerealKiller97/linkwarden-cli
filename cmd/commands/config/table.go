package config

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderConfigTable(ctx context.Context, cfg *linkwarden.PublicConfig) {
	t := utils.NewSummary(ctx, "Instance Configuration")
	t.Header("Setting", "Value")

	t.Row("DISABLE_REGISTRATION", cfg.DisableRegistration)
	t.Row("ADMIN", cfg.Admin)
	t.Row("RSS_POLLING_INTERVAL_MINUTES", cfg.RSSPollingIntervalMinutes)
	t.Row("EMAIL_PROVIDER", cfg.EmailProvider)
	t.Row("MAX_FILE_BUFFER", cfg.MaxFileBuffer)
	t.Row("USER_CONTENT_DOMAIN", cfg.UserContentDomain)
	t.Row("AI_ENABLED", cfg.AIEnabled)
	t.Row("INSTANCE_VERSION", cfg.InstanceVersion)

	t.JSON(cfg)
	t.Render()
}
