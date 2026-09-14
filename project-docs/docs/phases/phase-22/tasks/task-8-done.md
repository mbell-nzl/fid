# Phase 22 - Task 8: Interrupted Attempts and derived Run closure - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: FQ-06 repair complete; all eight unchanged final commands passed.
- Recovery attempts: 5
- Next action: Task10 FQ-07, then Task1 FQ-01 and independent final-review-2.
- Human attention required: no
- Historical blocker resolved: DEC-195 resolved ESC-22-004 / TRG-22-009 by assigning the two exact `execution_receipts` baseline replacements, with task-1 provenance, to this same-task continuation. The preserved prior BLOCKED evidence below records the formerly missing authority.

## Summary of what was implemented

- Persisted the gate idempotency claim with the pre-spawn Attempt, and replace that pending result only after a fid-observed Receipt commits. A lost terminal observation now preserves one Attempt without a Receipt and the same key cannot spawn again.
- Added the observed `cancelled` terminal outcome for a caller-cancelled process.
- Repaired the inherited R-14 subprocess fixtures: the signal child blocks after self-signalling, and the timeout child writes a readiness marker then blocks under a five-second real gate deadline.
- Added `run.close` and `run.abandon` domain operations. Close accepts only a trusted terminated observation, rejects outstanding Attempts, derives `succeeded` only when the shared Obligation evaluator has no reasons, and otherwise records `failed`. Abandon records a distinct terminal Run, releases its active Claim with `abandoned` provenance, accepts ordinary terminal proof or exact Human force authority over a live observation, and creates no Evidence or successor Run.
- Admitted the two task-8 actions while retaining `run.finish`, removed exactly the `run.closed` and `run.abandoned` DEC-192 availability fence keys, and added live/replay parity coverage for both Events.

## Files created or modified

- `internal/evidence/execution.go`
- `internal/evidence/live.go`
- `internal/evidence/live_test.go`
- `internal/workflow/run_close.go`
- `internal/workflow/run_close_test.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `internal/journal/checkpoint_projection_test.go`
- `internal/journal/recovery_schema_test.go`
- `cmd/fid/gate_run_test.go`
- `cmd/fid/receipt_test.go`
- `internal/db/operational_kernel_migration_test.go`
- `internal/evidence/receipt_read.go`
- `migrations/0013_p21_operational_kernel.up.sql`
- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`

## Validation Contract Coverage

- VAL-031: PASS - `TestRunLiveRecordsCancellationAsCancelled` proves the fid-observed `cancelled` process outcome; `TestMigration0013CancelledReceiptTerminalConstraint` admits cancelled Receipts carrying an exit code, signal, or neither and rejects `timed_out = 1` and a simultaneous exit code/signal; `TestReceiptReadFailuresAreTypedAndClosed` exercises the real `receipt list --outcome cancelled` CLI acceptance while retaining unknown-outcome rejection. The repeated R-14 gate command passed fifteen times.
- VAL-032: PASS - `TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn` proves one durable Attempt, no Receipt, and no same-key second spawn after terminal-observation loss.
- VAL-033: PASS - `TestCloseRunDerivesTerminalStateFromTrustedObservation` proves trusted terminal observation, shared-evaluator-derived success, versioning, and idempotent replay; the complete ordinary and tagged race gates passed.
- VAL-034: PASS - `TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority` proves ordinary intentional abandonment, Claim release provenance, no Evidence, and unknown-liveness rejection; `TestRunClosureEventsProjectIdenticallyLiveAndReplay` proves the `run.abandoned` Event's live/replay parity.

## Tests run

- `go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: same-key retry created `2/0` Attempts/Receipts.
- `gofmt -w internal/evidence/execution.go cmd/fid/gate_run_test.go && go test -race -run '^TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.507s` (terminated `exit_code: 0`)
- `go test -race -run '^(TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn|TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame|TestGateRunRecordsExactNonSuccessTerminalOutcomes)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	12.909s` (terminated `exit_code: 0`)
- `gofmt -w internal/evidence/live_test.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); red evidence: cancellation was incorrectly observed as `timed-out`.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`); the process signal overwrote the newly derived `cancelled` outcome.
- `gofmt -w internal/evidence/live.go && go test -race -run '^TestRunLiveRecordsCancellationAsCancelled$' ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.031s` (terminated `exit_code: 0`)
- `gofmt -w cmd/fid/gate_run_test.go && go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	17.140s` (terminated `exit_code: 0`)

