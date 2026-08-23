package main

import (
	"encoding/json"
	"testing"
)

func TestTask8IssueCreateRequiresStandaloneActorAndIdempotencyKey(t *testing.T) {
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
