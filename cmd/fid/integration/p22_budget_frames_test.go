//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestP22BudgetCrossingFrameOrder(t *testing.T) {
	a := p22NewFIDRecoveryAdapterWithLease(t, "agent", 3600)
	p22AcceptedCall(t, a, "frame-baseline", "run.observe-budget", map[string]any{"value": 0})
	for attempt := 1; attempt <= 3; attempt++ {
		fixture := p22ConformanceCase{Name: fmt.Sprintf("frame-attempt-%d", attempt), Operation: "gate.run"}
		args, input := a.p22Args(t, fixture)
		args = append(args, "--input", p21Input(t, a.workspace, fixture.Name, input), "--idempotency-key", fixture.Name)
		stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
		if err != nil || stderr != "" {
			t.Fatalf("admitted attempt%d: %v %q %q", attempt, err, stdout, stderr)
		}
		lines := strings.Split(strings.TrimSpace(stdout), "\n")
		budgetFrames := 0
		for i, line := range lines {
			var frame map[string]any
			if err := json.Unmarshal([]byte(line), &frame); err != nil {
				t.Fatal(err)
			}
			if frame["sequence"] != float64(i+1) || frame["frame_schema"] != "fid.live.v1" {
				t.Fatalf("frame ordering/schema: %#v", frame)
			}
			data := frame["data"].(map[string]any)
			if frame["stream"] == "control" && data["kind"] == "budget-stop" {
				budgetFrames++
				p22AssertKeys(t, data, "kind stop_ids attention_ids budget")
				if i == len(lines)-1 || p22ResultPath(data, "stop_ids.#") != float64(1) || p22ResultPath(data, "attention_ids.#") != float64(1) || p22ResultPath(data, "budget.state") != "exhausted" {
					t.Fatalf("budget crossing frame: %#v", frame)
				}
			}
			if i == len(lines)-1 && (frame["stream"] != "terminal" || data["outcome"] != "succeeded" || p22ResultPath(data, "response.ok") != true) {
				t.Fatalf("normal observed terminal lost: %#v", frame)
			}
		}
		want := 0
		if attempt == 3 {
			want = 1
		}
		if budgetFrames != want {
			t.Fatalf("admitted attempt%d budget-stop frames = %d, want %d", attempt, budgetFrames, want)
		}
		if attempt == 3 {
			before := p21JournalEventCount(t, a.workspace)
			replay, replayErr, err := closureRunCommand(a.workspace, a.binary, args...)
			if err != nil || replayErr != "" {
				t.Fatalf("terminal replay: %v %q", err, replayErr)
			}
			var frame map[string]any
			if err := json.Unmarshal([]byte(replay), &frame); err != nil {
				t.Fatalf("replay must be one frame: %v %q", err, replay)
			}
			if frame["stream"] != "terminal" || p22ResultPath(frame, "data.outcome") != "succeeded" || p22ResultPath(frame, "data.response.result.replayed") != true {
				t.Fatalf("terminal-only replay: %#v", frame)
			}
			if p21JournalEventCount(t, a.workspace) != before {
				t.Fatal("terminal replay appended an Event")
			}
		}
	}
}
