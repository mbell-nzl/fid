package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// backupFixture is one machine: a home, user config and cache in temp
// directories, and a local bare repository as the backup target. Commands
// get this environment through the backup seam, never the process's own.
type backupFixture struct {
	t                                 *testing.T
	base, home, configHome, cacheHome string
	target                            string
	env                               []string
}

type backupRecordForTest struct {
	Schema            string  `json:"schema"`
	StartedAt         string  `json:"started_at"`
	FinishedAt        string  `json:"finished_at"`
	Outcome           string  `json:"outcome"`
	ProjectID         *string `json:"project_id"`
	Commit            *string `json:"commit"`
	ArchiveDigest     *string `json:"archive_digest"`
	EventCount        *int    `json:"event_count"`
	TerminalEventHash *string `json:"terminal_event_hash"`
	Failure           *struct {
		Kind    string `json:"kind"`
		Message string `json:"message"`
	} `json:"failure"`
}

type backupFolderManifestForTest struct {
	Schema            string `json:"schema"`
	ProjectID         string `json:"project_id"`
	ArchiveDigest     string `json:"archive_digest"`
	ManifestSchema    string `json:"manifest_schema"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	ProjectionSchema  int    `json:"projection_schema"`
	EntryCount        int    `json:"entry_count"`
	GenesisEventHash  string `json:"genesis_event_hash"`
}

func newBackupFixture(t *testing.T) *backupFixture {
	t.Helper()
	base := t.TempDir()
	fixture := &backupFixture{t: t, base: base, home: filepath.Join(base, "home"),
		configHome: filepath.Join(base, "config"), cacheHome: filepath.Join(base, "cache"), target: filepath.Join(base, "backups.git")}
	if err := os.Mkdir(fixture.home, 0o700); err != nil {
		t.Fatal(err)
	}
	fixture.env = backupTestEnvironment(fixture.home, fixture.configHome, fixture.cacheHome)
	fixture.git("init", "--bare", "-q", "--template=", fixture.target)
	return fixture
}

// backupTestEnvironment is the process environment with HOME and the XDG
// directories replaced and every inherited GIT_ variable but the discovery
// guard removed, so no test reads the real user config, cache or git config.
func backupTestEnvironment(home, configHome, cacheHome string) []string {
	env := []string{}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name == "HOME" || name == "XDG_CONFIG_HOME" || name == "XDG_CACHE_HOME" ||
			strings.HasPrefix(name, "GIT_") && name != gittest.CeilingDirectoriesEnv {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "HOME="+home, "XDG_CONFIG_HOME="+configHome, "XDG_CACHE_HOME="+cacheHome)
}

// project makes and initializes a fresh project and returns its root and ID.
func (fixture *backupFixture) project() (string, string) {
	fixture.t.Helper()
	workspace := gittest.NewWorkspace(fixture.t)
	runInitCommand(fixture.t, workspace, "init", "--json")
	id, err := config.ReadProjectIdentity(workspace)
	if err != nil {
		fixture.t.Fatal(err)
	}
	return workspace, id
}

func (fixture *backupFixture) userConfigPath() string {
	return filepath.Join(fixture.configHome, "fid", "config.toml")
}

func (fixture *backupFixture) setTarget(target string) {
	fixture.t.Helper()
	if err := config.SetBackupTarget(fixture.userConfigPath(), target); err != nil {
		fixture.t.Fatal(err)
	}
}

func (fixture *backupFixture) backup(workspace string, args ...string) commandExecutionResult {
	return fixture.backupWaiting(workspace, time.Minute, args...)
}

func (fixture *backupFixture) backupWaiting(workspace string, wait time.Duration, args ...string) commandExecutionResult {
	return runBackupWith(fixture.env, workspace, wait, args...)
}

// runBackupWith runs fid backup in-process with environ as its environment.
func runBackupWith(environ []string, workspace string, wait time.Duration, args ...string) commandExecutionResult {
	root := newRootCommandAt(workspace)
	root.SetContext(context.WithValue(root.Context(), backupEnvironmentKey{}, backupEnvironment{environ: environ, lockWait: wait}))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(append(args, "backup"))
	err := root.Execute()
	return commandExecutionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

// git runs git with the fixture's environment, so the user's real git
// config never takes part.
func (fixture *backupFixture) git(args ...string) string {
	fixture.t.Helper()
	output, err := fixture.gitBytes(args...)
	if err != nil {
		var stderr []byte
		if exit, ok := err.(*exec.ExitError); ok {
			stderr = exit.Stderr
		}
		fixture.t.Fatalf("git %v: %v: %s%s", args, err, output, stderr)
	}
	return strings.TrimSpace(string(output))
}

func (fixture *backupFixture) gitBytes(args ...string) ([]byte, error) {
	command := exec.Command("git", args...)
	command.Env = fixture.env
	command.Dir = fixture.base
	return command.Output()
}

func (fixture *backupFixture) targetGit(args ...string) string {
	fixture.t.Helper()
	return fixture.git(append([]string{"--git-dir=" + fixture.target}, args...)...)
}

func (fixture *backupFixture) targetMain() string {
	fixture.t.Helper()
	output, err := fixture.gitBytes("--git-dir="+fixture.target, "rev-parse", "--verify", "-q", "refs/heads/main")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func decodeBackupSuccess(t *testing.T, result commandExecutionResult) (backupRecordForTest, string) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("backup = %v, stderr=%q", result.err, result.stderr)
	}
	var envelope struct {
		OK        bool            `json:"ok"`
		Operation string          `json:"operation"`
		Event     json.RawMessage `json:"event"`
		Result    struct {
			Backup     json.RawMessage `json:"backup"`
			RecordPath string          `json:"record_path"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode backup envelope %q: %v", result.stdout, err)
	}
	if !envelope.OK || envelope.Operation != "project.backup" || string(envelope.Event) != "null" {
		t.Fatalf("backup envelope = %s, want ok project.backup with no event", result.stdout)
	}
	return strictBackupRecord(t, envelope.Result.Backup), envelope.Result.RecordPath
}

