package export

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func workDecisionFixture(t *testing.T) (*semanticValidator, journal.Event) {
	t.Helper()
	v, err := newSemanticValidator(testMigrationsFS(t), 37)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = v.close() })
	if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,lane) VALUES('FID-a','task','Work','task-default','todo','standard'),('FID-b','task','Other work','task-default','todo','standard')`); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 22, 1, 0, 0, 0, time.UTC)
	return v, journal.Event{EventSchema: 2, EventID: "EVT-20261122-000001", EventType: "work.accepted", CommandName: "task.accept", OccurredAt: at, ActorV2: journal.EventActor{ActorType: "human", ActorID: "mark", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{Type: "issue", ID: "FID-a", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-accept"}, Values: map[string]any{"kind": "accept", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": "FID-a", "reason": "Approved work", "decided_at": at.Format(time.RFC3339)}}}}
}

func TestWorkDecisionAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, name := range []string{"valid", "agent", "service", "actor", "subject", "binding", "kind", "time", "blank reason", "command", "aggregate", "extra decision", "missing decision", "extra column", "legacy", "other family", "creator forgery", "closed subject", "cancelled subject", "other task decision"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := workDecisionFixture(t)
				d := e.Effects[0].Values
				switch name {
				case "closed subject", "cancelled subject":
					state := map[string]string{"closed subject": "done", "cancelled subject": "cancelled"}[name]
					if _, err := v.conn.Exec(`UPDATE issues SET status=? WHERE id='FID-a'`, state); err != nil {
						t.Fatal(err)
					}
				case "agent", "service":
					e.ActorV2.ActorType = name
					d["actor"] = name + ":mark"
				case "actor":
					d["actor"] = "human:other"
				case "other task decision":
					d["subject_id"], d["binding_id"] = "FID-b", "FID-b"
				case "subject":
					d["subject_id"] = "FID-other"
				case "binding":
					d["binding_id"] = "FID-other"
				case "kind":
					d["kind"] = "answer"
				case "time":
					d["decided_at"] = "2026-11-22T01:00:01Z"
				case "blank reason":
					d["reason"] = "\t"
				case "command":
					e.CommandName = "issue.update"
				case "aggregate":
					e.Aggregate.ID = "FID-other"
				case "extra decision":
					extra := e.Effects[0]
					extra.PK = map[string]any{"id": "PDN-extra"}
					e.Effects = append(e.Effects, extra)
				case "missing decision":
					e.Effects = nil
				case "extra column":
					d["binding_json"] = "{}"
				case "legacy":
					e.EventSchema = 1
				case "other family":
					e.EventType = "task.answered"
					e.Aggregate.Type = "task"
					e.CommandName = "task.answer"
				case "creator forgery":
					e.Effects = append(e.Effects, journal.Effect{Table: "issue_creators", Op: "insert", PK: map[string]any{"issue_id": "FID-a"}, Values: map[string]any{"actor": "human:mark"}})
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyWorkDecisionEvent(t, v, mode, e)
				if (err == nil) != (name == "valid") {
					t.Fatalf("%s admission=%v", name, err)
				}
				if name != "valid" {
					if after := handoffProjectionSnapshot(t, v); after != before {
						t.Fatal("refusal changed projection")
					}
					var count int
					if err := v.conn.QueryRow(`SELECT count(*) FROM person_decisions`).Scan(&count); err != nil || count != 0 {
						t.Fatalf("refusal changed decisions: %d/%v", count, err)
					}
					_, valid := workDecisionFixture(t)
					valid.Aggregate.ID = "FID-b"
					valid.Effects[0].Values["subject_id"], valid.Effects[0].Values["binding_id"] = "FID-b", "FID-b"
					// Failed import validation retains no authority; a later valid Event works.
					if err := applyWorkDecisionEvent(t, v, mode, valid); err != nil {
						t.Fatalf("valid after refusal: %v", err)
					}
				}
			})
		}
	}
}

func TestReworkAllowanceAdmissionEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, name := range []string{"valid", "agent", "subject", "binding", "kind", "time", "blank reason", "command", "aggregate", "duplicate", "zero", "negative", "fraction", "missing rounds", "extra member", "missing binding", "null binding", "wrong subject kind", "closed subject", "cancelled subject"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				v, e := workDecisionFixture(t)
				e.EventType, e.CommandName, e.Aggregate.Type = "rework.allowed", "task.rework-allow", "task"
				d := e.Effects[0].Values
				d["kind"], d["binding_json"] = "rework-allowance", `{"rounds":2}`
				switch name {
				case "closed subject", "cancelled subject":
					state := map[string]string{"closed subject": "done", "cancelled subject": "cancelled"}[name]
					if _, err := v.conn.Exec(`UPDATE issues SET status=? WHERE id='FID-a'`, state); err != nil {
						t.Fatal(err)
					}
				case "agent":
					e.ActorV2.ActorType = "agent"
					d["actor"] = "agent:mark"
				case "subject":
					d["subject_id"] = "FID-other"
				case "binding":
					d["binding_id"] = "patch-fingerprint"
				case "kind":
					d["kind"] = "accept"
				case "time":
					d["decided_at"] = "2026-11-22T01:00:01Z"
				case "blank reason":
					d["reason"] = "\t"
				case "command":
					e.CommandName = "task.accept"
				case "aggregate":
					e.Aggregate.ID = "FID-other"
				case "duplicate":
					extra := e.Effects[0]
					extra.PK = map[string]any{"id": "PDN-extra"}
					e.Effects = append(e.Effects, extra)
				case "zero":
					d["binding_json"] = `{"rounds":0}`
				case "negative":
					d["binding_json"] = `{"rounds":-1}`
				case "fraction":
					d["binding_json"] = `{"rounds":1.5}`
				case "missing rounds":
					d["binding_json"] = `{}`
				case "extra member":
					d["binding_json"] = `{"rounds":1,"other":true}`
				case "null binding":
					d["binding_json"] = nil
				case "missing binding":
					delete(d, "binding_json")
				case "wrong subject kind":
					d["subject_type"] = "epic"
				}
				before := handoffProjectionSnapshot(t, v)
				err := applyWorkDecisionEvent(t, v, mode, e)
				if (err == nil) != (name == "valid") {
					t.Fatalf("admission=%v", err)
				}
				if name != "valid" {
					if handoffProjectionSnapshot(t, v) != before {
						t.Fatal("refusal changed projection")
					}
					var n int
					if err := v.conn.QueryRow(`SELECT count(*) FROM person_decisions`).Scan(&n); err != nil || n != 0 {
						t.Fatal(n, err)
					}
					// A later valid write is still admitted after this malformed allowance.
					_, valid := workDecisionFixture(t)
					valid.Aggregate.ID = "FID-b"
					valid.Effects[0].Values["subject_id"], valid.Effects[0].Values["binding_id"] = "FID-b", "FID-b"
					valid.EventType, valid.CommandName, valid.Aggregate.Type = "rework.allowed", "task.rework-allow", "task"
					valid.Effects[0].Values["kind"], valid.Effects[0].Values["binding_json"] = "rework-allowance", `{"rounds":1}`
					if err := applyWorkDecisionEvent(t, v, mode, valid); err != nil {
						t.Fatal(err)
					}
				}
			})
		}
	}
}

func applyWorkDecisionEvent(t *testing.T, v *semanticValidator, mode string, e journal.Event) (err error) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("admission panicked instead of refusing: %v", recovered)
		}
	}()
	return applyReworkEvent(v, mode, e)
}
