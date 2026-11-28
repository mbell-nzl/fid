package journal

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/modelid"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// ReviewClaimEndEffects ends the task/stream reservation in the round's own
// Event, even when the recorder couriers a review produced without a claim.
func ReviewClaimEndEffects(conn reviewrecord.Queryer, task, stream, at string) ([]Effect, error) {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='review_claims'`).Scan(&installed); err != nil || installed == 0 {
		return nil, err
	}
	claims, err := reviewrecord.Claims(conn)
	if err != nil {
		return nil, err
	}
	effects := []Effect{}
	for _, c := range claims {
		if c.TaskID == task && c.Stream == stream && c.EndedAt == nil {
			effects = append(effects, EndReviewClaim(c.ID, at))
		}
	}
	return effects, nil
}

func EndReviewClaim(id, at string) Effect {
	return Effect{Op: "update", Table: "review_claims", PK: map[string]any{"id": id}, Values: map[string]any{"ended_at": at}, Preconditions: map[string]any{"ended_at": nil}}
}

func validateReviewClaim(tx *sql.Tx, event Event) error {
	stamp := event.OccurredAt.UTC().Format(time.RFC3339)
	claims, err := reviewrecord.Claims(tx)
	if err != nil {
		return err
	}
	if event.EventType == "review.heartbeat" {
		if event.CommandName != "hold.heartbeat" || len(event.Effects) != 1 {
			return fmt.Errorf("review heartbeat requires one renewal")
		}
		e := event.Effects[0]
		for _, c := range claims {
			if c.ID == event.Aggregate.ID && e.PK["id"] == c.ID && c.EndedAt == nil && c.Actor == event.ActorV2.ActorType+":"+event.ActorV2.ActorID && stamp >= c.LastHeartbeatAt && stamp < c.ExpiresAt {
				if !currentValuesEqual(e.Values, map[string]any{"last_heartbeat_at": stamp, "expires_at": event.OccurredAt.Add(time.Hour).UTC().Format(time.RFC3339)}) || !currentValuesEqual(e.Preconditions, map[string]any{"ended_at": nil, "expires_at": c.ExpiresAt}) {
					return fmt.Errorf("invalid review heartbeat renewal")
				}
				return reviewClaimTask(tx, c.TaskID, c.HandoffID)
			}
		}
		return fmt.Errorf("review session is no longer current")
	}
	if event.CommandName != "next.claim" || event.Aggregate.Sequence != 1 || event.Idempotency != nil {
		return fmt.Errorf("invalid review reservation command")
	}
	e := event.Effects[len(event.Effects)-1]
	if err := validateCurrentEffect(event, e); err != nil {
		return err
	}
	v := e.Values
	task, _ := v["task_id"].(string)
	stream, _ := v["stream"].(string)
	handoff, _ := v["handoff_id"].(string)
	root, _ := v["worktree_root"].(string)
	actor := event.ActorV2.ActorType + ":" + event.ActorV2.ActorID
	model, _ := v["model"].(string)
	if e.PK["id"] != event.Aggregate.ID || v["actor"] != actor || root == "" || v["worktree_id"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root))) || model != "" && !modelid.Valid(model) || event.ActorV2.ActorType == "agent" && model == "" || v["since"] != stamp || v["last_heartbeat_at"] != stamp || v["expires_at"] != event.OccurredAt.Add(time.Hour).UTC().Format(time.RFC3339) || v["ended_at"] != nil {
		return fmt.Errorf("review reservation differs from Event identity or times")
	}
	if err := reviewClaimTask(tx, task, handoff); err != nil {
		return err
	}
	built, err := reviewrecord.BuiltBy(tx, task, actor)
	if err != nil {
		return fmt.Errorf("read review builder history: %w", err)
	}
	if built {
		return fmt.Errorf("review reservation requires an independent actor")
	}
	// Admission uses only recorded bindings; it never requires the original
	// repository's objects on replay or import. Live selection computes these.
	var fp fingerprint.Result
	if err := decodeCheckJSON(v["binding_json"], &fp); err != nil {
		return err
	}
	if fp.Algorithm != fingerprint.Algorithm || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(fp.Fingerprint) {
		return fmt.Errorf("review reservation requires a fingerprint")
	}
	s, err := doneeval.GatherReviews("", tx, task, &fp)
	if err != nil {
		return err
	}
	if fp.Base != s.Handoff.Base || fp.Head != s.Handoff.Head {
		return fmt.Errorf("review reservation fingerprint differs from hand-off")
	}
	due := doneeval.DueReview(s)
	if stream == "judge" {
		eligible, err := reviewrecord.JudgeDue(tx, *s.Handoff, reviewrecord.Bind(fp))
		if err != nil {
			return err
		}
		independent, err := reviewrecord.JudgeIndependent(tx, task, actor, "")
		if err != nil {
			return err
		}
		if eligible && independent {
			due = &doneeval.ReviewWork{Stream: "judge"}
		}
	}
	if due == nil || due.Stream != stream {
		return fmt.Errorf("review reservation stream is not due")
	}
	var busy int
	if err := tx.QueryRow(`SELECT count(*) FROM holds h JOIN issues i ON i.id=h.task_id WHERE h.worktree_id=? AND (h.kept=1 OR h.ended_at IS NULL AND (h.expires_at>? OR i.status IN ('in_progress','in_review')))`, v["worktree_id"], stamp).Scan(&busy); err != nil {
		return fmt.Errorf("read review worktree reservations: %w", err)
	}
	if busy != 0 {
		return fmt.Errorf("review worktree already has build work")
	}
	want := []Effect{}
	for _, c := range claims {
		if c.EndedAt != nil || c.TaskID != task && c.WorktreeID != v["worktree_id"] {
			continue
		}
		if c.ExpiresAt > stamp && reviewClaimTask(tx, c.TaskID, c.HandoffID) == nil {
			return fmt.Errorf("review is held elsewhere")
		}
		want = append(want, EndReviewClaim(c.ID, stamp))
	}
	if !currentValuesEqual(event.Effects[:len(event.Effects)-1], want) {
		return fmt.Errorf("review reservation must release exactly its lapsed predecessors")
	}
	return nil
}

func reviewClaimTask(conn reviewrecord.Queryer, task, handoff string) error {
	var state string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id=? AND kind='task'`, task).Scan(&state); err != nil {
		return err
	}
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		return err
	}
	if state != "in_review" || h.ID != handoff {
		return fmt.Errorf("review reservation needs the latest in-review hand-off")
	}
	return nil
}

func validateBuildReviewReservation(conn reviewrecord.Queryer, effect Effect, at time.Time) error {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='review_claims'`).Scan(&installed); err != nil || installed == 0 {
		return err
	}
	claims, err := reviewrecord.Claims(conn)
	if err != nil {
		return err
	}
	for _, c := range claims {
		if c.EndedAt == nil && c.WorktreeID == effect.Values["worktree_id"] && c.ExpiresAt > at.UTC().Format(time.RFC3339) && reviewClaimTask(conn, c.TaskID, c.HandoffID) == nil {
			return fmt.Errorf("this worktree already has a live review reservation")
		}
	}
	return nil
}
