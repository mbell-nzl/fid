package show

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

const (
	boxWidth = 72
	// rowWidth is how far a row runs: to two columns inside the box.
	rowWidth    = boxWidth - 2
	minTitle    = 10
	progressBar = 10
	indent      = "  "
	// headingIndent is where the heading's title begins, after its mark.
	headingIndent = "    "
	// valueIndent is where a labelled row's value begins: the label is 10
	// columns wide, then two spaces.
	valueIndent = "              "
)

// WriteText writes the result as a box headed by the issue's ID and kind:
// the heading, the labelled rows that have something to show, the
// acceptance criteria and the children. Text fid didn't write has its
// control characters escaped before it is fitted to a column or coloured.
func WriteText(w io.Writer, result Result, renderer display.Renderer) error {
	text := textWriter{renderer: renderer, idWidth: idWidth(result)}
	blocks := [][]string{text.heading(result)}
	if rows := text.rows(result); len(rows) > 0 {
		blocks = append(blocks, rows)
	}
	if len(result.Issue.AcceptanceCriteria) > 0 {
		blocks = append(blocks, text.criteria(result.Issue.AcceptanceCriteria))
	}
	if len(result.Children) > 0 {
		blocks = append(blocks, text.children(result))
	}
	lines := []string{display.TopBorder(renderer.Colorize(display.RoleID, cli.EscapeControlText(result.Issue.ID)), result.Issue.Kind, boxWidth)}
	for _, block := range blocks {
		lines = append(append(lines, ""), block...)
	}
	lines = append(lines, display.BottomBorder(boxWidth))
	if _, err := io.WriteString(w, strings.Join(lines, "\n")+"\n"); err != nil {
		return fmt.Errorf("write show: %w", err)
	}
	return nil
}

type textWriter struct {
	renderer display.Renderer
	idWidth  int
}

// heading is the issue's mark, its title and its state flush right. The
// title wraps and is never cut. When the state would leave it fewer than
// minTitle columns, or its first word doesn't fit beside the state, the
// state takes a line of its own beneath.
func (text textWriter) heading(result Result) []string {
	symbol, role := mark(result.Issue.Kind, result.Issue.StatusCategory)
	title, state := cli.EscapeControlText(result.Issue.Title), cli.EscapeControlText(headingState(result))
	width := rowWidth - len(headingIndent)
	first := width - 2 - display.VisibleWidth(state)
	titles := display.Wrap(title, first, width)
	beside := first >= minTitle && display.VisibleWidth(titles[0]) <= first
	if !beside {
		titles = display.Wrap(title, width, width)
	}
	lines := []string{}
	for index, line := range titles {
		prefix := headingIndent
		if index == 0 {
			prefix = indent + text.renderer.Colorize(role, string(symbol)) + " "
		}
		line = text.renderer.Colorize(display.RoleHeader, line)
		if index == 0 && beside {
			line += strings.Repeat(" ", first-display.VisibleWidth(titles[0])+2) + text.renderer.Colorize(display.RoleMuted, state)
		}
		lines = append(lines, prefix+line)
	}
	if !beside {
		lines = append(lines, strings.Repeat(" ", max(rowWidth-display.VisibleWidth(state), 0))+text.renderer.Colorize(display.RoleMuted, state))
	}
	return lines
}

// headingState is an open epic's reading, as fid status reads it, or the
// issue's workflow state.
func headingState(result Result) string {
	if result.State != nil && (*result.State == issues.EpicInProgress || *result.State == issues.EpicToDo) {
		return strings.ReplaceAll(*result.State, "-", " ")
	}
	return display.StateText(result.Issue.Status)
}

