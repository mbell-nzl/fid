// Package backup copies a project's ledger off the machine: it exports the
// ledger into the project's folder in the backup repository, through fid's
// own clone of that repository, commits and pushes. It records every
// attempt, successful or not, beside the ledger. It never writes the journal.
package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

const (
	// RecordFileName is the record of backup attempts, in the ledger.
	RecordFileName = "backups.jsonl"
	// RecordSchema versions each line of the record.
	RecordSchema = "fid.backup-record.v1"
	// DefaultLockWait bounds how long a backup waits for another to finish.
	DefaultLockWait = 2 * time.Minute
)

// Outcome is a backup attempt's result.
type Outcome string

const (
	OutcomePushed    Outcome = "pushed"
	OutcomeUnchanged Outcome = "unchanged"
	OutcomeFailed    Outcome = "failed"
)

// Kind names why a backup failed.
type Kind string

const (
	KindProjectIdentity   Kind = "project-identity"
	KindUserConfig        Kind = "user-config"
	KindNoTarget          Kind = "no-target"
	KindLockUnavailable   Kind = "lock-unavailable"
	KindTargetUnreachable Kind = "target-unreachable"
	KindPushRejected      Kind = "push-rejected"
	KindDifferentLedger   Kind = "different-ledger"
	KindExportFailed      Kind = "export-failed"
	KindInternal          Kind = "internal"
)

// Options identifies the project and the environment a backup runs in.
type Options struct {
	// WorkspaceRoot is the project root, as project.Resolve returns it.
	WorkspaceRoot string
	// Environment is the process environment. It locates user config and
	// fid's cache, and git runs with it, less variables naming another
	// repository, so the user's own git credentials and transport apply.
	Environment   []string
	LockWait      time.Duration
	MigrationsDir string
	Now           func() time.Time
}

// Record is one line of the record, and the result fid backup reports.
type Record struct {
	Schema            string           `json:"schema"`
	StartedAt         string           `json:"started_at"`
	FinishedAt        string           `json:"finished_at"`
	Outcome           Outcome          `json:"outcome"`
	ProjectID         *string          `json:"project_id"`
	Commit            *string          `json:"commit"`
	ArchiveDigest     *string          `json:"archive_digest"`
	EventCount        *int             `json:"event_count"`
	TerminalEventHash *string          `json:"terminal_event_hash"`
	Failure           *RecordedFailure `json:"failure"`
}

// RecordedFailure is why a recorded attempt failed.
type RecordedFailure struct {
	Kind    Kind   `json:"kind"`
	Message string `json:"message"`
}

// Error is a failed backup. Path names the config file at fault for the
// identity, user-config and no-target kinds, and the lock for
// lock-unavailable. Detail is git's own account of a git failure.
type Error struct {
	Kind    Kind
	Message string
	Detail  string
	Path    string
	Err     error
}

func (err *Error) Error() string {
	if err.Detail != "" {
		return err.Message + ": " + err.Detail
	}
	return err.Message
}

func (err *Error) Unwrap() error { return err.Err }

func failure(kind Kind, cause error, format string, arguments ...any) *Error {
	detail := ""
	if cause != nil {
		detail = cause.Error()
	}
	return &Error{Kind: kind, Message: fmt.Sprintf(format, arguments...), Detail: detail, Err: cause}
}

// Run backs up the project and appends the attempt to its record, whatever
// the outcome, while it still holds fid's clone. It returns the record and
// the record file's path. A failed backup returns an *Error.
func Run(options Options) (Record, string, error) {
	if options.Now == nil {
		options.Now = time.Now
	}
	recordPath := filepath.Join(config.LedgerDir(options.WorkspaceRoot), RecordFileName)
	record := Record{Schema: RecordSchema, StartedAt: stamp(options.Now())}
	run := &attempt{options: options, getenv: lookup(options.Environment), record: &record}
	unlock, err := run.prepare()
	if err == nil {
		err = run.push()
	}
	err = finish(recordPath, &record, options.Now(), err)
	if unlock != nil {
		unlock()
	}
	return record, recordPath, err
}

