//go:build integration

package integration

import (
	"strconv"
	"strings"
	"testing"
)

func TestFirstLoopConventionalPRTitle(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, commits, subject string
		noBrief                bool
	}{
		{"first-feat", "    subprocess.run(['git', 'commit', '--allow-empty', '-qm', 'feat(loop)!: readable reviews'], check=True)\n" +
			"    subprocess.run(['git', 'commit', '--allow-empty', '-qm', 'feat: later feature'], check=True)\n", "feat(loop)!: readable reviews", false},
		{"no-feat", "    subprocess.run(['git', 'commit', '--allow-empty', '-qm', 'fix: later repair'], check=True)\n", "test: stub implementation", false},
		{"no-brief", "", "test: stub implementation", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			scenario := "approve"
			if tt.noBrief {
				scenario = "no-brief"
			}
			j := newFirstLoop(t, scenario)
			stem := "test"
			if tt.noBrief {
				stem = j.task
			}
			if tt.commits != "" {
				firstLoopReplaceStub(t, j, "agent", `    Path(os.environ["LOOP_REPORT"]).write_text`, tt.commits+`    Path(os.environ["LOOP_REPORT"]).write_text`)
			}
			j.run(t, 0)
			var title string
			for _, call := range j.calls(t) {
				if call["kind"] == "create" {
					title = call["title"].(string)
				}
			}
			if want := tt.subject + " (RESET " + stem + ")"; title != want {
				t.Fatalf("PR title = %q, want %q", title, want)
			}
		})
	}
}

func TestFirstLoopBranchDropsBriefStem(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "approve")
	p22Write(t, j.root, "docs/briefs/2.6.md", "Approved numeric brief.\n")
	closureRun(t, j.root, "git", "add", ".")
	closureRun(t, j.root, "git", "commit", "-qm", "Numeric brief")
	closureRun(t, j.root, "git", "push", "-q", "origin", "HEAD:trunk")
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	w.call("issue", "update", j.task, "--title", "2.6 Usage capture", "--spec-md", "docs/briefs/2.6.md", "--actor", "human:planner", "--idempotency-key", "numeric-title")
	j.run(t, 0)
	branch, _ := closureRun(t, j.root, "git", "branch", "--show-current")
	if strings.TrimSpace(branch) != "reset-2.6-usage-capture" {
		t.Fatalf("task branch repeats brief stem: %s", branch)
	}
}

func TestFirstLoopFindingAnswers(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"delivered", "restart", "pending"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "answers")
			if scenario == "restart" {
				firstLoopReplaceStub(t, j, "gh", `if Path(sys.argv[0]).name == "gh":`, `if Path(sys.argv[0]).name == "gh":
    if args[:2] == ["pr", "edit"] and not (state / "edit-failed").exists():
        (state / "edit-failed").touch()
        print("temporary edit failure after answers", file=sys.stderr)
        sys.exit(1)`)
				if out := j.run(t, 1); !strings.Contains(out, "temporary edit failure after answers") {
					t.Fatal(out)
				}
			} else if scenario == "pending" {
				firstLoopReplaceStub(t, j, "gh", `        log(kind="comment",`, `        if "## Review:" not in Path(args[args.index("--body-file") + 1]).read_text():
            print("offline answer delivery", file=sys.stderr)
            sys.exit(1)
        log(kind="comment",`)
			}
			out := j.run(t, 0)
			if scenario == "pending" && !strings.Contains(out, "review-post-pending") {
				t.Fatalf("pending response delivery was not printed: %s", out)
			}
			firstLoopAssertAnswers(t, j, scenario != "pending")
		})
	}
}

