//go:build integration

package integration

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestP24RetainedKeyedLeavesRebuildInOneJournal(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	issue := leanFlagFixture(t, binary, workspace)
	p24StartKeyedCapture(workspace)
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	}
	call("label", "add", "parity", "--colour", "red", "--actor", "agent:developer", "--idempotency-key", "single-label")
	call("issue", "transition", issue, "--to", "in_progress", "--expected-revision", "1", "--actor", "agent:developer", "--idempotency-key", "single-transition")
	call("issue", "ac", "add", issue, "--text", "Parity criterion", "--ordinal", "1", "--actor", "agent:developer", "--idempotency-key", "single-ac")
	call("issue", "update", issue, "--title", "Parity update", "--expected-revision", "3", "--actor", "agent:developer", "--idempotency-key", "single-update")
	call("issue", "label", "add", issue, "parity", "--actor", "agent:developer", "--idempotency-key", "single-label-add")
	call("issue", "label", "remove", issue, "parity", "--actor", "agent:developer", "--idempotency-key", "single-label-remove")
	peer := call("issue", "create", "--kind", "task", "--title", "Parity peer", "--actor", "agent:developer", "--idempotency-key", "single-peer")
	peerID := p21PathString(t, peer, "result", "issue", "id")
	call("link", "add", issue, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "single-link-add")
	call("link", "remove", issue, "blocks", peerID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "single-link-remove")
	p24LifecycleCommands(t, binary, workspace)
	captured := p24EndKeyedCapture(workspace)
	want := retainedKeyedLeaves(t)
	got := p24CapturedNames(captured)
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("retained keyed inventory=%v, want %v", got, want)
	}
	for _, operation := range want {
		field := p24FlagFields[operation]
		if field == "" {
			field = p24LifecycleFields[operation]
		}
		if operation == "issue.create" {
			field = "--title"
		}
		p24ExerciseCapturedLeaf(t, workspace, binary, operation, field, captured[operation])
	}
	conn := p22FaultDB(t, workspace)
	rows, err := conn.Query("SELECT DISTINCT operation_scope FROM idempotency_results ORDER BY operation_scope")
	if err != nil {
		t.Fatal(err)
	}
	scopes := []string{}
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			t.Fatal(err)
		}
		operation, _, _ := strings.Cut(scope, ":")
		scopes = append(scopes, operation)
	}
	p22FaultOK(t, rows.Err())
	p22FaultOK(t, rows.Close())
	p22FaultOK(t, conn.Close())
	sort.Strings(scopes)
	if !reflect.DeepEqual(scopes, want) {
		t.Fatalf("stored retained scopes=%v, want %v", scopes, want)
	}
	assertLeanIntegrationJournal(t, workspace)
	p24AssertRebuildMapParity(t, workspace, binary)
}

func retainedKeyedLeaves(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("testdata/p24_keyed_leaves.tsv")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{}
	for _, row := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		name, _, ok := strings.Cut(row, "\t")
		if !ok {
			t.Fatalf("invalid keyed row %q", row)
		}
		family, _, _ := strings.Cut(name, " ")
		if family == "artifact" {
			continue
		} // Own command corpus moves with its adapter in 1.5c2.
		switch family {
		case "issue", "label", "link", "park", "answer", "takeover", "start", "check", "handoff":
		default:
			continue
		}
		operation := strings.ReplaceAll(name, " ", ".")
		if mapped := map[string]string{"check": "check.execute", "handoff": "task.handoff", "start": "task.start", "park": "hold.park", "answer": "task.answer", "takeover": "task.takeover"}[name]; mapped != "" {
			operation = mapped
		}
		want = append(want, operation)
	}
	sort.Strings(want)
	if len(want) != 15 {
		t.Fatalf("retained keyed leaves=%v, want 15 plus separate Artifact compatibility corpus", want)
	}
	return want
}
