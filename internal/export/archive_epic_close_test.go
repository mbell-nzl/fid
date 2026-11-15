package export

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func epicCloseArchiveFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, err := newSemanticValidator(testMigrationsDir(t), 35)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = v.close() })
	for _, id := range []string{"FID-epic", "FID-other"} {
		if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES(?,'epic','Epic','epic-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`, id); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := v.conn.Exec(`INSERT INTO epic_metadata(issue_id,scheduling) VALUES('FID-epic','prioritised'),('FID-other','prioritised')`); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 15, 1, 0, 0, 0, time.UTC)
	return v, journal.Event{EventSchema: 2, EventType: "epic.closed", EventID: "EVT-20261115-000001", CommandName: "task.done", OccurredAt: at, Aggregate: journal.Aggregate{Type: "epic", ID: "FID-epic", Sequence: 1}, ActorV2: journal.EventActor{ActorType: "human", ActorID: "mark", Producer: "fid-standalone", ProvenanceRef: "test"}, Effects: []journal.Effect{
		{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-123456"}, Values: map[string]any{"kind": "close-epic", "actor": "human:mark", "subject_type": "epic", "subject_id": "FID-epic", "binding_id": "FID-epic", "reason": "Finished", "decided_at": at.Format(time.RFC3339)}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-epic"}, Values: map[string]any{"status": "done", "revision": 2, "updated_at": at.Format(time.RFC3339)}, Preconditions: map[string]any{"status": "todo", "revision": 1}},
	}}
}

func TestEpicCloseAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, name := range []string{"valid", "in progress", "done child", "cancelled child", "todo child", "started child", "review child", "late open child", "open subtask", "agent", "service", "blank reason", "task subject", "wrong actor", "wrong binding", "wrong subject", "wrong time", "wrong kind", "missing decision", "missing move", "other move", "extra move", "extra decision", "injected retry", "wrong revision", "wrong from", "wrong target", "changed title", "wrong command", "wrong aggregate", "legacy", "other family", "case alias", "closed epic", "cancelled epic"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := epicCloseArchiveFixture(t)
				want := name == "valid" || name == "in progress" || name == "done child" || name == "cancelled child" || name == "open subtask"
				switch name {
				case "in progress", "closed epic", "cancelled epic":
					kind, workflow, state := "epic", "epic-default", "in_progress"
					if name == "closed epic" {
						state = "done"
					}
					if name == "cancelled epic" {
						state = "cancelled"
					}

					if _, err := v.conn.Exec(`UPDATE issues SET kind=?,workflow_id=?,status=? WHERE id='FID-epic'`, kind, workflow, state); err != nil {
						t.Fatal(err)
					}
					e.Effects[1].Preconditions["status"] = state
				case "task subject":
					if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-task','task','Task','task-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`); err != nil {
						t.Fatal(err)
					}
					e.Aggregate.ID = "FID-task"
					e.Effects[0].Values["subject_id"] = "FID-task"
					e.Effects[0].Values["binding_id"] = "FID-task"
					e.Effects[1].PK["id"] = "FID-task"
				case "done child", "cancelled child", "todo child", "started child", "review child":
					state := map[string]string{"done child": "done", "cancelled child": "cancelled", "todo child": "todo", "started child": "in_progress", "review child": "in_review"}[name]
					if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-child','task','FID-epic','Task','task-default',?,'developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`, state); err != nil {
						t.Fatal(err)
					}
				case "late open child":
					if _, err := v.conn.Exec(`WITH RECURSIVE n(x) AS (VALUES(1) UNION ALL SELECT x+1 FROM n WHERE x<101) INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) SELECT printf('FID-%06d',x),'task','FID-epic','Finished task','task-default','done','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z' FROM n`); err != nil {
						t.Fatal(err)
					}
					if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-zzzzzz','task','FID-epic','Last task','task-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`); err != nil {
						t.Fatal(err)
					}
				case "open subtask":
					for _, stmt := range []string{
						`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-child','task','FID-epic','Finished task','task-default','cancelled','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`,
						`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-sub','subtask','FID-child','Open subtask','subtask-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`,
					} {
						if _, err := v.conn.Exec(stmt); err != nil {
							t.Fatal(err)
						}
					}
				case "agent", "service":
					e.ActorV2.ActorType = name
					e.Effects[0].Values["actor"] = name + ":mark"
				case "blank reason":
					e.Effects[0].Values["reason"] = " \n "
				case "wrong actor":
					e.Effects[0].Values["actor"] = "human:other"
				case "wrong binding":
					e.Effects[0].Values["binding_id"] = "FID-other"
				case "wrong subject":
					e.Effects[0].Values["subject_id"] = "FID-other"
				case "wrong time":
					e.Effects[0].Values["decided_at"] = "2026-11-15T00:00:00Z"
				case "wrong kind":
					e.Effects[0].Values["kind"] = "cancel"
				case "missing decision":
					e.Effects = e.Effects[1:]
				case "missing move":
					e.Effects = e.Effects[:1]
				case "other move":
					e.Effects[1].PK["id"] = "FID-other"
				case "extra move":
					e.Effects = append(e.Effects, e.Effects[1])
				case "extra decision":
					e.Effects = append([]journal.Effect{e.Effects[0]}, e.Effects...)
				case "injected retry":
					e.Effects = append(e.Effects, journal.Effect{Op: "insert", Table: "idempotency_results"})
				case "wrong revision":
					e.Effects[1].Values["revision"] = 3
				case "wrong from":
					e.Effects[1].Preconditions["status"] = "in_progress"
				case "wrong target":
					e.Effects[1].Values["status"] = "cancelled"
				case "changed title":
					e.Effects[1].Values["title"] = "Changed"
				case "wrong command":
					e.CommandName = "issue.transition"
				case "wrong aggregate":
					e.Aggregate.Type = "issue"
				case "legacy":
					e.EventSchema = 1
				case "other family":
					e.EventType = "issue.transitioned"
					e.Aggregate.Type = "issue"
				case "case alias":
					e.Effects[0].Table = "PERSON_DECISIONS"
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, mode, e)
				if (err == nil) != want {
					t.Fatalf("admit=%v want=%v", err, want)
				}
				if !want && before != handoffProjectionSnapshot(t, v) {
					t.Fatal("refusal changed projection")
				}
			})
		}
	}
}

func TestEpicCloseStateAndDecisionOwnership(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, op := range []string{"update decision", "delete decision", "insert decision", "move without decision"} {
			t.Run(mode+"/"+op, func(t *testing.T) {
				t.Parallel()
				v, e := epicCloseArchiveFixture(t)
				if op != "move without decision" {
					acceptMergeArchiveEvent(t, v, e)
				}
				e.EventSchema = 1
				e.EventType = "issue.updated"
				e.EventID = "EVT-20261115-000002"
				e.Aggregate = journal.Aggregate{ID: "FID-epic", Type: "issue", Sequence: 1}
				effect := journal.Effect{Op: "update", Table: "PERSON_DECISIONS", PK: map[string]any{"ID": "PDN-123456"}, Values: map[string]any{"reason": "Changed"}}
				switch op {
				case "delete decision":
					effect.Op = "delete"
					effect.Values = nil
				case "insert decision":
					effect = e.Effects[0]
					effect.Table = "PERSON_DECISIONS"
					effect.PK = map[string]any{"id": "PDN-forged"}
					delete(effect.Values, "kind")
					effect.Values["KiNd"] = "close-epic"
				case "move without decision":
					effect = e.Effects[1]
					effect.Table = "ISSUES"
					delete(effect.Values, "status")
					effect.Values["StAtUs"] = "done"
				}
				e.Effects = []journal.Effect{effect}
				before := handoffProjectionSnapshot(t, v)
				if err := applyLeanEmptyEvent(v, mode, e); err == nil || before != handoffProjectionSnapshot(t, v) {
					t.Fatalf("%s accepted or changed projection: %v", op, err)
				}
			})
		}
	}
}
