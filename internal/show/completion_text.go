package show

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/display"
)

func (text textWriter) completion(result Result) []string {
	c := result.Completion
	lines := []string{}
	add := func(value string) {
		row := text.value("Completed", value)
		if len(lines) > 0 {
			row[0] = valueIndent + strings.TrimPrefix(row[0], text.label("Completed"))
		}
		lines = append(lines, row...)
	}
	if c.Record != nil {
		add(display.Clock(c.At, result.Now) + " · " + c.ID + " · " + c.Actor)
	}
	if m := c.Merge; m != nil {
		add(fmt.Sprintf("PR #%d merged as %s · %s", m.Number, short(m.MergeCommit), display.Clock(m.MergedAt, result.Now)))
	}
	for _, ci := range c.CIResults {
		mark := "✗"
		if ci.Conclusion == "success" {
			mark = "✓"
		}
		add(fmt.Sprintf("CI %s %s · #%d on %s", mark, ci.Conclusion, ci.RunID, short(ci.HeadSHA)))
	}
	if h := c.Handoff; h != nil {
		message := "Hand-off " + h.ID + " at " + short(h.Head)
		if c.Record != nil && c.Fingerprint.Fingerprint != "" {
			message += " · fingerprint " + short(c.Fingerprint.Fingerprint)
		}
		add(message)
	} else if c.Record != nil && c.Fingerprint.Fingerprint != "" {
		add("fingerprint " + short(c.Fingerprint.Fingerprint))
	}
	if w := c.Waiver; w != nil {
		add("Waiver to " + w.Binding.Lane + " · " + w.View().Text("waived an earlier change"))
	}
	for _, round := range c.Rounds {
		model := round.Reviewer.Model
		if model == "" {
			model = "model not recorded"
		}
		mark := "✓"
		if round.Verdict != "APPROVE" {
			mark = "✗"
		}
		add(fmt.Sprintf("Rounds %s %s %s %s · %s (%s)", mark, round.ID, round.Stream, round.Verdict, round.Reviewer.Actor, model))
	}
	if a := c.Approval; a != nil {
		add("Approval ✓ " + a.ID + " " + a.Actor + " · " + display.Clock(a.At, result.Now))
	} else if c.Waiver != nil {
		add("Approval " + (decisionview.View{State: "not-required"}).Text(""))
	}
	if c.Record != nil && c.NoOtherModelReview {
		add("no other-model review")
	}
	if c.assured() {
		rollback := "absent"
		if c.RollbackPresent {
			rollback = "present"
		}
		add("Rollback " + rollback)
	}
	for _, missing := range c.Missing {
		add("⚠ " + missing)
	}
	return lines
}

func short(value string) string { return value[:min(8, len(value))] }
