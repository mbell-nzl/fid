package cli

import (
	"fmt"
	"reflect"
)

type Code string

type Category string

const (
	CategoryInvalidInput Category = "invalid-input"
	CategoryNotFound     Category = "not-found"
	CategoryConflict     Category = "conflict"
	CategoryAuthority    Category = "authority"
	CategoryPrecondition Category = "precondition"
	CategoryIntegrity    Category = "integrity"
	CategoryInternal     Category = "internal"
)

const (
	CodeInvalidArgument        Code = "invalid-argument"
	CodeHoldNotCurrent         Code = "hold-not-current"
	CodeWorktreeBusy           Code = "worktree-busy"
	CodeNotWaiting             Code = "not-waiting"
	CodeChecksNotConfigured    Code = "checks-not-configured"
	CodeComparisonUnavailable  Code = "comparison-unavailable"
	CodeChecksUnsatisfied      Code = "checks-unsatisfied"
	CodeTreeDirty              Code = "tree-dirty"
	CodeTestPathsNotConfigured Code = "test-paths-not-configured"
	CodeCheckUnsuccessful      Code = "check-unsuccessful"
	CodeCheckNotQualified      Code = "check-not-qualified"
	CodeCheckStillRunning      Code = "check-still-running"
	CodeCheckIncomplete        Code = "check-incomplete"
	CodeLaneNotSet             Code = "lane-not-set"
	CodeAlreadyHeld            Code = "already-held"
	CodePolicyUnavailable      Code = "policy-unavailable"
	CodeOperationUnavailable   Code = "operation-unavailable"
	CodeUnsupportedKind        Code = "unsupported-kind"
	CodeUnsupportedRelation    Code = "unsupported-relation"
	CodeProjectNotFound        Code = "project-not-found"
	CodeProjectAlreadyInit     Code = "project-already-initialized"
	CodeIssueNotFound          Code = "issue-not-found"
	CodeLabelNotFound          Code = "label-not-found"
	CodeLinkNotFound           Code = "link-not-found"
	CodeDuplicateEntity        Code = "duplicate-entity"
	CodeStateConflict          Code = "state-conflict"
	CodeLockUnavailable        Code = "lock-unavailable"
	CodeHierarchyInvalid       Code = "hierarchy-invalid"
	CodeDependencyCycle        Code = "dependency-cycle"
	CodeTransitionNotPermitted Code = "transition-not-permitted"
	CodeAuthorityDenied        Code = "authority-denied"
	CodeJournalReplayStopped   Code = "journal-replay-stopped"
	CodeProjectionDivergence   Code = "projection-divergence"
	CodeResultLimitExceeded    Code = "result-limit-exceeded"
	CodeInternalError          Code = "internal-error"
	CodeCapabilityUnavailable  Code = "capability-unavailable"
	CodeArchiveInvalid         Code = "archive-invalid"
	CodeImportTargetConflict   Code = "import-target-conflict"
	CodeRepositoryUnavailable  Code = "repository-unavailable"
	CodeIdempotencyConflict    Code = "idempotency-conflict"
)

type SubjectType string

const (
	SubjectProject    SubjectType = "project"
	SubjectIssue      SubjectType = "issue"
	SubjectLabel      SubjectType = "label"
	SubjectLink       SubjectType = "link"
	SubjectCapability SubjectType = "capability"
	SubjectExport     SubjectType = "export"
)

type Subject struct {
	Type SubjectType `json:"type"`
	ID   any         `json:"id"`
}

type RedactionPolicy string

const RedactClosedDetails RedactionPolicy = "closed-details-only"

type ReasonSpec struct {
	Category        Category
	Retryable       bool
	SubjectTypes    []SubjectType
	DetailFields    []string
	RedactionPolicy RedactionPolicy
}

type Failure struct {
	Code      Code           `json:"code"`
	Category  Category       `json:"category"`
	Retryable bool           `json:"retryable"`
	Message   string         `json:"message"`
	Subject   *Subject       `json:"subject"`
	Details   map[string]any `json:"details"`
	// NextAction carries producer-supplied human guidance outside the JSON contract.
	NextAction string `json:"-"`
}

type detailType uint8

const (
	detailString detailType = iota
	detailInteger
	detailStrings
	detailBoolean
	detailObject
	detailObjects
)

type detailField struct {
	name     string
	typeOf   detailType
	nullable bool
	fields   []detailField
}

type reasonDefinition struct {
	spec    ReasonSpec
	details []detailField
}

func stringField(name string) detailField { return detailField{name: name, typeOf: detailString} }

func nullableStringField(name string) detailField {
	return detailField{name: name, typeOf: detailString, nullable: true}
}

