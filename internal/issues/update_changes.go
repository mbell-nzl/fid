package issues

import "github.com/mbell-nzl/fid/internal/journal"

func buildIssueUpdateChanges(current issueUpdateState, request normalizedUpdateRequest) (map[string]any, map[string]any, map[string]any, []FieldChange) {
	values := make(map[string]any)
	preconditions := make(map[string]any)
	prior := make(map[string]any)
	changes := make([]FieldChange, 0, 8)

	if request.title != nil && current.Title != *request.title {
		values["title"] = *request.title
		preconditions["title"] = current.Title
		prior["title"] = current.Title
		changes = append(changes, FieldChange{Field: "title", From: current.Title, To: *request.title})
	}
	if request.priority != nil && !sameOptionalInt(current.Priority, request.priority) {
		values["priority"] = *request.priority
		preconditions["priority"] = intPointerValue(current.Priority)
		prior["priority"] = intPointerValue(current.Priority)
		changes = append(changes, FieldChange{Field: "priority", From: intPointerValue(current.Priority), To: *request.priority})
	}
	if request.specMDPath != nil && !sameOptionalString(current.SpecMDPath, request.specMDPath) {
		values["spec_md_path"] = *request.specMDPath
		preconditions["spec_md_path"] = stringPointerValue(current.SpecMDPath)
		prior["spec_md_path"] = stringPointerValue(current.SpecMDPath)
		changes = append(changes, FieldChange{Field: "spec_md_path", From: stringPointerValue(current.SpecMDPath), To: *request.specMDPath})
	}
	if request.ownerRole != nil && !sameOptionalString(current.OwnerRole, request.ownerRole) {
		values["owner_role"] = *request.ownerRole
		preconditions["owner_role"] = stringPointerValue(current.OwnerRole)
		prior["owner_role"] = stringPointerValue(current.OwnerRole)
		changes = append(changes, FieldChange{Field: "owner_role", From: stringPointerValue(current.OwnerRole), To: *request.ownerRole})
	}
	if request.ownerActor != nil && !sameOptionalString(current.OwnerActor, request.ownerActor) {
		values["owner_actor"] = *request.ownerActor
		preconditions["owner_actor"] = stringPointerValue(current.OwnerActor)
		prior["owner_actor"] = stringPointerValue(current.OwnerActor)
		changes = append(changes, FieldChange{Field: "owner_actor", From: stringPointerValue(current.OwnerActor), To: *request.ownerActor})
	}

	return values, preconditions, prior, changes
}

func buildLabelEffects(issueID string, current []string, desired *[]string) ([]journal.Effect, bool, FieldChange) {
	if desired == nil {
		return nil, false, FieldChange{}
	}
	if sameStringSlice(current, *desired) {
		return nil, false, FieldChange{}
	}

	currentSet := make(map[string]bool, len(current))
	for _, label := range current {
		currentSet[label] = true
	}
	desiredSet := make(map[string]bool, len(*desired))
	for _, label := range *desired {
		desiredSet[label] = true
	}

	effects := make([]journal.Effect, 0, len(current)+len(*desired))
	for _, label := range current {
		if desiredSet[label] {
			continue
		}
		effects = append(effects, journal.Effect{
			Op:    "delete",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": issueID,
				"label":    label,
			},
		})
	}
	for _, label := range *desired {
		if currentSet[label] {
			continue
		}
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": issueID,
				"label":    label,
			},
			Values: map[string]any{},
		})
	}

	return effects, true, FieldChange{
		Field: "labels",
		From:  append([]string(nil), current...),
		To:    append([]string(nil), (*desired)...),
	}
}

func activityDetailChanges(changes []FieldChange) []map[string]any {
	details := make([]map[string]any, 0, len(changes))
	for _, change := range changes {
		details = append(details, map[string]any{
			"field": change.Field,
			"from":  change.From,
			"to":    change.To,
		})
	}

	return details
}

func stringPointerValue(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func intPointerValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func sameOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameOptionalInt(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameStringSlice(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}
