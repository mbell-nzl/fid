package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// ValidateReviewResult protects stored retry responses, including direct reads
// from an already-current projection.
func ValidateReviewResult(raw string) error {
	var stored struct {
		Round    reviewrecord.Round `json:"round"`
		At       string             `json:"occurred_at"`
		Replayed bool               `json:"replayed"`
	}
	if err := decodeCheckJSON(raw, &stored); err != nil {
		return err
	}
	if err := stored.Round.Validate(); err != nil {
		return err
	}
	if stored.At != stored.Round.At {
		return fmt.Errorf("review response time differs from round")
	}
	return nil
}

func validateReviewRound(tx *sql.Tx, event Event) error {
	if event.CommandName != "review.record" || event.Aggregate.Type != "review" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("invalid review command, work actor or aggregate")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	for _, e := range event.Effects {
		if e.Table == "issues" || e.Table == "holds" || e.Table == "review_claims" {
			continue // Exact state and reservation effects are checked below.
		}
		if err := validateCurrentEffect(event, e); err != nil {
			return err
		}
	}
	insert := event.Effects[0]
	var r reviewrecord.Round
	if err := decodeCheckJSON(insert.Values["record_json"], &r); err != nil {
		return err
	}
	if err := reviewrecord.Admit(tx, r); err != nil {
		return err
	}
	rework, err := ReviewReworkEffects(tx, r)
	if err != nil {
		return err
	}
	ends, err := ReviewClaimEndEffects(tx, r.TaskID, r.Stream, r.At)
	if err != nil {
		return err
	}
	rework = append(rework, ends...)
	if insert.PK["id"] != r.ID || insert.Values["task_id"] != r.TaskID || insert.Values["handoff_id"] != r.HandoffID || event.Aggregate.ID != r.ID || r.RecordedBy != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || r.At != event.OccurredAt.UTC().Format(time.RFC3339) || len(event.Effects) != len(r.Findings)+closingCount(r)+len(rework)+2 {
		return fmt.Errorf("review round differs from Event")
	}
	for i, f := range r.Findings {
		e := event.Effects[i+1]
		var recorded reviewrecord.Finding
		if err := decodeCheckJSON(e.Values["record_json"], &recorded); err != nil {
			return err
		}
		if e.PK["id"] != f.ID || e.Values["round_id"] != r.ID || !currentValuesEqual(recorded, f) {
			return fmt.Errorf("round finding differs from immutable review")
		}
	}
	index := len(r.Findings) + 1
	for _, d := range r.Dispositions {
		if !d.Closes() {
			continue
		}
		e := event.Effects[index]
		if e.PK["finding_id"] != d.FindingID || e.Values["round_id"] != r.ID || e.Values["decision_id"] != nil || e.Values["disposition"] != d.Disposition || e.Values["reason"] != d.Reason {
			return fmt.Errorf("closure differs from round disposition")
		}
		index++
	}
	for _, want := range rework {
		if !currentValuesEqual(event.Effects[index], want) {
			return fmt.Errorf("review rework differs from task or its reservation")
		}
		index++
	}
	raw, ok := event.Effects[len(event.Effects)-1].Values["result_json"].(string)
	if !ok {
		return fmt.Errorf("missing review retry response")
	}
	if err := ValidateReviewResult(raw); err != nil {
		return err
	}
	var response struct {
		Round reviewrecord.Round `json:"round"`
		At    string             `json:"occurred_at"`
	}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return err
	}
	if !currentValuesEqual(response.Round, r) || response.At != r.At {
		return fmt.Errorf("review response differs from round")
	}
	return nil
}

func closingCount(r reviewrecord.Round) int {
	n := 0
	for _, d := range r.Dispositions {
		if d.Closes() {
			n++
		}
	}
	return n
}
