# Phase 22 - Task 4: Read-only repository reconciliation - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The task cannot implement or exercise the required `run.checkpointed` mutation because `internal/journal/project_effects.go` rejects that Event before effects apply, while the accepted task-4 scope excludes that required file. The required full race/coverage gate also failed in retained `internal/workflow/stop_test.go` on its fixed 2026-09-05 Event date after the journal advanced to 2026-09-06; that assertion predates Phase 22 (`ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase-22/task-2) and is not displaced by this task's repository change.
- Recovery attempts: 0
- Next action: Phase Runner must obtain a bounded task-4 scope amendment for the `run.checkpointed` projection admission (and, separately, a decision-free retained-gate repair authority for the fixed-date Stop fixture if it remains required), then dispatch a same-task continuation.
- Human attention required: no

## Summary of what was implemented

- Added a read-only Git repository observer using NUL-delimited Git path output, canonical SHA-256 dirty-state JSON, byte-safe regular-file and symlink image hashing, `.fid/` and ignored-untracked exclusion, tracked changes, mode changes, additions, deletions, and rename-as-delete-plus-add capture.
- Added reconciliation construction that binds observed repository identity, HEAD/tree/dirty facts, ancestry, classified paths, and a digest; exact adoption rejects omitted paths and stale digests.
- Added exactly `run.checkpoint` to the activated capability acceptance set.
- The durable `run.checkpoint` Event writer and its test cannot be completed under the current allowed source list; no Git mutation, CLI wiring, or projection-control edit was performed.

## Files created or modified

- `internal/project/repository.go`
- `internal/project/repository_test.go`
- `internal/workflow/reconciliation.go`
- `internal/workflow/reconciliation_test.go`
- `internal/workflow/activation.go`
- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-013: BLOCKED - `TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths` in `internal/project/repository_test.go` qualifies a clean base observation and proves a newline-bearing untracked path is observed; `run.checkpoint` cannot persist the required observation while `internal/journal/project_effects.go` retains `run.checkpointed` in `unavailableRecoveryEvents` outside task scope.
- VAL-014: BLOCKED - `TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions` pins staged, unstaged, mode-only, deleted, renamed, untracked, symlink, and exclusion digest behavior; `TestReconcileRepositoryRequiresExactDirtyAdoption` accepts the exact set and rejects omitted/stale adoption. It does not yet prove a durable checkpoint/recovery path.
- VAL-015: BLOCKED - `ReconcileRepository` rejects missing base, changed identity, unmerged paths, in-progress operations, and unattributable history with `RecoveryStateConflictError`, but the required fixture matrix and durable recovery rejection cannot be completed without the checkpoint projection path.
- VAL-016: BLOCKED - `ObserveRepository` invokes read-only Git inspection and tests use disposable repositories, but the required before/after refs, index, stash, and worktree snapshots across checkpoint/recovery success and failure are unfinished because checkpoint projection is unavailable.
- QA-14 through QA-17 and EVID-13/EVID-14: these are QA Lead checklist/evidence requirements, not additional Developer commands beyond task 4's exact gate. The partial focused tests above supply only preliminary EVID-13 material; no QA outcome is claimed.

## Tests run

- `go test ./internal/project -run '^TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Required first red attempt before product implementation: `undefined: ObserveRepository`.
- `go test ./internal/project -run '^TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/project\t0.055s` (terminated `exit_code: 0`)
- `go test ./internal/project -run '^TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red vector attempt: `unsupported Git image for "deleted.txt"`; causal parser repair followed.
- `go test ./internal/project -run '^TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Preserved second red attempt reached path-order assertion before the test expectation was corrected to raw-byte order.
- `go test ./internal/project -run '^TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Preserved third red attempt emitted the independently pinned canonical digest vector.
- `go test ./internal/project -run '^(TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths|TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/project\t0.222s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestReconcileRepositoryRequiresExactDirtyAdoption$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red attempt: `undefined: ReconcileRepository` and `undefined: ValidateReconciliationAdoption`.
- `go test ./internal/workflow -run '^TestReconcileRepositoryRequiresExactDirtyAdoption$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.051s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestActivatedCapabilitiesAcceptRunCheckpoint$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red attempt: `activated capabilities reject run.checkpoint`.
- `gofmt -w internal/project/repository.go internal/project/repository_test.go internal/workflow/reconciliation.go internal/workflow/reconciliation_test.go internal/workflow/activation.go && go test ./internal/workflow -run '^(TestReconcileRepositoryRequiresExactDirtyAdoption|TestActivatedCapabilitiesAcceptRunCheckpoint)$' -count=1 && go test ./internal/project -run '^(TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths|TestObserveRepositoryPinsCanonicalDirtyVectorsAndExclusions)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/project\t0.217s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
  - Immutable failed attempt: `TestEquivalentStopsShareAttentionButKeepDistinctFacts` printed `system clock is behind latest journal date 2026-09-06; refusing to append out of order`. Root cause is a retained fixed-date fixture introduced by `ac827959f8359473cbd6546ac6bea1ff1fd20751` (phase-22/task-2), not the task-4 observer. No repair was authorized because that retained assertion is outside the named task scope and is not invalidated by an accepted task-4 change.
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t61.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t8.529s` (terminated `exit_code: 0`)

## Deviations from spec

- The task remains incomplete. The accepted `run.checkpoint` contract requires a `run.checkpointed` Event, but the current projection's `unavailableRecoveryEvents` rejects it before applying the required `run_observations` and `run_runtime` effects. Editing that file is outside task 4's explicit source scope, so this report requests bounded continuation authority rather than changing an accepted file boundary or faking an Event.
- The required race/coverage command has one immutable failed attempt. Its fixed-date Stop fixture is outside task scope and not causally invalidated by this task, so no test, command, policy, or gate definition was changed.

## Concerns, technical debt, or notes for the Tech Lead

- The observer currently supplies the task-4 repository facts and exact adoption comparison, but checkpoint persistence and the full unsafe/nonmutation fixture matrix require the bounded projection admission repair before VAL-013 through VAL-016 can be completed.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a`

## Commit subject

- `feat(recovery): add repository reconciliation observation`
