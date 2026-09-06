# Phase 22 — Task 5: Atomic linked Run recovery — Completion Report

## Workflow outcome

- Task outcome: COMPLETE

## Summary of what was implemented

- Added the internal `RecoverRun` mutation. It reconciles the repository before it appends one idempotent `run.recovered` Event that terminalizes the predecessor, disposes or releases its Claim as applicable, creates one linked successor Claim and Run, and records the recovery link.
- Preserved recovery identity across a lost response: the same idempotency key returns the original successor, while a changed request under that key is rejected.
- Added ordinary recovery checks for a trusted terminated observation and a non-active predecessor Claim. Added authorized force recovery checks for current Human authority, exact recovery and force-recovery grants, a live observation, and the current envelope revision and grant ordinal. Unknown liveness is rejected for both paths.
- Exposed predecessor and successor links in Run reads and retained interrupted Runs in `ListRuns`. The predecessor remains terminal and prior Run writes are fenced.
- Applied DEC-192 by admitting only `run.recovered` from the recovery-event fence. The remaining recovery Events remain fenced, and the retained Task 1 schema/replay assertion now exercises `run.budget-observed` while preserving its corrupt Stop-code branch.

## Files created or modified

- `internal/workflow/recovery.go`
- `internal/workflow/recovery_test.go`
- `internal/workflow/activation.go`
- `internal/workflow/run.go`
- `internal/workflow/run_read.go`
- `internal/journal/project_effects.go`
- `internal/journal/checkpoint_projection_test.go`
- `internal/journal/recovery_schema_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-5-done.md`

## Documentation updated

- Added this completion report with validation, red/green, convergence, and gate evidence.

## Validation Contract Coverage

- `VAL-017` — `TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse` in `internal/workflow/recovery_test.go` proves a distinct linked successor Claim and Run, an interrupted terminal predecessor, and that the predecessor remains visible through `ListRuns`.
- `VAL-018` — `TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse` proves same-key replay returns the original recovery and rejects a changed request under the same key. `TestRecoveredEventProjectsIdenticallyLiveAndReplay` in `internal/journal/checkpoint_projection_test.go` proves `run.recovered` has matching live-apply and replay projection results. The recovery mutation appends its predecessor transition, Claim disposition, successor rows, link, and reconciliation record in one Event.
- `VAL-019` — `TestRecoverRunRejectsUnknownLivenessForOrdinaryAndForce` in `internal/workflow/recovery_test.go` proves unknown liveness is rejected on both recovery paths. `TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse` supplies the accepted terminated-observation ordinary recovery case.
- `VAL-020` — `TestRecoverRunForceRequiresCurrentHumanAuthorityAndFencesPriorRun` in `internal/workflow/recovery_test.go` rejects non-Human force authority, accepts the current Human actor with exact current grants and envelope revision, retains the live observation, and proves a later write to the interrupted predecessor is rejected.
- `VAL-021` — `TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse` snapshots and compares the Issue status and immutable Envelope execution contract before and after recovery.
- `DEC-192` / `VAL-003` convergence — removed exactly `run.recovered` from `unavailableRecoveryEvents` in `internal/journal/project_effects.go`. The displaced assertion in `internal/journal/recovery_schema_test.go` was introduced by `72c5974f6698fa08a4fe1a24099b3a93b3882f1d` (`feat(recovery): add schema and replay foundation`, `Refs: phase-22/task-1`). Task 5's accepted `run.recovered` admission displaced that Event choice; the retained assertion now uses still-fenced `run.budget-observed`, keeps live/replay rejection parity, and retains `assertInvalidStopCodeRejectsLiveAndReplay` for the corrupt Stop-code branch. `TestRecoveredEventProjectsIdenticallyLiveAndReplay` proves the newly admitted Event's live/replay parity.

## Tests run

- `go test ./internal/workflow -run '^TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse$' -count=1` — expected RED — `FAIL` (terminated `exit_code: 1`). The test initially failed because `RecoverRun` and `RunRecoveryRequest` were undefined.
- `go test ./internal/workflow -run '^TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse$' -count=1` — PASS — `ok` (terminated `exit_code: 0`).
- `go test ./internal/workflow -run '^TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse$' -count=1` — expected RED — `FAIL` (terminated `exit_code: 1`). The interrupted predecessor was not accepted by `ListRuns`.
- `go test ./internal/workflow -run '^TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse$' -count=1` — PASS — `ok` (terminated `exit_code: 0`).
- `go test ./internal/workflow -run '^(TestRecoverRunRejectsUnknownLivenessForOrdinaryAndForce|TestRecoverRunForceRequiresCurrentHumanAuthorityAndFencesPriorRun)$' -count=1` — PASS — `ok` (terminated `exit_code: 0`).
- `go test ./internal/journal -run '^(TestRecoveredEventProjectsIdenticallyLiveAndReplay|TestRecoverySchemaFreshInitAndReplay|TestRemainingRecoveryEventsRemainFenced)$' -count=1` — PASS — `ok` (terminated `exit_code: 0`).
- `go test ./internal/workflow ./internal/journal` — PASS — `ok` (terminated `exit_code: 0`).
- `go test ./...` — PASS — `ok` (terminated `exit_code: 0`).
- `go build -o /tmp/fid ./cmd/fid && /tmp/fid run --help` — PASS — `Use "fid run [command] --help" for more information about a command.` (terminated `exit_code: 0`). This real-process exercise confirms the current public Run command surface remains limited to the existing commands; Task 6 owns recovery CLI wiring.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` — PASS — `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` — PASS — `ok` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` — PASS — `62.7%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` — PASS — `ok` (terminated `exit_code: 0`).

## Deviations from spec

- None. DEC-192's `run.recovered` fence amendment was pre-authorized for this task.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
