package config

import (
	"os"
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderConfigTable(cfg *linkwarden.PublicConfig) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Instance Configuration")
	t.AppendHeader(table.Row{"Key", "Value"})

	t.AppendRow(table.Row{"DISABLE_REGISTRATION", boolValue(cfg.DisableRegistration)})
	t.AppendRow(table.Row{"ADMIN", intValue(cfg.Admin)})
	t.AppendRow(table.Row{"RSS_POLLING_INTERVAL_MINUTES", intValue(cfg.RSSPollingIntervalMinutes)})
	t.AppendRow(table.Row{"EMAIL_PROVIDER", boolValue(cfg.EmailProvider)})
	t.AppendRow(table.Row{"MAX_FILE_BUFFER", intValue(cfg.MaxFileBuffer)})
	t.AppendRow(table.Row{"USER_CONTENT_DOMAIN", stringValue(cfg.UserContentDomain)})
	t.AppendRow(table.Row{"AI_ENABLED", boolValue(cfg.AIEnabled)})
	t.AppendRow(table.Row{"INSTANCE_VERSION", stringValue(cfg.InstanceVersion)})

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
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
