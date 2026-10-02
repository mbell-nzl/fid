package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteSuccessUsesExactContractKeys(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer

	err := WriteSuccess(&output, "issue.get", map[string]any{"issue": map[string]any{}}, nil)
	if err != nil {
		t.Fatalf("WriteSuccess() error = %v", err)
	}

	if got, want := output.String()[len(output.String())-1:], "\n"; got != want {
		t.Fatalf("WriteSuccess() trailing character = %q, want %q", got, want)
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	assertExactKeys(t, envelope, "contract_version", "ok", "operation", "result", "event")
	assertJSONValue(t, envelope["contract_version"], `"fid-cli.v1"`)
	assertJSONValue(t, envelope["ok"], `true`)
	assertJSONValue(t, envelope["operation"], `"issue.get"`)
	assertJSONValue(t, envelope["event"], `null`)
}

func TestWriteFailureUsesExactContractKeys(t *testing.T) {
	t.Parallel()
	failure, err := NewFailure(
		CodeIssueNotFound,
		"Issue FID-123 was not found",
		&Subject{Type: SubjectIssue, ID: "FID-123"},
		map[string]any{"id": "FID-123"},
	)
	if err != nil {
		t.Fatalf("NewFailure() error = %v", err)
	}

	var output bytes.Buffer
	if err := WriteFailure(&output, "issue.get", failure); err != nil {
		t.Fatalf("WriteFailure() error = %v", err)
	}

	if got, want := output.String()[len(output.String())-1:], "\n"; got != want {
		t.Fatalf("WriteFailure() trailing character = %q, want %q", got, want)
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	assertExactKeys(t, envelope, "contract_version", "ok", "operation", "error")
	assertJSONValue(t, envelope["contract_version"], `"fid-cli.v1"`)
	assertJSONValue(t, envelope["ok"], `false`)
	assertJSONValue(t, envelope["operation"], `"issue.get"`)

	var responseError map[string]json.RawMessage
	if err := json.Unmarshal(envelope["error"], &responseError); err != nil {
		t.Fatalf("json.Unmarshal(error) error = %v", err)
	}
	assertExactKeys(t, responseError, "code", "category", "retryable", "message", "subject", "details")
	assertJSONValue(t, responseError["code"], `"issue-not-found"`)
	assertJSONValue(t, responseError["category"], `"not-found"`)
	assertJSONValue(t, responseError["retryable"], `false`)
	assertJSONValue(t, responseError["subject"], `{"id":"FID-123","type":"issue"}`)
	assertJSONValue(t, responseError["details"], `{"id":"FID-123"}`)
}

func TestWriteFailurePreservesRequiredNullFields(t *testing.T) {
	t.Parallel()
	failure, err := NewFailure(
		CodeInvalidArgument,
		"prefix is invalid",
		nil,
		map[string]any{"argument": "prefix", "value": nil, "allowed": []string{"FID"}},
	)
	if err != nil {
		t.Fatalf("NewFailure() error = %v", err)
	}

	var output bytes.Buffer
	if err := WriteFailure(&output, "project.init", failure); err != nil {
		t.Fatalf("WriteFailure() error = %v", err)
	}

	var envelope struct {
		Error struct {
			Subject *Subject       `json:"subject"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if envelope.Error.Subject != nil {
		t.Fatalf("subject = %#v, want nil", envelope.Error.Subject)
	}
	if value, ok := envelope.Error.Details["value"]; !ok || value != nil {
		t.Fatalf("details value = %#v, want present nil", value)
	}
}

func assertExactKeys(t *testing.T, values map[string]json.RawMessage, want ...string) {
	t.Helper()
	if len(values) != len(want) {
		t.Fatalf("key count = %d, want %d; keys = %v", len(values), len(want), values)
	}
	for _, key := range want {
		if _, ok := values[key]; !ok {
			t.Fatalf("missing key %q in %v", key, values)
		}
	}
}

func assertJSONValue(t *testing.T, got json.RawMessage, want string) {
	t.Helper()
	if string(got) != want {
		t.Fatalf("JSON value = %s, want %s", got, want)
	}
}
