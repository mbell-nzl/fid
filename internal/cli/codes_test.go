package cli

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

type detailKind uint8

const (
	detailStringKind detailKind = iota
	detailIntegerKind
	detailStringsKind
	detailBooleanKind
	detailObjectKind
	detailObjectsKind
)

type detailExpectation struct {
	name     string
	kind     detailKind
	nullable bool
	children []detailExpectation
	value    any
}

type reasonCodeExpectation struct {
	code      Code
	category  Category
	retryable bool
	subjects  []SubjectType
	details   []detailExpectation
}

func TestReasonCodeDetailSchemaIsExactForEveryCode(t *testing.T) {
	expectations := []reasonCodeExpectation{
		{CodeInvalidArgument, CategoryInvalidInput, false, nil, argumentDetails()},
		{CodeUnsupportedKind, CategoryInvalidInput, false, nil, argumentDetails()},
		{CodeUnsupportedRelation, CategoryInvalidInput, false, nil, argumentDetails()},
		{CodeProjectNotFound, CategoryNotFound, false, []SubjectType{SubjectProject}, stringDetails("path")},
		{CodeProjectAlreadyInit, CategoryConflict, false, []SubjectType{SubjectProject}, stringDetails("path")},
		{CodeIssueNotFound, CategoryNotFound, false, []SubjectType{SubjectIssue}, stringDetails("id")},
		{CodeRequirementNotFound, CategoryNotFound, false, []SubjectType{SubjectRequirement}, stringDetails("id")},
		{CodeLabelNotFound, CategoryNotFound, false, []SubjectType{SubjectLabel}, stringDetails("id")},
		{CodeLinkNotFound, CategoryNotFound, false, []SubjectType{SubjectLink}, []detailExpectation{
			stringDetail("from_type"), stringDetail("from_id"), stringDetail("relation"), stringDetail("to_type"), stringDetail("to_id"),
		}},
		{CodeWorkflowNotFound, CategoryNotFound, false, []SubjectType{SubjectWorkflow}, stringDetails("id")},
		{CodeDuplicateEntity, CategoryConflict, false, []SubjectType{SubjectIssue, SubjectLabel, SubjectLink, SubjectWorkflow, SubjectRun}, []detailExpectation{
			stringDetail("entity_type"),
			{name: "identity", kind: detailObjectKind, children: []detailExpectation{
				{name: "fields", kind: detailStringsKind, value: []string{"name", "colour"}},
				{name: "values", kind: detailStringsKind, value: []string{"bug", "red"}},
			}},
		}},
		{CodeStateConflict, CategoryConflict, false, []SubjectType{SubjectIssue}, []detailExpectation{integerDetail("expected_revision"), integerDetail("actual_revision")}},
		{CodeLockUnavailable, CategoryConflict, true, nil, stringDetails("lock_path")},
		{CodeHierarchyInvalid, CategoryPrecondition, false, []SubjectType{SubjectIssue}, []detailExpectation{
			stringDetail("kind"), nullableStringDetail("parent_id"), nullableStringDetail("parent_kind"), stringDetail("rule"),
		}},
		{CodeDependencyCycle, CategoryPrecondition, false, nil, []detailExpectation{{name: "path", kind: detailStringsKind}}},
		{CodeTransitionNotPermitted, CategoryPrecondition, false, []SubjectType{SubjectIssue}, []detailExpectation{
			stringDetail("from_state"), stringDetail("to_state"), {name: "unsatisfied", kind: detailStringsKind},
		}},
		{CodeAuthorityMissing, CategoryAuthority, false, nil, authorityDetails()},
		{CodeAuthorityUnknownActor, CategoryAuthority, false, nil, authorityDetails()},
		{CodeAuthorityDenied, CategoryAuthority, false, nil, authorityDetails()},
		{CodeJournalReplayStopped, CategoryIntegrity, false, nil, []detailExpectation{
			stringDetail("journal_file"), nullableIntegerDetail("line"), nullableStringDetail("event_id"),
		}},
		{CodeProjectionDivergence, CategoryIntegrity, false, nil, []detailExpectation{
			stringDetail("status"),
			{name: "errors", kind: detailObjectsKind, children: findingDetails()},
			{name: "warnings", kind: detailObjectsKind, children: findingDetails()},
			{name: "info", kind: detailObjectKind, children: []detailExpectation{
				integerDetail("issue_count"), integerDetail("claim_count_live"), integerDetail("recent_activity_7d"),
			}},
		}},
		{CodeResultLimitExceeded, CategoryIntegrity, false, nil, []detailExpectation{stringDetail("collection"), integerDetail("limit"), integerDetail("observed")}},
		{CodeInternalError, CategoryInternal, false, nil, nil},
		{CodeActivationConflict, CategoryConflict, false, []SubjectType{SubjectActivation}, []detailExpectation{
			nullableStringDetail("active_id"), nullableStringDetail("nonterminal_id"), stringDetail("requested_direction"),
		}},
		{CodeActivationNotFound, CategoryNotFound, false, []SubjectType{SubjectActivation}, stringDetails("id")},
		{CodeActivationInvalid, CategoryPrecondition, false, []SubjectType{SubjectActivation}, []detailExpectation{
			stringDetail("status"), {name: "expected_statuses", kind: detailStringsKind}, {name: "component_errors", kind: detailObjectsKind, children: []detailExpectation{stringDetail("component_id"), stringDetail("code")}},
		}},
		{CodeActivationVersionConflict, CategoryConflict, false, []SubjectType{SubjectActivation}, []detailExpectation{stringDetail("activation_id"), integerDetail("expected_version"), integerDetail("actual_version")}},
		{CodeActivationDirectionConflict, CategoryPrecondition, false, []SubjectType{SubjectActivation}, []detailExpectation{stringDetail("activation_id"), {name: "point_of_no_return", kind: detailBooleanKind}, stringDetail("required_direction")}},
		{CodeActivationIncomplete, CategoryPrecondition, false, []SubjectType{SubjectActivation}, []detailExpectation{
			stringDetail("activation_id"), {name: "components", kind: detailStringsKind}, {name: "reasons", kind: detailObjectsKind, children: []detailExpectation{stringDetail("code"), stringDetail("detail"), stringDetail("subject_id"), stringDetail("subject_type")}},
		}},
		{CodeAuthorityUnavailable, CategoryAuthority, false, nil, []detailExpectation{
			nullableStringDetail("actor_id"), stringDetail("operation"), {name: "required_actions", kind: detailStringsKind}, {name: "required_capabilities", kind: detailStringsKind},
		}},
		{CodeCapabilityUnavailable, CategoryAuthority, false, []SubjectType{SubjectCapability}, []detailExpectation{
			stringDetail("capability_id"), {name: "available", kind: detailBooleanKind}, integerDetail("owner_phase"),
		}},
		{CodeExecutorPolicyUnsupported, CategoryInvalidInput, false, []SubjectType{SubjectGateDefinition}, []detailExpectation{
			{name: "fields", kind: detailStringsKind}, {name: "supported", kind: detailObjectKind, children: []detailExpectation{
				stringDetail("network"), {name: "shell", kind: detailBooleanKind}, stringDetail("sandbox"), stringDetail("write"),
			}},
		}},
		{CodeRegistrationNotFound, CategoryNotFound, false, []SubjectType{SubjectGateDefinition, SubjectGateRegistration}, stringDetails("id")},
		{CodeEnvelopeNotFound, CategoryNotFound, false, []SubjectType{SubjectEnvelope}, stringDetails("id")},
		{CodeClaimNotFound, CategoryNotFound, false, []SubjectType{SubjectClaim}, stringDetails("id")},
		{CodeClaimConflict, CategoryConflict, false, []SubjectType{SubjectClaim}, claimConflictDetails()},
		{CodeClaimExpired, CategoryConflict, false, []SubjectType{SubjectClaim}, claimConflictDetails()},
		{CodeClaimStale, CategoryConflict, false, []SubjectType{SubjectClaim}, claimConflictDetails()},
		{CodeClaimModeInvalid, CategoryPrecondition, false, []SubjectType{SubjectClaim}, []detailExpectation{stringDetail("claim_id"), stringDetail("mode"), stringDetail("required_mode")}},
		{CodeRunNotFound, CategoryNotFound, false, []SubjectType{SubjectRun}, stringDetails("id")},
		{CodeRunInactive, CategoryPrecondition, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), stringDetail("status")}},
		{Code("run-version-conflict"), CategoryConflict, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), integerDetail("expected_version"), integerDetail("actual_version")}},
		{CodeExecutorLivenessUnknown, CategoryPrecondition, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), stringDetail("executor_id"), nullableStringDetail("observation_id")}},
		{CodeExecutorStillLive, CategoryPrecondition, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), stringDetail("executor_id"), stringDetail("observation_id")}},
		{CodeRecoveryStateConflict, CategoryConflict, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), nullableStringDetail("observation_digest"), {name: "paths", kind: detailStringsKind}, {name: "reasons", kind: detailStringsKind}}},
		{CodeObservationNotFound, CategoryNotFound, false, []SubjectType{SubjectRun}, stringDetails("id")},
		{CodeRecoveryAlreadyLinked, CategoryConflict, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("run_id"), stringDetail("successor_run_id")}},
		{CodeReceiptNotFound, CategoryNotFound, false, []SubjectType{SubjectReceipt}, stringDetails("id")},
		{CodeEvidenceNotFound, CategoryNotFound, false, []SubjectType{SubjectEvidence}, stringDetails("id")},
		{CodeReviewChangesRequired, CategoryPrecondition, false, nil, []detailExpectation{stringDetail("review_id"), stringDetail("obligation_id"), {name: "blocking_finding_ids", kind: detailStringsKind}}},
		{CodeReviewReworkExhausted, CategoryPrecondition, false, []SubjectType{SubjectEnvelope}, budgetDenialDetails()},
		{CodeEvidenceTypeNotAccepted, CategoryPrecondition, false, nil, []detailExpectation{stringDetail("obligation_id"), stringDetail("evidence_type"), {name: "accepted_evidence_types", kind: detailStringsKind}}},
		{CodeOutcomeNotPermitted, CategoryPrecondition, false, nil, []detailExpectation{stringDetail("obligation_id"), stringDetail("obligation_type"), stringDetail("outcome_id"), stringDetail("result")}},
		{CodeGateUnsuccessful, CategoryGateUnsuccessful, false, []SubjectType{SubjectReceipt}, []detailExpectation{
			stringDetail("attempt_id"), stringDetail("receipt_id"), stringDetail("outcome"), nullableIntegerDetail("exit_code"), nullableStringDetail("signal"), {name: "timed_out", kind: detailBooleanKind},
		}},
		{CodeBudgetExhausted, CategoryPrecondition, false, []SubjectType{SubjectRun}, budgetDenialDetails()},
		{CodeBudgetDataUnavailable, CategoryPrecondition, false, []SubjectType{SubjectRun}, budgetDenialDetails()},
		{CodeBudgetSampleConflict, CategoryConflict, false, []SubjectType{SubjectRun}, []detailExpectation{stringDetail("budget_name"), stringDetail("source_sample_id"), nullableStringDetail("expected_previous_sample_id"), nullableStringDetail("actual_previous_sample_id")}},
		{CodeStopNotFound, CategoryNotFound, false, []SubjectType{SubjectStop}, stringDetails("id")},
		{CodeAttentionNotFound, CategoryNotFound, false, []SubjectType{SubjectAttention}, stringDetails("id")},
		{CodeAttentionVersionConflict, CategoryConflict, false, []SubjectType{SubjectAttention}, []detailExpectation{stringDetail("attention_id"), integerDetail("expected_version"), integerDetail("actual_version")}},
		{CodeStopUnresolved, CategoryPrecondition, false, []SubjectType{SubjectStop}, []detailExpectation{{name: "stop_ids", kind: detailStringsKind}, {name: "reasons", kind: detailObjectsKind, children: []detailExpectation{stringDetail("code"), stringDetail("detail"), stringDetail("subject_id"), stringDetail("subject_type")}}}},
		{CodeResolutionInvalid, CategoryPrecondition, false, []SubjectType{SubjectStop, SubjectAttention}, []detailExpectation{{name: "subject", kind: detailObjectKind, children: []detailExpectation{stringDetail("type"), stringDetail("id")}}, {name: "fact_event_ids", kind: detailStringsKind}, {name: "reasons", kind: detailObjectsKind, children: []detailExpectation{stringDetail("code"), stringDetail("detail"), stringDetail("subject_id"), stringDetail("subject_type")}}}},
		{CodeObligationUnsatisfied, CategoryPrecondition, false, []SubjectType{SubjectEnvelope}, []detailExpectation{
			{name: "reasons", kind: detailObjectsKind, children: []detailExpectation{
				stringDetail("code"), stringDetail("detail"), stringDetail("subject_id"), stringDetail("subject_type"),
			}},
		}},
		{CodeArchiveInvalid, CategoryIntegrity, false, []SubjectType{SubjectExport}, []detailExpectation{
			stringDetail("archive_path"), {name: "checks", kind: detailObjectsKind, children: []detailExpectation{
				stringDetail("code"), nullableStringDetail("entry"), stringDetail("detail"),
			}},
		}},
		{CodeImportTargetConflict, CategoryConflict, false, []SubjectType{SubjectExport}, []detailExpectation{stringDetail("target_root"), stringDetail("reason")}},
		{CodeFindingNotFound, CategoryNotFound, false, []SubjectType{SubjectFinding}, stringDetails("id")},
		{CodeFindingAlreadyDisposed, CategoryConflict, false, []SubjectType{SubjectFinding}, []detailExpectation{stringDetail("finding_id"), stringDetail("state")}},
		{CodeConcernNotFound, CategoryNotFound, false, []SubjectType{SubjectConcern}, stringDetails("id")},
		{CodeDeviationNotFound, CategoryNotFound, false, []SubjectType{SubjectDeviation}, stringDetails("id")},
		{CodeConcernAlreadyDisposed, CategoryConflict, false, []SubjectType{SubjectConcern}, []detailExpectation{stringDetail("concern_id"), stringDetail("disposed_at")}},
	}

	registered := RegisteredCodes()
	if len(registered) != len(expectations) {
		t.Fatalf("RegisteredCodes() length = %d, want %d", len(registered), len(expectations))
	}
	for index, expectation := range expectations {
		if registered[index] != expectation.code {
			t.Fatalf("RegisteredCodes()[%d] = %q, want %q", index, registered[index], expectation.code)
		}
	}

	for _, expectation := range expectations {
		expectation := expectation
		t.Run(string(expectation.code), func(t *testing.T) {
			spec, ok := LookupCode(expectation.code)
			if !ok {
				t.Fatalf("LookupCode(%q) found no specification", expectation.code)
			}
			if spec.Category != expectation.category || spec.Retryable != expectation.retryable {
				t.Fatalf("LookupCode(%q) = category %q retryable %t, want category %q retryable %t", expectation.code, spec.Category, spec.Retryable, expectation.category, expectation.retryable)
			}
			if !reflect.DeepEqual(spec.SubjectTypes, expectation.subjects) {
				t.Fatalf("LookupCode(%q) subject types = %v, want %v", expectation.code, spec.SubjectTypes, expectation.subjects)
			}
			if got, want := spec.DetailFields, detailNames(expectation.details); !reflect.DeepEqual(got, want) {
				t.Fatalf("LookupCode(%q) detail fields = %v, want %v", expectation.code, got, want)
			}
			if spec.RedactionPolicy != RedactClosedDetails {
				t.Fatalf("LookupCode(%q) redaction policy = %q, want %q", expectation.code, spec.RedactionPolicy, RedactClosedDetails)
			}

			validDetails := detailValues(expectation.details)
			if expectation.code == CodeImportTargetConflict {
				validDetails["reason"] = "populated"
			}
			validSubject := subjectFor(expectation.subjects)
			failure := mustNewFailure(t, expectation.code, validSubject, validDetails)
			assertSerializedFailureSchema(t, failure, expectation)

			assertSubjectRules(t, expectation, validDetails)
			assertDetailRules(t, expectation, validSubject, validDetails)
		})
	}

	for _, code := range []Code{"", "not-a-contract-code"} {
		if _, err := NewFailure(code, "message", nil, map[string]any{}); err == nil {
			t.Errorf("NewFailure(%q) error = nil, want unknown-code rejection", code)
		}
	}
}