func integerField(name string) detailField { return detailField{name: name, typeOf: detailInteger} }

func stringsField(name string) detailField { return detailField{name: name, typeOf: detailStrings} }

func booleanField(name string) detailField { return detailField{name: name, typeOf: detailBoolean} }

func objectField(name string, fields ...detailField) detailField {
	return detailField{name: name, typeOf: detailObject, fields: fields}
}

func objectsField(name string, fields ...detailField) detailField {
	return detailField{name: name, typeOf: detailObjects, fields: fields}
}

func definition(category Category, retryable bool, subjects []SubjectType, details ...detailField) reasonDefinition {
	fields := make([]string, len(details))
	for index, detail := range details {
		fields[index] = detail.name
	}
	return reasonDefinition{
		spec:    ReasonSpec{category, retryable, subjects, fields, RedactClosedDetails},
		details: details,
	}
}

var failureCodeDefinitions = map[Code]reasonDefinition{
	CodeTreeDirty:              definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, stringField("worktree")),
	CodeTestPathsNotConfigured: definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}),
	CodeChecksNotConfigured:    definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, stringField("lane")),
	CodeComparisonUnavailable:  definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, nullableStringField("ref"), stringField("problem")),
	CodeChecksUnsatisfied:      definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, stringsField("checks"), stringsField("latest_result_ids")),
	CodeCheckUnsuccessful:      definition(CategoryGateUnsuccessful, false, []SubjectType{SubjectIssue}, stringsField("check_ids"), stringsField("outcomes")),
	CodeCheckNotQualified:      definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, stringsField("check_ids"), stringsField("reasons")),
	CodeCheckStillRunning:      definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringField("check_id"), integerField("pid")),
	CodeCheckIncomplete:        definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringsField("check_ids")),
	CodeHoldNotCurrent:         definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringField("hold_id"), stringField("reason")),
	CodeWorktreeBusy:           definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringField("task_id"), stringField("worktree")),
	CodeLaneNotSet:             definition(CategoryPrecondition, false, []SubjectType{SubjectIssue}, stringField("task_id")),
	CodeNotWaiting:             definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringField("task_id")),
	CodeAlreadyHeld:            definition(CategoryConflict, false, []SubjectType{SubjectIssue}, stringField("hold_id")),
	CodePolicyUnavailable:      definition(CategoryPrecondition, false, nil, stringField("problem"), nullableStringField("ref"), nullableStringField("commit")),
	CodeOperationUnavailable: definition(CategoryInvalidInput, false, nil,
		stringField("argument"), nullableStringField("value"), stringsField("allowed")),

	CodeIdempotencyConflict: definition(CategoryConflict, false, nil,
		stringField("operation_scope"), stringField("idempotency_key"), stringField("stored_request_digest"), stringField("supplied_request_digest"), stringsField("differing_fields")),

	CodeRepositoryUnavailable: definition(CategoryPrecondition, true, nil,
		stringField("reason"), stringField("detail")),

	CodeInvalidArgument:     definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeUnsupportedKind:     definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeUnsupportedRelation: definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeProjectNotFound:     definition(CategoryNotFound, false, []SubjectType{SubjectProject}, stringField("path")),
	CodeProjectAlreadyInit:  definition(CategoryConflict, false, []SubjectType{SubjectProject}, stringField("path")),
	CodeIssueNotFound:       definition(CategoryNotFound, false, []SubjectType{SubjectIssue}, stringField("id")),

	CodeLabelNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectLabel}, stringField("id")),
	CodeLinkNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectLink},
		stringField("from_type"), stringField("from_id"), stringField("relation"), stringField("to_type"), stringField("to_id")),

	CodeDuplicateEntity: definition(CategoryConflict, false, []SubjectType{SubjectIssue, SubjectLabel, SubjectLink},
		stringField("entity_type"), objectField("identity", stringsField("fields"), stringsField("values"))),
	CodeStateConflict:   definition(CategoryConflict, false, []SubjectType{SubjectIssue}, integerField("expected_revision"), integerField("actual_revision")),
	CodeLockUnavailable: definition(CategoryConflict, true, nil, stringField("lock_path")),
	CodeHierarchyInvalid: definition(CategoryPrecondition, false, []SubjectType{SubjectIssue},
		stringField("kind"), nullableStringField("parent_id"), nullableStringField("parent_kind"), stringField("rule")),
	CodeDependencyCycle: definition(CategoryPrecondition, false, nil, stringsField("path")),
	CodeTransitionNotPermitted: definition(CategoryPrecondition, false, []SubjectType{SubjectIssue},
		stringField("from_state"), stringField("to_state"), stringsField("unsatisfied")),

	CodeAuthorityDenied: definition(CategoryAuthority, false, nil,
		nullableStringField("actor_id"), stringField("operation"), stringsField("required_actions"), stringsField("required_capabilities")),
	CodeJournalReplayStopped: definition(CategoryIntegrity, false, nil,
		stringField("journal_file"), nullableIntegerField("line"), nullableStringField("event_id")),
	CodeProjectionDivergence: definition(CategoryIntegrity, false, nil,
		stringField("status"), objectsField("errors", findingFields()...), objectsField("warnings", findingFields()...), objectField("info", infoFields()...)),
	CodeResultLimitExceeded: definition(CategoryIntegrity, false, nil, stringField("collection"), integerField("limit"), integerField("observed")),
	CodeInternalError:       definition(CategoryInternal, false, nil),

	CodeCapabilityUnavailable: definition(CategoryAuthority, false, []SubjectType{SubjectCapability},
		stringField("capability_id"), booleanField("available"), integerField("owner_phase")),

	CodeArchiveInvalid: definition(CategoryIntegrity, false, []SubjectType{SubjectExport},
		stringField("archive_path"), objectsField("checks", stringField("code"), nullableStringField("entry"), stringField("detail"))),
	CodeImportTargetConflict: definition(CategoryConflict, false, []SubjectType{SubjectExport},
		stringField("target_root"), stringField("reason")),
}

