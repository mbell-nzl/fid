package status

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

const (
	boxWidth = 72
	// rowWidth is how far a row runs: to two columns inside the box.
	rowWidth = boxWidth - 2
	// toDoRows caps the text's task rows in TO DO, unless --all.
	toDoRows    = 10
	minTitle    = 10
	progressBar = 10
	topIndent   = "  "
	taskIndent  = "      "
)

// WriteText writes the result as a box: what needs a person, the work in
// progress and to do, and the backup's line, leaving out empty sections.
// TO DO shows at most ten task rows unless all is set.
func WriteText(w io.Writer, result Result, renderer display.Renderer, all bool) error {
	text := textWriter{renderer: renderer, idWidth: idWidth(result)}
	blocks := [][]string{}
	if len(result.NeedsYou) > 0 {
		block := []string{text.header("NEEDS YOU")}
		for _, item := range result.NeedsYou {
			block = append(block, text.item(display.RoleError, display.SymbolError, item.Summary, item.NextAction)...)
		}
		blocks = append(blocks, block)
	}
	if result.Hint != nil {
		blocks = append(blocks, text.item(display.RoleNormal, display.SymbolUnstarted, result.Hint.Summary, &result.Hint.NextAction))
	}
	if len(result.InProgress.Epics)+len(result.InProgress.Tasks) > 0 {
		blocks = append(blocks, text.inProgress(result.InProgress))
	}
	if len(result.ToDo.Epics)+len(result.ToDo.Tasks) > 0 {
		blocks = append(blocks, text.toDo(result.ToDo, all))
	}
	if result.Backup.State != string(backup.HealthNeedsYou) {
		role, symbol := backupMark(result.Backup.State)
		blocks = append(blocks, text.item(role, symbol, result.Backup.Summary, result.Backup.NextAction))
	}
	lines := []string{display.TopBorder(renderer.Colorize(display.RoleHeader, "fid status"), result.GeneratedAt.Format("2006-01-02"), boxWidth)}
	for _, block := range blocks {
		lines = append(append(lines, ""), block...)
	}
	lines = append(lines, display.BottomBorder(boxWidth))
	if _, err := io.WriteString(w, strings.Join(lines, "\n")+"\n"); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	return nil
}

func backupMark(state string) (display.Role, display.Symbol) {
	switch state {
	case string(backup.HealthCurrent):
		return display.RoleSuccess, display.SymbolSuccess
	case string(backup.HealthNotBackedUpYet):
		return display.RoleNormal, display.SymbolUnstarted
	}
	return display.RoleWarning, display.SymbolWarning
}

type textWriter struct {
	renderer display.Renderer
	idWidth  int
}

func (text textWriter) header(name string) string {
	return topIndent + text.renderer.Colorize(display.RoleHeader, name)
}

func (text textWriter) item(role display.Role, symbol display.Symbol, summary string, next *string) []string {
	lines := []string{topIndent + text.renderer.Line(role, symbol, summary)}
	if next != nil {
		lines = append(lines, taskIndent+string(display.SymbolNext)+" "+*next)
	}
	return lines
}

func (text textWriter) inProgress(section issues.OverviewSection) []string {
	lines := []string{text.header("IN PROGRESS")}
	for _, epic := range section.Epics {
		bar := display.ProgressBar(epic.Progress.Done, epic.Progress.Total, progressBar) + " done"
		lines = append(lines, text.row(topIndent, display.SymbolEpic, epic.ID, epic.Title, bar, display.RoleNormal))
		for _, task := range epic.Tasks {
			lines = append(lines, text.row(taskIndent, display.SymbolStarted, task.ID, task.Title, stateText(task.Status), display.RoleMuted))
		}
	}
	for _, task := range section.Tasks {
		lines = append(lines, text.row(topIndent, display.SymbolStarted, task.ID, task.Title, stateText(task.Status), display.RoleMuted))
	}
	return lines
}

// toDo lists every unstarted task under its epic. Without all, once ten task
// rows are shown the rest are counted in one line; a hidden epic with no
// tasks counts as one.
func (text textWriter) toDo(section issues.OverviewSection, all bool) []string {
	lines := []string{text.header("TO DO")}
	shown, hidden := 0, 0
	task := func(indent string, task issues.OverviewTask) {
		if !all && shown >= toDoRows {
			hidden++
			return
		}
		lines = append(lines, text.row(indent, display.SymbolUnstarted, task.ID, task.Title, blockedText(task.BlockedBy), display.RoleWarning))
		shown++
	}
	for _, epic := range section.Epics {
		if !all && shown >= toDoRows {
			hidden += max(len(epic.Tasks), 1)
			continue
		}
		lines = append(lines, text.row(topIndent, display.SymbolEpic, epic.ID, epic.Title, epicCount(epic), display.RoleMuted))
		for _, listed := range epic.Tasks {
			task(taskIndent, listed)
		}
	}
	for _, listed := range section.Tasks {
		task(topIndent, listed)
	}
	if hidden > 0 {
		lines = append(lines, topIndent+text.renderer.Colorize(display.RoleMuted, fmt.Sprintf("+%d more: fid status --all", hidden)))
	}
	return lines
}

// row is a symbol, an ID and a title, with right set flush to the row's end.
// A title too long for its column is cut to end in an ellipsis.
func (text textWriter) row(indent string, symbol display.Symbol, id, title, right string, role display.Role) string {
	line := indent + string(symbol) + " " + text.renderer.Colorize(display.RoleID, id) + strings.Repeat(" ", text.idWidth-len([]rune(id))) + "  "
	width := rowWidth - display.VisibleWidth(line)
	if right == "" {
		return line + strings.TrimRight(display.PadRight(title, width), " ")
	}
	width -= 2 + display.VisibleWidth(right)
	return line + display.PadRight(title, max(width, minTitle)) + "  " + text.renderer.Colorize(role, right)
}

func idWidth(result Result) int {
	width := 0
	for _, section := range []issues.OverviewSection{result.InProgress, result.ToDo} {
		for _, epic := range section.Epics {
			width = max(width, len([]rune(epic.ID)))
			for _, task := range epic.Tasks {
				width = max(width, len([]rune(task.ID)))
			}
		}
		for _, task := range section.Tasks {
			width = max(width, len([]rune(task.ID)))
		}
	}
	return width
}

func stateText(status string) string { return strings.ReplaceAll(status, "_", " ") }

func epicCount(epic issues.OverviewEpic) string {
	count := len(epic.Tasks)
	switch {
	case epic.State == issues.EpicInProgress:
		return fmt.Sprintf("%d to do", count)
	case count == 0:
		return "no tasks"
	case count == 1:
		return "1 task"
	}
	return fmt.Sprintf("%d tasks", count)
}

func blockedText(blockers []string) string {
	switch len(blockers) {
	case 0:
		return ""
	case 1:
		return "blocked by " + blockers[0]
	}
	return fmt.Sprintf("blocked by %s +%d", blockers[0], len(blockers)-1)
}