func TestImportTargetConflictAcceptsOnlyClosedReasons(t *testing.T) {
	for _, reason := range []string{"populated", "different-import", "unsafe-path"} {
		t.Run(reason, func(t *testing.T) {
			_, err := NewFailure(CodeImportTargetConflict, "conflict", &Subject{Type: SubjectExport, ID: "archive.zip"}, map[string]any{"target_root": "/target", "reason": reason})
			if err != nil {
				t.Fatalf("NewFailure(%q) error = %v", reason, err)
			}
		})
	}
	if _, err := NewFailure(CodeImportTargetConflict, "conflict", &Subject{Type: SubjectExport, ID: "archive.zip"}, map[string]any{"target_root": "/target", "reason": "other"}); err == nil {
		t.Fatal("NewFailure accepted an open import-target-conflict reason")
	}
}

func budgetDenialDetails() []detailExpectation {
	return []detailExpectation{
		{name: "budget", kind: detailObjectKind, children: []detailExpectation{stringDetail("name"), stringDetail("scope"), stringDetail("unit"), integerDetail("maximum"), nullableIntegerDetail("warning_threshold"), stringDetail("exhaustion_action"), nullableIntegerDetail("consumed"), nullableIntegerDetail("remaining"), stringDetail("state")}},
		{name: "stop_ids", kind: detailStringsKind}, {name: "attention_ids", kind: detailStringsKind},
		{name: "event", kind: detailObjectKind, nullable: true, children: []detailExpectation{stringDetail("id"), stringDetail("occurred_at")}}, {name: "replayed", kind: detailBooleanKind},
	}
}

