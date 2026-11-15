// Package donework observes GitHub outside the lock and commits task completion.
package donework

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergework"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/policysnapshot"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, TaskID, Actor, IdempotencyKey string
	Clock                                                       func() time.Time
	Runner                                                      github.Runner
}

type Refusal struct {
	Code, Message, Next string
	Details             map[string]any
}

func (r *Refusal) Error() string { return r.Message }

type observed struct {
	snapshot    doneeval.Snapshot
	sessions    string
	reservation bool
}

func Finish(ctx context.Context, o Options) (completionrecord.Result, error) {
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return completionrecord.Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if _, err := journal.EnsureProjectionFresh(o.WorkspaceRoot, o.MigrationsDir); err != nil {
		return completionrecord.Result{}, err
	}
	conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
	if err != nil {
		return completionrecord.Result{}, err
	}
	defer func() { _ = conn.Close() }()
	request, err := requestFor(o)
	if err != nil {
		return completionrecord.Result{}, err
	}
	if r, found, err := retry(conn, request); err != nil || found {
		return r, err
	}
	before, err := readObserved(o, conn, doneeval.Target{})
	if err != nil {
		return completionrecord.Result{}, err
	}
	s := before.snapshot
	if s.Task.ID == "" {
		return completionrecord.Result{}, &Refusal{Code: "issue-not-found", Message: "task not found", Details: map[string]any{"id": o.TaskID}}
	}
	if err := preflight(s); err != nil {
		return completionrecord.Result{}, err
	}
	if !before.reservation {
		return completionrecord.Result{}, refusal(doneeval.HandoffGate, string(doneeval.HandoffMismatch), nil)
	}
	// Current policy controls authority, while the hand-off freezes every gate.
	if _, err := policy.Load(o.WorkspaceRoot); err != nil {
		return completionrecord.Result{}, err
	}
	target, err := localTarget(ctx, o.WorkspaceRoot)
	if err != nil {
		return completionrecord.Result{}, err
	}
	remote, remoteErr := observe(ctx, o, target, s.Handoff.PR)
	var result completionrecord.Result
	err = journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid done"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsDir); err != nil {
			return err
		}
		if r, found, err := retry(conn, request); err != nil || found {
			result = r
			return err
		}
		after, err := readObserved(o, conn, doneeval.Target{})
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(before, after) {
			return refusal(doneeval.TaskGate, "observation-stale", nil)
		}
		// Origin metadata may also change during the external observation.
		current, err := localTarget(ctx, o.WorkspaceRoot)
		if err != nil {
			return err
		}
		if current != target {
			return refusal(doneeval.PRGate, "observation-stale", nil)
		}
		loaded, err := policy.Load(o.WorkspaceRoot)
		if err != nil {
			return err
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		s = after.snapshot
		s.Target, s.GitHub = target, remote
		// A merge is independently confirmed only once the exact head tree matches.
		if remote.PR != nil && remote.PR.Merged && remote.HeadTree == s.Handoff.Tree {
			if err := policysnapshot.Ensure(lock, conn, loaded, now); err != nil {
				return err
			}
			opts := mergework.Options{WorkspaceRoot: o.WorkspaceRoot, MigrationsDir: o.MigrationsDir, TaskID: o.TaskID, Actor: o.Actor, Clock: func() time.Time { return now }}
			m, err := mergework.RecordMergeLocked(lock, conn, opts, s.Handoff.ID, *remote.PR, remote.HeadTree)
			if err != nil {
				return err
			}
			s.Merge = &m
			if remote.CI != nil {
				if _, err := mergework.RecordCILocked(lock, conn, opts, m.ID, target.Repository, *remote.CI); err != nil {
					return err
				}
			}
			s.CIResults, err = mergework.CIResults(conn, o.TaskID, m.ID)
			if err != nil {
				return err
			}
		}
		if remoteErr != nil {
			return observationFailure(remoteErr, s.Merge != nil)
		}
		outcome := doneeval.Evaluate(s)
		if outcome.Status != doneeval.Complete {
			return outcomeRefusal(outcome)
		}
		if err := policysnapshot.Ensure(lock, conn, loaded, now); err != nil {
			return err
		}
		result, err = complete(lock, conn, o, request, s, outcome, now)
		return err
	})
	return result, err
}

func readObserved(o Options, conn *sql.DB, target doneeval.Target) (observed, error) {
	tx, err := conn.Begin()
	if err != nil {
		return observed{}, err
	}
	defer func() { _ = tx.Rollback() }()
	s, err := doneeval.Gather(o.WorkspaceRoot, tx, o.TaskID, target, doneeval.GitHubFacts{})
	if err != nil {
		var unavailable *fingerprint.Failure
		if errors.As(err, &unavailable) {
			return observed{}, refusal(doneeval.HandoffGate, "fingerprint-unavailable", nil)
		}
		return observed{}, fmt.Errorf("gather completion: %w", err)
	}
	// Takeovers do not bump the Issue revision. Include every task hold and
	// question so even a takeover followed by another hand-off cannot hide.
	var sessions string
	err = tx.QueryRow(`SELECT json_object('holds',(SELECT json_group_array(json_array(id,ended_at,end_reason,kept,last_heartbeat_at,rework_at)) FROM holds WHERE task_id=?),'questions',(SELECT json_group_array(json_array(id,resolved_at,resolution,decision_id)) FROM questions WHERE task_id=?))`, o.TaskID, o.TaskID).Scan(&sessions)
	if err != nil {
		return observed{}, err
	}
	reserved := false
	if s.Handoff != nil {
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM holds WHERE id=? AND task_id=? AND kept=1 AND ended_at IS NOT NULL AND end_reason='review')`, s.Handoff.HoldID, o.TaskID).Scan(&reserved)
	}
	return observed{s, sessions, reserved}, err
}

func preflight(s doneeval.Snapshot) error {
	if s.Task.ID == "" {
		return &Refusal{Code: "issue-not-found", Message: "task not found", Details: map[string]any{"id": ""}}
	}
	if s.Task.Kind != "task" {
		message := "fid done closes tasks; closing an epic arrives in 1.6k"
		if s.Task.Kind != "epic" {
			message = "fid done closes tasks; use issue transition for subtasks"
		}
		return &Refusal{Code: "invalid-argument", Message: message, Details: map[string]any{"argument": "task", "value": s.Task.ID, "allowed": []string{"task"}}}
	}
	if s.Task.State != "in_review" {
		return refusal(doneeval.TaskGate, string(doneeval.TaskNotInReview), nil)
	}
	if s.Handoff == nil {
		return refusal(doneeval.HandoffGate, string(doneeval.HandoffMissing), nil)
	}
	return nil
}

func observationFailure(err error, merged bool) error {
	var mismatch *github.MismatchError
	if errors.As(err, &mismatch) {
		reason := doneeval.RepositoryMismatch
		if mismatch.Kind == github.BranchMismatch {
			reason = doneeval.BranchMismatch
		}
		return refusal(doneeval.PRGate, string(reason), nil)
	}
	reason := "network-or-api"
	var unavailable *github.UnavailableError
	if errors.As(err, &unavailable) {
		reason = string(unavailable.Reason)
	}
	message := "GitHub confirmation unavailable"
	if merged {
		message = "merged, waiting for confirmation"
	}
	return &Refusal{Code: "github-unavailable", Message: message, Next: "check gh authentication and connectivity, then retry fid done", Details: map[string]any{"reason": reason, "merged": merged}}
}
