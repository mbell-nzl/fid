package journal

import (
	"fmt"
	"strings"
)

// RowIdentityError refuses a malformed Event's row identity, cardinality or
// immutable facts before any of its Effects are applied.
type RowIdentityError struct{ EventID, Table, Reason string }

func (err *RowIdentityError) Error() string {
	return fmt.Sprintf("Event %q: malformed %s Effect: %s", err.EventID, err.Table, err.Reason)
}

func (*RowIdentityError) RowAdmissionRefusal() {}

// Issue columns have the same domains throughout the retained profiles.
// Optional columns may be absent on creation: SQLite supplies their defaults.
var issueRowTables = map[string]currentTable{
	"issues":              currentTableColumns("id", "id:s kind:s parent_id:s? title:s workflow_id:s status:s priority:i? sequence:i? spec_md_path:s? owner_role:s? owner_actor:s? created_at:s? updated_at:s? revision:i lane:s?"),
	"acceptance_criteria": currentTableColumns("id", "id:s issue_id:s ordinal:i text:s"),
	"epic_metadata":       currentTableColumns("issue_id", "issue_id:s scheduling:s brief_md_path:s? target_start:s? target_end:s?"),
	"issue_labels":        currentTableColumns("issue_id label", "issue_id:s label:s"),
	"labels":              currentTableColumns("name", "name:s colour:s? description:s?"),
	"links":               currentTableColumns("from_type from_id relation to_type to_id", "from_type:s from_id:s relation:s to_type:s to_id:s created_at:s?"),
	"activity_log":        currentEventTables["activity_log"],
}

func rowFailure(event Event, table, reason string) error {
	return &RowIdentityError{EventID: event.EventID, Table: table, Reason: reason}
}

func validateIssueColumns(event Event, effect Effect) error {
	table, ok := issueRowTables[effect.Table]
	if !ok {
		return nil
	}
	if !currentKeys(effect.PK, strings.Join(table.primary, " ")) {
		return rowFailure(event, effect.Table, "must key exactly one row by its primary key")
	}
	for _, values := range []map[string]any{effect.PK, effect.Values, effect.Prior, effect.Preconditions} {
		for col, value := range values {
			kind, found := table.columns[col]
			if effect.Table == "activity_log" && col == "details" {
				continue
			} // The projector encodes structured details.
			if !found || !currentScalar(kind, currentEventValue(value)) {
				return rowFailure(event, effect.Table, "invalid column or scalar "+col)
			}
		}
	}
	if effect.Op == "insert" || effect.Op == "activity" {
		if _, err := mergeEffectValues(effect.PK, effect.Values); err != nil {
			return rowFailure(event, effect.Table, err.Error())
		}
	}
	return nil
}

// No command renames, deletes or reparents an Issue, or changes its kind,
// workflow or creation facts. This also covers Event-v1 history, whose
// families and keys otherwise retain their historical grammar.
func validateIssueIdentity(event Event) error {
	for _, effect := range event.Effects {
		if effect.Table != "issues" {
			continue
		}
		if effect.Op == "delete" {
			return rowFailure(event, "issues", "an Issue cannot be deleted or recreated")
		}
		if effect.Op != "update" {
			continue
		}
		for _, col := range strings.Fields("id kind parent_id workflow_id created_at sequence owner_role owner_actor") {
			if _, sets := effect.Values[col]; sets {
				return rowFailure(event, "issues", col+" is immutable")
			}
		}
	}
	return nil
}
