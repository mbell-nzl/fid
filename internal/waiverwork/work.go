// Package waiverwork records a person's waiver under the journal lock.
package waiverwork

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, TaskID, Actor, Reason, Lane, IdempotencyKey string
	Clock                                                                     func() time.Time
	Confirm                                                                   func(Preview) (bool, error)
}

type Cancelled struct{}

func (*Cancelled) Error() string { return "Waiver cancelled; nothing recorded" }

func Record(o Options) (waiverrecord.Result, error) {
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return waiverrecord.Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return waiverrecord.Result{}, &persondecision.ReasonRequiredError{}
	}
	var result waiverrecord.Result
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid waive"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsDir); err != nil {
			return err
		}
		conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		request, err := requestFor(o)
		if err != nil {
			return err
		}
		if o.IdempotencyKey != "" {
			prior, found, err := journal.ReplayIdempotent(conn, request)
			if err != nil {
				return err
			}
			if found {
				result, err = journal.ReadWaiverRetry(conn, prior.ResultJSON, prior.EntityID)
				if err == nil && (result.Decision.TaskID != o.TaskID || result.Decision.Actor != o.Actor || result.Decision.Reason != o.Reason || result.Decision.Binding.Lane != o.Lane) {
					return fmt.Errorf("waiver retry differs from requested task, actor, reason or lane")
				}
				result.EventID = prior.EventID
				return err
			}
		}
		h, floor, err := waiverrecord.Context(conn, o.TaskID, o.Lane)
		if err != nil {
			return err
		}
		fp, err := fingerprint.Compute(o.WorkspaceRoot, h.Base, h.Head)
		if err != nil {
			var failure *fingerprint.Failure
			if errors.As(err, &failure) && failure.Kind == fingerprint.MissingObject {
				return fmt.Errorf("%w; restore the recorded base and head objects locally, then retry", err)
			}
			return err
		}
		preview := Preview{TaskID: o.TaskID, Handoff: h, Floor: floor, Lane: o.Lane}
		if o.Confirm == nil {
			return fmt.Errorf("waiver requires interactive confirmation")
		}
		yes, err := o.Confirm(preview)
		if err != nil {
			return err
		}
		if !yes {
			return &Cancelled{}
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		id, effect, err := persondecision.Effect(conn, "waive", o.Actor, "task", o.TaskID, o.Reason, now)
		if err != nil {
			return err
		}
		r := waiverrecord.Record{ID: id, TaskID: o.TaskID, Actor: o.Actor, Reason: o.Reason, At: now.Format(time.RFC3339), Binding: waiverrecord.Binding{Lane: o.Lane, Algorithm: fp.Algorithm, Fingerprint: fp.Fingerprint, HandoffID: h.ID, Base: h.Base, Head: h.Head}}
		if err := waiverrecord.Admit(conn, r); err != nil {
			return err
		}
		raw, err := json.Marshal(r.Binding)
		if err != nil {
			return err
		}
		effect.Values["binding_id"] = fp.Fingerprint
		effect.Values["binding_json"] = string(raw)
		actorType, actorID, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: 2, EventType: "task.waived", CommandName: "task.waive", OccurredAt: now, ActorV2: journal.EventActor{ActorType: actorType, ActorID: actorID, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "decision", ID: id, Sequence: 1}, Effects: []journal.Effect{effect}}
		result = waiverrecord.Result{Decision: r, At: r.At}
		if o.IdempotencyKey == "" {
			appended, err := journal.CommitEvent(lock, conn, event)
			result.EventID = appended.EventID
			return err
		}
		raw, err = json.Marshal(result)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: id, ResultJSON: string(raw)}, nil
		})
		result.EventID = outcome.EventID
		return err
	})
	return result, err
}

func requestFor(o Options) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "task.waive", "project_id": cfg.ProjectID, "request": map[string]any{"task": o.TaskID, "reason": o.Reason, "lane": o.Lane}}, o.Actor, map[string]string{"task": "task", "reason": "--reason", "lane": "--lane"})
	return journal.IdempotencyRequest{OperationScope: "task.waive:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}
