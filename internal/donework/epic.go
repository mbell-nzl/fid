package donework

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/epicrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func CloseEpic(o Options) (epicrecord.Result, error) { return closeEpic(o, nil) }
func closeEpic(o Options, afterRead func()) (epicrecord.Result, error) {
	var result epicrecord.Result
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return result, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return result, &persondecision.ReasonRequiredError{}
	}
	if _, err := journal.EnsureProjectionFresh(o.WorkspaceRoot, o.MigrationsFS); err != nil {
		return result, err
	}
	conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
	if err != nil {
		return result, err
	}
	defer func() { _ = conn.Close() }()
	request, err := epicRequest(o)
	if err != nil {
		return result, err
	}
	if r, found, err := epicRetry(conn, request); err != nil || found {
		return r, err
	}
	if _, err := epicrecord.ReadEligibility(conn, o.TaskID); err != nil {
		return result, err
	}
	if afterRead != nil {
		afterRead()
	}
	err = journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid done"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
			return err
		}
		if r, found, err := epicRetry(conn, request); err != nil || found {
			result = r
			return err
		}
		eligible, err := epicrecord.ReadEligibility(conn, o.TaskID)
		if err != nil {
			return err
		}
		if !eligible.CanClose() {
			code, reason, message := "done-epic-not-open", "epic-not-open", "Only an open epic can be closed"
			if eligible.Kind == "epic" && (eligible.State == "todo" || eligible.State == "in_progress") {
				code, reason, message = "done-epic-open-tasks", "open-tasks", "Finish or cancel every task before closing the epic"
			}
			return &Refusal{Code: code, Message: message, Next: "fid show " + o.TaskID, Details: map[string]any{"gate": "epic", "reason": reason, "record_ids": eligible.OpenTasks()}}
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		id, decision, err := persondecision.Effect(conn, "close-epic", o.Actor, "epic", o.TaskID, o.Reason, now)
		if err != nil {
			return err
		}
		result = epicrecord.Result{Status: "done", Revision: eligible.Revision + 1, Decision: epicrecord.Decision{ID: id, Kind: "close-epic", Actor: o.Actor, SubjectType: "epic", SubjectID: o.TaskID, BindingID: o.TaskID, Reason: o.Reason, At: now.Format(time.RFC3339)}}
		aggregate := journal.Aggregate{Type: "epic", ID: o.TaskID}
		aggregate.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.JournalDirName), aggregate)
		if err != nil {
			return err
		}
		typ, actor, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: 2, EventType: "epic.closed", CommandName: "task.done", OccurredAt: now, Aggregate: aggregate, ActorV2: journal.EventActor{ActorType: typ, ActorID: actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Effects: []journal.Effect{decision, {Op: "update", Table: "issues", PK: map[string]any{"id": o.TaskID}, Values: map[string]any{"status": "done", "revision": result.Revision, "updated_at": result.Decision.At}, Preconditions: map[string]any{"status": eligible.State, "revision": eligible.Revision}}}}
		if o.IdempotencyKey == "" {
			appended, err := journal.CommitEvent(lock, conn, event)
			result.EventID = appended.EventID
			return err
		}
		raw, err := json.Marshal(result)
		if err != nil {
			return err
		}
		appended, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: o.TaskID, ResultJSON: string(raw)}, nil
		})
		result.EventID = appended.EventID
		return err
	})
	return result, err
}
func epicRequest(o Options) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "task.done", "project_id": cfg.ProjectID, "request": map[string]any{"epic": o.TaskID, "reason": o.Reason}}, o.Actor, map[string]string{"epic": "epic", "reason": "--reason"})
	return journal.IdempotencyRequest{OperationScope: "task.done:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}
func epicRetry(conn *sql.DB, request journal.IdempotencyRequest) (epicrecord.Result, bool, error) {
	if request.IdempotencyKey == "" {
		return epicrecord.Result{}, false, nil
	}
	prior, found, err := journal.ReplayIdempotent(conn, request)
	if err != nil || !found {
		return epicrecord.Result{}, found, err
	}
	var r epicrecord.Result
	if err := json.Unmarshal([]byte(prior.ResultJSON), &r); err != nil {
		return r, true, err
	}
	d, err := epicrecord.ReadDecision(conn, r.Decision.ID)
	if err == nil && (r.Status != "done" || r.Decision != d || prior.EntityID != d.SubjectID) {
		err = fmt.Errorf("epic retry differs from immutable decision")
	}
	r.EventID = prior.EventID
	return r, true, err
}

// EpicTarget consults a retry's stored result before the requested Issue. Each
// writer still validates the full request, so a changed ID under an existing
// key conflicts even when the newly requested Issue does not exist.
func EpicTarget(o Options) (bool, error) {
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return false, &persondecision.DeniedError{Actor: o.Actor}
	}
	if _, err := journal.EnsureProjectionFresh(o.WorkspaceRoot, o.MigrationsFS); err != nil {
		return false, err
	}
	conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
	if err != nil {
		return false, err
	}
	defer func() { _ = conn.Close() }()
	if o.IdempotencyKey != "" {
		request, err := requestFor(o)
		if err != nil {
			return false, err
		}
		var raw string
		err = conn.QueryRow(`SELECT result_json FROM idempotency_results WHERE operation_scope=? AND idempotency_key=?`, request.OperationScope, o.IdempotencyKey).Scan(&raw)
		if err == nil {
			var stored struct {
				Decision epicrecord.Decision `json:"decision"`
			}
			if err := json.Unmarshal([]byte(raw), &stored); err != nil {
				return false, fmt.Errorf("read done retry kind: %w", err)
			}
			return stored.Decision.Kind == "close-epic", nil
		}
		if err != sql.ErrNoRows {
			return false, err
		}
	}
	var kind string
	err = conn.QueryRow(`SELECT kind FROM issues WHERE id=?`, o.TaskID).Scan(&kind)
	if err == sql.ErrNoRows {
		return false, nil
	} // Finish owns the existing not-found refusal.
	return kind == "epic", err
}
