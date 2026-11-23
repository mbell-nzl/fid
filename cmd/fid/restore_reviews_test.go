package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func recordV8Reviews(t *testing.T, root, task string) {
	t.Helper()
	approveSuccess(t, approvePublicForTest(t, root, task, task+"\n", true, true, true, "--idempotency-key", "v8-approval-before-edit"))
	var ids []string
	reviewer := gittest.AddWorktree(t, root)
	claimV8Review(t, reviewer, task)
	for i, stream := range []string{"primary", "separate"} {
		in := reviewInput()
		in["stream"] = stream
		in["verdict"] = "CHANGES"
		in["reviewer"].(map[string]any)["session_id"] = "v8-" + stream
		if i == 0 {
			in["findings"] = []any{reviewFinding("major"), reviewFinding("major"), reviewFinding("major"), reviewFinding("minor"), reviewFinding("nit")}
		}
		if i == 1 {
			in["dispositions"] = []any{disposition(ids[0], "resolved"), disposition(ids[3], "not-applicable")}
		}
		raw, err := json.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "review.json")
		if err := os.WriteFile(path, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		fields, _ := lifecycleResult(t, executeCommandForTest(root, "--json", "review", "record", task, "--input", path, "--actor", "human:courier", "--idempotency-key", "v8-"+stream), "review.record")
		if i == 0 {
			for _, finding := range fields["round"].(map[string]any)["findings"].([]any) {
				ids = append(ids, finding.(map[string]any)["id"].(string))
			}
		}
	}
	lifecycleResult(t, executeCommandForTest(root, "--json", "finding", "close", ids[1], "--actor", "human:mark", "--reason", "Accepted risk", "--idempotency-key", "v8-close"), "finding.close")
	recordV8Rework(t, root, task, ids[2])
	claimV8Review(t, reviewer, task)
	in := reviewInput()
	in["stream"] = "separate"
	in["reviewer"].(map[string]any)["session_id"] = "v8-fixed-separate"
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	lifecycleResult(t, recordReviewTask(t, f, in, "v8-fixed-separate", "human:courier", task), "review.record")
	approveSuccess(t, approvePublicForTest(t, root, task, task+"\n", true, true, true, "--idempotency-key", "v8-approval-after-edit"))
	recordV8DefaultConflict(t, root)
}

