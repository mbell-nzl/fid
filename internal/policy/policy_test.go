package policy

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestRequireKnownRoleRejectsMissingAndUnknownRolesAndAcceptsSeededRole(t *testing.T) {
	t.Parallel()

	workspace := newPolicyWorkspace(t)

	err := RequireKnownRole(workspace, "   ")
	if err == nil {
		t.Fatal("RequireKnownRole(missing) error = nil, want missing role error")
	}

	var missingErr *MissingRoleError
	if !errors.As(err, &missingErr) {
		t.Fatalf("RequireKnownRole(missing) error = %T, want *MissingRoleError", err)
	}
	if missingErr.Error() != "--by is required" {
		t.Fatalf("MissingRoleError.Error() = %q, want %q", missingErr.Error(), "--by is required")
	}

	err = RequireKnownRole(workspace, "intern")
	if err == nil {
		t.Fatal("RequireKnownRole(unknown) error = nil, want unknown role error")
	}

	var unknownErr *UnknownRoleError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("RequireKnownRole(unknown) error = %T, want *UnknownRoleError", err)
	}
	if unknownErr.Role != "intern" {
		t.Fatalf("UnknownRoleError.Role = %q, want %q", unknownErr.Role, "intern")
	}
	if unknownErr.Error() != `role "intern" not in project profile` {
		t.Fatalf("UnknownRoleError.Error() = %q, want %q", unknownErr.Error(), `role "intern" not in project profile`)
	}

	if err := RequireKnownRole(workspace, "developer"); err != nil {
		t.Fatalf("RequireKnownRole(known) error = %v", err)
	}
}

func TestKnownRolesValidatesWorkspaceAndJournalDirectory(t *testing.T) {
	t.Parallel()

	_, err := KnownRoles("  ")
	if err == nil {
		t.Fatal("KnownRoles(empty) error = nil, want workspace root error")
	}
	if err.Error() != "known roles: workspace root is required" {
		t.Fatalf("KnownRoles(empty) error = %q, want %q", err.Error(), "known roles: workspace root is required")
	}

	workspace := t.TempDir()
	_, err = KnownRoles(workspace)
	if err == nil {
		t.Fatal("KnownRoles(missing journal dir) error = nil, want read error")
	}
	if want := filepath.ToSlash(filepath.Join(workspace, ".fid", "journal")); !containsNormalizedPath(err.Error(), want) {
		t.Fatalf("KnownRoles(missing journal dir) error = %q, want journal path %q", err.Error(), want)
	}
}

func TestKnownRolesCollectsInsertUpdateDeleteAndActiveFalse(t *testing.T) {
	t.Parallel()

	workspace := newPolicyJournalWorkspace(t)
	writePolicyJournalFile(t, workspace, "2026-06-07.jsonl",
		policyRoleEvent("EVT-20260607-000001", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "developer"},
			Values: map[string]any{
				"description": "Developer",
				"active":      true,
			},
		}),
		policyRoleEvent("EVT-20260607-000002", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "qa_reviewer"},
			Values: map[string]any{
				"description": "QA Reviewer",
				"active":      true,
			},
		}),
		policyRoleEvent("EVT-20260607-000003", journal.Effect{
			Op:    "update",
			Table: "roles",
			PK:    map[string]any{"name": "qa_reviewer"},
			Values: map[string]any{
				"active": false,
			},
		}),
		policyRoleEvent("EVT-20260607-000004", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "qa_lead"},
			Values: map[string]any{
				"description": "QA Lead",
			},
		}),
		policyRoleEvent("EVT-20260607-000005", journal.Effect{
			Op:    "delete",
			Table: "roles",
			PK:    map[string]any{"name": "qa_lead"},
		}),
	)

	got, err := KnownRoles(workspace)
	if err != nil {
		t.Fatalf("KnownRoles() error = %v", err)
	}

	want := []string{"developer"}
	if !equalStrings(got, want) {
		t.Fatalf("KnownRoles() = %#v, want %#v", got, want)
	}
}

