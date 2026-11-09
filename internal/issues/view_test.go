package issues

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func readView(t *testing.T, workspace, id string) IssueView {
	t.Helper()
	view, err := ReadIssueView(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: id})
	if err != nil {
		t.Fatalf("ReadIssueView(%s) error = %v", id, err)
	}
	return view
}

// describeView writes a view as one line per field, references as
// "ID kind status category".
func describeView(view IssueView) []string {
	reference := func(ref IssueReference) string {
		return fmt.Sprintf("%s %s %s %s", ref.ID, ref.Kind, ref.Status, ref.StatusCategory)
	}
	references := func(refs []IssueReference) string {
		parts := []string{}
		for _, ref := range refs {
			parts = append(parts, reference(ref))
		}
		return strings.Join(parts, ", ")
	}
	lines := []string{"issue: " + view.Issue.ID + " " + view.Issue.Kind}
	optional := func(name string, set bool, value func() string) {
		if !set {
			lines = append(lines, name+": null")
			return
		}
		lines = append(lines, name+": "+value())
	}
	optional("state", view.State != nil, func() string { return *view.State })
	optional("epic", view.Epic != nil, func() string {
		return fmt.Sprintf("%s %d/%d", reference(view.Epic.IssueReference), view.Epic.Progress.Done, view.Epic.Progress.Total)
	})
	optional("task", view.Task != nil, func() string { return reference(*view.Task) })
	lines = append(lines, "blocked by: "+references(view.BlockedBy), "blocks: "+references(view.Blocks))
	optional("progress", view.Progress != nil, func() string { return fmt.Sprintf("%d/%d", view.Progress.Done, view.Progress.Total) })
	children := []string{}
	for _, child := range view.Children {
		part := child.ID + " " + child.Status
		if len(child.BlockedBy) > 0 {
			part += " blocked by " + strings.Join(child.BlockedBy, "+")
		}
		children = append(children, part)
	}
	return append(lines, "children: "+strings.Join(children, ", "))
}

// viewWork is an epic with tasks in every state, a task with subtasks, and
// blocks links recorded from either end, one of them to a done task. T1
// blocks T5 and T6, which fid status's order lists the other way round
// from their IDs.
func viewWork(t *testing.T) string {
	t.Helper()
	n := intPointer
	return seedOverview(t, []overviewRow{
		{"E1", "epic", "", "in_progress", n(1), nil, 1},
		{"T3", "task", "E1", "done", n(3), nil, 2},
		{"T1", "task", "E1", "in_progress", n(1), nil, 3},
		{"T2", "task", "E1", "todo", n(2), nil, 4},
		{"T6", "task", "E1", "todo", n(5), nil, 5},
		{"T4", "task", "E1", "cancelled", n(4), nil, 6},
		{"S2", "subtask", "T1", "in_progress", nil, n(1), 7},
		{"S1", "subtask", "T1", "done", nil, n(0), 8},
		{"S3", "subtask", "T1", "todo", nil, n(2), 9},
		{"S4", "subtask", "T1", "cancelled", nil, nil, 10},
		{"T5", "task", "", "todo", nil, nil, 11},
		{"E2", "epic", "", "todo", n(2), nil, 12},
		{"T7", "task", "E2", "in_progress", nil, nil, 13},
		{"E3", "epic", "", "todo", n(3), nil, 14},
		{"T8", "task", "E3", "todo", nil, nil, 15},
		{"E4", "epic", "", "closed", n(4), nil, 16},
		{"E5", "epic", "", "cancelled", n(5), nil, 17},
		{"E6", "epic", "", "in_progress", n(6), nil, 18},
	}, [][3]string{
		{"T1", "blocked_by", "T3"},
		{"T1", "blocks", "T2"}, {"T2", "blocked_by", "T1"},
		{"T5", "blocked_by", "T1"},
		{"T1", "blocks", "T6"},
		{"T6", "blocks", "S3"},
	})
}

