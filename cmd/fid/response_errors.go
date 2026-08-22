package main

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func authorityFailure(operation, actor string, err error) (cli.Failure, bool, error) {
	actor = strings.TrimSpace(actor)
	details := map[string]any{"operation": operation, "actor": nil}

	var missingErr *policy.MissingRoleError
	if errors.As(err, &missingErr) {
		failure, failureErr := cli.NewFailure(cli.CodeAuthorityMissing, "acting authority is required", nil, details)
		return failure, true, failureErr
	}

	if actor != "" {
		details["actor"] = actor
	}
	var unknownErr *policy.UnknownRoleError
	if errors.As(err, &unknownErr) {
		failure, failureErr := cli.NewFailure(cli.CodeAuthorityUnknownActor, "acting authority is unknown", nil, details)
		return failure, true, failureErr
	}

	var deniedErr *policy.PermissionDeniedError
	if errors.As(err, &deniedErr) {
		failure, failureErr := cli.NewFailure(cli.CodeAuthorityDenied, "acting authority is denied", nil, details)
		return failure, true, failureErr
	}

	return cli.Failure{}, false, nil
}

func lockUnavailableFailure(err error) (cli.Failure, bool, error) {
	var unavailableErr *journal.LockUnavailableError
	if !errors.As(err, &unavailableErr) {
		return cli.Failure{}, false, nil
	}
	failure, failureErr := cli.NewFailure(cli.CodeLockUnavailable, "project write lock is unavailable", nil, map[string]any{
		"lock_path": filepath.ToSlash(filepath.Join(".fid", filepath.Base(unavailableErr.LockPath))),
	})
	return failure, true, failureErr
}

func journalReplayStoppedFailure(err error) (cli.Failure, bool, error) {
	var corruptErr *journal.CorruptJournalError
	if errors.As(err, &corruptErr) {
		failure, failureErr := cli.NewFailure(cli.CodeJournalReplayStopped, "journal replay stopped", nil, map[string]any{
			"journal_file": journalErrorPath(corruptErr.File),
			"line":         corruptErr.Line,
			"event_id":     nil,
		})
		return failure, true, failureErr
	}

	var conflictErr *journal.PreconditionConflictError
	if errors.As(err, &conflictErr) {
		failure, failureErr := cli.NewFailure(cli.CodeJournalReplayStopped, "journal replay stopped", nil, map[string]any{
			"journal_file": journalErrorPath(conflictErr.File),
			"line":         conflictErr.Line,
			"event_id":     conflictErr.EventID,
		})
		return failure, true, failureErr
	}

	return cli.Failure{}, false, nil
}
