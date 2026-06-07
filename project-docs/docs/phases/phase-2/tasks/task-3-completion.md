## Summary of what was implemented

- Added migration `0004_v01_kind_triggers` to enforce the Phase 2 kind backstops:
  `issues.parent_id` hierarchy, kind guards for `acceptance_criteria`,
  `task_completions`, `epic_metadata`, `triage_events.epic_id`,
  `decisions.epic_id`, and `issues.kind` immutability.
- Added raw-SQL trigger tests in `internal/db` that prove both reject and accept
  states for VAL-009..012, plus a journal-projection control showing a valid
  issue event still projects with the new triggers active.
- Kept `internal/db.ProjectionSchemaVersion` at `3`. I did not bump it because
  Task 3 adds trigger DDL only; existing projected rows do not require a fresh
  generation to become structurally valid, and the migration applies the new
  backstops in place.
- Implemented the accepted `task_completions` kind-guard resolution as
  **task or subtask**, following `03_TRD_SCHEMA.md` §14.1 rather than the
  indicative `02_TRD.md` §4.1 inventory note.

## Files created or modified

- `migrations/0004_v01_kind_triggers.up.sql`
- `migrations/0004_v01_kind_triggers.down.sql`
- `internal/db/kind_triggers_test.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/issues/issue_test.go`
- `project-docs/docs/phases/phase-2/tasks/task-3-completion.md`

## Documentation updated

- Added this completion report at
  `project-docs/docs/phases/phase-2/tasks/task-3-completion.md`.

## Validation Contract Coverage

- `VAL-009` PASS:
  `TestIssueParentKindTriggerRejectsInvalidHierarchiesAndAllowsValidRows`
  rejects epic-with-parent, task-with-non-epic-parent, subtask-without-parent,
  and subtask-with-epic-parent; valid epic, standalone task, task-under-epic,
  and subtask-under-task inserts succeed.
- `VAL-010` PASS:
  `TestIssueKindGuardTriggersRejectWrongKindsAndAllowValidRows` rejects
  `acceptance_criteria` and `task_completions` rows on an epic and rejects
  `epic_metadata` on a task; it accepts `acceptance_criteria` and
  `task_completions` on both task and subtask, and accepts `epic_metadata` on
  an epic. `task_completions` is guarded as task-or-subtask per §14.1.
- `VAL-011` PASS:
  `TestEpicReferenceKindGuardTriggersRejectWrongKindsAndAllowValidRows`
  rejects `triage_events.epic_id` and `decisions.epic_id` when they point at a
  task, accepts valid epic references for both tables, and accepts
  `decisions.epic_id = NULL`.
- `VAL-012` PASS:
  `TestIssueKindImmutabilityTriggerRejectsKindChangesAndAllowsOtherUpdates`
  rejects `UPDATE issues SET kind = ...` and accepts a non-kind update on the
  same row.
- Trigger false-positive control:
  `TestValidJournalProjectedIssueStillAppliesWithKindTriggersActive` proves a
  valid journal-applied epic plus task still projects under the new backstops.

## Tests run

- `go test ./internal/db -run TestIssueParentKindTriggerRejectsInvalidHierarchiesAndAllowsValidRows`
  Red: failed before migration `0004` existed because invalid hierarchy inserts
  were still accepted.
- `go test ./internal/db -run 'TestIssueParentKindTriggerRejectsInvalidHierarchiesAndAllowsValidRows|TestIssueKindGuardTriggersRejectWrongKindsAndAllowValidRows|TestEpicReferenceKindGuardTriggersRejectWrongKindsAndAllowValidRows|TestIssueKindImmutabilityTriggerRejectsKindChangesAndAllowsOtherUpdates|TestValidJournalProjectedIssueStillAppliesWithKindTriggersActive'`
  Green after adding `0004`.
- `go test ./internal/db`
- `go test ./internal/issues -run TestGetMapsPopulatedNullableFields`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `NOT REQUIRED` `go test -tags integration -race ./...`
  Task 3 trigger tests were implemented as ordinary temp-DB tests in
  `internal/db`, not behind the `integration` build tag.
- `NOT REQUIRED` `go build -o /tmp/fid ./cmd/fid`
  Task 3 does not change CLI wiring or require the conditional build step.

Coverage notes from `go tool cover -func=coverage.out`:

- `internal/journal` 80.6%
- `internal/issues` 85.1%
- `cmd/fid` 70.5%
- total 72.6%

## Deviations from spec

- Updated `internal/issues/issue_test.go` even though it was outside the Task 3
  file list in the spec. The existing fixture inserted a task under a task,
  which is invalid under VAL-009; I changed only the parent fixture row to an
  epic so the existing `Get()` behavior test remains valid under the new
  backstop without weakening the trigger.
- Hardened `internal/db/migrate_v01_tables_test.go` from `Steps(-2)` to
  `Migrate(1)` so the Task 1 Phase-1-schema regression test continues to target
  its intended schema boundary after adding migration `0004`.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e38ae45`

## Commit subject

- `feat(db): add kind backstop triggers`
