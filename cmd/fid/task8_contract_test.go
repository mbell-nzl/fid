package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	result := executeCommandForTest(
		workspace,
		"--project", workspace,
		"--json",
		"issue", "create",
		"--kind", "task",
		"--title", "Task 8 retained Issue",
		"--actor", "agent:task-8",
		"--idempotency-key", "task-8-create",
	)
	if result.err != nil {
		t.Fatalf("issue create error = %v, stderr = %q", result.err, result.stderr)
	}

	var envelope struct {
		Event struct {
			ID string `json:"id"`
		} `json:"event"`
		Result struct {
			IdempotencyKey string `json:"idempotency_key"`
			Replayed       bool   `json:"replayed"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode issue create response: %v; stdout = %q", err, result.stdout)
	}
	if envelope.Event.ID == "" || envelope.Result.IdempotencyKey != "task-8-create" || envelope.Result.Replayed {
		t.Fatalf("issue create response = %#v, want Event and non-replayed idempotency result", envelope)
	}
}

func TestTask8IssueCreateReplaysTheOriginalMutation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	args := []string{
		"--project", workspace,
		"--json",
		"issue", "create",
		"--kind", "task",
		"--title", "Task 8 replay",
		"--actor", "human:task-8",
		"--idempotency-key", "task-8-replay",
	}

	first := executeCommandForTest(workspace, args...)
	if first.err != nil {
		t.Fatalf("first issue create error = %v, stderr = %q", first.err, first.stderr)
	}
	second := executeCommandForTest(workspace, args...)
	if second.err != nil {
		t.Fatalf("replayed issue create error = %v, stderr = %q", second.err, second.stderr)
	}

	var firstEnvelope, replayEnvelope struct {
		Event struct {
			ID string `json:"id"`
		} `json:"event"`
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
			Replayed bool `json:"replayed"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(first.stdout), &firstEnvelope); err != nil {
		t.Fatalf("decode first issue create: %v", err)
	}
	if err := json.Unmarshal([]byte(second.stdout), &replayEnvelope); err != nil {
		t.Fatalf("decode replayed issue create: %v", err)
	}
	if firstEnvelope.Result.Replayed || !replayEnvelope.Result.Replayed || firstEnvelope.Event.ID != replayEnvelope.Event.ID || firstEnvelope.Result.Issue.ID != replayEnvelope.Result.Issue.ID {
		t.Fatalf("create/replay envelopes = (%#v, %#v), want original Event and Issue with replayed=true", firstEnvelope, replayEnvelope)
	}
}

