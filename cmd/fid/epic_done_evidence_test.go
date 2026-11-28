package main

import (
	"strings"
	"testing"
	"time"
)

func TestEpicDoneEvidenceRefusals(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"missing approval", "stale approval", "missing pass", "stale pass", "failed pass"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			if state != "missing approval" {
				verifier := "review:code"
				if state == "failed pass" {
					appendPolicy(t, f.workspace, "\n[checks.failure]\nrun=['false']\ntimeout='5s'\n")
					commit := policyCommit(t, f.workspace)
					identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
					verifier = "check:failure"
				}
				lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Claim", "--verifier", verifier, "--actor", "agent:builder", "--idempotency-key", "claim"), "epic.validation.add")
				approveTestEpic(t, f, epic)
				if state == "stale pass" {
					lifecycleResult(t, f.run("--json", "epic", "validation", "pass", epic, "V1", "--actor", "human:mark", "--reason", "Reviewed"), "epic.validation.pass")
				}
				if state == "stale approval" || state == "stale pass" {
					lifecycleResult(t, f.run("--json", "epic", "validation", "update", epic, "V1", "--claim", "Edited claim", "--verifier", verifier, "--actor", "agent:builder", "--idempotency-key", "edit"), "epic.validation.update")
				}
				if state == "stale pass" {
					approveTestEpic(t, f, epic)
				}
				if state == "failed pass" {
					lifecycleResult(t, f.run("--json", "epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "verify"), "epic.validation.verify")
				}
			}
			before := snapshotWorkspacePersistence(t, f.workspace)
			r := f.run("--json", "done", epic, "--actor", "human:mark", "--reason", "Finished")
			assertJSONFailureCode(t, r.stderr, "task.done", "done-epic-evidence-required")
			want := map[string]string{"missing approval": "approval: missing", "stale approval": "approval: stale", "missing pass": "V1 review:code: missing", "stale pass": "V1 review:code: stale", "failed pass": "V1 check:failure: failed"}[state]
			if !strings.Contains(r.stderr, want) {
				t.Fatalf("missing explanation %s: %#v", want, r)
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("refused close changed ledger")
			}
		})
	}
}
