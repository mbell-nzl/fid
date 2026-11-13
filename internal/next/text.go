package next

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

const (
	boxWidth = 72
	// rowWidth is how far a row runs: to two columns inside the box.
	rowWidth = boxWidth - 2
	minTitle = 10
	indent   = "  "
	// proseIndent is where an item's text begins after its symbol, and where
	// it continues.
	proseIndent = "    "
	// stepIndent is where a next step's arrow sits; its command begins two
	// columns later and continues two further in.
	stepIndent    = "      "
	commandIndent = "          "
	// labelWidth is a labelled row's label; its value begins two columns on.
	labelWidth = 10
)

// WriteText writes the result as fid status's box, headed fid next and
// today's date: the task handed out with its rows and next step, or what
// is in the way when nothing is ready, and with explain the whole order
// with a verdict for each task. Text fid didn't write has its control
// characters escaped before it is fitted to a column or coloured. It never
// names the kernel's words, the flag that takes a task included.
func WriteText(w io.Writer, result Result, renderer display.Renderer, explain bool) error {
	text := textWriter{renderer: renderer, now: result.Now}
	blocks := [][]string{}
	if result.Served >= 0 {
		served := result.Queue[result.Served]
		blocks = append(blocks, text.row(indent, served.OverviewTask, len([]rune(cli.EscapeControlText(served.ID))), display.StateText(served.Status), display.RoleMuted))
		rows := text.rows(result, served)
		message := text.message(result, explain)
		if message == nil && (result.Claimed || !explain) {
			rows = append(rows, text.step("fid show "+served.ID)...)
		}
		for _, block := range [][]string{rows, message} {
			if len(block) > 0 {
				blocks = append(blocks, block)
			}
		}
	} else {
		blocks = append(blocks, text.nothing(result, explain))
	}
	if explain && len(result.Queue) > 0 {
		blocks = append(blocks, text.inOrder(result))
	}
	lines := []string{display.TopBorder(renderer.Colorize(display.RoleHeader, "fid next"), result.Now.Format("2006-01-02"), boxWidth)}
	for _, block := range blocks {
		lines = append(append(lines, ""), block...)
	}
	lines = append(lines, display.BottomBorder(boxWidth))
	if _, err := io.WriteString(w, strings.Join(lines, "\n")+"\n"); err != nil {
		return fmt.Errorf("write next: %w", err)
	}
	return nil
}

type textWriter struct {
	renderer display.Renderer
	now      time.Time
}

// rows are the served task's epic, priority and, while its hold is live,
// where and since when it is being worked.
func (text textWriter) rows(result Result, served Entry) []string {
	lines := []string{}
	if epic := result.Epic; epic != nil {
		reference := issues.OverviewTask{ID: epic.ID, Title: epic.Title}
		progress := fmt.Sprintf("%d/%d done", epic.Progress.Done, epic.Progress.Total)
		lines = append(lines, text.markedRow(text.label("Epic"), display.SymbolEpic, reference, len([]rune(cli.EscapeControlText(epic.ID))), progress, display.RoleMuted)...)
	}
	if served.Priority != nil {
		lines = append(lines, text.label("Priority")+"P"+strconv.Itoa(*served.Priority))
	}
	if hold := served.Working; hold != nil && hold.Here && !hold.Lapsed {
		lines = append(lines, text.label("Working")+fmt.Sprintf("here, since %s, until %s", text.clock(hold.Since), text.clock(hold.ExpiresAt)))
	}
	return lines
}

// message says what a resume carries on with, or how many tasks fid next
// passed over before the one it would hand out. It is nil otherwise. Work
// on an interrupted task stopped at its lapsed hold's last sign of life:
// its since, and its last heartbeat from 1.5a.
func (text textWriter) message(result Result, explain bool) []string {
	id := result.Queue[result.Served].ID
	switch {
	case result.Kind == KindResume && result.Answer != nil:
		return append(text.prose(indent, "Answer: "+result.Answer.Text), text.step("fid show "+id)...)
	case result.Kind == KindResume && result.Interrupted != nil:
		detail := "work stopped at " + text.clock(result.Interrupted.Since) + " without a park"
		if result.Interrupted.Reason != nil {
			detail = "parked at " + text.clock(result.Interrupted.Since) + ": " + *result.Interrupted.Reason
		}
		lines := text.item(display.RoleWarning, display.SymbolWarning, "Interrupted here: "+detail+".")
		lines = append(lines, text.prose(proseIndent, "Carry on from what's in this worktree.")...)
		return append(lines, text.step("fid show "+id)...)
	case result.Kind == KindResume:
		return append(text.prose(indent, "Already yours in this worktree: carry on with it"), text.step("fid show "+id)...)
	case !result.Claimed && !explain && result.Served > 0:
		return append(text.prose(indent, "Passed over "+plural(result.Served, "task", "tasks")+" ahead of this one"), text.step("fid next --explain")...)
	}
	return nil
}

