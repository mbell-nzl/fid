// Package approvalwork records a person's approval under the journal lock.
package approvalwork

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/fingerprint"
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
	Confirm                               func(Preview) (bool, error)
}

type Cancelled struct{}

func (*Cancelled) Error() string { return "Approval cancelled; nothing recorded" }

func Record(o Options) (approvalrecord.Result, error) {
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return approvalrecord.Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return approvalrecord.Result{}, &persondecision.ReasonRequiredError{}
	}
	var result approvalrecord.Result
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid approve"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
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
				result, err = journal.ReadApprovalRetry(conn, prior.ResultJSON, prior.EntityID)
				if err == nil && (result.Decision.TaskID != o.TaskID || result.Decision.Actor != o.Actor || result.Decision.Reason != o.Reason) {
					return fmt.Errorf("approval retry differs from requested task, actor or reason")
				}
				result.EventID = prior.EventID
				return err
			}
		}
		var kind, state string
		if err := conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, o.TaskID).Scan(&kind, &state); err != nil {
			return err
		}
		if kind != "task" || state != "in_review" {
			return fmt.Errorf("approval task must be in_review")
		}
		h, err := reviewrecord.LatestHandoff(conn, o.TaskID)
		if err != nil {
			return err
		}
		fp, stat, err := fingerprint.Inspect(o.WorkspaceRoot, h.Base, h.Head)
		if err != nil {
			var failure *fingerprint.Failure
			if errors.As(err, &failure) && failure.Kind == fingerprint.MissingObject {
				return fmt.Errorf("%w; restore the recorded base and head objects locally, then retry", err)
			}
			return err
		}
		preview, err := previewFor(conn, o.TaskID, h, fp, stat)
		if err != nil {
			return err
		}
		if o.Confirm == nil {
			return fmt.Errorf("approval requires interactive confirmation")
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
		id, effect, err := persondecision.Effect(conn, "approve", o.Actor, "task", o.TaskID, o.Reason, now)
		if err != nil {
			return err
		}
		r := approvalrecord.Record{ID: id, TaskID: o.TaskID, Actor: o.Actor, Reason: o.Reason, At: now.Format(time.RFC3339), Binding: approvalrecord.Binding{Algorithm: fp.Algorithm, Fingerprint: fp.Fingerprint, HandoffID: h.ID, Base: h.Base, Head: h.Head}}
		if err := approvalrecord.Admit(conn, r); err != nil {
			return err
		}
		raw, err := json.Marshal(r.Binding)
		if err != nil {
			return err
		}
		effect.Values["binding_id"] = fp.Fingerprint
		effect.Values["binding_json"] = string(raw)
		actorType, actorID, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: 2, EventType: "task.approved", CommandName: "task.approve", OccurredAt: now, ActorV2: journal.EventActor{ActorType: actorType, ActorID: actorID, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "decision", ID: id, Sequence: 1}, Effects: []journal.Effect{effect}}
		result = approvalrecord.Result{Decision: r, At: r.At}
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
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": "task.approve", "project_id": cfg.ProjectID, "request": map[string]any{"task": o.TaskID, "reason": o.Reason}}, o.Actor, map[string]string{"task": "task", "reason": "--reason"})
	return journal.IdempotencyRequest{OperationScope: "task.approve:" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}
