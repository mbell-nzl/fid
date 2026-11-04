package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestMaintenanceCheckJSONMatchesContract(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	healthy := executeCommandForTest(workspace, "--json", "maintenance", "check")
	if healthy.err != nil || healthy.stderr != "" {
		t.Fatalf("healthy maintenance.check = %+v, want successful stdout-only envelope", healthy)
	}
	healthyEnvelope := decodeMaintenanceEnvelope(t, healthy.stdout)
	assertExactMaintenanceEnvelopeKeys(t, healthyEnvelope, []string{"contract_version", "ok", "operation", "result", "event"})
	if healthyEnvelope["contract_version"] != "fid-cli.v1" || healthyEnvelope["ok"] != true || healthyEnvelope["operation"] != "maintenance.check" || healthyEnvelope["event"] != nil {
		t.Fatalf("healthy maintenance.check envelope = %#v, want fid-cli.v1 read success", healthyEnvelope)
	}
	healthyResult := maintenanceResultFromEnvelope(t, healthyEnvelope, "result")
	if healthyResult["status"] != "pass" {
		t.Fatalf("healthy maintenance.check result = %#v, want status pass", healthyResult)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	seedCorruptWorkflowFixture(t, conn, "FID-bad001")

	failing := executeCommandForTest(workspace, "--json", "maintenance", "check")
	if failing.err == nil || failing.stdout != "" {
		t.Fatalf("failing maintenance.check = %+v, want failing stderr-only envelope", failing)
	}
	failingEnvelope := decodeMaintenanceEnvelope(t, failing.stderr)
	assertExactMaintenanceEnvelopeKeys(t, failingEnvelope, []string{"contract_version", "ok", "operation", "error"})
	if failingEnvelope["contract_version"] != "fid-cli.v1" || failingEnvelope["ok"] != false || failingEnvelope["operation"] != "maintenance.check" {
		t.Fatalf("failing maintenance.check envelope = %#v, want fid-cli.v1 failure", failingEnvelope)
	}
	errorObject, ok := failingEnvelope["error"].(map[string]any)
	if !ok {
		t.Fatalf("failing maintenance.check error = %#v, want object", failingEnvelope["error"])
	}
	if errorObject["code"] != "projection-divergence" || errorObject["category"] != "integrity" || errorObject["retryable"] != false || errorObject["subject"] != nil {
		t.Fatalf("maintenance.check failure error = %#v, want projection-divergence integrity failure", errorObject)
	}
	details, ok := errorObject["details"].(map[string]any)
	if !ok {
		t.Fatalf("maintenance.check failure details = %#v, want result object", errorObject["details"])
	}
	if details["status"] != "fail" {
		t.Fatalf("maintenance.check failure details = %#v, want status fail", details)
	}
	assertMaintenanceFindings(t, details, "errors")
}

func TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	seedStartedIssue(t, conn, "FID-warn001", time.Now().UTC().Add(-301*time.Second))

	warn := executeCommandForTest(workspace, "--json", "maintenance", "check")
	if warn.err != nil || warn.stderr != "" {
		t.Fatalf("warn-only maintenance.check = %+v, want successful stdout envelope", warn)
	}
	warnResult := maintenanceResultFromEnvelope(t, decodeMaintenanceEnvelope(t, warn.stdout), "result")
	if warnResult["status"] != "warn" {
		t.Fatalf("warn-only result = %#v, want status warn", warnResult)
	}

	promoted := executeCommandForTest(workspace, "--json", "maintenance", "check", "--warn-as-error")
	if promoted.err == nil || promoted.stdout != "" {
		t.Fatalf("promoted maintenance.check = %+v, want failing stderr envelope", promoted)
	}
	promotedEnvelope := decodeMaintenanceEnvelope(t, promoted.stderr)
	promotedError := promotedEnvelope["error"].(map[string]any)
	promotedResult, ok := promotedError["details"].(map[string]any)
	if !ok || promotedError["code"] != "projection-divergence" || promotedResult["status"] != "fail" {
		t.Fatalf("promoted maintenance.check = %#v, want projection-divergence with status fail", promotedEnvelope)
	}
	for _, key := range []string{"errors", "warnings", "info"} {
		if !reflect.DeepEqual(warnResult[key], promotedResult[key]) {
			t.Fatalf("maintenance %s differs across warning branches: warn=%#v promoted=%#v", key, warnResult[key], promotedResult[key])
		}
	}
}