var registeredCodes = []Code{CodeTreeDirty, CodeTestPathsNotConfigured, CodeComparisonUnavailable, CodeChecksUnsatisfied, CodeChecksNotConfigured, CodeCheckUnsuccessful, CodeCheckNotQualified, CodeCheckStillRunning, CodeCheckIncomplete, CodeLaneNotSet, CodeNotWaiting, CodeHoldNotCurrent, CodeWorktreeBusy, CodeAlreadyHeld,
	CodeInvalidArgument, CodePolicyUnavailable, CodeOperationUnavailable, CodeUnsupportedKind, CodeUnsupportedRelation, CodeProjectNotFound, CodeProjectAlreadyInit,
	CodeIssueNotFound, CodeLabelNotFound, CodeLinkNotFound, CodeDuplicateEntity, CodeStateConflict, CodeLockUnavailable, CodeHierarchyInvalid, CodeDependencyCycle,
	CodeTransitionNotPermitted, CodeAuthorityDenied,
	CodeJournalReplayStopped, CodeProjectionDivergence, CodeResultLimitExceeded, CodeInternalError,
	CodeCapabilityUnavailable, CodeArchiveInvalid,
	CodeImportTargetConflict,
	CodeRepositoryUnavailable}

func findingFields() []detailField {
	return []detailField{stringField("check"), stringField("message"), stringField("entity_id")}
}

func infoFields() []detailField {
	return []detailField{integerField("issue_count"), integerField("claim_count_live"), integerField("recent_activity_7d")}
}

func RegisteredCodes() []Code { return append([]Code(nil), registeredCodes...) }

func LookupCode(code Code) (ReasonSpec, bool) {
	definition, ok := failureCodeDefinitions[code]
	if !ok {
		return ReasonSpec{}, false
	}
	spec := definition.spec
	spec.SubjectTypes = append([]SubjectType(nil), spec.SubjectTypes...)
	spec.DetailFields = append([]string{}, spec.DetailFields...)
	return spec, true
}

func NewFailure(code Code, message string, subject *Subject, details map[string]any) (Failure, error) {
	definition, ok := failureCodeDefinitions[code]
	if !ok {
		return Failure{}, fmt.Errorf("new failure: unknown code %q", code)
	}
	failure := Failure{Code: code, Category: definition.spec.Category, Retryable: definition.spec.Retryable,
		Message: message, Subject: subject, Details: details}
	if err := validateFailure(failure, definition); err != nil {
		return Failure{}, fmt.Errorf("new failure: %w", err)
	}
	return failure, nil
}

func validateFailure(failure Failure, definition reasonDefinition) error {
	if failure.Category != definition.spec.Category || failure.Retryable != definition.spec.Retryable {
		return fmt.Errorf("code %q metadata does not match its registry entry", failure.Code)
	}
	if err := validateSubject(failure.Subject, definition.spec.SubjectTypes); err != nil {
		return err
	}
	fields := definition.details
	// Removed history adds a closed schema-refusal variant. Existing corruption
	// and precondition failures retain their original three-field details.
	if failure.Code == CodeJournalReplayStopped && failure.Details["reason"] != nil {
		if failure.Details["reason"] != "schema-unsupported" || failure.Details["event_type"] == "" {
			return fmt.Errorf("invalid journal schema refusal")
		}
		fields = append(append([]detailField{}, fields...), stringField("reason"), stringField("event_type"))
	}
	if err := validateDetails(failure.Details, fields); err != nil {
		return err
	}
	if failure.Code == CodeHoldNotCurrent {
		switch failure.Details["reason"] {
		case "not-found", "ended", "expired", "actor-mismatch", "worktree-mismatch":
		default:
			return fmt.Errorf("unknown hold-not-current reason")
		}
	}
	if failure.Code == CodeDuplicateEntity {
		return validateDuplicateEntityIdentity(failure.Details)
	}
	if failure.Code == CodeImportTargetConflict {
		return validateImportTargetConflictReason(failure.Details)
	}
	return nil
}