- `go test ./internal/workflow -run '^TestActivationAdmitsRunCloseAndAbandonActions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: the accepted action admission rejected `run.abandon` and `run.close`.
- `gofmt -w internal/workflow/activation.go internal/workflow/run_close_test.go && go test ./internal/workflow -run '^TestActivationAdmitsRunCloseAndAbandonActions$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.004s` (terminated `exit_code: 0`)
- `gofmt -w internal/workflow/run_close_test.go && go test ./internal/workflow -run '^(TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: `RunCloseRequest`, `CloseRun`, `RunAbandonRequest`, and `AbandonRun` were undefined.
- `go test ./internal/journal -run '^TestRunClosureEventsProjectIdenticallyLiveAndReplay$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red evidence: both accepted closure Events were rejected by their DEC-192 availability fences.
- `gofmt -w internal/journal/project_effects.go internal/journal/checkpoint_projection_test.go internal/workflow/run_close.go internal/workflow/run_close_test.go && go test ./internal/journal ./internal/workflow -run '^(TestRunClosureEventsProjectIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestActivationAdmitsRunCloseAndAbandonActions|TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - First green attempt preserved: the Claim provenance `run-abandoned` violated the accepted DDL's closed `release_reason` domain.
- `gofmt -w internal/journal/project_effects.go internal/journal/checkpoint_projection_test.go internal/workflow/run_close.go internal/workflow/run_close_test.go && go test ./internal/journal ./internal/workflow -run '^(TestRunClosureEventsProjectIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestActivationAdmitsRunCloseAndAbandonActions|TestCloseRunDerivesTerminalStateFromTrustedObservation|TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.393s` (terminated `exit_code: 0`)
- `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' - FAIL - `FAIL` (terminated `exit_code: 1`)
  - DEC-192 made its retained `run.closed` rejection fixture projectable, reaching `run_recoveries.prior_run_id` validation.
- `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/journal 0.104s` (terminated `exit_code: 0`)
- `go test ./internal/workflow ./internal/journal ./internal/evidence ./cmd/fid -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Retained `TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError` expected no idempotency record after task 8's durable pre-spawn Attempt.
- `gofmt -w cmd/fid/gate_run_test.go && go test ./cmd/fid -run '^TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 0.260s` (terminated `exit_code: 0`)
- `go test ./internal/workflow ./internal/journal ./internal/evidence ./cmd/fid -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 14.569s` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-8-done.md project-docs/docs/phases/phase-22/tasks/task-8.md` - FAIL - `FAIL report-lint found 2 failure(s)` (terminated `exit_code: 2`)
  - Root cause: `Recovery attempts` was not a typed integer, and one new `go test` receipt was not an output token. Repair: retain attempt count `3` alone and record `ok` as the exact short terminal token.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-8-done.md project-docs/docs/phases/phase-22/tasks/task-8.md` - FAIL - `FAIL report-lint found 3 failure(s)` (terminated `exit_code: 3`)
  - Root cause: bare `ok` is too short for the report control's `go test` receipt grammar. Repair: restore the exact package output lines from the observed terminal results.