// Serial: the 299/300-second fixture ages against the wall clock across two
// command runs, leaving about one second of margin that parallel CPU contention
// can consume. The deterministic boundary is proven in internal/maintenance.
func TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat(t *testing.T) { //nolint:paralleltest // wall-clock boundary fixture, see above.
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	configPath := filepath.Join(config.LedgerDir(workspace), "config.toml")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	updatedConfig := strings.Replace(string(configBytes), "abandoned_attempt_grace_seconds = 300", "abandoned_attempt_grace_seconds = 0", 1)
	if err := os.WriteFile(configPath, []byte(updatedConfig), 0o644); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	fixtureNow := time.Now().UTC()
	seedStartedIssue(t, conn, "FID-299", fixtureNow.Add(-299*time.Second))
	seedStartedIssue(t, conn, "FID-300", fixtureNow.Add(-300*time.Second))

	checked := executeCommandForTest(workspace, "--json", "maintenance", "check")
	if checked.err != nil {
		t.Fatalf("maintenance.check = %+v, want warn success", checked)
	}
	result := maintenanceResultFromEnvelope(t, decodeMaintenanceEnvelope(t, checked.stdout), "result")
	if hasMaintenanceFinding(t, result, "warnings", "started_without_claim", "FID-299") || !hasMaintenanceFinding(t, result, "warnings", "started_without_claim", "FID-300") {
		t.Fatalf("started-without-claim warnings = %#v, want only fixed 300-second boundary", result["warnings"])
	}

	promoted := executeCommandForTest(workspace, "--json", "maintenance", "check", "--warn-as-error")
	if promoted.err == nil || promoted.stdout != "" {
		t.Fatalf("promoted maintenance.check = %+v, want failing stderr envelope", promoted)
	}
	promotedError := decodeMaintenanceEnvelope(t, promoted.stderr)["error"].(map[string]any)
	promotedResult, ok := promotedError["details"].(map[string]any)
	if !ok || promotedError["code"] != "projection-divergence" || promotedResult["status"] != "fail" {
		t.Fatalf("promoted maintenance.check = %#v, want projection-divergence with status fail", promotedError)
	}
	for _, key := range []string{"errors", "warnings", "info"} {
		if !reflect.DeepEqual(result[key], promotedResult[key]) {
			t.Fatalf("maintenance %s differs across fixed-grace branches: warn=%#v promoted=%#v", key, result[key], promotedResult[key])
		}
	}

	format := executeCommandForTest(workspace, "--json", "maintenance", "check", "--format", "json")
	if format.err == nil || format.stdout != "" {
		t.Fatalf("maintenance.check --format = %+v, want rejected request", format)
	}
	formatError := decodeMaintenanceEnvelope(t, format.stderr)["error"].(map[string]any)
	if formatError["code"] != "invalid-argument" {
		t.Fatalf("maintenance.check --format error = %#v, want invalid-argument", formatError)
	}
}

func TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	human := executeCommandForTest(workspace, "maintenance", "check")
	if human.err != nil || human.stderr != "" || !strings.Contains(human.stdout, "maintenance check: pass") {
		t.Fatalf("human maintenance.check = %+v, want stdout-only success", human)
	}
	vacuum := executeCommandForTest(workspace, "--json", "maintenance", "vacuum")
	if vacuum.err == nil || vacuum.stdout != "" {
		t.Fatalf("maintenance.vacuum = %+v, want typed unavailable rejection", vacuum)
	}
	vacuumError := decodeMaintenanceEnvelope(t, vacuum.stderr)["error"].(map[string]any)
	if vacuumError["code"] != "capability-unavailable" {
		t.Fatalf("maintenance.vacuum error = %#v, want capability-unavailable", vacuumError)
	}
}