// rows are the labelled rows that have something to show, in order.
func (text textWriter) rows(result Result) []string {
	issue, lines := result.Issue, []string{}
	if result.Task != nil {
		lines = append(lines, text.reference(text.label("Task"), *result.Task, display.StateText(result.Task.Status), display.RoleMuted)...)
	}
	if epic := result.Epic; epic != nil {
		lines = append(lines, text.reference(text.label("Epic"), epic.IssueReference, progressText(epic.Progress), display.RoleMuted)...)
	}
	if issue.Kind == "epic" && len(result.Children) > 0 {
		lines = append(lines, text.label("Progress")+display.ProgressBar(result.Progress.Done, result.Progress.Total, progressBar)+" tasks done")
	}
	if result.Scheduling != nil {
		scheduling := "by priority"
		if *result.Scheduling == "sequential" {
			scheduling = "one at a time, in sequence"
		}
		lines = append(lines, text.value("Scheduling", scheduling)...)
	}
	if issue.Priority != nil {
		lines = append(lines, text.label("Priority")+"P"+strconv.Itoa(*issue.Priority))
	}
	lines = append(lines, text.references("Blocked by", open(result.BlockedBy))...)
	lines = append(lines, text.references("Blocks", open(result.Blocks))...)
	if len(issue.Labels) > 0 {
		lines = append(lines, text.value("Labels", strings.Join(issue.Labels, ", "))...)
	}
	if issue.SpecMDPath != nil {
		lines = append(lines, text.value("Spec", *issue.SpecMDPath)...)
	}
	if hold := result.Working; hold != nil {
		clock := func(stamp string) string { return display.Clock(stamp, result.Now) }
		value := fmt.Sprintf("here, since %s, until %s", clock(hold.Since), clock(hold.ExpiresAt))
		if !hold.Here {
			value = fmt.Sprintf("in %s by %s, since %s, until %s", hold.Worktree, hold.Actor, clock(hold.Since), clock(hold.ExpiresAt))
		}
		if hold.Lapsed {
			value = fmt.Sprintf("⚠ interrupted in %s by %s: work stopped at %s without a park", hold.Worktree, hold.Actor, clock(hold.Since))
		}
		lines = append(lines, text.value("Working", value)...)
	}
	return lines
}

func (text textWriter) label(name string) string {
	return indent + text.renderer.Colorize(display.RoleLabel, display.PadRight(name, len(valueIndent)-len(indent)-2)) + "  "
}

// value is a labelled row's text, wrapped under the value column.
func (text textWriter) value(name, value string) []string {
	lines := []string{}
	for index, line := range display.Wrap(cli.EscapeControlText(value), rowWidth-len(valueIndent), rowWidth-len(valueIndent)) {
		prefix := valueIndent
		if index == 0 {
			prefix = text.label(name)
		}
		lines = append(lines, prefix+line)
	}
	return lines
}

// references is a labelled row with one issue per line, under the value
// column.
func (text textWriter) references(name string, refs []issues.IssueReference) []string {
	lines := []string{}
	for index, ref := range refs {
		prefix := valueIndent
		if index == 0 {
			prefix = text.label(name)
		}
		lines = append(lines, text.reference(prefix, ref, display.StateText(ref.Status), display.RoleMuted)...)
	}
	return lines
}

func (text textWriter) criteria(criteria []issues.ReadAcceptanceCriterion) []string {
	lines := []string{indent + text.renderer.Colorize(display.RoleHeader, "ACCEPTANCE CRITERIA")}
	width := 0
	for _, criterion := range criteria {
		width = max(width, len(strconv.Itoa(criterion.Ordinal)))
	}
	start := len(indent) + width + 2
	for _, criterion := range criteria {
		for index, line := range display.Wrap(cli.EscapeControlText(criterion.Text), rowWidth-start, rowWidth-start) {
			prefix := strings.Repeat(" ", start)
			if index == 0 {
				prefix = fmt.Sprintf("%s%*d  ", indent, width, criterion.Ordinal)
			}
			lines = append(lines, prefix+line)
		}
	}
	return lines
}

