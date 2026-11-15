// Package reviewwork records canonical review rounds under the journal lock.
package reviewwork

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, TaskID, Actor, IdempotencyKey string
	Input                                                       reviewrecord.Input
	Clock                                                       func() time.Time
}

type Result struct {
	Round    reviewrecord.Round `json:"round"`
	At       string             `json:"occurred_at"`
	Replayed bool               `json:"replayed"`
	EventID  string             `json:"-"`
}

func Record(o Options) (Result, error) {
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return Result{}, fmt.Errorf("review recording requires a work actor")
	}
	if o.IdempotencyKey == "" {
		return Result{}, fmt.Errorf("review recording requires an idempotency key")
	}
	if err := o.Input.Validate(); err != nil {
		return Result{}, err
	}
	var result Result
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid review record"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsDir); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		request, err := requestFor(o)
		if err != nil {
			return err
		}
		prior, found, err := journal.ReplayIdempotent(conn, request)
		if err != nil {
			return err
		}
		if found {
			result, err = decode(conn, prior.ResultJSON, prior.EventID, prior.EntityID)
			return err
		}
		h, err := reviewrecord.LatestHandoff(conn, o.TaskID)
		if err != nil {
			return err
		}
		binding, err := fingerprint.Compute(o.WorkspaceRoot, h.Base, h.Head)
		if err != nil {
			return err
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		id, err := journal.GenerateID("RND", func(id string) (bool, error) {
			var n int
			err := conn.QueryRow(`SELECT count(*) FROM review_rounds WHERE id=?`, id).Scan(&n)
			return n != 0, err
		})
		if err != nil {
			return err
		}
		in := o.Input
		r := reviewrecord.Round{ID: id, TaskID: o.TaskID, HandoffID: h.ID, Reviewer: in.Reviewer, Stream: in.Stream, Profiles: in.Profiles, Verdict: in.Verdict, Purpose: in.Purpose, Findings: []reviewrecord.Finding{}, Dispositions: in.Dispositions, Verification: in.Verification, Provenance: in.Provenance, Binding: reviewrecord.Bind(binding), RecordedBy: o.Actor, At: now.Format(time.RFC3339)}
		for i, f := range in.Findings {
			r.Findings = append(r.Findings, reviewrecord.Finding{ID: fmt.Sprintf("%s-%d", id, i+1), FindingInput: f})
		}
		if err := reviewrecord.Admit(conn, r); err != nil {
			return err
		}
		effects, err := Effects(r)
		if err != nil {
			return err
		}
		rework, err := journal.ReviewReworkEffects(conn, r)
		if err != nil {
			return err
		}
		effects = append(effects, rework...)
		kind, actor, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "review.round-recorded", CommandName: "review.record", OccurredAt: now, ActorV2: journal.EventActor{ActorID: actor, ActorType: kind, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: r.ID, Type: "review", Sequence: 1}, Effects: effects}
		raw, err := json.Marshal(Result{Round: r, At: r.At})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: r.ID, ResultJSON: string(raw)}, nil
		})
		if err != nil {
			return err
		}
		result, err = decode(conn, outcome.ResultJSON, outcome.EventID, outcome.EntityID)
		return err
	})
	return result, err
}

func decode(conn *sql.DB, raw, id, entity string) (Result, error) {
	var r Result
	if err := journal.ValidateReviewResult(raw); err != nil {
		return r, err
	}
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return r, err
	}
	if r.Round.ID != entity {
		return Result{}, fmt.Errorf("review retry has wrong round ID")
	}
	var rawRound string
	if err := conn.QueryRow(`SELECT record_json FROM review_rounds WHERE id=?`, entity).Scan(&rawRound); err != nil {
		return Result{}, err
	}
	var recorded reviewrecord.Round
	if err := json.Unmarshal([]byte(rawRound), &recorded); err != nil {
		return Result{}, err
	}
	if !reflect.DeepEqual(recorded, r.Round) {
		return Result{}, fmt.Errorf("review retry differs from recorded round")
	}
	r.EventID = id
	return r, nil
}

func requestFor(o Options) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	raw, err := json.Marshal(o.Input)
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	request := map[string]any{}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return journal.IdempotencyRequest{}, err
	}
	request["task"] = o.TaskID
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "review.record", "project_id": cfg.ProjectID, "request": request}, o.Actor, map[string]string{"task": "task"})
	return journal.IdempotencyRequest{OperationScope: "review.record:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}

// Effects keeps the round and its stable finding IDs in a single Event.
func Effects(r reviewrecord.Round) ([]journal.Effect, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	effects := []journal.Effect{{Op: "insert", Table: "review_rounds", PK: map[string]any{"id": r.ID}, Values: map[string]any{"task_id": r.TaskID, "handoff_id": r.HandoffID, "record_json": string(raw)}}}
	for _, f := range r.Findings {
		raw, err := json.Marshal(f)
		if err != nil {
			return nil, err
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "round_findings", PK: map[string]any{"id": f.ID}, Values: map[string]any{"round_id": r.ID, "record_json": string(raw)}})
	}
	for _, d := range r.Dispositions {
		if d.Closes() {
			effects = append(effects, ClosureEffect(d.FindingID, r.ID, "", d.Disposition, d.Reason))
		}
	}
	return effects, nil
}

func ClosureEffect(finding, round, decision, disposition, reason string) journal.Effect {
	var rnd, pdn any
	if round != "" {
		rnd = round
	}
	if decision != "" {
		pdn = decision
	}
	return journal.Effect{Op: "insert", Table: "finding_closures", PK: map[string]any{"finding_id": finding}, Values: map[string]any{"round_id": rnd, "decision_id": pdn, "disposition": disposition, "reason": reason}}
}

// OpenProjection is shared by read adapters; queries do not mutate records.
func OpenProjection(root string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
}
