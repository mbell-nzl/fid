package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func completedShowFixture(t *testing.T, lane, model string) (*statusFixture, string, map[string]any, map[string]any) {
	t.Helper()
	f, task, hold := doneBuildFixture(t, lane, false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	for _, stream := range []string{"primary", "separate"} {
		if stream == "separate" && lane != "assured" {
			continue
		}
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["model"] = model
		in["reviewer"].(map[string]any)["session_id"] = stream
		lifecycleResult(t, recordReviewTask(t, f, in, stream, "human:courier", task), "review.record")
	}
	var approval map[string]any
	if lane == "assured" {
		approval = approveSuccess(t, approveForTest(f, task, task+"\n", true))["decision"].(map[string]any)
	}
	fields, _ := lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, task), nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	return f, task, fields["completion"].(map[string]any), approval
}

func TestShowCompletionReceipt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		lane, model string
		flag        bool
	}{
		{"standard", "claude/opus-5.5", false},
		{"assured", "claude/opus-5.5", false},
		{"assured", "codex/gpt-6.1-sol", true},
	} {
		t.Run(tc.lane+tc.model, func(t *testing.T) {
			t.Parallel()
			f, task, completion, approval := completedShowFixture(t, tc.lane, tc.model)
			// Move HEAD and make current policy/comparison unavailable. Finished work
			// must still describe its historical passes and approval, offline.
			writeBackupTestFile(t, filepath.Join(f.workspace, "later.txt"), "Later work\n")
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-qm", "Later work")
			runGitForTest(t, f.workspace, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			head := completion["fingerprint"].(map[string]any)["head"].(string)
			objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
			if err := os.Remove(filepath.Join(objects, head[:2], head[2:])); err != nil {
				t.Fatal(err)
			}

			before := startEventCount(t, f.workspace)
			shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
			receipt, ok := shown["completion"].(map[string]any)
			if !ok {
				t.Fatalf("completion receipt missing: %v", shown)
			}
			for key, want := range completion {
				if !reflect.DeepEqual(receipt[key], want) {
					t.Fatalf("receipt %s=%v, want %v", key, receipt[key], want)
				}
			}
			merge := receipt["merge_observation"].(map[string]any)
			handoff := receipt["handoff"].(map[string]any)
			if merge["number"] != float64(63) || merge["merge_commit"] != strings.Repeat("a", 40) || handoff["id"] != completion["handoff_id"] {
				t.Fatalf("wrong completed change: %v", receipt)
			}
			ci := receipt["ci_results"].([]any)
			if len(ci) != 1 || ci[0].(map[string]any)["conclusion"] != "success" || ci[0].(map[string]any)["head_sha"] != merge["merge_commit"] || ci[0].(map[string]any)["url"] != "https://github.com/mbell-nzl/fid/actions/runs/101" {
				t.Fatalf("wrong CI: %v", ci)
			}
			rounds := receipt["rounds"].([]any)
			wantRounds := 1
			if tc.lane == "assured" {
				wantRounds = 2
			}
			if len(rounds) != wantRounds || rounds[0].(map[string]any)["reviewer"].(map[string]any)["model"] != tc.model {
				t.Fatalf("wrong rounds: %v", rounds)
			}
			if (receipt["approval"] != nil) != (tc.lane == "assured") || receipt["no_other_model_review"] != tc.flag {
				t.Fatalf("approval/flag: %v", receipt)
			}
			// Pin the top-level historical view to the decision returned by approve,
			// even when its git objects and today's comparison are unavailable.
			var wantApproval, wantDisclosure any
			if tc.lane == "assured" {
				if approval["decision_id"] != completion["approval_id"] {
					t.Fatalf("completion names approval %v, want %v", completion["approval_id"], approval["decision_id"])
				}
				wantApproval = map[string]any{
					"state":       "current",
					"decision_id": approval["decision_id"],
					"actor":       approval["actor"],
					"at":          approval["at"],
					"fingerprint": approval["binding"].(map[string]any)["fingerprint"],
				}
				wantDisclosure = tc.flag
			}
			if got, ok := shown["approval"]; !ok || !reflect.DeepEqual(got, wantApproval) {
				t.Fatalf("historical top-level approval=%v (present=%t), want %v", got, ok, wantApproval)
			}
			if got, ok := shown["no_other_model_review"]; !ok || got != wantDisclosure {
				t.Fatalf("historical top-level disclosure=%v (present=%t), want %v", got, ok, wantDisclosure)
			}
			checks := shown["checks"].([]any)
			if len(checks) != 1 || checks[0].(map[string]any)["satisfied"] != true || checks[0].(map[string]any)["start"].(map[string]any)["id"] != handoff["result_ids"].([]any)[0] {
				t.Fatalf("not the completed checks: %v", checks)
			}
			if len(shown["warnings"].([]any)) != 0 {
				t.Fatalf("finished task compared current evidence: %v", shown["warnings"])
			}
			out := f.run("show", task)
			if out.err != nil {
				t.Fatal(out)
			}
			text := strings.Join(strings.Fields(out.stdout), " ")
			for _, want := range []string{"Completed", completion["id"].(string), "agent:finisher", "PR #63 merged as aaaaaaaa", "CI ✓ success", "fingerprint", "primary APPROVE", tc.model, "passed at", "counts at"} {
				if !strings.Contains(text, want) {
					t.Fatalf("missing %q in %s", want, text)
				}
			}
			if strings.Contains(text, "HEAD") || strings.Contains(text, "no other-model review") != tc.flag || strings.Contains(text, "Rollback present") != (tc.lane == "assured") {
				t.Fatalf("wrong historical text: %s", text)
			}
			if before != startEventCount(t, f.workspace) {
				t.Fatal("show wrote an Event")
			}
		})
	}
}

