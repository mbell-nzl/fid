## Summary of what was implemented

Implemented Phase 8 Task 3 as a single end-to-end integration test over the
real completion write path.

- Added `cmd/fid/e2e_v01c_complete_test.go` with one integration test that
  proves `task-default` blocks `in_progress -> dev_complete` on
  `has_completion_report` until a real `fid issue complete` command succeeds.
- The same test proves validate-before-append on the reject path by asserting
  unchanged journal line count and unchanged issue JSON.
- The same test then seeds attempt-2 `gate_runs`, records a second real
  completion with human intervention, deletes the SQLite DB, runs `fid rebuild`,
  and asserts identical completions JSON plus identical `gate_runs` rows before
  and after rebuild.

## Files created or modified

- `cmd/fid/e2e_v01c_complete_test.go`
- `project-docs/docs/phases/phase-8/tasks/task-3-done.md`

## Documentation updated

- Added this completion report at
  `project-docs/docs/phases/phase-8/tasks/task-3-done.md`.
- Recorded that the Phase 5 fixture-based `has_completion_report` gating
  evidence is now superseded by a real `fid issue complete` command flow.

## Validation Contract Coverage

- `VAL-006` PASS.
  Evidence:
  `TestV01CHasCompletionReportGateAndRebuildParity`
  The test creates a zero-AC task, transitions it to `in_progress`, rejects
  `dev_complete` before any completion, asserts stderr contains
  `has_completion_report` plus the issue ID, asserts the journal line count is
  unchanged, asserts the issue JSON is unchanged, then runs a real
  `fid issue complete` command and proves the same transition succeeds.
- `VAL-007` PASS.
  Evidence:
  `TestV01CHasCompletionReportGateAndRebuildParity`
  The test seeds attempt-2 `gate_runs`, records a second real completion with
  `--human-intervention --intervention-note`, snapshots completions JSON and
  ordered `gate_runs` rows, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and
  asserts exact parity after rebuild.

Observed command transcript excerpt from the integration test:

```text
$ fid issue transition <task> --to in_progress --by developer
transitioned <task>: todo → in_progress

$ fid issue transition <task> --to dev_complete --by developer
error: transition issue "<task>" from "in_progress" to "dev_complete" blocked by workflow rules:
  - has_completion_report: entities <task>
(journal lines unchanged; issue JSON unchanged)

$ fid issue complete <task> --by developer --notes-md project-docs/docs/phases/phase-8/tasks/task-3-done.md --tokens 1600
completed <task> (attempt 1)

$ fid issue transition <task> --to dev_complete --by developer
transitioned <task>: in_progress → dev_complete

$ fid issue complete <task> --by developer --notes-md project-docs/docs/phases/phase-8/tasks/task-3-done.md --tokens 1700 --human-intervention --intervention-note "paired with qa on rebuild parity"
completed <task> (attempt 2)

$ fid rebuild --by developer
rebuild complete
```

Rebuild parity snapshot covered:

- completion attempt ordering `2, 1`
- latest status `dev_complete`, prior status `superseded`
- latest `human_intervention = true`
- latest `human_intervention_note = "paired with qa on rebuild parity"`
- latest `gates_passed = 1`, `gates_failed = 1`
- latest `gate_results` entries for `go test -race ./...` and
  `golangci-lint run ./...`
- ordered `gate_runs` rows for attempt 2 before and after rebuild

## Tests run

- `go test -tags integration ./cmd/fid -run TestV01CHasCompletionReportGateAndRebuildParity -count=1`
  First red run failed at build time because the new test helper name
  `stringPointer` collided with an existing helper in `ready_test.go`.
- `go test -tags integration ./cmd/fid -run TestV01CHasCompletionReportGateAndRebuildParity -count=1`
  PASS after renaming the local helper to `task3StringPointer`.
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go test -tags integration -race ./...`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
  Coverage highlights:
  `cmd/fid` 78.7%
  `internal/issues` 81.9%
  `internal/workflow` 95.5%
- `go tool cover -func=coverage.out` PASS
  Overall statement coverage: 77.0%
- `go test -tags integration -race ./...` PASS
  Reason run:
  required because `cmd/fid/e2e_v01c_complete_test.go` uses the `integration`
  build tag.
- `go build -o /tmp/fid ./cmd/fid` PASS
  Reason run:
  used to support transcriptable CLI evidence for the report.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d0bbb83`

## Commit subject

- `test: add completion gate rebuild parity e2e`