- DEC-195 continuation red: `go test ./internal/db ./cmd/fid -run '^(TestMigration0013CancelledReceiptTerminalConstraint|TestMigration0013OperationalKernelInventoryAndRoundTrip|TestReceiptReadFailuresAreTypedAndClosed)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). `execution_receipts` rejected a valid cancelled Receipt and `receipt list --outcome cancelled` returned `invalid-argument`.
- DEC-195 continuation green: `gofmt -w internal/db/operational_kernel_migration_test.go cmd/fid/receipt_test.go internal/evidence/receipt_read.go && go test ./internal/db ./cmd/fid -run '^(TestMigration0013CancelledReceiptTerminalConstraint|TestMigration0013OperationalKernelInventoryAndRoundTrip|TestReceiptReadFailuresAreTypedAndClosed)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.091s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `golangci-lint run ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go tool cover -func=coverage.out` - BLOCKED - not run because its required coverage input is blocked by the protected cancelled-Receipt contradiction.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `go test -bench . -benchmem ./...` - BLOCKED - not run because the protected cancelled-Receipt contradiction prevents task completion.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	77.478s	coverage: 48.4% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				62.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.709s` (terminated `exit_code: 0`)
- `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	192.448s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	5.097s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	5.093s` (terminated `exit_code: 0`)

## Deviations from spec

- DEC-195 accepts the two exact 0013 baseline replacements for `execution_receipts`: the closed `terminal_status` list gains `cancelled` and its four-clause terminal CHECK becomes the schema.md section-1 five-clause block. `git blame` attributes both prior migration lines and the retained golden DDL string to the task-1 baseline (`72c5974f6698fa08a4fe1a24099b3a93b3882f1d`); DEC-195 transferred their execution to task 8. No other 0013 byte, migration, DDL surface, wire shape, error code, Event semantic, or gate definition changed.
- The first closure implementation attempt used Claim `release_reason: "run-abandoned"`; the accepted DDL rejected it. The causal repair uses the existing closed `"abandoned"` value and preserves explicit Claim-release provenance.
- Convergence provenance: `internal/journal/recovery_schema_test.go:94-107` was introduced by `72c5974f6698fa08a4fe1a24099b3a93b3882f1d` (phase-22/task-1) and was moved to its `run.closed` fence subject by task 7's accepted retained-row convergence. DEC-192 and task 8's accepted `run.closed` producer invalidate only that unavailable-Event identity. The subject-preserving repair keeps live/replay rejection coverage by using still-fenced `activation.resource-observed`; executing command: `go test ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1`.
- Convergence provenance: `cmd/fid/gate_run_test.go:180-209` was introduced by `da6970b4ac5562e20c62f6c72ed8568c10d36f08` (phase-21/task-15). Task 8 implementation commit `bead7efb134f80726022f1e38d7b6749fe550073` persistently fences the pre-spawn Attempt and its pending idempotency result, invalidating only the old `0/0` Receipt/idempotency count. The subject-preserving repair retains the frame-only `fid-error` and no-Receipt assertions, and adds the required `1/0/1` Attempt/Receipt/idempotency proof; executing command: `go test ./cmd/fid -run '^TestGateRunReceiptAppendFailureEmitsFrameOnlyFidError$' -count=1`.

## Concerns, technical debt, or notes for the Tech Lead

