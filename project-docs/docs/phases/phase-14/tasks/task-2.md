# Phase 14 - Task 2: `fid run record` write command

## Objective

Implement `fid run record`: append one `agent_runs` row capturing role, model,
tokens in/out, start/end timestamps, snapshot pricing, and computed `cost_usd`,
plus optional provider reconciliation fields. Warn (without rejecting) when the
model is absent from pricing config. Add the `run.record` permission row to the
in-code matrix. This is the first telemetry write surface and the foundation the
`fid stats` reads (Task 3) and the Phase 15 `fid dev *` composites build on.

## Requirements

- PRD refs: PRD-TELEM-001, PRD-TELEM-002, PRD-TELEM-003, PRD-TELEM-003a
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` section 7 (`agent_runs`),
  section 13 (`run.record`); `02_TRD.md` sections 2-3 (package boundaries)
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid run record`
- Quality refs: `08_QUALITY_PLAN.md` (validate-before-append, size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-14/qa/validation-contract.md`.

- VAL-002: one `run.recorded` event + one `agent_runs` row with all fields and
  snapshot prices for a known model; no `activity_log` row (PI-1/PI-2).
- VAL-003: `cost_usd` equals the formula over the row's snapshot prices.
- VAL-004: unknown model records with NULL cost/prices + stderr warning + exit 0;
  known model populates cost/prices with no warning (PI-3).
- VAL-005: optional provider fields persist when supplied, NULL when omitted,
  without altering normalized fields.
- VAL-006: in-code matrix has exactly one `run.record` row matching section 13;
  show matches enforce; unknown role rejected before append (CF-007).

## Scope

### Files To Create Or Modify

- `internal/runs/record.go` (new package; PI-1) - `Record(options)`: authorize
  `run.record`, `EnsureProjectionFresh`, `WithLock`, generate `RUN-<hash>` id,
  read pricing config, snapshot prices + compute cost (or NULL + warning),
  build the `run.recorded` event with one `agent_runs` insert effect,
  `DryRunProjectEvent` → `AppendEvent` → `ProjectEvent`.
- `internal/runs/record_test.go` (new) - domain tests for cost, unknown-model,
  provider fields, single-event, no-activity-row.
- `internal/policy/matrix.go` - append the `run.record` row (all seven roles
  true) copying `03_TRD_SCHEMA.md` section 13 verbatim (DEC-026).
- `internal/policy/matrix_test.go` - extend completeness/coverage test for
  `run.record`.
- `cmd/fid/run.go` (new) - `newRunCommand()` + `run record` cobra subcommand
  with flags per the CLI spec; stderr confirmation output; `renderIssueError`
  on failure.
- `cmd/fid/run_test.go` (new) - CLI tests incl. allow/deny role and warning.
- `cmd/fid/root.go` - register `newRunCommand()`.
- `project-docs/specs/02_TRD.md` - sections 2-3: add `internal/runs` domain
  package (agent run recording). Bump any "Last updated".
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - confirm/adjust the `fid run
  record` output example to match the implemented stderr confirmation.

### Files To Read For Context

- `internal/issues/complete.go` (closest write-command analog: authorize →
  EnsureProjectionFresh → WithLock → generate id → effect → DryRun → Append →
  Project)
- `internal/policy/matrix.go`, `internal/policy/authorize.go`
- `internal/config/config.go` (`ReadPricingConfig`, `WorkspacePathsForRoot`)
- `internal/journal/project_effects.go` (generic insert effect projection)
- `cmd/fid/claim.go`, `cmd/fid/concern.go` (cobra command + error-render pattern)

### Out Of Scope

- `fid run list` / `fid run show`, run `status`/lifecycle, and git context
  (all v0.2c).
- `fid stats` reads (Task 3).
- Any `activity_log` projection for telemetry (PI-2).

## Acceptance Criteria

- [ ] `fid run record --issue --role --model [--tokens-in --tokens-out] --by`
  appends one `run.recorded` event and one `agent_runs` row with role/model/
  tokens/started_at/ended_at and snapshot prices for a known model; zero
  `activity_log` rows. (VAL-002)
- [ ] `cost_usd` = `(tokens_in*price_in + tokens_out*price_out)/1e6` from the
  row's snapshot prices for a known model. (VAL-003)
- [ ] Unknown model: row recorded with NULL `cost_usd` and NULL snapshot prices,
  warning naming the model on stderr, exit 0; known model: no warning. (VAL-004)
- [ ] Optional provider fields persist when supplied and are NULL when omitted,
  with normalized fields unchanged. (VAL-005)
- [ ] In-code matrix has exactly one `run.record` row (all seven roles) matching
  section 13; `fid policy show` includes it; unknown role rejected before append
  (non-zero exit, zero new journal lines, error names role + `run.record`).
  (VAL-006)
- [ ] `02_TRD.md` records the `internal/runs` boundary. (doc)

## TDD Plan

- First failing test: `internal/runs` `Record` for a known model returns a
  result with computed `cost_usd` and writes one `agent_runs` row — fails red
  (package absent).
- Additional tests: unknown-model NULL cost + warning; provider fields
  populate/NULL; exactly one journal line and zero `activity_log` rows; matrix
  completeness for `run.record`; CLI allow path (each role) and unknown-role
  reject-before-append.
- Regression risks: lock ordering (refresh projection inside `WithLock`,
  DEC-028); keep `record.go` and `run.go` under 400 lines and functions under
  60; do not weaken the CF-007 validate-before-append seam.

## Documentation Updates

- `project-docs/specs/02_TRD.md` - add `internal/runs` domain package.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - align `fid run record` output
  example with the implementation if it diverges.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `fid policy show` must reflect the new
  `run.record` row (show-vs-enforce parity).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/runs`, `internal/policy`, and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if any test is integration-tagged.
- `go build -o /tmp/fid ./cmd/fid` for the CLI build used by command tests.

## Notes For Developer

- PI-1: `internal/runs` is a new inward-pointing domain boundary; record it in
  `02_TRD.md` in this task and in the completion report.
- PI-2: emit one `run.recorded` event with a single `agent_runs` insert effect;
  do not write `activity_log`.
- PI-3: unknown model is a warning, not a rejection — exit 0, NULL cost/prices.
- The only reject-before-append path here is an unknown/denied role (authorize).
  Keep the CF-007 seam: validate before any journal append.
- `RUN-<hash>` ids via `journal.GenerateID("RUN", existsFn)`.