func TestTask8IssueUpdateUsesStandaloneActorAndReplaysOriginalMutation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	issueID := createTask8Issue(t, workspace, "Task 8 update")
	args := []string{
		"--project", workspace,
		"--json",
		"issue", "update", issueID,
		"--title", "Task 8 update replayed",
		"--expected-revision", "1",
		"--actor", "service:task-8",
		"--idempotency-key", "task-8-update",
	}

	first := executeCommandForTest(workspace, args...)
	if first.err != nil {
		t.Fatalf("first issue update error = %v, stderr = %q", first.err, first.stderr)
	}
	replay := executeCommandForTest(workspace, args...)
	if replay.err != nil {
		t.Fatalf("replayed issue update error = %v, stderr = %q", replay.err, replay.stderr)
	}

	var firstEnvelope, replayEnvelope struct {
		Event struct {
			ID string `json:"id"`
		} `json:"event"`
		Result struct {
			IdempotencyKey string `json:"idempotency_key"`
			Replayed       bool   `json:"replayed"`
			Issue          struct {
				ID       string `json:"id"`
				Revision int    `json:"revision"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(first.stdout), &firstEnvelope); err != nil {
		t.Fatalf("decode first issue update: %v", err)
	}
	if err := json.Unmarshal([]byte(replay.stdout), &replayEnvelope); err != nil {
		t.Fatalf("decode replayed issue update: %v", err)
	}
	if firstEnvelope.Result.Replayed || !replayEnvelope.Result.Replayed || firstEnvelope.Event.ID != replayEnvelope.Event.ID || firstEnvelope.Result.Issue != replayEnvelope.Result.Issue || firstEnvelope.Result.IdempotencyKey != "task-8-update" || replayEnvelope.Result.IdempotencyKey != "task-8-update" {
		t.Fatalf("update/replay envelopes = (%#v, %#v), want original Event and result with replayed=true", firstEnvelope, replayEnvelope)
	}
}

func TestTask8IssueTransitionUsesStandaloneActorAndReplaysOriginalMutation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	issueID := createTask8Issue(t, workspace, "Task 8 transition")
	args := []string{
		"--project", workspace,
		"--json",
		"issue", "transition", issueID,
		"--to", "cancelled", "--reason", "Test cancellation",
		"--expected-revision", "1",
		"--actor", "human:task-8",
		"--idempotency-key", "task-8-transition",
	}

	first := executeCommandForTest(workspace, args...)
	if first.err != nil {
		t.Fatalf("first issue transition error = %v, stderr = %q", first.err, first.stderr)
	}
	replay := executeCommandForTest(workspace, args...)
	if replay.err != nil {
		t.Fatalf("replayed issue transition error = %v, stderr = %q", replay.err, replay.stderr)
	}

	var firstEnvelope, replayEnvelope struct {
		Event struct {
			ID string `json:"id"`
		} `json:"event"`
		Result struct {
			IdempotencyKey string `json:"idempotency_key"`
			Replayed       bool   `json:"replayed"`
			Issue          struct {
				ID       string `json:"id"`
				Revision int    `json:"revision"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(first.stdout), &firstEnvelope); err != nil {
		t.Fatalf("decode first issue transition: %v", err)
	}
	if err := json.Unmarshal([]byte(replay.stdout), &replayEnvelope); err != nil {
		t.Fatalf("decode replayed issue transition: %v", err)
	}
	if firstEnvelope.Result.Replayed || !replayEnvelope.Result.Replayed || firstEnvelope.Event.ID != replayEnvelope.Event.ID || firstEnvelope.Result.Issue != replayEnvelope.Result.Issue || firstEnvelope.Result.IdempotencyKey != "task-8-transition" || replayEnvelope.Result.IdempotencyKey != "task-8-transition" {
		t.Fatalf("transition/replay envelopes = (%#v, %#v), want original Event and result with replayed=true", firstEnvelope, replayEnvelope)
	}
}

func TestTask8IssueMutationsRejectLegacyByFlag(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	for _, args := range [][]string{
		{"issue", "create", "--kind", "task", "--title", "Legacy flag", "--actor", "agent:task-8", "--idempotency-key", "legacy-create", "--by", "developer"},
		{"issue", "update", "FID-missing", "--title", "Legacy flag", "--expected-revision", "1", "--actor", "agent:task-8", "--idempotency-key", "legacy-update", "--by", "developer"},
		{"issue", "transition", "FID-missing", "--to", "cancelled", "--expected-revision", "1", "--actor", "agent:task-8", "--idempotency-key", "legacy-transition", "--by", "developer"},
	} {
		result := executeCommandForTest(workspace, append([]string{"--project", workspace, "--json"}, args...)...)
		if result.err == nil || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
			t.Fatalf("%q result = (%v, %q), want rejected --by invalid-argument", args, result.err, result.stderr)
		}
	}
}

func createTask8Issue(t *testing.T, workspace, title string) string {
	t.Helper()
	result := executeCommandForTest(
		workspace,
		"--project", workspace,
		"--json",
		"issue", "create",
		"--kind", "task",
		"--title", title,
		"--actor", "agent:task-8",
		"--idempotency-key", "task-8-create-"+title,
	)
	if result.err != nil {
		t.Fatalf("issue create error = %v, stderr = %q", result.err, result.stderr)
	}
	var envelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode issue create response: %v", err)
	}
	return envelope.Result.Issue.ID
}
