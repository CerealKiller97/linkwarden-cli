package tui

import "charm.land/lipgloss/v2"

// The palette is Linkwarden's own blue, sampled from its web UI so the CLI
// looks like the same product. Every component in this package takes its
// colors from here.
var (
	// Sky is the light blue of Linkwarden's primary button and icon outlines.
	Sky = lipgloss.Color("#8FD1F8")
	// Ocean is the deeper blue of its folder icons.
	Ocean = lipgloss.Color("#4AA3E2")
	// Slate is the background of the selected sidebar row.
	Slate = lipgloss.Color("#3F4C54")
	// Ink is the dark text on Sky buttons.
	Ink = lipgloss.Color("#16232B")
	// Muted is for secondary text: hints, footers, alternate table rows.
	Muted = lipgloss.Color("245")

	// Status colors stay distinct from the blues so a problem still stands out.
	Warning = lipgloss.Color("214")
	Danger  = lipgloss.Color("203")
)
