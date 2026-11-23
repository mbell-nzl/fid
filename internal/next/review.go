package next

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type reviewCandidate struct {
	conflict    *defaultconflict.Record
	unavailable string
	work        *doneeval.ReviewWork
	claim       *reviewrecord.Claim
	built       bool
	capped      bool
}

func readReviews(conn reviewrecord.Queryer, options Options, work issues.Work, now time.Time) (map[string]reviewCandidate, error) {
	claims, err := reviewrecord.Claims(conn)
	if err != nil {
		return nil, err
	}
	reviews := map[string]reviewCandidate{}
	for _, task := range work.Tasks {
		if task.Status == "in_progress" {
			limit, err := doneeval.ReadReworkLimit(conn, task.ID)
			if err != nil {
				return nil, err
			}
			if limit != nil && limit.Reached() {
				reviews[task.ID] = reviewCandidate{capped: true}
			}
		}
		if task.Status != "in_review" {
			continue
		}

		s, err := doneeval.GatherReviews(options.WorkspaceRoot, conn, task.ID, nil)
		candidate := reviewCandidate{}
		for _, c := range claims {
			if s.Handoff != nil && c.TaskID == task.ID && c.HandoffID == s.Handoff.ID && c.EndedAt == nil && c.ExpiresAt > holds.Timestamp(now) {
				candidate.claim = &c
			}
		}
		if err != nil {
			// Keep an existing reviewer's worktree occupied even when its
			// comparison objects are unavailable. Never start another task.
			if _, ok := err.(*fingerprint.Failure); ok {
				candidate.unavailable = err.Error()
				reviews[task.ID] = candidate
				continue
			}
			return nil, err
		}
		candidate.work = doneeval.DueReview(s)
		awaiting, err := doneeval.AwaitingMerge(conn, s)
		if err != nil {
			return nil, err
		}
		if awaiting {
			evidence, err := defaultconflict.Inspect(options.WorkspaceRoot, s.Handoff.ID, s.Handoff.Head)
			if err != nil {
				candidate.unavailable = err.Error()
			}
			if evidence != nil {
				candidate.conflict = &defaultconflict.Record{Evidence: *evidence, Binding: s.Fingerprint}
			}
		}
		if options.Actor != "" {
			candidate.built, err = reviewrecord.BuiltBy(conn, task.ID, options.Actor)
			if err != nil {
				return nil, err
			}
		}

		reviews[task.ID] = candidate
	}
	return reviews, nil
}

func claimHold(c reviewrecord.Claim) *holds.Hold {
	since, _ := time.Parse(time.RFC3339, c.Since)
	expires, _ := time.Parse(time.RFC3339, c.ExpiresAt)
	heartbeat, _ := time.Parse(time.RFC3339, c.LastHeartbeatAt)
	return &holds.Hold{ID: c.ID, TaskID: c.TaskID, WorktreeID: c.WorktreeID, Worktree: c.Worktree, Actor: c.Actor, Model: c.Model, Since: since, ExpiresAt: expires, LastHeartbeatAt: heartbeat}
}

func takeReview(lock *journal.WriteLock, conn *sql.DB, options Options, result *Result, candidate reviewCandidate) error {
	if candidate.claim != nil {
		result.Hold = claimHold(*candidate.claim)
		return nil
	}
	id, err := journal.GenerateID("RCL", func(id string) (bool, error) {
		var n int
		err := conn.QueryRow(`SELECT count(*) FROM review_claims WHERE id=?`, id).Scan(&n)
		return n != 0, err
	})
	if err != nil {
		return err
	}
	fp, err := fingerprint.Compute(options.WorkspaceRoot, result.Review.Handoff.Base, result.Review.Handoff.Head)
	if err != nil {
		return err
	}
	binding, err := json.Marshal(fp)
	if err != nil {
		return err
	}
	h := holds.New(id, result.Task.ID, project.WorktreeID(options.WorkspaceRoot), options.WorkspaceRoot, options.Actor, result.Now, nil)
	h.Model = options.Model
	claims, err := reviewrecord.Claims(conn)
	if err != nil {
		return err
	}
	effects := []journal.Effect{}
	for _, c := range claims {
		if c.EndedAt == nil && (c.TaskID == h.TaskID || c.WorktreeID == h.WorktreeID) {
			effects = append(effects, journal.EndReviewClaim(c.ID, holds.Timestamp(result.Now)))
		}
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "review_claims", PK: map[string]any{"id": id}, Values: map[string]any{
		"task_id": h.TaskID, "handoff_id": result.Review.Handoff.ID, "stream": result.Review.Stream,
		"worktree_id": h.WorktreeID, "worktree_root": h.Worktree, "actor": h.Actor, "model": h.Model,
		"since": holds.Timestamp(h.Since), "expires_at": holds.Timestamp(h.ExpiresAt), "last_heartbeat_at": holds.Timestamp(h.Since), "ended_at": nil, "binding_json": string(binding),
	}})
	actorType, actorID, _ := strings.Cut(options.Actor, ":")
	actor := journal.EventActor{ActorID: actorID, ActorType: actorType, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
	appended, err := journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "review.claimed", CommandName: "next.claim", OccurredAt: h.Since, ActorV2: actor, Aggregate: journal.Aggregate{ID: id, Type: "review-claim", Sequence: 1}, Effects: effects})
	if err != nil {
		return fmt.Errorf("reserve review: %w", err)
	}
	result.Hold, result.Event = &h, &Event{ID: appended.EventID, OccurredAt: appended.OccurredAt}
	result.Queue[result.Served].Working = working(h, h.WorktreeID, result.Now)
	return nil
}
