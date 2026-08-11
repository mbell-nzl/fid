package format

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/issues"
)

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
