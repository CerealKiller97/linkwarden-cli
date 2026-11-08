package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// OutputMode is how commands print their results on stdout. It is set from the
// global --json and --quiet flags.
type OutputMode int

const (
	// OutputTable prints the CLI's usual tables (the default).
	OutputTable OutputMode = iota
	// OutputJSON prints results as indented JSON, using the API's field names.
	OutputJSON
	// OutputIDs prints only the ID of each result, one per line, for piping
	// into another command. Tables without an ID column print nothing.
	OutputIDs
)

var outputMode = OutputTable

// SetOutputMode sets how results are printed; see OutputMode.
func SetOutputMode(m OutputMode) {
	outputMode = m
}

// GetOutputMode reports how results are printed; see OutputMode.
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
func PrintRecord(v any, id int) error {
	return printRecord(os.Stdout, outputMode, v, id)
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
func PrintFileInfo(contentType string, size int) error {
	return printFileInfo(os.Stdout, outputMode, contentType, size)
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
