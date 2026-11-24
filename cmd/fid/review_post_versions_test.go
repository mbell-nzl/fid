package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestReviewPostVersionsSurviveRebuild(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	in := reviewInput()
	in["findings"] = []any{reviewFinding("minor")}
	first, _ := lifecycleResult(t, recordReviewTask(t, f, in, "legacy-round", "human:courier", task), "review.record")
	id := first["round"].(map[string]any)["id"].(string)
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
	e := journal.Event{EventSchema: 2, EventType: "review.posted", CommandName: "review.post", OccurredAt: now,
		ActorV2: journal.EventActor{ActorID: "fid-review-post", ActorType: "service", Producer: "fid", ProvenanceRef: "fid:review-post"}, Aggregate: journal.Aggregate{ID: id, Type: "review-post", Sequence: 1},
		Effects: []journal.Effect{{Op: "insert", Table: "review_posts", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": task, "handoff_id": handoff, "pr": "https://github.com/mbell-nzl/fid/pull/63", "url": "https://github.com/mbell-nzl/fid/issues/63#issuecomment-999", "render_version": 1, "digest": reviewcomment.Digest(body), "posted_at": now.Format(time.RFC3339)}}}}
	if err := journal.WithLock(f.workspace, journal.LockOwner{Command: "test legacy receipt"}, func(lock *journal.WriteLock) error {
		_, err := journal.CommitEvent(lock, conn, e)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	second, _ := lifecycleResult(t, recordReviewTask(t, f, in, "new-round", "human:courier", task), "review.record")
	newID := second["round"].(map[string]any)["id"].(string)
	run, _, _ := recoveryCommentRunner(t, false)
	posted, err := reviewpost.Round(context.Background(), f.workspace, issueMigrationsDir(), newID, run)
	if err != nil || len(posted.Pending) != 0 {
		t.Fatalf("post: %+v / %v", posted, err)
	}
	if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
		t.Fatal("rebuild mixed render versions", err)
	}
	rebuilt, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, rebuilt)
	for subject, version := range map[string]int{id: 1, newID: 2} {
		var got int
		if err := rebuilt.QueryRow(`SELECT render_version FROM review_posts WHERE id=?`, subject).Scan(&got); err != nil || got != version {
			t.Fatalf("receipt %s version=%d / %v, want %d", subject, got, err, version)
		}
	}
}
