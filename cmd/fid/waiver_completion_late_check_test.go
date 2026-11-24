package main

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/journal"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestShowWaivedCompletionKeepsCheckBeforeDone(t *testing.T) {
	t.Parallel()
	for _, timing := range []string{"same-time", "next-file"} {
		t.Run(timing, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "assured", false, func(f *statusFixture) {
				path := filepath.Join(f.workspace, ".fid.toml")
				raw := string(mustReadCommandFile(t, path))
				raw = strings.Replace(raw, "[lanes.standard]\nchecks=['evidence']", "[lanes.standard]\nchecks=['lower']", 1)
				raw += "\n[checks.lower]\nrun=['true']\ntimeout='1m'\n"
				writeBackupTestFile(t, path, raw)
			})
			lifecycleResult(t, f.run("--json", "check", "lower", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "lower-pass"), "check.execute")
			var before map[string]any
			var beforeText string
			stamp := time.Now().UTC().Truncate(time.Second)
			// This check starts before done and finishes after. Its caller clock will be
			// equal to completion, or advance to another journal file: neither its
			// start order nor its timestamps may decide the historical result.
			options := checkrun.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, HoldID: hold, Actor: "human:builder", IdempotencyKey: "pending-lower", Names: []string{"lower"}, Clock: func() time.Time { return stamp }, Version: "test"}
			options.Before = func(string) error {
				lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
				waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", "standard"))
				lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
				done, _ := lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, task), nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
				var err error
				stamp, err = time.Parse(time.RFC3339, done["completion"].(map[string]any)["at"].(string))
				if err != nil {
					t.Fatal(err)
				}
				if timing == "next-file" {
					stamp = stamp.AddDate(0, 0, 1)
				}
				showTime := f.clock()
				f.clock = func() time.Time { return showTime }
				before, _ = lifecycleResult(t, f.run("--json", "show", task), "issue.show")
				checks := before["checks"].([]any)
				if len(checks) != 1 || checks[0].(map[string]any)["name"] != "lower" || checks[0].(map[string]any)["satisfied"] != true {
					t.Fatalf("fixture did not complete with a lower-lane-only pass: %v", checks)
				}
				beforeText = waivedCompletionText(t, f, task)
				return nil
			}
			response, err := checkrun.Run(context.Background(), options)
			if err == nil || len(response.Checks) != 1 || response.Checks[0].Counts || !strings.Contains(strings.Join(response.Checks[0].Reasons, ","), "ownership-lost") {
				t.Fatalf("expected a later non-counting ownership-lost result: %+v / %v", response, err)
			}
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
						t.Fatal(err)
					}
				}
				shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
				for _, key := range []string{"completion", "checks", "waivers", "approval", "no_other_model_review", "floor", "handoff", "warnings"} {
					if !reflect.DeepEqual(shown[key], before[key]) {
						t.Errorf("rebuild=%t: completed %s changed after the late result: before=%v, after=%v", rebuild, key, before[key], shown[key])
					}
				}
				text := waivedCompletionText(t, f, task)
				if text != beforeText || !strings.Contains(text, "lower passed") || strings.Contains(text, "no recorded pass") || strings.Contains(text, "ownership was lost") {
					t.Errorf("rebuild=%t: historical text changed: before=%s, after=%s", rebuild, beforeText, text)
				}
			}
		})
	}
}

func waivedCompletionText(t *testing.T, f *statusFixture, task string) string {
	t.Helper()
	out := f.run("show", task)
	if out.err != nil {
		t.Fatal(out)
	}
	return strings.Join(strings.Fields(out.stdout), " ")
}

func TestShowWaivedCompletionRequiresCanonicalBoundary(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "assured", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", "standard"))
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, task), nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	// A corrupted projection must not silently select unbounded evidence when
	// its receipt cannot be located in canonical history. Rebuild repairs it.
	startSQL(t, f, `UPDATE completions SET record_json=json_set(record_json,'$.id','CMP-missing')`)
	for _, args := range [][]string{{"show", task}, {"--json", "show", task}} {
		out := f.run(args...)
		if out.err == nil || !strings.Contains(out.stderr, "completion CMP-missing missing from journal") {
			t.Fatalf("missing completion boundary accepted: %#v", out)
		}
	}
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	waivedCompletionText(t, f, task)
}
