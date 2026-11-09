package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestInitCommandCreatesWorkspaceFiles(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"init", "--prefix", "FID"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	identity, err := config.ReadProjectIdentity(workspace)
	if err != nil {
		t.Fatal(err)
	}
	root := gittest.ResolvedPath(t, workspace)
	wantOutput := fmt.Sprintf("initialized fid project in %s/\nproject ID: %s\ncommit %s to keep this identity in clones\n",
		filepath.ToSlash(config.LedgerDir(root)), identity, filepath.Join(root, config.ProjectIdentityFileName))
	if got, want := stdout.String(), wantOutput; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	requiredPaths := []string{
		filepath.Join(workspace, config.ProjectIdentityFileName),
		filepath.Join(config.LedgerDir(workspace), "journal"),
		filepath.Join(config.LedgerDir(workspace), "artifacts"),
		filepath.Join(config.LedgerDir(workspace), "config.toml"),
		filepath.Join(config.LedgerDir(workspace), "fid.sqlite"),
		filepath.Join(config.LedgerDir(workspace), ".gitignore"),
	}
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("Stat(%q) error = %v", path, err)
		}
	}
	assertNoInitStagingResidue(t, workspace)
	configPath := filepath.Join(config.LedgerDir(workspace), "config.toml")
	configPayload, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, err)
	}
	for _, want := range []string{
		"config_version = 1",
		"[project]\nid = ",
		"issue_prefix = ",
		"[storage]\nartifact_root = \"artifacts\"",
		"[maintenance]\nabandoned_attempt_grace_seconds = 300",
	} {
		if !strings.Contains(string(configPayload), want) {
			t.Fatalf("config.toml = %q, want %q", configPayload, want)
		}
	}
	if strings.Contains(string(configPayload), "\nprefix =") || strings.Contains(string(configPayload), "started_without_claim") {
		t.Fatalf("config.toml = %q, contains retired key", configPayload)
	}

	gitignorePath := filepath.Join(config.LedgerDir(workspace), ".gitignore")
	payload, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", gitignorePath, err)
	}
	if got := string(payload); !strings.Contains(got, "fid.sqlite") {
		t.Fatalf(".gitignore = %q, want fid.sqlite entry", got)
	}

	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if projectConfig.IssuePrefix != "FID" {
		t.Fatalf("project issue prefix = %q, want %q", projectConfig.IssuePrefix, "FID")
	}
	if projectConfig.ProjectID == "" {
		t.Fatal("project identity is empty")
	}
	if projectConfig.AbandonedAttemptGraceSeconds != config.DefaultAbandonedAttemptGraceSeconds {
		t.Fatalf("abandoned-attempt grace = %d, want %d", projectConfig.AbandonedAttemptGraceSeconds, config.DefaultAbandonedAttemptGraceSeconds)
	}
}

