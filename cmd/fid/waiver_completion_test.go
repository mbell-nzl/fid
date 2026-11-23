package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestShowWaivedCompletionReceipt(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"standard", "routine", "standard-extra"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			lane := strings.TrimSuffix(name, "-extra")
			f, task, hold := doneBuildFixture(t, "assured", false, func(f *statusFixture) {
				path := filepath.Join(f.workspace, ".fid.toml")
				raw := string(mustReadCommandFile(t, path))
				raw = strings.Replace(raw, "checks=['evidence']", "checks=[]", 1)
				raw = strings.Replace(raw, "[lanes.standard]\nchecks=['evidence']", "[lanes.standard]\nchecks=['zunit','conditional']", 1)
				raw = strings.Replace(raw, "[lanes.assured]\nchecks=['evidence']", "[lanes.assured]\nchecks=['evidence','zunit']", 1)
				raw += "\n[checks.zunit]\nrun=['true']\ntimeout='1m'\n[checks.conditional]\nrun=['true']\ntimeout='1m'\nwhen_changed=['**/*.go']\n"
				if name == "standard-extra" {
					raw = strings.Replace(raw, "checks=['zunit','conditional']", "checks=['lower','zunit','conditional']", 1)
					raw += "\n[checks.lower]\nrun=['true']\ntimeout='1m'\n"
				}
				writeBackupTestFile(t, path, raw)
			})
			if name == "standard-extra" {
				lifecycleResult(t, f.run("--json", "check", "lower", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "lower-check"), "check.execute")
			}
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			// Two matching decisions pin the same newest-first selection as done.
			waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", map[string]string{"standard": "routine", "routine": "standard"}[lane]))
			waiver := waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", lane))["decision"].(map[string]any)
			if lane == "standard" {
				lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
			}
			facts := doneFacts(t, f, task)
			facts.PR.Body = "No rollback section"
			done, _ := lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
			completion := done["completion"].(map[string]any)
			before := startEventCount(t, f.workspace)
			for _, phase := range []string{"completed", "offline", "rebuilt"} {
				if phase == "offline" {
					hideWaivedCompletionObjects(t, f, completion)
				}
				if phase == "rebuilt" {
					if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
						t.Fatal(err)
					}
				}
				assertWaivedCompletionReceipt(t, f, task, name, phase, completion, waiver)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("show or rebuild wrote an Event")
			}
		})
	}
}

