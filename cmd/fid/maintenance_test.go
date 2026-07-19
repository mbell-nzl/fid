package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestMaintenanceCheckCommandRendersHealthyAndFailingResults(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "check")
	if err != nil {
		t.Fatalf("healthy maintenance check error = %v", err)
	}
	if stderr != "" {
		t.Fatalf("healthy maintenance check stderr = %q, want empty", stderr)
	}

	var healthy struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(stdout), &healthy); err != nil {
		t.Fatalf("decode healthy maintenance output: %v", err)
	}
	if healthy.Status != "pass" {
		t.Fatalf("healthy status = %q, want pass", healthy.Status)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	seedCorruptWorkflowFixture(t, conn, "FID-bad001")

	stdout, stderr, err = runFIDCommandExpectError(t, workspace, "maintenance", "check")
	if err == nil {
		t.Fatal("failing maintenance check error = nil, want non-zero outcome")
	}
	if stderr != "" {
		t.Fatalf("failing maintenance check stderr = %q, want empty", stderr)
	}

	var failing struct {
		Status string `json:"status"`
		Errors []struct {
			Check string `json:"check"`
		} `json:"errors"`
	}
	if err := json.Unmarshal([]byte(stdout), &failing); err != nil {
		t.Fatalf("decode failing maintenance output: %v", err)
	}
	if failing.Status != "fail" {
		t.Fatalf("failing status = %q, want fail", failing.Status)
	}
	gotChecks := make([]string, len(failing.Errors))
	for index, finding := range failing.Errors {
		gotChecks[index] = finding.Check
	}
	if want := []string{"foreign_key_integrity", "foreign_key_integrity", "workflow_state"}; !reflect.DeepEqual(gotChecks, want) {
		t.Fatalf("failing checks = %#v, want %#v", gotChecks, want)
	}
}

func TestMaintenanceCheckCommandRendersAllFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	for _, formatName := range []string{"json", "toon", "markdown", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "check", "--format", formatName)
			if err != nil {
				t.Fatalf("healthy maintenance check error = %v", err)
			}
			if stderr != "" {
				t.Fatalf("healthy maintenance check stderr = %q, want empty", stderr)
			}
			switch formatName {
			case "json":
				if got := decodeCommandJSONStructure(t, stdout); got.(map[string]any)["status"] != "pass" {
					t.Fatalf("json status = %v, want pass", got.(map[string]any)["status"])
				}
			case "toon":
				if got := decodeCommandTOONStructure(t, stdout); got.(map[string]any)["status"] != "pass" {
					t.Fatalf("toon status = %v, want pass", got.(map[string]any)["status"])
				}
			case "yaml":
				if got := decodeCommandYAMLStructure(t, stdout); got.(map[string]any)["status"] != "pass" {
					t.Fatalf("yaml status = %v, want pass", got.(map[string]any)["status"])
				}
			case "markdown":
				for _, want := range []string{"# Maintenance Check", "**Status:** pass", "## Errors", "## Warnings", "## Info"} {
					if !strings.Contains(stdout, want) {
						t.Fatalf("markdown output = %q, want substring %q", stdout, want)
					}
				}
			}
		})
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	seedCorruptWorkflowFixture(t, conn, "FID-bad002")

	for _, formatName := range []string{"json", "toon", "markdown", "yaml"} {
		t.Run("failing "+formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "check", "--format", formatName)
			if err == nil {
				t.Fatal("failing maintenance check error = nil, want non-zero outcome")
			}
			if stderr != "" {
				t.Fatalf("failing maintenance check stderr = %q, want empty", stderr)
			}
			switch formatName {
			case "json":
				if got := decodeCommandJSONStructure(t, stdout).(map[string]any)["status"]; got != "fail" {
					t.Fatalf("json status = %v, want fail", got)
				}
			case "toon":
				if got := decodeCommandTOONStructure(t, stdout).(map[string]any)["status"]; got != "fail" {
					t.Fatalf("toon status = %v, want fail", got)
				}
			case "yaml":
				if got := decodeCommandYAMLStructure(t, stdout).(map[string]any)["status"]; got != "fail" {
					t.Fatalf("yaml status = %v, want fail", got)
				}
			case "markdown":
				if !strings.Contains(stdout, "**Status:** fail") || !strings.Contains(stdout, "workflow_state") {
					t.Fatalf("markdown output = %q, want fail status and workflow finding", stdout)
				}
			}
		})
	}
}

func seedCorruptWorkflowFixture(t *testing.T, conn *sql.DB, issueID string) {
	t.Helper()

	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("bypass fixture foreign keys: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA ignore_check_constraints = ON`); err != nil {
		t.Fatalf("bypass fixture checks: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at) VALUES (?, 'task', 'Corrupt workflow', 'missing', 'todo', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`, issueID); err != nil {
		t.Fatalf("seed corrupt issue: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA ignore_check_constraints = OFF`); err != nil {
		t.Fatalf("restore fixture checks: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("restore fixture foreign keys: %v", err)
	}
}

func TestMaintenanceCheckCommandPromotesWarningsWithoutRewritingThem(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at) VALUES ('FID-warn001', 'task', 'Started without claim', 'task-default', 'in_progress', '2020-01-01T00:00:00Z', '2020-01-01T00:00:00Z')`); err != nil {
		t.Fatalf("seed warning issue: %v", err)
	}

	defaultOutput, defaultStderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "check")
	if err != nil {
		t.Fatalf("default warning maintenance check error = %v", err)
	}
	if defaultStderr != "" {
		t.Fatalf("default warning maintenance check stderr = %q, want empty", defaultStderr)
	}

	promotedOutput, promotedStderr, err := runFIDCommandExpectError(t, workspace, "maintenance", "check", "--warn-as-error")
	if err == nil {
		t.Fatal("promoted warning maintenance check error = nil, want non-zero outcome")
	}
	if promotedStderr != "" {
		t.Fatalf("promoted warning maintenance check stderr = %q, want empty", promotedStderr)
	}

	var defaultResult, promotedResult struct {
		Status   string            `json:"status"`
		Errors   []json.RawMessage `json:"errors"`
		Warnings []json.RawMessage `json:"warnings"`
	}
	if err := json.Unmarshal([]byte(defaultOutput), &defaultResult); err != nil {
		t.Fatalf("decode default warning output: %v", err)
	}
	if err := json.Unmarshal([]byte(promotedOutput), &promotedResult); err != nil {
		t.Fatalf("decode promoted warning output: %v", err)
	}
	if defaultResult.Status != "warn" || promotedResult.Status != "fail" {
		t.Fatalf("warning statuses = %q, %q; want warn, fail", defaultResult.Status, promotedResult.Status)
	}
	if len(defaultResult.Errors) != 0 || len(promotedResult.Errors) != 0 {
		t.Fatalf("warning errors = %v, %v; want no errors", defaultResult.Errors, promotedResult.Errors)
	}
	if len(defaultResult.Warnings) == 0 || !reflect.DeepEqual(defaultResult.Warnings, promotedResult.Warnings) {
		t.Fatalf("warning payloads = %s, %s; want identical non-empty warnings", defaultResult.Warnings, promotedResult.Warnings)
	}
}
