package journal

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicrecord"
)

func validateEpicResults(tx *sql.Tx, event Event, count int) error {
	list, _, _, err := epicplan.Read(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	merges, err := epicplan.Merges(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	person := event.EventType == "epic.validation.pass"
	if !person {
		eligible, err := epicrecord.ReadEligibility(tx, event.Aggregate.ID)
		if err != nil {
			return err
		}
		if !eligible.CanClose() {
			return fmt.Errorf("validation verification requires no open tasks")
		}
	}
	seen := map[string]bool{}
	usedArtifacts := map[string]bool{}
	for _, e := range event.Effects[:count] {
		if e.Table == "person_decisions" || e.Table == "referenced_artifacts" {
			continue
		}
		if e.Table != "epic_validation_results" || !currentKeys(e.Values, "epic_id ordinal generation record_json") {
			return fmt.Errorf("invalid validation result row")
		}
		var p epicplan.Pass
		if err := decodeCheckJSON(e.Values["record_json"], &p); err != nil {
			return err
		}
		var v epicplan.Validation
		for _, candidate := range list {
			if candidate.ID == p.ValidationID {
				v = candidate
			}
		}
		if v.ID == "" || seen[v.ID] || p.ID != e.PK["id"] || p.EpicID != event.Aggregate.ID || e.Values["epic_id"] != p.EpicID || !currentValuesEqual(e.Values["ordinal"], v.Ordinal) || !currentValuesEqual(e.Values["generation"], v.Generation) || p.Generation != v.Generation || p.Actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || p.At != event.OccurredAt.UTC().Format(time.RFC3339) || !slices.Equal(p.Merges, merges) || (len(p.Commit) != 40 || strings.Trim(p.Commit, "0123456789abcdef") != "") {
			return fmt.Errorf("validation result differs from current validation, actor or merges")
		}
		seen[v.ID] = true
		if person {
			if count != 2 || strings.HasPrefix(v.Verifier, "check:") || p.Outcome != "passed" || p.Check != nil || p.Definition != nil {
				return fmt.Errorf("person passes apply only to review or human validations")
			}
			d := event.Effects[0]
			if d.Table != "person_decisions" || !currentKeys(d.Values, "kind actor subject_type subject_id reason binding_id decided_at binding_json") || d.PK["id"] != p.DecisionID || d.Values["kind"] != "pass-validation" || d.Values["actor"] != p.Actor || d.Values["subject_type"] != "epic" || d.Values["subject_id"] != p.EpicID || d.Values["binding_id"] != v.ID || d.Values["decided_at"] != p.At {
				return fmt.Errorf("validation pass requires its person decision")
			}
			var binding epicplan.Pass
			if err := decodeCheckJSON(d.Values["binding_json"], &binding); err != nil {
				return err
			}
			binding.ID = p.ID
			if epicChangeJSON(binding) != epicChangeJSON(p) {
				return fmt.Errorf("validation pass decision differs from result")
			}
			if reason, ok := d.Values["reason"].(string); !ok || strings.TrimSpace(reason) == "" {
				return fmt.Errorf("validation pass requires reason")
			}
		} else {
			if !strings.HasPrefix(v.Verifier, "check:") || p.Check == nil || p.Definition == nil || p.DecisionID != "" || (p.Outcome != "passed" && p.Outcome != "failed") {
				return fmt.Errorf("check validation requires execution evidence")
			}
			if len(p.Definition.Run) == 0 || p.Definition.Timeout <= 0 || p.Check.Counts && p.Check.Tree.Head != p.Commit || p.Check.Counts != (p.Outcome == "passed") || p.Check.Counts && (p.Check.Outcome != "passed" || p.Check.ExitCode == nil || *p.Check.ExitCode != 0 || p.Check.Signal != nil || p.Check.Tree.Dirty || len(p.Check.Reasons) > 0 || len(p.Check.RefusedArtifacts) > 0) {
				return fmt.Errorf("validation check result does not qualify")
			}
			if err := validateEpicFiles(tx, event, *p.Check, *p.Definition, usedArtifacts); err != nil {
				return err
			}
		}
	}
	for _, e := range event.Effects {
		if e.Table == "referenced_artifacts" {
			id, _ := e.PK["id"].(string)
			if !usedArtifacts[id] {
				return fmt.Errorf("verification cannot change unrelated artifacts")
			}
		}
	}
	if person && len(seen) != 1 {
		return fmt.Errorf("person pass requires exactly one validation result")
	}
	if !person {
		for _, v := range list {
			if strings.HasPrefix(v.Verifier, "check:") && !seen[v.ID] {
				return fmt.Errorf("verify must record each check validation")
			}
		}
	}
	return nil
}
