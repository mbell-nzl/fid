// Command race_packages prints, one per line, the import path of each package
// whose non-test Go files contain a go statement. CI and `make test` run the
// race detector on exactly these packages (decision G, task 0.7). It parses
// Go rather than matching text, so comments, strings and goto don't count.
// Files other platforms build count too.
//
// Usage, from the module root: go run .github/scripts/race_packages.go ./...
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type listedPackage struct {
	ImportPath     string
	Dir            string
	GoFiles        []string
	CgoFiles       []string
	IgnoredGoFiles []string
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "race_packages:", err)
		os.Exit(1)
	}
}

func run(patterns []string, out io.Writer) error {
	listed, err := exec.Command("go", append([]string{"list", "-json", "--"}, patterns...)...).Output()
	if err != nil {
		return fmt.Errorf("go list: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(listed))
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("decode go list: %w", err)
		}
		starts, err := startsGoroutines(pkg)
		if err != nil {
			return err
		}
		if starts {
			if _, err := fmt.Fprintln(out, pkg.ImportPath); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		}
	}
}

func startsGoroutines(pkg listedPackage) (bool, error) {
	files := append(append(append([]string{}, pkg.GoFiles...), pkg.CgoFiles...), pkg.IgnoredGoFiles...)
	fset := token.NewFileSet()
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return false, fmt.Errorf("parse %s: %w", pkg.ImportPath, err)
		}
		found := false
		ast.Inspect(file, func(node ast.Node) bool {
			if _, ok := node.(*ast.GoStmt); ok {
				found = true
			}
			return !found
		})
		if found {
			return true, nil
		}
	}
	return false, nil
}
