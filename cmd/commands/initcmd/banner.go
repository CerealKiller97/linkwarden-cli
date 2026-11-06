package initcmd

import (
	"fmt"
	"strings"
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

// brandBlue is the blue of Linkwarden's logo, as a 24-bit ANSI foreground.
const brandBlue = "\x1b[38;2;74;163;226m"

const (
	bold  = "\x1b[1m"
	dim   = "\x1b[2m"
	reset = "\x1b[0m"
)

// banner renders the logo with a welcome message beside it. With colorOn false
// (NO_COLOR, or output isn't a terminal) it's plain text.
func banner(colorOn bool) string {
	paint := func(code, s string) string {
		if !colorOn {
			return s
		}
		return code + s + reset
	}

	side := map[int]string{
		4: paint(bold+brandBlue, "Welcome to Linkwarden CLI"),
		5: paint(dim, "Let's connect you to your Linkwarden account."),
	}

	var b strings.Builder
	for i, row := range logo {
		line := paint(brandBlue, row)
		if msg, ok := side[i]; ok {
			line += strings.Repeat(" ", logoWidth-len([]rune(row))) + msg
		}
		fmt.Fprintln(&b, line)
	}

	return strings.TrimSuffix(b.String(), "\n")
}
