package journal

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
)

func validateEpicRetry(event Event) error {
	if event.Idempotency == nil {
		return nil
	}
	var r epicplan.Result
	if err := decodeCheckJSON(event.Effects[len(event.Effects)-1].Values["result_json"], &r); err != nil {
		return err
	}
	if r.Replayed || r.EpicID != event.Aggregate.ID || r.At != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("epic retry result differs from event")
	}
	accepted := []string{}
	checks := []epicplan.Pass{}
	decision := ""
	for _, e := range event.Effects {
		switch e.Table {
		case "person_decisions":
			decision, _ = e.PK["id"].(string)
		case "epic_plan_acceptances":
			id, _ := e.PK["issue_id"].(string)
			accepted = append(accepted, id)
		case "epic_validation_results":
			if event.EventType == "epic.validation.verify" {
				var p epicplan.Pass
				if err := decodeCheckJSON(e.Values["record_json"], &p); err != nil {
					return err
				}
				checks = append(checks, p)
			}
		case "epic_validation_changes":
			if r.Validation == nil {
				return fmt.Errorf("validation retry is missing its validation")
			}
			v := r.Validation
			if !currentValuesEqual(e.PK["generation"], v.Generation) || !currentValuesEqual(e.Values["ordinal"], v.Ordinal) || v.ID != fmt.Sprintf("V%d", v.Ordinal) || e.Values["claim"] != v.Claim || e.Values["verifier"] != v.Verifier || e.Values["recipe"] != v.Recipe {
				return fmt.Errorf("validation retry differs from revision")
			}
		}
	}
	if decision != r.DecisionID || strings.Join(accepted, "\n") != strings.Join(r.Accepted, "\n") {
		return fmt.Errorf("approval retry differs from accepted plan")
	}
	if event.EventType == "epic.validation.verify" && !reflect.DeepEqual(checks, r.Checks) {
		return fmt.Errorf("verification retry differs from results")
	}
	return nil
}
