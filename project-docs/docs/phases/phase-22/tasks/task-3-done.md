# Phase 22 - Task 3: Lease lifecycle and executor observations - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added explicit `claim.expire` mutation handling with strict `now < expires_at` lease validity and current-version fencing. Expired or released Claims with an unresolved prior Run now reject a fresh mutation Claim.
- Bound Run start to a project-global opaque executor identity and optional Envelope obligation. Reuse returns the exact `duplicate-entity` Run subject envelope; stale runtime mutations return the declared `run-version-conflict` envelope.
- Added trusted activated-harness-only `run.observe`, immutable Run executor/obligation read models, monotonic terminal liveness, and the task-owned `claim.expired` and `run.observed` journal effects.
- Added deterministic expiry, executor-binding, trusted-observer, and two-real-process concurrency regression coverage while retaining the earlier red test.

## Files created or modified

- `cmd/fid/claim.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/generic_claim_test.go`
- `cmd/fid/generic_run_test.go`
- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/run.go`
- `cmd/fid/task12_execution_transition_test.go`
- `internal/claims/generic.go`
- `internal/claims/generic_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/journal/project_effects.go`
- `internal/workflow/activation.go`
- `internal/workflow/run.go`
- `internal/workflow/run_observation.go`
- `internal/workflow/run_read.go`
- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`

## Validation Contract Coverage

- VAL-008: PASS - `TestGenericClaimExpiresAtEqualityAndNeverBefore` at `internal/claims/generic_test.go:8` proves valid-before, expired-at-equality, and expired-after boundaries. `TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired` at `cmd/fid/generic_claim_test.go:229` proves explicit expiry and the unresolved-Run reacquisition fence.
- VAL-009: PASS - `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` at `cmd/fid/integration/p21_conformance_test.go:74` starts two actual `fid` processes, asserts exactly one Claim winner, a typed `claim-conflict` loser, and one added Event.
- VAL-010: PASS - `TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion` at `cmd/fid/generic_run_test.go:155` accepts the bound activated-harness observation and rejects an ordinary Builder and a forged executor without recording an observation.
- VAL-011: PASS - `TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion` at `cmd/fid/generic_run_test.go:155` rejects a stale Run version with `run-version-conflict` and no new Event; the same test confirms a Run keeps its bound executor and obligation across observation versions.
- VAL-012: PASS - `TestRunStartBindsOptionalObligationAndRejectsReusedExecutor` at `cmd/fid/generic_run_test.go:205` proves matching optional-obligation binding, rejects mismatches, and asserts the reused executor returns the exact duplicate Run subject envelope.

## Tests run

- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Prior scope-gap red evidence retained: `generic_claim_test.go:250: expire Claim = main.commandExecutionResult{stdout:"", stderr:"", err:(*pflag.NotExistError)(0x814ac82f8c0)}, want explicit expiry success`.
- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - After adding expiry, the preserved test reached the next causal failure: `generic_claim_test.go:264: reacquire ... success, want typed claim conflict`.
- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.281s` (terminated `exit_code: 0`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red contract test: `codes_test.go:129: RegisteredCodes() length = 41, want 42`.
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t0.014s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red binding test: the start request received the typed invalid-argument executor rejection before runtime binding existed.
- `go test ./cmd/fid -run '^TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.256s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red observer test: `pflag.NotExistError` reported that `run observe` was not registered.
- `go test ./cmd/fid -run '^TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.277s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t3.183s` (terminated `exit_code: 0`)
  - This is the real CLI entrypoint exercise: the test launches two concurrent external `fid` processes, not goroutines.
- `go test ./... -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Diagnostic attempt only: `TestGateRunRecordsExactNonSuccessTerminalOutcomes/signaled` observed its helper process exit successfully. The isolated subtest rerun passed (`ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.350s`, terminated `exit_code: 0`); no test or gate definition was changed, and the required race suite below passed.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t39.316s\tcoverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t60.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t8.112s` (terminated `exit_code: 0`)

## Deviations from spec

- The earlier recovery found that `internal/cli/codes.go` was absent from the original task file list even though the accepted Task 3 public interface uniquely requires its Run duplicate-entity subject and `run-version-conflict` definition. This ordinary task-scope causal repair added only those definitions and their direct regression test; no accepted command, assertion semantics, spec, decision, dependency, or unavailable boundary changed.
- The retained diagnostic `go test ./... -count=1` failure above was not a task quality gate and did not authorize a test or gate change. The mandatory byte-exact race suite subsequently passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `7e419bfd8c233f933680d4a3b3957814c29dbd5e`

## Commit subject

- `feat(workflow): add lease expiry and executor observations`
