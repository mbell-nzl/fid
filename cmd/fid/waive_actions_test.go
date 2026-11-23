package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func waiveForTest(f *statusFixture, task, answer string, terminal bool, extra ...string) commandExecutionResult {
	root := newRootCommandAt(f.workspace)
	now := f.clock
	if now == nil {
		now = time.Now
	}
	root.SetContext(context.WithValue(root.Context(), displayEnvironmentKey{}, displayEnvironment{now: now}))
	for _, c := range root.Commands() {
		if c.Name() == "waive" {
			root.RemoveCommand(c)
		}
	}
	root.AddCommand(newWaiveCommandWithTerminal(func(any) bool { return terminal }))
	var out, prompt bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&prompt)
	root.SetIn(strings.NewReader(answer))
	root.SetArgs(append([]string{"--json", "waive", task, "--actor", "human:mark", "--reason", "Reviewed risk", "--lane", "routine"}, extra...))
	err := root.Execute()
	return commandExecutionResult{stdout: out.String(), stderr: prompt.String(), err: err}
}

func TestWaiveRefusalsAndCancellation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, answer, want string
		terminal           bool
		args               []string
		prepare            func(*testing.T, *statusFixture)
	}{
		{name: "agent", terminal: true, args: []string{"--actor", "agent:builder"}, want: "only a person"},
		{name: "service", terminal: true, args: []string{"--actor", "service:ci"}, want: "only a person"},
		{name: "pipe", want: "standard input and standard error must be terminals"},
		{name: "empty reason", terminal: true, args: []string{"--reason", "   "}, want: "--reason is required"},
		{name: "cancel", terminal: true, answer: "yes\n", want: "Waiver cancelled"},
		{name: "EOF", terminal: true, answer: "FID-9f13aa", want: "Waiver cancelled"},
		{name: "missing objects", terminal: true, want: "restore the recorded base and head objects locally", prepare: func(t *testing.T, f *statusFixture) {
			head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
			objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
			if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			if tc.prepare != nil {
				tc.prepare(t, f)
			}
			before := reviewProjection(t, f)
			got := waiveForTest(f, "FID-9f13aa", tc.answer, tc.terminal, tc.args...)
			if !strings.Contains(got.stderr, tc.want) || got.stdout != "" || reviewProjection(t, f) != before {
				t.Fatalf("refusal=%#v want %s", got, tc.want)
			}
		})
	}
}

func TestWaiveRequiresReviewAndHandoff(t *testing.T) {
	t.Parallel()
	for _, reviewed := range []bool{false, true} {
		t.Run(map[bool]string{false: "not in review", true: "no handoff"}[reviewed], func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "[]")
			want := "task must be in_review"
			if reviewed {
				seedReviewState(t, f.workspace, "FID-9f13aa")
				want = "needs a hand-off"
			}
			before := reviewProjection(t, f)
			got := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true)
			if got.err == nil || !strings.Contains(got.stderr, want) || reviewProjection(t, f) != before {
				t.Fatalf("refusal=%#v", got)
			}
		})
	}
}

func TestWaivePromptAndRetry(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	got := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "waiver")
	fields := waiveSuccess(t, got)
	if fields["decision"].(map[string]any)["decision_id"] == "" {
		t.Fatal(fields)
	}
	for _, part := range []string{"FID-9f13aa · standard → routine", "Floor:", "Hard floor:", "New lane: routine", "Type FID-9f13aa to waive:"} {
		if !strings.Contains(got.stderr, part) {
			t.Errorf("prompt missing %q: %s", part, got.stderr)
		}
	}
	before := reviewProjection(t, f)
	again := waiveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "waiver")
	if again.err != nil || again.stdout != got.stdout || again.stderr != "" || reviewProjection(t, f) != before {
		t.Fatalf("retry=%#v", again)
	}
	conflict := waiveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "waiver", "--reason", "Changed reason")
	assertJSONFailureCode(t, conflict.stderr, "task.waive", "idempotency-conflict")
	laneConflict := waiveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "waiver", "--lane", "standard")
	assertJSONFailureCode(t, laneConflict.stderr, "task.waive", "idempotency-conflict")
}

func waiveSuccess(t *testing.T, got commandExecutionResult) map[string]any {
	t.Helper()
	if got.err != nil {
		t.Fatalf("waive failed: %#v", got)
	}
	got.stderr = ""
	fields, _ := lifecycleResult(t, got, "task.waive")
	return fields
}

func TestWaiveMustLowerHandoffLane(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	for _, lane := range []string{"standard", "assured", "bogus"} {
		before := reviewProjection(t, f)
		got := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", lane)
		if got.err == nil || !strings.Contains(got.stderr, "waiver-not-lower") || reviewProjection(t, f) != before {
			t.Fatalf("lane %s: %#v", lane, got)
		}
	}
}

func TestWaiveRefusesHardFloorPaths(t *testing.T) {
	t.Parallel()
	for _, path := range []string{".fid.toml", ".github/workflows/ci.yml"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			f := laneFloorFixture(t, "standard")
			if path == ".fid.toml" {
				old := string(mustReadCommandFile(t, filepath.Join(f.workspace, path)))
				writeBackupTestFile(t, filepath.Join(f.workspace, path), old+"\n# Task edit\n")
			} else {
				writeBackupTestFile(t, filepath.Join(f.workspace, path), "name: task\n")
			}
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Hard floor change")
			requiredNames(t, checkRequired(t, f, "hard-checks"))
			reviewHandoff(t, f, "hard-handoff")
			before := reviewProjection(t, f)
			got := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", "standard")
			assertJSONFailureCode(t, got.stderr, "task.waive", "waiver-hard-floor")
			if !strings.Contains(got.stderr, path) || got.stdout != "" || reviewProjection(t, f) != before {
				t.Fatal(got)
			}
		})
	}
}

func TestWaiveRetryWithoutObjectsAndCorruptResponses(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"missing objects", "decision", "occurred_at", "reason", "binding"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			first := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "retry")
			waiveSuccess(t, first)
			if field == "missing objects" {
				head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
				objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
				path := filepath.Join(objects, head[:2], head[2:])
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.WriteFile(path, raw, 0o444); err != nil {
						t.Error(err)
					}
				})
			} else {
				mutateStoredRetry(t, f, "retry", func(result map[string]any) {
					switch field {
					case "binding":
						result["decision"].(map[string]any)["binding"].(map[string]any)["head"] = strings.Repeat("e", 40)
					case "reason":
						result["decision"].(map[string]any)["reason"] = "Forged reason"
					default:
						delete(result, field)
					}
				})
			}
			before := reviewProjection(t, f)
			got := waiveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "retry")
			if field == "missing objects" {
				if got.err != nil || got.stdout != first.stdout || got.stderr != "" {
					t.Fatalf("retry=%#v", got)
				}
			} else if got.err == nil {
				t.Fatal("accepted corrupt retry", got)
			}
			if reviewProjection(t, f) != before {
				t.Fatal("retry changed state")
			}
		})
	}
}

func TestIssueLaneLoweringPointsToWaive(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	got := f.run("issue", "update", "FID-9f13aa", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "lower")
	if got.err == nil || !strings.Contains(got.stderr, "fid waive FID-9f13aa") {
		t.Fatal("lowering refusal lacks waiver guidance", got)
	}
}
