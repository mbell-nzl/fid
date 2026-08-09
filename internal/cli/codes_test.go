package cli

import (
	"reflect"
	"testing"
)

func TestReasonCodeRegistryIsClosedAndComplete(t *testing.T) {
	expected := map[Code]struct {
		category  Category
		retryable bool
		subjects  []SubjectType
		details   []string
	}{
		CodeInvalidArgument:        {CategoryInvalidInput, false, nil, []string{"argument", "value", "allowed"}},
		CodeUnsupportedKind:        {CategoryInvalidInput, false, nil, []string{"argument", "value", "allowed"}},
		CodeUnsupportedRelation:    {CategoryInvalidInput, false, nil, []string{"argument", "value", "allowed"}},
		CodeProjectNotFound:        {CategoryNotFound, false, []SubjectType{SubjectProject}, []string{"path"}},
		CodeProjectAlreadyInit:     {CategoryConflict, false, []SubjectType{SubjectProject}, []string{"path"}},
		CodeIssueNotFound:          {CategoryNotFound, false, []SubjectType{SubjectIssue}, []string{"id"}},
		CodeRequirementNotFound:    {CategoryNotFound, false, []SubjectType{SubjectRequirement}, []string{"id"}},
		CodeLabelNotFound:          {CategoryNotFound, false, []SubjectType{SubjectLabel}, []string{"id"}},
		CodeLinkNotFound:           {CategoryNotFound, false, []SubjectType{SubjectLink}, []string{"from_type", "from_id", "relation", "to_type", "to_id"}},
		CodeWorkflowNotFound:       {CategoryNotFound, false, []SubjectType{SubjectWorkflow}, []string{"id"}},
		CodeDuplicateEntity:        {CategoryConflict, false, []SubjectType{SubjectIssue, SubjectLabel, SubjectLink}, []string{"entity_type", "identity"}},
		CodeStateConflict:          {CategoryConflict, false, []SubjectType{SubjectIssue}, []string{"expected_revision", "actual_revision"}},
		CodeLockUnavailable:        {CategoryConflict, true, nil, []string{"lock_path"}},
		CodeHierarchyInvalid:       {CategoryPrecondition, false, []SubjectType{SubjectIssue}, []string{"kind", "parent_id", "parent_kind", "rule"}},
		CodeDependencyCycle:        {CategoryPrecondition, false, nil, []string{"path"}},
		CodeTransitionNotPermitted: {CategoryPrecondition, false, []SubjectType{SubjectIssue}, []string{"from_state", "to_state", "unsatisfied"}},
		CodeAuthorityMissing:       {CategoryAuthority, false, nil, []string{"operation", "actor"}},
		CodeAuthorityUnknownActor:  {CategoryAuthority, false, nil, []string{"operation", "actor"}},
		CodeAuthorityDenied:        {CategoryAuthority, false, nil, []string{"operation", "actor"}},
		CodeJournalReplayStopped:   {CategoryIntegrity, false, nil, []string{"journal_file", "line", "event_id"}},
		CodeProjectionDivergence:   {CategoryIntegrity, false, nil, []string{"status", "errors", "warnings", "info"}},
		CodeResultLimitExceeded:    {CategoryIntegrity, false, nil, []string{"collection", "limit", "observed"}},
		CodeInternalError:          {CategoryInternal, false, nil, []string{}},
	}

	registered := RegisteredCodes()
	if len(registered) != len(expected) {
		t.Fatalf("RegisteredCodes() length = %d, want %d", len(registered), len(expected))
	}

	for code, want := range expected {
		spec, ok := LookupCode(code)
		if !ok {
			t.Fatalf("LookupCode(%q) found no specification", code)
		}
		if spec.Category != want.category || spec.Retryable != want.retryable {
			t.Fatalf("LookupCode(%q) = category %q retryable %t, want category %q retryable %t", code, spec.Category, spec.Retryable, want.category, want.retryable)
		}
		if !reflect.DeepEqual(spec.SubjectTypes, want.subjects) {
			t.Fatalf("LookupCode(%q) subject types = %v, want %v", code, spec.SubjectTypes, want.subjects)
		}
		if !reflect.DeepEqual(spec.DetailFields, want.details) {
			t.Fatalf("LookupCode(%q) detail fields = %v, want %v", code, spec.DetailFields, want.details)
		}
		if spec.RedactionPolicy != RedactClosedDetails {
			t.Fatalf("LookupCode(%q) redaction policy = %q, want %q", code, spec.RedactionPolicy, RedactClosedDetails)
		}
	}

	if _, ok := LookupCode("not-a-contract-code"); ok {
		t.Fatal("LookupCode() accepted an unknown code")
	}
}

func TestNewFailureRejectsInvalidSubjectDetailsAndRedaction(t *testing.T) {
	tests := []struct {
		name    string
		code    Code
		subject *Subject
		details map[string]any
	}{
		{
			name:    "unknown code",
			code:    "not-a-contract-code",
			details: map[string]any{},
		},
		{
			name:    "wrong subject type",
			code:    CodeIssueNotFound,
			subject: &Subject{Type: SubjectLabel, ID: "bug"},
			details: map[string]any{"id": "bug"},
		},
		{
			name:    "missing required detail",
			code:    CodeStateConflict,
			subject: &Subject{Type: SubjectIssue, ID: "FID-123"},
			details: map[string]any{"expected_revision": 3},
		},
		{
			name:    "redacted extra detail",
			code:    CodeProjectNotFound,
			subject: &Subject{Type: SubjectProject, ID: "project"},
			details: map[string]any{"path": "/work/project", "credential": "secret"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewFailure(test.code, "message", test.subject, test.details); err == nil {
				t.Fatal("NewFailure() error = nil, want validation error")
			}
		})
	}
}

func TestNewFailurePreservesRequiredNullDetailAndSubjectFields(t *testing.T) {
	failure, err := NewFailure(
		CodeHierarchyInvalid,
		"A task cannot have this parent",
		&Subject{Type: SubjectIssue, ID: "FID-123"},
		map[string]any{
			"kind":        "task",
			"parent_id":   nil,
			"parent_kind": nil,
			"rule":        "task-parent",
		},
	)
	if err != nil {
		t.Fatalf("NewFailure() error = %v", err)
	}
	if failure.Subject == nil || failure.Details["parent_id"] != nil || failure.Details["parent_kind"] != nil {
		t.Fatalf("NewFailure() = %#v, want required null detail and subject fields", failure)
	}
}
