package main

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestStartRetryRejectsMalformedResult(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"lane", "task_id", "hold_id", "status"} {
		for _, row := range []struct {
			name  string
			value any
		}{
			{"missing", nil}, {"null", nil}, {"number", 42}, {"empty", ""},
			{"array", []any{}}, {"object", map[string]any{}}, {"boolean", true},
			{"unknown", "urgent"},
		} {
			if row.name == "unknown" && field != "lane" && field != "status" {
				continue
			}
			for _, mode := range []string{"text", "json"} {
				t.Run(field+"/"+row.name+"/"+mode, func(t *testing.T) {
					t.Parallel()
					f := plannedStartFixture(t, "todo", stringPointer("standard"))
					args := []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}
					lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.start")
					mutateStoredRetry(t, f, "start", func(result map[string]any) {
						fields := result["fields"].(map[string]any)
						fields[field] = row.value
						if row.name == "missing" {
							delete(fields, field)
						}
					})
					if mode == "json" {
						args = append([]string{"--json"}, args...)
					}
					assertMalformedRetryRefused(t, f, args, "task start result")
				})
			}
		}
	}
}

func TestStartRetryRejectsMalformedEnvelope(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"fields", "occurred_at"} {
		for _, row := range []struct {
			name  string
			value any
		}{
			{"missing", nil}, {"null", nil}, {"empty", ""}, {"number", 42}, {"array", []any{}},
		} {
			for _, mode := range []string{"text", "json"} {
				t.Run(field+"/"+row.name+"/"+mode, func(t *testing.T) {
					t.Parallel()
					f := plannedStartFixture(t, "todo", stringPointer("standard"))
					args := []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}
					lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.start")
					mutateStoredRetry(t, f, "start", func(result map[string]any) {
						result[field] = row.value
						if row.name == "missing" {
							delete(result, field)
						}
					})
					if mode == "json" {
						args = append([]string{"--json"}, args...)
					}
					assertMalformedRetryRefused(t, f, args, "task start result")
				})
			}
		}
	}
}

func TestStartRetryPreservesTextAndJSONAcrossRebuild(t *testing.T) {
	t.Parallel()
	f := newJournalStatusFixture(t)
	now := time.Now().Add(time.Minute)
	f.clock = func() time.Time { return now }
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--title", "Planned epic", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "epic"), "issue.create")
	epic := created["issue"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--epic", epic, "--title", "Planned task", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "task"), "issue.create")
	taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
	hold := taken["hold"].(map[string]any)["id"].(string)
	args := []string{"start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start"}
	original, event := lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.start")
	events := startEventCount(t, f.workspace)
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if got := f.run("rebuild"); got.err != nil {
				t.Fatalf("rebuild = %#v", got)
			}
		}
		want := "✓ Started " + original["task_id"].(string) + " · Standard\n"
		if got := f.run(args...); got.err != nil || got.stderr != "" || got.stdout != want {
			t.Fatalf("start text retry = %#v, want %q", got, want)
		}
		retried, retryEvent := lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.start")
		if !reflect.DeepEqual(retried, original) || !reflect.DeepEqual(retryEvent, event) || events != startEventCount(t, f.workspace) {
			t.Fatalf("start retry = %v/%v, want %v/%v with no new Event", retried, retryEvent, original, event)
		}
	}
	archive := filepath.Join(t.TempDir(), "start.zip")
	if got := f.run("export", "create", archive); got.err != nil {
		t.Fatalf("export create = %#v", got)
	}
	if got := f.run("export", "validate", archive); got.err != nil {
		t.Fatalf("export validate = %#v", got)
	}
}