func strictBackupRecord(t *testing.T, payload []byte) backupRecordForTest {
	t.Helper()
	var record backupRecordForTest
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&record); err != nil {
		t.Fatalf("decode backup record %q: %v", payload, err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(payload, &keys); err != nil || len(keys) != 10 {
		t.Fatalf("backup record %s has %d keys, want all 10 present: %v", payload, len(keys), err)
	}
	if record.Schema != "fid.backup-record.v1" {
		t.Fatalf("record schema = %q", record.Schema)
	}
	for _, stamp := range []string{record.StartedAt, record.FinishedAt} {
		if _, err := time.Parse(time.RFC3339Nano, stamp); err != nil || !strings.HasSuffix(stamp, "Z") {
			t.Fatalf("record time %q is not UTC RFC 3339: %v", stamp, err)
		}
	}
	return record
}

// backupRecords reads every line of the project's backup record.
func backupRecords(t *testing.T, workspace string) []backupRecordForTest {
	t.Helper()
	payload, err := os.ReadFile(filepath.Join(config.LedgerDir(workspace), "backups.jsonl"))
	if err != nil {
		t.Fatalf("read backup record: %v", err)
	}
	if !bytes.HasSuffix(payload, []byte("\n")) {
		t.Fatalf("backup record %q does not end in a newline", payload)
	}
	records := []backupRecordForTest{}
	for _, line := range bytes.Split(bytes.TrimSuffix(payload, []byte("\n")), []byte("\n")) {
		records = append(records, strictBackupRecord(t, line))
	}
	return records
}

func backupOutcomes(records []backupRecordForTest) []string {
	outcomes := []string{}
	for _, record := range records {
		outcomes = append(outcomes, record.Outcome)
	}
	return outcomes
}

func TestBackupPushesTheExportIntoTheProjectFolder(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, id := fixture.project()
	fixture.setTarget(fixture.target)

	record, recordPath := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	main := fixture.targetMain()
	if record.Outcome != "pushed" || record.Failure != nil || record.ProjectID == nil || *record.ProjectID != id ||
		record.Commit == nil || *record.Commit != main || main == "" {
		t.Fatalf("backup record = %+v, want pushed %s at target main %q", record, id, main)
	}
	if want := filepath.Join(config.LedgerDir(workspace), "backups.jsonl"); recordPath != want {
		t.Fatalf("record_path = %q, want %q", recordPath, want)
	}
	if got := fixture.targetGit("ls-tree", "-r", "--name-only", "main"); got != id+"/backup.json\n"+id+"/journal/00000001.jsonl\n"+id+"/manifest.json" {
		t.Fatalf("target tree = %q, want only the project's folder", got)
	}
	if got := fixture.targetGit("log", "-1", "--format=%an <%ae>|%cn <%ce>|%s", "main"); got != "fid backup <backup@fid.invalid>|fid backup <backup@fid.invalid>|fid backup "+id {
		t.Fatalf("backup commit = %q, want fid's own identity and subject", got)
	}

	var manifest backupFolderManifestForTest
	decoder := json.NewDecoder(strings.NewReader(fixture.targetGit("cat-file", "blob", "main:"+id+"/backup.json")))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		t.Fatal(err)
	}
	archive, _ := exportForTest(t, workspace)
	digest := sha256.Sum256(archive)
	archiveDigest := "sha256:" + hex.EncodeToString(digest[:])
	if manifest.Schema != "fid.backup.v2" || manifest.ProjectID != id ||
		manifest.ArchiveDigest != archiveDigest || manifest.ManifestSchema != "fid.export-manifest.v2" || manifest.EventCount < 1 ||
		manifest.GenesisEventHash != ledgerGenesisHash(t, workspace) {
		t.Fatalf("folder manifest = %+v, want export members with archive digest %s, written by this ledger", manifest, archiveDigest)
	}
	if *record.ArchiveDigest != archiveDigest || *record.EventCount != manifest.EventCount || *record.TerminalEventHash != manifest.TerminalEventHash {
		t.Fatalf("record %+v disagrees with folder manifest %+v", record, manifest)
	}

	// The digest still identifies the deterministic archive export creates.
	archivePath := filepath.Join(t.TempDir(), "export.zip")
	if err := os.WriteFile(archivePath, archive, 0o600); err != nil {
		t.Fatal(err)
	}
	if validated := executeCommandForTest(workspace, "--json", "export", "validate", archivePath); validated.err != nil || !strings.Contains(validated.stdout, `"status":"valid"`) {
		t.Fatalf("export validate = %#v, want a valid archive", validated)
	}

	records := backupRecords(t, workspace)
	if len(records) != 1 || mustJSON(t, records[0]) != mustJSON(t, record) {
		t.Fatalf("backup record file = %+v, want the one reported record %+v", records, record)
	}
	if leftovers, err := filepath.Glob(filepath.Join(fixture.cacheHome, "fid", "backup", "*", "export")); err != nil || len(leftovers) != 0 {
		t.Fatalf("backup left its export in the cache: %v, %v", leftovers, err)
	}
	// fid's clone tracks what it pushed, so the next backup needs no fetch
	// and git never prunes the pushed objects.
	if tracked := fixture.git("--git-dir="+fixture.clone(), "rev-parse", "refs/remotes/origin/main"); tracked != main {
		t.Fatalf("fid's clone tracks %s, want the pushed %s", tracked, main)
	}
}

