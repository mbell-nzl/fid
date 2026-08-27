package main

import (
	"bytes"
	"fmt"
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

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestLiveLeafCommandsMatchContractInventories(t *testing.T) {
	retained := []string{
		"capabilities",
		"activation prepare",
		"activation component",
		"activation commit",
		"activation abort",
		"activation finalize",
		"activation status",
		"gate register",
		"gate resolve",
		"envelope create",
		"envelope get",
		"envelope activate",
		"claim acquire",
		"claim renew",
		"claim release",
		"claim list",
		"init",
		"rebuild",
		"issue create",
		"issue get",
		"issue list",
		"issue update",
		"issue transition",
		"issue ac add",
		"issue ac list",
		"issue ac verify",
		"issue label add",
		"issue label remove",
		"label add",
		"label list",
		"link add",
		"link remove",
		"link list",
		"ready",
		"workflow define",
		"workflow assign",
		"workflow list",
		"maintenance check",
		"maintenance vacuum",
		"activity list",
		"concern list",
		"concern address",
		"concern raise",
		"requirement sync",
		"stats tokens",
		"stats cost",
	}
	legacy := []string{
		"dev start",
		"dev gate",
		"dev finish",
		"dev block",
		"run record",
		"issue complete",
		"sql",
		"journal inspect",
	}

	if len(retained) != 46 || len(legacy) != 8 {
		t.Fatalf("contract inventory counts = retained %d, legacy %d; want 46 and 8", len(retained), len(legacy))
	}
	want := append(append([]string{}, retained...), legacy...)
	sort.Strings(want)
	got := liveExecutableLeafCommands(newRootCommandAt(t.TempDir()))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("live executable leaves = %v, want exact retained/legacy inventory %v", got, want)
	}

	wantCapabilities := make(map[string]rootFlagCapability, len(retained))
	for _, command := range retained {
		wantCapabilities[command] = rootFlagCapability{Project: true, JSON: true}
	}
	if !reflect.DeepEqual(rootFlagCapabilities, wantCapabilities) {
		t.Fatalf("root flag capabilities = %#v, want exact retained inventory %#v", rootFlagCapabilities, wantCapabilities)
	}

	for _, fixture := range legacyLeafInvocations() {
		for _, rootFlag := range []string{"--json", "--project"} {
			t.Run(fixture.command+" rejects "+rootFlag, func(t *testing.T) {
				assertLegacyLeafRejectsRootFlag(t, fixture, rootFlag)
			})
		}
	}
}

type legacyLeafInvocation struct {
	command string
	args    []string
}

func legacyLeafInvocations() []legacyLeafInvocation {
	return []legacyLeafInvocation{
		{command: "dev start", args: []string{"dev", "start", "FID-legacy"}},
		{command: "dev gate", args: []string{"dev", "gate", "FID-legacy"}},
		{command: "dev finish", args: []string{"dev", "finish", "FID-legacy"}},
		{command: "dev block", args: []string{"dev", "block", "FID-legacy"}},
		{command: "run record", args: []string{"run", "record"}},
		{command: "issue complete", args: []string{"issue", "complete", "FID-legacy"}},
		{command: "sql", args: []string{"sql", "SELECT 1"}},
		{command: "journal inspect", args: []string{"journal", "inspect", "EVT-legacy"}},
	}
}

func assertLegacyLeafRejectsRootFlag(t *testing.T, fixture legacyLeafInvocation, rootFlag string) {
	t.Helper()
	args := append([]string{}, fixture.args...)
	if rootFlag == "--project" {
		args = append([]string{rootFlag, t.TempDir()}, args...)
	} else {
		args = append([]string{rootFlag}, args...)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(t.TempDir())
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("fid %s with %s succeeded, want pre-handler rejection", fixture.command, rootFlag)
	}
	if stdout.Len() != 0 {
		t.Fatalf("fid %s with %s stdout = %q, want empty", fixture.command, rootFlag, stdout.String())
	}
	want := "error: " + rootFlag + " is not available for " + fixture.command
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("fid %s with %s stderr = %q, want %q", fixture.command, rootFlag, stderr.String(), want)
	}
}

func liveExecutableLeafCommands(root *cobra.Command) []string {
	var leaves []string
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if command != root && (command.Run != nil || command.RunE != nil) {
			leaves = append(leaves, rootCommandKey(command))
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	sort.Strings(leaves)
	return leaves
}

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

func TestProductSourceContainsNoUniversalSerializers(t *testing.T) {
	repoRoot := commandTestRepositoryRoot(t)
	violations, err := universalSerializerSourceViolations(repoRoot)
	if err != nil {
		t.Fatalf("universalSerializerSourceViolations() error = %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("product source universal serializer violations = %v, want none", violations)
	}
}

func TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags(t *testing.T) {
	t.Run("inherited persistent flag", func(t *testing.T) {
		root := newRootCommandAt(t.TempDir())
		root.PersistentFlags().String("format", "", "injected test flag")

		violations := liveUniversalFormatFlagViolations(root)
		if !containsUniversalFormatViolation(violations, "fid issue --format") {
			t.Fatalf("injected inherited format flag violations = %v, want descendant fid issue --format", violations)
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

func universalSerializerSourceViolations(repoRoot string) ([]string, error) {
	formatterImportSuffix := "/internal/" + "format"
	formatterDirectory := filepath.Join(repoRoot, "internal", "format")
	serializerNames := map[string]struct{}{
		"write" + "TOON":     {},
		"write" + "YAML":     {},
		"write" + "Markdown": {},
	}
	yamlPackage := "ya" + "ml"
	yamlMarshal := "Mar" + "shal"
	violations := []string{}
	if _, err := os.Stat(formatterDirectory); err == nil {
		violations = append(violations, formatterDirectory)
	}
	for _, directory := range []string{"cmd", "internal"} {
		root := filepath.Join(repoRoot, directory)
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				return err
			}
			for _, imported := range parsed.Imports {
				pathValue := strings.Trim(imported.Path.Value, "\"")
				if strings.HasSuffix(pathValue, formatterImportSuffix) || pathValue == "gopkg.in/yaml.v3" {
					violations = append(violations, path+":import "+pathValue)
				}
			}
			ast.Inspect(parsed, func(node ast.Node) bool {
				switch typed := node.(type) {
				case *ast.Ident:
					if _, forbidden := serializerNames[typed.Name]; forbidden {
						violations = append(violations, path+":"+typed.Name)
					}
				case *ast.SelectorExpr:
					if packageName, ok := typed.X.(*ast.Ident); ok && packageName.Name == yamlPackage && typed.Sel.Name == yamlMarshal {
						violations = append(violations, path+":"+yamlPackage+"."+yamlMarshal)
					}
				}
				return true
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(violations)
	return violations, nil
}

func commandTestRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}
