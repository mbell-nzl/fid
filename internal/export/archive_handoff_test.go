package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func handoffArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	validator, err := newSemanticValidator(testMigrationsDir(t), 26)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = validator.close() })
	_, err = validator.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,lane,owner_role,created_at,updated_at) VALUES('FID-a','task','Work','task-default','in_progress','standard','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(func(*journal.Event) {}))
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.accept(raw); err != nil {
		t.Fatal(err)
	}
	file, err := config.DecodeProject([]byte("[lanes.standard]\nchecks=[]\n[review]\ntest_adequacy=[]\n"))
	if err != nil {
		t.Fatal(err)
	}
	effective := policy.Resolve(file)
	encoded, err := effective.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	_, err = validator.conn.Exec(`INSERT INTO policy_snapshots(sequence,source_ref,digest,effective_json) VALUES(1,'built-in defaults',?,?)`, effective.Digest(), string(encoded))
	if err != nil {
		t.Fatal(err)
	}
	event := archiveStartEvent(t)
	event.EventType, event.CommandName = "task.handed-off", "task.handoff"
	record := handoffrecord.Record{View: handoffrecord.View{ID: "HOF-100001", PR: "https://github.com/owner/repo/pull/63", Head: strings.Repeat("a", 40), Base: strings.Repeat("b", 40), RequiredChecks: []string{}, ResultIDs: []string{}, At: "2026-11-10T09:05:00Z"}, TaskID: "FID-a", HoldID: "HLD-1", Actor: "agent:builder", Tree: strings.Repeat("c", 40), Lane: "standard", PolicySnapshot: 1, DefinitionDigests: map[string]string{}, ChangedPaths: []string{}}
	encoded, err = json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	event.Effects = []journal.Effect{
		{Op: "insert", Table: "handoffs", PK: map[string]any{"id": record.ID}, Values: map[string]any{"task_id": record.TaskID, "hold_id": record.HoldID, "policy_snapshot": 1, "record_json": string(encoded)}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": record.TaskID}, Values: map[string]any{"status": "in_review", "revision": 2, "updated_at": record.At}, Preconditions: map[string]any{"status": "in_progress", "revision": 1}},
		{Op: "update", Table: "holds", PK: map[string]any{"id": record.HoldID}, Values: map[string]any{"ended_at": record.At, "end_reason": "review", "kept": true}, Preconditions: map[string]any{"ended_at": nil}},
	}
	takeoverArchiveRetry(&event)
	event.Idempotency.Scope = "task.handoff:project-test"
	last := &event.Effects[3]
	last.PK["operation_scope"] = event.Idempotency.Scope
	encoded, err = journal.MarshalCanonicalJSON(map[string]any{"fields": record.Fields(), "occurred_at": record.At, "replayed": false, "worktree": ""})
	if err != nil {
		t.Fatal(err)
	}
	last.Values["result_json"] = string(encoded)
	return validator, event
}

func TestHandoffArchiveAndReplayValidateAtomicReceipt(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"archive", "replay"} {
		for _, row := range []struct {
			name, sql string
			mutate    func(*journal.Event)
			valid     bool
		}{
			{name: "complete", valid: true},
			{name: "wrong actor", mutate: func(e *journal.Event) { e.ActorV2.ActorID = "other" }},
			{name: "service actor", mutate: func(e *journal.Event) { e.ActorV2.ActorType = "service" }},
			{name: "wrong command", mutate: func(e *journal.Event) { e.CommandName = "task.start" }},
			{name: "no reservation", mutate: func(e *journal.Event) { e.Effects[2].Values["kept"] = false }},
			{name: "no transition", mutate: func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2:]...) }},
			{name: "wrong revision", mutate: func(e *journal.Event) { e.Effects[1].Values["revision"] = 9 }},
			{name: "wrong task", mutate: func(e *journal.Event) { e.Effects[0].Values["task_id"] = "FID-other" }},
			{name: "malformed receipt", mutate: func(e *journal.Event) { e.Effects[0].Values["record_json"] = "{}" }},
			{name: "malformed retry", mutate: func(e *journal.Event) { e.Effects[3].Values["result_json"] = "{}" }},
			{name: "not started", sql: `UPDATE issues SET status='todo' WHERE id='FID-a'`},
			{name: "wrong lane", sql: `UPDATE issues SET lane='assured' WHERE id='FID-a'`},
			{name: "ended hold", sql: `UPDATE holds SET ended_at='2026-11-10T09:01:00Z',end_reason='takeover' WHERE id='HLD-1'`},
			{name: "expired hold", sql: `UPDATE holds SET expires_at='2026-11-10T09:05:00Z' WHERE id='HLD-1'`},
			{name: "policy changed", sql: `INSERT INTO policy_snapshots SELECT 2,source_ref,source_commit,digest,effective_json FROM policy_snapshots WHERE sequence=1`},
			{name: "missing policy", sql: `DELETE FROM policy_snapshots`},
		} {
			t.Run(mode+"/"+row.name, func(t *testing.T) {
				t.Parallel()
				validator, event := handoffArchiveFixture(t)
				if row.sql != "" {
					if _, err := validator.conn.Exec(row.sql); err != nil {
						t.Fatal(err)
					}
				}
				if row.mutate != nil {
					row.mutate(&event)
				}
				var before string
				if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-a'`).Scan(&before); err != nil {
					t.Fatal(err)
				}
				var err error
				if mode == "archive" {
					var raw []byte
					raw, err = journal.MarshalCanonicalJSON(event)
					if err == nil {
						err = validator.accept(raw)
					}
				} else {
					err = journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: "sha256:handoff", File: "journal", Line: 2, Offset: 2})
				}
				if (err == nil) != row.valid {
					t.Fatalf("accept=%v, valid=%v", err, row.valid)
				}
				var state string
				var count, bindings int
				if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-a'`).Scan(&state); err != nil {
					t.Fatal(err)
				}
				if err := validator.conn.QueryRow(`SELECT count(*) FROM handoffs`).Scan(&count); err != nil {
					t.Fatal(err)
				}
				if err := validator.conn.QueryRow(`SELECT count(*) FROM idempotency_results`).Scan(&bindings); err != nil {
					t.Fatal(err)
				}
				if row.valid {
					if state != "in_review" || count != 1 || bindings != 1 {
						t.Fatalf("state=%s/%d/%d", state, count, bindings)
					}
				} else if state != before || count != 0 || bindings != 0 {
					t.Fatalf("partial rejected state=%s/%d/%d", state, count, bindings)
				}
			})
		}
	}
}
