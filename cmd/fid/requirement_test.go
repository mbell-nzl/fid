package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestRequirementSyncCreatesPointersInOneJournalEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	writeRequirementFixture(t, workspace, "project-docs/specs/sample_requirements.md", `# Requirements

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | The system shall support requirement pointer records. | Must | v0.1 |
| PRD-LINK-002 | The system shall support links from issues to requirements. | Must | v0.1 |

## 4.3 Persistence Rules

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| TRD-ARCH-010 | Every mutation flows through the journal write path. | Must | v0.1 |
`)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"requirement", "sync",
		"--source-doc", filepath.ToSlash(filepath.Join("project-docs", "specs", "sample_requirements.md")),
		"--by", "tl",
	)
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(stdout, "synced requirements from") {
		t.Fatalf("stdout = %q, want sync confirmation", stdout)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after sync = %d, want %d", linesAfter, linesBefore+1)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows := loadRequirementRows(t, conn)
	if len(rows) != 3 {
		t.Fatalf("requirements row count = %d, want 3", len(rows))
	}

	assertRequirementRow(t, rows[0], requirementRow{
		ID:             "PRD-GOV-006",
		Type:           "PRD",
		Summary:        "The system shall support requirement pointer records.",
		SourceDoc:      "project-docs/specs/sample_requirements.md",
		SourceSection:  stringPointer("3.11 Governance"),
		Status:         "active",
		HasSyncedHash:  true,
		LastSyncedHash: rows[0].LastSyncedHash,
	})
	assertRequirementRow(t, rows[1], requirementRow{
		ID:             "PRD-LINK-002",
		Type:           "PRD",
		Summary:        "The system shall support links from issues to requirements.",
		SourceDoc:      "project-docs/specs/sample_requirements.md",
		SourceSection:  stringPointer("3.11 Governance"),
		Status:         "active",
		HasSyncedHash:  true,
		LastSyncedHash: rows[1].LastSyncedHash,
	})
	assertRequirementRow(t, rows[2], requirementRow{
		ID:             "TRD-ARCH-010",
		Type:           "TRD",
		Summary:        "Every mutation flows through the journal write path.",
		SourceDoc:      "project-docs/specs/sample_requirements.md",
		SourceSection:  stringPointer("4.3 Persistence Rules"),
		Status:         "active",
		HasSyncedHash:  true,
		LastSyncedHash: rows[2].LastSyncedHash,
	})

	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "requirement.synced" {
		t.Fatalf("event_type = %q, want %q", event.EventType, "requirement.synced")
	}
	if len(event.Effects) != 4 {
		t.Fatalf("effect count = %d, want requirement effects plus idempotency result", len(event.Effects))
	}
	for _, effect := range event.Effects {
		if effect.Op != "insert" {
			t.Fatalf("effect op = %q, want insert", effect.Op)
		}
		if effect.Table != "requirements" {
			if effect.Table != "idempotency_results" {
				t.Fatalf("effect table = %q, want requirements or idempotency_results", effect.Table)
			}
		}
	}
}