func argumentDetails() []detailExpectation {
	return []detailExpectation{stringDetail("argument"), nullableStringDetail("value"), {name: "allowed", kind: detailStringsKind}}
}

func authorityDetails() []detailExpectation {
	return []detailExpectation{stringDetail("operation"), nullableStringDetail("actor")}
}

func claimConflictDetails() []detailExpectation {
	return []detailExpectation{
		stringDetail("issue_id"), nullableStringDetail("claim_id"), nullableIntegerDetail("expected_version"),
		nullableIntegerDetail("actual_version"), nullableStringDetail("expires_at"),
	}
}

func findingDetails() []detailExpectation {
	return []detailExpectation{stringDetail("check"), stringDetail("message"), stringDetail("entity_id")}
}

func stringDetails(name string) []detailExpectation { return []detailExpectation{stringDetail(name)} }

func stringDetail(name string) detailExpectation {
	return detailExpectation{name: name, kind: detailStringKind}
}

func nullableStringDetail(name string) detailExpectation {
	return detailExpectation{name: name, kind: detailStringKind, nullable: true}
}

func integerDetail(name string) detailExpectation {
	return detailExpectation{name: name, kind: detailIntegerKind}
}

func nullableIntegerDetail(name string) detailExpectation {
	return detailExpectation{name: name, kind: detailIntegerKind, nullable: true}
}

