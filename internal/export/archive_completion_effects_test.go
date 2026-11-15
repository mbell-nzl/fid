package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCompletionQuestionEffectsEveryPath(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"live", "replay", "import"} {
		for _, keyed := range []bool{false, true} {
			for _, name := range []string{"valid", "release substitutes question", "retry substitutes question", "other task question", "resolved question"} {
				label := mode + "/unkeyed/" + name
				if keyed {
					label = mode + "/keyed/" + name
				}
				t.Run(label, func(t *testing.T) {
					t.Parallel()
					v, e := completionArchiveFixture(t)
					for _, query := range []string{
						`INSERT INTO issues(id,kind,title,workflow_id,status,lane,owner_role,created_at,updated_at) VALUES('FID-b','task','Other task','task-default','in_review','routine','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`,
						`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at,ended_at,end_reason,kept) SELECT 'HLD-other','FID-b','other-worktree','/other',actor,since,expires_at,ended_at,end_reason,1 FROM holds WHERE id='HLD-1'`,
						`INSERT INTO questions(id,task_id,hold_id,text,asker,asked_at) VALUES('QST-open','FID-a','HLD-1','Question?','agent:builder','2026-11-10T09:00:00Z'),('QST-other','FID-b','HLD-other','Other question?','agent:builder','2026-11-10T09:00:00Z')`,
						`INSERT INTO questions(id,task_id,hold_id,text,asker,asked_at,resolved_at,resolution) VALUES('QST-resolved','FID-a','HLD-1','Earlier question?','agent:builder','2026-11-10T08:00:00Z','2026-11-10T08:30:00Z','cancelled')`,
					} {
						if _, err := v.conn.Exec(query); err != nil {
							t.Fatal(err)
						}
					}
					question := journal.Effect{Op: "update", Table: "questions", PK: map[string]any{"id": "QST-open"}, Values: map[string]any{"resolved_at": e.OccurredAt.Format(time.RFC3339), "resolution": "finished-before-answer", "decision_id": nil}, Preconditions: map[string]any{"resolved_at": nil}}
					switch name {
					case "release substitutes question":
						question = journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-other"}, Values: map[string]any{"kept": false}, Preconditions: map[string]any{"kept": true}}
					case "other task question":
						question.PK["id"] = "QST-other"
					case "resolved question":
						question.PK["id"] = "QST-resolved"
					}
					e.Effects = append(e.Effects, question)
					if name == "retry substitutes question" {
						e.Effects = e.Effects[:3]
						appendCompletionRetry(t, &e)
						e.Idempotency = nil
					}
					if keyed {
						appendCompletionRetry(t, &e)
					}
					before := handoffProjectionSnapshot(t, v)
					want := name == "valid"
					if mode == "live" {
						err := journal.DryRunProjectEvent(v.conn, e)
						if (err == nil) != want {
							t.Errorf("dry-run=%v want accepted=%t", err, want)
						}
						if handoffProjectionSnapshot(t, v) != before {
							t.Fatal("dry-run changed projection")
						}
					}
					err := applyLeanEmptyEvent(v, mode, e)
					if (err == nil) != want {
						t.Errorf("admit=%v want accepted=%t", err, want)
					}
					if !want && handoffProjectionSnapshot(t, v) != before {
						t.Fatal("refused completion changed projection")
					}
					if want {
						assertCompletionQuestionEffects(t, v)
					}
				})
			}
		}
	}
}

func appendCompletionRetry(t *testing.T, e *journal.Event) {
	t.Helper()
	var r completionrecord.Record
	if err := json.Unmarshal([]byte(e.Effects[0].Values["record_json"].(string)), &r); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(completionrecord.Result{Completion: r, Status: "done"})
	if err != nil {
		t.Fatal(err)
	}
	digest := "sha256:" + strings.Repeat("a", 64)
	e.Idempotency = &journal.EventIdempotency{Key: "done-1", RequestDigest: digest, Scope: "task.done:project-test"}
	e.Effects = append(e.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": e.Idempotency.Scope, "idempotency_key": e.Idempotency.Key}, Values: map[string]any{"request_digest": digest, "request_field_digests_json": `{"--actor":"` + digest + `"}`, "result_entity_id": e.Aggregate.ID, "result_event_id": e.EventID, "aggregate_sequence": e.Aggregate.Sequence, "result_json": string(raw), "created_at": r.At}})
}

func assertCompletionQuestionEffects(t *testing.T, v *semanticValidator) {
	t.Helper()
	var state, resolution string
	var kept, otherKept, otherOpen, completions int
	if err := v.conn.QueryRow(`SELECT (SELECT status FROM issues WHERE id='FID-a'),(SELECT kept FROM holds WHERE id='HLD-1'),(SELECT resolution FROM questions WHERE id='QST-open'),(SELECT kept FROM holds WHERE id='HLD-other'),(SELECT count(*) FROM questions WHERE id='QST-other' AND resolved_at IS NULL),(SELECT count(*) FROM completions WHERE task_id='FID-a')`).Scan(&state, &kept, &resolution, &otherKept, &otherOpen, &completions); err != nil {
		t.Fatal(err)
	}
	if state != "done" || kept != 0 || resolution != "finished-before-answer" || otherKept != 1 || otherOpen != 1 || completions != 1 {
		t.Fatalf("completion effects: state=%s kept=%d resolution=%s unrelated_kept=%d unrelated_open=%d completions=%d", state, kept, resolution, otherKept, otherOpen, completions)
	}
}
