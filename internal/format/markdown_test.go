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
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		},
		ParentTitle: &parentTitle,
		Links: []IssueMarkdownLink{
			{
				Relation:    "blocks",
				TargetID:    "FID-7c3e9f",
				TargetTitle: "Schema migration runner",
			},
			{
				Relation: "implements",
				TargetID: "PRD-ISSUE-001",
			},
		},
	}
}
