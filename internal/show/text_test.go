package show

import (
	"bytes"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
)

func reference(id, kind, status, category string) issues.IssueReference {
	return issues.IssueReference{ID: id, Kind: kind, Title: "Title of " + id, Status: status, StatusCategory: category}
}

func child(id, status, category string, blockers ...string) issues.ChildIssue {
	return issues.ChildIssue{IssueReference: reference(id, "subtask", status, category), BlockedBy: append([]string{}, blockers...)}
}

// textTask is a started task with an epic, a blocker, criteria and a
// cancelled subtask.
func textTask() Result {
	priority := 2
	spec := "docs/specs/02_TRD.md"
	return Result{IssueView: issues.IssueView{
		Issue: issues.IssueDetail{IssueSummary: issues.IssueSummary{ID: "FID-t1", Kind: "task", Title: "Title of FID-t1", Status: "in_progress",
			StatusCategory: "started", Priority: &priority, SpecMDPath: &spec, Labels: []string{"backup", "export"}},
			AcceptanceCriteria: []issues.ReadAcceptanceCriterion{{ID: "AC-1", Ordinal: 1, Text: "It works"}}},
		Epic:      &issues.EpicReference{IssueReference: reference("FID-e1", "epic", "in_progress", "started"), Progress: issues.Progress{Done: 2, Total: 7}},
		BlockedBy: []issues.IssueReference{reference("FID-b1", "task", "qa_review", "started"), reference("FID-b2", "task", "done", "completed")},
		Blocks:    []issues.IssueReference{},
		Progress:  &issues.Progress{Done: 1, Total: 2},
		Children:  []issues.ChildIssue{child("FID-s1", "done", "completed"), child("FID-s2", "todo", "unstarted", "FID-b1"), child("FID-s3", "cancelled", "cancelled")},
	}}
}

func writeText(t *testing.T, result Result, renderer display.Renderer) string {
	t.Helper()
	var output bytes.Buffer
	if err := WriteText(&output, result, renderer); err != nil {
		t.Fatalf("WriteText() error = %v", err)
	}
	return output.String()
}

// fidColors are the colour codes the renderer adds itself.
var fidColors = strings.NewReplacer("\033[0m", "", "\033[1m", "", "\033[36m", "", "\033[31m", "", "\033[33m", "", "\033[32m", "", "\033[2m", "", "\033[35m", "")

func assertRowsFitTheBox(t *testing.T, text string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if width := display.VisibleWidth(line); width > 72 || strings.HasSuffix(line, " ") {
			t.Fatalf("line is %d wide or ends in a space: %q\n%s", width, line, text)
		}
	}
}

func TestWriteTextColorsOnlyWhenTheRendererDoes(t *testing.T) {
	t.Parallel()
	colored := writeText(t, textTask(), display.Renderer{Color: true})
	plain := writeText(t, textTask(), display.Renderer{})
	if strings.Contains(plain, "\033") {
		t.Fatalf("plain output has ANSI codes:\n%q", plain)
	}
	if fidColors.Replace(colored) != plain {
		t.Fatalf("colored output differs from plain beyond its colours:\n%s\n%s", fidColors.Replace(colored), plain)
	}
	for _, want := range []string{
		"╭ \033[36mFID-t1\033[0m ─",
		"  \033[35mPriority  \033[0m  P2\n",
		"\033[32m✓\033[0m \033[36mFID-s1\033[0m",
		// Cancelled isn't an error, so its mark is muted, not red.
		"\033[2m✗\033[0m \033[36mFID-s3\033[0m",
	} {
		if !strings.Contains(colored, want) {
			t.Fatalf("colored output doesn't show %q:\n%q", want, colored)
		}
	}
}

