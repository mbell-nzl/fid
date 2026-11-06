package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// TestEveryProjectCommandRefusesOutsideAGitRepository runs every command that
// opens or creates a project in a directory that holds a ledger but is in no
// git repository. Each one refuses before touching the ledger.
func TestEveryProjectCommandRefusesOutsideAGitRepository(t *testing.T) {
	t.Parallel()
	outside := clonePristineWorkspace(t)
	if err := os.RemoveAll(filepath.Join(outside, ".git")); err != nil {
		t.Fatalf("remove .git: %v", err)
	}
	inputs := t.TempDir()
	cases := projectCommandRefusalCases(t, outside, inputs)
	if len(cases) != 93 {
		t.Fatalf("project commands = %d, want every live leaf except export validate and maintenance vacuum (93)", len(cases))
	}
	before := snapshotWorkspacePersistence(t, outside)
	for _, item := range cases {
		modes := []bool{false}
		if item.json {
			modes = append(modes, true)
		}
		for _, jsonMode := range modes {
			args := item.args
			if jsonMode {
				args = append([]string{"--json"}, args...)
			}
			result := grammarExecute(outside, args...)
			if cli.ExitCode(result.err) != 4 || result.stdout != "" {
				t.Errorf("%s (json %t): exit %d, result %#v; want refusal with exit 4", item.leaf, jsonMode, cli.ExitCode(result.err), result)
				continue
			}
			if !jsonMode {
				if !strings.HasPrefix(result.stderr, "error: fid needs a git repository\nreason: repository-unavailable; detail=") {
					t.Errorf("%s Human refusal = %q", item.leaf, result.stderr)
				}
				continue
			}
			assertNotARepositoryFailure(t, item.operation, outside, result.stderr)
		}
	}
	if !before.equal(snapshotWorkspacePersistence(t, outside)) {
		t.Fatal("a refused command changed the ledger")
	}
}

type projectCommandRefusalCase struct {
	leaf, operation string
	json            bool
	args            []string
}

// projectCommandRefusalCases builds, for every project command, arguments
// that pass grammar and input checks, so the command reaches discovery.
func projectCommandRefusalCases(t *testing.T, outside, inputs string) []projectCommandRefusalCase {
	t.Helper()
	root := newRootCommandAt(outside)
	leaves, _ := grammarCommands(root)
	for _, path := range [][]string{{"journal", "inspect"}, {"sql"}} {
		leaf, _, err := root.Find(path)
		if err != nil {
			t.Fatalf("find %v: %v", path, err)
		}
		leaves = append(leaves, leaf)
	}
	cases := []projectCommandRefusalCase{}
	for _, leaf := range leaves {
		key := rootCommandKey(leaf)
		if key == "export validate" || key == "maintenance vacuum" {
			continue // They touch no project.
		}
		args := append(strings.Fields(key), projectCommandPositionals(key, inputs, len(strings.Fields(leaf.Use))-1)...)
		if leaf.Flags().Lookup("input") != nil {
			path := filepath.Join(inputs, strings.ReplaceAll(key, " ", "-")+".json")
			writeActivationInput(t, path, projectCommandInput(t, key))
			args = append(args, "--input", path)
		}
		for _, name := range []string{"actor", "idempotency-key"} {
			if leaf.Flags().Lookup(name) != nil {
				args = append(args, "--"+name, grammarValidFlag(name))
			}
		}
		args = append(args, projectCommandFlags(key, outside)...)
		cases = append(cases, projectCommandRefusalCase{leaf: key, operation: leafOperation(leaf), json: !diagnosticLeaf(leaf), args: args})
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].leaf < cases[j].leaf })
	return cases
}

func projectCommandPositionals(key, inputs string, count int) []string {
	switch key {
	case "sql":
		return []string{"SELECT 1"}
	case "export create", "export import":
		return []string{filepath.Join(inputs, "archive.zip")}
	}
	values := []string{}
	for range count {
		values = append(values, "ID-1")
	}
	return values
}

// projectCommandFlags supplies the required flags of leaves without --input.
func projectCommandFlags(key, outside string) []string {
	switch key {
	case "issue create":
		return []string{"--kind", "task", "--title", "refusal"}
	case "issue transition":
		return []string{"--to", "done"}
	case "issue update":
		return []string{"--title", "refusal"}
	case "issue ac add":
		return []string{"--text", "refusal", "--ordinal", "1"}
	case "link add", "link remove":
		return []string{"--from-type", "issue", "--to-type", "issue"}
	case "link list":
		return []string{"--type", "issue"}
	case "gate resolve":
		return []string{"--version", "1"}
	case "export import":
		return []string{"--repository", outside}
	case "workflow define":
		return []string{"--name", "refusal", "--applies-to", "task", "--states", "todo,done", "--categories", "unstarted,completed", "--transitions", "todo->done"}
	case "workflow assign":
		return []string{"--name", "refusal", "--kind", "task"}
	}
	return nil
}

