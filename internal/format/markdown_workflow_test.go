package format

import (
	"bytes"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestWriteWorkflowListMarkdownMatchesCommittedLiteral(t *testing.T) {
	appliesTo := "task"
	items := []workflow.ListItem{
		{
			ID: "any-flow", Name: "Any flow", Active: true, States: 2, Transitions: 1,
			AssignedTo: []workflow.Assignment{{Scope: "kind", Kind: "task"}, {Scope: "project"}},
		},
		{
			ID: "task-flow", Name: "Task flow", AppliesTo: &appliesTo, Active: false, States: 3, Transitions: 4,
			AssignedTo: []workflow.Assignment{},
		},
	}

	var output bytes.Buffer
	if err := Write(&output, Markdown, items); err != nil {
		t.Fatalf("Write(Markdown) error = %v", err)
	}
	const want = "# Workflows\n\n" +
		"| id | name | applies_to | active | states | transitions | assigned_to |\n" +
		"| --- | --- | --- | ---: | ---: | ---: | --- |\n" +
		"| any-flow | Any flow | null | true | 2 | 1 | [{\"scope\":\"kind\",\"kind\":\"task\"},{\"scope\":\"project\"}] |\n" +
		"| task-flow | Task flow | task | false | 3 | 4 | [] |\n"
	if got := output.String(); got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
}

func TestWriteWorkflowListMarkdownRendersEmptyList(t *testing.T) {
	var output bytes.Buffer
	if err := Write(&output, Markdown, []workflow.ListItem{}); err != nil {
		t.Fatalf("Write(Markdown) error = %v", err)
	}
	if got, want := output.String(), "[]\n"; got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
}
