package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/require"
)

func TestColorEnabledFor(t *testing.T) {
	assert := require.New(t)

	assert.True(colorEnabledFor("", true), "no NO_COLOR, real terminal -> color on")
	assert.False(colorEnabledFor("", false), "no NO_COLOR, not a terminal -> color off")
	assert.False(colorEnabledFor("1", true), "NO_COLOR set -> color off even on a terminal")
	assert.False(colorEnabledFor("1", false), "NO_COLOR set, not a terminal -> color off")
}

func TestTableStyleFor_ColorDisabled(t *testing.T) {
	assert := require.New(t)

	style := tableStyleFor(false)
	assert.Equal(table.StyleRounded.Box, style.Box, "borders stay rounded regardless of color")
	assert.Empty(style.Color.Header, "no header color applied when color is disabled")
}

func TestTableStyleFor_ColorEnabled(t *testing.T) {
	assert := require.New(t)

	style := tableStyleFor(true)
	assert.Equal(table.StyleRounded.Box, style.Box, "borders stay rounded")
	assert.NotEmpty(style.Color.Header, "header gets a color when color is enabled")
}

func TestNewTable_ReturnsConfiguredWriter(t *testing.T) {
	assert := require.New(t)

	tbl := NewTable("Example")
	assert.NotNil(tbl)

	tbl.AppendHeader(table.Row{"A", "B"})
	tbl.AppendRow(table.Row{1, 2})
	assert.Equal(1, tbl.Length())
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		width int
		want  string
	}{
		{"fits", "hello", 5, "hello"},
		{"no limit", "hello world", 0, "hello world"},
		{"ascii cut", "hello world", 8, "hello w…"},
		{"trailing space before ellipsis is dropped", "hello world", 7, "hello…"},
		{"accented letters are one cell", "Ćacilend Čačak", 9, "Ćacilend…"},
		{"emoji is two cells and isn't split", "ab🫧cd", 4, "ab…"},
		{"emoji fits exactly", "ab🫧cd", 5, "ab🫧…"},
		{"cjk", "生命生命", 5, "生命…"},
		{"width one", "hello", 1, "…"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := require.New(t)

			got := Truncate(tc.in, tc.width)
			assert.Equal(tc.want, got)
			if tc.width > 0 {
				assert.LessOrEqual(text.StringWidthWithoutEscSequences(got), tc.width)
			}
		})
	}
}

func TestFitWidths(t *testing.T) {
	t.Run("fits unchanged", func(t *testing.T) {
		got := fitWidths([]int{2, 30, 40}, []int{2, 12, 12}, []bool{false, true, true}, 100, 0)
		require.Equal(t, []int{2, 30, 40}, got)
	})

	t.Run("flex columns share the remainder", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 72, 10}, []int{2, 20, 20, 10}, []bool{false, true, true, false}, 80, 0)
		require.Equal(t, []int{2, 34, 34, 10}, got)
	})

	t.Run("a narrow flex column gives its share away", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 10}, []int{2, 12, 10}, []bool{false, true, true}, 60, 0)
		require.Equal(t, []int{2, 48, 10}, got)
	})

	t.Run("minimums win when the terminal is too narrow", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 72}, []int{2, 20, 20}, []bool{false, true, true}, 20, 0)
		require.Equal(t, []int{2, 20, 20}, got)
	})

	t.Run("no flex columns", func(t *testing.T) {
		got := fitWidths([]int{50, 50}, []int{50, 50}, []bool{false, false}, 20, 0)
		require.Equal(t, []int{50, 50}, got)
	})

	t.Run("title width is respected", func(t *testing.T) {
		got := fitWidths([]int{2, 100}, []int{2, 12}, []bool{false, true}, 10, 30)
		require.Equal(t, []int{2, 28}, got)
	})
}

func TestListTable_FitsTerminal(t *testing.T) {
	assert := require.New(t)

	tbl := NewListTable("Links",
		Column{Header: "ID"},
		Column{Header: "Name", Flex: true, MinWidth: 20},
		Column{Header: "URL", Flex: true, MinWidth: 20},
		Column{Header: "Collection", Flex: true},
		Column{Header: "Tags", Flex: true, HideWhenEmpty: true},
	)
	tbl.AppendRow(43, "charmbracelet/bubbles: TUI components for Bubble Tea 🫧", "https://github.com/charmbracelet/bubbles", "Work", "")
	tbl.AppendRow(41, "refactor: move the whole TUI to Bubble Tea, Bubbles and Lip Gloss by CerealKiller97 · Pull Request #20 · CerealKiller97/linkwarden-cli · GitHub", "https://github.com/CerealKiller97/linkwarden-cli/pull/20/files?diff=split&w=1", "Work", "")
	tbl.AppendRow(7, "Ćacilend\nwith a line break", "https://example.com", "Ćacilend", "")

	var buf bytes.Buffer
	tbl.RenderTo(&buf, 120, false)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	assert.Len(lines, 9, "top, title, separator, header, separator, 3 rows, bottom")

	tableWidth := text.StringWidthWithoutEscSequences(lines[0])
	assert.LessOrEqual(tableWidth, 120)

	for _, line := range lines {
		assert.Equal(tableWidth, text.StringWidthWithoutEscSequences(line), "every line is as wide as the box: %q", line)
		assert.Equal(strings.TrimRight(line, " "), line, "no trailing whitespace")
	}

	out := buf.String()
	assert.NotContains(out, "TAGS", "empty TAGS column is dropped")
	assert.Contains(out, "…", "long values are shortened")
	assert.Contains(out, "Ćacilend with a line break", "line breaks become spaces")
}

func TestListTable_NoLimitKeepsFullValues(t *testing.T) {
	assert := require.New(t)

	long := strings.Repeat("very long name ", 20)
	tbl := NewListTable("Tags", Column{Header: "ID"}, Column{Header: "Name", Flex: true}, Column{Header: "Tags", HideWhenEmpty: true})
	tbl.AppendRow(1, long, "go, cli")

	var buf bytes.Buffer
	tbl.RenderTo(&buf, 0, false)

	assert.Contains(buf.String(), strings.TrimSpace(long))
	assert.NotContains(buf.String(), "…")
	assert.Contains(buf.String(), "TAGS", "non-empty columns are kept")
}

func TestListTable_ShortTitleTableIsNotPadded(t *testing.T) {
	assert := require.New(t)

	tbl := NewListTable("A much longer title than the cells", Column{Header: "ID"}, Column{Header: "Name", Flex: true})
	tbl.AppendRow(1, strings.Repeat("x", 200))

	var buf bytes.Buffer
	tbl.RenderTo(&buf, 30, false)

	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		assert.Equal(strings.TrimRight(line, " "), line, "no trailing whitespace")
	}
}
