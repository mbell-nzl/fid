package main

import (
	"testing"
	"time"
)

func addV8EpicValidations(t *testing.T, f *statusFixture, epic, child string) {
	t.Helper()
	appendPolicy(t, f.workspace, "\n[checks.epicproof]\nrun=['/bin/true']\ntimeout='5s'\nartifacts=['lean-artifact.txt']\n[checks.epicfailure]\nrun=['/bin/false']\ntimeout='5s'\n")
	commit := policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
	for _, v := range []struct{ claim, verifier, recipe string }{{"Inspect fixture", "human", ".fid.toml"}, {"Check default", "check:epicproof", ""}, {"Failed evidence", "check:epicfailure", ""}, {"Review evidence", "review:code", ""}} {
		args := []string{"--json", "epic", "validation", "add", epic, "--claim", v.claim, "--verifier", v.verifier, "--actor", "agent:builder", "--idempotency-key", "v8-" + v.verifier}
		if v.recipe != "" {
			args = append(args, "--recipe", v.recipe)
		}
		lifecycleResult(t, f.run(args...), "epic.validation.add")
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "v8-verify"), "epic.validation.verify")
	lifecycleResult(t, f.run("--json", "epic", "validation", "remove", epic, "V3", "--actor", "agent:builder", "--idempotency-key", "v8-remove-validation"), "epic.validation.remove")
	for _, id := range []string{"V1", "V4"} {
		lifecycleResult(t, f.run("--json", "epic", "validation", "pass", epic, id, "--actor", "human:mark", "--reason", "Inspected evidence", "--idempotency-key", "v8-pass-"+id), "epic.validation.pass")
	}
	lifecycleResult(t, f.run("--json", "epic", "validation", "update", epic, "V4", "--claim", "Updated review claim", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "v8-update-validation"), "epic.validation.update")
	lifecycleResult(t, f.run("--json", "epic", "validation", "pass", epic, "V4", "--actor", "human:mark", "--reason", "Inspected updated evidence", "--idempotency-key", "v8-pass-updated"), "epic.validation.pass")
	lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "subtask", "--parent", child, "--title", "Accepted by epic plan", "--actor", "agent:builder", "--idempotency-key", "v8-plan-task"), "issue.create")
	approveTestEpic(t, f, epic)
}

func captureV8EpicValidations(t *testing.T, root string, state *v8State) {
	t.Helper()
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	for _, table := range []string{"epic_validation_changes", "epic_validation_results", "epic_plan_acceptances"} {
		r := f.run("sql", "SELECT * FROM "+table+" ORDER BY rowid")
		if r.err != nil {
			t.Fatalf("%s=%#v", table, r)
		}
		if r.stdout == "[]\n" {
			t.Fatalf("V8 must retain %s", table)
		}
		state.reads[table] = r.stdout
	}
}
