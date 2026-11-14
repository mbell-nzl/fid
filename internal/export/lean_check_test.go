package export

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
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
	return task.ID
}

func TestLeanCheckExportImportPreservesRecordsAndArtifact(t *testing.T) {
	t.Parallel()
	var task string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string { task = addLeanExportFacts(t, root); return task })
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "lean-import"}); err != nil {
		t.Fatal(err)
	}
	var original string
	for _, root := range []string{fixture.workspace, target} {
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
		if err != nil {
			t.Fatal(err)
		}
		records, err := checkread.Read(conn, task)
		_ = conn.Close()
		if err != nil || len(records) != 1 || records[0].Result == nil || !records[0].Result.Counts || len(records[0].Result.Artifacts) != 1 {
			t.Fatalf("records=%#v/%v", records, err)
		}
		raw, err := json.Marshal(records)
		if err != nil {
			t.Fatal(err)
		}
		if original == "" {
			original = string(raw)
		} else if string(raw) != original {
			t.Fatalf("imported records differ: %s", raw)
		}
		artifact := records[0].Result.Artifacts[0]
		cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(root), config.ConfigFileName))
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(config.LedgerDir(root), cfg.ArtifactRoot, strings.TrimPrefix(artifact.Digest, "sha256:")))
		if err != nil || string(data) != "portable lean Artifact\x00\xff" {
			t.Fatalf("artifact=%v/%v", data, err)
		}
	}
}

func leanJournalLines(t *testing.T, root string) [][]byte {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(root), "journal"))
	if err != nil {
		t.Fatal(err)
	}
	lines := [][]byte{}
	for _, entry := range entries {
		file, err := os.Open(filepath.Join(config.LedgerDir(root), "journal", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 4096), 1<<20)
		for scanner.Scan() {
			lines = append(lines, append([]byte{}, scanner.Bytes()...))
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return lines
}

func TestLeanCheckReplayAndImportRejectForgedObservations(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanExportFacts)
	lines := leanJournalLines(t, fixture.workspace)
	for _, row := range []struct {
		name, family string
		mutate       func(*journal.Event)
	}{
		{"start digest", "check.started", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "check_starts" {
					var value map[string]any
					_ = json.Unmarshal([]byte(e.Effects[i].Values["start_json"].(string)), &value)
					value["definition_digest"] = "sha256:" + strings.Repeat("0", 64)
					e.Effects[i].Values["start_json"] = encodedLean(value)
				}
			}
		}},
		{"process identity", "check.process-started", func(e *journal.Event) {
			e.Effects[0].Values["process_json"] = `{"machine":"local","pid":0,"start_time":"0"}`
		}},
		{"result qualification", "check.result-recorded", func(e *journal.Event) {
			var value map[string]any
			_ = json.Unmarshal([]byte(e.Effects[0].Values["result_json"].(string)), &value)
			value["counts"] = false
			e.Effects[0].Values["result_json"] = encodedLean(value)
		}},
		{"artifact binding", "check.artifact-kept", func(e *journal.Event) { e.Effects[len(e.Effects)-1].PK["path"] = "not-listed" }},
		{"empty request lane", "check.completed", func(e *journal.Event) {
			var value map[string]any
			_ = json.Unmarshal([]byte(e.Effects[0].Values["response_json"].(string)), &value)
			value["lane"] = "assured"
			e.Effects[0].Values["response_json"] = encodedLean(value)
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), 25)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			for _, raw := range lines {
				var event journal.Event
				if err := json.Unmarshal(raw, &event); err != nil {
					t.Fatal(err)
				}
				if event.EventType == row.family {
					row.mutate(&event)
					raw, err = json.Marshal(event)
					if err != nil {
						t.Fatal(err)
					}
					if err := validator.accept(raw); err == nil {
						t.Fatal("import accepted forged observation")
					}
					if err := journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: "sha256:forged", File: "journal", Line: 1, Offset: 1}); err == nil {
						t.Fatal("replay accepted forged observation")
					}
					return
				}
				if err := validator.accept(raw); err != nil {
					t.Fatal(err)
				}
			}
			t.Fatal("fixture is missing family")
		})
	}
}

func encodedLean(value any) string { raw, _ := json.Marshal(value); return string(raw) }

func TestLeanCheckProjectionUpgradeRebuildsTheOldCache(t *testing.T) {
	t.Parallel()
	var task string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string { task = addLeanExportFacts(t, root); return task })
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{`DROP TABLE check_files`, `DROP TABLE check_results`, `DROP TABLE check_processes`, `DROP TABLE check_starts`, `DROP TABLE check_requests`, `UPDATE schema_migrations SET version=24`, `UPDATE projection_meta SET projection_schema_version=6`} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	result, err := journal.EnsureProjectionFresh(fixture.workspace, testMigrationsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	if result.Decision != fiddb.ProjectionRebuild {
		t.Fatalf("upgrade=%s, want a journal rebuild", result.Decision)
	}
	conn, err = fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	records, err := checkread.Read(conn, task)
	if err != nil || len(records) != 1 || records[0].Result == nil {
		t.Fatalf("rebuilt checks=%#v/%v", records, err)
	}
}
