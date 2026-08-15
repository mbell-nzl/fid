package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func TestLiveTestsContainNoUniversalFormatInvocation(t *testing.T) {
	repoRoot := commandTestRepositoryRoot(t)
	allowed := map[string]bool{
		"TestIssueGetRejectsRetiredFormatFlags":                                     true,
		"TestIssueGetRejectsAllFormatValues":                                        true,
		"TestIssueListRejectsRetiredFormatFlag":                                     true,
		"TestLabelListRejectsLegacyFormatFlag":                                      true,
		"TestReadyRejectsRetiredFormatFlag":                                         true,
		"TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues": true,
		"TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat":               true,
		"TestV01ReadCommandsRejectRetiredFormats":                                   true,
		"TestIssueReadRejectsRetiredAndMalformedInputs":                             true,
	}
	var violations []string
	term := "--for" + "mat"
	commandTests := filepath.Join(repoRoot, "cmd", "fid")
	if err := filepath.WalkDir(commandTests, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" || filepath.Base(path) == "contract_surface_test.go" {
			return nil
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Body == nil || allowed[function.Name.Name] {
				continue
			}
			ast.Inspect(function.Body, func(node ast.Node) bool {
				literal, ok := node.(*ast.BasicLit)
				if ok && literal.Kind == token.STRING && literal.Value == `"`+term+`"` {
					violations = append(violations, path+":"+function.Name.Name)
				}
				return true
			})
		}
		return nil
	}); err != nil {
		t.Fatalf("WalkDir(%q) error = %v", commandTests, err)
	}
	sort.Strings(violations)
	if len(violations) > 0 {
		t.Fatalf("live test format invocations = %v, want only retained rejection assertions", violations)
	}
}

func commandTestRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}
