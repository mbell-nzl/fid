package format

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/issues"
)

func TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue(t *testing.T) {
	view := populatedIssueMarkdownFixture()

	var markdownOutput bytes.Buffer
	if err := Write(&markdownOutput, Markdown, view); err != nil {
		t.Fatalf("Write(Markdown) error = %v", err)
	}

	golden, err := os.ReadFile(toonGoldenPath(t, "issue-populated.golden.md"))
	if err != nil {
		t.Fatalf("ReadFile() golden error = %v", err)
	}

	if markdownOutput.String() != string(golden) {
		t.Fatalf("markdown output = %q, want golden %q", markdownOutput.String(), string(golden))
	}
}

func TestWriteMarkdownOmitsEmptySectionsForEmptyIssue(t *testing.T) {
	view := emptyIssueMarkdownFixture()

	var markdownOutput bytes.Buffer
	if err := Write(&markdownOutput, Markdown, view); err != nil {
		t.Fatalf("Write(Markdown) error = %v", err)
	}

	want := "" +
		"# FID-empty01: Render empty issue\n\n" +
		"**Kind:** task | **Status:** todo | **Priority:** (none)\n" +
		"**Owner:** developer\n\n" +
		"---\n" +
		"Created: 2026-06-07T14:00:00Z | Updated: 2026-06-07T14:00:00Z\n"
	if markdownOutput.String() != want {
		t.Fatalf("markdown output = %q, want %q", markdownOutput.String(), want)
	}
	if strings.Contains(markdownOutput.String(), "## Acceptance Criteria") ||
		strings.Contains(markdownOutput.String(), "## Labels") ||
		strings.Contains(markdownOutput.String(), "## Links") {
		t.Fatalf("markdown output = %q, want empty sections omitted", markdownOutput.String())
	}
}

func TestWriteMarkdownRendersIssueGateLists(t *testing.T) {
	t.Run("populated gate list", func(t *testing.T) {
		var markdownOutput bytes.Buffer
		if err := Write(&markdownOutput, Markdown, []issues.IssueGate{
			{
				ID:             "GATE-000401",
				AttemptNumber:  1,
				AgentRunID:     stringPointer("RUN-000401"),
				RecordedByRole: "developer",
				Command:        "go test ./...",
				ExitCode:       0,
				OutputExcerpt:  stringPointer("ok"),
				RanAt:          "2026-05-27T09:20:00Z",
			},
		}); err != nil {
			t.Fatalf("Write(Markdown) error = %v", err)
		}

		output := markdownOutput.String()
		for _, want := range []string{"# Gates", "GATE-000401", "go test ./...", "RUN-000401", "ok"} {
			if !strings.Contains(output, want) {
				t.Fatalf("markdown output = %q, want substring %q", output, want)
			}
		}
	})

	t.Run("empty gate list", func(t *testing.T) {
		var markdownOutput bytes.Buffer
		if err := Write(&markdownOutput, Markdown, []issues.IssueGate{}); err != nil {
			t.Fatalf("Write(Markdown) error = %v", err)
		}

		want := "# Gates\n\n(no gates)\n"
		if markdownOutput.String() != want {
			t.Fatalf("markdown output = %q, want %q", markdownOutput.String(), want)
		}
	})
}

func emptyIssueMarkdownFixture() IssueMarkdownView {
	ownerRole := "developer"
	createdAt := "2026-06-07T14:00:00Z"
	updatedAt := "2026-06-07T14:00:00Z"

	return IssueMarkdownView{
		Issue: issues.Issue{
			ID:                 "FID-empty01",
			Kind:               "task",
			Title:              "Render empty issue",
			Status:             "todo",
			StatusCategory:     "unstarted",
			OwnerRole:          &ownerRole,
			Labels:             []string{},
			AcceptanceCriteria: []issues.AcceptanceCriterion{},
			CreatedAt:          &createdAt,
			UpdatedAt:          &updatedAt,
		},
	}
}

func populatedIssueMarkdownFixture() IssueMarkdownView {
	parentTitle := "Core persistence layer"
	parentID := "FID-e1f2g3"
	priority := 1
	sequence := 3
	specMDPath := "project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md"
	ownerRole := "developer"
	createdAt := "2026-05-25T09:00:00Z"
	updatedAt := "2026-05-27T10:20:00Z"

	return IssueMarkdownView{
		Issue: issues.Issue{
			ID:             "FID-a1b2c3",
			Kind:           "task",
			Title:          "Implement journal append with flush",
			Status:         "in_progress",
			StatusCategory: "started",
			ParentID:       &parentID,
			ParentTitle:    &parentTitle,
			Priority:       &priority,
			Sequence:       &sequence,
			SpecMDPath:     &specMDPath,
			OwnerRole:      &ownerRole,
			Labels:         []string{"core", "journal"},
			AcceptanceCriteria: []issues.AcceptanceCriterion{
				{
					ID:      "AC-d4e5f6",
					Ordinal: 1,
					Text:    "Journal file is flushed to disk before returning success",
					LatestVerifications: map[string]any{
						"developer": map[string]any{
							"status":      "pass",
							"verified_at": "2026-05-27T10:15:00Z",
						},
					},
				},
				{
					ID:      "AC-e5f6g7",
					Ordinal: 2,
					Text:    "Partial write leaves journal intact",
					LatestVerifications: map[string]any{
						"developer": nil,
					},
				},
			},
			Links: []issues.IssueLink{
				{
					Relation:    "blocks",
					TargetID:    "FID-7c3e9f",
					TargetTitle: stringPointer("Schema migration runner"),
				},
				{
					Relation: "implements",
					TargetID: "PRD-ISSUE-001",
				},
			},
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		},
	}
}

func stringPointer(value string) *string {
	return &value
}
