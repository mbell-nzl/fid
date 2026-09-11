package workflow

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestDefineAppendsOneAtomicWorkflowDefinedEvent(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	journalPath := workflowJournalPath(t, workspace)
	beforeJournal := readWorkflowTestFile(t, journalPath)

	_, err := Define(DefineOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: migrationsDir,
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
		Request: DefinitionRequest{
			Name:        "requirement-flow",
			AppliesTo:   "task",
			States:      []string{"todo", "done", "cancelled"},
			Categories:  []string{"unstarted", "completed", "cancelled"},
			Transitions: []Transition{{From: "todo", To: "done"}},
			Requirements: []Requirement{{
				From: "todo", To: "done", Rule: RuleHasCompletionReport,
			}},
		},
	})
	if err == nil || err.Error() != "workflow transition requirements are unavailable" {
		t.Fatalf("Define() error = %v, want unavailable transition requirements", err)
	}
	if afterRejectedJournal := readWorkflowTestFile(t, journalPath); !bytes.Equal(afterRejectedJournal, beforeJournal) {
		t.Fatal("rejected transition requirements changed journal bytes")
	}

	result, err := Define(DefineOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: migrationsDir,
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
		Request: DefinitionRequest{
			Name:        "direct-flow",
			AppliesTo:   "task",
			States:      []string{"todo", "done", "cancelled"},
			Categories:  []string{"unstarted", "completed", "cancelled"},
			Transitions: []Transition{{From: "todo", To: "done"}},
		},
	})
	if err != nil {
		t.Fatalf("Define() error = %v", err)
	}
	if result.Name != "direct-flow" || result.AppliesTo != "task" || result.States != 3 || result.Transitions != 3 || result.Requirements != 0 || result.EventID == "" || result.OccurredAt.IsZero() {
		t.Fatalf("Define() result = %#v, want projected counts and appended Event identity", result)
	}

	afterJournal := readWorkflowTestFile(t, journalPath)
	beforeLines := strings.Split(strings.TrimSpace(string(beforeJournal)), "\n")
	afterLines := strings.Split(strings.TrimSpace(string(afterJournal)), "\n")
	if got, want := len(afterLines), len(beforeLines)+1; got != want {
		t.Fatalf("journal lines = %d, want %d", got, want)
	}
	var event journal.Event
	if err := json.Unmarshal([]byte(afterLines[len(afterLines)-1]), &event); err != nil {
		t.Fatalf("decode appended event: %v", err)
	}
	if event.EventType != "workflow.defined" || event.EventSchema != journal.EventSchemaV2 {
		t.Fatalf("appended event = (%q, schema %d), want Event-v2 workflow.defined", event.EventType, event.EventSchema)
	}
	if got, want := len(event.Effects), 8; got != want || event.Effects[len(event.Effects)-1].Table != "idempotency_results" {
		t.Fatalf("event effects = %#v, want workflow effects plus idempotency result", event.Effects)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeRegistryDB(t, conn)
	registry, err := Load(conn)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if _, ok := registry.Workflow("direct-flow"); !ok {
		t.Fatal("Load() registry missing direct-flow")
	}
	var activityCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityCount); err != nil {
		t.Fatalf("count activity_log: %v", err)
	}
	if activityCount != 0 {
		t.Fatalf("activity_log rows = %d, want 0", activityCount)
	}
	if _, err := Define(DefineOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Actor: testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(), Request: DefinitionRequest{Name: "direct-flow", AppliesTo: "task", States: []string{"todo", "cancelled"}, Categories: []string{"unstarted", "cancelled"}, Transitions: []Transition{{From: "todo", To: "cancelled"}}}}); err == nil {
		t.Fatal("duplicate Define() error = nil, want rejection")
	}
	if finalJournal := readWorkflowTestFile(t, journalPath); !bytes.Equal(finalJournal, afterJournal) {
		t.Fatal("rejected duplicate Define() changed journal bytes")
	}
}

func initializedWorkflowWorkspace(t *testing.T) (string, string) {
	t.Helper()
	workspace := t.TempDir()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	migrationsDir := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", "..", "migrations"))
	if err := cloneWorkflowPristine(workspace); err != nil {
		t.Fatalf("clone initialized workflow workspace: %v", err)
	}
	return workspace, migrationsDir
}

func workflowJournalPath(t *testing.T, workspace string) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(workspace, ".fid", "journal", "*.jsonl"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("journal paths = %#v, error = %v, want exactly one", paths, err)
	}
	return paths[0]
}

func readWorkflowTestFile(t *testing.T, path string) []byte {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	return payload
}
