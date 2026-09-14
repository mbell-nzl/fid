# Phase 22 - Task 3: Lease lifecycle and executor observations - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: Phase Runner/QA independent review-1 re-review of MF-3.
- Human attention required: no

## Summary of what was implemented

- Added explicit `claim.expire` mutation handling with strict `now < expires_at` lease validity and current-version fencing. Expired or released Claims with an unresolved prior Run now reject a fresh mutation Claim.
- Bound Run start to a project-global opaque executor identity and optional Envelope obligation. Reuse returns the exact `duplicate-entity` Run subject envelope; stale runtime mutations return the declared `run-version-conflict` envelope.
- Added trusted activated-harness-only `run.observe`, immutable Run executor/obligation read models, monotonic terminal liveness, and the task-owned `claim.expired` and `run.observed` journal effects.
- Added deterministic expiry, executor-binding, trusted-observer, and two-real-process concurrency regression coverage while retaining the earlier red test.
- QA rework for review-1 MF-3 adds one observable CLI-bound Claim-renewal-after-Run scenario. It proves the Run's original lineage/binding remains immutable through Claim renewal, rejects a stale Claim version with no persistence change, and accepts the current Claim version for release.
- The original task-3 implementation remains `7e419bfd8c233f933680d4a3b3957814c29dbd5e` (`feat(workflow): add lease expiry and executor observations`); this rework's test commit is recorded below.

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
- VAL-011: PASS - `TestClaimRenewalAfterRunPreservesLineageAndFencesClaimVersions` at `cmd/fid/generic_claim_test.go:270` starts a valid Claim-bound Run, renews the Claim from version 1 to 2, and uses `run get` to prove the Run ID, Issue, Envelope/revision, Claim ID/original version, actor, executor, obligation, repository, and predecessor/successor links remain unchanged. It then rejects a stale version-1 `claim release` with typed `claim-stale` and an unchanged projection/journal snapshot, while version-2 release succeeds and the same Run lineage still reads unchanged.
- VAL-012: PASS - `TestRunStartBindsOptionalObligationAndRejectsReusedExecutor` at `cmd/fid/generic_run_test.go:205` proves matching optional-obligation binding, rejects mismatches, and asserts the reused executor returns the exact duplicate Run subject envelope.

## Tests run

- `go test -race -cover ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.7% of statements` (terminated `exit_code: 0`)
  - Required QA-rework baseline run before editing `cmd/fid/generic_claim_test.go`.
- `go test ./cmd/fid -run '^TestClaimRenewalAfterRunPreservesLineageAndFencesClaimVersions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Focused red receipt: `generic_claim_test.go:270: MF-3 Claim-renewal Run-lineage scenario is not implemented`.
- `go test ./cmd/fid -run '^TestClaimRenewalAfterRunPreservesLineageAndFencesClaimVersions$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.364s` (terminated `exit_code: 0`)
  - The green scenario exercises the registered `envelope create`, `envelope activate`, `claim acquire`, `run start`, `claim renew`, `run get`, and `claim release` command boundary with JSON request and response envelopes.
- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Prior scope-gap red evidence retained: `generic_claim_test.go:250: expire Claim = main.commandExecutionResult{stdout:"", stderr:"", err:(*pflag.NotExistError)(0x814ac82f8c0)}, want explicit expiry success`.
- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - After adding expiry, the preserved test reached the next causal failure: `generic_claim_test.go:264: reacquire ... success, want typed claim conflict`.
- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.281s` (terminated `exit_code: 0`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red contract test: `codes_test.go:129: RegisteredCodes() length = 41, want 42`.
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.014s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red binding test: the start request received the typed invalid-argument executor rejection before runtime binding existed.
- `go test ./cmd/fid -run '^TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.256s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red observer test: `pflag.NotExistError` reported that `run observe` was not registered.
- `go test ./cmd/fid -run '^TestRunObserveRequiresTrustedBoundHarnessAndCurrentVersion$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.277s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.183s` (terminated `exit_code: 0`)
  - This is the real CLI entrypoint exercise: the test launches two concurrent external `fid` processes, not goroutines.
- `go test ./... -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Diagnostic attempt only: `TestGateRunRecordsExactNonSuccessTerminalOutcomes/signaled` observed its helper process exit successfully. The isolated subtest rerun passed (`ok  	github.com/mbell-nzl/fid/cmd/fid	0.350s`, terminated `exit_code: 0`); no test or gate definition was changed, and the required race suite below passed.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	39.316s	coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)				60.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.112s` (terminated `exit_code: 0`)

- QA-rework repeat of every exact task-3 command from `/home/mark/projects/fid`:
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)				61.1%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.280s` (terminated `exit_code: 0`)

## Deviations from spec

- The earlier recovery found that `internal/cli/codes.go` was absent from the original task file list even though the accepted Task 3 public interface uniquely requires its Run duplicate-entity subject and `run-version-conflict` definition. This ordinary task-scope causal repair added only those definitions and their direct regression test; no accepted command, assertion semantics, spec, decision, dependency, or unavailable boundary changed.
- The retained diagnostic `go test ./... -count=1` failure above was not a task quality gate and did not authorize a test or gate change. The mandatory byte-exact race suite subsequently passed.
- QA rework source: review-1 MF-3, trigger `review-1:qa-rework:task-3-claim-renewal-run-lineage`, attempt 1, repair authority QA review. The focused red marker was replaced only by the required task-owned scenario after its observed failure. The repair changes no product behavior, public interface, schema, journal contract, task/QA artifact, dependency, workflow control, golden fixture, or gate definition.

## Concerns, technical debt, or notes for the Tech Lead

- The independent review-1 re-review of MF-3 remains a Phase Runner/QA action. This task report records the repair evidence and does not claim a re-review result.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d6d140a12bb36e9926f009aba8fc49c1446b39c2`

## Commit subject

- `test(fid): cover claim renewal run lineage`

## Root report display correction — 2026-09-08

The root delegate decoded escaped tab characters in historical Go-output
receipts following the QA-84 preflight. The original Developer sessions remain
the source of those receipts, including every recorded outcome, duration,
coverage value and observed termination. This formatting correction records
no new execution or attestation. Original encoded text remains in Git history.
