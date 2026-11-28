package main

import (
	"strings"
	"testing"
	"time"
)

func TestEpicPlanAndValidationDisplayTextAndJSON(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	task := epicDoneChild(t, f, epic, "planned")
	for _, format := range []string{"text", "json"} {
		args := []string{"status"}
		if format == "json" {
			args = append(args, "--json")
		}
		r := f.run(args...)
		if r.err != nil || !strings.Contains(r.stdout, "fid epic approve "+epic) {
			t.Fatalf("status %s=%#v", format, r)
		}
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Review complete", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "claim"), "epic.validation.add")
	approveTestEpic(t, f, epic)
	shown, _ := lifecycleResult(t, f.run("--json", "show", epic), "issue.show")
	if shown["epic_approval"].(map[string]any)["state"] != "current" || len(shown["plan"].([]any)) != 1 || shown["validations"].([]any)[0].(map[string]any)["state"] != "missing" {
		t.Fatalf("show=%v", shown)
	}
	text := f.run("show", epic)
	if text.err != nil {
		t.Fatal(text.err)
	}
	for _, part := range []string{"Approval", "current", "human:mark", "Plan", "VALIDATIONS", "V1", "Review complete", "review:code", "missing"} {
		if !strings.Contains(text.stdout, part) {
			t.Fatalf("show missing %s: %s", part, text.stdout)
		}
	}
	lifecycleResult(t, f.run("--json", "issue", "transition", task, "--to", "cancelled", "--actor", "human:mark", "--reason", "No delivery needed"), "issue.transition")
	status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	pass, closeIndex := -1, -1
	for i, item := range status["needs_you"].([]any) {
		v := item.(map[string]any)
		if v["kind"] == "epic-validation-pass" {
			pass = i
		}
		if v["kind"] == "epic-ready-to-close" {
			closeIndex = i
		}
	}
	if pass < 0 || closeIndex < pass {
		t.Fatalf("validation action must precede close: %v", status)
	}
	if r := f.run("status"); r.err != nil || !strings.Contains(r.stdout, "fid epic validation pass") {
		t.Fatalf("validation status=%#v", r)
	}
}
