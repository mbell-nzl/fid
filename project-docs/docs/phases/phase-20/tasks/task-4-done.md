# Phase 20 - Task 4: Rebuild, authority, and lock failures - Completion Report

## Summary of what was implemented

- Completed the inherited Task 4 implementation: `fid rebuild` now follows project selection and the `fid-cli.v1` Response Envelope, maps retained V1 authority failures, returns projection-only success with `event:null`, and reports write-lock contention as typed, retryable `lock-unavailable`.
- Deleted the dead private `describeLockHolder` formatter while retaining `.fid/write.lock/owner.json` metadata writing and `readLockMetadata` coverage.

## Files created or modified

- `cmd/fid/rebuild.go`, `cmd/fid/response_errors.go`, `cmd/fid/rebuild_test.go`, and `cmd/fid/root_flags_test.go`
- `internal/journal/lock.go`, `internal/journal/lock_test.go`, and `internal/journal/projection.go`
- Rebuild channel assertions in `cmd/fid/claim_test.go`, `cmd/fid/issue_ac_test.go`, and the affected `cmd/fid/e2e_v01*.go` tests.
- `project-docs/docs/phases/phase-20/tasks/task-4-done.md`

## Documentation updated

- None required; this required completion report replaces the recovery blocker report.

## Validation Contract Coverage

- VAL-001: PASS - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` asserts the exact success-envelope keys on stdout with empty stderr; the real `/tmp/fid` run emitted `{"contract_version":"fid-cli.v1",...,"event":null}`.
- VAL-002: PASS - `TestRebuildJSONMapsMissingAuthorityWithoutMutation`, `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation`, and `TestRebuildJSONMapsReplayStopWithClosedDetails` assert closed JSON failure envelopes, exact code-specific details, stderr-only output, and non-mutation.
- VAL-006: PASS - `TestRebuildJSONMapsMissingAuthorityWithoutMutation` and `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation` cover missing, unknown, and denied authority before journal/projection work; `go test -race -coverprofile=coverage.out ./...` passed.
- VAL-007: PASS - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` asserts `lock-unavailable`, `conflict`, `retryable:true`, closed `{lock_path}` details, no journal append, then success after release; `TestWithLockReturnsTypedUnavailableErrorForHeldLock` proves the typed journal error.
- VAL-008: PASS - `TestRebuildCommandReplaysJournalAndReportsCounts` asserts Human rebuild success on stdout with empty stderr; the JSON failure tests assert empty stdout and stderr-only failures.
- VAL-017: PASS - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` and `TestRootFlagsAllowRebuildForAnExplicitProject` assert `event == nil` for the projection-only rebuild read; the real CLI output also contained `"event":null`.
- VAL-044 (rebuild row): PASS - `TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork` rejects positional input and `--json --verbose` before mutation, while `TestRootFlagsAllowRebuildForAnExplicitProject` and `TestRebuildJSONMapsExplicitMissingProjectToSelectedPath` cover the closed root request/project behavior.

## Tests run

- `go test ./cmd/fid ./internal/journal -run '^(TestRebuildJSONMapsMissingAuthorityWithoutMutation|TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestRebuildJSONMapsExplicitMissingProjectToSelectedPath|TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease|TestRebuildCommandReplaysJournalAndReportsCounts|TestWithLockWritesOwnerMetadata|TestWithLockReturnsTypedUnavailableErrorForHeldLock)$' -count=1` - PASS - `0.009s`
- `/tmp/fid --project "$task4_workspace" init --prefix FID; /tmp/fid --project "$task4_workspace" --json rebuild --by workflow_runner` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":1,"upgraded_events":0},"event":null}`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `coverage: 80.5% of statements`
- `go tool cover -func=coverage.out` - PASS - `77.8%`
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `4.625s`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

## Deviations from spec

- The first session stopped at `internal/journal/lock.go:185:6: func describeLockHolder is unused (unused)`. In this later recovery session, the dead private formatter was deleted and the exact required gates passed; no contract, behavior, or test assertion was weakened.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
