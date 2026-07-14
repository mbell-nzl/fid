package format

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/mbell-nzl/fid/internal/links"
	"github.com/mbell-nzl/fid/internal/policy"
)

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

func writeConcernListMarkdown(writer io.Writer, items []issues.ConcernListItem) error {
	var builder strings.Builder

	builder.WriteString("# Concerns\n\n")
	if len(items) == 0 {
		builder.WriteString("(no concerns)\n")
	} else {
		for _, item := range items {
			builder.WriteString("- ")
			builder.WriteString(item.ID)
			builder.WriteString(" | issue ")
			builder.WriteString(item.IssueID)
			builder.WriteString(" | ")
			builder.WriteString(item.Status)
			builder.WriteString(" | by ")
			builder.WriteString(item.RaisedByRole)
			builder.WriteString(" | completion ")
			builder.WriteString(item.RaisedInCompletionID)
			builder.WriteString("\n")
			builder.WriteString("  ")
			builder.WriteString(item.Summary)
			builder.WriteString("\n")
			if item.BodyMDPath != nil && strings.TrimSpace(*item.BodyMDPath) != "" {
				builder.WriteString("  Body: ")
				builder.WriteString(*item.BodyMDPath)
				builder.WriteString("\n")
			}
		}
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writePolicyMatrixMarkdown(writer io.Writer, rows []policy.MatrixRow) error {
	var builder strings.Builder

	builder.WriteString("# Policy\n\n")
	builder.WriteString("| action | developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")

	for _, row := range rows {
		_, _ = fmt.Fprintf(&builder,
			"| %s | %t | %t | %t | %t | %t | %t | %t |\n",
			row.Action,
			row.Developer,
			row.QAReviewer,
			row.QALead,
			row.TLLite,
			row.TL,
			row.WorkflowRunner,
			row.Human,
		)
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeClaimListMarkdown(writer io.Writer, items []claims.ListItem) error {
	var builder strings.Builder

	builder.WriteString("# Claims\n\n")
	for _, item := range items {
		builder.WriteString("- ")
		builder.WriteString(item.RefType)
		builder.WriteString(" ")
		builder.WriteString(item.RefID)
		builder.WriteString(" | ")
		builder.WriteString(item.ClaimedBy)
		builder.WriteString(" | claimed ")
		builder.WriteString(item.ClaimedAt)
		builder.WriteString(" | ttl ")
		builder.WriteString(strconv.Itoa(item.TTLSeconds))
		builder.WriteString("s | expires ")
		builder.WriteString(item.ExpiresAt)
		builder.WriteString(" | stale: ")
		builder.WriteString(strconv.FormatBool(item.IsStale))
		builder.WriteString("\n")
		if item.Note != nil && strings.TrimSpace(*item.Note) != "" {
			builder.WriteString("  ")
			builder.WriteString(*item.Note)
			builder.WriteString("\n")
		}
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeIssueGateListMarkdown(writer io.Writer, items []issues.IssueGate) error {
	var builder strings.Builder

	builder.WriteString("# Gates\n\n")
	if len(items) == 0 {
		builder.WriteString("(no gates)\n")
	} else {
		builder.WriteString("| id | attempt | role | command | exit_code | agent_run_id | ran_at | output_excerpt |\n")
		builder.WriteString("| --- | ---: | --- | --- | ---: | --- | --- | --- |\n")
		for _, item := range items {
			builder.WriteString("| ")
			builder.WriteString(markdownTableCell(item.ID))
			builder.WriteString(" | ")
			builder.WriteString(strconv.Itoa(item.AttemptNumber))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.RecordedByRole))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.Command))
			builder.WriteString(" | ")
			builder.WriteString(strconv.Itoa(item.ExitCode))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(dereferenceMarkdownString(item.AgentRunID)))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.RanAt))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(dereferenceMarkdownString(item.OutputExcerpt)))
			builder.WriteString(" |\n")
		}
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeLabelListMarkdown(writer io.Writer, items []labels.Label) error {
	var builder strings.Builder

	builder.WriteString("# Labels\n\n")
	if len(items) == 0 {
		builder.WriteString("(no labels)\n")
	} else {
		builder.WriteString("| name | colour | description |\n")
		builder.WriteString("| --- | --- | --- |\n")
		for _, item := range items {
			builder.WriteString("| ")
			builder.WriteString(markdownTableCell(item.Name))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(dereferenceMarkdownString(item.Colour)))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(dereferenceMarkdownString(item.Description)))
			builder.WriteString(" |\n")
		}
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeLinkListMarkdown(writer io.Writer, items []links.ListItem) error {
	var builder strings.Builder

	builder.WriteString("# Links\n\n")
	if len(items) == 0 {
		builder.WriteString("(no links)\n")
	} else {
		builder.WriteString("| direction | relation | entity_type | entity_id | entity_title |\n")
		builder.WriteString("| --- | --- | --- | --- | --- |\n")
		for _, item := range items {
			builder.WriteString("| ")
			builder.WriteString(markdownTableCell(item.Direction))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.Relation))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.EntityType))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(item.EntityID))
			builder.WriteString(" | ")
			builder.WriteString(markdownTableCell(dereferenceMarkdownString(item.EntityTitle)))
			builder.WriteString(" |\n")
		}
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

func markdownTableCell(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
}

func dereferenceMarkdownString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
