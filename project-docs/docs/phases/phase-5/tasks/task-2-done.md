# Phase 5 - Task 2: `activity` effect op + activity-log projection discipline (CF-005) - Completion Report

## Summary of what was implemented

- Added `activity` effect support to the journal projector in `internal/journal/project.go`.
- Enforced the `activity_log` projection discipline in the projector by rejecting generic `insert`, `update`, and `delete` effects targeting `activity_log`; only `activity` may write that read-model table.
- Normalized `activity.values.details` to deterministic JSON text before insert so replay and rebuild preserve stable `activity_log` rows.
- Added `VAL-004` tests covering one-row projection, multiple activity effects, unresolved polymorphic refs, rebuild determinism, and the non-`activity` writer rejection.
- Updated existing synthetic journal fixtures in package tests to emit `activity` for `activity_log` rows so the full test suite remains aligned with the new discipline.

## Files created or modified

- `internal/journal/project.go`
- `internal/journal/project_test.go`
- `internal/journal/phase2_rebuild_test.go`
- `internal/db/ref_triggers_test.go`
- `project-docs/docs/phases/phase-5/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-5/tasks/task-2-done.md` - added the Task 2 completion report.
- `project-docs/specs/03_TRD_SCHEMA.md` - verified-no-change. Existing §14.6 already states that `activity_log` writes occur only during projection, and the task-specific deferred SQL guard note is captured here in the completion report.
- `project-docs/specs/02_TRD.md` - verified-no-change. Existing §4.3 already documents `activity` as a first-class effect op with fully resolved values.

## Validation Contract Coverage

- `VAL-004` - PASS. `TestProjectEventProjectsActivityEffectIntoActivityLog` proves one `activity` effect inserts exactly one `activity_log` row and stores `details` as JSON text.
- `VAL-004` - PASS. `TestProjectEventProjectsEachActivityEffectIntoOneRow` proves multiple `activity` effects in one event project one row each.
- `VAL-004` - PASS. `TestProjectEventRejectsActivityEffectWithUnresolvedRef` proves the §14.2 polymorphic-ref backstop still rejects unresolved `(ref_type, ref_id)` pairs during projection and leaves `activity_log` unchanged.
- `VAL-004` - PASS. `TestRebuildWorkspaceReplaysActivityLogDeterministically` proves replay/rebuild reproduces identical `activity_log` rows and projection metadata for a journal containing `activity` effects.
- `VAL-004` - PASS. `TestProjectEventRejectsGenericWritesToActivityLog` proves `activity_log` is populated only through the `activity` projection path; generic `insert`, `update`, and `delete` effects are rejected.
- `COV-001` - PASS. `go test -race -coverprofile=coverage.out ./...` reported `internal/journal` coverage at `80.6%`, meeting the `>= 80%` package gate named by the task spec.

## Tests run

- `go test ./internal/journal -run TestProjectEventProjectsActivityEffectIntoActivityLog -count=1` - FAIL before implementation: `project event "EVT-20260613-000101": unsupported effect op "activity"`
- `go test ./internal/journal -run TestProjectEventProjectsActivityEffectIntoActivityLog -count=1` - PASS
- `go test ./internal/journal -run 'TestProjectEventProjectsActivity|TestProjectEventRejectsActivity|TestProjectEventRejectsGenericWritesToActivityLog|TestRebuildWorkspaceReplaysActivityLogDeterministically' -count=1` - PASS
- `go test ./internal/journal ./internal/db -count=1` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- Touched `internal/journal/phase2_rebuild_test.go` and `internal/db/ref_triggers_test.go` in addition to the task spec's primary files. Both changes were limited to existing synthetic journal fixtures so the required package/repo tests reflect the new `activity_log` discipline; no product behavior outside Task 2 changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `27a2350`

## Commit subject

- `feat(journal): add activity projection discipline`
