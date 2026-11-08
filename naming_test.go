package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// TestTestFunctionNamesArePascalCase enforces the repo's test naming convention:
// test, benchmark and fuzz functions are PascalCase with no underscores
// (TestCreateUsesParentIDFlag, not TestCreate_UsesParentIDFlag). Neither gofmt
// nor golangci-lint's linters check test names, so this test does.
//
// Example functions are exempt: Go uses the underscore in their names to mark
// variants (ExampleClient_List).
func TestTestFunctionNamesArePascalCase(t *testing.T) {
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "vendor" || d.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}

			if isTestFunc(fn.Name.Name) && strings.Contains(fn.Name.Name, "_") {
				t.Errorf("%s: %s: use PascalCase without underscores, e.g. %s",
					fset.Position(fn.Pos()), fn.Name.Name, strings.ReplaceAll(fn.Name.Name, "_", ""))
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func isTestFunc(name string) bool {
	for _, prefix := range []string{"Test", "Benchmark", "Fuzz"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}
