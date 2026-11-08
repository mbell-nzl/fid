package issues

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// overviewRow is an issue to seed: an epic, a task or a subtask, by its
// parent, workflow state, sequence, priority and minute of creation.
type overviewRow struct {
	id, kind, parent, status string
	sequence, priority       *int
	minute                   int
}

func seedOverview(t *testing.T, rows []overviewRow, links [][3]string) string {
	t.Helper()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	seeded := []readyRow{}
	for _, row := range rows {
		var parent *string
		if row.parent != "" {
			parent = stringPointer(row.parent)
		}
		created := fmt.Sprintf("2026-11-01T10:%02d:00Z", row.minute)
		var scheduling *string
		if row.kind == "epic" {
			scheduling = stringPointer("prioritised")
			for _, task := range rows {
				if task.parent == row.id && task.sequence != nil {
					scheduling = stringPointer("sequential")
				}
			}
		}
		seeded = append(seeded, readyRow{id: row.id, kind: row.kind, parentID: parent, title: "Title of " + row.id, status: row.status,
			workflowID: row.kind + "-default", sequence: row.sequence, priority: row.priority, createdAt: created, updatedAt: created, scheduling: scheduling})
	}
	seedReadyRows(t, conn, seeded)
	for _, link := range links {
		if _, err := conn.Exec(`INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at) VALUES ('issue', ?, ?, 'issue', ?, '2026-11-01T11:00:00Z')`,
			link[0], link[1], link[2]); err != nil {
			t.Fatalf("seed link %v: %v", link, err)
		}
	}
	return workspace
}

func readOverview(t *testing.T, workspace string) Overview {
	t.Helper()
	overview, err := ReadOverview(OverviewOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t)})
	if err != nil {
		t.Fatalf("ReadOverview() error = %v", err)
	}
	return overview
}

// describeSection writes a section as one line per epic and task:
// "E1 in-progress 1/4: T1 in_progress, T2 qa_review" and "T9 todo".
func describeSection(section OverviewSection) []string {
	tasks := func(list []OverviewTask) string {
		parts := []string{}
		for _, task := range list {
			part := task.ID + " " + task.Status
			if len(task.BlockedBy) > 0 {
				part += " blocked by " + strings.Join(task.BlockedBy, "+")
			}
			parts = append(parts, part)
		}
		return strings.Join(parts, ", ")
	}
	lines := []string{}
	for _, epic := range section.Epics {
		lines = append(lines, fmt.Sprintf("%s %s %d/%d: %s", epic.ID, epic.State, epic.Progress.Done, epic.Progress.Total, tasks(epic.Tasks)))
	}
	for _, task := range section.Tasks {
		lines = append(lines, tasks([]OverviewTask{task}))
	}
	return lines
}

