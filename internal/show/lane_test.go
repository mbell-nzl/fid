package show

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestLaneTextExplainsOnlyTheRecordedRaise(t *testing.T) {
	t.Parallel()
	assured, standard, routine := "assured", "standard", "routine"
	for _, tc := range []struct {
		name    string
		handoff *handoffrecord.View
		want    string
	}{
		{"no hand-off", nil, "Assured"},
		{"historical", &handoffrecord.View{}, "Assured"},
		{"not raised", &handoffrecord.View{Floor: &policy.Floor{Lane: &assured, Paths: []string{"secure/a.go"}}}, "Assured"},
		{"unknown floor", &handoffrecord.View{RaisedFrom: &routine}, "Assured"},
		{"no setting paths", &handoffrecord.View{RaisedFrom: &routine, Floor: &policy.Floor{}}, "Assured"},
		{"empty explanation", &handoffrecord.View{RaisedFrom: &routine, Floor: &policy.Floor{Lane: &assured, Paths: []string{}}}, "Assured"},
		{"lane raised since hand-off", &handoffrecord.View{RaisedFrom: &routine, Floor: &policy.Floor{Lane: &standard, Paths: []string{"main.go"}}}, "Assured"},
		{"one setting path", &handoffrecord.View{RaisedFrom: &standard, Floor: &policy.Floor{Lane: &assured, Paths: []string{"secure/a.go"}}}, "Assured · floor: secure/a.go · raised from Standard"},
		{"two setting paths", &handoffrecord.View{RaisedFrom: &routine, Floor: &policy.Floor{Lane: &assured, Paths: []string{"secure/a.go", "secure/b.go"}}}, "Assured · floor: secure/a.go and 1 others · raised from Routine"},
		{"three setting paths", &handoffrecord.View{RaisedFrom: &routine, Floor: &policy.Floor{Lane: &assured, Paths: []string{"secure/a.go", "secure/b.go", "secure/c.go"}}}, "Assured · floor: secure/a.go and 2 others · raised from Routine"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := textTask()
			result.Issue.Lane = &assured
			if tc.handoff != nil {
				view := *tc.handoff
				view.At, view.Head, view.PR = "2026-11-22T00:00:00Z", strings.Repeat("a", 40), "https://github.com/owner/repo/pull/98"
				result.Handoff = &view
			}
			text := strings.Join(strings.Fields(writeText(t, result, display.Renderer{})), " ")
			if !strings.Contains(text, "Lane "+tc.want+" ") || strings.Contains(text, "floor:") != strings.Contains(tc.want, "floor:") {
				t.Fatalf("Lane row does not explain its hand-off: want %q in %s", tc.want, text)
			}
		})
	}
}
