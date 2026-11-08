package status

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

func textTask(id, status, category string, blockers ...string) issues.OverviewTask {
	return issues.OverviewTask{ID: id, Title: "Title of " + id, Status: status, StatusCategory: category, BlockedBy: append([]string{}, blockers...)}
}

func textEpic(id, state string, done, total int, tasks ...issues.OverviewTask) issues.OverviewEpic {
	return issues.OverviewEpic{ID: id, Title: "Title of " + id, Status: "todo", StatusCategory: "unstarted", State: state,
		Progress: issues.EpicProgress{Done: done, Total: total}, Tasks: tasks}
}

func todo(prefix string, count int) []issues.OverviewTask {
	tasks := []issues.OverviewTask{}
	for index := range count {
		tasks = append(tasks, textTask(fmt.Sprintf("%s-%02d", prefix, index+1), "todo", "unstarted"))
	}
	return tasks
}

func textResult(toDo issues.OverviewSection) Result {
	next := "fid backup"
	return Result{GeneratedAt: testAt, NeedsYou: []NeedsYouItem{},
		Backup: Backup{State: "warning", Summary: "Last backup failed (push-rejected) 1 hour ago", NextAction: &next},
		InProgress: issues.OverviewSection{
			Epics: []issues.OverviewEpic{textEpic("FID-e1", issues.EpicInProgress, 3, 8, textTask("FID-t1", "qa_review", "started"))},
			Tasks: []issues.OverviewTask{textTask("FID-s1", "in_progress", "started")}},
		ToDo: toDo}
}

func writeText(t *testing.T, result Result, renderer display.Renderer, all bool) string {
	t.Helper()
	var output bytes.Buffer
	if err := WriteText(&output, result, renderer, all); err != nil {
		t.Fatalf("WriteText() error = %v", err)
	}
	return output.String()
}

func TestWriteTextColorsOnlyWhenTheRendererDoes(t *testing.T) {
	t.Parallel()
	result := textResult(issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{textTask("FID-x", "todo", "unstarted", "FID-t1")}})
	colored := writeText(t, result, display.Renderer{Color: true}, false)
	plain := writeText(t, result, display.Renderer{}, false)
	if !strings.Contains(colored, "\033[") {
		t.Fatalf("colored output has no ANSI codes:\n%s", colored)
	}
	if strings.Contains(plain, "\033") {
		t.Fatalf("plain output has ANSI codes:\n%q", plain)
	}
	strip := strings.NewReplacer("\033[0m", "", "\033[1m", "", "\033[36m", "", "\033[31m", "", "\033[33m", "", "\033[32m", "", "\033[2m", "")
	if strip.Replace(colored) != plain {
		t.Fatalf("colored output differs from plain beyond its colours:\n%s\n%s", strip.Replace(colored), plain)
	}
	assertNoKernelWords(t, plain)
}

func TestWriteTextCapsToDoAtTenTaskRows(t *testing.T) {
	t.Parallel()
	section := func(epics []issues.OverviewEpic, tasks []issues.OverviewTask) issues.OverviewSection {
		return issues.OverviewSection{Epics: epics, Tasks: tasks}
	}
	for _, test := range []struct {
		name    string
		toDo    issues.OverviewSection
		shown   int
		more    string
		epicsIn int
	}{
		{"ten rows fit", section(nil, todo("FID-s", 10)), 10, "", 0},
		{"twelve rows", section(nil, todo("FID-s", 12)), 10, "+2 more: fid status --all", 0},
		{"epic lines don't count", section([]issues.OverviewEpic{textEpic("FID-e1", issues.EpicInProgress, 3, 8, todo("FID-a", 8)...)},
			todo("FID-s", 4)), 10, "+2 more: fid status --all", 1},
		{"a hidden epic with no tasks counts once", section([]issues.OverviewEpic{textEpic("FID-e2", issues.EpicToDo, 0, 10, todo("FID-a", 10)...),
			textEpic("FID-e3", issues.EpicToDo, 0, 0)}, todo("FID-s", 1)), 10, "+2 more: fid status --all", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			for _, all := range []bool{false, true} {
				text := writeText(t, textResult(test.toDo), display.Renderer{}, all)
				toDo := text[strings.Index(text, "  TO DO\n"):]
				rows := strings.Count(toDo, "○ FID-")
				epics := strings.Count(toDo, "◆ FID-")
				total := len(test.toDo.Tasks)
				for _, epic := range test.toDo.Epics {
					total += len(epic.Tasks)
				}
				wantRows, wantMore := test.shown, test.more
				if all {
					wantRows, wantMore = total, ""
				}
				if rows != wantRows || all && epics != len(test.toDo.Epics) || !all && epics != test.epicsIn {
					t.Fatalf("all=%v: %d task rows and %d epics, want %d rows:\n%s", all, rows, epics, wantRows, text)
				}
				if got := strings.Contains(toDo, "more"); got != (wantMore != "") || wantMore != "" && !strings.Contains(toDo, "  "+wantMore+"\n") {
					t.Fatalf("all=%v: more line, want %q:\n%s", all, wantMore, text)
				}
			}
		})
	}
}

func TestWriteTextFitsRowsInTheBox(t *testing.T) {
	t.Parallel()
	long := textTask("FID-long", "todo", "unstarted", "FID-b1", "FID-b2", "FID-b3")
	long.Title = strings.Repeat("A very long title ", 6)
	result := textResult(issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{long}})
	text := writeText(t, result, display.Renderer{}, false)
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if width := display.VisibleWidth(line); width > 72 || strings.HasSuffix(line, " ") {
			t.Fatalf("line is %d wide or ends in a space: %q", width, line)
		}
	}
	if !strings.Contains(text, "  ○ FID-long  A very long title A very long tit…  blocked by FID-b1 +2\n") {
		t.Fatalf("long title isn't cut to fit its column:\n%s", text)
	}
}
