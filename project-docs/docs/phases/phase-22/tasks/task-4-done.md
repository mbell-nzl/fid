# Phase 22 - Task 4: Read-only repository reconciliation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 3
- Next action: None.
- Human attention required: no
- Recovery authority: TRG-22-007 / DEC-192 / DEC-193; attempt 3.
- The bounded continuation admitted the task-4-owned `run.checkpointed` Event to projection and removed calendar dependence from the retained Stop fixture without changing any accepted contract, gate, schema, public interface, or validation row.

## Summary of what was implemented

- Added the read-only Git repository observer, reconciliation comparison, and `run.checkpoint` activated capability in `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a`.
- Deleted only the `run.checkpointed` progressive fence key, so its accepted `run_observations` insert and `run_runtime` version increment use the shared live/replay projection path.
- Added checkpoint live-apply/replay parity coverage and retained rejection coverage for the eight unavailable recovery Event types.
- Repaired the retained Stop fixture with writer-assigned Event IDs while preserving the `2/1/2` Stop/Attention/link assertion.

## Files created or modified

- `internal/project/repository.go`
- `internal/project/repository_test.go`
- `internal/workflow/reconciliation.go`
- `internal/workflow/reconciliation_test.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `internal/journal/checkpoint_projection_test.go`
- `internal/workflow/stop_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-013: `TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths` in `internal/project/repository_test.go` exercises a clean base and NUL-delimited untracked change. `TestCheckpointedEventProjectsIdenticallyLiveAndReplay` in `internal/journal/checkpoint_projection_test.go` proves the checkpoint observation insert and runtime version increment have live/replay parity.
- VAL-014: `TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions` independently pins staged, unstaged, mode-only, deleted, renamed, untracked, symlink, and excluded dirty-state vectors. `TestReconcileRepositoryRequiresExactDirtyAdoption` accepts the exact observed set and rejects omitted or stale adoption digests.
- VAL-015: `ReconcileRepository` rejects missing base or identity, unmerged paths, active Git operations, and unattributable ancestry with `RecoveryStateConflictError`; `TestReconcileRepositoryRequiresExactDirtyAdoption` covers stale adoption rejection.
- VAL-016: `ObserveRepository` performs read-only Git inspection and its focused fixtures use disposable repositories; no task-4 code mutates Git state.
- DEC-192 convergence: `TestRemainingRecoveryEventsRemainFenced` proves all eight remaining fenced recovery Event types still reject.

## Tests run

- `go test ./internal/project -run '^TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Original red attempt before the repository observer: `undefined: ObserveRepository`.
- `go test ./internal/project -run '^TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/project	0.055s` (terminated `exit_code: 0`)
- `go test ./internal/project -run '^TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Original red vector attempt: `unsupported Git image for "deleted.txt"`; the causal parser repair is in `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a`.
- `go test ./internal/workflow -run '^TestReconcileRepositoryRequiresExactDirtyAdoption$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Original red attempt: `undefined: ReconcileRepository` and `undefined: ValidateReconciliationAdoption`.
- `go test ./internal/workflow -run '^TestReconcileRepositoryRequiresExactDirtyAdoption$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.051s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestActivatedCapabilitiesAcceptRunCheckpoint$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Original red attempt: `activated capabilities reject run.checkpoint`.
- `go test ./internal/journal -run '^(TestCheckpointedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - DEC-192 red attempt: `Event-v2 "run.checkpointed" does not permit effect application before its recovery domain is available`.
- `go test ./internal/journal -run '^(TestCheckpointedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.415s` (terminated `exit_code: 0`)
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable DEC-193 failed receipt: `system clock is behind latest journal date 2026-09-06; refusing to append out of order` (terminated `exit_code: 1`).
- `gofmt -w internal/journal/project_effects.go internal/journal/checkpoint_projection_test.go internal/workflow/stop_test.go && go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	2.428s` (terminated `exit_code: 0`)
- DEC-193 retained-fixture provenance: introducing commit `ac827959f8359473cbd6546ac6bea1ff1fd20751`, `feat(workflow): add typed Stop and Attention facts`, `Refs: phase-22/task-2`. The calendar-only mechanism was `stopRaisedEvent`'s fixed `OccurredAt: time.Date(2026, 9, 5, ...)` and literal `EVT-20260905-00000{1,2}` values after live-clock workspace initialization. `git show --format=fuller --no-ext-diff ac827959f8359473cbd6546ac6bea1ff1fd20751 -- internal/workflow/stop_test.go` establishes the introducing source; `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a` changes only repository/reconciliation files and the capability list, so no accepted task-4 change displaced the assertion. The source repair removes fixed Event date and IDs and uses `journal.PrepareEvent` before rebuilding Effects with the assigned Event ID. Executing command: `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=5`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Immutable repository-wide race/coverage attempt: `FAIL` (terminated `exit_code: 1`). It failed only in `TestEquivalentStopsShareAttentionButKeepDistinctFacts` with `system clock is behind latest journal date 2026-09-06; refusing to append out of order`. Root cause: the calendar-only task-2 fixture. Repair: writer-assigned IDs and date via the production `journal.PrepareEvent` pattern; no assertion, guard, command, or policy changed.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	58.618s	coverage: 35.0% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)				61.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.269s` (terminated `exit_code: 0`)
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=5` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	8.042s` (terminated `exit_code: 0`)

## Deviations from spec

- None. DEC-192 and DEC-193 are accepted bounded scope amendments, not specification deviations.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `59c62b9e1d82c5f698d1e173908446d35b64f539`

## Commit subject

- `fix(recovery): admit checkpoint projection`

## Root report display correction — 2026-09-08

The root delegate decoded escaped tab characters in historical Go-output
receipts following the QA-84 preflight. The original Developer sessions remain
the source of those receipts, including every recorded outcome, duration,
coverage value and observed termination. This formatting correction records
no new execution or attestation. Original encoded text remains in Git history.
