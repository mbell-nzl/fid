package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTask10JSONRequestErrorsUseCanonicalEnvelope(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	journalPath := mustJournalPath(t, workspace)

	tests := []struct {
		name      string
		operation string
		args      []string
	}{
		{"issue ac add arity", "issue.ac.add", []string{"--json", "issue", "ac", "add"}},
		{"issue ac add unknown flag", "issue.ac.add", []string{"--json", "issue", "ac", "add", "FID-any", "--text", "criterion", "--ordinal", "1", "--by", "developer", "--unknown", "value"}},
		{"issue ac add missing text", "issue.ac.add", []string{"--json", "issue", "ac", "add", "FID-any", "--ordinal", "1", "--by", "developer"}},
		{"issue ac add missing ordinal", "issue.ac.add", []string{"--json", "issue", "ac", "add", "FID-any", "--text", "criterion", "--by", "developer"}},
		{"issue ac list arity", "issue.ac.list", []string{"--json", "issue", "ac", "list"}},
		{"issue ac list unknown flag", "issue.ac.list", []string{"--json", "issue", "ac", "list", "FID-any", "--unknown", "value"}},
		{"label add arity", "label.add", []string{"--json", "label", "add"}},
		{"label add unknown flag", "label.add", []string{"--json", "label", "add", "core", "--by", "developer", "--unknown", "value"}},
		{"label list arity", "label.list", []string{"--json", "label", "list", "extra"}},
		{"label list unknown flag", "label.list", []string{"--json", "label", "list", "--unknown", "value"}},
		{"issue label add arity", "issue.label.add", []string{"--json", "issue", "label", "add"}},
		{"issue label add unknown flag", "issue.label.add", []string{"--json", "issue", "label", "add", "FID-any", "core", "--by", "developer", "--unknown", "value"}},
		{"issue label remove arity", "issue.label.remove", []string{"--json", "issue", "label", "remove"}},
		{"issue label remove unknown flag", "issue.label.remove", []string{"--json", "issue", "label", "remove", "FID-any", "core", "--by", "developer", "--unknown", "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, tt.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only request rejection", result)
			}

			var failure struct {
				ContractVersion string `json:"contract_version"`
				OK              bool   `json:"ok"`
				Operation       string `json:"operation"`
				Error           struct {
					Code      string `json:"code"`
					Category  string `json:"category"`
					Retryable bool   `json:"retryable"`
					Message   string `json:"message"`
					Subject   any    `json:"subject"`
					Details   struct {
						Argument string   `json:"argument"`
						Value    any      `json:"value"`
						Allowed  []string `json:"allowed"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(stderr) error = %v; stderr = %q", err, result.stderr)
			}
			assertEnvelopeKeySet(t, []byte(result.stderr), []string{"contract_version", "error", "ok", "operation"})
			assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error", []string{"category", "code", "details", "message", "retryable", "subject"})
			assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error.details", []string{"allowed", "argument", "value"})
			if failure.ContractVersion != "fid-cli.v1" || failure.OK || failure.Operation != tt.operation || failure.Error.Code != "invalid-argument" || failure.Error.Category != "invalid-input" || failure.Error.Retryable || failure.Error.Message == "" || failure.Error.Subject != nil {
				t.Fatalf("failure = %+v, want canonical %s invalid-argument", failure, tt.operation)
			}
			if failure.Error.Details.Argument != "request" || failure.Error.Details.Value != nil || strings.Join(failure.Error.Details.Allowed, ",") != "accepted "+tt.operation+" input" {
				t.Fatalf("details = %+v, want request/null/accepted operation input", failure.Error.Details)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestTask10JSONMissingActorsUseAuthorityEnvelope(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	journalPath := mustJournalPath(t, workspace)

	for _, tt := range []struct {
		name      string
		operation string
		args      []string
	}{
		{"issue ac add missing by", "issue.ac.add", []string{"--json", "issue", "ac", "add", "FID-any", "--text", "criterion", "--ordinal", "1"}},
		{"issue ac add blank by", "issue.ac.add", []string{"--json", "issue", "ac", "add", "FID-any", "--text", "criterion", "--ordinal", "1", "--by", " "}},
		{"label add missing by", "label.add", []string{"--json", "label", "add", "core"}},
		{"label add blank by", "label.add", []string{"--json", "label", "add", "core", "--by", " "}},
		{"issue label add missing by", "issue.label.add", []string{"--json", "issue", "label", "add", "FID-any", "core"}},
		{"issue label add blank by", "issue.label.add", []string{"--json", "issue", "label", "add", "FID-any", "core", "--by", " "}},
		{"issue label remove missing by", "issue.label.remove", []string{"--json", "issue", "label", "remove", "FID-any", "core"}},
		{"issue label remove blank by", "issue.label.remove", []string{"--json", "issue", "label", "remove", "FID-any", "core", "--by", " "}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, tt.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only authority rejection", result)
			}

			var failure struct {
				ContractVersion string `json:"contract_version"`
				OK              bool   `json:"ok"`
				Operation       string `json:"operation"`
				Error           struct {
					Code      string `json:"code"`
					Category  string `json:"category"`
					Retryable bool   `json:"retryable"`
					Message   string `json:"message"`
					Subject   any    `json:"subject"`
					Details   struct {
						Operation string `json:"operation"`
						Actor     any    `json:"actor"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
				t.Fatalf("Unmarshal(stderr) error = %v; stderr = %q", err, result.stderr)
			}
			assertEnvelopeKeySet(t, []byte(result.stderr), []string{"contract_version", "error", "ok", "operation"})
			assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error", []string{"category", "code", "details", "message", "retryable", "subject"})
			assertNestedEnvelopeKeySet(t, []byte(result.stderr), "error.details", []string{"actor", "operation"})
			if failure.ContractVersion != "fid-cli.v1" || failure.OK || failure.Operation != tt.operation || failure.Error.Code != "authority-missing" || failure.Error.Category != "authority" || failure.Error.Retryable || failure.Error.Message != "acting authority is required" || failure.Error.Subject != nil {
				t.Fatalf("failure = %+v, want canonical %s authority-missing", failure, tt.operation)
			}
			if failure.Error.Details.Operation != tt.operation || failure.Error.Details.Actor != nil {
				t.Fatalf("details = %+v, want operation and null actor", failure.Error.Details)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestTask10RequiredActorHumanErrorsRemainProse(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	journalPath := mustJournalPath(t, workspace)

	for _, tt := range []struct {
		name string
		args []string
	}{
		{"issue ac add missing by", []string{"issue", "ac", "add", "FID-any", "--text", "criterion", "--ordinal", "1"}},
		{"issue ac add blank by", []string{"issue", "ac", "add", "FID-any", "--text", "criterion", "--ordinal", "1", "--by", " "}},
		{"label add missing by", []string{"label", "add", "core"}},
		{"label add blank by", []string{"label", "add", "core", "--by", " "}},
		{"issue label add missing by", []string{"issue", "label", "add", "FID-any", "core"}},
		{"issue label add blank by", []string{"issue", "label", "add", "FID-any", "core", "--by", " "}},
		{"issue label remove missing by", []string{"issue", "label", "remove", "FID-any", "core"}},
		{"issue label remove blank by", []string{"issue", "label", "remove", "FID-any", "core", "--by", " "}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, tt.args...)
			if result.err == nil || result.stdout != "" || result.stderr != "error: --by is required\n" {
				t.Fatalf("result = %#v, want concise Human required-actor error", result)
			}
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestTask10OperationsUseSelectedProjectOutsideWorkingDirectory(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	outside := t.TempDir()

	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "label", "add", "core", "--by", "developer"), "label.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "label", "list"), "label.list")

	created := executeCommandForTest(outside, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Selected project", "--actor", "agent:developer")
	if created.err != nil || created.stderr != "" {
		t.Fatalf("issue create = %#v, want selected-project success", created)
	}
	var createEnvelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil || createEnvelope.Result.Issue.ID == "" {
		t.Fatalf("issue create envelope = %#v, error = %v", createEnvelope, err)
	}
	issueID := createEnvelope.Result.Issue.ID

	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "ac", "add", issueID, "--text", "Selected project AC", "--ordinal", "1", "--by", "developer"), "issue.ac.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "ac", "list", issueID), "issue.ac.list")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "label", "add", issueID, "core", "--by", "developer"), "issue.label.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "label", "remove", issueID, "core", "--by", "developer"), "issue.label.remove")
}

func TestIssueACVerifyRemainsGuardedLegacyLeaf(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Legacy verify guard")
	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace, "--json", "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "pass")
	if result.err == nil || result.stdout != "" || result.stderr == "" {
		t.Fatalf("result = %#v, want guarded legacy rejection", result)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after guarded rejection = %d, want %d", got, before)
	}
}

func assertTask10ProjectOperation(t *testing.T, result commandExecutionResult, operation string) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("%s result = %#v, want selected-project JSON success", operation, result)
	}
	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", operation, err)
	}
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != operation {
		t.Fatalf("%s envelope = %#v, want selected-project success", operation, envelope)
	}
}
