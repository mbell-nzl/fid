package export

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func addLeanExportFacts(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, ".fid.toml")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\n[checks.test]\nrun=['true']\ntimeout='15m'\nartifacts=['coverage.out']\n[lanes.standard]\nchecks=[]\n"); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "coverage.out"), []byte("portable lean Artifact\x00\xff"), 0o600); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, root, "add", ".")
	fixtureGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-m", "Lean check policy")
	fixtureGit(t, root, "remote", "add", "origin", root)
	fixtureGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	fixtureGit(t, root, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	lane := "standard"
	task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Lane: &lane, Title: "Portable check", Actor: "human:mark", IdempotencyKey: "lean-task"})
	if err != nil {
		t.Fatal(err)
	}
	held, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), TaskID: task.ID, Actor: "human:mark", Reason: "Build export fixture"})
	if err != nil {
		t.Fatal(err)
	}
	options := holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: held.Fields["hold_id"].(string), Actor: "human:mark", IdempotencyKey: "lean-start"}
	if _, err := holdwork.Start(options); err != nil {
		t.Fatal(err)
	}
	if _, err := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: options.HoldID, Actor: options.Actor, Names: []string{"test"}, IdempotencyKey: "lean-check", Version: "test"}); err != nil {
		t.Fatal(err)
	}
	if _, err := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: options.HoldID, Actor: options.Actor, IdempotencyKey: "empty-check", Version: "test"}); err != nil {
		t.Fatal(err)
	}
	for _, raw := range leanJournalLines(t, root) {
		var event journal.Event
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatal(err)
		}
		switch {
		case event.EventType == "project.initialized", event.EventType == "issue.created", event.EventType == "policy.snapshotted", strings.HasPrefix(event.EventType, "task."), strings.HasPrefix(event.EventType, "check."):
		default:
			t.Fatalf("lean export fixture creates %s", event.EventType)
		}
	}
	return task.ID
}

// addLeanArtifactExportFacts returns the archive member of a genuinely kept
// check file. Generic archive attacks must never depend on kernel records.
func addLeanArtifactExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addLeanExportFacts(t, root)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	records, err := checkread.Read(conn, task)
	if err != nil || len(records) != 1 || records[0].Result == nil || len(records[0].Result.Artifacts) != 1 {
		t.Fatalf("lean fixture checks=%#v/%v", records, err)
	}
	artifact := records[0].Result.Artifacts[0]
	body, err := os.ReadFile(filepath.Join(root, "coverage.out"))
	if err != nil || artifact.Digest != digest(body) || artifact.SizeBytes != int64(len(body)) {
		t.Fatalf("lean fixture Artifact=%#v/%v", artifact, err)
	}
	return "artifacts/" + strings.TrimPrefix(artifact.Digest, "sha256:")
}

func fixtureGit(t *testing.T, workspace string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = workspace
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(arguments, " "), err, output)
	}
}
