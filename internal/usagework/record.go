// Package usagework records agent session output under the journal write lock.
package usagework

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/usagerecord"
)

type Options struct {
	WorkspaceRoot                                             string
	MigrationsFS                                              fs.FS
	HoldID, Actor, Format, Input, ToolVersion, IdempotencyKey string
	Clock                                                     func() time.Time
}
type Result struct {
	Usage    usagerecord.Record `json:"usage"`
	Replayed bool               `json:"replayed"`
}
type InputError struct{ Argument, Message string }

func (e *InputError) Error() string { return e.Message }

func Record(o Options) (Result, string, error) {
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return Result{}, "", &persondecision.DeniedError{Actor: o.Actor}
	}
	if o.IdempotencyKey == "" {
		return Result{}, "", &InputError{"--idempotency-key", "usage record needs an idempotency key"}
	}
	raw, err := os.ReadFile(o.Input)
	if err != nil {
		return Result{}, "", &InputError{"--input", fmt.Sprintf("read usage input: %v", err)}
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(raw))
	var result Result
	var eventID string
	err = journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid usage record"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		session, err := usagerecord.ReadSession(conn, o.HoldID)
		if err == sql.ErrNoRows {
			return &holdwork.Refusal{Code: "hold-not-current", Message: "usage session was not found", Details: map[string]any{"hold_id": o.HoldID, "reason": "missing"}}
		}
		if err != nil {
			return err
		}
		if session.Actor != o.Actor {
			return &persondecision.DeniedError{Actor: o.Actor}
		}
		request, err := requestFor(o, digest)
		if err != nil {
			return err
		}
		if previous, found, err := journal.ReplayIdempotent(conn, request); err != nil {
			return err
		} else if found {
			if err := journal.ValidateUsageRetry(conn, previous.ResultJSON, previous.EntityID); err != nil {
				return err
			}
			eventID = previous.EventID
			return json.Unmarshal([]byte(previous.ResultJSON), &result)
		}
		parsed, err := usagerecord.Parse(o.Format, raw, session.Model)
		if err != nil {
			argument := "--input"
			if o.Format != "codex" && o.Format != "claude" {
				argument = "--format"
			}
			return &InputError{argument, err.Error()}
		}
		if parsed.ToolVersion == "" && o.ToolVersion != "" {
			parsed.ToolVersion = o.ToolVersion
			parsed.VersionSource = "caller"
		}
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM usage_sessions WHERE hold_id=? AND input_digest=?`, o.HoldID, digest).Scan(&n); err != nil {
			return err
		}
		if n != 0 {
			return &holdwork.Refusal{Code: "duplicate-entity", TaskID: session.TaskID, Message: "this output is already recorded for this session", Details: map[string]any{"entity_type": "usage", "identity": map[string]any{"fields": []string{"hold_id", "input_digest"}, "values": []string{o.HoldID, digest}}}}
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		id, err := journal.GenerateID("USG", func(id string) (bool, error) {
			var n int
			err := conn.QueryRow(`SELECT count(*) FROM usage_sessions WHERE id=?`, id).Scan(&n)
			return n != 0, err
		})
		if err != nil {
			return err
		}
		r := usagerecord.Record{ID: id, TaskID: session.TaskID, HoldID: o.HoldID, Actor: o.Actor, Model: session.Model, Role: session.Role, Stream: session.Stream, InputDigest: digest, RecordedAt: now.Format(time.RFC3339), Parsed: parsed}
		effect, err := journal.UsageEffect(r)
		if err != nil {
			return err
		}
		actorType, actorID, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "usage.recorded", CommandName: "usage.record", OccurredAt: now, ActorV2: journal.EventActor{ActorType: actorType, ActorID: actorID, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "usage", ID: id, Sequence: 1}, Effects: []journal.Effect{effect}}
		result = Result{Usage: r}
		body, err := journal.MarshalCanonicalJSON(result)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: id, ResultJSON: string(body)}, nil
		})
		if err != nil {
			return err
		}
		eventID = outcome.EventID
		return json.Unmarshal([]byte(outcome.ResultJSON), &result)
	})
	if err != nil {
		return Result{}, "", fmt.Errorf("usage record: %w", err)
	}
	return result, eventID, nil
}
func requestFor(o Options, digest string) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	value, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "usage.record", "project_id": cfg.ProjectID, "request": map[string]any{"hold": o.HoldID, "format": o.Format, "input_digest": digest, "tool_version": o.ToolVersion}}, o.Actor, map[string]string{"hold": "--hold", "format": "--format", "input_digest": "--input", "tool_version": "--tool-version"})
	return journal.IdempotencyRequest{OperationScope: "usage.record:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: value, RequestFieldDigests: fields}, err
}
