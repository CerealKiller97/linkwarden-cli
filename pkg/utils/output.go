package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// OutputMode is how commands print their results; see app.OutputMode.
type OutputMode = app.OutputMode

const (
	OutputTable = app.OutputTable
	OutputJSON  = app.OutputJSON
	OutputIDs   = app.OutputIDs
)

// outputMode is the legacy output mode, used only when a command runs without
// an *app.App.
var outputMode = OutputTable

// SetOutputMode sets the legacy output mode.
//
// Legacy, removed once every command reads the App: run the command with an *app.App (app.WithMode in tests).
func SetOutputMode(m OutputMode) {
	outputMode = m
}

// GetOutputMode reports the legacy output mode.
//
// Legacy, removed once every command reads the App: read the App's Mode.
func GetOutputMode() OutputMode {
	return outputMode
}

// WriteJSON writes v to w as indented JSON. A nil slice is written as [] so
// an empty list is always a JSON array, never null.
func WriteJSON(w io.Writer, v any) error {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Slice && rv.IsNil() {
		v = []any{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	return enc.Encode(v)
}

// PrintRecord prints a single result of a command that otherwise only logs
// what it did (create, update): as JSON with --json, its ID with --quiet, and
// nothing in table mode, where the log line says it all.
func PrintRecord(app *app.App, v any, id int) error {
	app = orLegacy(app)
	return printRecord(app.Stdout(), app.Mode, v, id)
}

func printRecord(w io.Writer, mode OutputMode, v any, id int) error {
	switch mode {
	case OutputJSON:
		return WriteJSON(w, v)
	case OutputIDs:
		_, err := fmt.Fprintln(w, id)
		return err
	default:
		return nil
	}
}

// fileInfo is what --json prints for a downloaded file that wasn't saved.
type fileInfo struct {
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
}

// PrintFileInfo prints the content type and size of a file that was fetched
// but not saved (no --output): a line of text, a JSON object with --json, and
// nothing with --quiet.
func PrintFileInfo(app *app.App, contentType string, size int) error {
	app = orLegacy(app)
	return printFileInfo(app.Stdout(), app.Mode, contentType, size)
}

func printFileInfo(w io.Writer, mode OutputMode, contentType string, size int) error {
	switch mode {
	case OutputJSON:
		return WriteJSON(w, fileInfo{ContentType: contentType, Size: size})
	case OutputIDs:
		return nil
	default:
		_, err := fmt.Fprintf(w, "Content-Type: %s, Size: %d bytes\n", contentType, size)
		return err
	}
}