// projectCommandInput returns the leaf's --example template with its
// placeholders made structurally valid.
func projectCommandInput(t *testing.T, key string) map[string]any {
	t.Helper()
	if key == "activation prepare" {
		return activationPrepareInput("kit-refusal")
	}
	descriptor, ok := inputDescriptorRegistry[key]
	if !ok {
		t.Fatalf("%s has --input but no descriptor", key)
	}
	data, err := cli.InputExample(descriptor)
	if err != nil {
		t.Fatalf("%s example: %v", key, err)
	}
	var input map[string]any
	if err := json.Unmarshal(data, &input); err != nil {
		t.Fatalf("%s example: %v", key, err)
	}
	digest := "sha256:" + strings.Repeat("a", 64)
	var fix func(any) any
	fix = func(value any) any {
		switch typed := value.(type) {
		case map[string]any:
			for name, member := range typed {
				typed[name] = fix(member)
			}
			if _, ok := typed["component_digests"]; ok {
				typed["component_digests"] = []any{digest, "sha256:" + strings.Repeat("b", 64)}
			}
		case []any:
			for index, member := range typed {
				typed[index] = fix(member)
			}
		case string:
			switch typed {
			case "<head-oid>":
				return strings.Repeat("a", 40)
			case "<sha256-digest>":
				return digest
			case "<timestamp>":
				return "2026-11-04T00:00:00Z"
			}
		}
		return value
	}
	return fix(input).(map[string]any)
}

func assertNotARepositoryFailure(t *testing.T, operation, start, stderr string) {
	t.Helper()
	var envelope struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code      string         `json:"code"`
			Category  string         `json:"category"`
			Message   string         `json:"message"`
			Retryable bool           `json:"retryable"`
			Subject   any            `json:"subject"`
			Details   map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(stderr), &envelope); err != nil {
		t.Errorf("%s refusal: %v\nstderr=%q", operation, err, stderr)
		return
	}
	failure := envelope.Error
	detail, _ := failure.Details["detail"].(string)
	if envelope.OK || envelope.Operation != operation || failure.Code != "repository-unavailable" || failure.Category != "precondition" ||
		failure.Message != "fid needs a git repository" || !failure.Retryable || failure.Subject != nil || len(failure.Details) != 2 ||
		failure.Details["reason"] != "not-a-repository" || !strings.Contains(detail, start) {
		t.Errorf("%s refusal = %s, want repository-unavailable not-a-repository naming %q", operation, stderr, start)
	}
}

func TestCommandsThatTouchNoProjectWorkOutsideAGitRepository(t *testing.T) {
	t.Parallel()
	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(clonePristineWorkspace(t), "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	outside := t.TempDir()
	for _, args := range [][]string{
		{"--help"},
		{"issue", "create", "--help"},
		{"completion", "bash"},
		{"claim", "acquire", "--example"},
		{"--json", "export", "validate", archive},
	} {
		result := grammarExecute(outside, args...)
		if result.err != nil || result.stdout == "" || result.stderr != "" {
			t.Fatalf("fid %s outside a repository = %#v, want success", strings.Join(args, " "), result)
		}
	}
	vacuum := grammarExecute(outside, "--json", "maintenance", "vacuum")
	if cli.ExitCode(vacuum.err) != 5 || !strings.Contains(vacuum.stderr, `"code":"capability-unavailable"`) {
		t.Fatalf("maintenance vacuum outside a repository = %#v, want its unchanged capability-unavailable", vacuum)
	}
}

// TestProjectCommandsUseTheWorktreeTopLevel covers discovery from a
// subdirectory, --project as git -C, and the three commands that read the
// ledger without resolving a project.
func TestProjectCommandsUseTheWorktreeTopLevel(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	nested := filepath.Join(workspace, "nested", "deeper")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll(nested) error = %v", err)
	}
	created := executeCommandForTest(nested, "--json", "issue", "create", "--kind", "task", "--title", "From a subdirectory", "--actor", "agent:developer")
	if created.err != nil {
		t.Fatalf("issue create from a subdirectory = %#v", created)
	}
	var envelope struct {
		Result struct {
			Issue struct{ ID string } `json:"issue"`
		} `json:"result"`
		Event struct{ ID string } `json:"event"`
	}
	if err := json.Unmarshal([]byte(created.stdout), &envelope); err != nil {
		t.Fatalf("issue create response: %v\n%s", err, created.stdout)
	}
	issueID := envelope.Result.Issue.ID
	outside := t.TempDir()
	parent, name := filepath.Split(workspace)
	for _, item := range []struct {
		name, directory string
		args            []string
	}{
		{"subdirectory", nested, nil},
		{"--project subdirectory", outside, []string{"--project", nested}},
		{"relative --project", parent, []string{"--project", filepath.Join(name, "nested")}},
	} {
		listed := executeCommandForTest(item.directory, append(item.args, "--json", "issue", "list")...)
		if listed.err != nil || !strings.Contains(listed.stdout, `"id":"`+issueID+`"`) {
			t.Fatalf("%s: issue list = %#v, want %s", item.name, listed, issueID)
		}
		concerns := executeCommandForTest(item.directory, append(item.args, "--json", "concern", "list")...)
		if concerns.err != nil || !strings.Contains(concerns.stdout, `"operation":"concern.list"`) {
			t.Fatalf("%s: concern list = %#v", item.name, concerns)
		}
	}
	inspected := executeCommandForTest(nested, "journal", "inspect", envelope.Event.ID)
	if inspected.err != nil || !strings.Contains(inspected.stdout, `"event_id":"`+envelope.Event.ID+`"`) {
		t.Fatalf("journal inspect from a subdirectory = %#v", inspected)
	}
	queried := executeCommandForTest(nested, "sql", "SELECT id FROM issues")
	if queried.err != nil || !strings.Contains(queried.stdout, issueID) {
		t.Fatalf("sql from a subdirectory = %#v", queried)
	}
}