func TestReadIssueViewReadsAnIssueWithWhatSurroundsIt(t *testing.T) {
	t.Parallel()
	workspace := viewWork(t)
	for _, test := range []struct {
		id   string
		want []string
	}{
		{"T1", []string{"issue: T1 task", "state: null", "epic: E1 epic in_progress started 1/4", "task: null",
			"blocked by: T3 task done completed", "blocks: T2 task todo unstarted, T6 task todo unstarted, T5 task todo unstarted", "progress: 1/3",
			"children: S1 done, S2 in_progress, S3 todo blocked by T6, S4 cancelled"}},
		{"S3", []string{"issue: S3 subtask", "state: null", "epic: E1 epic in_progress started 1/4", "task: T1 task in_progress started",
			"blocked by: T6 task todo unstarted", "blocks: ", "progress: null", "children: "}},
		{"T2", []string{"issue: T2 task", "state: null", "epic: E1 epic in_progress started 1/4", "task: null",
			"blocked by: T1 task in_progress started", "blocks: ", "progress: null", "children: "}},
		{"T5", []string{"issue: T5 task", "state: null", "epic: null", "task: null",
			"blocked by: T1 task in_progress started", "blocks: ", "progress: null", "children: "}},
		{"E1", []string{"issue: E1 epic", "state: in-progress", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 1/4",
			"children: T1 in_progress, T2 todo blocked by T1, T3 done, T4 cancelled, T6 todo blocked by T1"}},
		{"E2", []string{"issue: E2 epic", "state: in-progress", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 0/1",
			"children: T7 in_progress"}},
		{"E3", []string{"issue: E3 epic", "state: to-do", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 0/1",
			"children: T8 todo"}},
		{"E4", []string{"issue: E4 epic", "state: completed", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 0/0", "children: "}},
		{"E5", []string{"issue: E5 epic", "state: cancelled", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 0/0", "children: "}},
		{"E6", []string{"issue: E6 epic", "state: in-progress", "epic: null", "task: null", "blocked by: ", "blocks: ", "progress: 0/0", "children: "}},
	} {
		t.Run(test.id, func(t *testing.T) {
			t.Parallel()
			if got := describeView(readView(t, workspace, test.id)); !slices.Equal(got, test.want) {
				t.Fatalf("view =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(test.want, "\n"))
			}
		})
	}
}

// The view's issue is issue.get's, and its lists are empty, never nil.
func TestReadIssueViewHoldsIssueGetsIssue(t *testing.T) {
	t.Parallel()
	workspace := viewWork(t)
	options := GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: "S3"}
	detail, err := ReadIssueDetail(options)
	if err != nil {
		t.Fatal(err)
	}
	view := readView(t, workspace, "S3")
	if !reflect.DeepEqual(view.Issue, detail) {
		t.Fatalf("view's issue = %#v, want issue.get's %#v", view.Issue, detail)
	}
	if view.Blocks == nil || view.Children == nil || view.BlockedBy == nil {
		t.Fatalf("view lists = %#v %#v %#v, want empty lists", view.BlockedBy, view.Blocks, view.Children)
	}
}

// Only an issue's exact ID finds it: not another case, a prefix or another
// record's ID.
func TestReadIssueViewFindsOnlyAnIssuesExactID(t *testing.T) {
	t.Parallel()
	workspace := viewWork(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`INSERT INTO acceptance_criteria (id, issue_id, ordinal, text) VALUES ('AC-1', 'T1', 1, 'A criterion')`); err != nil {
		t.Fatal(err)
	}
	closeIssueTestDB(t, conn)
	for _, id := range []string{"T9", "t1", "T", "AC-1"} {
		_, err := ReadIssueView(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: id})
		var notFound *NotFoundError
		if !errors.As(err, &notFound) || notFound.IssueID != id {
			t.Errorf("ReadIssueView(%q) error = %v, want not found", id, err)
		}
	}
}
