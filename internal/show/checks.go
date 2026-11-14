package show

import (
	"strings"

	"github.com/mbell-nzl/fid/internal/display"
)

func (text textWriter) checks(result Result) []string {
	lines := []string{}
	for _, record := range result.Checks {
		message := "○ " + record.Name + " required · not checked yet"
		if record.Start != nil {
			message = "⚠ " + record.Name + " never finished (" + record.Start.ID + ")"
		}
		if check := record.Result; check != nil {
			mark := "✗ "
			if check.Outcome == "passed" {
				mark = "✓ "
				if !check.Counts {
					mark = "⚠ "
				}
			}
			message = mark + record.Name + " " + check.Outcome + " at " + display.Clock(check.FinishedAt, result.Now)
			if !check.Counts && len(check.Reasons) > 0 {
				reasons := []string{}
				for _, reason := range check.Reasons {
					reasons = append(reasons, map[string]string{"tree-dirty": "the tree had uncommitted changes", "tree-changed": "the tree changed", "ownership-lost": "ownership was lost", "artifact-not-kept": "a listed file wasn't kept"}[reason])
				}
				message += ", doesn't count: " + strings.Join(reasons, ", ")
			}
		}
		if record.Start != nil {
			if record.Required != nil && *record.Required {
				message += " · required"
			}
			if record.Satisfied != nil {
				if *record.Satisfied {
					message += " · counts on HEAD"
				} else {
					message += " · no counting pass on HEAD"
				}
			}
		} else if record.Required == nil {
			message = "○ " + record.Name + " · not checked yet"
		}
		row := text.value("Checks", message)
		if len(lines) > 0 {
			row[0] = valueIndent + strings.TrimPrefix(row[0], text.label("Checks"))
		}
		lines = append(lines, row...)
	}
	if len(result.Warnings) > 0 {
		summaries := []string{}
		for _, warning := range result.Warnings {
			summaries = append(summaries, warning.Summary)
		}
		row := text.value("Checks", "⚠ "+strings.Join(summaries, "; "))
		if len(lines) > 0 {
			row[0] = valueIndent + strings.TrimPrefix(row[0], text.label("Checks"))
		}
		lines = append(lines, row...)
	}
	return lines
}
