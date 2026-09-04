package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestReceiptReadsAreExactFilteredOrderedAndReadOnly(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"missing.txt"}
	})
	first, firstRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "receipt-read-first")
	if firstRun.err != nil || first.Data.ReceiptID == "" {
		t.Fatalf("first gate run = %#v, terminal=%#v", firstRun, first)
	}
	second, secondRun := executeGateAgainstExistingRun(t, workspace, "receipt-read-first.json", "receipt-read-second")
	if secondRun.err != nil || second.Data.ReceiptID == "" {
		t.Fatalf("second gate run = %#v, terminal=%#v", secondRun, second)
	}

	beforeEvents := journalEventCount(t, workspace)
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "receipt", "get", first.Data.ReceiptID)
	if got.err != nil || got.stderr != "" {
		t.Fatalf("receipt get = %#v", got)
	}
	receipt := decodeCanonicalResult(t, got.stdout)["receipt"].(map[string]any)
	assertReceiptRejection(t, receipt, "missing.txt", "missing")
	runID := receipt["run_id"].(string)
	obligationID := receipt["obligation_id"].(string)

	listed := executeCommandForTest(workspace, "--project", workspace, "--json", "receipt", "list", "--run", runID, "--obligation", obligationID, "--outcome", "succeeded", "--limit", "2")
	if listed.err != nil || listed.stderr != "" {
		t.Fatalf("receipt list = %#v", listed)
	}
	items := decodeCanonicalResult(t, listed.stdout)["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("receipt list items = %#v, want two", items)
	}
	if !receiptOrderLessOrEqual(items[0].(map[string]any), items[1].(map[string]any)) {
		t.Fatalf("receipt list order = %#v, want started_at then id", items)
	}
	for _, item := range items {
		assertReceiptRejection(t, item.(map[string]any), "missing.txt", "missing")
	}
	assertReceiptListLength(t, workspace, 1, "--run", runID, "--limit", "1")
	assertReceiptListLength(t, workspace, 0, "--run", "RUN-missing")
	assertReceiptListLength(t, workspace, 0, "--obligation", "OBL-missing")
	assertReceiptListLength(t, workspace, 0, "--outcome", "failed")
	if afterEvents := journalEventCount(t, workspace); afterEvents != beforeEvents {
		t.Fatalf("Receipt reads changed Event count to %d, want %d", afterEvents, beforeEvents)
	}

	rebuild := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild")
	if rebuild.err != nil || rebuild.stderr != "" {
		t.Fatalf("rebuild = %#v", rebuild)
	}
	rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "receipt", "get", first.Data.ReceiptID)
	if rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuilt receipt get = %#v", rebuilt)
	}
	if got.stdout != rebuilt.stdout {
		t.Fatalf("rebuilt Receipt = %s, want byte-identical %s", rebuilt.stdout, got.stdout)
	}
}

func TestReceiptReadFailuresAreTypedAndClosed(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	notFound := executeCommandForTest(workspace, "--project", workspace, "--json", "receipt", "get", "RCP-missing")
	if notFound.err == nil || notFound.stdout != "" || !strings.Contains(notFound.stderr, `"code":"receipt-not-found"`) || !strings.Contains(notFound.stderr, `"category":"not-found"`) {
		t.Fatalf("receipt get missing = %#v, want typed not-found", notFound)
	}
	for _, args := range [][]string{
		{"receipt", "list", "--outcome", "unknown"},
		{"receipt", "list", "--limit", "0"},
		{"receipt", "list", "--limit", "10001"},
		{"receipt", "get"},
		{"receipt", "get", "RCP-one", "RCP-two"},
	} {
		result := executeCommandForTest(workspace, append([]string{"--project", workspace, "--json"}, args...)...)
		if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
			t.Fatalf("fid %s = %#v, want typed invalid-argument", strings.Join(args, " "), result)
		}
	}
}

func TestGateRunReplayRetainsRejectedArtifacts(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"missing.txt"}
	})
	terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "receipt-replay")
	if result.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v", result, terminal)
	}
	replayed, replayResult := executeGateAgainstExistingRun(t, workspace, "receipt-replay.json", "receipt-replay")
	if replayResult.err != nil || !replayed.Data.Replayed {
		t.Fatalf("gate replay = %#v, terminal=%#v", replayResult, replayed)
	}
	receipt := replayed.Data.Response["result"].(map[string]any)["receipt"].(map[string]any)
	assertReceiptRejection(t, receipt, "missing.txt", "missing")
}

func executeGateAgainstExistingRun(t *testing.T, workspace, inputName, key string) (gateTerminalFrame, commandExecutionResult) {
	t.Helper()
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", filepath.Join(workspace, inputName), "--idempotency-key", key)
	lines := strings.Split(strings.TrimSpace(result.stdout), "\n")
	var terminal gateTerminalFrame
	if len(lines) > 0 && lines[0] != "" {
		if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil {
			t.Fatalf("decode gate terminal: %v", err)
		}
	}
	return terminal, result
}

func assertReceiptRejection(t *testing.T, receipt map[string]any, path, reason string) {
	t.Helper()
	if got := receipt["rejected_artifacts"]; !reflect.DeepEqual(got, []any{map[string]any{"path": path, "reason": reason}}) {
		t.Fatalf("rejected_artifacts = %#v", got)
	}
}

func assertReceiptListLength(t *testing.T, workspace string, want int, args ...string) {
	t.Helper()
	command := append([]string{"--project", workspace, "--json", "receipt", "list"}, args...)
	result := executeCommandForTest(workspace, command...)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("receipt list %v = %#v", args, result)
	}
	if items := decodeCanonicalResult(t, result.stdout)["items"].([]any); len(items) != want {
		t.Fatalf("receipt list %v length = %d, want %d", args, len(items), want)
	}
}

func receiptOrderLessOrEqual(left, right map[string]any) bool {
	leftStarted, rightStarted := left["started_at"].(string), right["started_at"].(string)
	return leftStarted < rightStarted || leftStarted == rightStarted && left["id"].(string) <= right["id"].(string)
}
