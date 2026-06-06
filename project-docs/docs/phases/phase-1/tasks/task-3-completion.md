# Phase 1 Task 3 Completion Report

## Summary of what was implemented

- Implemented checkout-scoped `.fid/write.lock` acquisition/release in `internal/journal/lock.go` using atomic lock-directory creation plus owner metadata capture (`pid`, `hostname`, `started_at`, `command`, `actor_role` when supplied).
- Implemented `internal/journal/writer.go` append helpers so one locked write reads the UTC clock once, selects `.fid/journal/YYYY-MM-DD.jsonl`, allocates the next per-day `EVT-<date>-<sequence>` ID by scanning existing lines, appends one canonical JSONL line, and flushes it with `Sync()`.
- Added red/green tests for one-line append, per-day sequence increment, holder-named lock timeout errors, owner metadata, concurrent same-checkout writes, and clock-skew rejection without append.

## Files created or modified

- `internal/journal/lock.go`
- `internal/journal/writer.go`
- `internal/journal/lock_test.go`
- `internal/journal/writer_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-3-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-3-completion.md`

## Validation Contract Coverage

- `VAL-007`:
  `TestAppendEventWritesOneJSONLLineWithFirstSequence` proves one logical append writes exactly one v1 JSONL line with populated envelope fields.
  `TestAppendEventRejectsClockSkewWithoutAppending` proves the writer appends nothing when the chosen UTC date is behind the latest journal date.
- `VAL-013`:
  `TestWithLockWritesOwnerMetadata` proves `.fid/write.lock` records holder metadata and is released after the critical section.
  `TestWithLockTimesOutWithHolderNamedError` proves a waiting caller exits with an error naming the holder and lock path.
  `TestConcurrentWritesAppendDistinctEventsWithMonotonicSequences` proves two same-checkout writers both persist without data loss under the race detector.
- `VAL-015`:
  `TestAppendEventWritesOneJSONLLineWithFirstSequence` proves first-of-day allocation is `EVT-<date>-000001`.
  `TestAppendEventIncrementsSequenceWithinUTCJournalFile` proves the next append for the same UTC journal file allocates `000002`.
  `TestConcurrentWritesAppendDistinctEventsWithMonotonicSequences` proves concurrent same-day writes allocate distinct monotonic sequences under the lock.

## Tests run

- `go test ./internal/journal -run TestAppendEventWritesOneJSONLLineWithFirstSequence` - PASS
- `go test ./internal/journal -run 'TestWithLockWritesOwnerMetadata|TestWithLockTimesOutWithHolderNamedError'` - PASS
- `go test ./internal/journal -run 'TestAppendEventIncrementsSequenceWithinUTCJournalFile|TestAppendEventRejectsClockSkewWithoutAppending|TestConcurrentWritesAppendDistinctEventsWithMonotonicSequences'` - PASS
- `go test -race ./internal/journal` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- Conditional `go test -tags integration -race ./...` - not triggered; Task 3 lock/concurrency tests were kept in the default test suite.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `b0f3653`

## Commit subject

- `feat(journal): add write lock and append writer`