// nothing says nothing is ready, counting what stops each task by kind, in
// the order each kind first appears.
func (text textWriter) nothing(result Result, explain bool) []string {
	if waiting := result.Waiting; waiting != nil {
		return append(text.item(display.RoleNormal, display.SymbolUnstarted, "Nothing here until you answer "+waiting.TaskID+": "+waiting.Question), text.step("fid show "+waiting.TaskID)...)
	}
	if len(result.Queue) == 0 {
		return text.item(display.RoleNormal, display.SymbolUnstarted, "Nothing ready here: no tasks to do")
	}
	counts, order := map[string]int{}, []string{}
	for _, entry := range result.Queue {
		kind := stoppedBy(entry)
		if counts[kind] == 0 {
			order = append(order, kind)
		}
		counts[kind]++
	}
	// Each count stays whole on a line, as the words of one item.
	units := []string{"Nothing", "ready", "here:"}
	for index, kind := range order {
		unit := cli.EscapeControlText(countText(kind, counts[kind]))
		if index < len(order)-1 {
			unit += ","
		}
		units = append(units, unit)
	}
	lines := []string{}
	for _, unit := range units {
		if last := len(lines) - 1; last >= 0 && display.VisibleWidth(lines[last])+1+display.VisibleWidth(unit) <= rowWidth-len(proseIndent) {
			lines[last] += " " + unit
			continue
		}
		lines = append(lines, unit)
	}
	for index, line := range lines {
		lines[index] = proseIndent + line
		if index == 0 {
			lines[index] = indent + text.renderer.Line(display.RoleNormal, display.SymbolUnstarted, line)
		}
	}
	if !explain {
		lines = append(lines, text.step("fid next --explain")...)
	}
	return lines
}

// stoppedBy is what stops entry, as a kind to count: its reason, with the
// worktree directory for work held elsewhere.
func stoppedBy(entry Entry) string {
	reason := *entry.Reason
	if entry.Kept != nil && reason == ReasonInterruptedElsewhere {
		return reason + " " + filepath.Base(entry.Kept.Worktree)
	}
	if entry.Working != nil && (reason == ReasonHeldElsewhere || reason == ReasonInterruptedElsewhere) {
		return reason + " " + filepath.Base(entry.Working.Worktree)
	}
	return reason
}

func countText(kind string, count int) string {
	reason, directory, _ := strings.Cut(kind, " ")
	switch reason {
	case ReasonWaitingForYou:
		return plural(count, "waiting for you", "waiting for you")
	case ReasonHeldElsewhere:
		return fmt.Sprintf("%d being worked in %s", count, directory)
	case ReasonInterruptedElsewhere:
		return fmt.Sprintf("%d interrupted in %s", count, directory)
	case ReasonNotNextInSequence:
		return fmt.Sprintf("%d not next in %s epic's sequence", count, map[bool]string{true: "its", false: "their"}[count == 1])
	case ReasonEpicBlocked:
		return plural(count, "in a blocked epic", "in blocked epics")
	case ReasonEpicClosed:
		return plural(count, "in a closed epic", "in closed epics")
	}
	return fmt.Sprintf("%d %s", count, reason)
}

func plural(count int, one, many string) string {
	if count == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", count, many)
}

// inOrder lists the whole order, each task with its verdict.
func (text textWriter) inOrder(result Result) []string {
	width := 0
	for _, entry := range result.Queue {
		width = max(width, len([]rune(cli.EscapeControlText(entry.ID))))
	}
	lines := []string{indent + text.renderer.Colorize(display.RoleHeader, "IN ORDER")}
	for index, entry := range result.Queue {
		verdict, role := Verdict(entry, index == result.Served)
		lines = append(lines, text.row(indent, entry.OverviewTask, width, verdict, role)...)
	}
	return lines
}

