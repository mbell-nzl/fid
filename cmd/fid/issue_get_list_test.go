package main

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
)

var issueSummaryKeys = []string{
	"id", "kind", "title", "status", "status_category", "workflow_id", "parent_id",
	"priority", "sequence", "spec_md_path", "revision", "labels", "lane", "created_at", "updated_at",
}

var issueDetailKeys = append(append([]string{}, issueSummaryKeys...), "acceptance_criteria")

type readEnvelope struct {
	ContractVersion string         `json:"contract_version"`
	OK              bool           `json:"ok"`
	Operation       string         `json:"operation"`
	Result          map[string]any `json:"result"`
	Event           *any           `json:"event"`
}

func decodeReadEnvelope(t *testing.T, output, operation string) readEnvelope {
	t.Helper()
	var envelope readEnvelope
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Unmarshal(envelope) error = %v", err)
	}
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != operation || envelope.Event != nil {
		t.Fatalf("envelope = %#v, want successful %s read with event:null", envelope, operation)
	}
	return envelope
}

func assertJSONFailureCode(t *testing.T, output, operation, code string) {
	t.Helper()
	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Error           struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Unmarshal(failure envelope) error = %v; output = %q", err, output)
	}
	assertEnvelopeKeySet(t, []byte(output), []string{"contract_version", "error", "ok", "operation"})
	assertNestedEnvelopeKeySet(t, []byte(output), "error", []string{"category", "code", "details", "message", "retryable", "subject"})
	if envelope.ContractVersion != "fid-cli.v1" || envelope.OK || envelope.Operation != operation || envelope.Error.Code != code {
		t.Fatalf("failure envelope = %#v, want %s failure code %q", envelope, operation, code)
	}
}

func requireObject(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := object[key]
	if !ok {
		t.Fatalf("object = %#v, missing %q", object, key)
	}
	result, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("object[%q] = %#v, want object", key, value)
	}
	return result
}

func assertExactKeys(t *testing.T, object map[string]any, want []string) {
	t.Helper()
	gotKeys := make([]string, 0, len(object))
	for key := range object {
		gotKeys = append(gotKeys, key)
	}
	sort.Strings(gotKeys)
	wantKeys := append([]string(nil), want...)
	sort.Strings(wantKeys)
	if len(gotKeys) != len(wantKeys) {
		t.Fatalf("keys = %#v, want %#v", gotKeys, wantKeys)
	}
	for index := range wantKeys {
		if gotKeys[index] != wantKeys[index] {
			t.Fatalf("keys = %#v, want %#v", gotKeys, wantKeys)
		}
	}
}

func assertNull(t *testing.T, object map[string]any, key string) {
	t.Helper()
	if value, found := object[key]; !found || value != nil {
		t.Fatalf("object[%q] = %#v, want present null", key, value)
	}
}

func assertArray(t *testing.T, object map[string]any, key string, length int) {
	t.Helper()
	value, found := object[key]
	array, ok := value.([]any)
	if !found || !ok || len(array) != length {
		t.Fatalf("object[%q] = %#v, want array length %d", key, value, length)
	}
}

func TestIssueGetJSONMatchesClosedSchema(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Closed read", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, stderr)

	result := executeCommandForTest(workspace, "--json", "issue", "get", issueID)
	if result.err != nil {
		t.Fatalf("Execute() error = %v, want nil; stderr = %q", result.err, result.stderr)
	}
	if result.stderr != "" {
		t.Fatalf("stderr = %q, want empty", result.stderr)
	}

	envelope := decodeReadEnvelope(t, result.stdout, "issue.get")
	assertExactKeys(t, envelope.Result, []string{"issue"})
	issue := requireObject(t, envelope.Result, "issue")
	assertExactKeys(t, issue, issueDetailKeys)
	assertNull(t, issue, "parent_id")
	assertNull(t, issue, "priority")
	assertNull(t, issue, "sequence")
	assertNull(t, issue, "spec_md_path")
	assertArray(t, issue, "labels", 0)
	assertArray(t, issue, "acceptance_criteria", 0)
}

func TestIssueListJSONMatchesClosedSchemaAndRetainedFilters(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)
	seedIssueListCommandFixture(t, workspace)

	result := executeCommandForTest(workspace, "--json", "issue", "list", "--kind", "task", "--status-category", "started", "--order", "updated")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("list result = %#v, want JSON success on stdout only", result)
	}
	envelope := decodeReadEnvelope(t, result.stdout, "issue.list")
	assertExactKeys(t, envelope.Result, []string{"issues", "count"})
	items, ok := envelope.Result["issues"].([]any)
	if !ok || len(items) != 1 || envelope.Result["count"] != float64(1) {
		t.Fatalf("list result = %#v, want one started IssueSummary", envelope.Result)
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("issue list item = %#v, want object", items[0])
	}
	assertExactKeys(t, item, issueSummaryKeys)
	if item["id"] != "FID-task02" {
		t.Fatalf("issue list id = %#v, want FID-task02", item["id"])
	}
}

func TestIssueReadRejectsRetiredAndMalformedInputs(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	for _, args := range [][]string{
		{"--json", "issue", "get", "missing", "--section", "acceptance"},
		{"--json", "issue", "list", "--inbox"},
		{"--json", "issue", "list", "--kind", "story"},
		{"--json", "issue", "list", "--status-category", "paused"},
		{"--json", "issue", "list", "--priority", "0,0"},
		{"--json", "issue", "list", "--priority", "6"},
		{"--json", "issue", "list", "--order", "title"},
		{"--json", "ready", "--format", "json"},
	} {
		result := executeCommandForTest(workspace, args...)
		if result.err == nil {
			t.Fatalf("%v succeeded, want rejection", args)
		}
		if result.stdout != "" {
			t.Fatalf("%v stdout = %q, want empty on rejection", args, result.stdout)
		}
	}
}

func TestReadCollectionLimitRejectsOnlyOverLimit(t *testing.T) {
	t.Parallel()
	if err := issues.EnforceReadCollectionLimit("issues", 10000); err != nil {
		t.Fatalf("at limit error = %v", err)
	}
	if err := issues.EnforceReadCollectionLimit("issues", 10001); err == nil {
		t.Fatal("over limit error = nil")
	}
}
