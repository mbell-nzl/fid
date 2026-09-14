# Phase 23 - Task 1: Review and judgment schema baseline

## Objective

Land the Phase 23 persistence baseline: the accepted textual changes to
migrations 0013 and 0015, the new migration `0017_p23_review_judgment` with its
down, the `evidence_records` rebuild, and the progressive availability fence
that rejects every Phase 23 Event until its domain exists.

## Validation Contract Coverage

- VAL-001: Migration 0017 creates exactly the declared objects and nothing else.
- VAL-002: The 0013 and 0015 baseline changes are exactly the accepted text and
  the Obligation registration discriminator enforces both states.
- VAL-003: The progressive fence rejects every unimplemented Phase 23 Event.
- VAL-004: `foreign_key_check` and `integrity_check` are clean after 0017 and
  after a full rebuild.
- VAL-005: The `evidence_records` rebuild preserves every prior receipt row.

Full sources, pass and fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-23/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under the accepted slice contract; execute after the Runner's
  required QA readiness check.
- `schema.md` sections 1, 2, 3, 6.1, and 7 are the exact authority. Do not add,
  rename, reorder, or re-type a single column, index, trigger, or CHECK.
- DEC-178: shipped migration SQL and the canonical slice DDL must agree exactly.
  Every line of 0013 and 0015 outside the accepted replacement table stays
  byte-identical.
- This is a deliberate pre-adoption baseline replacement requiring fresh init
  and rebuild. It is not an upgrade of old fid data and needs no converter.
- DEC-192 fence protocol: add `unavailableJudgmentEvents` to
  `internal/journal/project_effects.go` rejecting all nine Phase 23 Event types
  listed in `schema.md` section 6. Delete no key in this task. Later owners
  delete exactly their own key; task 12 removes the empty map and its guard.
- Register the nine Event types and their operation mappings in
  `internal/journal/current_event_catalog.go` with the effect rules in
  `schema.md` section 6.2, so a fenced type is rejected by the fence and not by
  an unknown-type path.
- Manifest `projection_schema` becomes 17 in this task. It becomes 18 in task
  15; do not anticipate that value.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this accepted
  change, and every owned path executes under the commands below.

## Scope

### Files To Create Or Modify

- `migrations/0013_p21_operational_kernel.up.sql` - only the accepted rows.
- `migrations/0015_p22_recovery.up.sql` - only the accepted rows.
- `migrations/0017_p23_review_judgment.up.sql`, `.down.sql` - new.
- `internal/db/` - migration loading, latest-schema expectations, tests.
- `internal/journal/project_effects.go`, `current_event_catalog.go` and tests.
- `internal/export/` - manifest projection schema and its fixtures.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-23/schema.md` - sections 1, 2, 3, 6, 7.
- `project-docs/docs/phases/phase-23/contracts.md` - sections 1 and 2 only.
- `project-docs/docs/phases/phase-23/execution-scope.md`,
  `v1-convergence-inventory.md`, `contract-coverage.md`.
- `migrations/0015_p22_recovery.up.sql` and `0016_p22_remove_escalations.up.sql`
  as the closest precedent for shape and down behavior.

### Out Of Scope

- Any domain rule, CLI leaf, Human text, or capability advertisement.
- Migration 0018 and every V1 removal.
- Deleting any fence key.

## Acceptance Criteria

- [ ] Latest-schema object inventory equals the accepted DDL exactly. (VAL-001)
- [ ] Both invalid and both valid registration-discriminator inserts behave as
  accepted, and no other baseline line changed. (VAL-002)
- [ ] Every one of the nine Phase 23 Event types is rejected by the fence with
  a byte-identical projection afterwards. (VAL-003)
- [ ] `PRAGMA foreign_key_check` and `PRAGMA integrity_check` are clean after
  0017 and after a full journal rebuild. (VAL-004)
- [ ] Every pre-existing receipt Evidence row survives the rebuild unchanged
  with the three new columns null. (VAL-005)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit, phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

`TestMigration0017CreatesExactlyTheAcceptedObjects` fails against the absent
migration. Capture the failure before product edits, then implement the
smallest behavior that satisfies it. A build failure alone does not prove the
intended rejection.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md`: add section 8.4 incorporating
  `schema.md` by reference, per the accepted contract.
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md`: implementation,
  red/green, every exact command with output token and observed termination,
  assertion evidence, convergence provenance, and deviations.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory from this task.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

No exception to the full task gate. Commands whose expected pass is non-zero:
none. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination, as
required by `AGENTS.md`. An interrupted or hung command is never `PASS`;
preserve all failed attempts and rerun the byte-identical command after an
authorized causal repair. Report lint is a Runner control after the report
commit, not a Developer gate.

## Notes For Developer

- This task belongs to QA Review window 1.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-1`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