func TestShowCompletionMissingRecords(t *testing.T) {
	t.Parallel()
	for _, missing := range []string{"completion", "hand-off", "merge", "CI result", "review round", "approval", "check result"} {
		t.Run(missing, func(t *testing.T) {
			t.Parallel()
			f, task, _, _ := completedShowFixture(t, "assured", "codex/gpt-6.1-sol")
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			// Model old/incomplete projections without manufacturing journal history.
			queries := map[string]string{"completion": "DELETE FROM completions", "hand-off": "DELETE FROM handoffs", "merge": "DELETE FROM merge_observations", "CI result": "DELETE FROM ci_results", "review round": "DELETE FROM review_rounds", "approval": "DELETE FROM person_decisions WHERE kind='approve'", "check result": "DELETE FROM check_results"}
			if _, err := conn.Exec("PRAGMA foreign_keys=OFF"); err != nil {
				t.Fatal(err)
			}
			if _, err := conn.Exec(queries[missing]); err != nil {
				t.Fatal(err)
			}
			_ = conn.Close()
			shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
			if missing == "check result" {
				check := shown["checks"].([]any)[0].(map[string]any)
				if check["satisfied"] != false {
					t.Fatalf("missing check counted: %v", check)
				}
			}
			raw, _ := json.Marshal(shown["completion"])
			out := f.run("show", task)
			if missing == "completion" {
				missing = "completion record"
			}
			if out.err != nil || !strings.Contains(string(raw), "missing") || !strings.Contains(strings.Join(strings.Fields(out.stdout), " "), missing+" missing") {
				t.Fatalf("missing %s not explained: %s / %#v", missing, raw, out)
			}
			if missing != "completion record" && !strings.Contains(out.stdout, "agent:finisher") {
				t.Fatal("lost available completion facts")
			}
		})
	}
}

func TestShowCompletionNullOutsideDoneTasks(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, 0)
	showWork(f)
	for _, id := range []string{"FID-71a0c2", "FID-4c5d6e", "FID-e3a1c0", "FID-aa01b2"} {
		shown, _ := lifecycleResult(t, f.run("--json", "show", id), "issue.show")
		if value, ok := shown["completion"]; !ok || value != nil {
			t.Fatalf("%s completion=%v present=%t", id, value, ok)
		}
	}
	// A legacy done task has no completion row, and still renders successfully.
	out := f.run("show", "FID-1a2b3c")
	if out.err != nil || !strings.Contains(out.stdout, "completion record missing") {
		t.Fatal(out)
	}
}

