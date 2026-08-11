# Phase 20 - Task 4: Rebuild, authority, and lock failures - Completion Report

## Summary of what was implemented

- The committed Task 4 implementation migrates `fid rebuild` to project selection and the `fid-cli.v1` Response Envelope, maps retained V1 authority failures, returns projection-only success with `event:null`, and maps write-lock contention to retryable `lock-unavailable`.
- This documentation-only repair recaptured Task 4 evidence on the implementation tree. No product code, tests, migrations, or behavior changed in this session.

## Files created or modified

- `project-docs/docs/phases/phase-20/tasks/task-4-done.md` — rewritten with fresh DEC-136 evidence.
- The implementation commit `e56b3a0ce0433f1044f70b7969c14e19234dfdf8` modified `cmd/fid/rebuild.go`, `cmd/fid/response_errors.go`, `cmd/fid/rebuild_test.go`, `cmd/fid/root_flags_test.go`, `internal/journal/lock.go`, `internal/journal/lock_test.go`, and `internal/journal/projection.go`.

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-001: PASS - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` at `cmd/fid/rebuild_test.go:264` asserts JSON success is stdout-only; this session's real CLI entrypoint returned `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":1,"upgraded_events":0},"event":null}`.
- VAL-002: PASS - `TestRebuildJSONMapsMissingAuthorityWithoutMutation` (`cmd/fid/rebuild_test.go:17`), `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation` (`cmd/fid/rebuild_test.go:86`), and `TestRebuildJSONMapsReplayStopWithClosedDetails` (`cmd/fid/rebuild_test.go:352`) cover the closed failure envelopes, code-specific details, stderr-only failures, and no mutation; the fresh full race suite passed.
- VAL-006: PASS - `TestRebuildJSONMapsMissingAuthorityWithoutMutation` (`cmd/fid/rebuild_test.go:17`) and `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation` (`cmd/fid/rebuild_test.go:86`) cover missing, unknown, and denied authority before rebuild work; the fresh prescribed focused command passed for both command and journal packages.
- VAL-007: PASS - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` (`cmd/fid/rebuild_test.go:264`) proves `lock-unavailable`, `retryable:true`, and success after deterministic release; `TestWithLockReturnsTypedUnavailableErrorForHeldLock` (`internal/journal/lock_test.go:96`) proves the typed lock error. Both were included in this session's prescribed focused command.
- VAL-008: PASS - `TestRebuildCommandReplaysJournalAndReportsCounts` (`cmd/fid/rebuild_test.go:397`) covers Human success on stdout with empty stderr, while the JSON authority/lock failure tests at `cmd/fid/rebuild_test.go:17`, `:86`, and `:264` cover empty stdout and stderr-only failures; the fresh full race suite passed.
- VAL-017: PASS - `TestRootFlagsAllowRebuildForAnExplicitProject` (`cmd/fid/root_flags_test.go:79`) and `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` (`cmd/fid/rebuild_test.go:264`) assert a null Event for the projection-only read; the real CLI entrypoint output above contains `"event":null`.
- VAL-044 (rebuild row): PASS - `TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork` (`cmd/fid/rebuild_test.go:153`) rejects positional input and JSON-plus-verbose before work; `TestRootFlagsAllowRebuildForAnExplicitProject` (`cmd/fid/root_flags_test.go:79`) and `TestRebuildJSONMapsExplicitMissingProjectToSelectedPath` (`cmd/fid/rebuild_test.go:221`) cover accepted and missing explicit-project requests. The focused command passed in this session.

## Tests run

- `go test ./cmd/fid ./internal/journal -run '^(TestRebuildJSONMapsMissingAuthorityWithoutMutation|TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestRebuildJSONMapsExplicitMissingProjectToSelectedPath|TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease|TestRebuildCommandReplaysJournalAndReportsCounts|TestWithLockWritesOwnerMetadata|TestWithLockReturnsTypedUnavailableErrorForHeldLock)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal  0.005s`
- `/tmp/fid --project /tmp/fid-task4-S14hTp init --prefix FID` - PASS - `initialized fid project in .fid/`
- `/tmp/fid --project /tmp/fid-task4-S14hTp --json rebuild --by workflow_runner` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":1,"upgraded_events":0},"event":null}`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  8.316s  coverage: 80.5% of statements`
- `go tool cover -func=coverage.out` - PASS - `total:                    (statements)          77.8%`
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  5.215s`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

## Deviations from spec

- The first implementation session stopped at `internal/journal/lock.go:185:6: func describeLockHolder is unused (unused)`. The committed Task 4 implementation removed that dead private formatter; no contract, behavior, or test assertion was weakened.
- DEC-136 supersedes the report tokens `0.009s`, `coverage: 80.5% of statements`, and `4.625s`; all three were replaced by fresh captures from this session against the same Task 4 implementation tree, rather than expanded or reused from earlier evidence.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e56b3a0ce0433f1044f70b7969c14e19234dfdf8`

## Commit subject

- `feat(rebuild): map typed authority and lock failures`
