package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestInitProjectIdentitySurvivesClone(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	stdout, _ := runInitCommand(t, workspace, "--json", "init", "--prefix", "FID")
	var envelope struct {
		Result struct {
			Project             struct{ ID string } `json:"project"`
			ProjectIdentityPath string              `json:"project_identity_path"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	id := envelope.Result.Project.ID
	if !regexp.MustCompile(`^project-[0-9a-f]{32}$`).MatchString(id) {
		t.Fatalf("project ID = %q, want project- plus 128 random bits in hex", id)
	}
	path := filepath.Join(workspace, ".fid.toml")
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "[project]\nid = \"" + id + "\"\n"; string(payload) != want {
		t.Fatalf(".fid.toml = %q, want %q", payload, want)
	}
	if want := filepath.Join(gittest.ResolvedPath(t, workspace), ".fid.toml"); envelope.Result.ProjectIdentityPath != want {
		t.Fatalf("project_identity_path = %q, want %q", envelope.Result.ProjectIdentityPath, want)
	}
	if tracked := identityGit(t, workspace, "ls-files"); tracked != "" {
		t.Fatalf("init staged or committed files: %q", tracked)
	}
	payload = append([]byte("# project identity; preserve this comment\n"), payload...)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	identityGit(t, workspace, "add", ".fid.toml")
	identityGit(t, workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false", "commit", "-q", "--no-verify", "-m", "project identity")
	clone := filepath.Join(t.TempDir(), "clone")
	identityGit(t, workspace, "clone", "-q", "--no-hardlinks", workspace, clone)
	if err := gittest.CheckRepository(clone); err != nil {
		t.Fatal(err)
	}
	text, _ := runInitCommand(t, clone, "init", "--prefix", "FID")
	if !strings.Contains(text, id) || !strings.Contains(text, ".fid.toml") || !strings.Contains(text, "commit") {
		t.Fatalf("text output = %q, want file, ID and commit instruction", text)
	}
	if got, err := os.ReadFile(filepath.Join(clone, ".fid.toml")); err != nil || string(got) != string(payload) {
		t.Fatalf("clone identity = %q, %v, want unchanged %q", got, err, payload)
	}
	project, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(clone), "config.toml"))
	if err != nil || project.ProjectID != id {
		t.Fatalf("clone ledger identity = %q, %v, want %q", project.ProjectID, err, id)
	}
	if got := identityGit(t, clone, "status", "--porcelain"); got != "" {
		t.Fatalf("clone changed committed identity: %q", got)
	}
}

func TestInitWritesIdentityAtRepositoryRoot(t *testing.T) {
	t.Parallel()
	var ids []string
	for range 2 {
		workspace := gittest.NewNamedWorkspace(t, "same-name")
		subdir := filepath.Join(workspace, "subdir")
		if err := os.Mkdir(subdir, 0o755); err != nil {
			t.Fatal(err)
		}
		runInitCommand(t, subdir, "init", "--prefix", "FID")
		id, err := config.ReadProjectIdentity(workspace)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
		if _, err := os.Lstat(filepath.Join(subdir, ".fid.toml")); !os.IsNotExist(err) {
			t.Fatalf("subdirectory identity = %v, want absent", err)
		}
	}
	if ids[0] == ids[1] {
		t.Fatalf("independent repositories with the same name got identical IDs: %q", ids[0])
	}
}

func TestInitRefusesInvalidProjectIdentityWithoutWrites(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, contents string
		directory      bool
	}{
		{name: "missing ID", contents: "[project]\n"},
		{name: "empty ID", contents: "[project]\nid = \"\"\n"},
		{name: "unsafe ID", contents: "[project]\nid = \"../other\"\n"},
		{name: "invalid TOML", contents: "[project\nid = 12\n"},
		{name: "duplicate ID", contents: "[project]\nid = \"project-0123456789abcdef0123456789abcdef\"\nid = \"project-0123456789abcdef0123456789abcdef\"\n"},
		{name: "unknown key", contents: "[project]\nid = \"project-0123456789abcdef0123456789abcdef\"\nother = \"x\"\n"},
		{name: "unreadable file type", directory: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace := gittest.NewWorkspace(t)
			path := filepath.Join(workspace, ".fid.toml")
			if test.directory {
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(path, []byte(test.contents), 0o644); err != nil {
				t.Fatal(err)
			}
			stdout, stderr, err := executeInitCommand(t, "--project", workspace, "--json", "init", "--prefix", "FID")
			if err == nil || stdout != "" || !strings.Contains(stderr, ".fid.toml") {
				t.Fatalf("init = %q, %q, %v; want clear .fid.toml refusal", stdout, stderr, err)
			}
			if cli.ExitCode(err) != 2 {
				t.Errorf("exit = %d, want invalid-input exit 2", cli.ExitCode(err))
			}
			assertProjectInitFailureCode(t, stderr, "invalid-argument")
			var failure struct {
				Error struct {
					Category string
					Details  struct {
						Argument string
						Value    *string
						Allowed  []string
					}
				}
			}
			if err := json.Unmarshal([]byte(stderr), &failure); err != nil {
				t.Fatal(err)
			}
			if failure.Error.Category != "invalid-input" || failure.Error.Details.Argument != filepath.Join(gittest.ResolvedPath(t, workspace), ".fid.toml") || failure.Error.Details.Value != nil || len(failure.Error.Details.Allowed) == 0 {
				t.Errorf("invalid config failure = %+v", failure.Error)
			}
			if _, err := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(err) {
				t.Fatalf("ledger state = %v, want absent", err)
			}
			if !test.directory {
				got, err := os.ReadFile(path)
				if err != nil || string(got) != test.contents {
					t.Fatalf(".fid.toml changed: %q, %v", got, err)
				}
			}
		})
	}
}

func identityGit(t *testing.T, workspace string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
	return string(output)
}
