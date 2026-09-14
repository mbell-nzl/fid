# Phase 22 - Task 4: Read-only repository reconciliation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: FQ-03 repository safety repair and all six unchanged gates are complete; independent final re-review remains.
- Recovery attempts: 3
- Next action: Root continues Task5 FQ-04; final-review-2 follows all seven owners.
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

- `e4fa84cee9463e8d0159c570ce4a80f18d2d492d`

## Commit subject

- `fix(recovery): reject unsafe and changing repository captures`

## Root report display correction — 2026-09-08

The root delegate decoded escaped tab characters in historical Go-output
receipts following the QA-84 preflight. The original Developer sessions remain
the source of those receipts, including every recorded outcome, duration,
coverage value and observed termination. This formatting correction records
no new execution or attestation. Original encoded text remains in Git history.

## Final-QA FQ-03 same-task continuation — 2026-09-09

Authority: review-7, final-review.md FQ-03, task-4.md and contracts.md §3/3.1.
Root is the sole writer. Original receipts and prior task/DEC-192/193 history
remain unchanged. This appendix records the current repair; independent final
re-review remains required after all seven owners.

### Implementation and direct evidence

- VAL-015: the observer now rejects uppercase skip-worktree as well as
  assume-unchanged flags and all unmerged stages. It enumerates the union of
  index-versus-HEAD and worktree-versus-index paths, preserving staged changes
  even when the worktree matches HEAD; changed gitlinks therefore cannot vanish
  from observation. Two complete captures must agree, or the operation rejects.
- Reconciliation verifies that the recorded base still exists, including when
  HEAD is unchanged. No digest grammar, path classification, adoption policy or
  imported destination authority changed.
- Existing public recovery rendering maps wrapped ErrUnsafeRepository to the
  accepted recovery-state-conflict/conflict/exit 4 envelope, including the
  original subject, observation digest, paths and reasons. No new wire fields.
- TestObserveRepositoryRejectsUnsafeImages proves assume-unchanged, skip-worktree,
  unmerged, missing HEAD, changed gitlink, unreadable file, FIFO and escaping
  parent rejection. TestObserveRepositoryRejectsChangeDuringCapture uses real Git
  NUL output and a test-local executable wrapper that creates a source file just
  after untracked enumeration; it asserts rejection and the external write.
  Invalid UTF-8, NUL, absolute/traversal paths reject at the path grammar.
- The retained fixed staged/unstaged/mode/deleted/rename/untracked/symlink and
  exclusion vectors remain unchanged. A new independently encoded staged-change/
  restored-worktree vector pins dirty digest
  sha256:9926607bb8970a8199f50aadbeb618aa94def7c761d00fc86cf3ddf7b2f77b48.
- VAL-016: TestP22RecoveryPreservesRepositoryOnSuccessAndUnsafeRejection executes
  the actual built fid binary across 21 cases: clean, exact dirty in-scope,
  staged/restored, both tracked flags, unmerged, missing HEAD, unrelated history,
  foreign worktree, gitlink, unreadable, FIFO, escaping parent, protected adoption,
  and all seven merge/rebase/cherry-pick/revert/bisect markers. Each retains a real
  stash and additional ref. Complete source bytes/modes/symlink payloads and all
  Git files/directories (HEAD, complete refs, index, stash, reflogs and objects)
  are compared before/after the real operation. Only .fid operational state is
  excluded; caller input is outside the workspace. Unreadable-file mode is
  compared, then permissions are restored by the test and exact content checked.
  Success appends exactly one Event; each rejection has the exact typed conflict
  and appends none. No recovery source/Git mutation was observed.
- TestReconcileRepositoryRejectsMissingBaseEvenWithUnchangedHead covers null and
  nonexistent recorded bases with an otherwise real unchanged Git observation.
  Force/caller scope authority remains Task5's separate FQ-04 repair.

### Provenance and mandatory deviations

