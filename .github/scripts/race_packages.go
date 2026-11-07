// Command race_packages classifies packages for decision G (task 0.7): a
// package runs under the race detector exactly when its non-test Go files
// contain a go statement. It parses Go rather than matching text, so comments,
// strings and goto don't count, and files other platforms build count too.
//
// It lists packages with and without -race, because a package whose files all
// need //go:build race (or !race) exists in only one mode. It prints
// {"packages": [...], "race": [...]}. It refuses, exiting 1, when a package
// doesn't build in the mode it runs in, or has tests that build only in the
// other mode, since those would never run.
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
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type listedPackage struct {
	ImportPath     string
	Dir            string
	GoFiles        []string
	CgoFiles       []string
	IgnoredGoFiles []string
	TestGoFiles    []string
	XTestGoFiles   []string
}

type classification struct {
	Packages []string `json:"packages"`
	Race     []string `json:"race"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "race_packages:", err)
		os.Exit(1)
	}
}

func run(patterns []string, out io.Writer) error {
	plain, err := list(patterns, false)
	if err != nil {
		return err
	}
	race, err := list(patterns, true)
	if err != nil {
		return err
	}
	files := &parsed{fset: token.NewFileSet(), files: map[string]*ast.File{}}
	result := classification{Packages: []string{}, Race: []string{}}
	var problems []string
	for _, path := range union(plain, race) {
		starts, err := files.startsGoroutines(plain[path], race[path])
		if err != nil {
			return err
		}
		problem, err := files.unrunnable(path, starts, plain[path], race[path])
		if err != nil {
			return err
		}
		if problem != "" {
			problems = append(problems, problem)
			continue
		}
		result.Packages = append(result.Packages, path)
		if starts {
			result.Race = append(result.Race, path)
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "\n"))
	}
	if err := json.NewEncoder(out).Encode(result); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func list(patterns []string, race bool) (map[string]*listedPackage, error) {
	args := []string{"list", "-json"}
	if race {
		args = append(args, "-race")
	}
	listed, err := exec.Command("go", append(append(args, "--"), patterns...)...).Output()
	if err != nil {
		return nil, fmt.Errorf("go %s: %w", strings.Join(args, " "), err)
	}
	packages := map[string]*listedPackage{}
	decoder := json.NewDecoder(bytes.NewReader(listed))
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); errors.Is(err, io.EOF) {
			return packages, nil
		} else if err != nil {
			return nil, fmt.Errorf("decode go list: %w", err)
		}
		packages[pkg.ImportPath] = &pkg
	}
}

func union(plain, race map[string]*listedPackage) []string {
	var paths []string
	for path := range plain {
		paths = append(paths, path)
	}
	for path := range race {
		if plain[path] == nil {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}

type parsed struct {
	fset  *token.FileSet
	files map[string]*ast.File
}

func (p *parsed) parse(dir, name string) (*ast.File, error) {
	path := filepath.Join(dir, name)
	if file, ok := p.files[path]; ok {
		return file, nil
	}
	file, err := parser.ParseFile(p.fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	p.files[path] = file
	return file, nil
}

func (p *parsed) startsGoroutines(listings ...*listedPackage) (bool, error) {
	for _, pkg := range listings {
		if pkg == nil {
			continue
		}
		for _, name := range append(append(append([]string{}, pkg.GoFiles...), pkg.CgoFiles...), pkg.IgnoredGoFiles...) {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			file, err := p.parse(pkg.Dir, name)
			if err != nil {
				return false, err
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
	}
	return false, nil
}

// unrunnable describes why a package can't be tested exactly in the mode it
// runs in, or returns "".
func (p *parsed) unrunnable(path string, starts bool, plain, race *listedPackage) (string, error) {
	own, other, reason, otherMode := plain, race, "starts no goroutines, so it runs without -race", "with -race"
	if starts {
		own, other, reason, otherMode = race, plain, "starts goroutines, so it runs with -race", "without -race"
	}
	if own == nil {
		return fmt.Sprintf("%s %s, but it builds only %s", path, reason, otherMode), nil
	}
	ownTests, err := p.testNames(own)
	if err != nil {
		return "", err
	}
	otherTests, err := p.testNames(other)
	if err != nil {
		return "", err
	}
	var never []string
	for name := range otherTests {
		if !ownTests[name] {
			never = append(never, name)
		}
	}
	if len(never) == 0 {
		return "", nil
	}
	sort.Strings(never)
	return fmt.Sprintf("%s %s; these tests build only %s, so they would never run: %s", path, reason, otherMode, strings.Join(never, ", ")), nil
}

func (p *parsed) testNames(pkg *listedPackage) (map[string]bool, error) {
	names := map[string]bool{}
	if pkg == nil {
		return names, nil
	}
	for _, name := range append(append([]string{}, pkg.TestGoFiles...), pkg.XTestGoFiles...) {
		file, err := p.parse(pkg.Dir, name)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && isTestName(fn.Name.Name) {
				names[fn.Name.Name] = true
			}
		}
	}
	return names, nil
}

// isTestName follows go test: Test, Example or Fuzz, then nothing or a
// non-lower-case rune. TestMain is setup, not a test.
func isTestName(name string) bool {
	if name == "TestMain" {
		return false
	}
	for _, prefix := range []string{"Test", "Example", "Fuzz"} {
		if rest, ok := strings.CutPrefix(name, prefix); ok {
			r, _ := utf8.DecodeRuneInString(rest)
			return rest == "" || !unicode.IsLower(r)
		}
	}
	return false
}
