package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/maintenance"
)

type MaintenanceCheckView struct {
	Status string `json:"status"`
	maintenance.Result
}

func writeMaintenanceCheckMarkdown(writer io.Writer, view MaintenanceCheckView) error {
	var builder strings.Builder

	builder.WriteString("# Maintenance Check\n\n")
	fmt.Fprintf(&builder, "**Status:** %s\n", view.Status)
	writeMaintenanceFindingsMarkdown(&builder, "Errors", view.Errors, "(no errors)")
	writeMaintenanceFindingsMarkdown(&builder, "Warnings", view.Warnings, "(no warnings)")
	builder.WriteString("\n## Info\n\n")
	fmt.Fprintf(&builder, "- Issue count: %d\n", view.Info.IssueCount)
	fmt.Fprintf(&builder, "- Live claim count: %d\n", view.Info.ClaimCountLive)
	fmt.Fprintf(&builder, "- Recent activity (7d): %d\n", view.Info.RecentActivity)

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func writeMaintenanceFindingsMarkdown(builder *strings.Builder, heading string, findings []maintenance.Finding, empty string) {
	builder.WriteString("\n## ")
	builder.WriteString(heading)
	builder.WriteString("\n\n")
	if len(findings) == 0 {
		builder.WriteString(empty)
		builder.WriteString("\n")
		return
	}
	for _, finding := range findings {
		fmt.Fprintf(builder, "- `%s` %s: %s\n", finding.Check, finding.EntityID, finding.Message)
	}
}
