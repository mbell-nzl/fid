package main

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestWorkflowListRendersDocumentedJSONFieldSet(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(t, workspace, "workflow", "list")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	var items []map[string]any
	if err := json.Unmarshal([]byte(stdout), &items); err != nil {
		t.Fatalf("decode workflow list JSON: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("workflow list returned no seeded workflows")
	}

	wantFields := []string{"active", "applies_to", "assigned_to", "id", "name", "states", "transitions"}
	for _, item := range items {
		gotFields := make([]string, 0, len(item))
		for field := range item {
			gotFields = append(gotFields, field)
		}
		sort.Strings(gotFields)
		if !reflect.DeepEqual(gotFields, wantFields) {
			t.Fatalf("workflow %q fields = %v, want %v", item["id"], gotFields, wantFields)
		}
	}
}

func TestWorkflowListFormatsFiltersAndAssignments(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	conn := openWorkflowCommandDB(t, workspace)
	seedWorkflowListCommandFixture(t, conn)
	closeCommandDB(t, conn)

	jsonOutput, jsonStderr := runFIDCommand(t, workspace, "workflow", "list", "--kind", "qa_review", "--active")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	want := decodeCommandJSONStructure(t, jsonOutput)
	if got, wantIDs := workflowListCommandIDs(t, want), []string{"alpha-any", "bravo-qa"}; !reflect.DeepEqual(got, wantIDs) {
		t.Fatalf("active kind IDs = %v, want %v", got, wantIDs)
	}
	assertWorkflowListCommandShapes(t, want)

	for _, formatName := range []string{"toon", "markdown", "yaml"} {
		stdout, stderr := runFIDCommand(t, workspace, "workflow", "list", "--kind", "qa_review", "--active", "--format", formatName)
		if stderr != "" {
			t.Fatalf("%s stderr = %q, want empty", formatName, stderr)
		}
		var got any
		switch formatName {
		case "toon":
			got = decodeCommandTOONStructure(t, stdout)
		case "markdown":
			got = decodeWorkflowListMarkdown(t, stdout)
		case "yaml":
			got = decodeCommandYAMLStructure(t, stdout)
		}
		if !commandStructuresEqual(want, got) {
			t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
		}
	}

	secondJSON, _ := runFIDCommand(t, workspace, "workflow", "list", "--kind", "qa_review", "--active")
	if secondJSON != jsonOutput {
		t.Fatalf("consecutive JSON outputs differ: %q != %q", secondJSON, jsonOutput)
	}
	unfiltered, _ := runFIDCommand(t, workspace, "workflow", "list", "--kind", "qa_review")
	if got, wantIDs := workflowListCommandIDs(t, decodeCommandJSONStructure(t, unfiltered)), []string{"alpha-any", "bravo-qa", "charlie-qa-inactive"}; !reflect.DeepEqual(got, wantIDs) {
		t.Fatalf("unfiltered kind IDs = %v, want %v", got, wantIDs)
	}

	emptyWorkspace := t.TempDir()
	clonePristineWorkspaceInto(t, emptyWorkspace)
	for _, formatName := range []string{"json", "toon", "markdown", "yaml"} {
		stdout, stderr := runFIDCommand(t, emptyWorkspace, "workflow", "list", "--kind", "qa_review", "--format", formatName)
		if stderr != "" || stdout != "[]\n" {
			t.Fatalf("empty %s = (%q, %q), want ([]\\n, empty stderr)", formatName, stdout, stderr)
		}
	}
}

func seedWorkflowListCommandFixture(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, item := range []struct {
		id, appliesTo string
		active        bool
	}{
		{"alpha-any", "", true},
		{"bravo-qa", "qa_review", true},
		{"charlie-qa-inactive", "qa_review", false},
		{"delta-epic", "epic", true},
	} {
		var appliesTo any
		if item.appliesTo != "" {
			appliesTo = item.appliesTo
		}
		if _, err := conn.Exec(`INSERT INTO workflows (id, name, applies_to, active) VALUES (?, ?, ?, ?)`, item.id, item.id, appliesTo, item.active); err != nil {
			t.Fatalf("insert workflow %q: %v", item.id, err)
		}
	}
	for ordinal, state := range []string{"open", "closed"} {
		if _, err := conn.Exec(`INSERT INTO workflow_states (workflow_id, state, category, ordinal) VALUES ('alpha-any', ?, 'unstarted', ?)`, state, ordinal+1); err != nil {
			t.Fatalf("insert alpha-any state %q: %v", state, err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO workflow_transitions (workflow_id, from_state, to_state) VALUES ('alpha-any', 'open', 'closed')`); err != nil {
		t.Fatalf("insert alpha-any transition: %v", err)
	}
	for _, assignment := range []struct{ scope, kind string }{{"kind", "task"}, {"project", ""}} {
		var kind any
		if assignment.kind != "" {
			kind = assignment.kind
		}
		if _, err := conn.Exec(`INSERT INTO workflow_assignments (scope, kind, workflow_id) VALUES (?, ?, 'alpha-any')`, assignment.scope, kind); err != nil {
			t.Fatalf("insert alpha-any %s assignment: %v", assignment.scope, err)
		}
	}
}

func workflowListCommandIDs(t *testing.T, value any) []string {
	t.Helper()
	rows, ok := value.([]any)
	if !ok {
		t.Fatalf("workflow list = %#v, want array", value)
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		item, ok := row.(map[string]any)
		if !ok {
			t.Fatalf("workflow item = %#v, want object", row)
		}
		ids = append(ids, item["id"].(string))
	}
	return ids
}

func assertWorkflowListCommandShapes(t *testing.T, value any) {
	t.Helper()
	rows := value.([]any)
	alpha := rows[0].(map[string]any)
	if alpha["applies_to"] != nil || alpha["states"].(json.Number).String() != "2" || alpha["transitions"].(json.Number).String() != "1" {
		t.Fatalf("alpha-any fields = %#v, want null applicability and counts", alpha)
	}
	assignments := alpha["assigned_to"].([]any)
	if got, want := assignments[0].(map[string]any), map[string]any{"scope": "kind", "kind": "task"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("kind assignment = %#v, want %#v", got, want)
	}
	project := assignments[1].(map[string]any)
	if project["scope"] != "project" || project["kind"] != nil {
		t.Fatalf("project assignment = %#v, want only project scope", project)
	}
	if _, exists := project["kind"]; exists {
		t.Fatalf("project assignment kind key = %#v, want absent", project["kind"])
	}
	bravo := rows[1].(map[string]any)
	if assigned := bravo["assigned_to"]; assigned == nil || len(assigned.([]any)) != 0 {
		t.Fatalf("unassigned workflow assigned_to = %#v, want empty array", assigned)
	}
}

func decodeWorkflowListMarkdown(t *testing.T, value string) any {
	t.Helper()
	if value == "[]\n" {
		return []any{}
	}
	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) < 4 {
		t.Fatalf("markdown = %q, want workflow table", value)
	}
	rows := make([]any, 0, len(lines)-4)
	for _, line := range lines[4:] {
		cells := splitMarkdownTableRow(line)
		if len(cells) != 7 {
			t.Fatalf("markdown row = %q, want 7 cells", line)
		}
		states, err := strconv.Atoi(cells[4])
		if err != nil {
			t.Fatalf("parse markdown states %q: %v", cells[4], err)
		}
		transitions, err := strconv.Atoi(cells[5])
		if err != nil {
			t.Fatalf("parse markdown transitions %q: %v", cells[5], err)
		}
		active, err := strconv.ParseBool(cells[3])
		if err != nil {
			t.Fatalf("parse markdown active %q: %v", cells[3], err)
		}
		var assignments any
		if err := json.Unmarshal([]byte(cells[6]), &assignments); err != nil {
			t.Fatalf("parse markdown assigned_to %q: %v", cells[6], err)
		}
		var appliesTo any
		if cells[2] != "null" {
			appliesTo = cells[2]
		}
		rows = append(rows, map[string]any{
			"id": cells[0], "name": cells[1], "applies_to": appliesTo, "active": active,
			"states": states, "transitions": transitions, "assigned_to": assignments,
		})
	}
	return rows
}
