package format

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/stats"
)

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

func writeTokensStatsMarkdown(writer io.Writer, rows []stats.TokensRow) error {
	var builder strings.Builder

	builder.WriteString("# Token Stats\n\n")
	builder.WriteString("| Group | Total Tokens In | Total Tokens Out | Run Count |\n")
	builder.WriteString("| --- | ---: | ---: | ---: |\n")
	for _, row := range rows {
		fmt.Fprintf(&builder, "| %s | %d | %d | %d |\n", markdownStatsGroupLabel(row.Role, row.Model, row.EpicID), row.TotalTokensIn, row.TotalTokensOut, row.RunCount)
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeCostStatsMarkdown(writer io.Writer, rows []stats.CostRow) error {
	var builder strings.Builder

	builder.WriteString("# Cost Stats\n\n")
	builder.WriteString("| Group | Total Cost (USD) | Total Tokens In | Total Tokens Out | Run Count |\n")
	builder.WriteString("| --- | ---: | ---: | ---: | ---: |\n")
	for _, row := range rows {
		fmt.Fprintf(&builder, "| %s | %s | %d | %d | %d |\n", markdownStatsGroupLabel(row.Role, row.Model, row.EpicID), markdownStatsCost(row.TotalCostUSD), row.TotalTokensIn, row.TotalTokensOut, row.RunCount)
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func markdownStatsGroupLabel(role, model, epicID *string) string {
	switch {
	case role != nil && strings.TrimSpace(*role) != "":
		return *role
	case model != nil && strings.TrimSpace(*model) != "":
		return *model
	case epicID != nil && strings.TrimSpace(*epicID) != "":
		return *epicID
	default:
		return "(none)"
	}
}

func markdownStatsCost(value float64) string {
	text := strconv.FormatFloat(value, 'f', 6, 64)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" {
		return "0"
	}
	return text
}
