package journal

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/linkrelation"
)

// validateIssueRows binds the retained Issue writers' rows and bounds each
// stage. Event-v2 carries arguments in its row keys and values, rather than
// an argv payload: the aggregate identifies the subject and all subsidiary
// rows must agree. Event-v1 has no aggregate; only its immutable Issue facts
// are checked here (validateIssueIdentity).
func validateIssueRows(event Event) error {
	if !event.IsV2() {
		return nil
	}
	family := event.EventType
	if !strings.HasPrefix(family, "issue.") && family != "label.added" && family != "link.added" && family != "link.removed" {
		return nil
	}
	counts := map[string]int{}
	seen := map[string]bool{}
	for _, effect := range event.Effects {
		counts[effect.Table]++
		if err := validateIssueColumns(event, effect); err != nil {
			return err
		}
		if err := validateIssueRow(event, effect, seen); err != nil {
			return err
		}
	}
	maxima := map[string]int{"issues": 1, "epic_metadata": 1, "acceptance_criteria": 1, "labels": 1, "activity_log": 1, "person_decisions": 1, "idempotency_results": 1, "holds": 1, "questions": 1}
	switch family {
	case "issue.updated", "issue.transitioned":
		maxima["issues"] = 2
	case "link.added", "link.removed":
		maxima["issues"] = 2
		maxima["links"] = 2
	case "issue.label.added", "issue.label.removed", "issue.label-added", "issue.label-removed":
		maxima["issue_labels"] = 1
	}
	for table, maximum := range maxima {
		if counts[table] > maximum {
			return rowFailure(event, table, fmt.Sprintf("stage may write at most %d rows", maximum))
		}
	}
	if (family == "issue.updated" || family == "issue.transitioned") && (counts["issues"] != 2 || !seen[event.Aggregate.ID+"/revision"]) {
		return rowFailure(event, "issues", "update requires one scalar row and one revision row")
	}
	if family == "link.added" || family == "link.removed" {
		return validateLinkRows(event, counts)
	}
	return nil
}

func validateIssueRow(event Event, effect Effect, seen map[string]bool) error {
	subject := event.Aggregate.ID
	switch effect.Table {
	case "issues":
		if effect.Op == "insert" {
			if revision, present := effect.Values["revision"]; present && !currentValuesEqual(revision, 1) {
				return rowFailure(event, "issues", "a creation starts at revision one")
			}
		}
		if effect.Op == "update" {
			return validateIssueUpdateRow(event, effect, seen)
		}

	case "acceptance_criteria":
		if effect.Values["issue_id"] != subject {
			return rowFailure(event, effect.Table, "criterion must belong to the Event's Issue")
		}
	case "epic_metadata":
		if effect.PK["issue_id"] != subject || event.Effects[0].Values["kind"] != "epic" {
			return rowFailure(event, effect.Table, "metadata must belong to the Event's epic")
		}
	case "issue_labels":
		if effect.PK["issue_id"] != subject || len(effect.Values) != 0 {
			return rowFailure(event, effect.Table, "label key must belong to the Event's Issue and values must be empty")
		}
		key := fmt.Sprint("label/", effect.PK["label"])
		if seen[key] {
			return rowFailure(event, effect.Table, "a label changes at most once per Event")
		}
		seen[key] = true
	case "labels":
		if effect.PK["name"] != subject {
			return rowFailure(event, effect.Table, "label name must match the Event aggregate")
		}
	case "activity_log":
		refType, refID := "issue", subject
		if event.Aggregate.Type == "link" {
			refType, _ = event.Effects[0].PK["from_type"].(string)
			refID, _ = event.Effects[0].PK["from_id"].(string)
		}
		if effect.Values["journal_event_id"] != event.EventID || effect.Values["ref_type"] != refType || effect.Values["ref_id"] != refID {
			return rowFailure(event, effect.Table, "activity must name its own Event and subject")
		}
	}
	return nil
}

// The first row declares a link's arguments. The second must be its inverse;
// only those two endpoints receive one revision update each.
func validateLinkRows(event Event, counts map[string]int) error {
	if counts["links"] != 2 {
		return rowFailure(event, "links", "a link writes exactly its declared row and inverse")
	}
	declared, inverse := event.Effects[0].PK, event.Effects[1].PK
	from, _ := declared["from_id"].(string)
	to, _ := declared["to_id"].(string)
	relation, _ := declared["relation"].(string)
	inverseRelation := linkrelation.Inverse(relation)
	if from == "" || to == "" || from == to || declared["from_type"] != "issue" || declared["to_type"] != "issue" || inverseRelation == "" || event.Aggregate.ID != from+":"+relation+":"+to {
		return rowFailure(event, "links", "declared key must match the Event's link arguments")
	}
	if inverse["from_type"] != declared["to_type"] || inverse["to_type"] != declared["from_type"] || inverse["from_id"] != to || inverse["to_id"] != from || inverse["relation"] != inverseRelation {
		return rowFailure(event, "links", "second key must be the declared link's inverse")
	}
	ends := map[string]bool{from: false, to: false}
	for _, effect := range event.Effects {
		if effect.Table != "issues" {
			continue
		}
		id, _ := effect.PK["id"].(string)
		changed, endpoint := ends[id]
		if !endpoint || changed {
			return rowFailure(event, "issues", "only a link's two ends may change, once each")
		}
		ends[id] = true
	}
	if counts["issues"] != len(ends) {
		return rowFailure(event, "issues", "each link endpoint needs one revision update")
	}
	return nil
}

func validateIssueUpdateRow(event Event, effect Effect, seen map[string]bool) error {
	id, _ := effect.PK["id"].(string)
	for col := range effect.Values {
		key := id + "/" + col
		if seen[key] {
			return rowFailure(event, "issues", "a column changes at most once per Event")
		}
		seen[key] = true
	}
	if _, revision := effect.Values["revision"]; revision {
		if !currentKeys(effect.Values, "revision") || !currentKeys(effect.Prior, "revision") || !currentKeys(effect.Preconditions, "revision") || !currentValuesEqual(effect.Prior["revision"], effect.Preconditions["revision"]) || !currentNextVersion(effect.Preconditions["revision"], effect.Values["revision"]) {
			return rowFailure(event, "issues", "revision requires its guarded prior plus one, in its own row")
		}
		return nil
	}
	allowed := ""
	switch event.EventType {
	case "issue.updated":
		allowed = "title priority spec_md_path lane updated_at"
	case "issue.transitioned":
		allowed = "status updated_at"
	default:
		return rowFailure(event, "issues", "this stage changes only the revision")
	}
	for col := range effect.Values {
		if !strings.Contains(" "+allowed+" ", " "+col+" ") {
			return rowFailure(event, "issues", col+" is not an argument of this command")
		}
	}
	return nil
}