func validateImportTargetConflictReason(details map[string]any) error {
	reason, _ := details["reason"].(string)
	for _, allowed := range []string{"populated", "different-import", "unsafe-path"} {
		if reason == allowed {
			return nil
		}
	}
	return fmt.Errorf("import-target-conflict reason %q is not accepted", reason)
}

func validateDuplicateEntityIdentity(details map[string]any) error {
	identity, ok := details["identity"].(map[string]any)
	if !ok {
		return fmt.Errorf("duplicate-entity identity must be an object")
	}
	fields, ok := stringArrayLength(identity["fields"])
	if !ok {
		return fmt.Errorf("duplicate-entity identity fields must be a string array")
	}
	values, ok := stringArrayLength(identity["values"])
	if !ok {
		return fmt.Errorf("duplicate-entity identity values must be a string array")
	}
	if fields != values {
		return fmt.Errorf("duplicate-entity identity fields and values must have equal length")
	}
	return nil
}

func stringArrayLength(value any) (int, bool) {
	switch values := value.(type) {
	case []string:
		return len(values), true
	case []any:
		return len(values), true
	default:
		return 0, false
	}
}

func validateSubject(subject *Subject, allowed []SubjectType) error {
	if len(allowed) == 0 {
		if subject != nil {
			return fmt.Errorf("subject must be null")
		}
		return nil
	}
	if subject == nil {
		return fmt.Errorf("subject is required")
	}
	for _, subjectType := range allowed {
		if subject.Type == subjectType {
			return nil
		}
	}
	return fmt.Errorf("subject type %q is not allowed", subject.Type)
}

func validateDetails(values map[string]any, fields []detailField) error {
	if values == nil {
		return fmt.Errorf("details must be an object")
	}
	if len(values) != len(fields) {
		return fmt.Errorf("details must contain exactly %d fields", len(fields))
	}
	for _, field := range fields {
		value, ok := values[field.name]
		if !ok {
			return fmt.Errorf("details missing %q", field.name)
		}
		if err := validateDetailValue(value, field); err != nil {
			return fmt.Errorf("details %q: %w", field.name, err)
		}
	}
	return nil
}

func validateDetailValue(value any, field detailField) error {
	if value == nil {
		if field.nullable {
			return nil
		}
		return fmt.Errorf("must not be null")
	}
	switch field.typeOf {
	case detailString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("must be a string")
		}
	case detailInteger:
		if !isInteger(value) {
			return fmt.Errorf("must be an integer")
		}
	case detailStrings:
		if !isStrings(value) {
			return fmt.Errorf("must be a string array")
		}
	case detailBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("must be a boolean")
		}
	case detailObject:
		object, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("must be an object")
		}
		return validateDetails(object, field.fields)
	case detailObjects:
		return validateObjects(value, field.fields)
	}
	return nil
}

func isInteger(value any) bool {
	kind := reflect.ValueOf(value).Kind()
	return kind >= reflect.Int && kind <= reflect.Int64 || kind >= reflect.Uint && kind <= reflect.Uint64
}

func isStrings(value any) bool {
	values, ok := value.([]string)
	if ok {
		return values != nil
	}
	items, ok := value.([]any)
	if !ok || items == nil {
		return false
	}
	for _, item := range items {
		if _, ok := item.(string); !ok {
			return false
		}
	}
	return true
}

func validateObjects(value any, fields []detailField) error {
	items, ok := value.([]map[string]any)
	if ok {
		for _, item := range items {
			if err := validateDetails(item, fields); err != nil {
				return err
			}
		}
		return nil
	}
	values, ok := value.([]any)
	if !ok || values == nil {
		return fmt.Errorf("must be an object array")
	}
	for _, value := range values {
		item, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("must be an object array")
		}
		if err := validateDetails(item, fields); err != nil {
			return err
		}
	}
	return nil
}

func nullableIntegerField(name string) detailField {
	return detailField{name: name, typeOf: detailInteger, nullable: true}
}
