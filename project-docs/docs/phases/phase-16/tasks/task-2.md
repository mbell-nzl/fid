# Phase 16 - Task 2: CF-011 backfill part 1 - issue create / ac add / ac verify / complete

## Objective

Begin the CF-011 broad `activity_log` reconciliation (DC-A, approved at kickoff):
backfill the four
issue-lifecycle write commands that mutate or append an issue fact but do not yet
project an `activity_log` row, so the timeline (`fid activity list`, Task 1) is
consistent across them.

Each writer gains exactly **one** `activity` effect inside its existing single
journal event (PI-2), `ref_type=issue`, mirroring the existing per-command
emitters in shape. Additive and rebuild-safe: prior events replay unchanged, no
`event_schema` bump, no upgrader.

Writers in this task (all in `internal/issues`):

| Writer | Event type (existing) | New activity event name |
|---|---|---|
| `issue create` | `issue.created` | `issue_created` |
| `issue ac add` | `issue.ac_added` | `ac_added` |
| `issue ac verify` | `issue.ac_verified` | `ac_verified` |
| `issue complete` | (completion) | `task_completed` |

`issue complete` builds its effects via the shared `issues.BuildCompletionEffects`
(CF-013), so adding the `task_completed` activity there means `fid dev finish`/
`fid dev block` inherit it automatically (PI-6) — keep each composite at one
event and confirm that in the completion report.

## Requirements

- PRD refs: PRD-OUT-017 (timeline consistency); PRD-ISSUE-001..005 (create),
  PRD-AC-001..005 (ac), PRD-DEV-001..004/009 (complete)
- TRD refs: `03_TRD_SCHEMA.md` §9 (`activity_log`), §14.2 (ref integrity:
  `ref_type` must be `issue`/`requirement`/`decision`/`escalation`), §14.6
  (activity projection discipline)
- Design refs: `06_DESIGN_PLAN_CLI.md` §8 (activity entry shape)
- Quality refs: `08_QUALITY_PLAN.md` (validate-before-append, one event per write)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-005: `issue create`, `ac add`, `ac verify`, `complete` each append exactly
  one `activity_log` row (named events) inside their single event; rejected
  attempts add zero rows; `fid activity list --on issue/<id>` reflects each.

## Scope

### Files To Create Or Modify

- `internal/issues/create.go` - append one `activity` effect (`issue_created`,
  `ref_type=issue`, `ref_id`=new issue id, `details` = kind/title/parent) to the
  `issue.created` event's effects.
- `internal/issues/ac.go` - append one `activity` effect (`ac_added`,
  `ref_type=issue`, `details` = ac_id/ordinal) to the `issue.ac_added` event.
- `internal/issues/ac_verify.go` - append one `activity` effect (`ac_verified`,
  `ref_type=issue`, `details` = ac_id/verifier_role/status) to the
  `issue.ac_verified` event.
- `internal/issues/complete.go` (and/or the shared `BuildCompletionEffects`) -
  append one `activity` effect (`task_completed`, `ref_type=issue`, `details` =
  completion_id/status/attempt_number) so `fid issue complete` **and** the
  `internal/dev` composites that consume the shared builder emit it (PI-6).
- `internal/issues/write_helpers.go` - reuse the existing `activityIDExists` +
  `journal.GenerateID("ACT", ...)` helpers; add a small shared activity-effect
  builder if it removes duplication (keep functions ≤60 lines).
- `internal/issues/*_test.go` - extend create/ac/ac_verify/complete tests:
  exactly one new event, exactly one new `activity_log` row with the named event
  + `ref_type=issue`, rejection path adds zero activity rows, and the existing
  "no activity row" assertions in `ac_test.go`/`complete_test.go` are updated to
  the new expected behavior (decision-backed, not weakening).
- `internal/dev/*_test.go` - assert `fid dev finish`/`fid dev block` still emit
  exactly one event now carrying both `status_changed` and `task_completed`
  activity rows (PI-6).
- `cmd/fid/*_test.go` - where create/ac/complete CLI tests assert activity-log
  counts, update to the new expected counts.

