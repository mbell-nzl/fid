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
	CodeUnsupportedKind        Code = "unsupported-kind"
	CodeUnsupportedRelation    Code = "unsupported-relation"
	CodeProjectNotFound        Code = "project-not-found"
	CodeProjectAlreadyInit     Code = "project-already-initialized"
	CodeIssueNotFound          Code = "issue-not-found"
	CodeRequirementNotFound    Code = "requirement-not-found"
	CodeLabelNotFound          Code = "label-not-found"
	CodeLinkNotFound           Code = "link-not-found"
	CodeWorkflowNotFound       Code = "workflow-not-found"
	CodeDuplicateEntity        Code = "duplicate-entity"
	CodeStateConflict          Code = "state-conflict"
	CodeLockUnavailable        Code = "lock-unavailable"
	CodeHierarchyInvalid       Code = "hierarchy-invalid"
	CodeDependencyCycle        Code = "dependency-cycle"
	CodeTransitionNotPermitted Code = "transition-not-permitted"
	CodeAuthorityMissing       Code = "authority-missing"
	CodeAuthorityUnknownActor  Code = "authority-unknown-actor"
	CodeAuthorityDenied        Code = "authority-denied"
	CodeJournalReplayStopped   Code = "journal-replay-stopped"
	CodeProjectionDivergence   Code = "projection-divergence"
	CodeResultLimitExceeded    Code = "result-limit-exceeded"
	CodeInternalError          Code = "internal-error"
)

type SubjectType string

const (
	SubjectProject     SubjectType = "project"
	SubjectIssue       SubjectType = "issue"
	SubjectRequirement SubjectType = "requirement"
	SubjectLabel       SubjectType = "label"
	SubjectLink        SubjectType = "link"
	SubjectWorkflow    SubjectType = "workflow"
)

type Subject struct {
	Type SubjectType `json:"type"`
	ID   string      `json:"id"`
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
}

type detailType uint8

const (
	detailString detailType = iota
	detailInteger
	detailStrings
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

func nullableIntegerField(name string) detailField {
	return detailField{name: name, typeOf: detailInteger, nullable: true}
}

func stringsField(name string) detailField { return detailField{name: name, typeOf: detailStrings} }

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

var reasonCodeDefinitions = map[Code]reasonDefinition{
	CodeInvalidArgument:     definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeUnsupportedKind:     definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeUnsupportedRelation: definition(CategoryInvalidInput, false, nil, stringField("argument"), nullableStringField("value"), stringsField("allowed")),
	CodeProjectNotFound:     definition(CategoryNotFound, false, []SubjectType{SubjectProject}, stringField("path")),
	CodeProjectAlreadyInit:  definition(CategoryConflict, false, []SubjectType{SubjectProject}, stringField("path")),
	CodeIssueNotFound:       definition(CategoryNotFound, false, []SubjectType{SubjectIssue}, stringField("id")),
	CodeRequirementNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectRequirement}, stringField("id")),
	CodeLabelNotFound:       definition(CategoryNotFound, false, []SubjectType{SubjectLabel}, stringField("id")),
	CodeLinkNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectLink},
		stringField("from_type"), stringField("from_id"), stringField("relation"), stringField("to_type"), stringField("to_id")),
	CodeWorkflowNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectWorkflow}, stringField("id")),
	CodeDuplicateEntity: definition(CategoryConflict, false, []SubjectType{SubjectIssue, SubjectLabel, SubjectLink, SubjectWorkflow},
		stringField("entity_type"), objectField("identity", stringsField("fields"), stringsField("values"))),
	CodeStateConflict:   definition(CategoryConflict, false, []SubjectType{SubjectIssue}, integerField("expected_revision"), integerField("actual_revision")),
	CodeLockUnavailable: definition(CategoryConflict, true, nil, stringField("lock_path")),
	CodeHierarchyInvalid: definition(CategoryPrecondition, false, []SubjectType{SubjectIssue},
		stringField("kind"), nullableStringField("parent_id"), nullableStringField("parent_kind"), stringField("rule")),
	CodeDependencyCycle: definition(CategoryPrecondition, false, nil, stringsField("path")),
	CodeTransitionNotPermitted: definition(CategoryPrecondition, false, []SubjectType{SubjectIssue},
		stringField("from_state"), stringField("to_state"), stringsField("unsatisfied")),
	CodeAuthorityMissing:      definition(CategoryAuthority, false, nil, stringField("operation"), nullableStringField("actor")),
	CodeAuthorityUnknownActor: definition(CategoryAuthority, false, nil, stringField("operation"), nullableStringField("actor")),
	CodeAuthorityDenied:       definition(CategoryAuthority, false, nil, stringField("operation"), nullableStringField("actor")),
	CodeJournalReplayStopped: definition(CategoryIntegrity, false, nil,
		stringField("journal_file"), nullableIntegerField("line"), nullableStringField("event_id")),
	CodeProjectionDivergence: definition(CategoryIntegrity, false, nil,
		stringField("status"), objectsField("errors", findingFields()...), objectsField("warnings", findingFields()...), objectField("info", infoFields()...)),
	CodeResultLimitExceeded: definition(CategoryIntegrity, false, nil, stringField("collection"), integerField("limit"), integerField("observed")),
	CodeInternalError:       definition(CategoryInternal, false, nil),
}

var registeredCodes = []Code{
	CodeInvalidArgument, CodeUnsupportedKind, CodeUnsupportedRelation, CodeProjectNotFound, CodeProjectAlreadyInit,
	CodeIssueNotFound, CodeRequirementNotFound, CodeLabelNotFound, CodeLinkNotFound, CodeWorkflowNotFound,
	CodeDuplicateEntity, CodeStateConflict, CodeLockUnavailable, CodeHierarchyInvalid, CodeDependencyCycle,
	CodeTransitionNotPermitted, CodeAuthorityMissing, CodeAuthorityUnknownActor, CodeAuthorityDenied,
	CodeJournalReplayStopped, CodeProjectionDivergence, CodeResultLimitExceeded, CodeInternalError,
}

func findingFields() []detailField {
	return []detailField{stringField("check"), stringField("message"), stringField("entity_id")}
}

func infoFields() []detailField {
	return []detailField{integerField("issue_count"), integerField("claim_count_live"), integerField("recent_activity_7d")}
}

func RegisteredCodes() []Code { return append([]Code(nil), registeredCodes...) }

func LookupCode(code Code) (ReasonSpec, bool) {
	definition, ok := reasonCodeDefinitions[code]
	if !ok {
		return ReasonSpec{}, false
	}
	spec := definition.spec
	spec.SubjectTypes = append([]SubjectType(nil), spec.SubjectTypes...)
	spec.DetailFields = append([]string{}, spec.DetailFields...)
	return spec, true
}

func NewFailure(code Code, message string, subject *Subject, details map[string]any) (Failure, error) {
	definition, ok := reasonCodeDefinitions[code]
	if !ok {
		return Failure{}, fmt.Errorf("new failure: unknown code %q", code)
	}
	failure := Failure{code, definition.spec.Category, definition.spec.Retryable, message, subject, details}
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
	if err := validateDetails(failure.Details, definition.details); err != nil {
		return err
	}
	if failure.Code == CodeDuplicateEntity {
		return validateDuplicateEntityIdentity(failure.Details)
	}
	return nil
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