func TestShowCompletionSelectsNamedEvidence(t *testing.T) {
	t.Parallel()
	f, task, completion, _ := completedShowFixture(t, "assured", "claude/opus-5.5")
	before, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	// Later unrelated observations must not replace the evidence named by done.
	queries := []string{
		`PRAGMA foreign_keys=OFF`,
		`INSERT INTO handoffs SELECT 'HOF-decoy', task_id, 'HLD-decoy', policy_snapshot, json_set(record_json,'$.id','HOF-decoy','$.head','ffffffffffffffffffffffffffffffffffffffff') FROM handoffs`,
		`INSERT INTO merge_observations SELECT 'MRG-decoy',task_id,'HOF-decoy','ffffffffffffffffffffffffffffffffffffffff',json_set(record_json,'$.id','MRG-decoy','$.number',999) FROM merge_observations`,
		`INSERT INTO ci_results SELECT 'CIR-decoy',task_id,merge_id,repository,999,attempt,json_set(record_json,'$.id','CIR-decoy','$.run_id',999,'$.url','https://github.com/mbell-nzl/fid/actions/runs/999') FROM ci_results`,
		`INSERT INTO review_rounds SELECT 'RND-decoy',task_id,handoff_id,json_set(record_json,'$.id','RND-decoy','$.reviewer.actor','agent:decoy') FROM review_rounds LIMIT 1`,
		`INSERT INTO person_decisions SELECT 'PDN-decoy',kind,'human:decoy',subject_type,subject_id,reason,binding_id,decided_at,binding_json FROM person_decisions WHERE kind='approve'`,
		`INSERT INTO check_starts SELECT 'CHK-decoy','REQ-decoy',task_id,name,json_set(start_json,'$.id','CHK-decoy') FROM check_starts`,
		`INSERT INTO check_results SELECT 'CHK-decoy',json_set(result_json,'$.outcome','failed','$.counts',json('false')) FROM check_results`,
	}
	for _, query := range queries {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	after, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	if !reflect.DeepEqual(before["completion"], after["completion"]) || !reflect.DeepEqual(before["checks"], after["checks"]) || after["handoff"].(map[string]any)["id"] != completion["handoff_id"] {
		t.Fatalf("later records replaced completed evidence: %v", after)
	}
	if _, err := conn.Exec("DELETE FROM handoffs WHERE id=?", completion["handoff_id"]); err != nil {
		t.Fatal(err)
	}
	missing, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	if missing["handoff"] != nil {
		t.Fatalf("missing completed hand-off replaced by unrelated latest hand-off: %v", missing["handoff"])
	}

}

func TestShowCompletionLegacyFacts(t *testing.T) {
	t.Parallel()
	f, task, _, _ := completedShowFixture(t, "assured", "codex/gpt-6.1-sol")
	for _, query := range []string{
		`UPDATE completions SET record_json=json_set(record_json,'$.fingerprint',json('{}'),'$.rollback_present',json('false'))`,
		`UPDATE ci_results SET record_json=json_set(record_json,'$.url','')`,
		`UPDATE review_rounds SET record_json=json_remove(record_json,'$.reviewer.model')`,
	} {
		startSQL(t, f, query)
	}
	out := f.run("show", task)
	text := strings.Join(strings.Fields(out.stdout), " ")
	for _, want := range []string{"fingerprint missing", "model not recorded", "Rollback absent", "CI ✓ success"} {
		if out.err != nil || !strings.Contains(text, want) {
			t.Fatalf("missing %q: %#v", want, out)
		}
	}
	if strings.Contains(text, "actions/runs") {
		t.Fatal("invented an unrecorded CI URL")
	}
}

func TestShowCompletionAllCountedCIRuns(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "standard", false)
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	facts := doneFacts(t, f, task)
	for i, conclusion := range []string{"success", "skipped"} {
		run := facts.CI.Runs[0]
		run.ID, run.WorkflowID = int64(202+i), int64(2+i)
		run.Name, run.Path = "other", ".github/workflows/other.yml"
		run.Conclusion = conclusion
		run.URL = "https://github.com/mbell-nzl/fid/actions/runs/" + fmt.Sprint(run.ID)
		facts.CI.Runs = append(facts.CI.Runs, run)
	}
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	receipt := shown["completion"].(map[string]any)
	runs := receipt["ci_results"].([]any)
	if len(runs) != 2 || runs[0].(map[string]any)["run_id"] != float64(101) || runs[1].(map[string]any)["run_id"] != float64(202) {
		t.Fatalf("not all counted runs: %v", runs)
	}
	out := f.run("show", task)
	text := strings.Join(strings.Fields(out.stdout), " ")
	if out.err != nil || !strings.Contains(text, "#101") || !strings.Contains(text, "#202") || strings.Contains(text, "#203") {
		t.Fatal(out)
	}
}
