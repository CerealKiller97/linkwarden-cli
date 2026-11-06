package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/require"
)

func TestColorEnabledFor(t *testing.T) {
	assert := require.New(t)

	assert.True(colorEnabledFor("", true), "no NO_COLOR, real terminal -> color on")
	assert.False(colorEnabledFor("", false), "no NO_COLOR, not a terminal -> color off")
	assert.False(colorEnabledFor("1", true), "NO_COLOR set -> color off even on a terminal")
	assert.False(colorEnabledFor("1", false), "NO_COLOR set, not a terminal -> color off")
}

func TestShouldBrowse(t *testing.T) {
	assert := require.New(t)

	assert.False(shouldBrowse(5, 40), "fits on screen -> just print")
	assert.True(shouldBrowse(100, 40), "taller than the screen -> viewer")
	assert.False(shouldBrowse(100, 0), "unknown height -> just print")
}

func TestTable_RendersStaticWhenNotATerminal(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	tbl := NewTable("Example")
	tbl.out = &out
	tbl.Header("ID", "Name")
	tbl.Row(1, "first")
	tbl.Row(2, nil)

	tbl.Render()

	assert.Equal(2, tbl.Len())
	got := out.String()
	assert.True(strings.HasPrefix(got, "Example\n"))
	assert.Contains(got, "╭", "rounded border, same as the viewer")
	assert.Contains(got, "first")
	assert.Contains(got, "<nil>")
	assert.NotContains(got, "\x1b[", "test output isn't a terminal, so no ANSI")
}

func TestStartProgress_NoOpWithoutTerminal(t *testing.T) {
	p := StartProgress("x", 3)
	p.Increment()
	p.Stop()

	require.IsType(t, noProgress{}, p)
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
				assert.LessOrEqual(ansi.StringWidth(got), tc.width)
			}
		})
	}
}

func TestFitWidths(t *testing.T) {
	t.Run("fits unchanged", func(t *testing.T) {
		got := fitWidths([]int{2, 30, 40}, []int{2, 12, 12}, []bool{false, true, true}, 100)
		require.Equal(t, []int{2, 30, 40}, got)
	})

	t.Run("flex columns share the remainder", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 72, 10}, []int{2, 20, 20, 10}, []bool{false, true, true, false}, 80)
		require.Equal(t, []int{2, 34, 34, 10}, got)
	})

	t.Run("a narrow flex column gives its share away", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 10}, []int{2, 12, 10}, []bool{false, true, true}, 60)
		require.Equal(t, []int{2, 48, 10}, got)
	})

	t.Run("minimums win when the terminal is too narrow", func(t *testing.T) {
		got := fitWidths([]int{2, 136, 72}, []int{2, 20, 20}, []bool{false, true, true}, 20)
		require.Equal(t, []int{2, 20, 20}, got)
	})

	t.Run("no flex columns", func(t *testing.T) {
		got := fitWidths([]int{50, 50}, []int{50, 50}, []bool{false, false}, 20)
		require.Equal(t, []int{50, 50}, got)
	})
}

func newLinksLikeTable(out *bytes.Buffer) *Table {
	tbl := NewTable("Links")
	tbl.out = out
	tbl.Columns(
		Column{Header: "ID"},
		Column{Header: "Name", Flex: true, MinWidth: 20},
		Column{Header: "URL", Flex: true, MinWidth: 20},
		Column{Header: "Collection", Flex: true},
		Column{Header: "Tags", Flex: true, HideWhenEmpty: true},
	)
	tbl.Row(43, "charmbracelet/bubbles: TUI components for Bubble Tea 🫧", "https://github.com/charmbracelet/bubbles", "Work", "")
	tbl.Row(41, "refactor: move the whole TUI to Bubble Tea, Bubbles and Lip Gloss by CerealKiller97 · Pull Request #20 · CerealKiller97/linkwarden-cli · GitHub", "https://github.com/CerealKiller97/linkwarden-cli/pull/20/files?diff=split&w=1", "Work", "")
	tbl.Row(7, "Ćacilend\nwith a line break", "https://example.com", "Ćacilend", "")

	return tbl
}

func TestTable_FitsTerminal(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	newLinksLikeTable(&out).renderStatic(120, false)

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	assert.Len(lines, tui.TableChromeHeight+3, "title, borders, header, rule and 3 rows")
	assert.Equal("Links", lines[0])

	box := lines[1:]
	tableWidth := ansi.StringWidth(box[0])
	assert.LessOrEqual(tableWidth, 120)

	for _, line := range box {
		assert.Equal(tableWidth, ansi.StringWidth(line), "every line is as wide as the box: %q", line)
		assert.Equal(strings.TrimRight(line, " "), line, "no trailing whitespace")
	}

	got := out.String()
	assert.NotContains(got, "Tags", "empty Tags column is dropped")
	assert.Contains(got, "…", "long values are shortened")
	assert.Contains(got, "Ćacilend with a line break", "line breaks become spaces")
}

func TestTable_NoLimitKeepsFullValues(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	newLinksLikeTable(&out).renderStatic(0, false)

	assert.Contains(out.String(), "Pull Request #20 · CerealKiller97/linkwarden-cli · GitHub")
	assert.NotContains(out.String(), "…")
}

func TestTable_HiddenColumnKeepsTone(t *testing.T) {
	assert := require.New(t)

	tbl := NewTable("T")
	tbl.Columns(Column{Header: "A"}, Column{Header: "Empty", HideWhenEmpty: true}, Column{Header: "C"})
	tbl.Row("a", "", "c")
	tbl.Tone(func(_, col int, _ string) Tone {
		if col == 2 {
			return ToneDanger
		}

		return ToneNone
	})

	data := tbl.visible()
	assert.Equal([]string{"A", "C"}, data.Headers)
	assert.Equal([][]string{{"a", "c"}}, data.Rows)
	assert.Equal(ToneDanger, data.Tone(0, 1), "column 1 of the visible table is original column 2")
	assert.Equal(ToneNone, data.Tone(0, 0))
}
