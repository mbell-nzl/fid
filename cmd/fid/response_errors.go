package main

import (
	"errors"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
)

func lockUnavailableFailure(err error) (cli.Failure, bool, error) {
	var unavailableErr *journal.LockUnavailableError
	if !errors.As(err, &unavailableErr) {
		return cli.Failure{}, false, nil
	}
	failure, failureErr := cli.NewFailure(cli.CodeLockUnavailable, "project write lock is unavailable", nil, map[string]any{
		"lock_path": filepath.ToSlash(unavailableErr.LockPath),
	})
	return failure, true, failureErr
}

func journalReplayStoppedFailure(err error) (cli.Failure, bool, error) {
	var corruptErr *journal.CorruptJournalError
	if errors.As(err, &corruptErr) {
		failure, failureErr := cli.NewFailure(cli.CodeJournalReplayStopped, "journal replay stopped", nil, map[string]any{
			"journal_file": corruptErr.File,
			"line":         corruptErr.Line,
			"event_id":     nil,
		})
		return failure, true, failureErr
	}

	var conflictErr *journal.PreconditionConflictError
	if errors.As(err, &conflictErr) {
		failure, failureErr := cli.NewFailure(cli.CodeJournalReplayStopped, "journal replay stopped", nil, map[string]any{
			"journal_file": conflictErr.File,
			"line":         conflictErr.Line,
			"event_id":     conflictErr.EventID,
		})
		return failure, true, failureErr
	}

	return cli.Failure{}, false, nil
}