// Text lists only the blockers that aren't completed or cancelled, and a
// blocked unstarted child names its first blocker.
func TestWriteTextShowsWhatIsWaitingNow(t *testing.T) {
	t.Parallel()
	text := writeText(t, textTask(), display.Renderer{})
	assertRowsFitTheBox(t, text)
	for _, want := range []string{
		// Each row runs to column 70.
		"  Blocked by  ● FID-b1  Title of FID-b1" + strings.Repeat(" ", 70-39-9) + "qa review\n",
		"  ○ FID-s2  Title of FID-s2" + strings.Repeat(" ", 70-27-17) + "blocked by FID-b1\n",
		"  SUBTASKS" + strings.Repeat(" ", 70-10-8) + "1/2 done\n",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("text doesn't show %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "FID-b2") || strings.Contains(text, "Blocks") {
		t.Fatalf("text shows a cleared blocker or an empty row:\n%s", text)
	}
}

// The heading's title wraps and is never cut, its state staying on the
// first line; a state that leaves the title too little room takes a line
// of its own beneath.
func TestWriteTextWrapsTheHeadingsTitle(t *testing.T) {
	t.Parallel()
	result := textTask()
	result.Issue.Title = "Retry a rejected backup push once after fetching and rebasing fid's clone, and record the retry"
	text := writeText(t, result, display.Renderer{})
	assertRowsFitTheBox(t, text)
	if want := "  ● Retry a rejected backup push once after fetching and   in progress\n" +
		"    rebasing fid's clone, and record the retry\n"; !strings.Contains(text, want) {
		t.Fatalf("text doesn't show %q:\n%s", want, text)
	}
	result.Issue.Status = strings.Repeat("waiting_", 7) + "reviewer"
	text = writeText(t, result, display.Renderer{})
	assertRowsFitTheBox(t, text)
	if want := "  ● Retry a rejected backup push once after fetching and rebasing\n" +
		"    fid's clone, and record the retry\n" +
		strings.Repeat(" ", 6) + "waiting waiting waiting waiting waiting waiting waiting reviewer\n"; !strings.Contains(text, want) {
		t.Fatalf("text doesn't show %q:\n%s", want, text)
	}
}

// controls is text fid didn't write, carrying characters a terminal acts on.
const controls = "\x1b[2J\nINJECTED\r\t\a"

// escapedControls is controls as the text shows it.
const escapedControls = `\u001b[2J\nINJECTED\r\t\u0007`

func TestWriteTextEscapesControlsInWhatFidDidNotWrite(t *testing.T) {
	t.Parallel()
	result := textTask()
	result.Issue.Title = "T" + controls
	result.Issue.Labels = []string{"L" + controls}
	spec := "S" + controls
	result.Issue.SpecMDPath = &spec
	result.Issue.AcceptanceCriteria[0].Text = "C" + controls
	result.Epic.Title = "E" + controls
	result.Children[1].BlockedBy = []string{"FID-b" + controls}
	result.Now = time.Date(2026, 11, 11, 9, 10, 0, 0, time.UTC)
	result.Working = &holds.Working{Worktree: "/tmp/" + controls, Actor: "agent:test" + controls, Since: "2026-11-11T09:02:00Z", ExpiresAt: "2026-11-11T10:02:00Z"}
	for _, renderer := range []display.Renderer{{}, {Color: true}} {
		text := fidColors.Replace(writeText(t, result, renderer))
		for _, want := range []string{"  ● T" + escapedControls + " ", "  Labels      L" + escapedControls + "\n", "  Spec        S" + escapedControls + "\n",
			"  1  C" + escapedControls + "\n", "◆ FID-e1  E" + escapedControls + " ", "blocked by FID-b" + escapedControls + "\n",
			"/tmp/" + escapedControls, "agent:test" + escapedControls} {
			if !strings.Contains(text, want) {
				t.Fatalf("color=%v: text doesn't show %q:\n%s", renderer.Color, want, text)
			}
		}
		for _, value := range text {
			if value != '\n' && unicode.IsControl(value) {
				t.Fatalf("color=%v: text holds control %U:\n%q", renderer.Color, value, text)
			}
		}
	}
}

func TestWorkingUsesLocalDatesAndFollowsSpec(t *testing.T) {
	t.Parallel()
	result := textTask()
	result.Now = time.Date(2026, 11, 12, 0, 10, 0, 0, time.FixedZone("NZDT", 13*60*60))
	result.Working = &holds.Working{Here: true, Since: "2026-11-11T10:40:00Z", ExpiresAt: "2026-11-11T11:40:00Z"}
	text := writeText(t, result, display.Renderer{})
	want := "  Spec        docs/specs/02_TRD.md\n  Lane        not set\n  Working     here, since 2026-11-11 23:40, until 00:40\n\n  ACCEPTANCE CRITERIA\n"
	if !strings.Contains(text, want) {
		t.Fatalf("working isn't local or isn't in its reserved slot:\n%s", text)
	}
}

// JSON has every field, null rather than left out.
func TestFieldsHoldsEveryField(t *testing.T) {
	t.Parallel()
	result := textTask()
	result.Epic, result.Progress = nil, nil
	encoded, err := json.Marshal(result.Fields())
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for key := range fields {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	if want := []string{"blocked_by", "blocks", "checks", "checks_base", "children", "epic", "handoff", "issue", "kept", "progress", "questions", "scheduling", "state", "task", "warnings", "working"}; !slices.Equal(keys, want) {
		t.Fatalf("fields = %v, want %v", keys, want)
	}
	for _, name := range []string{"handoff", "checks_base", "epic", "progress", "scheduling", "state", "task", "working", "kept"} {
		if fields[name] != nil {
			t.Fatalf("%s = %v, want null", name, fields[name])
		}
	}
}