func TestALedgerBelowTheTopLevelIsNotAProject(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	nested := filepath.Join(workspace, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	clonePristineWorkspaceInto(t, nested)
	if err := os.RemoveAll(filepath.Join(nested, ".git")); err != nil {
		t.Fatalf("remove nested .git: %v", err)
	}
	for _, directory := range []string{workspace, nested} {
		result := executeCommandForTest(directory, "--json", "issue", "list")
		if cli.ExitCode(result.err) != 3 || !strings.Contains(result.stderr, `"code":"project-not-found"`) {
			t.Fatalf("issue list from %q = %#v, want project-not-found", directory, result)
		}
	}
}

func TestInitFromASubdirectoryInitializesTheTopLevel(t *testing.T) {
	t.Parallel()
	for _, viaProject := range []bool{false, true} {
		workspace := gittest.NewNamedWorkspace(t, "toplevel")
		nested := filepath.Join(workspace, "nested")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatalf("Mkdir(nested) error = %v", err)
		}
		directory, args := nested, []string{"--json", "init"}
		if viaProject {
			directory, args = t.TempDir(), []string{"--project", nested, "--json", "init"}
		}
		result := grammarExecute(directory, args...)
		if result.err != nil {
			t.Fatalf("init from a subdirectory (--project %t) = %#v", viaProject, result)
		}
		project := decodeCanonicalResult(t, result.stdout)["project"].(map[string]any)
		if want := gittest.ResolvedPath(t, workspace); project["root"] != want || project["issue_prefix"] != "TOPLEVEL" {
			t.Fatalf("init from a subdirectory (--project %t) project = %#v, want root %q and prefix TOPLEVEL", viaProject, project, want)
		}
		if _, err := os.Stat(filepath.Join(config.LedgerDir(workspace), "config.toml")); err != nil {
			t.Fatalf("top-level ledger: %v", err)
		}
		if _, err := os.Stat(config.LedgerDir(nested)); !os.IsNotExist(err) {
			t.Fatalf("subdirectory ledger = %v, want none", err)
		}
	}
}

func TestInitAndImportReportTheRealPathOfASymlinkedRepository(t *testing.T) {
	t.Parallel()
	links := t.TempDir()
	workspace := gittest.NewWorkspace(t)
	alias := filepath.Join(links, "workspace")
	if err := os.Symlink(workspace, alias); err != nil {
		t.Fatalf("Symlink(workspace) error = %v", err)
	}
	initialized := grammarExecute(t.TempDir(), "--project", alias, "--json", "init", "--prefix", "FID")
	if initialized.err != nil {
		t.Fatalf("init through a symlink = %#v", initialized)
	}
	if root, want := decodeCanonicalResult(t, initialized.stdout)["project"].(map[string]any)["root"], gittest.ResolvedPath(t, workspace); root != want {
		t.Fatalf("init project.root = %v, want real path %q", root, want)
	}

	archive := filepath.Join(t.TempDir(), "fid-export.zip")
	if created := executeCommandForTest(workspace, "--json", "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	for _, through := range []string{"symlink", "subdirectory"} {
		target := gittest.NewWorkspace(t)
		repository := filepath.Join(links, "target")
		if through == "symlink" {
			if err := os.Symlink(target, repository); err != nil {
				t.Fatalf("Symlink(target) error = %v", err)
			}
		} else {
			repository = filepath.Join(target, "nested")
			if err := os.Mkdir(repository, 0o755); err != nil {
				t.Fatalf("Mkdir(nested) error = %v", err)
			}
		}
		imported := grammarExecute(t.TempDir(), "--json", "export", "import", archive, "--repository", repository, "--idempotency-key", "import-"+through)
		if imported.err != nil {
			t.Fatalf("export import through a %s = %#v", through, imported)
		}
		if root, want := decodeCanonicalResult(t, imported.stdout)["import"].(map[string]any)["target_root"], gittest.ResolvedPath(t, target); root != want {
			t.Fatalf("import through a %s target_root = %v, want top level %q", through, root, want)
		}
		if _, err := os.Stat(filepath.Join(config.LedgerDir(target), "config.toml")); err != nil {
			t.Fatalf("imported ledger through a %s: %v", through, err)
		}
	}
}