// finish completes the record and appends it. A record that can't be
// written is reported in place of the outcome it would have recorded.
func finish(recordPath string, record *Record, now time.Time, runErr error) error {
	record.FinishedAt = stamp(now)
	var failed *Error
	if runErr != nil {
		if !errors.As(runErr, &failed) {
			failed = failure(KindInternal, runErr, "backup failed")
		}
		*record = Record{Schema: record.Schema, StartedAt: record.StartedAt, FinishedAt: record.FinishedAt, Outcome: OutcomeFailed,
			ProjectID: record.ProjectID, Failure: &RecordedFailure{Kind: failed.Kind, Message: bounded(failed.Error())}}
	}
	if err := appendRecord(recordPath, *record); err != nil {
		if failed != nil {
			failed.Message += fmt.Sprintf(" (and recording the failure failed: %v)", err)
			return failed
		}
		return failure(KindInternal, err, "backup %s, but recording it failed", record.Outcome)
	}
	if failed != nil {
		return failed
	}
	return nil
}

// attempt is one backup, from reading its configuration to its push.
type attempt struct {
	options   Options
	getenv    func(string) string
	record    *Record
	projectID string
	target    string
	directory string
	genesis   string
}

// prepare checks the project's identity and the configured target, then
// waits for fid's clone of the target. It returns the unlock function once
// the clone is held.
func (run *attempt) prepare() (func(), error) {
	root := run.options.WorkspaceRoot
	identityPath := filepath.Join(root, config.ProjectIdentityFileName)
	id, err := config.ReadProjectIdentity(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, identityFailure(identityPath, nil, "%s does not exist; it holds the project ID that names the backup folder", identityPath)
	} else if err != nil {
		return nil, identityFailure(identityPath, err, "invalid project identity in %s", identityPath)
	}
	ledger, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(root), config.ConfigFileName))
	if err != nil {
		return nil, failure(KindInternal, err, "read the ledger config")
	}
	if ledger.ProjectID != id {
		return nil, identityFailure(identityPath, nil, "the project ID in %s is %s, but the ledger's is %s", identityPath, id, ledger.ProjectID)
	}
	run.projectID = id
	run.record.ProjectID = &id

	configPath, err := config.UserConfigPathFrom(run.getenv)
	if err != nil {
		return nil, failure(KindInternal, err, "find the user config")
	}
	user, err := config.ReadUserConfig(configPath)
	if err != nil {
		return nil, &Error{Kind: KindUserConfig, Message: "invalid user config", Detail: err.Error(), Path: configPath, Err: err}
	}
	if strings.TrimSpace(user.BackupTarget) == "" {
		return nil, &Error{Kind: KindNoTarget, Message: "no backup target is set: add [backup] target to " + configPath, Path: configPath}
	}
	run.target = user.BackupTarget

	cache, err := config.UserCacheDirFrom(run.getenv)
	if err != nil {
		return nil, failure(KindInternal, err, "find fid's cache")
	}
	run.directory = CloneDirectory(cache, run.target)
	if err := os.MkdirAll(run.directory, 0o700); err != nil {
		return nil, failure(KindInternal, err, "create fid's backup cache")
	}
	return lock(filepath.Join(run.directory, "lock"), run.options.LockWait, run.options.Now)
}

func identityFailure(path string, cause error, format string, arguments ...any) *Error {
	failed := failure(KindProjectIdentity, cause, format, arguments...)
	failed.Path = path
	return failed
}

// CloneDirectory is where fid keeps its clone of target, and the lock that
// serializes backups through it: a directory under fid's cache named by the
// first 128 bits of the target's SHA-256, so the path never holds the
// target, nor any credential in it.
func CloneDirectory(cache, target string) string {
	sum := sha256.Sum256([]byte(target))
	return filepath.Join(cache, "backup", hex.EncodeToString(sum[:16]))
}

func lookup(environment []string) func(string) string {
	values := map[string]string{}
	for _, entry := range environment {
		if name, value, ok := strings.Cut(entry, "="); ok {
			values[name] = value
		}
	}
	return func(name string) string { return values[name] }
}

func stamp(at time.Time) string { return at.UTC().Format(time.RFC3339Nano) }

// bounded keeps a recorded message to at most 512 characters.
func bounded(text string) string {
	runes := []rune(text)
	if len(runes) <= 512 {
		return text
	}
	return string(runes[:511]) + "…"
}
