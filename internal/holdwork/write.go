// Package holdwork implements session commands on lean holds under the
// shared journal lock. Every command names the exact hold or takeover task.
package holdwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, HoldID, TaskID, Actor, Reason, IdempotencyKey string
	Clock                                                                       func() time.Time
}

type Result struct {
	Fields              map[string]any
	EventID, OccurredAt string
}

type Refusal struct {
	Code, TaskID, Message string
	Details               map[string]any
}

func (err *Refusal) Error() string { return err.Message }

type writeContext struct {
	options Options
	lock    *journal.WriteLock
	conn    *sql.DB
	now     time.Time
}

func write(options Options, operation string, build func(writeContext) (journal.Event, map[string]any, error)) (Result, error) {
	set := policy.WorkSet
	if operation == "task.takeover" {
		set = policy.PersonDecisions
	}
	if !policy.Allows(options.Actor, set) {
		return Result{}, &persondecision.DeniedError{Actor: options.Actor}
	}
	var result Result
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + operation}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		var request journal.IdempotencyRequest
		if options.IdempotencyKey != "" {
			request, err = idempotencyRequest(options, operation)
			if err != nil {
				return err
			}
			outcome, found, err := journal.ReplayIdempotent(conn, request)
			if err != nil {
				return err
			}
			if found {
				result, err = decodeResult(outcome.ResultJSON, outcome.EventID)
				return err
			}
		}
		now := time.Now()
		if options.Clock != nil {
			now = options.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		event, fields, err := build(writeContext{options, lock, conn, now})
		if err != nil {
			return err
		}
		actorType, actorID, _ := strings.Cut(options.Actor, ":")
		event.EventSchema = journal.EventSchemaV2
		event.CommandName = operation
		event.OccurredAt = now
		event.ActorV2 = journal.EventActor{ActorID: actorID, ActorType: actorType, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
		event.Aggregate.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(options.WorkspaceRoot), "journal"), event.Aggregate)
		if err != nil {
			return err
		}
		if options.IdempotencyKey == "" {
			appended, err := journal.CommitEvent(lock, conn, event)
			if err != nil {
				return err
			}
			result = Result{fields, appended.EventID, holdsTimestamp(now)}
			return nil
		}
		raw, err := journal.MarshalCanonicalJSON(map[string]any{"fields": fields, "occurred_at": holdsTimestamp(now), "replayed": false})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: event.Aggregate.ID, ResultJSON: string(raw)}, nil
		})
		if err != nil {
			return err
		}
		result, err = decodeResult(outcome.ResultJSON, outcome.EventID)
		return err
	})
	if err != nil {
		return Result{}, fmt.Errorf("%s: %w", operation, err)
	}
	return result, nil
}

func holdsTimestamp(now time.Time) string { return now.UTC().Format(time.RFC3339) }

func decodeResult(raw, eventID string) (Result, error) {
	var stored struct {
		Fields     map[string]any `json:"fields"`
		OccurredAt string         `json:"occurred_at"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return Result{}, fmt.Errorf("decode hold result: %w", err)
	}
	return Result{stored.Fields, eventID, stored.OccurredAt}, nil
}

func idempotencyRequest(options Options, operation string) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	request := map[string]any{"reason": options.Reason}
	names := map[string]string{"reason": "--reason", "hold": "--hold", "task": "task"}
	if operation == "hold.park" {
		request["hold"] = options.HoldID
		request["mode"] = "interrupted"
		names["mode"] = "--interrupted"
	} else {
		request["task"] = options.TaskID
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": options.Actor, "operation": operation, "project_id": cfg.ProjectID, "request": request}, options.Actor, names)
	return journal.IdempotencyRequest{OperationScope: operation + ":" + cfg.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}