func assertSection(t *testing.T, name string, section OverviewSection, want []string) {
	t.Helper()
	if got := describeSection(section); !slices.Equal(got, want) {
		t.Fatalf("%s =\n%s\nwant\n%s", name, strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestReadOverviewSortsWorkIntoInProgressAndToDo(t *testing.T) {
	t.Parallel()
	n := intPointer
	rows := []overviewRow{
		{"E1", "epic", "", "in_progress", n(1), nil, 1},
		{"T1", "task", "E1", "in_progress", n(1), nil, 2},
		{"T2", "task", "E1", "qa_review", n(2), nil, 3},
		{"T3", "task", "E1", "todo", n(3), nil, 4},
		{"T4", "task", "E1", "done", n(4), nil, 5},
		{"T5", "task", "E1", "cancelled", n(5), nil, 6},
		{"X1", "subtask", "T1", "in_progress", nil, nil, 7},
		{"E2", "epic", "", "todo", n(2), nil, 8},
		{"T6", "task", "E2", "in_progress", n(1), nil, 9},
		{"E3", "epic", "", "todo", n(3), nil, 10},
		{"T7", "task", "E3", "qa_pass", n(1), nil, 11},
		{"T8", "task", "E3", "todo", n(2), nil, 12},
		{"E4", "epic", "", "todo", n(4), nil, 13},
		{"T9", "task", "E4", "todo", n(1), nil, 14},
		{"T10", "task", "E4", "todo", n(2), nil, 15},
		{"E5", "epic", "", "todo", n(5), nil, 16},
		{"E6", "epic", "", "closed", n(6), nil, 17},
		{"T11", "task", "E6", "in_progress", n(1), nil, 18},
		{"T12", "task", "E6", "todo", n(2), nil, 19},
		{"E7", "epic", "", "cancelled", n(7), nil, 20},
		{"T13", "task", "E7", "todo", n(1), nil, 21},
		{"E8", "epic", "", "wrapped", n(8), nil, 22},
		{"S1", "task", "", "failed_gate", nil, nil, 23},
		{"S2", "task", "", "todo", nil, nil, 24},
		{"S3", "task", "", "done", nil, nil, 25},
		{"S4", "task", "", "cancelled", nil, nil, 26},
	}
	overview := readOverview(t, seedOverview(t, rows, nil))
	if overview.Empty {
		t.Fatal("overview is empty, want work")
	}
	assertSection(t, "in progress", overview.InProgress, []string{
		"E1 in-progress 1/4: T1 in_progress, T2 qa_review",
		"E2 in-progress 0/1: T6 in_progress",
		"E3 in-progress 1/2: ",
		"E8 in-progress 0/0: ",
		"T11 in_progress",
		"S1 failed_gate",
	})
	assertSection(t, "to do", overview.ToDo, []string{
		"E1 in-progress 1/4: T3 todo",
		"E3 in-progress 1/2: T8 todo",
		"E4 to-do 0/2: T9 todo, T10 todo",
		"E5 to-do 0/0: ",
		"T13 todo",
		"T12 todo",
		"S2 todo",
	})
	epic := overview.InProgress.Epics[0]
	if epic.Title != "Title of E1" || epic.StatusCategory != "started" || epic.Sequence == nil || *epic.Sequence != 1 || epic.Priority != nil {
		t.Fatalf("epic = %+v", epic)
	}
	task := overview.InProgress.Epics[0].Tasks[1]
	if task.Title != "Title of T2" || task.StatusCategory != "started" || task.EpicID == nil || *task.EpicID != "E1" || task.BlockedBy == nil {
		t.Fatalf("task = %+v", task)
	}
	if flat := overview.InProgress.Tasks[0]; flat.EpicID == nil || *flat.EpicID != "E6" {
		t.Fatalf("task under a closed epic = %+v, want its epic kept", flat)
	}
}

// A sequential epic's tasks have a sequence, and other tasks have none, so
// tasks are ordered by sequence in one and by priority in the other.
func TestReadOverviewOrdersBySequencePriorityCreationAndID(t *testing.T) {
	t.Parallel()
	n := intPointer
	rows := []overviewRow{
		{"EA", "epic", "", "todo", nil, n(1), 1},
		{"EB", "epic", "", "todo", n(2), nil, 1},
		{"EC", "epic", "", "todo", n(2), n(0), 5},
		{"ED", "epic", "", "todo", nil, nil, 0},
		{"EF", "epic", "", "todo", nil, nil, 3},
		{"EE", "epic", "", "todo", nil, nil, 3},
		{"TB", "task", "EC", "todo", n(2), nil, 1},
		{"TC", "task", "EC", "todo", n(1), nil, 9},
		{"TA", "task", "EB", "todo", nil, n(1), 1},
		{"TD", "task", "EB", "todo", nil, n(1), 0},
		{"TE", "task", "EB", "todo", nil, nil, 0},
		{"TF", "task", "EB", "todo", nil, n(0), 9},
		{"SB", "task", "", "todo", nil, n(3), 2},
		{"SA", "task", "", "todo", nil, n(2), 2},
		{"SC", "task", "", "todo", nil, nil, 0},
	}
	overview := readOverview(t, seedOverview(t, rows, nil))
	assertSection(t, "to do", overview.ToDo, []string{
		"EC to-do 0/2: TC todo, TB todo",
		"EB to-do 0/4: TF todo, TD todo, TA todo, TE todo",
		"EA to-do 0/0: ",
		"ED to-do 0/0: ",
		"EE to-do 0/0: ",
		"EF to-do 0/0: ",
		"SA todo",
		"SB todo",
		"SC todo",
	})
}

func TestReadOverviewNamesOpenBlockers(t *testing.T) {
	t.Parallel()
	rows := []overviewRow{
		{"B1", "task", "", "todo", nil, nil, 1},
		{"B2", "task", "", "in_progress", nil, nil, 2},
		{"X", "task", "", "todo", nil, nil, 3},
		{"Y", "epic", "", "in_progress", nil, nil, 4},
		{"Z", "task", "", "done", nil, nil, 5},
		{"W", "task", "", "cancelled", nil, nil, 6},
		{"V", "subtask", "B2", "in_progress", nil, nil, 7},
	}
	links := [][3]string{
		{"X", "blocks", "B1"}, {"B1", "blocked_by", "X"},
		{"B1", "blocked_by", "Y"},
		{"Z", "blocks", "B1"}, {"W", "blocks", "B1"},
		{"V", "blocks", "B2"},
		{"B1", "relates_to", "B2"}, {"X", "duplicates", "B2"},
	}
	overview := readOverview(t, seedOverview(t, rows, links))
	assertSection(t, "in progress", overview.InProgress, []string{"Y in-progress 0/0: ", "B2 in_progress blocked by V"})
	assertSection(t, "to do", overview.ToDo, []string{"B1 todo blocked by X+Y", "X todo"})
}

func TestReadOverviewIsEmptyOnlyWithNoEpicsOrTasks(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		rows  []overviewRow
		empty bool
	}{
		{"a new project", nil, true},
		{"a closed epic", []overviewRow{{"E1", "epic", "", "closed", nil, nil, 1}}, false},
		{"a done task", []overviewRow{{"T1", "task", "", "done", nil, nil, 1}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			overview := readOverview(t, seedOverview(t, test.rows, nil))
			if overview.Empty != test.empty {
				t.Fatalf("Empty = %v, want %v", overview.Empty, test.empty)
			}
			if test.empty && (overview.InProgress.Epics == nil || overview.InProgress.Tasks == nil || overview.ToDo.Epics == nil || overview.ToDo.Tasks == nil) {
				t.Fatalf("overview = %+v, want empty lists rather than nil", overview)
			}
		})
	}
}
