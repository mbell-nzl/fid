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
	CodeInvalidArgument             Code = "invalid-argument"
	CodeUnsupportedKind             Code = "unsupported-kind"
	CodeUnsupportedRelation         Code = "unsupported-relation"
	CodeProjectNotFound             Code = "project-not-found"
	CodeProjectAlreadyInit          Code = "project-already-initialized"
	CodeIssueNotFound               Code = "issue-not-found"
	CodeRequirementNotFound         Code = "requirement-not-found"
	CodeLabelNotFound               Code = "label-not-found"
	CodeLinkNotFound                Code = "link-not-found"
	CodeWorkflowNotFound            Code = "workflow-not-found"
	CodeDuplicateEntity             Code = "duplicate-entity"
	CodeStateConflict               Code = "state-conflict"
	CodeLockUnavailable             Code = "lock-unavailable"
	CodeHierarchyInvalid            Code = "hierarchy-invalid"
	CodeDependencyCycle             Code = "dependency-cycle"
	CodeTransitionNotPermitted      Code = "transition-not-permitted"
	CodeAuthorityMissing            Code = "authority-missing"
	CodeAuthorityUnknownActor       Code = "authority-unknown-actor"
	CodeAuthorityDenied             Code = "authority-denied"
	CodeJournalReplayStopped        Code = "journal-replay-stopped"
	CodeProjectionDivergence        Code = "projection-divergence"
	CodeResultLimitExceeded         Code = "result-limit-exceeded"
	CodeInternalError               Code = "internal-error"
	CodeActivationConflict          Code = "activation-conflict"
	CodeActivationNotFound          Code = "activation-not-found"
	CodeActivationInvalid           Code = "activation-invalid"
	CodeActivationVersionConflict   Code = "activation-version-conflict"
	CodeActivationDirectionConflict Code = "activation-direction-conflict"
	CodeActivationIncomplete        Code = "activation-incomplete"
	CodeAuthorityUnavailable        Code = "authority-unavailable"
	CodeCapabilityUnavailable       Code = "capability-unavailable"
	CodeExecutorPolicyUnsupported   Code = "executor-policy-unsupported"
	CodeRegistrationNotFound        Code = "registration-not-found"
	CodeEnvelopeNotFound            Code = "envelope-not-found"
	CodeClaimNotFound               Code = "claim-not-found"
	CodeClaimConflict               Code = "claim-conflict"
	CodeClaimExpired                Code = "claim-expired"
	CodeClaimStale                  Code = "claim-stale"
	CodeClaimModeInvalid            Code = "claim-mode-invalid"
	CodeRunNotFound                 Code = "run-not-found"
	CodeRunInactive                 Code = "run-inactive"
	CodeRunVersionConflict          Code = "run-version-conflict"
	CodeExecutorLivenessUnknown     Code = "executor-liveness-unknown"
	CodeExecutorStillLive           Code = "executor-still-live"
	CodeRecoveryStateConflict       Code = "recovery-state-conflict"
	CodeObservationNotFound         Code = "observation-not-found"
	CodeRecoveryAlreadyLinked       Code = "recovery-already-linked"
	CodeReceiptNotFound             Code = "receipt-not-found"
	CodeEvidenceNotFound            Code = "evidence-not-found"
	CodeEvidenceInvalid             Code = "evidence-invalid"
	CodeEvidenceStale               Code = "evidence-stale"
	CodeArtifactUnavailable         Code = "artifact-unavailable"
	CodeGateUnsuccessful            Code = "gate-unsuccessful"
	CodeBudgetExhausted             Code = "budget-exhausted"
	CodeBudgetDataUnavailable       Code = "budget-data-unavailable"
	CodeBudgetSampleConflict        Code = "budget-sample-conflict"
	CodeStopNotFound                Code = "stop-not-found"
	CodeAttentionNotFound           Code = "attention-not-found"
	CodeAttentionVersionConflict    Code = "attention-version-conflict"
	CodeStopUnresolved              Code = "stop-unresolved"
	CodeResolutionInvalid           Code = "resolution-invalid"
	CodeObligationUnsatisfied       Code = "obligation-unsatisfied"
	CodeArchiveInvalid              Code = "archive-invalid"
	CodeImportTargetConflict        Code = "import-target-conflict"
	CodeReviewNotIndependent        Code = "review-not-independent"
	CodeReviewVerdictInvalid        Code = "review-verdict-invalid"
	CodeRepositoryMismatch          Code = "repository-mismatch"
	CodeReviewNotFound              Code = "review-not-found"
	CodeIdempotencyConflict         Code = "idempotency-conflict"
)

