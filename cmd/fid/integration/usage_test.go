//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

var p24UsageFields = map[string]string{"usage.record": "--tool-version"}

func TestUsageReplayExportImport(t *testing.T) {
	t.Parallel()
	j := newFirstLoop(t, "changes")
	j.run(t, 0)
	w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
	snapshot := func(root string) string {
		t.Helper()
		out, _ := closureRunFID(t, j.binary, root, "sql", "SELECT * FROM usage_sessions ORDER BY id")
		return out
	}
	before := snapshot(j.root)
	w.call("rebuild")
	if snapshot(j.root) != before {
		t.Fatal("rebuild changed usage")
	}
	archive := filepath.Join(t.TempDir(), "usage.zip")
	w.call("export", "create", archive)
	w.call("export", "validate", archive)
	target := gittest.NewWorkspace(t)
	p21Run(t, j.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "usage-import")
	if snapshot(target) != before {
		t.Fatal("import changed usage")
	}
	source := w.call("show", j.task)["result"].(map[string]any)["usage"]
	restored := p21Path(t, p21Run(t, j.binary, target, "--json", "show", j.task), "result", "usage")
	if !reflect.DeepEqual(source, restored) {
		t.Fatalf("usage summary changed: %v / %v", source, restored)
	}
}

func TestLoopUsageAfterFailedAndWaitingBuild(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"interrupt", "waiting"} {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, scenario)
			code := 1
			if scenario == "waiting" {
				code = 0
			}
			j.run(t, code)
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			shown := w.call("show", j.task)
			usage := p21Path(t, shown, "result", "usage").(map[string]any)
			if usage["sessions"] != float64(1) || usage["total"] != float64(15752) {
				t.Fatal(usage)
			}
		})
	}
}

func TestLoopUsageNoneAndRefusalContinue(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"none", "claude"} {
		t.Run(format, func(t *testing.T) {
			t.Parallel()
			j := newFirstLoop(t, "approve")
			j.env = append(j.env, "LOOP_BUILDER_USAGE="+format, "LOOP_REVIEWER_USAGE="+format)
			out := j.run(t, 0)
			if format == "claude" && !strings.Contains(out, "usage record failed:") {
				t.Fatal("usage refusal was not printed")
			}
			w := loopWorker{t, j.binary, j.root, "agent:sol-builder"}
			shown := w.call("show", j.task)
			if p21Path(t, shown, "result").(map[string]any)["usage"] != nil {
				t.Fatal("recorded disabled or refused usage")
			}
			if len(w.call("review", "list", j.task)["result"].(map[string]any)["rounds"].([]any)) != 1 {
				t.Fatal("usage failure stopped review")
			}
		})
	}
}

func TestUsageLoopContract(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Usage contract", "")
	work := w.call("next", "--claim", "--actor", w.actor, "--model", "openai/sol")
	hold := p21PathString(t, work, "result", "hold", "id")
	path := filepath.Join(t.TempDir(), "stdout.jsonl")
	raw := readP21File(t, "../../../internal/usagerecord/testdata/codex.jsonl")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	args := []string{"usage", "record", "--hold", hold, "--actor", w.actor, "--tool", "codex", "--tool-output", path, "--tool-version", "codex-cli 0.172.0", "--idempotency-key", "contract-usage"}
	p24StartKeyedCapture(root)
	reply := w.call(args...)
	captured := p24EndKeyedCapture(root)
	for operation, field := range p24UsageFields {
		entry, ok := captured[operation]
		if !ok {
			t.Fatalf("usage command missing from capture: %v", captured)
		}
		p24ExerciseCapturedLeaf(t, root, binary, operation, field, entry)
	}
	c.observe("usage-record", reply, 0)
	if p21PathString(t, reply, "result", "usage", "task_id") != task {
		t.Fatal(reply)
	}
	c.observe("usage-retry", w.call(args...), 0)
	args[len(args)-1] = "duplicate"
	c.failure("usage-duplicate", binary, root, args...)
	args[5] = "agent:other"
	c.failure("usage-denied", binary, root, args...)
	args[5] = w.actor
	args[7] = "bogus"
	c.failure("usage-invalid", binary, root, args...)
}

func TestUsageLargeCountsReplay(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	w := loopWorker{t, binary, root, "agent:builder"}
	task := w.create("task", "Large usage", "")
	claimed := w.call("next", "--claim", "--actor", w.actor, "--model", "openai/sol")
	hold := p21PathString(t, claimed, "result", "hold", "id")
	path := filepath.Join(t.TempDir(), "large.jsonl")
	raw := `{"type":"turn.completed","usage":{"input_tokens":9007199254740993,"cached_input_tokens":0,"cache_write_input_tokens":0,"output_tokens":0,"reasoning_output_tokens":0}}`
	if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
		t.Fatal(err)
	}
	w.call("usage", "record", "--hold", hold, "--actor", w.actor, "--tool", "codex", "--tool-output", path, "--idempotency-key", "large")
	w.call("rebuild")
	text, _ := closureRunFID(t, binary, root, "sql", "SELECT total FROM usage_sessions WHERE task_id='"+task+"'")
	if !strings.Contains(text, "9007199254740993") {
		t.Fatal(text)
	}
	archive := filepath.Join(t.TempDir(), "large.zip")
	w.call("export", "create", archive)
	w.call("export", "validate", archive)
}
