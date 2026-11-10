package utils

import (
	"context"

	"bytes"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"

	"github.com/stretchr/testify/require"
)

func TestWriteJSONNilSliceIsEmptyArray(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	var none []int
	assert.NoError(WriteJSON(&out, none))
	assert.Equal("[]\n", out.String())
}

func TestWriteJSONKeepsHTMLCharacters(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	assert.NoError(WriteJSON(&out, map[string]string{"url": "https://example.com/?a=1&b=<2>"}))
	assert.Contains(out.String(), `"https://example.com/?a=1&b=<2>"`)
}

func TestPrintRecordByMode(t *testing.T) {
	assert := require.New(t)

	record := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{ID: 7, Name: "Go"}

	var out bytes.Buffer
	assert.NoError(printRecord(&out, OutputJSON, record, record.ID))
	assert.JSONEq(`{"id":7,"name":"Go"}`, out.String())

	out.Reset()
	assert.NoError(printRecord(&out, OutputIDs, record, record.ID))
	assert.Equal("7\n", out.String())

	out.Reset()
	assert.NoError(printRecord(&out, OutputTable, record, record.ID))
	assert.Empty(out.String(), "table mode leaves it to the log line")
}

func TestPrintFileInfoByMode(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	assert.NoError(printFileInfo(&out, OutputTable, "application/pdf", 42))
	assert.Equal("Content-Type: application/pdf, Size: 42 bytes\n", out.String())

	out.Reset()
	assert.NoError(printFileInfo(&out, OutputJSON, "application/pdf", 42))
	assert.JSONEq(`{"contentType":"application/pdf","size":42}`, out.String())

	out.Reset()
	assert.NoError(printFileInfo(&out, OutputIDs, "application/pdf", 42))
	assert.Empty(out.String())
}

func newModeTable(out *bytes.Buffer, mode OutputMode) *Table {
	ctx := app.WithApp(context.Background(), &app.App{Out: out, Mode: mode})
	tbl := NewTable(ctx, "Example")
	tbl.Header("ID", "Name")
	tbl.Row(1, "first")
	tbl.Row(2, nil)

	return tbl
}

func TestTableJSONPrintsSource(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	tbl := newModeTable(&out, OutputJSON)
	tbl.JSON([]map[string]any{{"id": 1, "name": "first"}, {"id": 2, "name": nil}})
	tbl.Render()

	assert.JSONEq(`[{"id":1,"name":"first"},{"id":2,"name":null}]`, out.String())
	assert.NotContains(out.String(), "╭", "no table borders")
}

func TestTableJSONWithoutSourcePrintsRows(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	newModeTable(&out, OutputJSON).Render()

	assert.JSONEq(`[{"ID":"1","Name":"first"},{"ID":"2","Name":"-"}]`, out.String())
}

func TestTableQuietPrintsIDs(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	newModeTable(&out, OutputIDs).Render()

	assert.Equal("1\n2\n", out.String())
}

func TestTableQuietWithoutIDColumnPrintsNothing(t *testing.T) {
	assert := require.New(t)

	var out bytes.Buffer
	ctx := app.WithApp(t.Context(), &app.App{Out: &out, Mode: OutputIDs})
	tbl := NewSummary(ctx, "Totals")
	tbl.Header("Total", "Failed")
	tbl.Row(3, 0)
	tbl.Render()

	assert.Empty(out.String())
}
