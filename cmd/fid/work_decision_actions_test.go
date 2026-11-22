package main

import (
	"strings"
	"testing"
	"time"
)

func followWorkDecisionAction(t *testing.T, f *statusFixture, item map[string]any, index int, operation string) map[string]any {
	t.Helper()
	actions := strings.Split(item["next_action"].(string), " or ")
	if index >= len(actions) {
		t.Fatalf("missing action %d: %v", index, item)
	}
	words := strings.Fields(actions[index])
	if len(words) == 0 || words[0] != "fid" {
		t.Fatalf("not a fid command: %q", actions[index])
	}
	// Tokenize before substituting so a multiword reason stays one argument.
	fill := strings.NewReplacer("<name>", "mark", "<key>", "suggested-action", "<text>", "Cover the uncovered case", "<n>", "1")
	args := []string{"--json"}
	for _, word := range words[1:] {
		arg := fill.Replace(word)
		if strings.ContainsAny(arg, "<>") {
			t.Fatalf("unfilled placeholder in action: %q", arg)
		}
		args = append(args, arg)
	}
	fields, _ := lifecycleResult(t, f.run(args...), operation)
	return fields
}

func TestAcceptanceSuggestedActions(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "subtask"} {
		for _, action := range []struct {
			name, operation string
			index           int
		}{
			{"accept", "task.accept", 0},
			{"cancel", "issue.transition", 1},
		} {
			t.Run(kind+"/"+action.name, func(t *testing.T) {
				t.Parallel()
				f := newStatusFixture(t, time.Hour)
				f.clock = time.Now
				parent := ""
				if kind == "subtask" {
					parent = createAcceptanceTask(t, f, "human:mark", "task", "")
				}
				issue := createAcceptanceTask(t, f, "agent:orchestrator", kind, parent)
				item := selectionItem(t, f, "accept")
				if item == nil || item["subject"].(map[string]any)["id"] != issue {
					t.Fatalf("missing acceptance item for %s: %v", issue, item)
				}
				if actions := strings.Split(item["next_action"].(string), " or "); len(actions) != 2 {
					t.Fatalf("expected accept and cancel actions: %v", actions)
				}
				followWorkDecisionAction(t, f, item, action.index, action.operation)
				if remaining := selectionItem(t, f, "accept"); remaining != nil {
					t.Fatalf("%s did not clear acceptance item: %v", action.name, remaining)
				}
			})
		}
	}
}

func TestReworkAllowanceSuggestedAction(t *testing.T) {
	t.Parallel()
	for _, withSubtask := range []bool{false, true} {
		name := "task"
		if withSubtask {
			name = "task with subtask"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, task := allowanceFixture(t)
			if withSubtask {
				createAcceptanceTask(t, f, "human:mark", "subtask", task)
			}
			item := selectionItem(t, f, "rework-cap")
			if item == nil || item["subject"].(map[string]any)["id"] != task {
				t.Fatalf("missing task's rework-cap item: %v", item)
			}
			followWorkDecisionAction(t, f, item, 0, "task.rework-allow")
			if remaining := selectionItem(t, f, "rework-cap"); remaining != nil {
				t.Fatalf("suggested allowance did not clear cap: %v", remaining)
			}
			got := commandResultFields(t, f, "next")
			if got["kind"] != "rework" || got["task"].(map[string]any)["id"] != task {
				t.Fatalf("suggested allowance did not return rework: %v", got)
			}
		})
	}
}

func TestSubtaskHasNoReworkAllowanceAction(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, time.Hour)
	f.clock = time.Now
	parent := createAcceptanceTask(t, f, "human:mark", "task", "")
	subtask := createAcceptanceTask(t, f, "human:mark", "subtask", parent)
	if item := selectionItem(t, f, "rework-cap"); item != nil {
		t.Fatalf("subtasks have no rework cap: %v", item)
	}
	before := startEventCount(t, f.workspace)
	got := f.run("--json", "rework", "allow", subtask, "--rounds", "1", "--actor", "human:mark", "--reason", "Continue")
	assertJSONFailureCode(t, got.stderr, "task.rework-allow", "invalid-argument")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("refused subtask allowance changed journal")
	}
}
