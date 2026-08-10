//go:build integration

package integration

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestIntegrationClassification(t *testing.T) {
	dedicatedDir := classificationDirectory(t)
	repositoryRoot := filepath.Dir(filepath.Dir(filepath.Dir(dedicatedDir)))
	wantJourneys := []string{
		"TestMaintenanceCheckReportsIncompatibleWorkflowAssignment",
		"TestV01ClosureJourneyAndRebuildParity",
		"TestV01ParityGoldenOutputs",
		"TestV02aWorkflowAuthoringExitScenarioAndRebuildParity",
	}

	gotJourneys := topLevelTests(t, dedicatedDir,
		"TestClosureJourneyHasNoWallClockSleep",
		"TestIntegrationClassification",
		"TestIntegrationBinaryOwnedByPackage",
	)
	if !reflect.DeepEqual(gotJourneys, wantJourneys) {
		t.Fatalf("dedicated integration journeys = %v, want %v", gotJourneys, wantJourneys)
	}

	if tagged := taggedTestFilesOutside(t, repositoryRoot, dedicatedDir); len(tagged) != 0 {
		t.Fatalf("integration-tagged tests outside dedicated package = %v, want none", tagged)
	}
}

func TestClosureJourneyHasNoWallClockSleep(t *testing.T) {
	path := filepath.Join(classificationDirectory(t), "e2e_v01_closure_test.go")
	parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatalf("parse closure journey: %v", err)
	}

	ast.Inspect(parsed, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Sleep" {
			return true
		}
		identifier, ok := selector.X.(*ast.Ident)
		if ok && identifier.Name == "time" {
			t.Errorf("closure journey contains time.Sleep at %s", path)
		}
		return true
	})
}

func TestIntegrationBinaryOwnedByPackage(t *testing.T) {
	functions := topLevelFunctions(t, classificationDirectory(t))
	if !functions["TestMain"] {
		t.Fatal("dedicated integration package has no TestMain-owned binary lifecycle")
	}
	if functions["buildClosureFIDBinary"] {
		t.Fatal("per-test buildClosureFIDBinary helper remains")
	}
}

func classificationDirectory(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve classification test path")
	}
	return filepath.Dir(filename)
}

func topLevelTests(t *testing.T, directory string, excluded ...string) []string {
	t.Helper()
	excludedNames := make(map[string]bool, len(excluded))
	for _, name := range excluded {
		excludedNames[name] = true
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read integration directory: %v", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, entry.Name()), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if ok && isGoTestFunction(function) && !excludedNames[function.Name.Name] {
				names = append(names, function.Name.Name)
			}
		}
	}
	sort.Strings(names)
	return names
}

func isGoTestFunction(function *ast.FuncDecl) bool {
	if function.Recv != nil || !strings.HasPrefix(function.Name.Name, "Test") || function.Type.Params == nil {
		return false
	}
	parameters := function.Type.Params.List
	if len(parameters) != 1 || len(parameters[0].Names) != 1 {
		return false
	}
	pointer, ok := parameters[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := pointer.X.(*ast.SelectorExpr)
	return ok && selector.Sel.Name == "T"
}

func topLevelFunctions(t *testing.T, directory string) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read integration directory: %v", err)
	}
	functions := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, entry.Name()), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if ok && function.Recv == nil {
				functions[function.Name.Name] = true
			}
		}
	}
	return functions
}

func taggedTestFilesOutside(t *testing.T, repositoryRoot, dedicatedDir string) []string {
	t.Helper()
	var names []string
	err := filepath.WalkDir(repositoryRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path == dedicatedDir || (path != repositoryRoot && strings.HasPrefix(entry.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(contents), "//go:build integration\n") {
			relative, err := filepath.Rel(repositoryRoot, path)
			if err != nil {
				return err
			}
			names = append(names, filepath.ToSlash(relative))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("inventory repository integration-tagged tests: %v", err)
	}
	sort.Strings(names)
	return names
}
