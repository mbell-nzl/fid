package issuestate

import "testing"

func TestInitialState(t *testing.T) {
	t.Parallel()
	if got := InitialState([]State{{Name: "started", Category: "started"}, {Name: "todo", Category: "unstarted"}, {Name: "queued", Category: "unstarted"}}); got != "todo" {
		t.Fatalf("InitialState(first unstarted) = %q, want todo", got)
	}
	if got := InitialState([]State{{Name: "started", Category: "started"}, {Name: "done", Category: "completed"}}); got != "started" {
		t.Fatalf("InitialState(no unstarted) = %q, want started", got)
	}
}
