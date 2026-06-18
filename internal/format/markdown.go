package format

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/issues"
)

type IssueMarkdownView struct {
	Issue issues.Issue
}

type IssueListMarkdownView struct {
	Items []issues.ListItem
}

func writeMarkdown(writer io.Writer, value any) error {
	switch typed := value.(type) {
	case ReadyView:
		return writeReadyMarkdown(writer, typed)
	case *ReadyView:
		if typed == nil {
			return fmt.Errorf("markdown ready view is nil")
		}
		return writeReadyMarkdown(writer, *typed)
	case IssueMarkdownView:
		return writeIssueMarkdown(writer, typed)
	case *IssueMarkdownView:
		if typed == nil {
			return fmt.Errorf("markdown issue view is nil")
		}
		return writeIssueMarkdown(writer, *typed)
	case issues.Issue:
		return writeIssueMarkdown(writer, IssueMarkdownView{Issue: typed})
	case *issues.Issue:
		if typed == nil {
			return fmt.Errorf("markdown issue is nil")
		}
		return writeIssueMarkdown(writer, IssueMarkdownView{Issue: *typed})
	case IssueListMarkdownView:
		return writeIssueListMarkdown(writer, typed)
	case *IssueListMarkdownView:
		if typed == nil {
			return fmt.Errorf("markdown issue list view is nil")
		}
		return writeIssueListMarkdown(writer, *typed)
	case issues.AcceptanceCriteriaList:
		return writeAcceptanceCriteriaListMarkdown(writer, typed)
	case *issues.AcceptanceCriteriaList:
		if typed == nil {
			return fmt.Errorf("markdown acceptance criteria list is nil")
		}
		return writeAcceptanceCriteriaListMarkdown(writer, *typed)
	case []issues.ListItem:
		return writeIssueListMarkdown(writer, IssueListMarkdownView{Items: typed})
	default:
		return fmt.Errorf("unsupported markdown value type %T", value)
	}
}

