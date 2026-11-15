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

func approveForTest(f *statusFixture, task, answer string, terminal bool, extra ...string) commandExecutionResult {
	root := newRootCommandAt(f.workspace)
	now := f.clock
	if now == nil {
		now = time.Now
	}
	root.SetContext(context.WithValue(root.Context(), displayEnvironmentKey{}, displayEnvironment{now: now}))
	for _, c := range root.Commands() {
		if c.Name() == "approve" {
			root.RemoveCommand(c)
		}
	}
	root.AddCommand(newApproveCommandWithTerminal(func(any) bool { return terminal }))
	var out, prompt bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&prompt)
	root.SetIn(strings.NewReader(answer))
	root.SetArgs(append([]string{"--json", "approve", task, "--actor", "human:mark", "--reason", "Reviewed the full change"}, extra...))
	err := root.Execute()
	return commandExecutionResult{stdout: out.String(), stderr: prompt.String(), err: err}
}

func TestApproveRefusalsAndCancellation(t *testing.T) {
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
		{name: "cancel", terminal: true, answer: "yes\n", want: "Approval cancelled"},
		{name: "EOF", terminal: true, answer: "FID-9f13aa", want: "Approval cancelled"},
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
			got := approveForTest(f, "FID-9f13aa", tc.answer, tc.terminal, tc.args...)
			if !strings.Contains(got.stderr, tc.want) || got.stdout != "" || reviewProjection(t, f) != before {
				t.Fatalf("refusal=%#v want %s", got, tc.want)
			}
		})
	}
}

func TestApproveRequiresReviewAndHandoff(t *testing.T) {
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
			got := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true)
			if got.err == nil || !strings.Contains(got.stderr, want) || reviewProjection(t, f) != before {
				t.Fatalf("refusal=%#v", got)
			}
		})
	}
}

func TestApprovePromptAndRetry(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	got := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "approval")
	fields := approveSuccess(t, got)
	if fields["decision"].(map[string]any)["decision_id"] == "" {
		t.Fatal(fields)
	}
	for _, part := range []string{"FID-9f13aa · standard", "notes.md", "1 file changed", "Checks: none required", "primary: missing", "Earlier approvals: none", "Type FID-9f13aa to approve:"} {
		if !strings.Contains(got.stderr, part) {
			t.Errorf("prompt missing %q: %s", part, got.stderr)
		}
	}
	before := reviewProjection(t, f)
	again := approveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "approval")
	if again.err != nil || again.stdout != got.stdout || again.stderr != "" || reviewProjection(t, f) != before {
		t.Fatalf("retry=%#v", again)
	}
	conflict := approveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "approval", "--reason", "Changed reason")
	assertJSONFailureCode(t, conflict.stderr, "task.approve", "idempotency-conflict")
}

func TestApproveOtherModelDisclosure(t *testing.T) {
	t.Parallel()
	for _, model := range []string{"", "codex/gpt-6-astra", "claude/opus-5.5"} {
		t.Run(model, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			in := reviewInput()
			in["reviewer"].(map[string]any)["model"] = model
			lifecycleResult(t, recordReviewTask(t, f, in, "review", "human:courier", task), "review.record")
			got := approveForTest(f, task, task+"\n", true)
			approveSuccess(t, got)
			if strings.Contains(got.stderr, "no other-model review") != (model != "claude/opus-5.5") || !strings.Contains(got.stderr, "APPROVE") || !strings.Contains(got.stderr, "agent:reviewer") || !strings.Contains(got.stderr, "current") || !strings.Contains(got.stderr, "separate: missing") {
				t.Fatal(got.stderr)
			}
		})
	}
}

func approveSuccess(t *testing.T, got commandExecutionResult) map[string]any {
	t.Helper()
	if got.err != nil {
		t.Fatalf("approve failed: %#v", got)
	}
	got.stderr = ""
	fields, _ := lifecycleResult(t, got, "task.approve")
	return fields
}

func TestApproveRetryWithoutObjectsAndCorruptResponses(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"missing objects", "decision", "occurred_at", "reason"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			first := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "retry")
			approveSuccess(t, first)
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
					if field == "reason" {
						result["decision"].(map[string]any)["reason"] = "Forged reason"
					} else {
						delete(result, field)
					}
				})
			}
			before := reviewProjection(t, f)
			got := approveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "retry")
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

func TestApproveDisclosesUnknownHistoricalBuilderModel(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "assured")
	// Historical agent holds predate required model attribution. A human's empty
	// model means manual work; an agent's empty model cannot prove other-model review.
	startSQL(t, f, `UPDATE holds SET model='' WHERE actor='agent:builder'`)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "review", "human:courier", task), "review.record")
	got := approveForTest(f, task, task+"\n", true)
	approveSuccess(t, got)
	if !strings.Contains(got.stderr, "no other-model review") {
		t.Fatal(got.stderr)
	}
}

func TestApproveRefusesOversizedEvent(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	before := reviewProjection(t, f)
	got := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "oversized", "--reason", strings.Repeat("x", 600000))
	if got.err == nil || !strings.Contains(got.stderr, "approval Event exceeds") || reviewProjection(t, f) != before {
		t.Fatalf("oversized approval: %v %s", got.err, got.stderr)
	}
}

func TestApproveRetryCannotReturnAnotherDecision(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "first"))
	approveSuccess(t, approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "second", "--reason", "Reviewed again"))
	startSQL(t, f, `UPDATE idempotency_results SET result_json=(SELECT result_json FROM idempotency_results WHERE idempotency_key='second'), result_entity_id=(SELECT result_entity_id FROM idempotency_results WHERE idempotency_key='second') WHERE idempotency_key='first'`)
	before := reviewProjection(t, f)
	got := approveForTest(f, "FID-9f13aa", "", true, "--idempotency-key", "first")
	if got.err == nil || reviewProjection(t, f) != before {
		t.Fatal("retry returned a different decision", got)
	}
}

func TestApproveEmptyChangeDoesNotSuggestRestoringObjects(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "[]")
	reviewHandoff(t, f, "empty-handoff")
	before := reviewProjection(t, f)
	got := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true)
	if got.err == nil || !strings.Contains(got.stderr, "fingerprint empty-change") || strings.Contains(got.stderr, "restore the recorded base and head") || reviewProjection(t, f) != before {
		t.Fatalf("empty approval: %#v", got)
	}
}