func detailNames(fields []detailExpectation) []string {
	names := make([]string, len(fields))
	for index, field := range fields {
		names[index] = field.name
	}
	return names
}

func detailValues(fields []detailExpectation) map[string]any {
	values := make(map[string]any, len(fields))
	for _, field := range fields {
		values[field.name] = detailValue(field)
	}
	return values
}

func detailValue(field detailExpectation) any {
	if field.value != nil {
		return cloneDetailValue(field.value)
	}
	switch field.kind {
	case detailStringKind:
		return field.name + "-value"
	case detailIntegerKind:
		return int64(7)
	case detailStringsKind:
		return []string{field.name + "-one", field.name + "-two"}
	case detailBooleanKind:
		return true
	case detailObjectKind:
		return detailValues(field.children)
	case detailObjectsKind:
		return []map[string]any{detailValues(field.children)}
	default:
		panic("unknown detail expectation kind")
	}
}

func cloneDetails(details map[string]any) map[string]any {
	return cloneDetailValue(details).(map[string]any)
}

func cloneDetailValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		clone := make(map[string]any, len(typed))
		for key, item := range typed {
			clone[key] = cloneDetailValue(item)
		}
		return clone
	case []string:
		return append([]string(nil), typed...)
	case []map[string]any:
		clone := make([]map[string]any, len(typed))
		for index, item := range typed {
			clone[index] = cloneDetails(item)
		}
		return clone
	default:
		return value
	}
}

