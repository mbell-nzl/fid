package mergerecord

import (
	"strings"
	"testing"
)

func validMerge() Merge {
	return Merge{ID: "MRG-123456", TaskID: "FID-123456", HandoffID: "HOF-123456", PR: "https://github.com/owner/repo/pull/77", Number: 77, Repository: "owner/repo", BaseBranch: "main", MergeCommit: strings.Repeat("a", 40), HeadSHA: strings.Repeat("b", 40), HeadTree: strings.Repeat("c", 40), MergedAt: "2026-11-15T01:00:00Z", ObservedAt: "2026-11-15T02:00:00Z"}
}

func TestMergeValidation(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		change func(*Merge)
	}{
		{"valid", func(*Merge) {}},
		{"missing task", func(m *Merge) { m.TaskID = "" }},
		{"missing handoff", func(m *Merge) { m.HandoffID = "" }},
		{"wrong PR number", func(m *Merge) { m.Number++ }},
		{"wrong repository", func(m *Merge) { m.Repository = "other/repo" }},
		{"short merge SHA", func(m *Merge) { m.MergeCommit = "abc" }},
		{"nonhex head", func(m *Merge) { m.HeadSHA = strings.Repeat("z", 40) }},
		{"long tree", func(m *Merge) { m.HeadTree = strings.Repeat("a", 64) }},
		{"missing branch", func(m *Merge) { m.BaseBranch = "" }},
		{"missing merge time", func(m *Merge) { m.MergedAt = "" }},
		{"missing observed time", func(m *Merge) { m.ObservedAt = "" }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			m := validMerge()
			row.change(&m)
			if err := m.Validate(); (err == nil) != (row.name == "valid") {
				t.Fatalf("Validate=%v", err)
			}
		})
	}
}

func TestCIResultValidation(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		change func(*CIResult)
	}{
		{"valid", func(*CIResult) {}},
		{"missing task", func(c *CIResult) { c.TaskID = "" }},
		{"missing merge", func(c *CIResult) { c.MergeID = "" }},
		{"zero run", func(c *CIResult) { c.RunID = 0 }},
		{"zero attempt", func(c *CIResult) { c.Attempt = 0 }},
		{"zero workflow", func(c *CIResult) { c.WorkflowID = 0 }},
		{"missing name", func(c *CIResult) { c.Name = "" }},
		{"missing path", func(c *CIResult) { c.Path = "" }},
		{"missing event", func(c *CIResult) { c.Event = "" }},
		{"short head", func(c *CIResult) { c.HeadSHA = "a" }},
		{"pending", func(c *CIResult) { c.Status = "in_progress" }},
		{"unknown status", func(c *CIResult) { c.Status = "green" }},
		{"unknown conclusion", func(c *CIResult) { c.Conclusion = "green" }},
		{"wrong URL", func(c *CIResult) { c.URL = "https://github.com/other/repo/actions/runs/12" }},
		{"missing time", func(c *CIResult) { c.ObservedAt = "" }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			c := CIResult{ID: "CIR-123456", TaskID: "FID-123456", MergeID: "MRG-123456", Repository: "owner/repo", RunID: 12, Attempt: 1, WorkflowID: 3, Name: "CI", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: strings.Repeat("a", 40), Status: "completed", Conclusion: "success", URL: "https://github.com/owner/repo/actions/runs/12", ObservedAt: "2026-11-15T02:00:00Z"}
			row.change(&c)
			if err := c.Validate(); (err == nil) != (row.name == "valid") {
				t.Fatalf("Validate=%v", err)
			}
		})
	}
}
