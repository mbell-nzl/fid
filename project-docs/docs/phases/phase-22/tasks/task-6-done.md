# Phase 22 - Task 6: Recovery CLI and early process journey - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The task requires `run.checkpoint` and `run.inspect`, but predecessor task 4 did not deliver a checkpoint mutation or inspection read API in `internal/workflow`. The available recovery domain requires a reconciliation digest, and the canonical checkpoint Event has no producer. Implementing that domain outside the accepted task-6 file scope or fabricating it in the CLI would violate the contract.
- Recovery attempts: 0
- Next action: Runner must obtain a bounded same-task scope amendment for the missing canonical workflow checkpoint and inspection APIs, then resume task 6 from `TestP22EarlyRecoveryJourney`.
- Human attention required: no

## Summary of what was implemented

- Added the first red built-binary subprocess journey fixture. It initializes a real Git repository and fid project, activates a recovery-capable actor, creates one Issue and Envelope, acquires and starts a Run, records a terminated observation, and then reaches the absent `run.inspect` leaf.
- No production behavior was implemented because the required predecessor domain boundary is absent from the accepted task scope.

## Files created or modified

- `cmd/fid/integration/p22_recovery_journey_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-6-done.md`

## Documentation updated

- Added this durable recovery checkpoint report.

## Validation Contract Coverage

- VAL-022: BLOCKED - `TestP22EarlyRecoveryJourney` in `cmd/fid/integration/p22_recovery_journey_test.go:65` invokes the real built binary and reaches `fid run inspect`; it terminates with `exit_code: 1` because the current binary has no inspect leaf. The remaining checkpoint, expiry, recovery, and successor-inspection steps require a canonical missing predecessor API.
- VAL-023: BLOCKED - no recovery CLI leaf can receive or reject its closed wire input until the missing checkpoint/inspection domain APIs are supplied.
- VAL-024: BLOCKED - no task-6 process error mapping can be correctly implemented until the recovery leaves can invoke the canonical domain operations and their typed errors.

## Tests run

- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: the predecessor activation contract rejected `phase: 22` with `invalid-argument`; the fixture was corrected to its accepted pre-task-6 `phase: 21` shape.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: activation component digests were unsorted; the fixture now follows the accepted sorted order.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable fixture attempt: Envelope creation omitted the retained exact one-revision gate-evidence obligation and registration; the fixture now supplies both.
- `gofmt -w cmd/fid/integration/p22_recovery_journey_test.go && go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Intended red evidence: `fid --project /tmp/TestP22EarlyRecoveryJourney327508039/001 --json run inspect RUN-c620cb = err:exit status 2 stdout:"" stderr:""`. The focused test itself terminated `exit_code: 1` after the built binary rejected the absent leaf.

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- `golangci-lint run ./...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.
- `go tool cover -func=coverage.out` - BLOCKED - not started because its required coverage profile was not produced; no observed termination result.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not started because the missing canonical checkpoint/inspection domain API prevents a conforming task implementation; no observed termination result.

## Deviations from spec

- None. The task was not adapted; the missing predecessor API is recorded for bounded recovery.

## Concerns, technical debt, or notes for the Tech Lead

- `git log --all --oneline -- internal/workflow` identifies task 4's `ac5e6f9 feat(recovery): add repository reconciliation observation` and task 4's report claims a checkpoint projection test, but the checked-out `internal/workflow/` tree contains no `RunCheckpoint`, checkpoint Event producer, or inspection read API. Task 6 cannot expose canonical JSON for that behavior without an explicit scope amendment.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