### Files To Read For Context

- `internal/issues/transition.go` (canonical `activity` effect shape to mirror:
  `Op:"activity"`, `Table:"activity_log"`, `PK`, `Values` with
  `journal_event_id`/`ref_type`/`ref_id`/`event`/`details`/`actor`/`occurred_at`)
- `internal/issues/update.go`, `internal/issues/label.go` (more emitter examples
  incl. `details` shaping)
- `internal/issues/complete.go` + the `BuildCompletionEffects` builder (CF-013)
- `internal/issues/ac.go`, `internal/issues/ac_verify.go`, `internal/issues/create.go`
- `internal/dev/finish.go`, `internal/dev/block.go` (consume the shared builder)
- DC-A (approved) and PI-2/PI-6 in the validation contract

### Out Of Scope

- `concern raise`/`concern address`/`unclaim`/`link add` backfill (Task 3).
- Excluded writers: `fid label add` (label creation — invalid `ref_type`),
  `fid run record`, `fid requirement sync` (DC-A).
- Any new event type or `event_schema` bump (additive effect only, PI-2).

## Acceptance Criteria

- [ ] `fid issue create` appends exactly one new event whose effects include one
  `issue_created` `activity_log` row (`ref_type=issue`); `fid activity list --on
  issue/<id>` shows it. (VAL-005)
- [ ] `fid issue ac add` and `fid issue ac verify` each append exactly one event
  with one `ac_added` / `ac_verified` activity row. (VAL-005)
- [ ] `fid issue complete` appends exactly one event with one `task_completed`
  activity row; `fid dev finish`/`fid dev block` still emit exactly one event
  each, now carrying `status_changed` + `task_completed` activity. (VAL-005, PI-6)
- [ ] A rejected attempt of each writer (e.g. denied role / invalid input)
  appends zero new journal lines and zero activity rows. (VAL-005)

## TDD Plan

- First failing test: `fid issue create` now writes one `issue_created`
  `activity_log` row — fails red against current behavior (today zero rows).
- Additional tests: ac add/verify/complete each add one named row; reject paths
  add zero rows; one-event invariant per writer; `dev finish`/`dev block` still
  one event with the extra activity row.
- Regression risks: do **not** introduce a second event — the activity effect
  joins the existing event's effects (PI-2); preserve validate-before-append
  (build all effects, then append once); mirror the `transition.go` activity
  shape byte-for-byte so replay stays deterministic; update the existing
  "activity_log rows == 0" assertions to the new decision-backed counts rather
  than deleting them; keep functions ≤60 lines.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md` - extend the `activity_log.event`
  vocabulary comment (§9) to list the new event names (`issue_created`,
  `ac_added`, `ac_verified`, `task_completed`). Bump "Last updated".
- Completion report: record the DC-A scope (which writers gained rows) and the
  PI-6 composite confirmation.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: the `internal/dev` composites must stay one event
  each (PI-6); confirm `fid dev finish`/`block` effect counts.
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
`internal/issues`, `internal/dev`, and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required: this task changes
  `cmd/fid`-exercised write paths and the package set under test includes
  integration-tagged tests; record PASS/BLOCKED (Phase 14 ESC-001 lesson).
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the create/ac/
  complete CLI tests.

## Notes For Developer

- PI-2: each writer still emits **one** event. Add the `activity` effect to the
  existing effect slice; do not append a second journal line.
- The §14.2 trigger only allows `ref_type` in `issue`/`requirement`/`decision`/
  `escalation`; all four writers here use `ref_type=issue`.
- Updating the existing `ac_test.go`/`complete_test.go` "zero activity rows"
  assertions is a **decision-backed behavior change** (DC-A), not test
  weakening — cite DC-A/VAL-005 in the test and the completion report.
- Reuse `journal.GenerateID("ACT", ...)` + `activityIDExists` and the
  `transition.go` activity `Values` shape so rebuild stays deterministic.
- DC-A (CF-011 scope/vocabulary) is human-confirmed for Phase 16; do not reopen
  it unless implementation reveals a new spec contradiction.
