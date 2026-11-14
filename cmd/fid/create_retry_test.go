package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func createRetryFixture(t *testing.T, kind, lane string) (*statusFixture, []string, commandExecutionResult) {
	t.Helper()
	f := newJournalStatusFixture(t)
	args := []string{"issue", "create", "--kind", kind, "--title", "Plan work", "--actor", "human:mark", "--idempotency-key", "create"}
	if kind == "epic" {
		args = append(args, "--scheduling", "prioritised")
	}
	if lane != "" {
		args = append(args, "--lane", lane)
	}
	first := f.run(append([]string{"--json"}, args...)...)
	lifecycleResult(t, first, "issue.create")
	return f, args, first
}

func mutateStoredRetry(t *testing.T, f *statusFixture, key string, mutate func(map[string]any)) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	var raw string
	if err := conn.QueryRow(`SELECT result_json FROM idempotency_results WHERE idempotency_key=?`, key).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatal(err)
	}
	mutate(result)
	payload, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`UPDATE idempotency_results SET result_json=? WHERE idempotency_key=?`, string(payload), key); err != nil {
		t.Fatal(err)
	}
}

func assertMalformedRetryRefused(t *testing.T, f *statusFixture, args []string, message string) {
	t.Helper()
	defer func() {
		if caught := recover(); caught != nil {
			t.Errorf("malformed retry panicked: %v", caught)
		}
	}()
	query := "SELECT id,kind,status,revision,lane FROM issues ORDER BY id"
	before := f.run("sql", query)
	events := startEventCount(t, f.workspace)
	got := f.run(args...)
	if got.err == nil || got.stdout != "" || !strings.Contains(got.stderr, message) {
		t.Errorf("malformed retry = %#v, want refusal containing %q", got, message)
	}
	after := f.run("sql", query)
	if before.err != nil || after.err != nil || before.stdout != after.stdout || events != startEventCount(t, f.workspace) {
		t.Fatalf("refused retry changed issues or journal: before=%#v after=%#v", before, after)
	}
}

func TestCreateRetryRejectsMalformedLane(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "epic"} {
		for _, row := range []struct {
			name  string
			value any
		}{
			{"empty", ""}, {"unknown", "urgent"}, {"number", 42}, {"boolean", true},
			{"array", []any{}}, {"object", map[string]any{}}, {"epic-lane", "standard"},
		} {
			if kind == "task" && row.name == "epic-lane" {
				continue
			}
			for _, mode := range []string{"text", "json"} {
				t.Run(kind+"/"+row.name+"/"+mode, func(t *testing.T) {
					t.Parallel()
					f, args, _ := createRetryFixture(t, kind, "")
					mutateStoredRetry(t, f, "create", func(result map[string]any) {
						result["issue"].(map[string]any)["lane"] = row.value
					})
					if mode == "json" {
						args = append([]string{"--json"}, args...)
					}
					assertMalformedRetryRefused(t, f, args, "Issue create result")
				})
			}
		}
	}
}

func TestCreateRetryPreservesLaneAndUnsetAcrossRebuild(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ kind, lane string }{
		{"task", ""}, {"task", "routine"}, {"task", "standard"}, {"task", "assured"}, {"epic", ""},
	} {
		t.Run(row.kind+"/"+row.lane, func(t *testing.T) {
			t.Parallel()
			f, args, first := createRetryFixture(t, row.kind, row.lane)
			original, event := lifecycleResult(t, first, "issue.create")
			if row.lane == "" {
				mutateStoredRetry(t, f, "create", func(result map[string]any) {
					delete(result["issue"].(map[string]any), "lane")
				})
			}
			issue := original["issue"].(map[string]any)
			want := fmt.Sprintf("created issue %s (%s) %q", issue["id"], row.kind, issue["title"])
			if row.lane != "" {
				want += " · " + strings.ToUpper(row.lane[:1]) + row.lane[1:]
			}
			want += "\n"
			events := startEventCount(t, f.workspace)
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if got := f.run("rebuild"); got.err != nil {
						t.Fatalf("rebuild = %#v", got)
					}
				}
				if got := f.run(args...); got.err != nil || got.stderr != "" || got.stdout != want {
					t.Fatalf("creation text retry = %#v, want %q", got, want)
				}
				retried, retryEvent := lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "issue.create")
				original["replayed"] = true
				if !reflect.DeepEqual(retried, original) || !reflect.DeepEqual(retryEvent, event) || events != startEventCount(t, f.workspace) {
					t.Fatalf("creation retry = %v/%v, want %v/%v with no new Event", retried, retryEvent, original, event)
				}
			}
			archive := filepath.Join(t.TempDir(), "create.zip")
			if got := f.run("export", "create", archive); got.err != nil {
				t.Fatalf("export create = %#v", got)
			}
			if got := f.run("export", "validate", archive); got.err != nil {
				t.Fatalf("export validate = %#v", got)
			}
		})
	}
}
