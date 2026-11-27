package reviewwork

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
)

// WriteText uses the same box, escaping and wrapping as status and show.
func WriteText(w io.Writer, task string, views []RoundView, renderer display.Renderer, now time.Time) error {
	lines := []string{display.TopBorder(renderer.Colorize(display.RoleHeader, "fid review "+cli.EscapeControlText(task)), now.Format("2006-01-02"), 72), ""}
	prose := func(indent, text string) {
		for _, line := range display.Wrap(cli.EscapeControlText(text), 70-len(indent), 70-len(indent)) {
			lines = append(lines, indent+line)
		}
	}
	if len(views) == 0 {
		prose("  ", "No reviews recorded.")
	}
	for i, v := range views {
		if i > 0 {
			lines = append(lines, "")
		}
		if v.Stream == "judge" {
			prose("  ", v.JudgeText())
			continue
		}
		state, model := "stale", v.Reviewer.Model
		if v.Current {
			state = "current"
		}
		if model == "" {
			model = "unknown model"
		}
		prose("  ", fmt.Sprintf("%s · %s · %s · %s", v.ID, v.Stream, v.Verdict, state))
		prose("    ", v.Reviewer.Actor+" · "+model+fmt.Sprintf(" · other model: %t", v.OtherModel))
		prose("    ", "Session "+v.Reviewer.SessionID+" · profiles: "+strings.Join(v.Profiles, ", "))
		prose("    ", v.Binding.Algorithm+" "+v.Binding.Fingerprint)
		for _, f := range v.Findings {
			state := "open"
			if !f.Open {
				state = "closed"
			}
			prose("    ", f.ID+" · "+f.Severity+" · "+state)
			if f.Closure != nil {
				by := f.Closure.RoundID
				if by == "" {
					by = f.Closure.DecisionID
					if by == "" {
						by = f.Closure.ResponseID + " → " + f.Closure.DeferredTaskID
					}
				}
				prose("      ", "By "+by+" · "+f.Closure.Disposition+" · "+f.Closure.Reason)
			}
			for _, r := range f.Responses {
				if r.DeferredTask == nil {
					prose("      ", "Answer by "+r.Actor+" · "+r.Text)
				}
			}
			prose("      ", f.Summary)
			if f.Location != nil {
				prose("      ", fmt.Sprintf("%s:%d", f.Location.File, f.Location.Line))
			}
		}
	}
	lines = append(lines, display.BottomBorder(72))
	_, err := io.WriteString(w, strings.Join(lines, "\n")+"\n")
	return err
}
