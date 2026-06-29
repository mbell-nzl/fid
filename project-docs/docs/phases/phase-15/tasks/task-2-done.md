## Summary of what was implemented

Implemented the new `internal/dev` foundation for Phase 15 Task 2 and shipped
`fid dev start` plus `fid dev gate` as single-event composite writes.

- `internal/dev.Start` now authorizes `dev.start`, validates before append,
  acquires the caller claim in the same event, transitions `todo ->
  in_progress`, records the `status_changed` activity effect, and opens one
  `agent_runs` row with `ended_at = NULL`.
- `internal/dev.Gate` now authorizes `dev.gate`, records reported gate results
  without executing the command, derives the active attempt from prior
  completions, links to the open run when present, and appends one
  `gate_runs` row in one `dev.gate.recorded` event.
- Registered the new `fid dev` command tree from `cmd/fid/root.go` and added
  the required `--model`, `--exit-code`, and `--output-excerpt` CLI surface.
- Added the four `dev.*` permission rows now (`dev.start`, `dev.gate`,
  `dev.finish`, `dev.block`) and covered the `start`/`gate` show-vs-enforce
  parity in tests.

## Files created or modified

- `cmd/fid/dev.go`
- `cmd/fid/dev_test.go`
- `cmd/fid/root.go`
- `internal/dev/dev.go`
- `internal/dev/start.go`
- `internal/dev/gate.go`
- `internal/dev/start_test.go`
- `internal/dev/gate_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-15/tasks/task-2-done.md`

## Documentation updated

- `project-docs/specs/02_TRD.md`
  - Added the `internal/dev` domain boundary to the architecture map and module
    boundary table.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  - Added required `fid dev start --model <model>`.
  - Added required `fid dev gate --exit-code <n>` and optional
    `--output-excerpt <text>`.
  - Clarified that `fid dev gate` records, and does not execute, the reported
    command.
- `project-docs/docs/phases/phase-15/tasks/task-2-done.md`
  - Recorded implementation evidence, gate results, and coverage notes.

## Validation Contract Coverage

- `VAL-003` PASS
  - `internal/dev/start_test.go` proves one `dev.started` event with claim
    insert, issue transition, activity effect, and open run creation.
  - `cmd/fid/dev_test.go` proves the CLI confirmation and projected
    `in_progress` state.
- `VAL-004` PASS
  - `internal/dev/start_test.go` proves reject-before-append for a conflicting
    live claim and a non-`todo` issue, with unchanged journal line counts.
- `VAL-005` PASS
  - `internal/dev/gate_test.go` proves one `dev.gate.recorded` event, active
    attempt derivation, open-run linking, and reject-before-append for missing
    started state and missing live claim.
  - `cmd/fid/dev_test.go` proves the record-only CLI behavior by asserting the
    reported command string is stored and not executed.
- `VAL-006` PARTIAL PASS FOR TASK 2
  - `internal/policy/matrix_test.go` proves the four `dev.*` rows match section
    13.
  - `cmd/fid/dev_test.go` proves `fid policy show` includes those rows and that
    `dev.start`/`dev.gate` allow `developer` and deny `qa_reviewer` before
    append.
  - `dev.finish`/`dev.block` enforcement remains for Task 3, per the task spec.

## Tests run

- Red evidence:
  - `go test ./internal/dev`
    - FAIL expected: `undefined: Start` / `undefined: StartOptions`
  - `go test ./internal/policy -run 'TestMatrixContainsExactlyTheImplementedV01WriteActions|TestMatrixDevRowsMatchSchemaSection13'`
    - FAIL expected: missing `dev.*` rows
  - `go test ./cmd/fid -run TestDev`
    - FAIL expected: `unknown command "dev" for "fid"`
- Focused green checks:
  - `go test ./internal/dev -run TestStart`
  - `go test ./internal/dev -run TestGate`
  - `go test ./internal/dev`
  - `go test ./internal/policy -run 'TestMatrixContainsExactlyTheImplementedV01WriteActions|TestMatrixDevRowsMatchSchemaSection13'`
  - `go test ./cmd/fid -run TestDev`
- Size check:
  - `wc -l internal/dev/dev.go internal/dev/start.go internal/dev/gate.go internal/dev/start_test.go internal/dev/gate_test.go cmd/fid/dev.go cmd/fid/dev_test.go`
  - Result: `internal/dev/dev.go 304`, `internal/dev/start.go 201`,
    `internal/dev/gate.go 181`, `cmd/fid/dev.go 146`, `cmd/fid/dev_test.go 215`

## Exact quality-gate results

- `gofmt -l cmd internal`
  - PASS
  - Output: empty
- `golangci-lint run ./...`
  - PASS
  - Output excerpt: `0 issues.`
- `go test -race -coverprofile=coverage.out ./...`
  - PASS
  - Output excerpts:
    - `cmd/fid` coverage `77.3%`
    - `internal/dev` coverage `74.2%`
    - `internal/policy` coverage `84.6%`
- `go tool cover -func=coverage.out`
  - PASS
  - Output excerpt: total coverage `75.8%`
- `go test -tags integration -race ./...`
  - PASS
  - Output excerpt: integration suite green for `cmd/fid`, `internal/dev`,
    `internal/issues`, `internal/policy`, and the remaining packages.
- `go build -o /tmp/fid ./cmd/fid`
  - PASS
  - Output: empty

Package coverage percentages required by Task 2:

- `internal/dev`: `74.2%`
- `internal/policy`: `84.6%`
- `cmd/fid`: `77.3%`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `07ac0ae`

## Commit subject

- `feat(dev): add dev start and gate composites`
