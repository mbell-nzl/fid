# Phase 5 - Task 3: `fid issue update` with activity-log audit

## Objective

Implement `fid issue update <id> --by <role>` to mutate the scalar issue fields
(title, priority, labels, spec_md, est_tokens, due_date, owner_role,
owner_actor) on an existing issue, journaling one event whose effects update
`issues` (and `issue_labels` when labels change) and emit an `activity` effect
projected to `activity_log`. This is the first command to exercise the Task 1
validate-before-append seam and the Task 2 `activity` op end to end.

## Requirements

- PRD refs: PRD-ISSUE-015 (labels/priority/estimates/due dates), PRD-ISSUE-021
  (owner role/actor change journaled + projected to activity log)
- TRD refs: TRD §4.3 (write path, update effects with preconditions,
  `activity` effect), TRD-ARCH-004, TRD-ARCH-005 (`cmd/fid` thin adapter)
- Design refs: `06_DESIGN_PLAN_CLI.md` §`fid issue update`
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/issues` ≥ 80%),
  COV-002 (`cmd/fid` ≥ 70%)

## Validation Contract Coverage

- VAL-001: update mutates allowed fields; read-back reflects new values,
  preserves untouched fields; invalid update rejected and leaves issue
  unchanged.
- VAL-002: one update = one journal event with the documented effects; rebuild
  reproduces identical state.
- VAL-003: owner role/actor change journaled and projected to `activity_log`
  with actor + before/after.
- VAL-005 (advances): an invalid update appends zero journal lines (via the
  Task 1 seam).

## Scope

### Files To Create Or Modify

- `internal/issues/update.go` (new) - `Update` domain service: load current
  issue, validate requested changes (reuse field validators from `issue.go`),
  build the update effect with preconditions on changed columns, label
  add/remove effects, and one `activity` effect; route through the Task 1
  validate-before-append seam.
- `internal/issues/update_test.go` (new) - package tests.
- `cmd/fid/issue.go` - wire the `update` subcommand (flags per CLI design) as a
  thin adapter delegating to `internal/issues.Update`.
- `cmd/fid/issue_test.go` (or a new `issue_update_test.go`) - CLI tests.

### Files To Read For Context

- `internal/issues/issue.go` (`Create`, validators, label handling)
- `internal/journal/project.go` (effects, preconditions, activity op)
- `cmd/fid/issue.go` (existing subcommand wiring)
- `06_DESIGN_PLAN_CLI.md` §`fid issue update`

### Out Of Scope

- `--epic` / `--inbox` reparent/promote (v0.2, PRD-ISSUE-025).
- `fid issue label add/remove` standalone commands (v0.1d) — `--labels` here is
  replace-all per the CLI design.
- Status changes (Task 5 `transition`).

## Acceptance Criteria

- [ ] `fid issue update` changes each supported field and `fid issue get
  --format json` reads back the new value while untouched fields are preserved.
  (VAL-001)
- [ ] An invalid update (priority out of 0–5, blank title, unknown label,
  unknown owner role, malformed due_date) returns a non-zero actionable error
  and leaves the issue unchanged with zero journal lines appended.
  (VAL-001, VAL-005)
- [ ] One update produces exactly one journal event with an `update` effect on
  `issues` (+ label effects when labels change) and one `activity` effect;
  `fid rebuild` reproduces identical projected state. (VAL-002)
- [ ] `--owner-role`/`--owner-actor` changes are journaled and projected to one
  `activity_log` row recording actor + from/to. (VAL-003)
- [ ] The `issues.updated_at` column advances and the update effect carries a
  precondition on the prior value of each changed column. (VAL-002)
- [ ] `--by` is required and validated against known roles. (VAL-001)
- [ ] COV-001 (`internal/issues`) and COV-002 (`cmd/fid`) hold.

## TDD Plan

- First failing test: `internal/issues.Update` changes title + priority and a
  follow-up `Get` reflects them. Fails (no Update yet).
- Additional tests: each field individually; labels replace-all; owner change
  emits an activity row; invalid inputs rejected with unchanged read-back and
  unchanged journal line count; one-event-per-update; rebuild parity; CLI flag
  parsing and stderr confirmation line.
- Regression risks: ensure update preconditions match current values to avoid
  spurious replay conflicts; ensure label replace-all diffs correctly
  (add/remove only the delta); ensure no-op updates are handled (reject or
  no-event — choose and test).

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` - confirm the `fid issue update` documented contract
  matches the implementation; update the example output if needed.
- Completion report: cite VAL-001/002/003 and note CF-005 progress (first real
  activity-effect emitter).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `fid issue update` field set must match the
  fields `fid issue get`/`fid issue create` expose; check JSON read-back parity.
- Audit or release-readiness impact: None beyond CF-005 progress.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command as
`PASS` (or `BLOCKED` with reason) in the completion report.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

- `go build` if manual CLI smoke of `fid issue update` is used as evidence for
  VAL-001/003.

## Notes For Developer

- Decide and document the no-op semantics (update with no changed value):
  reject with a clear message or skip the event. Do not append an empty event.
- Reuse the Task 1 validate-before-append seam so an invalid update never
  appends (VAL-005). Do not re-introduce append-before-project.
- Generate the `activity_log` row's `ACT-` id at event-build time; the
  projector applies resolved values (Task 2 contract).
