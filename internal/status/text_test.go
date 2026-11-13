package status

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/next"
)

func textTask(id, status, category string, blockers ...string) issues.OverviewTask {
	return issues.OverviewTask{ID: id, Title: "Title of " + id, Status: status, StatusCategory: category, BlockedBy: append([]string{}, blockers...)}
}

func textEpic(id, state string, done, total int, tasks ...issues.OverviewTask) issues.OverviewEpic {
	return issues.OverviewEpic{ID: id, Title: "Title of " + id, Status: "todo", StatusCategory: "unstarted", State: state,
		Progress: issues.Progress{Done: done, Total: total}, Tasks: tasks}
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
		Ready: toDo}
}

func writeText(t *testing.T, result Result, renderer display.Renderer, all bool) string {
	t.Helper()
	if result.entries == nil {
		result.entries = map[string]next.Entry{}
		collect := func(task issues.OverviewTask) {
			entry := next.Entry{OverviewTask: task, EpicBlockedBy: []string{}}
			if len(task.BlockedBy) > 0 {
				reason := next.ReasonBlocked
				entry.Reason = &reason
			}
			result.entries[task.ID] = entry
		}
		for _, epic := range result.Ready.Epics {
			for _, task := range epic.Tasks {
				collect(task)
			}
		}
		for _, task := range result.Ready.Tasks {
			collect(task)
		}
	}
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
				toDo := text[strings.Index(text, "  READY\n"):]
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

func assertRowsFitTheBox(t *testing.T, text string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		if width := display.VisibleWidth(line); width > 72 || strings.HasSuffix(line, " ") {
			t.Fatalf("line is %d wide or ends in a space: %q\n%s", width, line, text)
		}
	}
}

func TestWriteTextFitsRowsInTheBox(t *testing.T) {
	t.Parallel()
	long := textTask("FID-long", "todo", "unstarted", "FID-b1", "FID-b2", "FID-b3")
	long.Title = strings.Repeat("A very long title ", 6)
	result := textResult(issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{long}})
	text := writeText(t, result, display.Renderer{}, false)
	assertRowsFitTheBox(t, text)
	if !strings.Contains(text, "  ○ FID-long  A very long title A very long tit…  blocked by FID-b1 +2\n") {
		t.Fatalf("long title isn't cut to fit its column:\n%s", text)
	}

	// The longest IDs fid makes, a 16-letter prefix and six hex digits, and
	// the longest workflow state leave a title too little room beside the
	// right-hand text, which then takes a line of its own.
	id := func(suffix string) string { return "ABCDEFGHIJKLMNOP-" + suffix }
	state := strings.Repeat("waiting_", 7) + "reviewer"
	started := textTask(id("a00001"), state, "started")
	epicTask := textTask(id("a00002"), "todo", "unstarted", id("b00001"), id("b00002"))
	standalone := textTask(id("a00003"), "todo", "unstarted", id("b00003"))
	result = textResult(issues.OverviewSection{Epics: []issues.OverviewEpic{textEpic(id("e00002"), issues.EpicToDo, 0, 1, epicTask)},
		Tasks: []issues.OverviewTask{standalone}})
	result.InProgress = issues.OverviewSection{Epics: []issues.OverviewEpic{textEpic(id("e00001"), issues.EpicInProgress, 3, 8, started)},
		Tasks: []issues.OverviewTask{}}
	for _, renderer := range []display.Renderer{{}, {Color: true}} {
		text = writeText(t, result, renderer, false)
		assertRowsFitTheBox(t, text)
		if renderer.Color {
			text = fidColors.Replace(text)
		}
		for _, want := range []string{
			"  ◆ ABCDEFGHIJKLMNOP-e00001  Title of ABCDEFGHI…  ███░░░░░░░  3/8 done\n",
			"      ● ABCDEFGHIJKLMNOP-a00001  Title of ABCDEFGHIJKLMNOP-a00001\n" +
				strings.Repeat(" ", 6) + "waiting waiting waiting waiting waiting waiting waiting reviewer\n",
			"      ○ ABCDEFGHIJKLMNOP-a00002  Title of ABCDEFGHIJKLMNOP-a00002\n" +
				strings.Repeat(" ", 33) + "blocked by ABCDEFGHIJKLMNOP-b00001 +1\n",
			"  ○ ABCDEFGHIJKLMNOP-a00003  Title of ABCDEFGHIJKLMNOP-a00003\n" +
				strings.Repeat(" ", 36) + "blocked by ABCDEFGHIJKLMNOP-b00003\n",
		} {
			if !strings.Contains(text, want) {
				t.Fatalf("color=%v: text doesn't show %q:\n%s", renderer.Color, want, text)
			}
		}
	}
}

// controls is text fid didn't write, carrying characters a terminal acts on:
// erase the screen, a line break, a carriage return, a tab, a bell and a C1
// control sequence introducer.
const controls = "\x1b[2J\nINJECTED\r\t\a\u009b"

// escapedControls is controls as the text shows it.
const escapedControls = `\u001b[2J\nINJECTED\r\t\u0007\u009b`

