package export

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func reworkArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, e := reviewArchiveFixture(t)
	if err := migrateSemanticProjection(v.conn, testMigrationsFS(t), 32); err != nil {
		t.Fatal(err)
	}
	var err error
	v.tables, err = loadSemanticTables(v.conn)
	if err != nil {
		t.Fatal(err)
	}
	v.profile = 32
	if _, err := v.conn.Exec(`UPDATE holds SET worktree_id=?`, project.WorktreeID("/w/one")); err != nil {
		t.Fatal(err)
	}
	index := len(e.Effects) - 1
	e.Effects = append(e.Effects[:index],
		journal.Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"status": "in_progress", "revision": 3, "updated_at": "2026-11-10T09:06:00Z"}, Preconditions: map[string]any{"status": "in_review", "revision": 2}},
		journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-1"}, Values: map[string]any{"rework_at": "2026-11-10T09:06:00Z"}, Preconditions: map[string]any{"kept": true, "end_reason": "review", "rework_at": nil}},
		e.Effects[index])
	return v, e
}

func TestReworkRoundAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name, sql string
			change    func(*testing.T, *journal.Event)
			valid     bool
		}{
			{name: "complete", valid: true},
			{name: "todo", sql: `UPDATE issues SET status='todo'`},
			{name: "done", sql: `UPDATE issues SET status='done'`},
			{name: "released reservation", sql: `UPDATE holds SET kept=0`},
			{name: "another task reservation", sql: `INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-other','task','Other','task-default','in_review','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z'); UPDATE holds SET task_id='FID-other'`},
			{name: "no transition", change: func(_ *testing.T, e *journal.Event) { e.Effects = append(e.Effects[:4], e.Effects[5:]...) }},
			{name: "no reservation change", change: func(_ *testing.T, e *journal.Event) { e.Effects = append(e.Effects[:5], e.Effects[6:]...) }},
			{name: "wrong revision", change: func(_ *testing.T, e *journal.Event) { e.Effects[4].Values["revision"] = 99 }},
			{name: "wrong time", change: func(_ *testing.T, e *journal.Event) { e.Effects[5].Values["rework_at"] = "2026-11-10T09:07:00Z" }},
			{name: "wrong reservation", change: func(_ *testing.T, e *journal.Event) { e.Effects[5].PK["id"] = "HLD-other" }},
			{name: "APPROVE with forged transition", change: func(t *testing.T, e *journal.Event) {
				var r reviewrecord.Round
				if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
					t.Fatal(err)
				}
				r.Verdict, r.Findings = "APPROVE", []reviewrecord.Finding{}
				transition, reservation := e.Effects[4], e.Effects[5]
				setReviewArchiveRound(t, e, r)
				retry := e.Effects[1]
				e.Effects = append(e.Effects[:1], transition, reservation, retry)
			}},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, e := reworkArchiveFixture(t)
				if tc.sql != "" {
					if _, err := v.conn.Exec(tc.sql); err != nil {
						t.Fatal(err)
					}
				}
				if tc.change != nil {
					tc.change(t, &e)
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != tc.valid {
					t.Fatalf("admission=%v want valid=%t", err, tc.valid)
				}
				if !tc.valid && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refused rework changed projection")
				}
			})
		}
	}
}

func applyReworkEvent(v *semanticValidator, mode string, e journal.Event) error {
	if mode == "live" {
		return journal.DryRunProjectEvent(v.conn, e)
	}
	return applyLeanEmptyEvent(v, mode, e)
}

func TestReworkHoldAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, tc := range []struct {
			name, root, replaces string
			legacy, noHeartbeat  bool
			keepWorktreeID       bool
			aliases              bool
			valid                bool
		}{
			{name: "complete", root: "/w/one", replaces: "HLD-1", valid: true},
			{name: "another worktree", root: "/w/another", replaces: "HLD-1"},
			{name: "missing link", root: "/w/one"},
			{name: "root differs from ID", root: "/w/another", replaces: "HLD-1", keepWorktreeID: true},
			{name: "missing heartbeat", root: "/w/one", replaces: "HLD-1", noHeartbeat: true},
			{name: "root differs without heartbeat", root: "/w/another", replaces: "HLD-1", keepWorktreeID: true, noHeartbeat: true},
			{name: "legacy another worktree", root: "/w/another", replaces: "HLD-1", legacy: true},
			{name: "legacy same worktree", root: "/w/one", replaces: "HLD-1", legacy: true},
			{name: "legacy missing link", root: "/w/one", legacy: true},
			{name: "legacy case aliases", root: "/w/another", replaces: "HLD-1", legacy: true, aliases: true},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				v, round := reworkArchiveFixture(t)
				acceptMergeArchiveEvent(t, v, round)
				now := time.Date(2026, 11, 10, 9, 7, 0, 0, time.UTC)
				var replaces *string
				if tc.replaces != "" {
					replaces = &tc.replaces
				}
				hold := holds.New("HLD-2", "FID-a", project.WorktreeID(tc.root), tc.root, "agent:fixer", now, replaces)
				hold.Model = "codex/gpt-6.1-sol"
				e := journal.Event{EventSchema: 2, EventID: "EVT-20261110-000005", EventType: "hold.taken", CommandName: "next.claim", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "fixer", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: hold.ID, Type: "hold", Sequence: 1}, Effects: []journal.Effect{
					{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-1"}, Values: map[string]any{"kept": false}, Preconditions: map[string]any{"kept": true}}, hold.Insert(),
				}}
				if tc.keepWorktreeID {
					e.Effects[1].Values["worktree_id"] = project.WorktreeID("/w/one")
				}
				if tc.noHeartbeat {
					delete(e.Effects[1].Values, "last_heartbeat_at")
				}
				if tc.legacy {
					e.EventSchema = journal.EventSchemaV1
					e.Actor = journal.Actor{Role: "developer", Actor: &hold.Actor}
					e.Command = journal.Command{Argv: []string{"fid", "next", "--claim"}, CWD: tc.root, FIDVersion: "test"}
					delete(e.Effects[1].Values, "model")
				}
				if tc.aliases {
					e.Effects[0].Table = "HOLDS"
					e.Effects[0].PK = map[string]any{"ID": "HLD-1"}
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyReworkEvent(v, mode, e)
				if (err == nil) != tc.valid {
					t.Errorf("rework hold=%v want valid=%t", err, tc.valid)
				}
				if !tc.valid && handoffProjectionSnapshot(t, v) != before {
					t.Fatal("refused rework hold changed projection")
				}
			})
		}
	}
}
