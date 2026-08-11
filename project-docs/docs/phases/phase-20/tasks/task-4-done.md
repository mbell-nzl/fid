# Phase 20 - Task 4: Rebuild, authority, and lock failures - Blocker Report

## Summary of what was implemented

- Implemented the in-progress Task 4 changes: rebuild root-flag/project selection, Response Envelope rendering, typed V1 authority mappings, projection-only JSON success, closed rebuild-request checks, and typed retryable rebuild lock handling.
- The required lint gate failed before task completion; per the task protocol, this report records the blocker and the task is not claimed complete.

## Files created or modified

- `cmd/fid/rebuild.go`
- `cmd/fid/response_errors.go`
- `cmd/fid/rebuild_test.go`
- `cmd/fid/root_flags_test.go`
- `internal/journal/lock.go`
- `internal/journal/lock_test.go`
- `internal/journal/projection.go`
- Direct rebuild channel assertions in `cmd/fid/claim_test.go`, `cmd/fid/issue_ac_test.go`, and the affected `cmd/fid/e2e_v01*.go` tests.
- `project-docs/docs/phases/phase-20/tasks/task-4-done.md`

## Documentation updated

- None required; this blocker report is the required task artifact.

## Validation Contract Coverage

- VAL-001: BLOCKED - focused rebuild JSON success/rejection tests passed, but the required lint gate failed.
- VAL-002: BLOCKED - `TestRebuildJSONMapsMissingAuthorityWithoutMutation`, `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation`, and `TestRebuildJSONMapsReplayStopWithClosedDetails` passed in the focused package run; task closure is blocked by lint.
- VAL-006: BLOCKED - the authority test matrix passed in `cmd/fid`; task closure is blocked by lint.
- VAL-007: BLOCKED - `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease` and `TestWithLockReturnsTypedUnavailableErrorForHeldLock` passed in focused runs; task closure is blocked by lint.
- VAL-008: BLOCKED - direct rebuild Human-channel assertions were updated and passed in `go test ./cmd/fid ./internal/journal -count=1`; task closure is blocked by lint.
- VAL-017: BLOCKED - rebuild JSON success asserts `event:null` in `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease`; task closure is blocked by lint.
- VAL-044 (rebuild row): BLOCKED - `TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork` and `TestRebuildJSONMapsExplicitMissingProjectToSelectedPath` passed in focused runs; task closure is blocked by lint.

## Tests run

- `go test ./cmd/fid --run '^TestRebuildJSONMapsMissingAuthorityWithoutMutation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.073s`
- `go test ./internal/journal --run '^TestWithLockReturnsTypedUnavailableErrorForHeldLock$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.004s`
- `go test ./cmd/fid --run '^TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.094s`
- `go test ./cmd/fid --run '^TestRebuildJSONMapsReplayStopWithClosedDetails$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.085s`
- `go test ./cmd/fid --run '^TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.073s`
- `go test ./cmd/fid --run '^TestRebuildJSONMapsExplicitMissingProjectToSelectedPath$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.041s`
- `go test ./cmd/fid ./internal/journal -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	1.740s`
- `go test ./cmd/fid ./internal/journal -count=1 && git diff --check && git diff --stat && git status --short` - PASS - `?? cmd/fid/response_errors.go`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - BLOCKED - `* unused: 1`
- `go test -race -coverprofile=coverage.out ./...` - not run - blocked by the failed required lint gate.
- `go tool cover -func=coverage.out` - not run - blocked because its required coverage producer did not run.
- `go test -tags integration -race ./cmd/fid/integration/...` - not run - blocked by the failed required lint gate.
- `go build -o /tmp/fid ./cmd/fid` - not run - blocked by the failed required lint gate.

## Deviations from spec

- None. The task is incomplete because the required lint gate failed.

## Concerns, technical debt, or notes for the Tech Lead

- Required lint failure: `internal/journal/lock.go:185:6: func describeLockHolder is unused (unused)`. The typed lock-unavailable path no longer reads or formats holder metadata. The task protocol prohibits repairing and rerunning a failed required gate in this session; a follow-up must resolve this before Task 4 can close.
- `project-docs/docs/phases/phase-20/run-ledger.md` was already modified when this session began and was not edited.

## Dependency-budget note

- No dependency changes.

## Commit hash

- No commit created - blocker report.

## Commit subject

- No commit created - blocker report.