// children lists an epic's tasks or a task's subtasks, each with its state
// or, while it's unstarted and blocked, its first open blocker.
func (text textWriter) children(result Result) []string {
	header := text.renderer.Colorize(display.RoleHeader, "TASKS")
	if result.Issue.Kind == "task" {
		count := progressText(*result.Progress)
		header = text.renderer.Colorize(display.RoleHeader, "SUBTASKS") +
			strings.Repeat(" ", rowWidth-len(indent)-len("SUBTASKS")-len(count)) + text.renderer.Colorize(display.RoleMuted, count)
	}
	lines := []string{indent + header}
	for _, child := range result.Children {
		right, role := display.StateText(child.Status), display.RoleMuted
		if child.StatusCategory == "unstarted" && len(child.BlockedBy) > 0 {
			right, role = display.BlockedText(child.BlockedBy), display.RoleWarning
		}
		lines = append(lines, text.reference(indent, child.IssueReference, right, role)...)
	}
	return lines
}

// reference is a line naming another issue after prefix: its mark, its ID
// and its title, with right set flush to the row's end, as fid status sets
// its rows. A title too long for its column is cut to end in an ellipsis,
// dropping an escape that crosses the cut whole. When right would leave the
// title fewer than minTitle columns, right takes a line of its own beneath.
func (text textWriter) reference(prefix string, ref issues.IssueReference, right string, role display.Role) []string {
	symbol, symbolRole := mark(ref.Kind, ref.StatusCategory)
	id, title, right := cli.EscapeControlText(ref.ID), cli.EscapeControlText(ref.Title), cli.EscapeControlText(right)
	line := prefix + text.renderer.Colorize(symbolRole, string(symbol)) + " " + text.renderer.Colorize(display.RoleID, id) +
		strings.Repeat(" ", text.idWidth-len([]rune(id))) + "  "
	width := rowWidth - display.VisibleWidth(line)
	rightWidth := display.VisibleWidth(right)
	if width-2-rightWidth >= minTitle {
		return []string{line + display.PadRight(cli.CutEscapedText(title, width-2-rightWidth), width-2-rightWidth) + "  " + text.renderer.Colorize(role, right)}
	}
	return []string{line + cli.CutEscapedText(title, width),
		strings.Repeat(" ", max(rowWidth-rightWidth, 0)) + text.renderer.Colorize(role, right)}
}

// mark is the symbol for an issue and its colour: ◆ for an epic, otherwise
// by its state's category. Cancelled isn't an error, so its mark is muted.
func mark(kind, category string) (display.Symbol, display.Role) {
	switch {
	case kind == "epic":
		return display.SymbolEpic, display.RoleNormal
	case category == "completed":
		return display.SymbolSuccess, display.RoleSuccess
	case category == "cancelled":
		return display.SymbolError, display.RoleMuted
	case category == "started":
		return display.SymbolStarted, display.RoleNormal
	}
	return display.SymbolUnstarted, display.RoleNormal
}

// open leaves out the issues that are completed or cancelled.
func open(refs []issues.IssueReference) []issues.IssueReference {
	listed := []issues.IssueReference{}
	for _, ref := range refs {
		if ref.StatusCategory != "completed" && ref.StatusCategory != "cancelled" {
			listed = append(listed, ref)
		}
	}
	return listed
}

func progressText(progress issues.Progress) string {
	return fmt.Sprintf("%d/%d done", progress.Done, progress.Total)
}

// idWidth is how wide the longest ID the text lists is, once escaped.
func idWidth(result Result) int {
	refs := append(open(result.BlockedBy), open(result.Blocks)...)
	if result.Task != nil {
		refs = append(refs, *result.Task)
	}
	if result.Epic != nil {
		refs = append(refs, result.Epic.IssueReference)
	}
	for _, child := range result.Children {
		refs = append(refs, child.IssueReference)
	}
	width := 0
	for _, ref := range refs {
		width = max(width, len([]rune(cli.EscapeControlText(ref.ID))))
	}
	return width
}
