package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fidexport "github.com/mbell-nzl/fid/internal/export"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestReviewPostVersionsSurviveRebuild(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	in := reviewInput()
	in["findings"] = []any{reviewFinding("minor")}
	first, _ := lifecycleResult(t, recordReviewTask(t, f, in, "legacy-round", "human:courier", task), "review.record")
	id := first["round"].(map[string]any)["id"].(string)
	finding := id + "-1"
	answer, _ := lifecycleResult(t, f.run("--json", "finding", "answer", finding, "--actor", "agent:builder", "--text", "Accepted explanation"), "finding.answer")
	deferred, _ := lifecycleResult(t, f.run("--json", "finding", "defer", finding, "--actor", "agent:builder", "--title", "Follow up"), "finding.defer")
	a, d := answer["response"].(map[string]any), deferred["response"].(map[string]any)
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	old, err := os.ReadFile(filepath.Join("testdata", "review", "comment_v1.md"))
	if err != nil {
		t.Fatal(err)
	}
	body := strings.ReplaceAll(strings.Split(string(old), "\n<!-- fid-review:")[0], "RND-<new>", id)
	var handoff string
	if err := conn.QueryRow(`SELECT handoff_id FROM review_rounds WHERE id=?`, id).Scan(&handoff); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	legacy := map[string]string{id: body,
		a["id"].(string): fmt.Sprintf("## Finding answer: %s\n\nBy agent:builder\n\n%s\n", finding, a["text"]),
		d["id"].(string): fmt.Sprintf("## Deferred finding: %s\n\nBy agent:builder · Discovered task %s: Follow up\n\n%s\n", finding, d["deferred_task"].(map[string]any)["id"], d["text"]),
	}
	comment := 999
	for subject, text := range legacy {
		comment++
		e := journal.Event{EventSchema: 2, EventType: "review.posted", CommandName: "review.post", OccurredAt: now,
			ActorV2: journal.EventActor{ActorID: "fid-review-post", ActorType: "service", Producer: "fid", ProvenanceRef: "fid:review-post"}, Aggregate: journal.Aggregate{ID: id, Type: "review-post", Sequence: 1},
			Effects: []journal.Effect{{Op: "insert", Table: "review_posts", PK: map[string]any{"id": subject}, Values: map[string]any{"task_id": task, "handoff_id": handoff, "pr": "https://github.com/mbell-nzl/fid/pull/63", "url": fmt.Sprintf("https://github.com/mbell-nzl/fid/issues/63#issuecomment-%d", comment), "render_version": 1, "digest": reviewcomment.Digest(text), "posted_at": now.Format(time.RFC3339)}}}}
		e.Aggregate.ID = subject
		if err := journal.WithLock(f.workspace, journal.LockOwner{Command: "test legacy receipt"}, func(lock *journal.WriteLock) error {
			_, err := journal.CommitEvent(lock, conn, e)
			return err
		}); err != nil {
			t.Fatal(err)
		}
	}
	// A caught-up schema-41 cache must upgrade even with no pending Events.
	down, err := fs.ReadFile(fidmigrations.FS, "0042_review_render_versions.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(string(down)); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE schema_migrations SET version=41; UPDATE projection_meta SET projection_schema_version=22`); err != nil {
		t.Fatal(err)
	}
	if _, err := journal.EnsureProjectionFresh(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	second, _ := lifecycleResult(t, recordReviewTask(t, f, in, "new-round", "human:courier", task), "review.record")
	newID := second["round"].(map[string]any)["id"].(string)
	run, _, _ := recoveryCommentRunner(t, false)
	posted, err := reviewpost.Round(context.Background(), f.workspace, fidmigrations.FS, newID, run)
	if err != nil || len(posted.Pending) != 0 {
		t.Fatalf("post: %+v / %v", posted, err)
	}
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal("rebuild mixed render versions", err)
	}
	rebuilt, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, rebuilt)
	versions := map[string]int{newID: 2}
	for subject := range legacy {
		versions[subject] = 1
	}
	for subject, version := range versions {
		var got int
		if err := rebuilt.QueryRow(`SELECT render_version FROM review_posts WHERE id=?`, subject).Scan(&got); err != nil || got != version {
			t.Fatalf("receipt %s version=%d / %v, want %d", subject, got, err, version)
		}
	}
	var migration int
	if err := rebuilt.QueryRow(`SELECT version FROM schema_migrations`).Scan(&migration); err != nil || migration != 44 {
		t.Fatalf("upgraded migration=%d / %v", migration, err)
	}
	archive := filepath.Join(t.TempDir(), "mixed.zip")
	if _, err := fidexport.Create(fidexport.CreateOptions{WorkspaceRoot: f.workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	if _, err := fidexport.Validate(fidexport.ValidateOptions{ArchivePath: archive, MigrationsFS: fidmigrations.FS}); err != nil {
		t.Fatal("validate mixed-version archive", err)
	}
	if _, err := fidexport.Import(fidexport.ImportOptions{ArchivePath: archive, RepositoryRoot: gittest.NewWorkspace(t), MigrationsFS: fidmigrations.FS, IdempotencyKey: "mixed-versions"}); err != nil {
		t.Fatal("import mixed-version archive", err)
	}
}
