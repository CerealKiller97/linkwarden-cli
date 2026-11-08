package tui

import (
	"regexp"

	"charm.land/lipgloss/v2"
)

// hexColor matches #rgb and #rrggbb, the forms Linkwarden stores colors in.
var hexColor = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

// Swatch is a dot drawn in hex, for showing a color instead of its code. ok is
// false, and the swatch empty, when hex isn't a #rgb or #rrggbb color.
func Swatch(hex string) (swatch string, ok bool) {
	if !hexColor.MatchString(hex) {
		return "", false
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render("●"), true
}
