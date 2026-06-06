## Summary of what was implemented

- Added journal replay and projection services that apply each event in its own SQLite transaction and persist the `projection_meta` checkpoint.
- Added freshness handling for clean reads, catch-up, and rebuild decisions based on `projection_meta` plus journal checkpoint inspection.
- Added structured corruption and precondition-conflict errors, plus the thin `fid rebuild` Cobra adapter that rebuilds under the write lock and reports replay counts.

## Files created or modified

- `cmd/fid/root.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/rebuild_test.go`
- `internal/db/projection_meta.go`
- `internal/journal/project.go`
- `internal/journal/project_test.go`
- `internal/journal/replay.go`
- `internal/journal/replay_test.go`
- `testdata/journal/replay-v1.golden.jsonl`
- `testdata/journal/corrupt-tail.jsonl`
- `testdata/journal/conflicting-precondition.jsonl`
- `testdata/journal/replay-catchup-event.json`
- `project-docs/docs/phases/phase-1/tasks/task-4-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-4-completion.md`

## Validation Contract Coverage

- `VAL-008` PASS: `TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase` proves delete-DB + rebuild recreates identical projected state and checkpoint from journal data only.
- `VAL-009` PASS: `TestReplayJournalAppliesGoldenV1JournalInFileOrder` replays the committed v1 golden journal through the upgrader path and projects the expected SQLite state.
- `VAL-010` PASS: `TestReplayJournalReportsMalformedTrailingLine` and `TestRebuildCommandReportsCorruptJournal` report corrupt trailing JSONL with journal file and line information.
- `VAL-011` PASS: `TestReplayJournalStopsOnPreconditionConflict` stops replay on a mismatched precondition and exposes the structured conflict details.
- `VAL-012` PASS: `TestProjectEventAppliesEffectsAndUpdatesProjectionMeta` and `TestEnsureProjectionFreshUsesProjectionMetaMatrix` cover checkpoint writes plus clean-match, missing-meta, schema-mismatch, position-mismatch, hash-mismatch, and catch-up decisions.

## Tests run

- `go test ./internal/journal -run TestProjectEventAppliesEffectsAndUpdatesProjectionMeta`
- `go test ./internal/journal -run 'TestReplayJournal|TestEnsureProjectionFreshUsesProjectionMetaMatrix|TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase|TestProjectEventAppliesEffectsAndUpdatesProjectionMeta'`
- `go test ./cmd/fid -run TestRebuildCommand`
- `go test ./...`
- `go test -race -coverprofile=coverage.out ./...`

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
- `go tool cover -func=coverage.out` PASS
- `go build -o /tmp/fid ./cmd/fid` PASS
- `go test -tags integration -race ./...` Not triggered: replay/rebuild tests were kept in the default test packages and not placed behind the `integration` tag.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- `02_TRD.md` §4.3 says an event-id/hash mismatch at the recorded position should stop as projection inconsistency, while Task 4 / `VAL-012` / the Phase 1 QA plan require rebuild on last-event-hash mismatch. This implementation follows the task and validation contract; the spec wording should be reconciled before later phases build on this behavior.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `dc49d3d`

## Commit subject

- `feat(journal): add replay and rebuild support`