- The public recovery renderer originated in
  69941fc1b6262fdd7efaaeecfe72760390dbd511 (Phase22/Task6), with its non-null paths
  convergence at f3be616f0587881ee3f5c3716ffbf6167e3f2a20 (Phase22/Task6).
  FQ-03 explicitly assigns this causal mapping repair; Task6's report is unchanged.
- Observer red ended `FAIL` (exit_code: 1): skip-worktree, unmerged, gitlink and
  capture-time changes were accepted. Root repaired the concrete predicates and
  enumeration; the byte-identical focused command then passed at 1.333s.
- Actual CLI red ended `FAIL` (exit_code: 1): unsafe errors produced
  invalid-argument/invalid-input/exit 2 instead of the accepted conflict/exit 4.
  The existing renderer mapping repaired this; the byte-identical command passed
  at 24.931s. Its first escaping-parent fixture also lacked a parent directory;
  creating that directory fixed construction without changing the assertion.
- Missing-base red ended `FAIL` (exit_code: 1) because a nonexistent base with
  unchanged HEAD was accepted. The current-base existence check repaired it;
  the byte-identical command passed at 1.066s.
- The initial observer test compile command failed with final token `go mod tidy`
  (exit_code: 1) because the test imported testify, which is not in the actual
  module despite the high-level stack description. The test now uses the existing
  standard-library style. No module, dependency or license policy changed.
- No accepted contract, Event, schema, dependency, gate command, threshold,
  protected control or product policy changed. All failed attempts remain below.

### Continuation command receipts

Working directory: `/home/mark/projects/fid`. Root captured these tokens and real
termination results in this recording session. Focused tests supplement the six
unchanged final commands, which retain their specified order.


- `go test -race ./internal/project -run '^(TestObserveRepositoryRejects|TestRepositoryPathGrammar)' -count=1` - FAIL - `	go mod tidy` (terminated `exit_code: 1`)
  - unsafe observer behavioral red.
- `go test -race ./internal/project -run '^(TestObserveRepositoryRejects|TestRepositoryPathGrammar)' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - unsafe observer red stdlib fixture.
- `go test -race ./internal/project -run '^(TestObserveRepositoryRejects|TestRepositoryPathGrammar)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/project	1.333s` (terminated `exit_code: 0`)
  - unsafe observer green.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22RecoveryPreservesRepositoryOnSuccessAndUnsafeRejection$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - actual recovery safety and typed mapping red.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22RecoveryPreservesRepositoryOnSuccessAndUnsafeRejection$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	24.931s` (terminated `exit_code: 0`)
  - actual recovery safety typed mapping green.
- `go test -race ./internal/workflow -run '^TestReconcileRepositoryRejectsMissingBaseEvenWithUnchangedHead$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - missing base behavioral red.
- `go test -race ./internal/workflow -run '^TestReconcileRepositoryRejectsMissingBaseEvenWithUnchangedHead$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	1.066s` (terminated `exit_code: 0`)
  - missing base green.
- `go test -race ./internal/project` - PASS - `ok  	github.com/mbell-nzl/fid/internal/project	1.827s` (terminated `exit_code: 0`)
  - complete project canonical and safety regression.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - required Task4 post-gate.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - required Task4 post-gate.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	236.610s	coverage: 60.7% of statements` (terminated `exit_code: 0`)
  - required full race coverage.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.1%` (terminated `exit_code: 0`)
  - required Task4 post-gate.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	128.972s` (terminated `exit_code: 0`)
  - required Task4 post-gate.
- `go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=5` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	9.692s` (terminated `exit_code: 0`)
  - required repeated Stop regression.

Prior primary binding: 59c62b9e1d82c5f698d1e173908446d35b64f539, `fix(recovery): admit checkpoint projection`. Observer origin: ac5e6f9bec80e06bfe5b7b33e3611a894d16038a, Phase22/Task4. Final ordinary coverage65.1%; all six final commands terminated exit0.