func firstLoopAssertAnswers(t *testing.T, j firstLoop, delivered bool) {
	t.Helper()
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	rounds := p21Path(t, w.call("review", "list", j.task), "result", "rounds").([]any)
	findings := rounds[0].(map[string]any)["findings"].([]any)
	if len(rounds) != 2 || len(findings) != 2 {
		t.Fatalf("expected CHANGES and APPROVE with two findings: %v", rounds)
	}
	for _, raw := range findings {
		finding := raw.(map[string]any)
		responses, ok := finding["responses"].([]any)
		if !ok || len(responses) != 1 || responses[0].(map[string]any)["actor"] != w.actor {
			t.Fatalf("finding lacks exactly one builder response: %v", finding)
		}
		response := responses[0].(map[string]any)
		if finding["severity"] == "major" && !strings.Contains(response["text"].(string), "`change_test.go`") {
			t.Fatalf("answer lost Markdown: %v", response)
		}
		if finding["severity"] == "minor" {
			deferred := response["deferred_task"].(map[string]any)
			if deferred["title"] != "Document the retry example" {
				t.Fatal(deferred)
			}
			shown := w.call("show", deferred["id"].(string))
			if p21PathString(t, shown, "result", "issue", "status") != "todo" {
				t.Fatalf("deferral did not create work awaiting acceptance: %v", shown)
			}
			loopQueueReason(t, w.call("next", "--explain"), deferred["id"].(string), "not-accepted")
		}
		var posts int
		var bodies []string
		for _, call := range j.calls(t) {
			if call["kind"] == "comment" && (strings.Contains(call["body"].(string), "## Finding answer:") || strings.Contains(call["body"].(string), "## Deferred finding:")) && strings.Contains(call["body"].(string), finding["id"].(string)) {
				posts++
				bodies = append(bodies, call["body"].(string))
			}
		}
		want := 0
		if delivered {
			want = 1
		}
		if posts != want {
			t.Fatalf("response comments for %s = %d, want %d: %v", finding["id"], posts, want, bodies)
		}
	}
}

func TestFirstLoopMalformedAnswers(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, expression string }{
		{"json", strconv.Quote("{bad-json")},
		{"object", strconv.Quote(`{}`)},
		{"entry", strconv.Quote(`[{}]`)},
		{"unknown", strconv.Quote(`[{"finding_id":"unknown","answer":"Fixed"}]`)},
		{"partial-invalid", `json.dumps(answers + [{}])`},
		{"duplicate", `json.dumps([answers[0], answers[0]])`},
		{"both", `json.dumps([dict(answers[0], defer="Later")])`},
		{"blank", `json.dumps([dict(answers[0], answer=" ")])`},
		{"major-defer", `json.dumps([{"finding_id":answers[0]["finding_id"],"defer":"Later"}])`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "answers")
			firstLoopReplaceStub(t, j, "agent", `write_text(json.dumps(answers))`, `write_text(`+tt.expression+`)`)
			if out := j.run(t, 1); !strings.Contains(out, "answers") {
				t.Fatalf("missing answers diagnostic: %s", out)
			}
			for _, call := range j.calls(t) {
				if call["kind"] == "edit" || call["kind"] == "ready" {
					t.Fatalf("malformed answers reached hand-off: %v", call)
				}
			}
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			for _, round := range p21Path(t, w.call("review", "list", j.task), "result", "rounds").([]any) {
				for _, finding := range round.(map[string]any)["findings"].([]any) {
					if finding.(map[string]any)["responses"] != nil {
						t.Fatalf("malformed file recorded partial responses: %v", finding)
					}
				}
			}
		})
	}
}

func TestLoopFindingResponseContract(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "answers")
	// Stop after the first recorded CHANGES round, before a fixer can answer it.
	firstLoopReplaceStub(t, j, "agent", `    merging = subprocess.run`, `    if context["reviews"]["rounds"]:
        sys.exit(42)
    merging = subprocess.run`)
	j.run(t, 1)
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	rounds := p21Path(t, w.call("review", "list", j.task), "result", "rounds").([]any)
	findings := rounds[0].(map[string]any)["findings"].([]any)
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	for i, verb := range []string{"answer", "defer"} {
		flag := "--text"
		if verb == "defer" {
			flag = "--title"
		}
		id := findings[i].(map[string]any)["id"].(string)
		c.observe("finding-"+verb, w.call("finding", verb, id, "--actor", w.actor, flag, "Follow up the retry"), 0)
	}
	c.observe("review-list-responses", w.call("review", "list", j.task), 0)
	c.observe("status-pending-responses", w.call("status"), 0)
}
