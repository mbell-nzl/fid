package status

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

const (
	boxWidth = 72
	// rowWidth is how far a row runs: to two columns inside the box.
	rowWidth = boxWidth - 2
	// readyRows caps the text's task rows in READY, unless --all.
	readyRows   = 10
	minTitle    = 10
	progressBar = 10
	topIndent   = "  "
	taskIndent  = "      "
	// summaryIndent and nextIndent are where an item's summary and next
	// step begin, after their symbols; commandIndent is two more.
	summaryIndent = "    "
	nextIndent    = "        "
	commandIndent = "          "
)

// WriteText writes the result as a box: what needs a person, the work in
// progress and to do, policy warnings, and the backup's line, leaving out
// empty sections.
// READY shows at most ten task rows unless all is set. The result's text,
// which holds titles, paths and record values fid didn't write, has its
// control characters escaped before it is fitted to a column or coloured;
// JSON keeps it as it is.
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
		blocks = append(blocks, text.inProgress(result))
	}
	if len(result.Ready.Epics)+len(result.Ready.Tasks) > 0 {
		blocks = append(blocks, text.ready(result, all))
	}
	for _, warning := range result.Warnings {
		blocks = append(blocks, text.item(display.RoleWarning, display.SymbolWarning, warning.Summary, &warning.NextAction))
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

// item is a summary and its next step, each wrapped to the row's end. A
// next step that is a fid command wraps with a shell continuation.
func (text textWriter) item(role display.Role, symbol display.Symbol, summary string, next *string) []string {
	lines := []string{}
	for index, line := range display.Wrap(cli.EscapeControlText(summary), rowWidth-len(summaryIndent), rowWidth-len(summaryIndent)) {
		if index == 0 {
			lines = append(lines, topIndent+text.renderer.Line(role, symbol, line))
			continue
		}
		lines = append(lines, summaryIndent+text.renderer.Colorize(role, line))
	}
	if next == nil {
		return lines
	}
	step, indent := cli.EscapeControlText(*next), nextIndent
	steps := display.Wrap(step, rowWidth-len(nextIndent), rowWidth-len(nextIndent))
	if strings.HasPrefix(step, "fid ") {
		steps, indent = display.WrapCommand(step, rowWidth-len(nextIndent), rowWidth-len(commandIndent)), commandIndent
	}
	for index, line := range steps {
		if index == 0 {
			lines = append(lines, taskIndent+string(display.SymbolNext)+" "+line)
			continue
		}
		lines = append(lines, indent+line)
	}
	return lines
}

func (text textWriter) inProgress(result Result) []string {
	section := result.InProgress
	lines := []string{text.header("IN PROGRESS")}
	for _, epic := range section.Epics {
		bar := display.ProgressBar(epic.Progress.Done, epic.Progress.Total, progressBar) + " done"
		lines = append(lines, text.row(topIndent, display.SymbolEpic, epic.ID, epic.Title, bar, display.RoleNormal)...)
		for _, task := range epic.Tasks {
			lines = append(lines, text.row(taskIndent, display.SymbolStarted, task.ID, task.Title, result.startedVerdict(task), display.RoleMuted)...)
		}
	}
	for _, task := range section.Tasks {
		lines = append(lines, text.row(topIndent, display.SymbolStarted, task.ID, task.Title, result.startedVerdict(task), display.RoleMuted)...)
	}
	return lines
}

// ready lists every unstarted task under its epic. Without all, once ten task
// rows are shown the rest are counted in one line; a hidden epic with no
// tasks counts as one.
func (text textWriter) ready(result Result, all bool) []string {
	section := result.Ready
	lines := []string{text.header("READY")}
	shown, hidden := 0, 0
	task := func(indent string, task issues.OverviewTask) {
		if !all && shown >= readyRows {
			hidden++
			return
		}
		verdict, role := result.readyVerdict(task)
		lines = append(lines, text.row(indent, display.SymbolUnstarted, task.ID, task.Title, verdict, role)...)
		shown++
	}
	for _, epic := range section.Epics {
		if !all && shown >= readyRows {
			hidden += max(len(epic.Tasks), 1)
			continue
		}
		lines = append(lines, text.row(topIndent, display.SymbolEpic, epic.ID, epic.Title, epicCount(epic), display.RoleMuted)...)
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
// A title too long for its column is cut to end in an ellipsis, dropping an
// escape that crosses the cut whole. When right
// would leave the title fewer than minTitle columns, as the longest IDs and
// workflow states can, right takes a line of its own beneath, still flush
// to the row's end.
func (text textWriter) row(indent string, symbol display.Symbol, id, title, right string, role display.Role) []string {
	id, title, right = cli.EscapeControlText(id), cli.EscapeControlText(title), cli.EscapeControlText(right)
	line := indent + string(symbol) + " " + text.renderer.Colorize(display.RoleID, id) + strings.Repeat(" ", text.idWidth-len([]rune(id))) + "  "
	width := rowWidth - display.VisibleWidth(line)
	rightWidth := display.VisibleWidth(right)
	switch {
	case right == "":
		return []string{line + cli.CutEscapedText(title, width)}
	case width-2-rightWidth >= minTitle:
		return []string{line + display.PadRight(cli.CutEscapedText(title, width-2-rightWidth), width-2-rightWidth) + "  " + text.renderer.Colorize(role, right)}
	}
	return []string{line + cli.CutEscapedText(title, width),
		strings.Repeat(" ", max(rowWidth-rightWidth, 0)) + text.renderer.Colorize(role, right)}
}

// idWidth is how wide the longest ID shown is, once escaped.
func idWidth(result Result) int {
	width := 0
	measure := func(id string) { width = max(width, len([]rune(cli.EscapeControlText(id)))) }
	for _, section := range []issues.OverviewSection{result.InProgress, result.Ready} {
		for _, epic := range section.Epics {
			measure(epic.ID)
			for _, task := range epic.Tasks {
				measure(task.ID)
			}
		}
		for _, task := range section.Tasks {
			measure(task.ID)
		}
	}
	return width
}

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
