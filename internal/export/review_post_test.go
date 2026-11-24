package export

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
)

func postArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, e := responseArchiveFixture(t, false)
	if err := migrateSemanticProjection(v.conn, testMigrationsDir(t), 41); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 41
	s, err := reviewcomment.ReadSubject(v.conn, "RND-123456")
	if err != nil {
		t.Fatal(err)
	}
	e.EventType, e.CommandName = "review.posted", "review.post"
	e.ActorV2 = journal.EventActor{ActorType: "service", ActorID: "fid-review-post", Producer: "fid", ProvenanceRef: "fid:review-post"}
	e.Aggregate = journal.Aggregate{ID: s.ID, Type: "review-post", Sequence: 1}
	e.Effects = []journal.Effect{{Op: "insert", Table: "review_posts", PK: map[string]any{"id": s.ID}, Values: map[string]any{"task_id": s.TaskID, "handoff_id": s.HandoffID, "pr": s.PR, "url": "https://github.com/owner/repo/issues/63#issuecomment-123", "render_version": reviewcomment.RenderVersion, "digest": reviewcomment.Digest(s.Body), "posted_at": "2026-11-10T09:07:00Z"}}}
	return v, e
}

func TestReviewPostAdmission(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"valid", "actor", "command", "subject", "task", "handoff", "pr", "url", "digest", "version", "time", "extra", "legacy", "update"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v, e := postArchiveFixture(t)
			switch name {
			case "actor":
				e.ActorV2.ActorID = "other-service"
			case "command":
				e.CommandName = "review.record"
			case "subject":
				e.Aggregate.ID = "RND-other"
			case "task":
				e.Effects[0].Values["task_id"] = "FID-other"
			case "handoff":
				e.Effects[0].Values["handoff_id"] = "HOF-other"
			case "pr":
				e.Effects[0].Values["pr"] = "https://github.com/owner/repo/pull/64"
			case "url":
				e.Effects[0].Values["url"] = "https://github.com/owner/repo/issues/64#issuecomment-123"
			case "digest":
				e.Effects[0].Values["digest"] = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			case "version":
				e.Effects[0].Values["render_version"] = 99
			case "time":
				e.Effects[0].Values["posted_at"] = "2026-11-10T09:08:00Z"
			case "extra":
				e.Effects = append(e.Effects, e.Effects[0])
			case "legacy":
				e.EventSchema = 1
			case "update":
				e.Effects[0].Op = "update"
			}
			before := handoffProjectionSnapshot(t, v)
			err := applyLeanEmptyEvent(v, "import", e)
			if name == "valid" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("forged post admitted")
			}
			if handoffProjectionSnapshot(t, v) != before {
				t.Fatal("refused post changed projection")
			}
			_, control := postArchiveFixture(t)
			if err := applyLeanEmptyEvent(v, "import", control); err != nil {
				t.Fatal("valid write after refusal", err)
			}
		})
	}
}

func TestReviewCommentKeepsRecordedHandoff(t *testing.T) {
	t.Parallel()
	v, e := responseArchiveFixture(t, false)
	if err := applyLeanEmptyEvent(v, "replay", e); err != nil {
		t.Fatal(err)
	}
	if _, err := v.conn.Exec(`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason) SELECT 'HLD-later',task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,'review' FROM holds WHERE id='HLD-1'`); err != nil {
		t.Fatal(err)
	}
	if _, err := v.conn.Exec(`INSERT INTO handoffs SELECT 'HOF-later',task_id,'HLD-later',policy_snapshot,json_set(record_json,'$.id','HOF-later','$.hold_id','HLD-later','$.pr','https://github.com/owner/repo/pull/64') FROM handoffs WHERE id='HOF-100001'`); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"RND-123456", "FRS-one"} {
		s, err := reviewcomment.ReadSubject(v.conn, id)
		if err != nil || s.HandoffID != "HOF-100001" || s.PR != "https://github.com/owner/repo/pull/63" {
			t.Fatalf("pending delivery redirected: %v/%v", s, err)
		}
	}
}
