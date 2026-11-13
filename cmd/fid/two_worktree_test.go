package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// v3Answer is what one fid --json next --claim answered.
type v3Answer struct {
	exit  int
	kind  string
	task  string
	hold  string
	event bool
}

func v3Next(t *testing.T, binary, worktree, actor string) v3Answer {
	t.Helper()
	stdout, stderr, exit := runAdoptedFID(t, binary, worktree, "--json", "next", "--claim", "--actor", actor)
	var envelope struct {
		Result struct {
			Kind string `json:"kind"`
			Task *struct {
				ID string `json:"id"`
			} `json:"task"`
			Hold *struct {
				ID     string `json:"id"`
				TaskID string `json:"task_id"`
			} `json:"hold"`
		} `json:"result"`
		Event json.RawMessage `json:"event"`
	}
	if exit != 0 || json.Unmarshal([]byte(stdout), &envelope) != nil {
		t.Fatalf("next --claim in %s = exit %d\n%s%s", worktree, exit, stdout, stderr)
	}
	answer := v3Answer{exit: exit, kind: envelope.Result.Kind, event: string(envelope.Event) != "null"}
	if envelope.Result.Task != nil {
		answer.task = envelope.Result.Task.ID
	}
	if envelope.Result.Hold != nil {
		answer.hold = envelope.Result.Hold.ID
		if envelope.Result.Hold.TaskID != answer.task {
			t.Fatalf("next --claim in %s held %s for task %s", worktree, envelope.Result.Hold.TaskID, answer.task)
		}
	}
	return answer
}

// v3Issue is an issue as fid show's JSON lists it among an epic's tasks.
type v3Issue struct {
	ID             string   `json:"id"`
	StatusCategory string   `json:"status_category"`
	BlockedBy      []string `json:"blocked_by"`
}

// assertV3Servable checks that nothing should have stopped task: it has no
// open blocker, and in a sequential epic it is the next unfinished task in
// sequence, or that task waits on it through blocks links in the epic. It
// returns the task's epic when that epic is sequential, else "".
func assertV3Servable(t *testing.T, workspace, task string, sequential map[string]bool) string {
	t.Helper()
	var shown struct {
		Result struct {
			BlockedBy []struct {
				ID             string `json:"id"`
				StatusCategory string `json:"status_category"`
			} `json:"blocked_by"`
			Epic *struct {
				ID string `json:"id"`
			} `json:"epic"`
		} `json:"result"`
	}
	result := executeCommandForTest(workspace, "--json", "show", task)
	if result.err != nil || json.Unmarshal([]byte(result.stdout), &shown) != nil {
		t.Fatalf("show %s = %#v", task, result)
	}
	for _, blocker := range shown.Result.BlockedBy {
		if blocker.StatusCategory != "completed" && blocker.StatusCategory != "cancelled" {
			t.Fatalf("served %s while %s, which blocks it, is %s", task, blocker.ID, blocker.StatusCategory)
		}
	}
	if shown.Result.Epic == nil || !sequential[shown.Result.Epic.ID] {
		return ""
	}
	var epic struct {
		Result struct {
			Children []v3Issue `json:"children"`
		} `json:"result"`
	}
	result = executeCommandForTest(workspace, "--json", "show", shown.Result.Epic.ID)
	if result.err != nil || json.Unmarshal([]byte(result.stdout), &epic) != nil {
		t.Fatalf("show %s = %#v", shown.Result.Epic.ID, result)
	}
	children := map[string]v3Issue{}
	next := ""
	for _, child := range epic.Result.Children {
		children[child.ID] = child
		if next == "" && child.StatusCategory != "completed" && child.StatusCategory != "cancelled" {
			next = child.ID
		}
	}
	waitsOn := map[string]bool{}
	for pending := []string{next}; len(pending) > 0; pending = pending[1:] {
		if !waitsOn[pending[0]] {
			waitsOn[pending[0]] = true
			for _, blocker := range children[pending[0]].BlockedBy {
				if _, inEpic := children[blocker]; inEpic {
					pending = append(pending, blocker)
				}
			}
		}
	}
	if !waitsOn[task] {
		t.Fatalf("served %s out of its epic's sequence; %s is next and doesn't wait on it", task, next)
	}
	return shown.Result.Epic.ID
}

