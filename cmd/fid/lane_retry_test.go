package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func laneRetryFixture(t *testing.T) (*statusFixture, []string, commandExecutionResult) {
	t.Helper()
	f := newJournalStatusFixture(t)
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Plan lane", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "plan"), "issue.create")
	id := created["issue"].(map[string]any)["id"].(string)
	args := []string{"issue", "update", id, "--lane", "standard", "--actor", "agent:builder", "--idempotency-key", "raise"}
	first := f.run(append([]string{"--json"}, args...)...)
	lifecycleResult(t, first, "issue.update")
	return f, args, first
}

func TestLaneUpdateRetryRejectsMalformedResult(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name  string
		value any
		from  bool
	}{
		{"null", nil, false}, {"number", 42, false}, {"empty", "", false}, {"unknown", "urgent", false}, {"missing", nil, false},
		{"from-number", 42, true}, {"from-empty", "", true}, {"from-unknown", "urgent", true},
	} {
		for _, mode := range []string{"text", "json"} {
			t.Run(row.name+"/"+mode, func(t *testing.T) {
				t.Parallel()
				defer func() {
					if caught := recover(); caught != nil {
						t.Errorf("malformed lane retry panicked: %v", caught)
					}
				}()
				f, args, _ := laneRetryFixture(t)
				conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
				if err != nil {
					t.Fatal(err)
				}
				defer closeCommandDB(t, conn)
				var raw string
				if err := conn.QueryRow(`SELECT result_json FROM idempotency_results WHERE idempotency_key='raise'`).Scan(&raw); err != nil {
					t.Fatal(err)
				}
				var result map[string]any
				if err := json.Unmarshal([]byte(raw), &result); err != nil {
					t.Fatal(err)
				}
				change := result["changes"].([]any)[0].(map[string]any)
				field := "To"
				if row.from {
					field = "From"
				}
				change[field] = row.value
				if row.name == "missing" {
					delete(change, "To")
				}
				payload, err := json.Marshal(result)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := conn.Exec(`UPDATE idempotency_results SET result_json=? WHERE idempotency_key='raise'`, string(payload)); err != nil {
					t.Fatal(err)
				}
				before := f.run("--json", "issue", "get", args[2])
				events := startEventCount(t, f.workspace)
				if mode == "json" {
					args = append([]string{"--json"}, args...)
				}
				got := f.run(args...)
				if got.err == nil || got.stdout != "" || !strings.Contains(got.stderr, "invalid lane change result") {
					t.Fatalf("malformed lane retry = %#v", got)
				}
				after := f.run("--json", "issue", "get", result["issue"].(map[string]any)["id"].(string))
				if before.err != nil || after.err != nil || before.stdout != after.stdout || events != startEventCount(t, f.workspace) {
					t.Fatalf("refused retry changed issue or journal: before=%#v after=%#v", before, after)
				}
			})
		}
	}
}

func TestLaneUpdateRetryPreservesTextAndJSONAcrossRebuild(t *testing.T) {
	t.Parallel()
	f, args, first := laneRetryFixture(t)
	original, event := lifecycleResult(t, first, "issue.update")
	events := startEventCount(t, f.workspace)
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if result := f.run("rebuild"); result.err != nil {
				t.Fatalf("rebuild = %#v", result)
			}
		}
		text := f.run(args...)
		if want := "updated issue " + args[2] + " · Standard\n"; text.err != nil || text.stderr != "" || text.stdout != want {
			t.Fatalf("lane text retry = %#v, want %q", text, want)
		}
		retried, retryEvent := lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "issue.update")
		original["replayed"] = true
		if !reflect.DeepEqual(retried, original) || !reflect.DeepEqual(retryEvent, event) || events != startEventCount(t, f.workspace) {
			t.Fatalf("lane JSON retry = %v/%v, want %v/%v with no new Event", retried, retryEvent, original, event)
		}
	}
	archive := filepath.Join(t.TempDir(), "lane.zip")
	if result := f.run("export", "create", archive); result.err != nil {
		t.Fatalf("export create = %#v", result)
	}
	if result := f.run("export", "validate", archive); result.err != nil {
		t.Fatalf("export validate = %#v", result)
	}
}
