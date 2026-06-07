## Summary of what was implemented

- Added migration `0006_v01b_sequence_scheduling` to enforce the Phase 4 `issues.sequence` backstop for tasks against `epic_metadata.scheduling`.
- Covered the raw SQL reject/accept matrix for sequential, prioritised, unattached, metadata-less, and duplicate-sequence task states.
- Added down-migration evidence that rolling back `0006` removes only the new sequence backstop while preserving the earlier Phase 2 kind triggers.
- Added journal projection evidence proving valid sequenced task projection succeeds with the trigger active and invalid projection fails with the sequence-specific constraint.
- Updated existing repo fixtures that legitimately create tasks under epics so they now seed matching `epic_metadata` rows required by the new invariant.

## Files created or modified

- `migrations/0006_v01b_sequence_scheduling.up.sql`
- `migrations/0006_v01b_sequence_scheduling.down.sql`
- `internal/db/kind_triggers_test.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/db/ref_triggers_test.go`
- `internal/journal/project_test.go`
- `internal/journal/v01_tables_test.go`
- `internal/journal/phase2_rebuild_test.go`
- `cmd/fid/issue_test.go`
- `internal/issues/issue_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-4/tasks/task-1-completion.md`

## Validation Contract Coverage

- `VAL-001` PASS.
- Exact evidence:
  `internal/db/kind_triggers_test.go`: `TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows` proves:
  sequential epic task with sequence succeeds; sequential epic task without sequence fails; prioritised epic task with sequence fails; unattached task with sequence fails; task under metadata-less epic fails; duplicate sibling sequence fails.
  `internal/db/migrate_v01_tables_test.go`: `TestMigration0006DownRemovesSequenceSchedulingBackstopOnly` proves the `0006` down path removes only the new backstop and leaves the earlier parent-kind trigger active.
  `internal/journal/project_test.go`: `TestProjectEventHonoursSequenceSchedulingBackstop` proves journal projection respects the trigger: a valid sequenced task projects cleanly when `epic_metadata` is projected first, and an invalid sequential task without sequence is rejected during projection.
  `go test -race -coverprofile=coverage.out ./...`: PASS, including `cmd/fid` coverage `74.3%`, `internal/journal` coverage `80.7%`, and `internal/issues` coverage `85.1%`.
  `go tool cover -func=coverage.out`: total statement coverage `75.2%`.

## Tests run

- `go test ./internal/db -run TestIssueSequenceSchedulingTriggerRequiresSequenceForSequentialEpicTasks`
  Red: failed with `error = nil, want substring "issues: task under sequential epic requires sequence"`.
- `go test ./internal/db ./internal/journal -run 'TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows|TestMigration0006DownRemovesSequenceSchedulingBackstopOnly|TestProjectEventHonoursSequenceSchedulingBackstop'`
  Red: failed before migration with missing backstop errors.
- `go test ./internal/db ./internal/journal -run 'TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows|TestMigration0006DownRemovesSequenceSchedulingBackstopOnly|TestProjectEventHonoursSequenceSchedulingBackstop'`
  Green after adding `0006`: PASS.
- `go test ./internal/db ./internal/journal`
  Green after fixture updates: PASS.
- `gofmt -l cmd internal`
  PASS.
- `golangci-lint run ./...`
  PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...`
  PASS.
- `go tool cover -func=coverage.out`
  PASS.

## Exact quality-gate results

- `gofmt -l cmd internal`: PASS.
- `golangci-lint run ./...`: PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...`: PASS.
- `go tool cover -func=coverage.out`: PASS (`total: (statements) 75.2%`).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- TBD

## Commit subject

- TBD