func subjectFor(subjects []SubjectType) *Subject {
	if len(subjects) == 0 {
		return nil
	}
	return &Subject{Type: subjects[0], ID: "subject-id"}
}

func mustNewFailure(t *testing.T, code Code, subject *Subject, details map[string]any) Failure {
	t.Helper()
	failure, err := NewFailure(code, "message", subject, details)
	if err != nil {
		t.Fatalf("NewFailure(%q) error = %v", code, err)
	}
	return failure
}

func assertSubjectRules(t *testing.T, expectation reasonCodeExpectation, details map[string]any) {
	t.Helper()
	allSubjectTypes := []SubjectType{SubjectProject, SubjectIssue, SubjectRequirement, SubjectLabel, SubjectLink, SubjectWorkflow}
	if len(expectation.subjects) == 0 {
		assertNewFailureRejected(t, expectation.code, &Subject{Type: SubjectIssue, ID: "subject-id"}, details)
		return
	}

	assertNewFailureRejected(t, expectation.code, nil, details)
	for _, subjectType := range allSubjectTypes {
		subject := &Subject{Type: subjectType, ID: "subject-id"}
		if containsSubjectType(expectation.subjects, subjectType) {
			mustNewFailure(t, expectation.code, subject, cloneDetails(details))
			continue
		}
		assertNewFailureRejected(t, expectation.code, subject, details)
	}
	assertNewFailureRejected(t, expectation.code, &Subject{Type: "unknown", ID: "subject-id"}, details)
}

