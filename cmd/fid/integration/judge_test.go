//go:build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestLoopShadowJudge(t *testing.T) {
	t.Parallel()
	for _, tool := range []string{"codex", "claude"} {
		t.Run(tool, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "judge")
			j.env = append(j.env, "LOOP_JUDGE_USAGE="+tool)
			if tool == "claude" {
				j.env = append(j.env, "LOOP_JUDGE_MODEL=claude/opus-4-6", "LOOP_JUDGE_ACTOR=agent:opus-judge", "FIRST_LOOP_JUDGE_CLAUDE=1")
			}
			j.run(t, 0)
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			c := &loopContract{t: t, decl: readLoopDeclaration(t)}
			sealed := w.call("review", "list", j.task)
			c.observe("judge-sealed-"+tool, sealed, 0)
			raw, _ := json.Marshal(sealed)
			if strings.Contains(string(raw), "judge sentinel") || !strings.Contains(string(raw), `"sealed":true`) {
				t.Fatalf("unsealed: %s", raw)
			}
			assertJudgeArchive(t, j, sealed, "sealed")
			var judges, comments int
			for _, call := range j.calls(t) {
				if call["kind"] == "judge" {
					judges++
					context := call["context"].(map[string]any)
					if context["receipts"] == nil || !strings.Contains(call["prompt"].(string), "whole diff") {
						t.Fatalf("missing judge evidence: %v", call)
					}
				}
				if call["kind"] == "comment" {
					comments++
					if strings.Contains(call["body"].(string), "judge sentinel") {
						t.Fatal("posted sealed judge")
					}
				}
			}
			if judges != 1 || comments != 1 {
				t.Fatalf("judges %d comments %d", judges, comments)
			}
			heldRun := backupJourneyFID(t, j.env, j.root, "--json", "hold", j.task, "--actor", "human:mark", "--reason", "Wait for a closer look", "--idempotency-key", "mark-call")
			heldRun.expectExit(t, 0)
			held := p23Envelope(t, heldRun.stdout.String())
			c.observe("judge-hold-"+tool, held, 0)
			revealed := w.call("review", "list", j.task)
			c.observe("judge-revealed-"+tool, revealed, 0)
			raw, _ = json.Marshal(revealed)
			if !strings.Contains(string(raw), "judge sentinel") || !strings.Contains(string(raw), `"agrees":true`) {
				t.Fatalf("not revealed: %s", raw)
			}
			assertJudgeArchive(t, j, revealed, "revealed")
			shown := w.call("show", j.task)
			byRole := p21Path(t, shown, "result", "usage", "by_role").(map[string]any)
			if byRole["judge.judge"] == nil {
				t.Fatalf("no judge usage: %v", shown)
			}
			usage, _ := closureRunFID(t, j.binary, j.root, "sql", "SELECT role,stream,format,model FROM usage_sessions WHERE role='judge'")
			raw = []byte(usage)
			if !strings.Contains(string(raw), tool) {
				t.Fatalf("wrong judge parser: %s", raw)
			}
			posted := 0
			for _, call := range j.calls(t) {
				if call["kind"] == "comment" && strings.Contains(call["body"].(string), "judge sentinel") {
					posted++
					body := call["body"].(string)
					if !strings.Contains(body, "Assess merge readiness") || !strings.Contains(body, "**Lane** `standard`") {
						t.Fatalf("judge comment lost purpose or lane: %s", body)
					}
				}
			}
			if posted != 1 {
				t.Fatalf("revealed judge posted %d times", posted)
			}
		})
	}
}

func assertJudgeArchive(t *testing.T, j firstLoop, want map[string]any, phase string) {
	t.Helper()
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	w.call("rebuild")
	if got := w.call("review", "list", j.task); !reflect.DeepEqual(got["result"], want["result"]) {
		t.Fatalf("%s rebuild changed sealing: %v / %v", phase, got, want)
	}
	archive := filepath.Join(t.TempDir(), "judge.zip")
	w.call("export", "create", archive)
	w.call("export", "validate", archive)
	target := gittest.NewWorkspace(t)
	p21Run(t, j.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "judge-import")
	got := p21Run(t, j.binary, target, "--json", "review", "list", j.task)
	if !reflect.DeepEqual(got["result"], want["result"]) {
		t.Fatalf("%s import changed sealing: %v / %v", phase, got, want)
	}
}

func TestLoopJudgeSealedInReworkContext(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "judge")
	j.run(t, 0)
	w := loopWorker{t, j.binary, j.root, "agent:sol-reviewer"}
	loopRound(w, j.task, "CHANGES", "after-judge", []string{"code", "test-adequacy"})
	j.run(t, 0)
	found := false
	for _, call := range j.calls(t) {
		if call["kind"] != "build" {
			continue
		}
		ctx := call["context"].(map[string]any)
		if ctx["work"].(map[string]any)["kind"] != "rework" {
			continue
		}
		found = true
		raw, _ := json.Marshal(ctx)
		if strings.Contains(string(raw), "judge sentinel") || !strings.Contains(string(raw), `"sealed":true`) {
			t.Fatalf("loop context leaked: %s", raw)
		}
	}
	if !found {
		t.Fatal("no rework context inspected")
	}
}

var p24HoldFields = map[string]string{"task.hold": "--reason"}

func TestHoldKeyedContract(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Hold retry", "")
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	c.failure("hold-not-in-review", binary, root, "hold", task, "--actor", "human:mark", "--reason", "Wait")
	w.build("", task, "hold-retry")
	c.failure("hold-denied", binary, root, "hold", task, "--actor", "agent:builder", "--reason", "Wait")
	c.failure("hold-reason-required", binary, root, "hold", task, "--actor", "human:mark", "--reason", " ")
	p24StartKeyedCapture(root)
	w.call("hold", task, "--actor", "human:mark", "--reason", "Wait for evidence", "--idempotency-key", "hold-contract")
	captured := p24EndKeyedCapture(root)
	entry, ok := captured["task.hold"]
	if !ok {
		t.Fatal("hold command was not captured")
	}
	p24ExerciseCapturedLeaf(t, root, binary, "task.hold", "--reason", entry)
}

func TestLoopRejectsJudgeReviewStream(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "judge")
	const output = `    Path(os.environ["LOOP_REVIEW"]).write_text(json.dumps(review, indent=2))`
	firstLoopReplaceStub(t, j, "agent", output, `    review.update(stream="primary", profiles=["code", "test-adequacy"], verdict="CHANGES", findings=[{"profile":"code","severity":"major","summary":"Wrong role","location":None,"evidence":"Judge returned a review","fix":"Stop before recording"}])
    review.pop("reason")
`+output)
	out := j.run(t, 1)
	if !strings.Contains(out, "review output stream differs from requested work") {
		t.Fatalf("missing role refusal: %s", out)
	}
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	if state := p21PathString(t, w.call("show", j.task), "result", "issue", "status"); state != "in_review" {
		t.Fatalf("judge created rework: %s", state)
	}
	rounds := p21Path(t, w.call("review", "list", j.task), "result", "rounds").([]any)
	if len(rounds) != 1 {
		t.Fatalf("misrouted judge was recorded: %v", rounds)
	}
	builds := 0
	for _, call := range j.calls(t) {
		if call["kind"] == "build" {
			builds++
		}
	}
	if builds != 1 {
		t.Fatalf("judge triggered %d builds", builds)
	}
}