func TestRequirementSyncRejectsMalformedInputBeforeAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		body          string
		args          []string
		wantSubstring string
	}{
		{
			name: "no parseable ids",
			body: `# Notes

This file contains prose only.
`,
			args:          []string{"--source-doc", "project-docs/specs/bad.md", "--by", "tl"},
			wantSubstring: `error: requirement sync: no parseable requirement ids found in "project-docs/specs/bad.md"`,
		},
		{
			name: "duplicate ids",
			body: `# Duplicate

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | First summary. | Must | v0.1 |
| PRD-GOV-006 | Duplicate summary. | Must | v0.1 |
`,
			args:          []string{"--source-doc", "project-docs/specs/bad.md", "--by", "tl"},
			wantSubstring: `error: requirement sync: duplicate requirement id "PRD-GOV-006" in source document`,
		},
		{
			name: "unknown type",
			body: `# Unknown Type

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer summary. | Must | v0.1 |
`,
			args:          []string{"--source-doc", "project-docs/specs/bad.md", "--type", "OPS", "--by", "tl"},
			wantSubstring: `error: requirement sync: --type must be PRD or TRD, got "OPS"`,
		},
		{
			name: "missing id from source",
			body: `# Missing ID

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer summary. | Must | v0.1 |
`,
			args:          []string{"--source-doc", "project-docs/specs/bad.md", "--id", "TRD-ARCH-010", "--by", "tl"},
			wantSubstring: `error: requirement sync: requirement id "TRD-ARCH-010" not found in "project-docs/specs/bad.md"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			writeRequirementFixture(t, workspace, "project-docs/specs/bad.md", tt.body)

			journalPath := mustJournalPath(t, workspace)
			linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

			commandArgs := append([]string{"requirement", "sync"}, tt.args...)
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, commandArgs...)
			if err == nil {
				t.Fatal("Execute() error = nil, want rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, tt.wantSubstring) {
				t.Fatalf("stderr = %q, want substring %q", stderr, tt.wantSubstring)
			}

			linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			if linesAfter != linesBefore {
				t.Fatalf("journal lines after rejection = %d, want %d", linesAfter, linesBefore)
			}

			conn, openErr := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
			if openErr != nil {
				t.Fatalf("Open() error = %v", openErr)
			}
			defer closeCommandDB(t, conn)
			if got := len(loadRequirementRows(t, conn)); got != 0 {
				t.Fatalf("requirements row count after rejection = %d, want 0", got)
			}
		})
	}
}

func TestRequirementSyncAuthorizationAndPolicyShowParity(t *testing.T) {
	t.Parallel()
	t.Run("human allowed", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)
		writeRequirementFixture(t, workspace, "project-docs/specs/auth.md", `# Auth

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer summary. | Must | v0.1 |
`)

		stdout, stderr := runFIDCommand(t, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/auth.md", "--by", "human")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}
		for _, want := range []string{
			"synced requirements from project-docs/specs/auth.md",
			"created: 1",
			"updated: 0",
			"unchanged: 0",
		} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("stdout = %q, want substring %q", stdout, want)
			}
		}
	})

	t.Run("denied and unknown reject before append", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)
		writeRequirementFixture(t, workspace, "project-docs/specs/auth.md", `# Auth

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer summary. | Must | v0.1 |
`)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		for _, actor := range []string{"agent:developer", "agent:intern"} {
			stdout, stderr := runFIDCommand(t, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/auth.md", "--actor", actor)
			if stderr != "" || !strings.Contains(stdout, "synced requirements") {
				t.Fatalf("%s actor result = (%q, %q), want stdout success", actor, stdout, stderr)
			}
		}
		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
			t.Fatalf("journal lines after project-local actors = %d, want %d", got, linesBefore+1)
		}
	})

	t.Run("policy show includes requirement sync row", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		rows := legacyPolicyRows(t)
		if !commandPolicyCellAllowed(t, rows, "requirement.sync", "tl") {
			t.Fatal("policy show tl/requirement.sync = false, want true")
		}
		if !commandPolicyCellAllowed(t, rows, "requirement.sync", "human") {
			t.Fatal("policy show human/requirement.sync = false, want true")
		}
		if commandPolicyCellAllowed(t, rows, "requirement.sync", "developer") {
			t.Fatal("policy show developer/requirement.sync = true, want false")
		}
	})
}

type requirementRow struct {
	ID             string
	Type           string
	Summary        string
	SourceDoc      string
	SourceSection  *string
	LastSyncedHash string
	Status         string
	HasSyncedHash  bool
}

func loadRequirementRows(t *testing.T, conn *sql.DB) []requirementRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, type, summary, source_doc, source_section, COALESCE(last_synced_hash, ''), status
		FROM requirements
		ORDER BY id
	`)
	if err != nil {
		t.Fatalf("Query(requirements) error = %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Fatalf("Rows.Close() error = %v", err)
		}
	}()

	var items []requirementRow
	for rows.Next() {
		var (
			item          requirementRow
			sourceSection sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Summary,
			&item.SourceDoc,
			&sourceSection,
			&item.LastSyncedHash,
			&item.Status,
		); err != nil {
			t.Fatalf("Scan(requirements) error = %v", err)
		}
		if sourceSection.Valid {
			item.SourceSection = stringPointer(sourceSection.String)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows(requirements) error = %v", err)
	}
	return items
}

func assertRequirementRow(t *testing.T, got, want requirementRow) {
	t.Helper()

	if got.ID != want.ID {
		t.Fatalf("row id = %q, want %q", got.ID, want.ID)
	}
	if got.Type != want.Type {
		t.Fatalf("row type for %s = %q, want %q", got.ID, got.Type, want.Type)
	}
	if got.Summary != want.Summary {
		t.Fatalf("row summary for %s = %q, want %q", got.ID, got.Summary, want.Summary)
	}
	if got.SourceDoc != want.SourceDoc {
		t.Fatalf("row source_doc for %s = %q, want %q", got.ID, got.SourceDoc, want.SourceDoc)
	}
	if dereferenceCommandString(got.SourceSection) != dereferenceCommandString(want.SourceSection) {
		t.Fatalf("row source_section for %s = %#v, want %#v", got.ID, got.SourceSection, want.SourceSection)
	}
	if got.Status != want.Status {
		t.Fatalf("row status for %s = %q, want %q", got.ID, got.Status, want.Status)
	}
	if want.HasSyncedHash && got.LastSyncedHash == "" {
		t.Fatalf("row last_synced_hash for %s = empty, want populated hash", got.ID)
	}
}

func mustReadLastJournalEvent(t *testing.T, journalPath string) journal.Event {
	t.Helper()

	lines := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(lines) == 0 {
		t.Fatal("journal lines = 0, want at least one line")
	}

	var event journal.Event
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal(last journal event) error = %v", err)
	}
	return event
}

func dereferenceCommandString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func writeRequirementFixture(t *testing.T, workspace, relativePath, body string) {
	t.Helper()

	path := filepath.Join(workspace, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
