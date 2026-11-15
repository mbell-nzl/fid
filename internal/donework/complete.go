package donework

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/journal"
)

func requestFor(o Options) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "task.done", "project_id": cfg.ProjectID, "request": map[string]any{"task": o.TaskID}}, o.Actor, map[string]string{"task": "task"})
	return journal.IdempotencyRequest{OperationScope: "task.done:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}
func retry(conn *sql.DB, request journal.IdempotencyRequest) (completionrecord.Result, bool, error) {
	if request.IdempotencyKey == "" {
		return completionrecord.Result{}, false, nil
	}
	prior, found, err := journal.ReplayIdempotent(conn, request)
	if err != nil || !found {
		return completionrecord.Result{}, found, err
	}
	r, err := journal.ReadCompletionRetry(conn, prior.ResultJSON, prior.EntityID)
	r.EventID = prior.EventID
	return r, true, err
}

func complete(lock *journal.WriteLock, conn *sql.DB, o Options, request journal.IdempotencyRequest, s doneeval.Snapshot, outcome doneeval.Outcome, now time.Time) (completionrecord.Result, error) {
	r := completionrecord.From(s, outcome)
	r.Actor, r.At = o.Actor, now.Format(time.RFC3339)
	id, err := journal.GenerateID("CMP", func(id string) (bool, error) {
		var n int
		err := conn.QueryRow(`SELECT count(*) FROM completions WHERE id=?`, id).Scan(&n)
		return n > 0, err
	})
	if err != nil {
		return completionrecord.Result{}, err
	}
	r.ID = id
	raw, err := json.Marshal(r)
	if err != nil {
		return completionrecord.Result{}, err
	}
	effects := []journal.Effect{
		{Op: "insert", Table: "completions", PK: map[string]any{"id": r.ID}, Values: map[string]any{"task_id": r.TaskID, "handoff_id": r.HandoffID, "merge_id": r.MergeID, "record_json": string(raw)}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": r.TaskID}, Values: map[string]any{"status": "done", "revision": s.Task.Revision + 1, "updated_at": r.At}, Preconditions: map[string]any{"status": "in_review", "revision": s.Task.Revision}},
		{Op: "update", Table: "holds", PK: map[string]any{"id": s.Handoff.HoldID}, Values: map[string]any{"kept": false}, Preconditions: map[string]any{"kept": true}},
	}
	rows, err := conn.Query(`SELECT id FROM questions WHERE task_id=? AND resolved_at IS NULL`, r.TaskID)
	if err != nil {
		return completionrecord.Result{}, err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return completionrecord.Result{}, err
		}
		effects = append(effects, journal.Effect{Op: "update", Table: "questions", PK: map[string]any{"id": id}, Values: map[string]any{"resolved_at": r.At, "resolution": "finished-before-answer", "decision_id": nil}, Preconditions: map[string]any{"resolved_at": nil}})
	}
	err = rows.Err()
	_ = rows.Close()
	if err != nil {
		return completionrecord.Result{}, err
	}
	actorType, actorID, _ := strings.Cut(o.Actor, ":")
	aggregate := journal.Aggregate{Type: "task", ID: r.TaskID}
	aggregate.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.JournalDirName), aggregate)
	if err != nil {
		return completionrecord.Result{}, err
	}
	event := journal.Event{EventSchema: 2, EventType: "task.finished", CommandName: "task.done", OccurredAt: now, Aggregate: aggregate, ActorV2: journal.EventActor{ActorType: actorType, ActorID: actorID, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Effects: effects}
	result := completionrecord.Result{Completion: r, Status: "done"}
	if o.IdempotencyKey == "" {
		appended, err := journal.CommitEvent(lock, conn, event)
		result.EventID = appended.EventID
		return result, err
	}
	raw, err = json.Marshal(result)
	if err != nil {
		return result, err
	}
	committed, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: r.TaskID, ResultJSON: string(raw)}, nil
	})
	result.EventID = committed.EventID
	return result, err
}
