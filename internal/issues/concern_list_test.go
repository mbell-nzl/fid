package issues

import "testing"

func TestListConcernsFiltersByIssueAndStatus(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	items, err := ListConcerns(ConcernListOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t)})
	if err != nil {
		t.Fatalf("ListConcerns() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("ListConcerns() = %#v, want empty frozen inventory", items)
	}
}

func TestListConcernsOrdersSameSecondRowsByAppendOrder(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	items, err := ListConcerns(ConcernListOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Status: "open"})
	if err != nil {
		t.Fatalf("ListConcerns(open) error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("ListConcerns(open) = %#v, want empty frozen inventory", items)
	}
}