type SubjectType string

const (
	SubjectProject          SubjectType = "project"
	SubjectIssue            SubjectType = "issue"
	SubjectRequirement      SubjectType = "requirement"
	SubjectLabel            SubjectType = "label"
	SubjectLink             SubjectType = "link"
	SubjectWorkflow         SubjectType = "workflow"
	SubjectActivation       SubjectType = "activation"
	SubjectGateDefinition   SubjectType = "gate-definition"
	SubjectGateRegistration SubjectType = "gate-registration"
	SubjectEnvelope         SubjectType = "envelope"
	SubjectClaim            SubjectType = "claim"
	SubjectRun              SubjectType = "run"
	SubjectReceipt          SubjectType = "receipt"
	SubjectEvidence         SubjectType = "evidence"
	SubjectArtifact         SubjectType = "artifact"
	SubjectCapability       SubjectType = "capability"
	SubjectExport           SubjectType = "export"
	SubjectStop             SubjectType = "stop"
	SubjectAttention        SubjectType = "attention"
	SubjectReview           SubjectType = "review"
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

func nullableIntegerField(name string) detailField {
	return detailField{name: name, typeOf: detailInteger, nullable: true}
}

func stringsField(name string) detailField { return detailField{name: name, typeOf: detailStrings} }

func booleanField(name string) detailField { return detailField{name: name, typeOf: detailBoolean} }

func objectField(name string, fields ...detailField) detailField {
	return detailField{name: name, typeOf: detailObject, fields: fields}
}

func nullableObjectField(name string, fields ...detailField) detailField {
	return detailField{name: name, typeOf: detailObject, nullable: true, fields: fields}
}

func objectsField(name string, fields ...detailField) detailField {
	return detailField{name: name, typeOf: detailObjects, fields: fields}
}

func repositoryFields() []detailField {
	return []detailField{
		nullableStringField("base_oid"), stringField("capture_source"), booleanField("dirty"), nullableStringField("dirty_digest"),
		stringField("head_oid"), stringField("project_id"), stringField("tree_oid"), stringField("worktree_id"),
	}
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
	CodeReviewNotIndependent: definition(CategoryPrecondition, false, nil,
		stringField("review_actor_id"), objectField("conflicting_subject", stringField("type"), stringField("id")), stringField("conflicting_actor_id")),
	CodeReviewVerdictInvalid: definition(CategoryInvalidInput, false, nil,
		stringField("verdict"), stringField("reason")),
	CodeRepositoryMismatch: definition(CategoryConflict, false, nil,
		objectField("expected", repositoryFields()...), objectField("observed", repositoryFields()...)),
	CodeReviewNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectReview}, stringField("id")),
	CodeIdempotencyConflict: definition(CategoryConflict, false, nil,
		stringField("operation_scope"), stringField("idempotency_key"), stringField("stored_request_digest"), stringField("supplied_request_digest")),
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
	CodeDuplicateEntity: definition(CategoryConflict, false, []SubjectType{SubjectIssue, SubjectLabel, SubjectLink, SubjectWorkflow, SubjectRun},
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
	CodeActivationConflict: definition(CategoryConflict, false, []SubjectType{SubjectActivation},
		nullableStringField("active_id"), nullableStringField("nonterminal_id"), stringField("requested_direction")),
	CodeActivationNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectActivation},
		stringField("id")),
	CodeActivationInvalid: definition(CategoryPrecondition, false, []SubjectType{SubjectActivation},
		stringField("status"), stringsField("expected_statuses"), objectsField("component_errors", stringField("component_id"), stringField("code"))),
	CodeActivationVersionConflict: definition(CategoryConflict, false, []SubjectType{SubjectActivation},
		stringField("activation_id"), integerField("expected_version"), integerField("actual_version")),
	CodeActivationDirectionConflict: definition(CategoryPrecondition, false, []SubjectType{SubjectActivation},
		stringField("activation_id"), booleanField("point_of_no_return"), stringField("required_direction")),
	CodeActivationIncomplete: definition(CategoryPrecondition, false, []SubjectType{SubjectActivation},
		stringField("activation_id"), stringsField("components"), objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeAuthorityUnavailable: definition(CategoryAuthority, false, nil,
		nullableStringField("actor_id"), stringField("operation"), stringsField("required_actions"), stringsField("required_capabilities")),
	CodeCapabilityUnavailable: definition(CategoryAuthority, false, []SubjectType{SubjectCapability},
		stringField("capability_id"), booleanField("available"), integerField("owner_phase")),
	CodeExecutorPolicyUnsupported: definition(CategoryInvalidInput, false, []SubjectType{SubjectGateDefinition},
		stringsField("fields"), objectField("supported", stringField("network"), booleanField("shell"), stringField("sandbox"), stringField("write"))),
	CodeRegistrationNotFound: definition(CategoryNotFound, false, []SubjectType{SubjectGateDefinition, SubjectGateRegistration}, stringField("id")),
	CodeEnvelopeNotFound:     definition(CategoryNotFound, false, []SubjectType{SubjectEnvelope}, stringField("id")),
	CodeClaimNotFound:        definition(CategoryNotFound, false, []SubjectType{SubjectClaim}, stringField("id")),
	CodeClaimConflict:        definition(CategoryConflict, false, []SubjectType{SubjectClaim}, stringField("issue_id"), nullableStringField("claim_id"), nullableIntegerField("expected_version"), nullableIntegerField("actual_version"), nullableStringField("expires_at")),
	CodeClaimExpired:         definition(CategoryConflict, false, []SubjectType{SubjectClaim}, stringField("issue_id"), nullableStringField("claim_id"), nullableIntegerField("expected_version"), nullableIntegerField("actual_version"), nullableStringField("expires_at")),
	CodeClaimStale:           definition(CategoryConflict, false, []SubjectType{SubjectClaim}, stringField("issue_id"), nullableStringField("claim_id"), nullableIntegerField("expected_version"), nullableIntegerField("actual_version"), nullableStringField("expires_at")),
	CodeClaimModeInvalid:     definition(CategoryPrecondition, false, []SubjectType{SubjectClaim}, stringField("claim_id"), stringField("mode"), stringField("required_mode")),
	CodeRunNotFound:          definition(CategoryNotFound, false, []SubjectType{SubjectRun}, stringField("id")),
	CodeRunInactive:          definition(CategoryPrecondition, false, []SubjectType{SubjectRun}, stringField("run_id"), stringField("status")),
	CodeRunVersionConflict:   definition(CategoryConflict, false, []SubjectType{SubjectRun}, stringField("run_id"), integerField("expected_version"), integerField("actual_version")),
	CodeExecutorLivenessUnknown: definition(CategoryPrecondition, false, []SubjectType{SubjectRun},
		stringField("run_id"), stringField("executor_id"), nullableStringField("observation_id")),
	CodeExecutorStillLive: definition(CategoryPrecondition, false, []SubjectType{SubjectRun},
		stringField("run_id"), stringField("executor_id"), stringField("observation_id")),
	CodeRecoveryStateConflict: definition(CategoryConflict, false, []SubjectType{SubjectRun},
		stringField("run_id"), nullableStringField("observation_digest"), stringsField("paths"), stringsField("reasons")),
	CodeObservationNotFound:   definition(CategoryNotFound, false, []SubjectType{SubjectRun}, stringField("id")),
	CodeRecoveryAlreadyLinked: definition(CategoryConflict, false, []SubjectType{SubjectRun}, stringField("run_id"), stringField("successor_run_id")),
	CodeReceiptNotFound:       definition(CategoryNotFound, false, []SubjectType{SubjectReceipt}, stringField("id")),
	CodeEvidenceNotFound:      definition(CategoryNotFound, false, []SubjectType{SubjectEvidence}, stringField("id")),
	CodeEvidenceInvalid:       definition(CategoryPrecondition, false, []SubjectType{SubjectEvidence}, objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeEvidenceStale:         definition(CategoryPrecondition, false, []SubjectType{SubjectEvidence}, objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeArtifactUnavailable:   definition(CategoryPrecondition, false, []SubjectType{SubjectArtifact}, objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeGateUnsuccessful: definition(CategoryGateUnsuccessful, false, []SubjectType{SubjectReceipt},
		stringField("attempt_id"), stringField("receipt_id"), stringField("outcome"), nullableIntegerField("exit_code"), nullableStringField("signal"), booleanField("timed_out")),
	CodeBudgetExhausted: definition(CategoryPrecondition, false, []SubjectType{SubjectRun},
		objectField("budget", stringField("name"), stringField("scope"), stringField("unit"), integerField("maximum"), nullableIntegerField("warning_threshold"), stringField("exhaustion_action"), nullableIntegerField("consumed"), nullableIntegerField("remaining"), stringField("state")),
		stringsField("stop_ids"), stringsField("attention_ids"), nullableObjectField("event", stringField("id"), stringField("occurred_at")), booleanField("replayed")),
	CodeBudgetSampleConflict: definition(CategoryConflict, false, []SubjectType{SubjectRun},
		stringField("budget_name"), stringField("source_sample_id"), nullableStringField("expected_previous_sample_id"), nullableStringField("actual_previous_sample_id")),
	CodeBudgetDataUnavailable: definition(CategoryPrecondition, false, []SubjectType{SubjectRun},
		objectField("budget", stringField("name"), stringField("scope"), stringField("unit"), integerField("maximum"), nullableIntegerField("warning_threshold"), stringField("exhaustion_action"), nullableIntegerField("consumed"), nullableIntegerField("remaining"), stringField("state")),
		stringsField("stop_ids"), stringsField("attention_ids"), nullableObjectField("event", stringField("id"), stringField("occurred_at")), booleanField("replayed")),
	CodeStopNotFound:             definition(CategoryNotFound, false, []SubjectType{SubjectStop}, stringField("id")),
	CodeAttentionNotFound:        definition(CategoryNotFound, false, []SubjectType{SubjectAttention}, stringField("id")),
	CodeAttentionVersionConflict: definition(CategoryConflict, false, []SubjectType{SubjectAttention}, stringField("attention_id"), integerField("expected_version"), integerField("actual_version")),
	CodeStopUnresolved:           definition(CategoryPrecondition, false, []SubjectType{SubjectStop}, stringsField("stop_ids"), objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeResolutionInvalid:        definition(CategoryPrecondition, false, []SubjectType{SubjectStop, SubjectAttention}, objectField("subject", stringField("type"), stringField("id")), stringsField("fact_event_ids"), objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeObligationUnsatisfied: definition(CategoryPrecondition, false, []SubjectType{SubjectEnvelope},
		objectsField("reasons", stringField("code"), stringField("detail"), stringField("subject_id"), stringField("subject_type"))),
	CodeArchiveInvalid: definition(CategoryIntegrity, false, []SubjectType{SubjectExport},
		stringField("archive_path"), objectsField("checks", stringField("code"), nullableStringField("entry"), stringField("detail"))),
	CodeImportTargetConflict: definition(CategoryConflict, false, []SubjectType{SubjectExport},
		stringField("target_root"), stringField("reason")),
}

var registeredCodes = []Code{
	CodeInvalidArgument, CodeUnsupportedKind, CodeUnsupportedRelation, CodeProjectNotFound, CodeProjectAlreadyInit,
	CodeIssueNotFound, CodeRequirementNotFound, CodeLabelNotFound, CodeLinkNotFound, CodeWorkflowNotFound,
	CodeDuplicateEntity, CodeStateConflict, CodeLockUnavailable, CodeHierarchyInvalid, CodeDependencyCycle,
	CodeTransitionNotPermitted, CodeAuthorityMissing, CodeAuthorityUnknownActor, CodeAuthorityDenied,
	CodeJournalReplayStopped, CodeProjectionDivergence, CodeResultLimitExceeded, CodeInternalError,
	CodeActivationConflict,
	CodeActivationNotFound,
	CodeActivationInvalid,
	CodeActivationVersionConflict,
	CodeActivationDirectionConflict,
	CodeActivationIncomplete,
	CodeAuthorityUnavailable, CodeCapabilityUnavailable, CodeExecutorPolicyUnsupported, CodeRegistrationNotFound, CodeEnvelopeNotFound,
	CodeClaimNotFound, CodeClaimConflict, CodeClaimExpired, CodeClaimStale, CodeClaimModeInvalid, CodeRunNotFound, CodeRunInactive, CodeRunVersionConflict,
	CodeExecutorLivenessUnknown, CodeExecutorStillLive, CodeRecoveryStateConflict, CodeObservationNotFound, CodeRecoveryAlreadyLinked,
	CodeReceiptNotFound, CodeEvidenceNotFound,
	CodeGateUnsuccessful,
	CodeBudgetExhausted, CodeBudgetDataUnavailable, CodeBudgetSampleConflict,
	CodeStopNotFound, CodeAttentionNotFound, CodeAttentionVersionConflict, CodeStopUnresolved, CodeResolutionInvalid,
	CodeObligationUnsatisfied,
	CodeArchiveInvalid,
	CodeImportTargetConflict,
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
