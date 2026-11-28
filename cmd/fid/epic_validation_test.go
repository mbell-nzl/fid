package main

import (
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestEpicValidationLifecycle(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	add := func(key string) map[string]any {
		t.Helper()
		r, _ := lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "A person verifies the fixture", "--verifier", "human", "--recipe", "docs/fixture.md", "--actor", "agent:builder", "--idempotency-key", key), "epic.validation.add")
		return r["validation"].(map[string]any)
	}
	if got := add("first")["id"]; got != "V1" {
		t.Fatalf("first ID=%v", got)
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "update", epic, "V1", "--claim", "An updated claim", "--verifier", "human", "--recipe", "docs/updated.md", "--actor", "agent:builder", "--idempotency-key", "update"), "epic.validation.update")
	lifecycleResult(t, f.run("--json", "epic", "validation", "remove", epic, "V1", "--actor", "agent:builder", "--idempotency-key", "remove"), "epic.validation.remove")
	if got := add("second")["id"]; got != "V2" {
		t.Fatalf("ID reused: %v", got)
	}
	if r := f.run("rebuild"); r.err != nil {
		t.Fatalf("rebuild=%#v", r)
	}
	show := f.run("--json", "show", epic)
	if show.err != nil || !strings.Contains(show.stdout, `"validations"`) || !strings.Contains(show.stdout, `"V2"`) || strings.Contains(show.stdout, `"V1"`) {
		t.Fatalf("show=%#v", show)
	}
}

func approveTestEpic(t *testing.T, f *statusFixture, epic string) {
	t.Helper()
	_, err := epicwork.Approve(epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "human:mark", Reason: "Approve test plan", Confirm: func(epicplan.Plan) (bool, error) { return true, nil }})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEpicApprovalGatesOnlyUnstartedWork(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	task := epicDoneChild(t, f, epic, "planned-task")
	r, _ := lifecycleResult(t, f.run("--json", "next", "--epic", epic), "next.get")
	if r["kind"] != "none" || !strings.Contains(f.run("--json", "next", "--explain").stdout, "epic-not-approved") {
		t.Fatalf("unapproved work served: %v", r)
	}
	approveTestEpic(t, f, epic)
	r, _ = lifecycleResult(t, f.run("--json", "next", "--epic", epic), "next.get")
	if r["kind"] != "new" || r["task"].(map[string]any)["id"] != task {
		t.Fatalf("approved work not served: %v", r)
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "New claim", "--verifier", "human", "--recipe", "recipe.md", "--actor", "agent:builder", "--idempotency-key", "stale"), "epic.validation.add")
	r, _ = lifecycleResult(t, f.run("--json", "next", "--epic", epic), "next.get")
	if r["kind"] != "none" {
		t.Fatalf("stale plan served: %v", r)
	}
}

func TestEpicValidationRejectsInvalidVerifier(t *testing.T) {
	t.Parallel()
	for _, verifier := range []string{"bogus", "check:missing", "review:missing", "human"} {
		t.Run(verifier, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			r := f.run("--json", "epic", "validation", "add", epic, "--claim", "Claim", "--verifier", verifier, "--actor", "agent:builder", "--idempotency-key", "invalid")
			if r.err == nil {
				t.Fatalf("accepted invalid validation: %#v", r)
			}
		})
	}
}

func TestEpicValidationPersonPassAndClose(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Review complete", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "claim"), "epic.validation.add")
	approveTestEpic(t, f, epic)
	before := f.run("--json", "done", epic, "--actor", "human:mark", "--reason", "Finished")
	if before.err == nil {
		t.Fatal("closed without validation pass")
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "pass", epic, "V1", "--actor", "human:mark", "--reason", "Reviewed the evidence", "--idempotency-key", "pass"), "epic.validation.pass")
	lifecycleResult(t, f.run("--json", "done", epic, "--actor", "human:mark", "--reason", "Finished"), "task.done")
}