- Retained R-14 fixture provenance: `cmd/fid/gate_run_test.go` signal helper and timeout row were introduced by `b69499116e70d4c88423a30daa2010a4ff221f14` (phase-21/task-15); the timeout extension was last touched by `76ad781727b9e8545e069c4ac0e9f926012e6af9` (phase-21/task-21). DEC-186 and task 8’s cancellation/observed-terminal acceptance invalidate only their fixture mechanisms. The repaired subject remains exact signal and timeout terminal Receipt assertions, executed by `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' ./cmd/fid`.
- DEC-195 resolution provenance: task-1 owns the accepted 0013 baseline and its golden DDL inventory; DEC-195 assigned only the two omitted cancellation replacements and corresponding golden convergence to task 8. `TestMigration0013CancelledReceiptTerminalConstraint` is the executing proof.
- Recovery accounting: attempt 1 made interrupted-attempt and R-14 material progress; attempt 2 was report-only; attempt 3 made closure, action-admission, Event-parity, and retained-fixture convergence progress. The bound no-progress counter remains 0/2.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c0003e0c5df7f2dffc46e05d04ddd7bdfa53f307`

## Commit subject

- `fix(workflow): verify terminal repository and exact abandon authority`

## Root report display correction — 2026-09-08

The root delegate decoded escaped tab characters in historical Go-output
receipts following the QA-84 preflight. The original Developer sessions remain
the source of those receipts, including every recorded outcome, duration,
coverage value and observed termination. This formatting correction records
no new execution or attestation. Original encoded text remains in Git history.


## Final-QA FQ-06 same-task continuation — 2026-09-09

Authority: review-7, final-review.md FQ-06, task-8.md, contracts.md §§3/7 and
DEC-208. Root is the sole writer. All original failures, BLOCKED receipts and
DEC-195 history remain above. Prior primary binding:
c5dc59e4bb9cf9db20a4aa157daeef68d132bb60,
`fix(evidence): persist cancelled receipts`.

### Implementation and validation evidence

- VAL-033: CloseRun validates the asserted terminal repository against the
  actual current observer while holding the write lock, before shared Obligation
  evaluation. Base remains the historical comparison origin; a non-null base
  must exist, and nullable base remains supported. Every other identity/dirty
  field must match observation. Caller assertions are never refreshed into
  current proof. AbandonRun uses the same current terminal assertion check.
- VAL-034: exact force abandonment additionally checks the referenced grant's
  Issue scope and required force capability. Another valid grant cannot make an
  unrelated referenced ordinal sufficient. Existing Human type, current Envelope,
  Actor, live observation, grant action and reference checks remain.
- New run_close_authority_test.go executes 39 abandonment cases across ordinary
  and Human force, plus ten direct CloseRun observation/version/authority cases.
  It covers valid controls, missing/nonexistent/wrong-Run/wrong-executor/wrong-kind
  observations, unknown/live states, stale Run/Claim versions, absent grants,
  closed Envelope, Actor mismatch, expired Claim, non-Human force, capability,
  reference/Envelope/revision/ordinal and four separately referenced-grant defects.
  Every rejection compares all projection tables and journal bytes unchanged.
- Successful abandonment checks distinct abandoned status, current version,
  Claim release reason/version/time, immutable Claim identity/scope, durable
  Actor/Run Event provenance, no successor identity, no new Receipt/Evidence or
  changed other context, and exact public terminal replay without new facts.
- New p22_close_proof_test.go executes 13 actual CLI cases against a required
  Obligation. Qualifying fid execution plus linked Evidence succeeds; missing,
  stale and failed proof derive failure. A real Receipt append fault supplies
  an outstanding Attempt, which close rejects. Live/unknown, missing/null
  observation, stale version, forged old repository and missing base reject
  unchanged. A null base is handled without panic or manufactured proof.
  Stale proof is tested by changing real base.txt bytes; the forged old-target
  test proves the caller cannot choose an old target to make that proof current.
  Accepted close changes only terminal Run/runtime and persistence bookkeeping;
  it leaves the Claim, Issue, Envelope and all proof unchanged. Exact CLI replay
  retains the terminal Run and Event once.
- New p22_abandon_authority_test.go exercises 14 actual CLI cases. Ordinary
  terminated and exact Human known-live controls succeed. Both modes reject
  unknown, missing/null/nonexistent observation and forged target; ordinary live
  and force without a reference reject. It proves the accepted typed failures,
  no-effect rejection, distinct release provenance, no new proof/successor and
  byte-equivalent public replay apart from replayed=true.
- DEC-208 allocates only cmd/fid/run.go's close/abandon domain-error returns to
  the existing recovery renderer. This restores accepted liveness/observation/
  repository error codes, categories and exits without editing the registry or
  renderer. Introducing Task9 commit:
  a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460, verified with git log -S for both
  domain calls. Task9's report remains COMPLETE and unchanged.

### Mandatory deviations and failed attempts

- Authority red ended `FAIL` (exit_code: 1), 111.545s. Wrong-Issue and missing
  force capability on the specifically referenced grant were accepted. The
  exact-grant checks repair these causal defects. The identical command passed
  at 112.277s. Two initial valid replay assertions compared ScopeJSON, an internal
  field explicitly excluded from serialized Claim results. The comparison now
  checks the complete public Claim JSON and preserves the separate immutable
  database Claim scope assertion; no product field or wire contract changed.
- Actual close red ended `FAIL` (exit_code: 1), 19.394s. Reusing the old repository
  assertion after a real change produced succeeded; live/unknown domain failures
  became generic invalid-argument/2. Actual observation validation plus DEC-208's
  existing renderer delegation repaired them. The identical command passed at
  25.438s, including two added nullable/missing-base controls.
- Actual public abandonment matrix passed at 17.355s (exit_code: 0).
- No cancellation, signal, timeout, Receipt outcome/DDL, Event grammar, dependency,
  protected gate/control/threshold, validation assertion or product policy changed.
  All synchronized executor cases and both required benchmark commands remain.
  Independent final-review-2 judges sufficiency after all owners; Human judgments
  remain blank.

### Continuation command receipts

Working directory: `/home/mark/projects/fid`. Root captured each following
command token and nested observed termination in this recording session.

Retained caller convergence: TestRunStartBindsOptionalObligationAndRejectsReusedExecutor initially failed `FAIL` (exit_code: 1),47.663s, because its synthetic non-Git repository could not supply the now-required current terminal observation. Its original Task3 introducing commit is 7e419bfd8c233f933680d4a3b3957814c29dbd5e. Replacing only its setup with observedEvidenceEnvelopeForTest(t,nil) preserves optional Obligation binding, duplicate executor rejection and exact null subject/details assertions. The byte-identical four-caller command passed51.094s (exit_code: 0). New primary implementation c0003e0, fix(workflow): verify terminal repository and exact abandon authority.

Recovery accounting: four original attempts plus FQ-06 attempt1; consecutive no-progress remains0/2. All eight final commands passed. Aggregate coverage65.4%; first/second 100k rebuild26718/27316ms. Required cancellation cases passed in the ordinary race suite; fifteen-repeat command preserves failed/spawn-failed/timed-out/signaled cases.

- `go test -race ./internal/workflow -run '^(TestAbandonRunRequiresExactCurrentAuthorityAndObservation|TestCloseRunRejectsUnqualifiedObservationVersionAndAuthority)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 1: close and abandonment authority matrix red.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22CloseRunUsesObservedRepositoryAndQualifyingProof$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 2: actual close proof and repository assertion red.
- `go test -race ./internal/workflow -run '^(TestAbandonRunRequiresExactCurrentAuthorityAndObservation|TestCloseRunRejectsUnqualifiedObservationVersionAndAuthority)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	112.277s` (terminated `exit_code: 0`). Attempt 3: close and abandonment authority matrix green.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22CloseRunUsesObservedRepositoryAndQualifyingProof$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	25.438s` (terminated `exit_code: 0`). Attempt 4: actual close proof and repository assertion green.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22AbandonUsesExactHumanForceAndTypedObservationFailures$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	17.355s` (terminated `exit_code: 0`). Attempt 5: actual abandon authority and typed errors.
- `go test -race ./cmd/fid -run '^(TestRunStartBindsOptionalObligationAndRejectsReusedExecutor|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 6: retained close caller convergence check.
- `go test -race ./cmd/fid -run '^(TestRunStartBindsOptionalObligationAndRejectsReusedExecutor|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	51.094s` (terminated `exit_code: 0`). Attempt 7: retained close caller convergence rerun.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 8: full task-8 gate 1.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 9: full task-8 gate 2.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	466.648s	coverage: 62.1% of statements` (terminated `exit_code: 0`). Attempt 10: full task-8 gate 3.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.4%` (terminated `exit_code: 0`). Attempt 11: full task-8 gate 4.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	211.237s` (terminated `exit_code: 0`). Attempt 12: full task-8 gate 5.
- `go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	218.693s` (terminated `exit_code: 0`). Attempt 13: full task-8 gate 6.
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	33.609s` (terminated `exit_code: 0`). Attempt 14: full task-8 gate 7 first benchmark.
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	33.762s` (terminated `exit_code: 0`). Attempt 15: full task-8 gate 8 second benchmark.