func TestMaintenanceCheckUsesExplicitProjectRoot(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(t.TempDir())
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", workspace, "--json", "maintenance", "check"})
	if err := cmd.Execute(); err != nil || stderr.Len() != 0 {
		t.Fatalf("explicit-root maintenance.check = (%v, %q), want stdout-only success", err, stderr.String())
	}
	envelope := decodeMaintenanceEnvelope(t, stdout.String())
	if envelope["operation"] != "maintenance.check" || envelope["event"] != nil {
		t.Fatalf("explicit-root maintenance envelope = %#v, want maintenance.check read response", envelope)
	}
}

func TestMaintenanceCheckRejectsFindingLimitWithoutTruncation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	startedAt := time.Now().UTC().Add(-301 * time.Second)
	for index := 0; index < 1001; index++ {
		seedStartedIssue(t, conn, "FID-limit-"+strconv.Itoa(index), startedAt)
	}

	checked := executeCommandForTest(workspace, "--json", "maintenance", "check")
	if checked.err == nil || checked.stdout != "" {
		t.Fatalf("maintenance.check over limit = %+v, want stderr-only failure", checked)
	}
	errorObject := decodeMaintenanceEnvelope(t, checked.stderr)["error"].(map[string]any)
	details := errorObject["details"].(map[string]any)
	if errorObject["code"] != "result-limit-exceeded" || details["collection"] != "maintenance.warnings" || details["limit"] != json.Number("1000") || details["observed"] != json.Number("1001") {
		t.Fatalf("maintenance.check over-limit error = %#v, want exact finding limit details", errorObject)
	}
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM issues`).Scan(&count); err != nil || count != 1001 {
		t.Fatalf("issues after over-limit maintenance check = %d, err = %v; want 1001 unchanged rows", count, err)
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

func seedStartedIssue(t *testing.T, conn *sql.DB, id string, startedAt time.Time) {
	t.Helper()
	value := startedAt.UTC().Format(time.RFC3339Nano)
	if _, err := conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at) VALUES (?, 'task', 'Started without claim', 'task-default', 'in_progress', ?, ?)`, id, value, value); err != nil {
		t.Fatalf("seed warning issue: %v", err)
	}
}

func decodeMaintenanceEnvelope(t *testing.T, output string) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(output))
	decoder.UseNumber()
	var envelope map[string]any
	if err := decoder.Decode(&envelope); err != nil {
		t.Fatalf("decode maintenance envelope %q: %v", output, err)
	}
	return envelope
}

func assertExactMaintenanceEnvelopeKeys(t *testing.T, envelope map[string]any, want []string) {
	t.Helper()
	got := make([]string, 0, len(envelope))
	for key := range envelope {
		got = append(got, key)
	}
	if !reflect.DeepEqual(sortedStrings(got), sortedStrings(want)) {
		t.Fatalf("envelope keys = %#v, want %#v", got, want)
	}
}

func maintenanceResultFromEnvelope(t *testing.T, envelope map[string]any, key string) map[string]any {
	t.Helper()
	result, ok := envelope[key].(map[string]any)
	if !ok {
		t.Fatalf("maintenance envelope %s = %#v, want object", key, envelope[key])
	}
	return result
}

func assertMaintenanceFindings(t *testing.T, result map[string]any, key string) {
	t.Helper()
	findings, ok := result[key].([]any)
	if !ok || len(findings) == 0 {
		t.Fatalf("maintenance result %s = %#v, want non-empty finding list", key, result[key])
	}
	for _, value := range findings {
		finding, ok := value.(map[string]any)
		if !ok || !reflect.DeepEqual(sortedStrings(mapKeys(finding)), []string{"check", "entity_id", "message"}) {
			t.Fatalf("maintenance finding = %#v, want exact Finding schema", value)
		}
	}
}

func hasMaintenanceFinding(t *testing.T, result map[string]any, severity, check, entityID string) bool {
	t.Helper()
	findings, ok := result[severity].([]any)
	if !ok {
		t.Fatalf("maintenance result %s = %#v, want array", severity, result[severity])
	}
	for _, value := range findings {
		finding := value.(map[string]any)
		if finding["check"] == check && finding["entity_id"] == entityID {
			return true
		}
	}
	return false
}

func mapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func sortedStrings(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}
