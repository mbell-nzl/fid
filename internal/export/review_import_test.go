package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func addReviewExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addLeanExportFacts(t, root)
	file, err := os.OpenFile(filepath.Join(root, ".fid.toml"), os.O_APPEND|os.O_WRONLY, 0)
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
	fixtureGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Review policy")
	fixtureGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "review-change.txt"), []byte("reviewed change\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, root, "add", ".")
	fixtureGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Reviewed change")
	taken, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:coded", Model: "codex/gpt-6.1-sol", Reason: "Build with model"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = holdwork.Handoff(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), HoldID: taken.Fields["hold_id"].(string), Actor: "human:coded", PR: "https://github.com/owner/repo/pull/63", IdempotencyKey: "review-handoff"})
	if err != nil {
		t.Fatal(err)
	}
	in := reviewrecord.Input{Reviewer: reviewrecord.Reviewer{Actor: "agent:reviewer", Model: "claude/opus-5.5", Effort: "high", SessionID: "portable-session"}, Stream: "primary", Profiles: []string{"code", "test-adequacy"}, Verdict: "CHANGES", Purpose: "Persist canonical reviews", Verification: "go test ./... passed; a probe failed", Provenance: &reviewrecord.Provenance{Prompt: "Review the change", Command: "review-command", Output: "## Review: CHANGES\nClaim: Persist canonical reviews\nVerification: passed"}, Findings: []reviewrecord.FindingInput{}}
	for _, severity := range []string{"major", "minor", "nit"} {
		in.Findings = append(in.Findings, reviewrecord.FindingInput{Profile: "code", Severity: severity, Summary: severity + " finding", Location: &reviewrecord.Location{File: "internal/example.go", Line: 42}, Evidence: "Observed behaviour", Fix: "Add coverage"})
	}
	if _, err := reviewwork.Record(reviewwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:courier", IdempotencyKey: "review-record", Input: in}); err != nil {
		t.Fatal(err)
	}
	return task
}

func reviewExportSnapshot(t *testing.T, root string) string {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var state strings.Builder
	for _, query := range []string{`SELECT record_json FROM review_rounds ORDER BY rowid`, `SELECT record_json FROM round_findings ORDER BY rowid`, `SELECT id || ':' || model FROM holds ORDER BY rowid`, `SELECT result_json FROM idempotency_results WHERE operation_scope LIKE 'review.record:%' ORDER BY rowid`} {
		rows, err := conn.Query(query)
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

func TestReviewExportImportRebuildPreservesCanonicalReview(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addReviewExportFacts)
	before := reviewExportSnapshot(t, f.workspace)
	for _, text := range []string{"major finding", "minor finding", "nit finding", "portable-session", "review-command", "codex/gpt-6.1-sol", "claude/opus-5.5"} {
		if !strings.Contains(before, text) {
			t.Fatalf("fixture missing %s", text)
		}
	}
	if _, err := journal.RebuildWorkspace(f.workspace, testMigrationsFS(t)); err != nil {
		t.Fatal(err)
	}
	if reviewExportSnapshot(t, f.workspace) != before {
		t.Fatal("rebuild changed review")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "review-import"}); err != nil {
		t.Fatal(err)
	}
	if reviewExportSnapshot(t, target) != before {
		t.Fatal("import changed review")
	}
	assertImportedJournalEquals(t, f.workspace, target)
	conn, err := reviewwork.OpenProjection(target)
	if err != nil {
		t.Fatal(err)
	}
	var task string
	if err := conn.QueryRow(`SELECT task_id FROM review_rounds`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	// Import does not bring repository objects; canonical current reviews remain readable.
	views, err := reviewwork.List(target, task)
	if err != nil || len(views) != 1 || !views[0].Current || !views[0].OtherModel {
		t.Fatalf("imported read=%v/%v", views, err)
	}
}

func TestReviewArchiveImportRefusesForgedBindingWithoutPublication(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addReviewExportFacts)
	for _, field := range []string{"reviewer", "head", "retry"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, f.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], "review.round-recorded", func(e *journal.Event) {
				if field == "retry" {
					e.Effects[len(e.Effects)-1].Values["result_json"] = "{}"
					return
				}
				var r reviewrecord.Round
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				if field == "reviewer" {
					r.Reviewer.Actor = "human:coded"
				} else {
					r.Binding.Head = strings.Repeat("e", 40)
				}
				e.Effects[0].Values["record_json"] = encodedLean(r)
				e.Effects[len(e.Effects)-1].Values["result_json"] = encodedLean(reviewwork.Result{Round: r, At: r.At})
			})
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "forged.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "refused"})
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
