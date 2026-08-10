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
	parentDir := filepath.Dir(dedicatedDir)
	wantJourneys := []string{
		"TestMaintenanceCheckReportsIncompatibleWorkflowAssignment",
		"TestV01ClosureJourneyAndRebuildParity",
		"TestV01ParityGoldenOutputs",
		"TestV02aWorkflowAuthoringExitScenarioAndRebuildParity",
	}

	gotJourneys := topLevelTests(t, dedicatedDir, "TestIntegrationClassification")
	if !reflect.DeepEqual(gotJourneys, wantJourneys) {
		t.Fatalf("dedicated integration journeys = %v, want %v", gotJourneys, wantJourneys)
	}

	if tagged := taggedTestFiles(t, parentDir); len(tagged) != 0 {
		t.Fatalf("integration-tagged tests outside dedicated package = %v, want none", tagged)
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
			if ok && function.Recv == nil && strings.HasPrefix(function.Name.Name, "Test") && !excludedNames[function.Name.Name] {
				names = append(names, function.Name.Name)
			}
		}
	}
	sort.Strings(names)
	return names
}

func taggedTestFiles(t *testing.T, directory string) []string {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read parent test directory: %v", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		if strings.HasPrefix(string(contents), "//go:build integration\n") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}
