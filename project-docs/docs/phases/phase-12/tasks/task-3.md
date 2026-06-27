# Phase 12 - Task 3: `fid link remove`

## Objective

Implement `fid link remove`: remove a declared link and its auto-created inverse
in one journal event, project both removals to `activity_log` (PRD-LINK-006a),
and reject removal of a nonexistent link before append. Reuses the Task 2 links
domain registry and endpoint helpers.

## Requirements

- PRD refs: PRD-LINK-006a, PRD-ERR-003
- TRD refs: TRD-ARCH-010; `03_TRD_SCHEMA.md` sections 10, 13
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid link remove`
- Quality refs: `08_QUALITY_PLAN.md` COV-001/COV-002; CF-007 validate-before-append
- Decisions: DEC-028 (held-lock projection freshness); CF-011 (activity-log
  discipline; broad reconciliation stays v0.1e)
- Depends on: Task 2 (`internal/links` add + registry)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-12/qa/validation-contract.md`.

- VAL-009: `fid link remove` removes declared + inverse in one event, projects
  activity-log removal evidence, and confirms both directions; removing a
  nonexistent link rejects before append.

## Scope

### Files To Create Or Modify

- `internal/links/remove.go` (or equivalent) and tests - remove service building
  a single journal event with two `delete` effects (declared + inverse) plus the
  `activity_log` removal effect; nonexistent-link rejection.
- `cmd/fid/link.go` and `cmd/fid/link_test.go` - add the `remove` subcommand
  with `--by`; stderr confirmation naming both removed directions.

### Files To Read For Context

- `internal/claims/acquire.go` for the activity-log effect shape and ACT id
  generation.
- `internal/journal/project_effects.go` for `delete` effect + precondition and
  `ensureSingleRowAffected` semantics.
- `internal/links/` (Task 2) for the relation/inverse registry and
  endpoint-resolution helpers.
- `migrations/0005_v01_ref_triggers.up.sql` `activity_log_ref_integrity_guard`
  (the activity ref must resolve to an existing issue/requirement).

### Out Of Scope

- `fid link list` (Task 4).
- Compensating-event `fid event revert` (v0.2).
- Schema migrations.

## Acceptance Criteria

- [ ] `internal/links` owns removal business rules; `cmd/fid` contains only flag
  wiring and confirmation output (VAL-009).
- [ ] `fid link remove <from> <relation> <to> --by <role>` deletes the declared
  row and its inverse in exactly one journal event (two `delete` effects)
  (VAL-009).
- [ ] The remove event projects an `activity_log` row recording the removal
  (referencing the declared link's `from` entity, which resolves under the
  activity ref-integrity trigger) (VAL-009).
- [ ] The stderr confirmation names both removed directions (declared +
  inverse) (VAL-009).
- [ ] Removing a nonexistent link rejects before append with non-zero exit, zero
  new journal lines, and an error naming the missing relation and entities
  (VAL-009).
- [ ] `link.remove` authorizes via the matrix (added in Task 2); an unknown role
  is rejected before append (VAL-009).

## TDD Plan

- First failing test: after adding `blocks` between two issues, `Remove` deletes
  both `links` rows in one event and leaves zero rows for that pair.
- Additional tests: an `activity_log` removal row exists after remove; remove of
  a nonexistent link rejects with zero append; confirmation text names both
  directions; remove of an issue -> requirement `implements` link removes the
  `implemented_by` inverse; unknown-role rejection.
- Regression risks: both deletes must be effects of one event; activity-log
  writes must use `Op: "activity"`; validate before append (dry-run project);
  use `EnsureProjectionFreshWithLock` inside `WithLock` (DEC-028); do not let a
  partial delete leave a dangling inverse.

## Documentation Updates

- `06_DESIGN_PLAN_CLI.md` `fid link remove` - reconcile only if implemented
  output differs; otherwise record verified-no-change.
- `project-docs/docs/phases/phase-12/decisions.md` - record any deviation found.
- `project-docs/docs/phases/phase-12/tasks/task-3-done.md` - completion report
  with command results, coverage, and commit hash.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None new (matrix rows landed in Task 2).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
may be used during the red/green loop, but it does not substitute for the full
gate at close.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

- `go build -o /tmp/fid ./cmd/fid` - required because this task adds a CLI
  subcommand; record PASS.
- `go test -tags integration -race ./...` - run if any integration fixture is
  touched; otherwise record N/A with reason.

## Notes For Developer

- Use red/green TDD. Do not add a schema migration unless a spec-backed blocker
  is recorded first.
- One logical write = one event: declared delete, inverse delete, and the
  activity-log row are all effects of the single remove event.
- Completion report commit fields must be filled from `git rev-parse HEAD`.
