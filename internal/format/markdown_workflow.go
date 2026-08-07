package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func writeWorkflowListMarkdown(writer io.Writer, items []workflow.ListItem) error {
	var builder strings.Builder
	if len(items) == 0 {
		builder.WriteString("[]\n")
	} else {
		builder.WriteString("# Workflows\n\n")
		builder.WriteString("| id | name | applies_to | active | states | transitions | assigned_to |\n")
		builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | --- |\n")
		for _, item := range items {
			assignments, err := json.Marshal(item.AssignedTo)
			if err != nil {
				return fmt.Errorf("marshal workflow assignments: %w", err)
			}
			fmt.Fprintf(&builder, "| %s | %s | %s | %t | %d | %d | %s |\n",
				markdownTableCell(item.ID), markdownTableCell(item.Name), workflowAppliesToMarkdown(item.AppliesTo),
				item.Active, item.States, item.Transitions, markdownTableCell(string(assignments)))
		}
	}
	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func workflowAppliesToMarkdown(appliesTo *string) string {
	if appliesTo == nil {
		return "null"
	}
	return markdownTableCell(*appliesTo)
}