func TestInitCommandUsesExactProjectInitResponseEnvelope(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", workspace, "--json", "init", "--prefix", "FID"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	assertEnvelopeKeySet(t, stdout.Bytes(), []string{"contract_version", "event", "ok", "operation", "result"})
	assertNestedEnvelopeKeySet(t, stdout.Bytes(), "result", []string{"journal_path", "project", "project_identity_path"})
	assertNestedEnvelopeKeySet(t, stdout.Bytes(), "result.project", []string{"id", "issue_prefix", "root"})
	assertNestedEnvelopeKeySet(t, stdout.Bytes(), "event", []string{"id", "occurred_at"})

	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Result          struct {
			Project struct {
				ID          string `json:"id"`
				IssuePrefix string `json:"issue_prefix"`
				Root        string `json:"root"`
			} `json:"project"`
			JournalPath string `json:"journal_path"`
		} `json:"result"`
		Event *struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v\nstdout=%q", err, stdout.String())
	}
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != "project.init" {
		t.Fatalf("envelope header = %+v, want fid-cli.v1 successful project.init", envelope)
	}
	root := gittest.ResolvedPath(t, workspace)
	if envelope.Result.Project.ID == "" || envelope.Result.Project.IssuePrefix != "FID" || envelope.Result.Project.Root != root {
		t.Fatalf("project result = %+v, want generated ID, FID, and root %q", envelope.Result.Project, root)
	}
	if envelope.Event == nil || envelope.Event.ID == "" || envelope.Event.OccurredAt == "" {
		t.Fatalf("mutation event = %+v, want event identity", envelope.Event)
	}
	journalPath := filepath.ToSlash(filepath.Join(config.LedgerDir(root), config.JournalDirName, filepath.Base(mustJournalPath(t, workspace))))
	if envelope.Result.JournalPath != journalPath {
		t.Fatalf("journal_path = %q, want absolute path %q", envelope.Result.JournalPath, journalPath)
	}
	journalPayload, err := os.ReadFile(mustJournalPath(t, workspace))
	if err != nil {
		t.Fatalf("ReadFile(journal) error = %v", err)
	}
	var journalEvent journal.Event
	if err := json.Unmarshal([]byte(splitJournalLines(journalPayload)[0]), &journalEvent); err != nil {
		t.Fatalf("Unmarshal(journal event) error = %v", err)
	}
	if envelope.Event.ID != journalEvent.EventID || envelope.Event.OccurredAt != journalEvent.OccurredAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("envelope event = %+v, want journal event %s at %s", envelope.Event, journalEvent.EventID, journalEvent.OccurredAt.UTC().Format(time.RFC3339Nano))
	}

	stdout.Reset()
	stderr.Reset()
	cmd = newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", workspace, "--json", "init", "--prefix", "FID"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want initialized-project rejection")
	}
	if stdout.Len() != 0 {
		t.Fatalf("failure stdout = %q, want empty", stdout.String())
	}
	var failure struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Error           struct {
			Code     string `json:"code"`
			Category string `json:"category"`
			Details  struct {
				Path string `json:"path"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &failure); err != nil {
		t.Fatalf("Unmarshal(stderr) error = %v\nstderr=%q", err, stderr.String())
	}
	assertEnvelopeKeySet(t, stderr.Bytes(), []string{"contract_version", "error", "ok", "operation"})
	assertNestedEnvelopeKeySet(t, stderr.Bytes(), "error", []string{"category", "code", "details", "message", "retryable", "subject"})
	assertNestedEnvelopeKeySet(t, stderr.Bytes(), "error.details", []string{"path"})
	assertNestedEnvelopeKeySet(t, stderr.Bytes(), "error.subject", []string{"id", "type"})
	if failure.ContractVersion != "fid-cli.v1" || failure.OK || failure.Operation != "project.init" ||
		failure.Error.Code != "project-already-initialized" || failure.Error.Category != "conflict" || failure.Error.Details.Path != root {
		t.Fatalf("failure = %+v, want typed project-already-initialized envelope for %q", failure, root)
	}
}

func assertEnvelopeKeySet(t *testing.T, payload []byte, want []string) {
	t.Helper()
	var value map[string]json.RawMessage
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("Unmarshal key set error = %v", err)
	}
	assertJSONKeys(t, value, want)
}

func assertNestedEnvelopeKeySet(t *testing.T, payload []byte, path string, want []string) {
	t.Helper()
	var current json.RawMessage = payload
	for _, name := range strings.Split(path, ".") {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(current, &object); err != nil {
			t.Fatalf("Unmarshal %s error = %v", path, err)
		}
		var ok bool
		current, ok = object[name]
		if !ok {
			t.Fatalf("%s is absent", path)
		}
	}
	var value map[string]json.RawMessage
	if err := json.Unmarshal(current, &value); err != nil {
		t.Fatalf("Unmarshal %s object error = %v", path, err)
	}
	assertJSONKeys(t, value, want)
}

func assertJSONKeys(t *testing.T, value map[string]json.RawMessage, want []string) {
	t.Helper()
	got := make([]string, 0, len(value))
	for key := range value {
		got = append(got, key)
	}
	sort.Strings(got)
	sortedWant := append([]string(nil), want...)
	sort.Strings(sortedWant)
	if !reflect.DeepEqual(got, sortedWant) {
		t.Fatalf("JSON keys = %v, want %v", got, sortedWant)
	}
}

func TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging(t *testing.T) {
	t.Parallel()
	validWorkspace := gittest.NewWorkspace(t)
	runInitCommand(t, validWorkspace, "init", "--prefix", "A123456789012345")
	if _, err := os.Stat(config.LedgerDir(validWorkspace)); err != nil {
		t.Fatalf("valid-prefix workspace error = %v, want activated workspace", err)
	}

	invalidWorkspace := gittest.NewWorkspace(t)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(invalidWorkspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", invalidWorkspace, "--json", "init", "--prefix", "fid"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want malformed-prefix rejection")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	var envelope struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &envelope); err != nil {
		t.Fatalf("Unmarshal(stderr) error = %v\nstderr=%q", err, stderr.String())
	}
	if envelope.Error.Code != "invalid-argument" {
		t.Fatalf("error code = %q, want invalid-argument", envelope.Error.Code)
	}
	if _, err := os.Stat(config.LedgerDir(invalidWorkspace)); !os.IsNotExist(err) {
		t.Fatalf("invalid-prefix workspace state = %v, want no .fid", err)
	}
}

func TestInitCommandPrefixBoundaryMatrix(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		workspaceName string
		prefix        string
		valid         bool
		wantPrefix    string
	}{
		{name: "explicit one character", workspaceName: "project", prefix: "A", valid: true, wantPrefix: "A"},
		{name: "explicit 16 characters", workspaceName: "project", prefix: "A123456789012345", valid: true, wantPrefix: "A123456789012345"},
		{name: "explicit 17 characters", workspaceName: "project", prefix: "A1234567890123456"},
		{name: "explicit lowercase", workspaceName: "project", prefix: "fid"},
		{name: "explicit leading digit", workspaceName: "project", prefix: "1FID"},
		{name: "derived uppercase basename", workspaceName: "fid", valid: true, wantPrefix: "FID"},
		{name: "derived invalid basename", workspaceName: "invalid_prefix"},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			workspace := gittest.NewNamedWorkspace(t, test.workspaceName)
			args := []string{"--project", workspace, "--json", "init"}
			if test.prefix != "" {
				args = append(args, "--prefix", test.prefix)
			}
			stdout, stderr, err := executeInitCommand(t, args...)
			if test.valid {
				if err != nil {
					t.Fatalf("Execute(%v) error = %v", args, err)
				}
				if stderr != "" {
					t.Fatalf("stderr = %q, want empty", stderr)
				}
				assertProjectInitSuccess(t, stdout)
				assertActivatedConfigRoundTrip(t, workspace, test.wantPrefix)
				assertNoInitStagingResidue(t, workspace)
				return
			}
			if err == nil {
				t.Fatalf("Execute(%v) error = nil, want invalid-argument rejection", args)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			assertProjectInitFailureCode(t, stderr, "invalid-argument")
			if _, statErr := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(statErr) {
				t.Fatalf(".fid state = %v, want absent", statErr)
			}
			assertNoInitStagingResidue(t, workspace)
		})
	}
}

func TestInitCommandPreservesDanglingWorkspaceSymlinkCollision(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	danglingTarget := filepath.Join(workspace, "missing-fid-target")
	workspacePath := config.LedgerDir(workspace)
	if err := os.Symlink(danglingTarget, workspacePath); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	stdout, stderr, err := executeInitCommand(t, "--project", workspace, "--json", "init", "--prefix", "FID")
	if err == nil {
		t.Fatal("Execute() error = nil, want project-already-initialized rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	assertProjectInitFailureCode(t, stderr, "project-already-initialized")
	if got, readErr := os.Readlink(workspacePath); readErr != nil || got != danglingTarget {
		t.Fatalf("Readlink(%q) = %q, %v; want unchanged dangling target %q", workspacePath, got, readErr, danglingTarget)
	}
	if info, statErr := os.Lstat(workspacePath); statErr != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%q) = %v, %v; want unchanged symlink", workspacePath, info, statErr)
	}
	assertNoInitStagingResidue(t, workspace)
}

func TestInitCommandStoresConfiguredOrDefaultPrefix(t *testing.T) {
	t.Parallel()
	defaultWorkspace := gittest.NewWorkspace(t)
	runInitCommand(t, defaultWorkspace, "init", "--prefix", "FID")

	defaultConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(defaultWorkspace), "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() default error = %v", err)
	}
	if defaultConfig.IssuePrefix != "FID" {
		t.Fatalf("default issue prefix = %q, want %q", defaultConfig.IssuePrefix, "FID")
	}

	customWorkspace := gittest.NewWorkspace(t)
	runInitCommand(t, customWorkspace, "init", "--prefix", "ABC")

	customConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(customWorkspace), "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() custom error = %v", err)
	}
	if customConfig.IssuePrefix != "ABC" {
		t.Fatalf("custom issue prefix = %q, want %q", customConfig.IssuePrefix, "ABC")
	}
}

func TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewNamedWorkspace(t, "invalid_prefix")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"init"})

	if err := config.ValidateIssuePrefix(bootstrap.DefaultPrefix(workspace)); err == nil {
		t.Fatal("test workspace unexpectedly has a valid default prefix")
	}

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want explicit-prefix rejection")
	}
	if got, want := stderr.String(), "error: invalid project initialization argument\nreason: invalid-argument; argument=prefix, allowed=valid project.init input\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
	if _, statErr := os.Stat(config.LedgerDir(workspace)); !os.IsNotExist(statErr) {
		t.Fatalf(".fid state = %v, want no workspace", statErr)
	}
}

func TestInitCommandSeedsExactlyOneRoleFreeBootstrapEvent(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init", "--prefix", "FID")

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
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	if event.EventSchema != journal.EventSchemaV2 || event.ActorV2.ActorID != "fid-bootstrap" ||
		event.ActorV2.ActorType != "service" || event.ActorV2.Producer != "fid" ||
		event.ActorV2.ProvenanceRef != "fid:init" || event.ActorV2.SessionRef != nil ||
		event.Aggregate.ID != projectConfig.ProjectID || event.Aggregate.Type != "project" || event.Aggregate.Sequence != 1 ||
		event.CommandName != "project.init" || event.CausationID != nil || event.CorrelationID != nil ||
		event.Repository != nil || event.Idempotency != nil || len(event.Effects) != 0 {
		t.Fatalf("bootstrap Event = %#v, want exact role-free Event-v2 genesis", event)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	assertRoleTableAbsent(t, conn)
}

func TestInitCommandRejectsReinitWithoutChanges(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	before := snapshotWorkspaceState(t, workspace)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"init", "--prefix", "FID"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want re-init rejection")
	}
	if got, want := stderr.String(), "error: project is already initialized\n"+p24ExpectedReason("reason: project-already-initialized; path="+gittest.ResolvedPath(t, workspace))+"\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	after := snapshotWorkspaceState(t, workspace)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("workspace changed on re-init\nbefore=%v\nafter=%v", before, after)
	}
}

func TestInitCommandRebuildRoundTripPreservesRoleFreeBootstrapAndProjectionMeta(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	journalPath := mustJournalPath(t, workspace)

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	assertRoleTableAbsent(t, conn)
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

	runInitCommand(t, workspace, "rebuild")

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeCommandDB(t, rebuiltConn)

	assertRoleTableAbsent(t, rebuiltConn)
	afterMeta, found, err := fiddb.ReadProjectionMeta(rebuiltConn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() rebuilt error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() rebuilt found = false, want true")
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

func executeInitCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(t.TempDir())
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func assertProjectInitSuccess(t *testing.T, payload string) {
	t.Helper()
	var envelope struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
	}
	if err := json.Unmarshal([]byte(payload), &envelope); err != nil {
		t.Fatalf("Unmarshal(success) error = %v\npayload=%q", err, payload)
	}
	if !envelope.OK || envelope.Operation != "project.init" {
		t.Fatalf("success envelope = %+v, want successful project.init", envelope)
	}
}

func assertProjectInitFailureCode(t *testing.T, payload, wantCode string) {
	t.Helper()
	var envelope struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(payload), &envelope); err != nil {
		t.Fatalf("Unmarshal(failure) error = %v\npayload=%q", err, payload)
	}
	if envelope.OK || envelope.Operation != "project.init" || envelope.Error.Code != wantCode {
		t.Fatalf("failure envelope = %+v, want failed project.init with %q", envelope, wantCode)
	}
}

func assertActivatedConfigRoundTrip(t *testing.T, workspace, wantPrefix string) {
	t.Helper()
	path := filepath.Join(config.LedgerDir(workspace), "config.toml")
	project, err := config.ReadProjectConfig(path)
	if err != nil {
		t.Fatalf("ReadProjectConfig(%q) error = %v", path, err)
	}
	if project.IssuePrefix != wantPrefix {
		t.Fatalf("IssuePrefix = %q, want %q", project.IssuePrefix, wantPrefix)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	want := fmt.Sprintf(`config_version = 1

[project]
id = %q
issue_prefix = %q

[storage]
artifact_root = "artifacts"

[executor_limits]
max_output_bytes = %d
max_excerpt_bytes = %d

[maintenance]
abandoned_attempt_grace_seconds = %d
`, project.ProjectID, project.IssuePrefix, config.DefaultMaxOutputBytes, config.DefaultMaxExcerptBytes, config.DefaultAbandonedAttemptGraceSeconds)
	if string(payload) != want {
		t.Fatalf("activated config contents = %q, want exact target keys %q", payload, want)
	}
}

// assertNoInitStagingResidue checks that init left nothing behind: no staging
// directory beside the ledger in the git common directory, and nothing in the
// fresh worktree but .git and the project identity.
func assertNoInitStagingResidue(t *testing.T, workspace string) {
	t.Helper()
	commonDir := filepath.Dir(config.LedgerDir(workspace))
	entries, err := os.ReadDir(commonDir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", commonDir, err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "fid-staging-") {
			t.Fatalf("staging residue %q remains", entry.Name())
		}
	}
	entries, err = os.ReadDir(workspace)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", workspace, err)
	}
	for _, entry := range entries {
		if entry.Name() != ".git" && entry.Name() != config.ProjectIdentityFileName {
			t.Fatalf("init left %q in the worktree, want only .git and .fid.toml", entry.Name())
		}
	}
}

func runInitCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()
	if len(args) > 0 && args[0] == "init" && !containsArgument(args, "--prefix") {
		args = append(args, "--prefix", "FID")
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	return stdout.String(), stderr.String()
}

func containsArgument(args []string, value string) bool {
	for _, arg := range args {
		if arg == value {
			return true
		}
	}
	return false
}

func mustJournalPath(t *testing.T, workspace string) string {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), "journal", "*.jsonl"))
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

func assertRoleTableAbsent(t *testing.T, conn *sql.DB) {
	t.Helper()

	var name string
	err := conn.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'roles'`).Scan(&name)
	if err != sql.ErrNoRows {
		t.Fatalf("roles table lookup error = %v, want no roles table", err)
	}
}
