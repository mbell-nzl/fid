package reviewcomment

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func renderRoundV2(r reviewrecord.Round, number int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Review: %s\n\n**Reviewer** %s", r.Verdict, codeSpan(r.Reviewer.Actor))
	attribution := []string{}
	if r.Reviewer.Model != "" {
		attribution = append(attribution, codeSpan(r.Reviewer.Model))
	}
	if r.Reviewer.Effort != "" {
		attribution = append(attribution, "effort "+r.Reviewer.Effort)
	}
	if len(attribution) > 0 {
		fmt.Fprintf(&b, " (%s)", strings.Join(attribution, ", "))
	}
	fmt.Fprintf(&b, " · **Profiles** %s · **Round** %d  \n**Session** %s · **Stream** %s · %s  \n**Claim** %s\n\n### Findings\n\n", strings.Join(r.Profiles, ", "), number, codeSpan(r.Reviewer.SessionID), r.Stream, codeSpan(r.ID), r.Purpose)
	if len(r.Findings) == 0 {
		b.WriteString("None.\n")
	}
	for _, f := range r.Findings {
		action := map[string]string{"major": "fix before merge", "minor": "fix or answer", "nit": "optional"}[f.Severity]
		fmt.Fprintf(&b, "- **%s** · **%s** · %s — %s  \n", codeSpan(f.ID), f.Severity, action, f.Summary)
		if f.Location != nil {
			fmt.Fprintf(&b, "  **Location** %s  \n", codeSpan(fmt.Sprintf("%s:%d", f.Location.File, f.Location.Line)))
		}
		fmt.Fprintf(&b, "  **Evidence** %s  \n  **Fix** %s\n", f.Evidence, f.Fix)
	}
	if len(r.Dispositions) > 0 {
		b.WriteString("\n### Earlier findings\n\n")
		for _, d := range r.Dispositions {
			fmt.Fprintf(&b, "- %s · **%s**: %s\n", codeSpan(d.FindingID), d.Disposition, d.Reason)
		}
	}
	fmt.Fprintf(&b, "\n### Verification\n\n%s\n", r.Verification)
	return b.String()
}

// codeSpan fences literal backticks and pads boundaries that CommonMark would
// otherwise parse as part of the fence or strip from the value.
func codeSpan(value string) string {
	longest, run := 0, 0
	for _, c := range value {
		if c == '`' {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	fence := strings.Repeat("`", longest+1)
	if strings.HasPrefix(value, "`") || strings.HasSuffix(value, "`") ||
		(strings.HasPrefix(value, " ") && strings.HasSuffix(value, " ") && strings.Trim(value, " ") != "") {
		value = " " + value + " "
	}
	return fence + value + fence
}
