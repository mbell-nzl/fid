package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func renderJournalReplayError(writeLine func(string, ...any) error, err error, stopMessage string) (bool, error) {
	var corruptErr *journal.CorruptJournalError
	if errors.As(err, &corruptErr) {
		if writeErr := writeLine("error: malformed journal entry at %s:%d\n", journalErrorPath(corruptErr.File), corruptErr.Line); writeErr != nil {
			return true, writeErr
		}
		if writeErr := writeLine("  parse error: %v\n", corruptErr.Err); writeErr != nil {
			return true, writeErr
		}
		return true, writeLine("  %s\n", stopMessage)
	}

	var conflictErr *journal.PreconditionConflictError
	if errors.As(err, &conflictErr) {
		if writeErr := writeLine("error: replay precondition conflict at %s:%d\n", journalErrorPath(conflictErr.File), conflictErr.Line); writeErr != nil {
			return true, writeErr
		}
		if writeErr := writeLine("  event: %s (%s)\n", conflictErr.EventID, conflictErr.EventType); writeErr != nil {
			return true, writeErr
		}
		if writeErr := writeLine("  entity: %s pk=%s\n", conflictErr.Table, formatJournalMap(conflictErr.PK)); writeErr != nil {
			return true, writeErr
		}
		if writeErr := writeLine("  expected: %s\n", formatJournalMap(conflictErr.Expected)); writeErr != nil {
			return true, writeErr
		}
		if writeErr := writeLine("  actual: %s\n", formatJournalMap(conflictErr.Actual)); writeErr != nil {
			return true, writeErr
		}
		return true, writeLine("  %s\n", stopMessage)
	}

	return false, nil
}

func journalErrorPath(file string) string {
	if strings.HasPrefix(file, ".fid/") || strings.HasPrefix(file, "/") {
		return file
	}
	return ".fid/journal/" + file
}

func formatJournalMap(values map[string]any) string {
	if len(values) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s:%v", key, values[key]))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