// Verdict is what fid next makes of a task in the order: next for the one
// it hands out, ready for any other that nothing stops, and otherwise what
// stops it. A held task names the worktree's directory.
func Verdict(entry Entry, served bool) (string, display.Role) {
	switch {
	case served:
		return "next", display.RoleSuccess
	case entry.Reason == nil:
		return "ready", display.RoleMuted
	}
	switch *entry.Reason {
	case ReasonWaitingForYou:
		return "waiting for you", display.RoleWarning
	case ReasonHeldElsewhere:
		return "working in " + filepath.Base(entry.Working.Worktree), display.RoleMuted
	case ReasonInterruptedElsewhere:
		return "interrupted in " + interruptedDirectory(entry), display.RoleWarning
	case ReasonNotNextInSequence:
		return "after " + *entry.After, display.RoleMuted
	case ReasonEpicBlocked:
		return "epic " + display.BlockedText(entry.EpicBlockedBy), display.RoleWarning
	case ReasonEpicClosed:
		return "epic closed", display.RoleMuted
	}
	return display.BlockedText(entry.BlockedBy), display.RoleWarning
}

// clock is a hold's time in the shown time zone: HH:MM today, with the date
// on any other day.
func (text textWriter) clock(stamp string) string {
	return display.Clock(stamp, text.now)
}

func (text textWriter) label(name string) string {
	return indent + text.renderer.Colorize(display.RoleLabel, display.PadRight(name, labelWidth)) + "  "
}

// item is a symbol and its text, wrapped under the text.
func (text textWriter) item(role display.Role, symbol display.Symbol, summary string) []string {
	lines := []string{}
	for index, line := range display.Wrap(cli.EscapeControlText(summary), rowWidth-len(proseIndent), rowWidth-len(proseIndent)) {
		if index == 0 {
			lines = append(lines, indent+text.renderer.Line(role, symbol, line))
			continue
		}
		lines = append(lines, proseIndent+text.renderer.Colorize(role, line))
	}
	return lines
}

// prose is text wrapped at prefix.
func (text textWriter) prose(prefix, summary string) []string {
	lines := []string{}
	for _, line := range display.Wrap(cli.EscapeControlText(summary), rowWidth-len(prefix), rowWidth-len(prefix)) {
		lines = append(lines, prefix+line)
	}
	return lines
}

// step is a next step: an arrow and a fid command, wrapped with a shell
// continuation.
func (text textWriter) step(command string) []string {
	lines := []string{}
	for index, line := range display.WrapCommand(cli.EscapeControlText(command), rowWidth-len(stepIndent)-2, rowWidth-len(commandIndent)) {
		if index == 0 {
			lines = append(lines, stepIndent+string(display.SymbolNext)+" "+line)
			continue
		}
		lines = append(lines, commandIndent+line)
	}
	return lines
}

// row is a task's mark, its ID padded to idWidth and its title, with right
// set flush to the row's end.
func (text textWriter) row(prefix string, task issues.OverviewTask, idWidth int, right string, role display.Role) []string {
	symbol := display.SymbolUnstarted
	if task.StatusCategory == "started" {
		symbol = display.SymbolStarted
	}
	return text.markedRow(prefix, symbol, task, idWidth, right, role)
}

// markedRow is row with its mark given. A title too long for its column is
// cut to end in an ellipsis, dropping an escape that crosses the cut whole.
// When right would leave the title fewer than minTitle columns, right takes
// a line of its own beneath, still flush to the row's end.
func (text textWriter) markedRow(prefix string, symbol display.Symbol, task issues.OverviewTask, idWidth int, right string, role display.Role) []string {
	id, title, right := cli.EscapeControlText(task.ID), cli.EscapeControlText(task.Title), cli.EscapeControlText(right)
	line := prefix + string(symbol) + " " + text.renderer.Colorize(display.RoleID, id) + strings.Repeat(" ", max(idWidth-len([]rune(id)), 0)) + "  "
	width := rowWidth - display.VisibleWidth(line)
	rightWidth := display.VisibleWidth(right)
	if width-2-rightWidth >= minTitle {
		return []string{line + display.PadRight(cli.CutEscapedText(title, width-2-rightWidth), width-2-rightWidth) + "  " + text.renderer.Colorize(role, right)}
	}
	return []string{line + cli.CutEscapedText(title, width),
		strings.Repeat(" ", max(rowWidth-rightWidth, 0)) + text.renderer.Colorize(role, right)}
}

func interruptedDirectory(entry Entry) string {
	if entry.Kept != nil {
		return filepath.Base(entry.Kept.Worktree)
	}
	return filepath.Base(entry.Working.Worktree)
}
