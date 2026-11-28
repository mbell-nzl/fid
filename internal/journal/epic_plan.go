package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
)

func validateEpicPlan(tx *sql.Tx, event Event) error {
	if err := validateEpicRetry(event); err != nil {
		return err
	}
	if event.CommandName != event.EventType || event.Repository != nil {
		return fmt.Errorf("epic plan Event requires its command")
	}
	p, err := epicplan.Preview(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	if p.State == "done" || p.State == "cancelled" {
		return fmt.Errorf("epic plan requires an open epic")
	}
	count := len(event.Effects)
	if event.Idempotency != nil {
		count--
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	for _, e := range event.Effects {
		if err := validateCurrentEffect(event, e); err != nil {
			return err
		}
	}
	if event.EventType == "epic.validation.pass" || event.EventType == "epic.validation.verify" {
		return validateEpicResults(tx, event, count)
	}
	if strings.HasPrefix(event.EventType, "epic.validation.") {
		return validateValidationChange(tx, event, count)
	}
	if event.EventType != "epic.approve" {
		return fmt.Errorf("unknown epic plan event")
	}
	e := event.Effects[0]
	d := e.Values
	var binding epicplan.Binding
	if err := decodeCheckJSON(d["binding_json"], &binding); err != nil {
		return err
	}
	if !currentKeys(d, "kind actor subject_type subject_id reason binding_id decided_at binding_json") || d["kind"] != "approve-epic" || d["actor"] != "human:"+event.ActorV2.ActorID || d["subject_type"] != "epic" || d["subject_id"] != p.EpicID || d["binding_id"] != binding.Digest || d["decided_at"] != event.OccurredAt.UTC().Format(time.RFC3339) || binding != p.Binding {
		return fmt.Errorf("epic approval must sign the current validation list")
	}
	if reason, ok := d["reason"].(string); !ok || strings.TrimSpace(reason) == "" {
		return fmt.Errorf("epic approval requires a reason")
	}
	index := 1
	for _, t := range p.Tasks {
		if !t.Accept {
			continue
		}
		if index >= count {
			return fmt.Errorf("epic approval must accept all unaccepted agent work")
		}
		a := event.Effects[index]
		if a.Table != "epic_plan_acceptances" || a.PK["issue_id"] != t.ID || !currentKeys(a.Values, "decision_id") || a.Values["decision_id"] != e.PK["id"] {
			return fmt.Errorf("epic acceptance differs from plan")
		}
		index++
	}
	if index != count {
		return fmt.Errorf("epic approval has unrelated acceptances")
	}
	return nil
}

func validateValidationChange(tx *sql.Tx, event Event, count int) error {
	if count != 1 {
		return fmt.Errorf("validation change requires exactly one revision")
	}
	e := event.Effects[0]
	if !currentKeys(e.Values, "ordinal claim verifier recipe removed") || e.PK["epic_id"] != event.Aggregate.ID {
		return fmt.Errorf("validation change differs from epic")
	}
	list, generation, ordinal, err := epicplan.Read(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	if !currentValuesEqual(e.PK["generation"], generation+1) {
		return fmt.Errorf("validation generation must increase")
	}
	var change struct {
		epicplan.Validation
		Removed bool `json:"removed"`
	}
	if err := decodeCheckJSON(epicChangeJSON(e.Values), &change); err != nil {
		return err
	}
	v := change.Validation
	// The SQL bool is checked independently of its JSON representation.
	removed := currentValuesEqual(e.Values["removed"], true) || currentValuesEqual(e.Values["removed"], 1)
	if event.EventType == "epic.validation.add" {
		if v.Ordinal != ordinal+1 || removed {
			return fmt.Errorf("validation IDs cannot be reused")
		}
	} else {
		found := false
		for _, old := range list {
			if old.Ordinal == v.Ordinal {
				found = true
				if event.EventType == "epic.validation.remove" && (v.Claim != old.Claim || v.Verifier != old.Verifier || v.Recipe != old.Recipe) {
					return fmt.Errorf("removal must retain validation identity")
				}
			}
		}
		if !found || removed != (event.EventType == "epic.validation.remove") {
			return fmt.Errorf("validation change requires a live validation")
		}
	}
	normalized, err := epicplan.Normalize(v)
	if err != nil {
		return err
	}
	if normalized != v {
		return fmt.Errorf("validation must be normalized")
	}
	return nil
}

func epicChangeJSON(v any) string { raw, _ := json.Marshal(v); return string(raw) }
