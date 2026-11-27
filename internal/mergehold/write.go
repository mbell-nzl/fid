// Package mergehold records a person's hold without moving or returning work.
package mergehold

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

type Options struct {
	WorkspaceRoot                         string
	MigrationsFS                          fs.FS
	TaskID, Actor, Reason, IdempotencyKey string
	Clock                                 func() time.Time
}
type Result struct {
	Fields   map[string]any `json:"fields"`
	Replayed bool           `json:"replayed"`
	EventID  string         `json:"-"`
}

func Record(o Options) (Result, error) {
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return Result{}, &persondecision.ReasonRequiredError{}
	}
	var result Result
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid hold"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
			return err
		}
		conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
		if err != nil {
			return err
		}
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "task.hold", "project_id": cfg.ProjectID, "request": map[string]any{"task": o.TaskID, "reason": o.Reason}}, o.Actor, map[string]string{"task": "task", "reason": "--reason"})
		if err != nil {
			return err
		}
		req := journal.IdempotencyRequest{OperationScope: "task.hold:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if o.IdempotencyKey != "" {
			prior, found, err := journal.ReplayIdempotent(conn, req)
			if err != nil {
				return err
			}
			if found {
				result.Fields, err = journal.ReadMergeHoldRetry(conn, prior.ResultJSON, prior.EntityID)
				if err != nil {
					return err
				}
				if result.Fields["task_id"] != o.TaskID || result.Fields["actor"] != o.Actor || result.Fields["reason"] != o.Reason {
					return fmt.Errorf("hold retry differs from requested task, actor or reason")
				}
				result.EventID = prior.EventID
				return nil
			}
		}
		var kind, state string
		if err := conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, o.TaskID).Scan(&kind, &state); err != nil {
			return err
		}
		if kind != "task" || state != "in_review" {
			return fmt.Errorf("hold task must be in_review")
		}
		h, err := reviewrecord.LatestHandoff(conn, o.TaskID)
		if err != nil {
			return err
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		at := now.Format(time.RFC3339)
		id, err := journal.GenerateID("PDN", func(id string) (bool, error) {
			var n int
			err := conn.QueryRow(`SELECT (SELECT count(*) FROM merge_holds WHERE id=?)+(SELECT count(*) FROM person_decisions WHERE id=?)`, id, id).Scan(&n)
			return n > 0, err
		})
		if err != nil {
			return err
		}
		result.Fields = map[string]any{"decision_id": id, "task_id": o.TaskID, "handoff_id": h.ID, "actor": o.Actor, "reason": o.Reason, "at": at, "verdict": "HOLD"}
		typ, actor, _ := strings.Cut(o.Actor, ":")
		e := journal.Event{EventSchema: 2, EventType: "task.held", CommandName: "task.hold", OccurredAt: now, ActorV2: journal.EventActor{ActorType: typ, ActorID: actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "decision", ID: id, Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "merge_holds", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": o.TaskID, "handoff_id": h.ID, "actor": o.Actor, "reason": o.Reason, "at": at}}}}
		if o.IdempotencyKey == "" {
			appended, err := journal.CommitEvent(lock, conn, e)
			result.EventID = appended.EventID
			return err
		}
		raw, err := json.Marshal(result)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, req, e, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: id, ResultJSON: string(raw)}, nil
		})
		result.EventID = outcome.EventID
		return err
	})
	return result, err
}
