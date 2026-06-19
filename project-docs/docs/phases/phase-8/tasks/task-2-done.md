## Summary of what was implemented

Implemented Phase 8 Task 2 gate snapshotting on `fid issue complete`.

- `internal/issues.Complete` now snapshots matching `gate_runs` rows for the derived `(issue_id, attempt_number)` into `task_completions.gate_results`, derives `gates_passed` from `exit_code == 0`, and derives `gates_failed` from non-zero exits.
- The snapshot is a pure read of `gate_runs`; source `gate_runs` rows remain unchanged.
- The issue detail read path now exposes structured `gate_results` objects plus `gates_passed` / `gates_failed`, while remaining tolerant of older string-array fixtures already present in tests.
- Added journal-fixture-seeded Task 2 tests covering pass+fail snapshotting, attempt filtering, preserved source rows, and the no-gate `[]` / `0` / `0` case.

## Files created or modified

- `internal/issues/complete.go`
- `internal/issues/detail.go`
- `internal/issues/complete_test.go`
- `project-docs/docs/phases/phase-8/tasks/task-2-done.md`

## Documentation updated

- Verified `project-docs/specs/03_TRD_SCHEMA.md` gate snapshot comment and `project-docs/specs/06_DESIGN_PLAN_CLI.md` snapshot note; no spec text change was required.
- Added this completion report at `project-docs/docs/phases/phase-8/tasks/task-2-done.md`.

## Validation Contract Coverage

- `VAL-005` PASS.
  Evidence:
  `TestCompleteSnapshotsMatchingGateRunsAndPreservesSourceRows`
  `TestCompleteLeavesGateSnapshotEmptyWhenNoMatchingGateRunsExist`
  Focused red/green evidence used journal-fixture-seeded `gate_runs` because `fid dev gate` remains v0.1e per DEC-C.
  The matching-attempt snapshot asserted this stored/read-back shape:

```json
[
  {
    "cmd": "go test ./...",
    "exit_code": 0,
    "ran_at": "2026-06-20T09:00:00Z",
    "output_excerpt": "ok\tgithub.com/mbell-nzl/fid/internal/issues\t0.321s"
  },
  {
    "cmd": "golangci-lint run ./...",
    "exit_code": 1,
    "ran_at": "2026-06-20T09:01:00Z",
    "output_excerpt": "internal/issues/complete.go:128: line is 61 chars"
  }
]
```

  Derived counts asserted:
  `gates_passed = 1`
  `gates_failed = 1`
  The same test also proved that attempt-2 fixture row `GATE-000103` was excluded from attempt 1 and that all three `gate_runs` rows remained present after completion.
  The no-gate test proved `gate_results = []` and `gates_passed = gates_failed = 0`.

## Tests run

- `go test ./internal/issues -run TestCompleteSnapshotsMatchingGateRunsAndPreservesSourceRows`
  First red run failed at build time because `IssueCompletion` did not yet expose `gates_passed` / `gates_failed`.
- `go test ./internal/issues -run 'TestComplete(SnapshotsMatchingGateRunsAndPreservesSourceRows|LeavesGateSnapshotEmptyWhenNoMatchingGateRunsExist)$'`
- `go test ./internal/issues`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
  Coverage highlights:
  `cmd/fid` 78.7%
  `internal/issues` 81.9%
  `internal/journal` 80.3%
  `internal/policy` 85.0%
  `internal/workflow` 95.5%
- `go tool cover -func=coverage.out` PASS
  Overall statement coverage: 77.0%
- `go test -tags integration -race ./...` NOT RUN
  Reason:
  Task 2 snapshot tests remained in the default test build; no `integration` tag was introduced.
- `go build -o /tmp/fid ./cmd/fid` NOT RUN
  Reason:
  no CLI transcript was needed for this task because VAL-005 was satisfied through deterministic `internal/issues` tests.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Recorded after commit creation in the terminal handoff. The tracked report intentionally does not embed the final hash because doing so would make the committed file self-referential and require a forbidden amend.

## Commit subject

- Recorded after commit creation in the terminal handoff.