// assertV3JournalNeverDoubleHeld replays the journal's hold Effects and
// checks that no task, and no worktree, ever had two open holds.
func assertV3JournalNeverDoubleHeld(t *testing.T, workspace string, wantHolds int) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), config.JournalDirName, "*.jsonl"))
	if err != nil || len(files) == 0 {
		t.Fatalf("journal files = %v, %v", files, err)
	}
	sort.Strings(files)
	type held struct{ task, worktree string }
	open, byTask, byWorktree, taken := map[string]held{}, map[string]int{}, map[string]int{}, 0
	for _, file := range files {
		journalFile, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(journalFile)
		scanner.Buffer(make([]byte, 1<<20), 1<<24)
		for scanner.Scan() {
			var event struct {
				EventType string `json:"event_type"`
				Effects   []struct {
					Op     string         `json:"op"`
					Table  string         `json:"table"`
					PK     map[string]any `json:"pk"`
					Values map[string]any `json:"values"`
				} `json:"effects"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				t.Fatal(err)
			}
			if event.EventType == "hold.taken" {
				taken++
			}
			for _, effect := range event.Effects {
				id, _ := effect.PK["id"].(string)
				switch {
				case effect.Table == "holds" && effect.Op == "insert":
					hold := held{effect.Values["task_id"].(string), effect.Values["worktree_id"].(string)}
					open[id] = hold
					byTask[hold.task]++
					byWorktree[hold.worktree]++
					if byTask[hold.task] > 1 || byWorktree[hold.worktree] > 1 {
						t.Fatalf("hold %s left task %s with %d open holds and its worktree with %d", id, hold.task, byTask[hold.task], byWorktree[hold.worktree])
					}
				case effect.Table == "holds" && effect.Op == "update" && effect.Values["ended_at"] != nil:
					hold, found := open[id]
					if !found {
						t.Fatalf("event %s ended hold %s, which wasn't open", event.EventType, id)
					}
					delete(open, id)
					byTask[hold.task]--
					byWorktree[hold.worktree]--
				}
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
		if err := journalFile.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if taken != wantHolds || len(open) != 0 {
		t.Fatalf("journal took %d holds and left %d open, want %d taken and none open", taken, len(open), wantHolds)
	}
}

// V3 (check:two-worktree): two worktrees share one ledger and never claim
// the same task; next respects dependencies. Three worktrees of one
// repository take tasks at the same moment, as separate fid processes,
// round after round, until nothing is left. No two take the same task,
// none takes a blocked task or one out of its epic's sequence, and the
// journal never holds two open holds on a task. A sequential epic serves
// one task at a time, across worktrees, even when links branch: no round
// leaves two open holds in one. The checks don't depend on which worktree
// wins.
func TestV3TwoWorktreesNeverHoldTheSameTask(t *testing.T) {
	t.Parallel()
	binary := buildAdoptedFID(t)
	main := gittest.NewWorkspace(t)
	runFIDCommand(t, main, "init", "--prefix", "FID")
	worktrees := []string{main, gittest.AddWorktree(t, main), gittest.AddWorktree(t, main)}
	create := func(args ...string) string {
		_, created := runFIDCommand(t, main, append([]string{"issue", "create"}, append(args, "--actor", "human:mark")...)...)
		return parseCreatedIssueID(t, created)
	}
	link := func(from, relation, to string) {
		runFIDCommand(t, main, "link", "add", from, relation, to, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark")
	}
	prioritised := create("--kind", "epic", "--title", "Prioritised", "--scheduling", "prioritised")
	sequential := create("--kind", "epic", "--title", "Sequential", "--scheduling", "sequential")
	linked := create("--kind", "epic", "--title", "Links beat sequence", "--scheduling", "sequential")
	branching := create("--kind", "epic", "--title", "Links that branch", "--scheduling", "sequential")
	tasks := []string{}
	for priority := 1; priority <= 3; priority++ {
		tasks = append(tasks, create("--kind", "task", "--title", fmt.Sprintf("Prioritised %d", priority), "--epic", prioritised, "--priority", fmt.Sprint(priority)))
	}
	inSequence := map[string][]string{}
	for _, epic := range []string{sequential, linked, branching} {
		for sequence := 1; sequence <= 3; sequence++ {
			id := create("--kind", "task", "--title", fmt.Sprintf("In sequence %d", sequence), "--epic", epic, "--sequence", fmt.Sprint(sequence))
			inSequence[epic] = append(inSequence[epic], id)
			tasks = append(tasks, id)
		}
	}
	link(inSequence[linked][0], "blocked_by", inSequence[linked][2])
	// The first task waits on two others, which nothing else stops.
	link(inSequence[branching][0], "blocked_by", inSequence[branching][1])
	link(inSequence[branching][0], "blocked_by", inSequence[branching][2])
	blocker, blocked := create("--kind", "task", "--title", "Blocker"), create("--kind", "task", "--title", "Blocked", "--priority", "0")
	link(blocker, "blocks", blocked)
	// A chain whose blocked middle comes first in fid status's order, so a
	// worktree would reach it in the round its blocker is taken.
	head, middle := create("--kind", "task", "--title", "Chain head"), create("--kind", "task", "--title", "Chain middle", "--priority", "0")
	end := create("--kind", "task", "--title", "Chain end")
	link(middle, "blocked_by", head)
	link(middle, "blocks", end)
	tasks = append(tasks, blocker, blocked, head, middle, end, create("--kind", "task", "--title", "Standalone"))
	isSequential := map[string]bool{sequential: true, linked: true, branching: true}

	served := map[string]int{}
	for round := 1; ; round++ {
		answers := make([]v3Answer, len(worktrees))
		release := make(chan struct{})
		var wait sync.WaitGroup
		for index, worktree := range worktrees {
			wait.Add(1)
			go func() {
				defer wait.Done()
				<-release
				answers[index] = v3Next(t, binary, worktree, fmt.Sprintf("agent:wt%d", index+1))
			}()
		}
		close(release)
		wait.Wait()
		// Each round's tasks are cancelled before the next, so the holds
		// taken in a round are all the open holds.
		took, epicHeld := map[string]int{}, map[string]string{}
		for index, answer := range answers {
			if answer.kind == "none" {
				if answer.task != "" || answer.hold != "" || answer.event {
					t.Fatalf("round %d: wt%d got none with %+v", round, index+1, answer)
				}
				continue
			}
			if answer.kind != "new" || answer.hold == "" || !answer.event {
				t.Fatalf("round %d: wt%d got %+v, want a new hold in an Event", round, index+1, answer)
			}
			if other, taken := took[answer.task]; taken {
				t.Fatalf("round %d: wt%d and wt%d both took %s", round, other+1, index+1, answer.task)
			}
			took[answer.task] = index
			served[answer.task]++
			if epic := assertV3Servable(t, main, answer.task, isSequential); epic != "" {
				if other, held := epicHeld[epic]; held {
					t.Fatalf("round %d: %s and %s of sequential epic %s are both held", round, other, answer.task, epic)
				}
				epicHeld[epic] = answer.task
			}
		}
		if len(took) == 0 {
			break
		}
		for task, index := range took {
			again := v3Next(t, binary, worktrees[index], fmt.Sprintf("agent:wt%d", index+1))
			if again.kind != "resume" || again.hold != answers[index].hold || again.event {
				t.Fatalf("round %d: wt%d's re-run = %+v, want resume of %s with no Event", round, index+1, again, answers[index].hold)
			}
			runFIDCommand(t, main, "issue", "transition", task, "--to", "cancelled", "--actor", "human:mark", "--reason", "Test cancellation")
		}
		if round > len(tasks) {
			t.Fatalf("still taking tasks after %d rounds", round)
		}
	}
	for _, task := range tasks {
		if served[task] != 1 {
			t.Errorf("task %s was taken %d times, want once", task, served[task])
		}
	}
	assertV3JournalNeverDoubleHeld(t, main, len(tasks))
	if snapshots := policySnapshots(t, main); len(snapshots) != 1 {
		t.Fatalf("concurrent claims wrote %d policy snapshots, want one", len(snapshots))
	}
	if remaining := executeCommandForTest(main, "--json", "next"); !strings.Contains(remaining.stdout, `"kind":"none"`) || remaining.err != nil {
		t.Fatalf("next with an empty backlog = %#v, want none", remaining)
	}
}