func writeReadyMarkdown(writer io.Writer, view ReadyView) error {
	var builder strings.Builder

	builder.WriteString("# Ready Task\n\n")
	if view.Task == nil {
		builder.WriteString("(no ready task)\n")
	} else {
		builder.WriteString("- ID: ")
		builder.WriteString(view.Task.ID)
		builder.WriteString("\n")
		builder.WriteString("- Title: ")
		builder.WriteString(view.Task.Title)
		builder.WriteString("\n")
		builder.WriteString("- Status: ")
		builder.WriteString(view.Task.Status)
		builder.WriteString("\n")
		builder.WriteString("- Epic: ")
		builder.WriteString(markdownReadyEpic(view.Task))
		builder.WriteString("\n")
		builder.WriteString("- Sequence: ")
		builder.WriteString(markdownReadyOptionalInt(view.Task.Sequence))
		builder.WriteString("\n")
		builder.WriteString("- Priority: ")
		builder.WriteString(markdownPriority(view.Task.Priority))
		builder.WriteString("\n")
		builder.WriteString("- Est. tokens: ")
		builder.WriteString(markdownReadyOptionalInt(view.Task.EstTokens))
		builder.WriteString("\n")
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeIssueMarkdown(writer io.Writer, view IssueMarkdownView) error {
	var builder strings.Builder

	builder.WriteString("# ")
	builder.WriteString(view.Issue.ID)
	builder.WriteString(": ")
	builder.WriteString(view.Issue.Title)
	builder.WriteString("\n\n")

	builder.WriteString("**Kind:** ")
	builder.WriteString(view.Issue.Kind)
	builder.WriteString(" | **Status:** ")
	builder.WriteString(view.Issue.Status)
	builder.WriteString(" | **Priority:** ")
	builder.WriteString(markdownPriority(view.Issue.Priority))
	builder.WriteString("\n")

	if line := markdownEpicLine(view.Issue); line != "" {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	if line := markdownOwnerSequenceLine(view.Issue); line != "" {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	if view.Issue.SpecMDPath != nil && strings.TrimSpace(*view.Issue.SpecMDPath) != "" {
		builder.WriteString("**Spec:** ")
		builder.WriteString(*view.Issue.SpecMDPath)
		builder.WriteString("\n")
	}

	writeAcceptanceCriteriaMarkdown(&builder, view.Issue.AcceptanceCriteria)
	writeLabelsMarkdown(&builder, view.Issue.Labels)
	writeLinksMarkdown(&builder, view.Issue.Links)
	writeFooterMarkdown(&builder, view.Issue)

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func markdownPriority(priority *int) string {
	if priority == nil {
		return "(none)"
	}
	return fmt.Sprintf("P%d", *priority)
}

func markdownEpicLine(issue issues.Issue) string {
	if issue.ParentID == nil || strings.TrimSpace(*issue.ParentID) == "" {
		return ""
	}

	line := "**Epic:** " + *issue.ParentID
	if issue.ParentTitle != nil && strings.TrimSpace(*issue.ParentTitle) != "" {
		line += " (" + *issue.ParentTitle + ")"
	}
	return line
}

func markdownReadyEpic(task *issues.ReadyTask) string {
	if task == nil || task.EpicID == nil || strings.TrimSpace(*task.EpicID) == "" {
		return "(none)"
	}

	line := *task.EpicID
	if task.EpicTitle != nil && strings.TrimSpace(*task.EpicTitle) != "" {
		line += " (" + *task.EpicTitle + ")"
	}
	return line
}

func markdownReadyOptionalInt(value *int) string {
	if value == nil {
		return "(none)"
	}
	return strconv.Itoa(*value)
}

func markdownOwnerSequenceLine(issue issues.Issue) string {
	parts := make([]string, 0, 2)
	if issue.OwnerRole != nil && strings.TrimSpace(*issue.OwnerRole) != "" {
		parts = append(parts, "**Owner:** "+*issue.OwnerRole)
	}
	if issue.Sequence != nil {
		parts = append(parts, fmt.Sprintf("**Sequence:** %d", *issue.Sequence))
	}
	return strings.Join(parts, " | ")
}

func writeAcceptanceCriteriaMarkdown(builder *strings.Builder, criteria []issues.AcceptanceCriterion) {
	if len(criteria) == 0 {
		return
	}

	sorted := slices.Clone(criteria)
	sort.Slice(sorted, func(left, right int) bool {
		if sorted[left].Ordinal != sorted[right].Ordinal {
			return sorted[left].Ordinal < sorted[right].Ordinal
		}
		return sorted[left].ID < sorted[right].ID
	})

	builder.WriteString("\n## Acceptance Criteria\n\n")
	writeAcceptanceCriteriaMarkdownItems(builder, sorted)
}

func writeAcceptanceCriteriaMarkdownItems(builder *strings.Builder, criteria []issues.AcceptanceCriterion) {
	for _, criterion := range criteria {
		builder.WriteString("- [")
		if criterionChecked(criterion.LatestVerifications) {
			builder.WriteString("x")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString("] ")
		builder.WriteString(criterion.ID)
		builder.WriteString(": ")
		builder.WriteString(criterion.Text)
		builder.WriteString("\n")

		for _, line := range markdownVerificationLines(criterion.LatestVerifications) {
			builder.WriteString("  - ")
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	}
}

func writeAcceptanceCriteriaListMarkdown(writer io.Writer, list issues.AcceptanceCriteriaList) error {
	var builder strings.Builder

	builder.WriteString("# Acceptance Criteria\n\n")
	if len(list.Criteria) == 0 {
		builder.WriteString("(no acceptance criteria)\n")
	} else {
		writeAcceptanceCriteriaMarkdownItems(&builder, list.Criteria)
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func criterionChecked(latest map[string]any) bool {
	for _, line := range markdownVerificationLines(latest) {
		if strings.Contains(line, ": pass") {
			return true
		}
	}
	return false
}

func markdownVerificationLines(latest map[string]any) []string {
	if len(latest) == 0 {
		return nil
	}

	roles := make([]string, 0, len(latest))
	for role := range latest {
		roles = append(roles, role)
	}
	sort.Strings(roles)

	lines := make([]string, 0, len(roles))
	for _, role := range roles {
		line := role + ": "
		entry := latest[role]
		if entry == nil {
			lines = append(lines, line+"(not verified)")
			continue
		}

		typed, ok := entry.(map[string]any)
		if !ok {
			lines = append(lines, line+fmt.Sprintf("%v", entry))
			continue
		}

		status := strings.TrimSpace(fmt.Sprintf("%v", typed["status"]))
		if status == "" || status == "<nil>" {
			lines = append(lines, line+"(not verified)")
			continue
		}

		verifiedAt := strings.TrimSpace(fmt.Sprintf("%v", typed["verified_at"]))
		if verifiedAt != "" && verifiedAt != "<nil>" {
			lines = append(lines, line+status+" ("+verifiedAt+")")
			continue
		}
		lines = append(lines, line+status)
	}
	return lines
}

func writeLabelsMarkdown(builder *strings.Builder, labels []string) {
	if len(labels) == 0 {
		return
	}

	builder.WriteString("\n## Labels\n\n")
	builder.WriteString(strings.Join(labels, ", "))
	builder.WriteString("\n")
}

func writeLinksMarkdown(builder *strings.Builder, links []issues.IssueLink) {
	if len(links) == 0 {
		return
	}

	builder.WriteString("\n## Links\n\n")
	for _, link := range links {
		builder.WriteString("- ")
		builder.WriteString(link.Relation)
		builder.WriteString(" ")
		builder.WriteString(link.TargetID)
		if link.TargetTitle != nil && strings.TrimSpace(*link.TargetTitle) != "" {
			builder.WriteString(" (")
			builder.WriteString(*link.TargetTitle)
			builder.WriteString(")")
		}
		builder.WriteString("\n")
	}
}

func writeFooterMarkdown(builder *strings.Builder, issue issues.Issue) {
	createdAt := ""
	updatedAt := ""
	if issue.CreatedAt != nil {
		createdAt = *issue.CreatedAt
	}
	if issue.UpdatedAt != nil {
		updatedAt = *issue.UpdatedAt
	}

	builder.WriteString("\n---\n")
	builder.WriteString("Created: ")
	builder.WriteString(createdAt)
	builder.WriteString(" | Updated: ")
	builder.WriteString(updatedAt)
	builder.WriteString("\n")
}

func writeIssueListMarkdown(writer io.Writer, view IssueListMarkdownView) error {
	var builder strings.Builder

	builder.WriteString("# Issues\n\n")
	for _, item := range view.Items {
		builder.WriteString("- ")
		builder.WriteString(item.ID)
		builder.WriteString(" | ")
		builder.WriteString(item.Kind)
		builder.WriteString(" | ")
		builder.WriteString(item.Status)
		builder.WriteString(" | ")
		builder.WriteString(item.StatusCategory)
		builder.WriteString(" | ")
		builder.WriteString(markdownListEpic(item.ParentID))
		builder.WriteString(" | ")
		builder.WriteString(markdownPriority(item.Priority))
		if item.InInbox {
			builder.WriteString(" | inbox")
		}
		builder.WriteString(" | concerns: ")
		builder.WriteString(strconv.Itoa(item.OpenConcernsCount))
		builder.WriteString("\n")
		builder.WriteString("  ")
		builder.WriteString(item.Title)
		builder.WriteString("\n")
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func markdownListEpic(parentID *string) string {
	if parentID == nil || strings.TrimSpace(*parentID) == "" {
		return "no epic"
	}
	return "epic " + *parentID
}