func containsSubjectType(subjects []SubjectType, want SubjectType) bool {
	for _, subject := range subjects {
		if subject == want {
			return true
		}
	}
	return false
}

func assertDetailRules(t *testing.T, expectation reasonCodeExpectation, subject *Subject, validDetails map[string]any) {
	t.Helper()
	paths := detailPaths(expectation.details, nil)
	if len(paths) == 0 {
		assertNewFailureRejected(t, expectation.code, subject, map[string]any{"unexpected": "value"})
	}
	for _, path := range paths {
		path := path
		t.Run("details_"+strings.Join(path.names, "_"), func(t *testing.T) {
			missing := cloneDetails(validDetails)
			delete(detailParent(missing, path.names), path.names[len(path.names)-1])
			assertNewFailureRejected(t, expectation.code, subject, missing)

			extra := cloneDetails(validDetails)
			detailParent(extra, path.names)["unexpected"] = "value"
			assertNewFailureRejected(t, expectation.code, subject, extra)

			wrongType := cloneDetails(validDetails)
			wrongValue := any(true)
			if path.kind == detailBooleanKind {
				wrongValue = "not-a-boolean"
			}
			detailParent(wrongType, path.names)[path.names[len(path.names)-1]] = wrongValue
			assertNewFailureRejected(t, expectation.code, subject, wrongType)

			if path.nullable {
				nullable := cloneDetails(validDetails)
				detailParent(nullable, path.names)[path.names[len(path.names)-1]] = nil
				mustNewFailure(t, expectation.code, subject, nullable)
				return
			}
			null := cloneDetails(validDetails)
			detailParent(null, path.names)[path.names[len(path.names)-1]] = nil
			assertNewFailureRejected(t, expectation.code, subject, null)
		})
	}

	if expectation.code == CodeDuplicateEntity {
		unequal := cloneDetails(validDetails)
		identity := unequal["identity"].(map[string]any)
		identity["values"] = []string{"bug"}
		assertNewFailureRejected(t, expectation.code, subject, unequal)
	}

	for _, redacted := range []string{"environment", "config", "journal-payload", "SQL", "file-content", "credential", "stack-trace"} {
		redacted := redacted
		t.Run("redacts_"+redacted, func(t *testing.T) {
			details := cloneDetails(validDetails)
			details[redacted] = "secret"
			assertNewFailureRejected(t, expectation.code, subject, details)
		})
	}
}

type detailPath struct {
	names    []string
	nullable bool
	kind     detailKind
}

func detailPaths(fields []detailExpectation, prefix []string) []detailPath {
	var paths []detailPath
	for _, field := range fields {
		path := append(append([]string(nil), prefix...), field.name)
		paths = append(paths, detailPath{names: path, nullable: field.nullable, kind: field.kind})
		paths = append(paths, detailPaths(field.children, path)...)
	}
	return paths
}