// fidColors are the colour codes the renderer adds itself.
var fidColors = strings.NewReplacer("\033[0m", "", "\033[1m", "", "\033[36m", "", "\033[31m", "", "\033[33m", "", "\033[32m", "", "\033[2m", "", "\033[35m", "")

// assertNoRawControls fails if text holds a control character other than
// the line breaks between its lines, or a line an injected break began.
func assertNoRawControls(t *testing.T, text string) {
	t.Helper()
	for _, value := range text {
		if value != '\n' && unicode.IsControl(value) {
			t.Fatalf("text holds control %U:\n%q", value, text)
		}
	}
	if strings.Contains(text, "\nINJECTED") {
		t.Fatalf("an injected line break reached the text:\n%q", text)
	}
}

func TestWriteTextEscapesControlsInWhatFidDidNotWrite(t *testing.T) {
	t.Parallel()
	needsNext := "fix /cfg" + controls + "/fid/config.toml, then fid backup"
	backupNext := "fix or remove /led" + controls + "/backups.jsonl, then fid backup"
	started := textTask("FID-t1", "in_progress", "started")
	started.Title = "T" + controls
	epic := textEpic("FID-e1", issues.EpicInProgress, 1, 3, started)
	epic.Title = "E" + controls
	blocked := textTask("FID-t2", "todo", "unstarted", "FID-b"+controls)
	result := Result{GeneratedAt: testAt,
		NeedsYou:   []NeedsYouItem{{Kind: NeedsYouBackup, Summary: "Can't read user config /cfg" + controls + "/fid/config.toml", NextAction: &needsNext}},
		Backup:     Backup{State: "warning", Summary: "Last backup failed (" + controls + ")", NextAction: &backupNext},
		InProgress: issues.OverviewSection{Epics: []issues.OverviewEpic{epic}, Tasks: []issues.OverviewTask{}},
		Ready:      issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{blocked}}}
	plain := writeText(t, result, display.Renderer{}, false)
	for _, want := range []string{
		"  ✗ Can't read user config\n    /cfg" + escapedControls + "/fid/config.toml\n",
		"      → fix /cfg" + escapedControls + "/fid/config.toml,\n        then fid backup\n",
		"  ⚠ Last backup failed (" + escapedControls + ")\n",
		"      → fix or remove\n        /led" + escapedControls + "/backups.jsonl, then\n        fid backup\n",
		"  ◆ FID-e1  E" + escapedControls + "  ",
		"      ● FID-t1  T" + escapedControls + "  ",
		"blocked by FID-b" + escapedControls + "\n",
	} {
		if !strings.Contains(plain, want) {
			t.Fatalf("text doesn't show %q:\n%s", want, plain)
		}
	}
	assertNoRawControls(t, plain)
	colored := writeText(t, result, display.Renderer{Color: true}, false)
	if fidColors.Replace(colored) != plain {
		t.Fatalf("colored output differs from plain beyond its colours:\n%q\n%q", fidColors.Replace(colored), plain)
	}
	assertNoRawControls(t, fidColors.Replace(colored))
}

// Prose wraps at spaces, its continuation lines indented to where it began;
// a next step that is a fid command wraps with a shell continuation.
func TestWriteTextWrapsLongLines(t *testing.T) {
	t.Parallel()
	path := "/home/mark/.config/fid/config.toml"
	next := "add [backup] target to " + path + ", then fid backup"
	result := Result{GeneratedAt: testAt, NeedsYou: []NeedsYouItem{},
		Backup:     Backup{State: "warning", Summary: "Can't read user config " + path + "; last backed up 1 hour ago", NextAction: &next},
		InProgress: issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{}},
		Ready:      issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{}},
		Hint:       &Hint{Summary: "No tasks yet", NextAction: CreateTaskCommand}}
	for _, renderer := range []display.Renderer{{}, {Color: true}} {
		text := writeText(t, result, renderer, false)
		assertRowsFitTheBox(t, text)
		text = fidColors.Replace(text)
		for _, want := range []string{
			"  ○ No tasks yet\n" +
				"      → fid issue create --kind task --title <title> \\\n" +
				"          --actor human:<name> --idempotency-key <key>\n",
			"  ⚠ Can't read user config " + path + "; last\n" +
				"    backed up 1 hour ago\n" +
				"      → add [backup] target to " + path + ",\n" +
				"        then fid backup\n",
		} {
			if !strings.Contains(text, want) {
				t.Fatalf("color=%v: text doesn't show %q:\n%s", renderer.Color, want, text)
			}
		}
	}
}

// A title cut to its column drops an escape that crosses the cut whole.
func TestWriteTextCutsATitleBeforeAnEscape(t *testing.T) {
	t.Parallel()
	task := textTask("FID-x", "todo", "unstarted")
	task.Title = strings.Repeat("a", 55) + controls
	result := textResult(issues.OverviewSection{Epics: []issues.OverviewEpic{}, Tasks: []issues.OverviewTask{task}})
	text := writeText(t, result, display.Renderer{}, false)
	if want := "  ○ FID-x   " + strings.Repeat("a", 55) + "…\n"; !strings.Contains(text, want) {
		t.Fatalf("text doesn't show %q:\n%s", want, text)
	}
}
