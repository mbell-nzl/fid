# Phase 4 - Task 1: Sequence Scheduling Backstop

## Objective

Add the cross-table `issues.sequence` backstop that Phase 2 carried forward as
CF-004. The trigger must enforce sequence semantics using `epic_metadata` while
still allowing valid journal projection and rebuild. This isolates the schema
change before Phase 4 expands `fid issue create`.

## Requirements

- PRD refs: PRD-ISSUE-013, PRD-ISSUE-014
- TRD refs: `03_TRD_SCHEMA.md` section 14.1 (`issues.sequence`), TRD-ARCH-006,
  TRD-SEC-004
- Design refs: `07_DESIGN_PLAN_FLOWS.md` section 1.1 (sequential epic pattern)
- Quality refs: `08_QUALITY_PLAN.md` section 5.4, DEC-007 COV-001 carry-in
- Carry-forward: CF-004 in `project-docs/docs/current/decisions_active.md`

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- VAL-001: database and projection backstops enforce sequential/prioritised
  task sequence semantics.

## Scope

### Files To Create Or Modify

- `migrations/0006_v01b_sequence_scheduling.up.sql`
- `migrations/0006_v01b_sequence_scheduling.down.sql`
- `internal/db/kind_triggers_test.go` or a new focused `internal/db/*sequence*_test.go`
- `internal/journal/*test.go` for projection/rebuild evidence with the new
  trigger active
- `project-docs/docs/phases/phase-4/tasks/task-1-completion.md`

### Files To Read For Context

- `migrations/0004_v01_kind_triggers.up.sql`
- `migrations/0005_v01_ref_triggers.up.sql`
- `internal/db/migrate.go`
- `internal/db/kind_triggers_test.go`
- `internal/journal/project.go`
- `project-docs/specs/03_TRD_SCHEMA.md` section 14.1

### Out Of Scope

- `fid issue create` flags and CLI behavior (tasks 2-3).
- `fid ready` scheduling selection (Phase 6).
- Any workflow transition logic (Phase 5).

## Acceptance Criteria

- [ ] Migration `0006` adds a reversible trigger/backstop for `issues.sequence`.
  (VAL-001)
- [ ] Raw SQL tests prove all states: sequential epic task with sequence
  succeeds; sequential epic task without sequence fails; prioritised epic task
  with sequence fails; unattached task with sequence fails; task under an epic
  missing `epic_metadata` fails; duplicate sibling sequence fails. (VAL-001)
- [ ] Journal projection/rebuild tests prove valid event-projected rows pass with
  the trigger active. (VAL-001)
- [ ] The migration down path removes only the new backstop and preserves earlier
  Phase 2 triggers. (VAL-001)

## TDD Plan

- First failing test: raw SQL insert of a task under a sequential epic without a
  sequence must fail with a constraint-specific error.
- Additional tests: valid sequential task with sequence; invalid prioritised
  task with sequence; invalid unattached task with sequence; invalid
  metadata-less parent epic; duplicate sibling sequence; valid journal
  projection control.
- Regression risks: false positives for subtasks, stale migration step tests,
  and trigger ordering with the existing kind-hierarchy trigger.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-1-completion.md` - record
  VAL-001 evidence and exact command outcomes.
- Do not edit product specs unless implementation proves `03_TRD_SCHEMA.md`
  section 14.1 wrong; if that happens, stop and escalate.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: None yet; this is a persistence backstop.
- Audit or release-readiness impact: migration/trigger behavior is a
  compatibility surface; record any deviation.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted commands
are failed gates.

Conditional commands:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- None. No `integration` build tag is currently in use; CLI build evidence is
  required by later user-facing tasks.

## Output Requirements

- Write `project-docs/docs/phases/phase-4/tasks/task-1-completion.md`.
- Completion report must cite VAL-001 and include all required command outcomes.

## Notes For Developer

- Keep this task schema-focused. Do not add issue-create CLI flags here.
- Use existing migration naming/order patterns.
- If a reliable SQL trigger cannot express the missing-`epic_metadata` case,
  stop and report the gap rather than weakening VAL-001.