func TestKnownRolesTRG002CorruptTailTolerance(t *testing.T) {
	t.Parallel()

	t.Run("valid role before malformed trailing line reports corruption", func(t *testing.T) {
		t.Parallel()

		workspace := newPolicyJournalWorkspace(t)
		writePolicyJournalFile(t, workspace, "2026-06-07.jsonl",
			policyRoleEvent("EVT-20260607-000001", journal.Effect{
				Op:    "insert",
				Table: "roles",
				PK:    map[string]any{"name": "developer"},
				Values: map[string]any{
					"active": true,
				},
			}),
		)
		appendMalformedPolicyLine(t, workspace, "2026-06-07.jsonl", "{broken-json")

		_, err := KnownRoles(workspace)
		if err == nil {
			t.Fatal("KnownRoles() error = nil, want corrupt journal error")
		}

		assertPolicyCorruptJournalError(t, err, "2026-06-07.jsonl", 2)
	})

	t.Run("malformed first line surfaces an error", func(t *testing.T) {
		t.Parallel()

		workspace := newPolicyJournalWorkspace(t)
		writeRawPolicyJournalFile(t, workspace, "2026-06-07.jsonl", "{broken-json\n")

		_, err := KnownRoles(workspace)
		if err == nil {
			t.Fatal("KnownRoles() error = nil, want corrupt journal error")
		}

		assertPolicyCorruptJournalError(t, err, "2026-06-07.jsonl", 1)
	})

	t.Run("valid role before a future-schema trailing event reports corruption", func(t *testing.T) {
		t.Parallel()

		workspace := newPolicyJournalWorkspace(t)
		writePolicyJournalFile(t, workspace, "2026-06-07.jsonl",
			policyRoleEvent("EVT-20260607-000001", journal.Effect{
				Op:    "insert",
				Table: "roles",
				PK:    map[string]any{"name": "developer"},
				Values: map[string]any{
					"active": true,
				},
			}),
		)
		appendMalformedPolicyLine(t, workspace, "2026-06-07.jsonl", `{"event_id":"EVT-20260607-000002","event_schema":3}`)

		_, err := KnownRoles(workspace)
		if err == nil {
			t.Fatal("KnownRoles() error = nil, want corrupt journal error")
		}
		assertPolicyCorruptJournalError(t, err, "2026-06-07.jsonl", 2)
	})
}

func TestKnownRolesCollectsCleanRolesAcrossFilesAndLines(t *testing.T) {
	t.Parallel()

	workspace := newPolicyJournalWorkspace(t)
	writePolicyJournalFile(t, workspace, "2026-06-07.jsonl",
		policyRoleEvent("EVT-20260607-000001", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "developer"},
			Values: map[string]any{
				"active": true,
			},
		}),
		policyRoleEvent("EVT-20260607-000002", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "qa_reviewer"},
			Values: map[string]any{
				"active": true,
			},
		}),
	)
	writePolicyJournalFile(t, workspace, "2026-06-08.jsonl",
		policyRoleEvent("EVT-20260608-000001", journal.Effect{
			Op:    "update",
			Table: "roles",
			PK:    map[string]any{"name": "qa_reviewer"},
			Values: map[string]any{
				"active": false,
			},
		}),
		policyRoleEvent("EVT-20260608-000002", journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "qa_lead"},
			Values: map[string]any{
				"active": true,
			},
		}),
	)

	got, err := KnownRoles(workspace)
	if err != nil {
		t.Fatalf("KnownRoles() error = %v", err)
	}
	if !equalStrings(got, []string{"developer", "qa_lead"}) {
		t.Fatalf("KnownRoles() = %#v, want %#v", got, []string{"developer", "qa_lead"})
	}
}

func TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		want []string
	}{
		{
			name: "exact bootstrap recognizes the closed V1 vocabulary",
			want: []string{
				"developer",
				"human",
				"qa_lead",
				"qa_reviewer",
				"tl",
				"tl_lite",
				"workflow_runner",
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newPolicyJournalWorkspace(t)
			writePolicyJournalFile(t, workspace, "2026-08-23.jsonl", policyBootstrapEvent())
			journalPath := filepath.Join(workspace, ".fid", "journal", "2026-08-23.jsonl")
			before, err := os.ReadFile(journalPath)
			if err != nil {
				t.Fatalf("ReadFile(before lookup) error = %v", err)
			}

			got, err := KnownRoles(workspace)
			if err != nil {
				t.Fatalf("KnownRoles() error = %v", err)
			}
			if !equalStrings(got, tt.want) {
				t.Fatalf("KnownRoles() = %#v, want %#v", got, tt.want)
			}

			after, err := os.ReadFile(journalPath)
			if err != nil {
				t.Fatalf("ReadFile(after lookup) error = %v", err)
			}
			if string(after) != string(before) {
				t.Fatal("KnownRoles() modified the bootstrap journal")
			}
		})
	}
}

func TestKnownRolesKeepsExplicitRoleFactsAuthoritative(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		effect journal.Effect
		want   []string
	}{
		{
			name: "active explicit role does not receive fallback roles",
			effect: journal.Effect{
				Op:     "insert",
				Table:  "roles",
				PK:     map[string]any{"name": "qa_lead"},
				Values: map[string]any{"active": true},
			},
			want: []string{"qa_lead"},
		},
		{
			name: "inactive explicit role still blocks the fallback",
			effect: journal.Effect{
				Op:     "insert",
				Table:  "roles",
				PK:     map[string]any{"name": "developer"},
				Values: map[string]any{"active": false},
			},
			want: []string{},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newPolicyJournalWorkspace(t)
			writePolicyJournalFile(t, workspace, "2026-08-23.jsonl",
				policyBootstrapEvent(),
				policyRoleEvent("EVT-20260823-000002", tt.effect),
			)

			got, err := KnownRoles(workspace)
			if err != nil {
				t.Fatalf("KnownRoles() error = %v", err)
			}
			if !equalStrings(got, tt.want) {
				t.Fatalf("KnownRoles() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestKnownRolesRejectsCompatibilityForUninitializedJournal(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
	}{
		{name: "empty journal directory"},
		{name: "empty journal file"},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newPolicyJournalWorkspace(t)
			if tt.name == "empty journal file" {
				writeRawPolicyJournalFile(t, workspace, "2026-08-23.jsonl", "")
			}

			got, err := KnownRoles(workspace)
			if err != nil {
				t.Fatalf("KnownRoles() error = %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("KnownRoles() = %#v, want no roles", got)
			}
		})
	}
}

func TestKnownRolesRejectsCompatibilityForCorruptJournal(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		contents string
	}{
		{name: "malformed event after bootstrap", contents: "{broken-json\n"},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newPolicyJournalWorkspace(t)
			writePolicyJournalFile(t, workspace, "2026-08-23.jsonl", policyBootstrapEvent())
			appendMalformedPolicyLine(t, workspace, "2026-08-23.jsonl", strings.TrimSpace(tt.contents))

			_, err := KnownRoles(workspace)
			if err == nil {
				t.Fatal("KnownRoles() error = nil, want corrupt journal error")
			}
			assertPolicyCorruptJournalError(t, err, "2026-08-23.jsonl", 2)
		})
	}
}

