package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestEpicCancellationRecordsPersonsDecision(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"", "cancel-epic"} {
		t.Run("key="+key, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			args := []string{"--json", "issue", "transition", epic, "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"}
			if key != "" {
				args = append(args, "--idempotency-key", key)
			}
			result, event := lifecycleResult(t, f.run(args...), "issue.transition")
			id, ok := result["decision_id"].(string)
			if !ok || !strings.HasPrefix(id, "PDN-") || len(id) != 10 {
				t.Fatalf("decision ID = %v, want PDN- and six hex digits", result["decision_id"])
			}
			wantIssue := map[string]any{"id": epic, "revision": float64(2)}
			if result["from_state"] != "todo" || result["to_state"] != "cancelled" || result["idempotency_key"] != key || result["replayed"] != false || !reflect.DeepEqual(result["issue"], wantIssue) {
				t.Fatalf("cancel result = %v", result)
			}
			at, err := time.Parse(time.RFC3339Nano, event.(map[string]any)["occurred_at"].(string))
			if err != nil {
				t.Fatal(err)
			}
			want := map[string]any{"id": id, "kind": "cancel", "actor": "human:mark", "subject_type": "epic",
				"subject_id": epic, "reason": "No longer needed", "binding_id": epic, "decided_at": at.UTC().Format(time.RFC3339)}
			assertEpicCancellationDecision(t, f, want)
			assertEpicCancellationEvent(t, f, event, want)
			before := snapshotEpicCancellationPersistence(t, f.workspace)
			assertEpicCancellationRetry(t, f, args, key, result, event)
			if after := snapshotEpicCancellationPersistence(t, f.workspace); !before.equal(after) {
				t.Fatal("retry changed the projection or journal")
			}
			if rebuilt := f.run("rebuild"); rebuilt.err != nil {
				t.Fatalf("rebuild = %#v", rebuilt)
			}
			assertEpicCancellationDecision(t, f, want)
			shown, _ := lifecycleResult(t, f.run("--json", "show", epic), "issue.show")
			issue := shown["issue"].(map[string]any)
			if issue["status"] != "cancelled" || issue["revision"] != float64(2) {
				t.Fatalf("rebuilt epic = %v", issue)
			}
			before = snapshotEpicCancellationPersistence(t, f.workspace)
			assertEpicCancellationRetry(t, f, args, key, result, event)
			assertEpicCancellationDecision(t, f, want)
			if after := snapshotEpicCancellationPersistence(t, f.workspace); !before.equal(after) {
				t.Fatal("retry after rebuild changed the projection or journal")
			}
		})
	}
}

func TestEpicCancellationRefusalsLeaveLedgerUnchanged(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, actor, reason, key, code string
		exit                           int
	}{
		{"agent keyed", "agent:builder", "No longer needed", "cancel", "authority-denied", 5},
		{"agent unkeyed", "agent:builder", "No longer needed", "", "authority-denied", 5},
		{"service keyed", "service:ci", "No longer needed", "cancel", "authority-denied", 5},
		{"service unkeyed", "service:ci", "No longer needed", "", "authority-denied", 5},
		{"reason missing keyed", "human:mark", "", "cancel", "invalid-argument", 2},
		{"reason missing unkeyed", "human:mark", "", "", "invalid-argument", 2},
		{"reason blank", "human:mark", "   ", "cancel", "invalid-argument", 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			before := snapshotEpicCancellationPersistence(t, f.workspace)
			args := []string{"--json", "issue", "transition", epic, "--to", "cancelled", "--actor", test.actor}
			if test.reason != "" {
				args = append(args, "--reason", test.reason)
			}
			if test.key != "" {
				args = append(args, "--idempotency-key", test.key)
			}
			answer := f.run(args...)
			if cli.ExitCode(answer.err) != test.exit || answer.stdout != "" {
				t.Fatalf("cancel refusal = %#v, want exit %d", answer, test.exit)
			}
			assertJSONFailureCode(t, answer.stderr, "issue.transition", test.code)
			if test.code == "invalid-argument" && !strings.Contains(answer.stderr, `"argument":"--reason"`) {
				t.Fatalf("missing reason refusal = %s", answer.stderr)
			}
			if after := snapshotEpicCancellationPersistence(t, f.workspace); !before.equal(after) {
				t.Fatal("refused epic cancellation changed the projection or journal")
			}
		})
	}
}

func newEpicCancellationFixture(t *testing.T) (*statusFixture, string) {
	t.Helper()
	f := newLifecycleFixture(t)
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--title", "Epic to cancel",
		"--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "create-epic"), "issue.create")
	return f, created["issue"].(map[string]any)["id"].(string)
}

func assertEpicCancellationDecision(t *testing.T, f *statusFixture, want map[string]any) {
	t.Helper()
	answer := f.run("sql", "SELECT * FROM person_decisions ORDER BY id")
	var rows []map[string]any
	if answer.err != nil || answer.stderr != "" || json.Unmarshal([]byte(answer.stdout), &rows) != nil || !reflect.DeepEqual(rows, []map[string]any{want}) {
		t.Fatalf("complete epic cancellation decision = %#v, want %v", answer, want)
	}
}

func assertEpicCancellationEvent(t *testing.T, f *statusFixture, response any, want map[string]any) {
	t.Helper()
	event := lastJournalEvent(t, f.workspace)
	if event["event_type"] != "issue.transitioned" || event["event_id"] != response.(map[string]any)["id"] || event["aggregate"].(map[string]any)["id"] != want["subject_id"] {
		t.Fatalf("cancel Event = %v", event)
	}
	decisions := 0
	for _, raw := range event["effects"].([]any) {
		effect := raw.(map[string]any)
		if effect["table"] == "person_decisions" {
			decisions++
			row := effect["values"].(map[string]any)
			row["id"] = effect["pk"].(map[string]any)["id"]
			if effect["op"] != "insert" || !reflect.DeepEqual(row, want) {
				t.Fatalf("cancel decision Effect = %v, want %v", effect, want)
			}
		}
	}
	if decisions != 1 {
		t.Fatalf("cancellation Event has %d decision Effects, want 1", decisions)
	}
}

func assertEpicCancellationRetry(t *testing.T, f *statusFixture, args []string, key string, original map[string]any, originalEvent any) {
	t.Helper()
	answer := f.run(args...)
	if key == "" {
		if cli.ExitCode(answer.err) != 4 || answer.stdout != "" {
			t.Fatalf("unkeyed repeat = %#v, want state refusal", answer)
		}
		assertJSONFailureCode(t, answer.stderr, "issue.transition", "transition-not-permitted")
		return
	}
	retried, event := lifecycleResult(t, answer, "issue.transition")
	want := make(map[string]any, len(original))
	for name, value := range original {
		want[name] = value
	}
	want["replayed"] = true
	if !reflect.DeepEqual(retried, want) || !reflect.DeepEqual(event, originalEvent) {
		t.Fatalf("keyed retry = %v, %v; want %v, %v", retried, event, want, originalEvent)
	}
}

func snapshotEpicCancellationPersistence(t *testing.T, workspace string) workspacePersistenceSnapshot {
	t.Helper()
	return snapshotWorkspacePersistence(t, workspace)
}
