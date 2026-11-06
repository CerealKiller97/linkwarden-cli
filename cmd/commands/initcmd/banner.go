package initcmd

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
)

// logo is Linkwarden's lightning-bolt mark, drawn with box-drawing characters.
var logo = []string{
	`       ╱│`,
	`      ╱ │`,
	`     ╱  │`,
	`    ╱   │`,
	`   ╱    ╰───╮`,
	`  ╰───╮    ╱`,
	`      │   ╱`,
	`      │  ╱`,
	`      │ ╱`,
	`      │╱`,
}

// logoWidth pads every logo row so the text beside it lines up.
const logoWidth = 16

var (
	logoStyle  = lipgloss.NewStyle().Foreground(tui.Ocean)
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(tui.Sky)
	hintStyle  = lipgloss.NewStyle().Foreground(tui.Muted)
)

// banner renders the logo with a welcome message beside it. With colorOn false
// (NO_COLOR, or output isn't a terminal) it's plain text.
func banner(colorOn bool) string {
	paint := func(style lipgloss.Style, s string) string {
		if !colorOn {
			return s
		}
		return style.Render(s)
	}

	side := map[int]string{
		4: paint(titleStyle, "Welcome to Linkwarden CLI"),
		5: paint(hintStyle, "Let's connect you to your Linkwarden account."),
	}

	var b strings.Builder
	for i, row := range logo {
		line := paint(logoStyle, row)
		if msg, ok := side[i]; ok {
			line += strings.Repeat(" ", logoWidth-len([]rune(row))) + msg
		}
		fmt.Fprintln(&b, line)
	}

	return strings.TrimSuffix(b.String(), "\n")
}
