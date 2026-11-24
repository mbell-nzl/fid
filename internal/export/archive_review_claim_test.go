package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func reviewClaimArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, _ := reviewArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 36); err != nil {
		t.Fatal(err)
	}
	var err error
	v.profile = 36
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.conn.Exec(`UPDATE holds SET worktree_id=? WHERE worktree_root=?`, project.WorktreeID("/w/one"), "/w/one"); err != nil {
		t.Fatal(err)
	}
	stamp := "2026-11-10T09:06:00Z"
	fp := fingerprint.Result{Algorithm: fingerprint.Algorithm, Fingerprint: strings.Repeat("f", 64), Base: strings.Repeat("b", 40), Head: strings.Repeat("a", 40)}
	e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000004", EventType: "review.claimed", CommandName: "next.claim", OccurredAt: time.Date(2026, 11, 10, 9, 6, 0, 0, time.UTC), ActorV2: journal.EventActor{ActorType: "agent", ActorID: "reviewer", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "review-claim", ID: "RCL-123456", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "review_claims", PK: map[string]any{"id": "RCL-123456"}, Values: map[string]any{
		"task_id": "FID-a", "handoff_id": "HOF-100001", "stream": "primary", "worktree_id": project.WorktreeID("/w/review"), "worktree_root": "/w/review", "actor": "agent:reviewer", "model": "claude/opus-5.5", "since": stamp, "last_heartbeat_at": stamp, "expires_at": "2026-11-10T10:06:00Z", "ended_at": nil, "binding_json": encodedLean(fp),
	}}}}
	return v, e
}

func TestReviewClaimAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			mutate func(*journal.Event)
		}{
			{"valid", nil},
			{"builder", func(e *journal.Event) { e.ActorV2.ActorID = "builder"; e.Effects[0].Values["actor"] = "agent:builder" }},
			{"actor", func(e *journal.Event) { e.Effects[0].Values["actor"] = "agent:other" }},
			{"task", func(e *journal.Event) { e.Effects[0].Values["task_id"] = "FID-other" }},
			{"handoff", func(e *journal.Event) { e.Effects[0].Values["handoff_id"] = "HOF-other" }},
			{"stream", func(e *journal.Event) { e.Effects[0].Values["stream"] = "separate" }},
			{"root", func(e *journal.Event) { e.Effects[0].Values["worktree_root"] = "/w/other" }},
			{"busy worktree", func(e *journal.Event) {
				e.Effects[0].Values["worktree_root"] = "/w/one"
				e.Effects[0].Values["worktree_id"] = project.WorktreeID("/w/one")
			}},
			{"lapsed in-review build", func(e *journal.Event) {
				e.Effects[0].Values["worktree_root"] = "/w/one"
				e.Effects[0].Values["worktree_id"] = project.WorktreeID("/w/one")
				e.OccurredAt = e.OccurredAt.Add(2 * time.Hour)
				stamp := e.OccurredAt.Format(time.RFC3339)
				e.Effects[0].Values["since"], e.Effects[0].Values["last_heartbeat_at"], e.Effects[0].Values["expires_at"] = stamp, stamp, e.OccurredAt.Add(time.Hour).Format(time.RFC3339)
			}},
			{"time", func(e *journal.Event) { e.Effects[0].Values["expires_at"] = "2026-11-11T10:06:00Z" }},
			{"model", func(e *journal.Event) { e.Effects[0].Values["model"] = "" }},
			{"command", func(e *journal.Event) { e.CommandName = "review.record" }},
			{"aggregate", func(e *journal.Event) { e.Aggregate.ID = "RCL-other" }},
			{"duplicate insert", func(e *journal.Event) { e.Effects = append(e.Effects, e.Effects[0]) }},
			{"v1", func(e *journal.Event) { e.EventSchema = 1 }},
			{"other family", func(e *journal.Event) { e.EventType = "hold.taken" }},
			{"fingerprint hex", func(e *journal.Event) {
				e.Effects[0].Values["binding_json"] = strings.Replace(e.Effects[0].Values["binding_json"].(string), strings.Repeat("f", 64), strings.Repeat("z", 64), 1)
			}},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reviewClaimArchiveFixture(t)
				if tc.mutate != nil {
					tc.mutate(&e)
				}
				if tc.name == "lapsed in-review build" {
					if _, err := v.conn.Exec(`UPDATE holds SET kept=0`); err != nil {
						t.Fatal(err)
					}
					if _, err := v.conn.Exec(`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-takeover','FID-a',?, '/w/one','human:mark','2026-11-10T09:06:00Z','2026-11-10T10:06:00Z')`, project.WorktreeID("/w/one")); err != nil {
						t.Fatal(err)
					}
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != (tc.mutate == nil) {
					t.Fatalf("admission=%v", err)
				}
				if tc.mutate != nil {
					var count int
					if err := v.conn.QueryRow(`SELECT count(*) FROM review_claims`).Scan(&count); err != nil || count != 0 || handoffProjectionSnapshot(t, v) != before {
						t.Fatalf("refusal changed state: claims=%d err=%v", count, err)
					}
				}
			})
		}
	}
}

func TestReviewClaimReplacementAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name           string
			elapsed        time.Duration
			omitEnd, valid bool
		}{
			{"live cannot be replaced", time.Minute, false, false},
			{"lapsed can be replaced", time.Hour, false, true},
			{"lapsed must be ended", time.Hour, true, false},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reviewClaimArchiveFixture(t)
				if err := applyLeanEmptyEvent(v, "replay", e); err != nil {
					t.Fatal(err)
				}
				e.EventID = "EVT-20261110-000005"
				e.Aggregate.ID = "RCL-654321"
				e.Effects[0].PK["id"] = e.Aggregate.ID
				e.OccurredAt = e.OccurredAt.Add(tc.elapsed)
				stamp := e.OccurredAt.Format(time.RFC3339)
				e.Effects[0].Values["since"], e.Effects[0].Values["last_heartbeat_at"], e.Effects[0].Values["expires_at"] = stamp, stamp, e.OccurredAt.Add(time.Hour).Format(time.RFC3339)
				if !tc.omitEnd {
					e.Effects = append([]journal.Effect{journal.EndReviewClaim("RCL-123456", stamp)}, e.Effects...)
				}
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != tc.valid {
					t.Fatalf("admission=%v want valid=%t", err, tc.valid)
				}
				if !tc.valid {
					var n int
					if err := v.conn.QueryRow(`SELECT count(*) FROM review_claims WHERE ended_at IS NULL`).Scan(&n); err != nil || n != 1 {
						t.Fatalf("refusal changed reservation %d/%v", n, err)
					}
				}
			})
		}
	}
}

func TestLiveReviewBlocksBuildAdmission(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, claim := reviewClaimArchiveFixture(t)
			if err := applyLeanEmptyEvent(v, "replay", claim); err != nil {
				t.Fatal(err)
			}
			if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,lane,created_at,updated_at) VALUES('FID-b','task','New work','task-default','todo','standard','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`); err != nil {
				t.Fatal(err)
			}
			at := claim.OccurredAt.Add(time.Minute)
			h := holds.New("HLD-new", "FID-b", project.WorktreeID("/w/review"), "/w/review", "agent:new-builder", at, nil)
			event := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "hold.taken", CommandName: "next.claim", OccurredAt: at, ActorV2: journal.EventActor{ActorType: "agent", ActorID: "new-builder", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: h.ID, Type: "hold", Sequence: 1}, Effects: []journal.Effect{h.Insert()}}
			if err := applyReworkEvent(v, mode, event); err == nil {
				t.Fatal("a worktree already reviewing took build work")
			}
		})
	}
}

func TestReviewHeartbeatAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name   string
			change func(*journal.Event)
		}{
			{"valid", nil},
			{"actor", func(e *journal.Event) { e.ActorV2.ActorID = "other" }},
			{"expired", func(e *journal.Event) { e.OccurredAt = e.OccurredAt.Add(time.Hour) }},
			{"wrong renewal", func(e *journal.Event) { e.Effects[0].Values["expires_at"] = "2026-11-11T10:07:00Z" }},
			{"immutable field", func(e *journal.Event) { e.Effects[0].Values["stream"] = "separate" }},
			{"command", func(e *journal.Event) { e.CommandName = "next.claim" }},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, claim := reviewClaimArchiveFixture(t)
				if err := applyLeanEmptyEvent(v, "replay", claim); err != nil {
					t.Fatal(err)
				}
				at := claim.OccurredAt.Add(time.Minute)
				e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "review.heartbeat", CommandName: "hold.heartbeat", OccurredAt: at, ActorV2: claim.ActorV2, Aggregate: journal.Aggregate{ID: claim.Aggregate.ID, Type: "review-claim", Sequence: 2}, Effects: []journal.Effect{{Op: "update", Table: "review_claims", PK: map[string]any{"id": claim.Aggregate.ID}, Values: map[string]any{"last_heartbeat_at": at.Format(time.RFC3339), "expires_at": at.Add(time.Hour).Format(time.RFC3339)}, Preconditions: map[string]any{"ended_at": nil, "expires_at": "2026-11-10T10:06:00Z"}}}}
				if tc.change != nil {
					tc.change(&e)
				}
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != (tc.change == nil) {
					t.Fatalf("admission=%v", err)
				}
				if tc.change != nil {
					var stamp string
					if err := v.conn.QueryRow(`SELECT last_heartbeat_at FROM review_claims`).Scan(&stamp); err != nil || stamp != "2026-11-10T09:06:00Z" {
						t.Fatalf("refusal changed heartbeat %s/%v", stamp, err)
					}
				}
			})
		}
	}
}

func TestReviewClaimRejectsUnrelatedRelease(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, e := reviewClaimArchiveFixture(t)
			if err := applyLeanEmptyEvent(v, "replay", e); err != nil {
				t.Fatal(err)
			}
			h, err := reviewrecord.LatestHandoff(v.conn, "FID-a")
			if err != nil {
				t.Fatal(err)
			}
			h.ID, h.TaskID, h.HoldID = "HOF-peer", "FID-peer", "HLD-peer"
			raw, err := json.Marshal(h)
			if err != nil {
				t.Fatal(err)
			}
			for _, statement := range []string{
				`INSERT INTO issues(id,kind,title,workflow_id,status,lane,created_at,updated_at) VALUES('FID-peer','task','Peer','task-default','in_review','standard','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`,
				`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason,kept) VALUES('HLD-peer','FID-peer','peer-builder','/w/peer-builder','agent:peer-builder','2026-11-10T09:00:00Z','2026-11-10T10:00:00Z','2026-11-10T09:05:00Z','review',1)`,
			} {
				if _, err := v.conn.Exec(statement); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := v.conn.Exec(`INSERT INTO handoffs VALUES(?,?,?,?,?)`, h.ID, h.TaskID, h.HoldID, h.PolicySnapshot, string(raw)); err != nil {
				t.Fatal(err)
			}
			bytes, err := json.Marshal(e)
			if err != nil {
				t.Fatal(err)
			}
			var peer journal.Event
			if err := json.Unmarshal(bytes, &peer); err != nil {
				t.Fatal(err)
			}
			peer.EventID, peer.Aggregate.ID, peer.ActorV2.ActorID = "EVT-20261110-000005", "RCL-peer", "peer-reviewer"
			peer.Effects[0].PK["id"] = "RCL-peer"
			values := peer.Effects[0].Values
			values["task_id"], values["handoff_id"], values["actor"] = "FID-peer", "HOF-peer", "agent:peer-reviewer"
			values["worktree_root"], values["worktree_id"] = "/w/peer-review", project.WorktreeID("/w/peer-review")
			if err := applyLeanEmptyEvent(v, "replay", peer); err != nil {
				t.Fatal(err)
			}
			e.EventID, e.Aggregate.ID = "EVT-20261110-000006", "RCL-new"
			e.OccurredAt = e.OccurredAt.Add(time.Hour)
			stamp := e.OccurredAt.Format(time.RFC3339)
			e.Effects[0].PK["id"] = e.Aggregate.ID
			e.Effects[0].Values["since"], e.Effects[0].Values["last_heartbeat_at"], e.Effects[0].Values["expires_at"] = stamp, stamp, e.OccurredAt.Add(time.Hour).Format(time.RFC3339)
			e.Effects = append([]journal.Effect{journal.EndReviewClaim("RCL-123456", stamp), journal.EndReviewClaim("RCL-peer", stamp)}, e.Effects...)
			if err := applyReworkEvent(v, mode, e); err == nil {
				t.Fatal("a review released an unrelated task/worktree reservation")
			}
			var n int
			if err := v.conn.QueryRow(`SELECT count(*) FROM review_claims WHERE ended_at IS NULL`).Scan(&n); err != nil || n != 2 {
				t.Fatalf("refusal changed reservations: %d/%v", n, err)
			}
		})
	}
}
