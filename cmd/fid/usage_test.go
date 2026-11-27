package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUsageRecord(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, 100*time.Hour)
	f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Build", "todo", nil, nil, nil)}, nil)
	next, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "openai/sol"), "next.claim")
	hold := next["hold"].(map[string]any)["id"].(string)
	path := filepath.Join(t.TempDir(), "stdout.jsonl")
	raw, err := os.ReadFile("../../internal/usagerecord/testdata/codex.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"--json", "usage", "record", "--hold", hold, "--actor", "agent:builder", "--format", "codex", "--input", path, "--tool-version", "codex-cli 0.172.0", "--idempotency-key", "usage-1"}
	usageGolden := func(name string) {
		for _, format := range []string{"txt", "json"} {
			args := []string{"show", "FID-aaaaaa"}
			if format == "json" {
				args = append([]string{"--json"}, args...)
			}
			out := f.run(args...)
			if out.err != nil {
				t.Fatal(out)
			}
			assertGolden(t, filepath.Join("show", name+"."+format), strings.ReplaceAll(f.placeholders(out.stdout), hold, "HLD-session"))
		}
	}
	usageGolden("usage-none")
	before := startEventCount(t, f.workspace)
	wrong := append([]string{}, args...)
	wrong[6] = "agent:other"
	denied := f.run(wrong...)
	assertJSONFailureCode(t, denied.stderr, "usage.record", "authority-denied")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("refusal wrote an Event")
	}
	first := f.run(args...)
	result, _ := lifecycleResult(t, first, "usage.record")
	usage := result["usage"].(map[string]any)
	if usage["role"] != "build" || usage["model"] != "openai/sol" || usage["task_id"] != "FID-aaaaaa" || usage["hold_id"] != hold || usage["counts"].(map[string]any)["total"] != float64(15752) || usage["version_source"] != "caller" {
		t.Fatal(usage)
	}
	if again := f.run(args...); again.err != nil || again.stdout != first.stdout {
		t.Fatalf("retry: %#v", again)
	}
	if startEventCount(t, f.workspace) != before+1 {
		t.Fatal("retry wrote an Event")
	}
	duplicate := append([]string{}, args...)
	duplicate[len(duplicate)-1] = "usage-2"
	got := f.run(duplicate...)
	assertJSONFailureCode(t, got.stderr, "usage.record", "duplicate-entity")
	if startEventCount(t, f.workspace) != before+1 {
		t.Fatal("duplicate wrote an Event")
	}
	changed := append([]string{}, args...)
	changed[12] = "different tool version"
	conflict := f.run(changed...)
	assertJSONFailureCode(t, conflict.stderr, "usage.record", "idempotency-conflict")
	if startEventCount(t, f.workspace) != before+1 {
		t.Fatal("conflicting retry wrote an Event")
	}

	for _, tc := range []struct {
		index int
		value string
	}{{8, "unknown"}, {10, "/missing/input"}, {8, "claude"}, {4, "HLD-missing"}} {
		bad := append([]string{}, args...)
		bad[tc.index] = tc.value
		bad[len(bad)-1] = "bad-" + tc.value
		before := startEventCount(t, f.workspace)
		got := f.run(bad...)
		if got.err == nil || startEventCount(t, f.workspace) != before {
			t.Fatalf("refusal changed ledger: %#v", got)
		}
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
	if shown["usage"].(map[string]any)["total"] != float64(15752) {
		t.Fatal(shown)
	}
	usageGolden("usage")
	text := f.run("show", "FID-aaaaaa")
	if !strings.Contains(text.stdout, "Usage") || !strings.Contains(text.stdout, "15752") {
		t.Fatal(text)
	}
}

func TestUsageAfterWaitingPark(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, 100*time.Hour)
	f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Build", "in_progress", nil, nil, nil)}, nil)
	f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
	hold := "HLD-aaaaaa"
	lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", "agent:claude-builder", "--waiting", "Which rule?", "--idempotency-key", "wait"), "hold.park")
	path := filepath.Join(t.TempDir(), "out")
	raw, err := os.ReadFile("../../internal/usagerecord/testdata/codex.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	lifecycleResult(t, f.run("--json", "usage", "record", "--hold", hold, "--actor", "agent:claude-builder", "--format", "codex", "--input", path, "--idempotency-key", "parked-usage"), "usage.record")
}

func TestUsagePreservesLargeIntegerCounts(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Large usage", "in_progress", nil, nil, nil)}, nil)
	f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
	path := filepath.Join(t.TempDir(), "large.jsonl")
	raw := `{"type":"turn.completed","usage":{"input_tokens":9007199254740993,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":0,"reasoning_output_tokens":0}}`
	if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
		t.Fatal(err)
	}
	got := f.run("--json", "usage", "record", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--format", "codex", "--input", path, "--idempotency-key", "large")
	if got.err != nil || !strings.Contains(got.stdout, `"total":9007199254740993`) {
		t.Fatalf("large usage: %#v", got)
	}
}

func TestUsageVersionSources(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, header, caller, want, source string }{
		{"output", `{"type":"system","claude_code_version":"2.1.332"}`, "caller v0", "2.1.332", "output"},
		{"caller", "", "caller v0", "caller v0", "caller"},
		{"unavailable", "", "", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Versions", "in_progress", nil, nil, nil)}, nil)
			f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
			raw, err := os.ReadFile("../../internal/usagerecord/testdata/claude.json")
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(t.TempDir(), "claude.stdout")
			if err := os.WriteFile(path, append([]byte(tc.header+"\n"), raw...), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"--json", "usage", "record", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--format", "claude", "--input", path, "--idempotency-key", "version"}
			if tc.caller != "" {
				args = append(args, "--tool-version", tc.caller)
			}
			got, _ := lifecycleResult(t, f.run(args...), "usage.record")
			r := got["usage"].(map[string]any)
			if r["tool_version"] != tc.want || r["version_source"] != tc.source {
				t.Fatal(r)
			}
		})
	}
}
