package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestInitCommandCreatesWorkspaceFiles(t *testing.T) {
	workspace := t.TempDir()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"init"})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	requiredPaths := []string{
		filepath.Join(workspace, ".fid", "journal"),
		filepath.Join(workspace, ".fid", "config.toml"),
		filepath.Join(workspace, ".fid", "pricing.toml"),
		filepath.Join(workspace, ".fid", "fid.sqlite"),
		filepath.Join(workspace, ".fid", ".gitignore"),
	}
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("Stat(%q) error = %v", path, err)
		}
	}

	gitignorePath := filepath.Join(workspace, ".fid", ".gitignore")
	payload, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", gitignorePath, err)
	}
	if got := string(payload); !strings.Contains(got, "fid.sqlite") {
		t.Fatalf(".gitignore = %q, want fid.sqlite entry", got)
	}

	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if projectConfig.Prefix != strings.ToUpper(filepath.Base(workspace)) {
		t.Fatalf("project prefix = %q, want %q", projectConfig.Prefix, strings.ToUpper(filepath.Base(workspace)))
	}

	pricingConfig, err := config.ReadPricingConfig(filepath.Join(workspace, ".fid", "pricing.toml"))
	if err != nil {
		t.Fatalf("ReadPricingConfig() error = %v", err)
	}
	if len(pricingConfig.Models) != 3 {
		t.Fatalf("pricing models = %d, want %d", len(pricingConfig.Models), 3)
	}
	if got := stderr.String(); !strings.Contains(got, "initialized fid project in .fid/") {
		t.Fatalf("stderr = %q, want init confirmation", got)
	}
}

func TestInitCommandStoresConfiguredOrDefaultPrefix(t *testing.T) {
	defaultWorkspace := t.TempDir()
	runInitCommand(t, defaultWorkspace, "init")

	defaultConfig, err := config.ReadProjectConfig(filepath.Join(defaultWorkspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() default error = %v", err)
	}
	wantDefault := strings.ToUpper(filepath.Base(defaultWorkspace))
	if defaultConfig.Prefix != wantDefault {
		t.Fatalf("default prefix = %q, want %q", defaultConfig.Prefix, wantDefault)
	}

	customWorkspace := t.TempDir()
	runInitCommand(t, customWorkspace, "init", "--prefix", "ABC")

	customConfig, err := config.ReadProjectConfig(filepath.Join(customWorkspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() custom error = %v", err)
	}
	if customConfig.Prefix != "ABC" {
		t.Fatalf("custom prefix = %q, want %q", customConfig.Prefix, "ABC")
	}
}

func TestInitCommandSeedsExactlyOneBootstrapEventAndSevenRoles(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	payload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	lines := splitJournalLines(payload)
	if len(lines) != 1 {
		t.Fatalf("journal lines = %d, want %d", len(lines), 1)
	}

	var event journal.Event
	if err := json.Unmarshal([]byte(lines[0]), &event); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if event.EventType != "project.initialized" {
		t.Fatalf("event_type = %q, want %q", event.EventType, "project.initialized")
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var projectedRoles []string
	rows, err := conn.Query(`SELECT name FROM roles ORDER BY name`)
	if err != nil {
		t.Fatalf("Query() roles error = %v", err)
	}
	defer closeRows(t, rows)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Scan() role error = %v", err)
		}
		projectedRoles = append(projectedRoles, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows.Err() = %v", err)
	}

	wantRoles := []string{
		"developer",
		"human",
		"qa_lead",
		"qa_reviewer",
		"tl",
		"tl_lite",
		"workflow_runner",
	}
	if !reflect.DeepEqual(projectedRoles, wantRoles) {
		t.Fatalf("projected roles = %v, want %v", projectedRoles, wantRoles)
	}
}

func TestInitCommandRejectsReinitWithoutChanges(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	before := snapshotWorkspaceState(t, filepath.Join(workspace, ".fid"))

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"init"})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	err = cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want re-init rejection")
	}
	if got := stderr.String(); !strings.Contains(got, "error: .fid/ already exists") {
		t.Fatalf("stderr = %q, want already exists error", got)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	after := snapshotWorkspaceState(t, filepath.Join(workspace, ".fid"))
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("workspace changed on re-init\nbefore=%v\nafter=%v", before, after)
	}
}

func TestInitCommandRebuildRoundTripRecreatesRolesAndProjectionMeta(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	journalPath := mustJournalPath(t, workspace)

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	beforeRoles := readRoleNames(t, conn)
	beforeMeta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}
	closeCommandDB(t, conn)

	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runInitCommand(t, workspace, "rebuild", "--by", "developer")

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeCommandDB(t, rebuiltConn)

	afterRoles := readRoleNames(t, rebuiltConn)
	afterMeta, found, err := fiddb.ReadProjectionMeta(rebuiltConn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() rebuilt error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() rebuilt found = false, want true")
	}

	if !reflect.DeepEqual(afterRoles, beforeRoles) {
		t.Fatalf("rebuilt roles = %v, want %v", afterRoles, beforeRoles)
	}

	if beforeMeta.JournalFile != afterMeta.JournalFile ||
		beforeMeta.JournalOffset != afterMeta.JournalOffset ||
		beforeMeta.LastEventID != afterMeta.LastEventID ||
		beforeMeta.LastEventHash != afterMeta.LastEventHash ||
		beforeMeta.ReplayStatus != afterMeta.ReplayStatus {
		t.Fatalf("rebuilt projection_meta = %+v, want %+v", afterMeta, beforeMeta)
	}

	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) rebuilt error = %v", journalPath, err)
	}
	if !bytes.Equal(afterJournal, beforeJournal) {
		t.Fatal("journal bytes changed after rebuild")
	}
}

func runInitCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	return stdout.String(), stderr.String()
}

func mustJournalPath(t *testing.T, workspace string) string {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(workspace, ".fid", "journal", "*.jsonl"))
	if err != nil {
		t.Fatalf("Glob() error = %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("journal files = %d, want %d", len(matches), 1)
	}
	return matches[0]
}

func splitJournalLines(payload []byte) []string {
	lines := strings.Split(strings.TrimSpace(string(payload)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	return lines
}

func snapshotWorkspaceState(t *testing.T, root string) map[string]string {
	t.Helper()

	state := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		if info.IsDir() {
			state[filepath.ToSlash(relative)+"/"] = "dir"
			return nil
		}

		payload, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		state[filepath.ToSlash(relative)] = string(payload)
		return nil
	})
	if err != nil {
		t.Fatalf("Walk(%q) error = %v", root, err)
	}

	return state
}

func readRoleNames(t *testing.T, conn *sql.DB) []string {
	t.Helper()

	rows, err := conn.Query(`SELECT name FROM roles ORDER BY name`)
	if err != nil {
		t.Fatalf("Query() roles error = %v", err)
	}
	defer closeRows(t, rows)

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Scan() role error = %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows.Err() = %v", err)
	}

	sort.Strings(names)
	return names
}

func closeRows(t *testing.T, rows interface{ Close() error }) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Close() rows error = %v", err)
	}
}
