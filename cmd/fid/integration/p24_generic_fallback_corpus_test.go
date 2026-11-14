//go:build integration

package integration

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

type p24FallbackRow struct {
	site, operation, argument, allowed, message string
	args                                        []string
	project                                     bool
}

// Each site is the final R20 mapper, not a grammar or typed invalid-argument.
// Missing-project rows reach real project-domain failures with valid grammar.
var p24FallbackRows = []p24FallbackRow{
	{"artifact.go:renderArtifactError", "artifact.verify", "input", "accepted gate registration input", "verify Artifact: Artifact \"ART-missing\" was not found", []string{"artifact", "verify", "--actor", "agent:builder", "--idempotency-key", "fallback"}, true},
	{"issue_mutation_response.go:issueMutationFailure", "issue.create", "request", "accepted issue.create input", "create issue: epic scheduling is required", []string{"issue", "create", "--kind", "epic", "--title", "test", "--actor", "agent:a", "--idempotency-key", "fallback"}, true},
	{"link_response.go:linkFailure", "link.list", "request", "accepted link.list input", "link list: entity id is required", []string{"link", "list", " ", "--type", "issue"}, true},
	{"issue_get_list.go:issueReadFailure", "activity.list", "request", "accepted issue read input", "activity list: --since must be YYYY-MM-DD or RFC3339", []string{"activity", "list", "--since", "invalid"}, true},
}

// Empty: all current final fallbacks have reachable process rows. An
// exclusion must name a specific site and explain why the CLI cannot reach it.
var p24FallbackExclusions = map[string]string{}

func TestP24GenericFallbackCorpusCompleteness(t *testing.T) {
	t.Parallel()
	sites := map[string]bool{}
	files, err := filepath.Glob(filepath.Join(integrationRepositoryRoot(), "cmd/fid/*.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				name, ok := call.Fun.(*ast.Ident)
				if !ok || (name.Name != "genericFailure" && name.Name != "writeGenericFailure") {
					return true
				}
				for _, arg := range call.Args {
					literal, ok := arg.(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						continue
					}
					value, err := strconv.Unquote(literal.Value)
					if err != nil {
						t.Fatal(err)
					}
					if strings.HasPrefix(value, "invalid ") {
						sites[filepath.Base(path)+":"+fn.Name.Name] = true
					}
				}
				return true
			})
		}
	}
	seen := map[string]bool{}
	for _, row := range p24FallbackRows {
		if seen[row.site] || !sites[row.site] {
			t.Fatalf("duplicate or stale row %s", row.site)
		}
		seen[row.site] = true
	}
	for site, reason := range p24FallbackExclusions {
		if seen[site] || !sites[site] || strings.TrimSpace(reason) == "" {
			t.Fatalf("invalid exclusion %s: %s", site, reason)
		}
		seen[site] = true
	}
	for site := range sites {
		if !seen[site] {
			t.Errorf("R20 fallback has no process row or pinned exclusion: %s", site)
		}
	}
}

func TestP24GenericFallbackDomainPairs(t *testing.T) {
	t.Parallel()
	for _, row := range p24FallbackRows {
		t.Run(row.site, func(t *testing.T) {
			t.Parallel()
			binary, root := integrationFIDBinary(t), gittest.NewWorkspace(t)
			if row.project {
				p21Run(t, binary, root, "init", "--prefix", "P24", "--json")
			}
			args := append([]string{"--project", root}, row.args...)
			if row.operation == "artifact.verify" {
				args = append(args, "--input", p21Input(t, root, "missing-artifact", map[string]any{"artifact_id": "ART-missing"}))
			}
			before := p24FailureFilesystem(t, root)
			out, stderr, err := closureRunCommand(root, binary, append([]string{"--json"}, args...)...)
			if err == nil || processExitCode(t, err) != 2 || out != "" {
				t.Fatalf("domain refusal: %v %q %q", err, out, stderr)
			}
			payload := p23Envelope(t, stderr)
			f := payload["error"].(map[string]any)
			message := strings.ReplaceAll(row.message, "$ROOT", gittest.ResolvedPath(t, root))
			var value any
			if row.site == "issue_get_list.go:issueReadFailure" {
				value = ""
			}
			details := map[string]any{"argument": row.argument, "value": value, "allowed": []any{row.allowed}}
			if payload["operation"] != row.operation || f["code"] != "invalid-argument" || f["category"] != "invalid-input" || f["retryable"] != false || f["subject"] != nil || f["message"] != message || !reflect.DeepEqual(f["details"], details) {
				t.Fatalf("did not reach pinned R20 domain fallback: %#v", payload)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, root)) {
				t.Fatal("JSON refusal changed filesystem")
			}
			out, stderr, err = closureRunCommand(root, binary, args...)
			want := "error: " + message + "\nreason: invalid-argument; argument=" + row.argument + ", allowed=" + row.allowed + "\n"
			if row.operation == "artifact.verify" {
				want += "next: see fid artifact verify --help; template: fid artifact verify --example\n"
			}
			if err == nil || processExitCode(t, err) != 2 || out != "" || stderr != want {
				t.Fatalf("Human domain refusal: %v %q %q want %q", err, out, stderr, want)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, root)) {
				t.Fatal("Human refusal changed filesystem")
			}
		})
	}
}
