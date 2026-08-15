package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

func TestLiveCobraTreeContainsNoUniversalFormatFlags(t *testing.T) {
	root := newRootCommandAt(t.TempDir())
	if violations := liveUniversalFormatFlagViolations(root); len(violations) != 0 {
		t.Fatalf("live Cobra tree universal format flags = %v, want none", violations)
	}
}

func TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags(t *testing.T) {
	t.Run("inherited persistent flag", func(t *testing.T) {
		root := newRootCommandAt(t.TempDir())
		root.PersistentFlags().String("format", "", "injected test flag")

		violations := liveUniversalFormatFlagViolations(root)
		if !containsUniversalFormatViolation(violations, "fid --format") {
			t.Fatalf("injected inherited format flag violations = %v, want fid --format", violations)
		}
	})

	t.Run("local flag", func(t *testing.T) {
		root := newRootCommandAt(t.TempDir())
		var ready *cobra.Command
		for _, command := range root.Commands() {
			if command.Name() == "ready" {
				ready = command
				break
			}
		}
		if ready == nil {
			t.Fatal("live Cobra tree has no ready command")
		}
		ready.Flags().String("yaml", "", "injected test flag")

		violations := liveUniversalFormatFlagViolations(root)
		if !containsUniversalFormatViolation(violations, "fid ready --yaml") {
			t.Fatalf("injected local yaml flag violations = %v, want fid ready --yaml", violations)
		}
	})
}

var universalFormatFlagNames = map[string]struct{}{
	"format":   {},
	"toon":     {},
	"yaml":     {},
	"markdown": {},
}

func liveUniversalFormatFlagViolations(root *cobra.Command) []string {
	violations := []string{}
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		inspect := func(flags *pflag.FlagSet) {
			flags.VisitAll(func(flag *pflag.Flag) {
				if _, universal := universalFormatFlagNames[flag.Name]; universal {
					violations = append(violations, fmt.Sprintf("%s --%s", command.CommandPath(), flag.Name))
				}
			})
		}
		inspect(command.LocalNonPersistentFlags())
		inspect(command.PersistentFlags())
		inspect(command.InheritedFlags())
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	sort.Strings(violations)
	return violations
}

func containsUniversalFormatViolation(violations []string, want string) bool {
	for _, violation := range violations {
		if strings.HasPrefix(violation, want) {
			return true
		}
	}
	return false
}

func commandTestRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}