func detailParent(details map[string]any, path []string) map[string]any {
	parent := details
	for _, name := range path[:len(path)-1] {
		switch value := parent[name].(type) {
		case map[string]any:
			parent = value
		case []map[string]any:
			parent = value[0]
		default:
			panic("detail path does not resolve to an object")
		}
	}
	return parent
}

func assertNewFailureRejected(t *testing.T, code Code, subject *Subject, details map[string]any) {
	t.Helper()
	if _, err := NewFailure(code, "message", subject, cloneDetails(details)); err == nil {
		t.Fatalf("NewFailure(%q) error = nil, want rejection for subject %#v details %#v", code, subject, details)
	}
}

func assertSerializedFailureSchema(t *testing.T, failure Failure, expectation reasonCodeExpectation) {
	t.Helper()
	var output bytes.Buffer
	if err := WriteFailure(&output, "operation.name", failure); err != nil {
		t.Fatalf("WriteFailure() error = %v", err)
	}
	if !strings.HasSuffix(output.String(), "\n") || strings.HasSuffix(output.String(), "\n\n") {
		t.Fatalf("WriteFailure() output = %q, want one trailing newline", output.String())
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("json.Unmarshal(envelope) error = %v", err)
	}
	assertExactKeys(t, envelope, "contract_version", "ok", "operation", "error")

	var responseError map[string]json.RawMessage
	if err := json.Unmarshal(envelope["error"], &responseError); err != nil {
		t.Fatalf("json.Unmarshal(error) error = %v", err)
	}
	assertExactKeys(t, responseError, "code", "category", "retryable", "message", "subject", "details")
	assertJSONValue(t, responseError["code"], `"`+string(expectation.code)+`"`)
	assertJSONValue(t, responseError["category"], `"`+string(expectation.category)+`"`)
	assertJSONValue(t, responseError["retryable"], map[bool]string{true: "true", false: "false"}[expectation.retryable])
	if len(expectation.subjects) == 0 {
		assertJSONValue(t, responseError["subject"], `null`)
	} else {
		var subject map[string]json.RawMessage
		if err := json.Unmarshal(responseError["subject"], &subject); err != nil {
			t.Fatalf("json.Unmarshal(subject) error = %v", err)
		}
		assertExactKeys(t, subject, "type", "id")
		assertJSONValue(t, subject["type"], `"`+string(expectation.subjects[0])+`"`)
		assertJSONValue(t, subject["id"], `"subject-id"`)
	}
	assertSerializedDetails(t, responseError["details"], expectation.details)
}

func assertSerializedDetails(t *testing.T, raw json.RawMessage, fields []detailExpectation) {
	t.Helper()
	var values map[string]json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		t.Fatalf("json.Unmarshal(details) error = %v", err)
	}
	assertExactKeys(t, values, detailNames(fields)...)
	for _, field := range fields {
		assertSerializedDetailValue(t, values[field.name], field)
	}
}

func assertSerializedDetailValue(t *testing.T, raw json.RawMessage, field detailExpectation) {
	t.Helper()
	switch field.kind {
	case detailStringKind:
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("json.Unmarshal(%q string) error = %v", field.name, err)
		}
	case detailIntegerKind:
		var value int64
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("json.Unmarshal(%q integer) error = %v", field.name, err)
		}
	case detailStringsKind:
		var value []string
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("json.Unmarshal(%q string array) error = %v", field.name, err)
		}
	case detailBooleanKind:
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("json.Unmarshal(%q boolean) error = %v", field.name, err)
		}
	case detailObjectKind:
		assertSerializedDetails(t, raw, field.children)
	case detailObjectsKind:
		var values []json.RawMessage
		if err := json.Unmarshal(raw, &values); err != nil {
			t.Fatalf("json.Unmarshal(%q object array) error = %v", field.name, err)
		}
		if len(values) != 1 {
			t.Fatalf("%q object array length = %d, want 1", field.name, len(values))
		}
		assertSerializedDetails(t, values[0], field.children)
	default:
		t.Fatalf("unknown detail kind for %q", field.name)
	}
}