// clone returns fid's one clone of the target in the fixture's cache.
func (fixture *backupFixture) clone() string {
	fixture.t.Helper()
	clones, err := filepath.Glob(filepath.Join(fixture.cacheHome, "fid", "backup", "*", "repo.git"))
	if err != nil || len(clones) != 1 {
		fixture.t.Fatalf("fid's clones = %v, %v; want one", clones, err)
	}
	return clones[0]
}

func TestBackupTextOutputNamesFolderAndCommit(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, id := fixture.project()
	fixture.setTarget(fixture.target)

	pushed := fixture.backup(workspace)
	if pushed.err != nil || pushed.stderr != "" {
		t.Fatalf("backup = %#v", pushed)
	}
	main := fixture.targetMain()
	if want := "backup pushed: " + id + "/ at commit " + main + "\n"; pushed.stdout != want {
		t.Fatalf("text output = %q, want %q", pushed.stdout, want)
	}
	unchanged := fixture.backup(workspace)
	if want := "backup unchanged: " + id + "/ is current at commit " + main + "\n"; unchanged.err != nil || unchanged.stdout != want {
		t.Fatalf("unchanged text output = %#v, want %q", unchanged, want)
	}
}

func TestBackupWithoutChangesPushesNothingAndRecordsUnchanged(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)

	first, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	second, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	if second.Outcome != "unchanged" || *second.Commit != *first.Commit || *second.ArchiveDigest != *first.ArchiveDigest {
		t.Fatalf("second backup = %+v, want unchanged at %s", second, *first.Commit)
	}
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "1" {
		t.Fatalf("target commits = %s, want 1 after an unchanged backup", count)
	}

	runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Backed up later", "--actor", "agent:developer")
	third, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	if third.Outcome != "pushed" || *third.Commit == *first.Commit || *third.EventCount <= *first.EventCount {
		t.Fatalf("third backup = %+v, want a new push with more events than %+v", third, first)
	}
	if parent := fixture.targetGit("rev-parse", "main^"); parent != *first.Commit {
		t.Fatalf("new backup's parent = %s, want %s", parent, *first.Commit)
	}
	if got := backupOutcomes(backupRecords(t, workspace)); strings.Join(got, ",") != "pushed,unchanged,pushed" {
		t.Fatalf("recorded outcomes = %v", got)
	}
}

