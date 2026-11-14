package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergework"
)

func addMergeExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addLeanExportFacts(t, root)
	path := filepath.Join(root, ".fid.toml")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\n[review]\ntest_adequacy=[]\n"); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, root, "add", ".")
	fixtureGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Configure review")
	fixtureGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	var hold string
	if err := conn.QueryRow(`SELECT id FROM holds WHERE task_id=? AND ended_at IS NULL`, task).Scan(&hold); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	result, err := holdwork.Handoff(holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), HoldID: hold, Actor: "human:mark", PR: "https://github.com/owner/repo/pull/63", IdempotencyKey: "merge-handoff"})
	if err != nil {
		t.Fatal(err)
	}
	conn, err = fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE id=?`, result.Fields["handoff_id"]).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	var handoff handoffrecord.Record
	if err := json.Unmarshal([]byte(raw), &handoff); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	options := mergework.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), TaskID: task, Actor: "agent:observer"}
	repo := github.Repository{Owner: "owner", Name: "repo"}
	pr := github.PullRequest{Number: 63, URL: handoff.PR, State: "MERGED", Merged: true, MergedAt: &now, MergeCommit: strings.Repeat("a", 40), HeadSHA: handoff.Head, BaseRef: "trunk", BaseRepository: repo}
	m, err := mergework.RecordMerge(options, handoff.ID, pr, handoff.Tree)
	if err != nil {
		t.Fatal(err)
	}
	run := github.WorkflowRun{ID: 9007199254740993, Attempt: 1, WorkflowID: 20, Name: "CI", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: m.MergeCommit, Status: "completed", Conclusion: "failure", URL: "https://github.com/owner/repo/actions/runs/9007199254740993"}
	if _, err := mergework.RecordCI(options, m.ID, repo, github.CI{Commit: m.MergeCommit, Runs: []github.WorkflowRun{run}}); err != nil {
		t.Fatal(err)
	}
	run.Attempt = 2
	run.Conclusion = "success"
	if _, err := mergework.RecordCI(options, m.ID, repo, github.CI{Commit: m.MergeCommit, Runs: []github.WorkflowRun{run}}); err != nil {
		t.Fatal(err)
	}
	return task
}

func mergeExportSnapshot(t *testing.T, root string) string {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var state strings.Builder
	for _, table := range []string{"merge_observations", "ci_results"} {
		rows, err := conn.Query("SELECT record_json FROM " + table + " ORDER BY rowid")
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var raw string
			if err := rows.Scan(&raw); err != nil {
				t.Fatal(err)
			}
			state.WriteString(raw)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return state.String()
}

func TestMergeCIExportImportAndRebuildPreserveEveryAttempt(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addMergeExportFacts)
	before := mergeExportSnapshot(t, fixture.workspace)
	if !strings.Contains(before, `"conclusion":"success"`) || !strings.Contains(before, `"conclusion":"failure"`) || !strings.Contains(before, `"attempt":2`) {
		t.Fatal("fixture lacks completed attempts")
	}
	if _, err := journal.RebuildWorkspace(fixture.workspace, testMigrationsDir(t)); err != nil {
		t.Fatal(err)
	}
	if got := mergeExportSnapshot(t, fixture.workspace); got != before {
		t.Fatal("rebuild changed observations")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "merge-import"}); err != nil {
		t.Fatal(err)
	}
	if got := mergeExportSnapshot(t, target); got != before {
		t.Fatal("import changed observations")
	}
	assertImportedJournalEquals(t, fixture.workspace, target)
}

func TestMergeCIImportRefusesMalformedBindingsWithoutPublication(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addMergeExportFacts)
	for _, name := range []string{"wrong task", "wrong handoff", "wrong CI head"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, fixture.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			family := "merge.observed"
			if name == "wrong CI head" {
				family = "ci.result-recorded"
			}
			data, hash := mutateArchiveEvent(t, entries[path], family, func(e *journal.Event) {
				switch name {
				case "wrong task":
					mutateObservationJSON(t, e, "task_id", "FID-other")
					e.Effects[0].Values["task_id"] = "FID-other"
				case "wrong handoff":
					mutateObservationJSON(t, e, "handoff_id", "HOF-other")
					e.Effects[0].Values["handoff_id"] = "HOF-other"
				case "wrong CI head":
					mutateObservationJSON(t, e, "head_sha", strings.Repeat("f", 40))
				}
			})
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "invalid.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "refused"})
			assertArchiveCheck(t, err, "reference-invalid")
			if !reflect.DeepEqual(before, workspaceTreeDigest(t, target)) {
				t.Fatal("refused import changed target")
			}
			if _, err := os.Stat(config.LedgerDir(target)); !os.IsNotExist(err) {
				t.Fatalf("refused import published ledger: %v", err)
			}
		})
	}
}