func TestPolicyHelpersCoverActiveValueDefaults(t *testing.T) {
	t.Parallel()

	if !activeValue(nil) {
		t.Fatal("activeValue(nil) = false, want true")
	}
	if !activeValue(map[string]any{}) {
		t.Fatal("activeValue(empty) = false, want true")
	}
	if activeValue(map[string]any{"active": false}) {
		t.Fatal("activeValue(false) = true, want false")
	}
	if !activeValue(map[string]any{"active": "not-a-bool"}) {
		t.Fatal("activeValue(non-bool) = false, want true")
	}
}

func newPolicyWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testPolicyMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func newPolicyJournalWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	journalDir := filepath.Join(workspace, ".fid", "journal")
	if err := os.MkdirAll(journalDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", journalDir, err)
	}

	return workspace
}

func writePolicyJournalFile(t *testing.T, workspace, fileName string, events ...journal.Event) {
	t.Helper()

	lines := make([][]byte, 0, len(events))
	for _, event := range events {
		line, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		lines = append(lines, line)
	}

	writeRawPolicyJournalFile(t, workspace, fileName, joinLines(lines))
}

func writeRawPolicyJournalFile(t *testing.T, workspace, fileName, payload string) {
	t.Helper()

	path := filepath.Join(workspace, ".fid", "journal", fileName)
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func appendMalformedPolicyLine(t *testing.T, workspace, fileName, payload string) {
	t.Helper()

	path := filepath.Join(workspace, ".fid", "journal", fileName)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("OpenFile(%q) error = %v", path, err)
	}
	defer closePolicyFile(t, file)

	if _, err := file.WriteString(payload + "\n"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
}

func policyRoleEvent(eventID string, effects ...journal.Effect) journal.Event {
	return journal.Event{
		EventID:     eventID,
		EventSchema: journal.EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.June, 7, 9, 0, 0, 0, time.UTC),
		Actor: journal.Actor{
			Role: "human",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "init"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: effects,
	}
}

func policyBootstrapEvent() journal.Event {
	event := bootstrap.BuildProjectInitializedEvent("project-policy")
	event.EventID = "EVT-20260823-000001"
	event.OccurredAt = time.Date(2026, time.August, 23, 0, 0, 0, 0, time.UTC)
	return event
}

func joinLines(lines [][]byte) string {
	if len(lines) == 0 {
		return ""
	}

	payload := make([]byte, 0, len(lines)*64)
	for _, line := range lines {
		payload = append(payload, line...)
		payload = append(payload, '\n')
	}

	return string(payload)
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}

	sortedGot := append([]string(nil), got...)
	sortedWant := append([]string(nil), want...)
	sort.Strings(sortedGot)
	sort.Strings(sortedWant)
	for index := range sortedGot {
		if sortedGot[index] != sortedWant[index] {
			return false
		}
	}

	return true
}

func containsNormalizedPath(message, wantPath string) bool {
	return strings.Contains(filepath.ToSlash(message), filepath.ToSlash(wantPath))
}

func closePolicyFile(t *testing.T, file interface{ Close() error }) {
	t.Helper()

	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func assertPolicyCorruptJournalError(t *testing.T, err error, fileName string, line int) {
	t.Helper()

	var corruptErr *journal.CorruptJournalError
	if !errors.As(err, &corruptErr) {
		t.Fatalf("KnownRoles() error = %T, want *journal.CorruptJournalError", err)
	}

	wantFile := filepath.ToSlash(filepath.Join(".fid", "journal", fileName))
	if corruptErr.File != wantFile {
		t.Fatalf("corrupt file = %q, want %q", corruptErr.File, wantFile)
	}
	if corruptErr.Line != line {
		t.Fatalf("corrupt line = %d, want %d", corruptErr.Line, line)
	}
	if !strings.Contains(corruptErr.Error(), wantFile) {
		t.Fatalf("CorruptJournalError.Error() = %q, want path %q", corruptErr.Error(), wantFile)
	}
}

func testPolicyMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}