func TestBackupKeepsOtherFoldersInTheTarget(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	first, firstID := fixture.project()
	second, secondID := fixture.project()
	fixture.setTarget(fixture.target)
	decodeBackupSuccess(t, fixture.backup(first, "--json"))

	// Someone else adds to the backup repository between fid's backups.
	human := filepath.Join(fixture.base, "human")
	fixture.git("clone", "-q", "-b", "main", fixture.target, human)
	if err := os.WriteFile(filepath.Join(human, "README.md"), []byte("ledgers\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fixture.git("-C", human, "add", "README.md")
	fixture.git("-C", human, "-c", "user.name=human", "-c", "user.email=human@example.test", "commit", "-q", "-m", "readme")
	fixture.git("-C", human, "push", "-q", "origin", "HEAD:main")
	// A branch whose name ends like main's. git ls-remote's pattern matches
	// it too, and it names the commit fid's clone already tracks.
	fixture.git("-C", human, "push", "-q", "origin", "HEAD~1:refs/heads/x/refs/heads/main")

	decodeBackupSuccess(t, fixture.backup(second, "--json"))
	runFIDCommand(t, first, "issue", "create", "--kind", "task", "--title", "Changed again", "--actor", "agent:developer")
	decodeBackupSuccess(t, fixture.backup(first, "--json"))
	firstFiles := []string{firstID + "/backup.json", firstID + "/journal/00000001.jsonl", firstID + "/manifest.json"}
	secondFiles := []string{secondID + "/backup.json", secondID + "/journal/00000001.jsonl", secondID + "/manifest.json"}
	want := strings.Join(append(append([]string{"README.md"}, firstFiles...), secondFiles...), "\n")
	if firstID > secondID {
		want = strings.Join(append(append([]string{"README.md"}, secondFiles...), firstFiles...), "\n")
	}
	if got := fixture.targetGit("ls-tree", "-r", "--name-only", "main"); got != want {
		t.Fatalf("target tree = %q, want %q", got, want)
	}
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "4" {
		t.Fatalf("target commits = %s, want 4", count)
	}
}

func TestBackupJSONFailureIsTypedAndExitsNonZero(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()

	result := fixture.backup(workspace, "--json")
	if cli.ExitCode(result.err) != 2 || result.stdout != "" {
		t.Fatalf("backup without target = %#v, want exit 2", result)
	}
	assertJSONFailureCode(t, result.stderr, "project.backup", "invalid-argument")
}

// A backup killed mid-export leaves its export behind in fid's cache. The
// next backup clears it and carries on.
func TestBackupClearsAnExportLeftByACrashedBackup(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)
	first, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	stale := filepath.Join(filepath.Dir(fixture.clone()), "export", "export.zip")
	writeBackupTestFile(t, stale, "half an archive")

	record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	if record.Outcome != "unchanged" || *record.ArchiveDigest != *first.ArchiveDigest {
		t.Fatalf("backup = %+v, want the same export as before the crash", record)
	}
	if _, err := os.Lstat(filepath.Dir(stale)); !os.IsNotExist(err) {
		t.Fatalf("the crashed backup's export is still in the cache: %v", err)
	}
}