func claimV8Review(t *testing.T, root, task string) {
	t.Helper()
	got, _ := lifecycleResult(t, executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:v8-reviewer", "--model", "claude/opus-5.5"), "next.claim")
	if got["kind"] != "review" || got["task"].(map[string]any)["id"] != task {
		t.Fatalf("V8 review reservation=%v", got)
	}
	lifecycleResult(t, executeCommandForTest(root, "--json", "heartbeat", "--hold", got["hold"].(map[string]any)["id"].(string), "--actor", "agent:v8-reviewer"), "hold.heartbeat")
}

func recordV8Rework(t *testing.T, root, task, finding string) {
	t.Helper()
	taken, _ := lifecycleResult(t, executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:fixer", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if taken["kind"] != "rework" || taken["task"].(map[string]any)["id"] != task || taken["hold"].(map[string]any)["replaces"] == nil {
		t.Fatalf("V8 needs public rework: %v", taken)
	}
	hold := taken["hold"].(map[string]any)["id"].(string)
	writeBackupTestFile(t, filepath.Join(root, "rework.txt"), "Fix the major review finding\n")
	runGitForTest(t, root, "add", ".")
	runGitForTest(t, root, "commit", "-qm", "Rework the reviewed change")
	lifecycleResult(t, executeCommandForTest(root, "--json", "check", "--hold", hold, "--actor", "agent:fixer", "--idempotency-key", "v8-rework-check"), "check.execute")
	lifecycleResult(t, executeCommandForTest(root, "--json", "park", "--hold", hold, "--actor", "agent:fixer", "--waiting", "Does this need an answer before finishing?", "--idempotency-key", "v8-finished-question"), "hold.park")
	finishTakeover, _ := lifecycleResult(t, executeCommandForTest(root, "--json", "takeover", task, "--actor", "human:mark", "--reason", "Finish before an answer", "--model", "codex/gpt-6.1-sol"), "task.takeover")
	hold = finishTakeover["hold_id"].(string)
	lifecycleResult(t, executeCommandForTest(root, "--json", "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "v8-rework-handoff"), "task.handoff")
	in := reviewInput()
	in["reviewer"].(map[string]any)["session_id"] = "v8-fixed-primary"
	in["dispositions"] = []any{disposition(finding, "resolved")}
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	lifecycleResult(t, recordReviewTask(t, f, in, "v8-fixed-review", "human:courier", task), "review.record")
	// Preserve V8's open-major coverage on a separate task. An approving
	// round cannot coexist with an unclosed major on its own task.
	created, _ := lifecycleResult(t, executeCommandForTest(root, "--json", "issue", "create", "--kind", "task", "--title", "Open review finding at backup", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "v8-open-task"), "issue.create")
	openTask := created["issue"].(map[string]any)["id"].(string)
	other := gittest.AddWorktree(t, root)
	took, _ := lifecycleResult(t, executeCommandForTest(other, "--json", "takeover", openTask, "--actor", "human:mark", "--reason", "Build the second review subject"), "task.takeover")
	hold = took["hold_id"].(string)
	lifecycleResult(t, executeCommandForTest(other, "--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "v8-open-start"), "task.start")
	lifecycleResult(t, executeCommandForTest(other, "--json", "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/64", "--idempotency-key", "v8-open-handoff"), "task.handoff")
	f.workspace = other
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "v8-open-review", "human:courier", openTask), "review.record")
}

func assertV8Reviews(t *testing.T, workspace string) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var conflicts int
	if err := conn.QueryRow(`SELECT count(*) FROM default_conflicts`).Scan(&conflicts); err != nil || conflicts != 1 {
		t.Fatalf("V8 default-conflict receipts=%d/%v", conflicts, err)
	}
	var rounds, findings, models, manual, closures int
	var claims int
	if err := conn.QueryRow(`SELECT count(*) FROM review_claims WHERE ended_at IS NOT NULL`).Scan(&claims); err != nil || claims != 2 {
		t.Fatalf("V8 review reservations=%d/%v", claims, err)
	}
	for _, q := range []struct {
		sql string
		to  *int
	}{
		{`SELECT count(*) FROM finding_closures`, &closures},
		{`SELECT count(*) FROM review_rounds`, &rounds}, {`SELECT count(*) FROM round_findings`, &findings},
		{`SELECT count(*) FROM holds WHERE model<>''`, &models}, {`SELECT count(*) FROM holds WHERE model=''`, &manual},
	} {
		if err := conn.QueryRow(q.sql).Scan(q.to); err != nil {
			t.Fatal(err)
		}
	}
	if rounds != 6 || findings != 6 || closures != 4 || models == 0 || manual == 0 {
		t.Fatalf("V8 reviews=%d findings=%d modeled holds=%d manual holds=%d closures=%d", rounds, findings, models, manual, closures)
	}
	for _, severity := range []string{"major", "minor", "nit"} {
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM round_findings WHERE json_extract(record_json,'$.severity')=?`, severity).Scan(&n); err != nil || n != map[string]int{"major": 4, "minor": 1, "nit": 1}[severity] {
			t.Fatalf("V8 %s=%d/%v", severity, n, err)
		}
	}
	var task string
	if err := conn.QueryRow(`SELECT task_id FROM review_rounds WHERE json_extract(record_json,'$.reviewer.session_id')='v8-primary'`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	open, err := reviewrecord.OpenMajors(conn, task)
	if err != nil || len(open) != 0 {
		t.Fatalf("V8 open majors=%v/%v", open, err)
	}
	views, err := reviewwork.List(workspace, task)
	if err != nil || len(views) != 4 || views[0].Current || views[1].Current || !views[2].Current || views[2].Verdict != "APPROVE" || !views[3].Current || views[3].Verdict != "APPROVE" {
		t.Fatalf("V8 stale and fixed rounds=%v/%v", views, err)
	}
	var outstanding int
	if err := conn.QueryRow(`SELECT count(*) FROM round_findings f LEFT JOIN finding_closures c ON c.finding_id=f.id WHERE json_extract(f.record_json,'$.severity')='major' AND c.finding_id IS NULL`).Scan(&outstanding); err != nil || outstanding != 1 {
		t.Fatalf("V8 must retain an open major: %d/%v", outstanding, err)
	}
	for _, stream := range []string{"primary", "separate"} {
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM review_rounds WHERE json_extract(record_json,'$.stream')=?`, stream).Scan(&n); err != nil || n != map[string]int{"primary": 4, "separate": 2}[stream] {
			t.Fatalf("V8 %s=%d/%v", stream, n, err)
		}
	}
}

func recordV8DefaultConflict(t *testing.T, root string) {
	t.Helper()
	builder := gittest.AddWorktree(t, root)
	f := &statusFixture{t: t, workspace: builder, clock: time.Now}
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Approved default conflict at backup", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "v8-conflict-create"), "issue.create")
	task := created["issue"].(map[string]any)["id"].(string)
	took, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Build the update subject"), "task.takeover")
	hold := took["hold_id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "v8-conflict-start"), "task.start")
	writeBackupTestFile(t, filepath.Join(builder, "notes.md"), "Approved content\n")
	runGitForTest(t, builder, "add", ".")
	runGitForTest(t, builder, "commit", "-qm", "Approved conflict subject")
	reviewHandoffTask(t, f, "v8-conflict-handoff", task)
	approveConflictReviews(t, f, task, "standard")
	moveConflictDefault(t, f, true)
	updated, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	if updated["kind"] != "rework" || updated["rework"].(map[string]any)["cause"] != "default-conflict" {
		t.Fatalf("V8 conflict return=%v", updated)
	}
	hold = updated["hold"].(map[string]any)["id"].(string)
	lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", "human:mark", "--interrupted", "--reason", "Keep conflict update across restore", "--idempotency-key", "v8-conflict-park"), "hold.park")
}