func hideWaivedCompletionObjects(t *testing.T, f *statusFixture, completion map[string]any) {
	t.Helper()
	// Today's policy and changed paths disagree with the completed hand-off.
	path := filepath.Join(f.workspace, ".fid.toml")
	raw := strings.ReplaceAll(string(mustReadCommandFile(t, path)), "checks=[]", "checks=['evidence']")
	raw = strings.ReplaceAll(raw, "checks=['zunit','conditional']", "checks=['evidence']")
	writeBackupTestFile(t, path, raw)
	writeBackupTestFile(t, filepath.Join(f.workspace, "later.go"), "package later\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-qm", "Later policy and code")
	runGitForTest(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	runGitForTest(t, f.workspace, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
	objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
	fp := completion["fingerprint"].(map[string]any)
	for _, key := range []string{"base", "head"} {
		oid := fp[key].(string)
		if err := os.Remove(filepath.Join(objects, oid[:2], oid[2:])); err != nil {
			t.Fatal(err)
		}
	}
}

func assertWaivedCompletionReceipt(t *testing.T, f *statusFixture, task, name, phase string, completion, waiver map[string]any) {
	t.Helper()
	lane := strings.TrimSuffix(name, "-extra")
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	receipt := shown["completion"].(map[string]any)
	for key, want := range completion {
		if !reflect.DeepEqual(receipt[key], want) {
			t.Errorf("%s completion %s=%v, want %v", phase, key, receipt[key], want)
		}
	}
	waivers, ok := shown["waivers"].([]any)
	if !ok || len(waivers) != 2 {
		t.Errorf("%s waivers=%v, want both matching decisions", phase, shown["waivers"])
	} else {
		for _, value := range waivers {
			if value.(map[string]any)["current"] != true {
				t.Errorf("%s matching waiver marked stale: %v", phase, value)
			}
		}
	}
	wantWaiver := map[string]any{"current": true}
	for key, value := range waiver {
		wantWaiver[key] = value
	}
	if !reflect.DeepEqual(receipt["waiver"], wantWaiver) {
		t.Errorf("%s completed waiver=%v, want %v", phase, receipt["waiver"], wantWaiver)
	}
	if !reflect.DeepEqual(shown["approval"], map[string]any{"state": "not-required"}) || shown["no_other_model_review"] != nil || receipt["approval"] != nil || len(receipt["missing"].([]any)) != 0 {
		t.Errorf("%s waived approval still required: %v / %v", phase, shown["approval"], receipt)
	}
	h := receipt["handoff"].(map[string]any)
	if shown["issue"].(map[string]any)["lane"] != "assured" || h["lane"] != "assured" || !reflect.DeepEqual(h["required_checks"], []any{"evidence", "zunit"}) {
		t.Errorf("%s waiver changed recorded lanes or hand-off requirements: %v", phase, shown)
	}
	checks := shown["checks"].([]any)
	wantCount := map[string]int{"standard": 1, "routine": 0, "standard-extra": 2}[name]
	if len(checks) != wantCount {
		t.Errorf("%s checks=%v, want %d required by waived %s", phase, checks, wantCount, lane)
	} else {
		for i, value := range checks {
			check := value.(map[string]any)
			wantName := "zunit"
			if name == "standard-extra" && i == 0 {
				wantName = "lower"
			}
			if check["name"] != wantName || check["required"] != true || check["satisfied"] != true || check["start"] == nil {
				t.Errorf("%s wrong waived check: %v", phase, check)
			} else if start := check["start"].(map[string]any); start["name"] != wantName || wantName == "zunit" && start["id"] != h["result_ids"].([]any)[1] {
				t.Errorf("%s wrong waived check evidence ID: %v", phase, check)
			}
		}
	}
	if len(shown["warnings"].([]any)) != 0 {
		t.Errorf("%s consulted current policy or git objects: %v", phase, shown["warnings"])
	}
	out := f.run("show", task)
	if out.err != nil {
		t.Fatal(out)
	}
	text := strings.Join(strings.Fields(out.stdout), " ")
	completed := strings.SplitN(text, "Completed ", 2)
	for _, want := range []string{"Lane Assured · waiver to " + lane, "Approval ○ not required · current waiver", "PR #63 merged as aaaaaaaa"} {
		if !strings.Contains(text, want) {
			t.Errorf("%s text missing %q: %s", phase, want, text)
		}
	}
	if len(completed) != 2 || !strings.Contains(completed[1], "Waiver to "+lane+" · ✓ "+waiver["decision_id"].(string)) || !strings.Contains(completed[1], "human:mark") {
		t.Errorf("%s waiver absent from Completed block: %s", phase, text)
	}
	for _, absent := range []string{"Rollback", "approval missing", "no other-model review", "HEAD", "evidence passed", "conditional"} {
		if strings.Contains(text, absent) {
			t.Errorf("%s text contains obsolete requirement %q: %s", phase, absent, text)
		}
	}
}

func TestShowCompletedWaiverCurrency(t *testing.T) {
	t.Parallel()
	for _, edited := range []bool{false, true} {
		t.Run(map[bool]string{false: "unchanged", true: "edited"}[edited], func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			waiver := waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", "standard"))["decision"].(map[string]any)
			changes, _ := lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
			finding := changes["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
			prepareSecondReviewHandoff(t, f, task)
			if edited {
				writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "An edited task\n")
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-qm", "Edit the waived change")
			}
			reviewHandoffTask(t, f, "next-handoff", task)
			for _, stream := range []string{"primary", "separate"} {
				if stream == "separate" && !edited {
					continue
				}
				in := reviewInput()
				in["stream"] = stream
				in["reviewer"].(map[string]any)["session_id"] = stream
				if stream == "primary" {
					in["dispositions"] = []any{disposition(finding, "resolved")}
				}
				lifecycleResult(t, recordReviewTask(t, f, in, stream, "human:courier", task), "review.record")
			}
			if edited {
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			done, _ := lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, task), nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
			completion := done["completion"].(map[string]any)
			if completion["handoff_id"] == waiver["binding"].(map[string]any)["handoff_id"] {
				t.Fatal("fixture did not replace the hand-off")
			}
			hideWaivedCompletionObjects(t, f, completion)
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
						t.Fatal(err)
					}
				}
				shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
				receipt := shown["completion"].(map[string]any)
				waivers, ok := shown["waivers"].([]any)
				if !ok || len(waivers) != 1 || waivers[0].(map[string]any)["current"] != !edited {
					t.Fatalf("historical waiver currency=%v, edited=%t", shown["waivers"], edited)
				}
				if (receipt["waiver"] != nil) != !edited || (receipt["approval"] != nil) != edited || len(receipt["missing"].([]any)) != 0 {
					t.Fatalf("waiver displaced the completed requirements: %v", receipt)
				}
				out := f.run("show", task)
				text := strings.Join(strings.Fields(out.stdout), " ")
				if out.err != nil || strings.Contains(text, "⚠ stale") != edited || strings.Contains(text, "Rollback present") != edited || strings.Contains(text, "Waiver to standard") != !edited || !strings.Contains(text, waiver["decision_id"].(string)) {
					t.Fatalf("wrong historical currency text: %#v", out)
				}
			}
		})
	}
}

func TestShowWaivedCompletionMissingPolicy(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "assured", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", "standard"))
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, task), nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	// Incomplete older projections must explain the missing frozen policy,
	// while retaining the available waiver and completion facts.
	for _, query := range []string{"PRAGMA foreign_keys=OFF", "DELETE FROM policy_snapshots"} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	receipt := shown["completion"].(map[string]any)
	if receipt["waiver"] == nil || !reflect.DeepEqual(receipt["missing"], []any{"policy snapshot missing · 1"}) || len(shown["checks"].([]any)) != 0 {
		t.Fatalf("missing frozen requirements invented or unexplained: %v", receipt)
	}
	out := f.run("show", task)
	if out.err != nil || !strings.Contains(strings.Join(strings.Fields(out.stdout), " "), "policy snapshot missing · 1") {
		t.Fatal(out)
	}
}
