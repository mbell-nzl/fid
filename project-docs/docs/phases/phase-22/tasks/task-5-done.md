# Phase 22 — Task 5: Atomic linked Run recovery — Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: FQ-04 same-task repair complete; independent final re-review remains due.
- Recovery attempts: 1
- Next action: Runner closes this owner and continues Task 7; final QA follows all queued repairs.
- Human attention required: no

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

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	68.149s	coverage: 44.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `62.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.861s` (terminated `exit_code: 0`)

## Deviations from spec

- None. DEC-192's `run.recovered` fence amendment was pre-authorized for this task.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `11e8a84a2ec910597e7d152f2105948ecdd2a99f`

## Commit subject

- `test(recovery): keep fake checkpoint updates process-free`


## Final-QA FQ-04 same-task continuation — 2026-09-09

Authority: review-7, final-review.md FQ-04, task-5.md, contracts.md §§3–4.
Root is the sole writer. Original receipts and historical claims remain above;
this appendix supplies the previously missing direct proof and current repair.
Prior primary binding: b7dce8e747a3f2d720002faee3db7f95b3cb4019,
`feat(recovery): add atomic linked run recovery`.

### Implementation and validation evidence

- VAL-014/015/021/080: RecoverRun reads immutable Claim paths and the latest
  accepted checkpoint. It recomputes the digest with those bindings instead of
  treating caller adoption as approved scope. Ordinary recovery cannot relabel
  outside work. After exact Human authority is verified, force may adopt the
  precise outside-scope set; that classification and the original Claim scope
  remain unchanged. Protected paths remain forbidden. No scope/digest/wire policy
  changed. The existing public ValidateReconciliationAdoption remains ordinary.
- VAL-020: the specifically referenced force grant must match the Actor, Issue
  scope, force action and required force capability. A different valid grant
  cannot make an invalid referenced ordinal sufficient.
- TestRecoveryUsesImmutableClaimScopeForOrdinaryAndForceAdoption covers seven
  ordinary/force, forged-scope, stale/omitted and protected cases with actual
  observed Git digests, unchanged rejected journal, retained Claim scope and
  immutable Envelope. TestRecoveryUsesLatestDurableCheckpoint executes a real
  checkpoint and subsequent inspection, then verifies recovery retains that ID.
- VAL-017/019/020: TestRecoverRunRequiresExactObservationVersionsAndAuthority
  covers 39 finite cases across ordinary and force: valid controls, unknown,
  missing/nonexistent/wrong-Run/wrong-executor/wrong-kind observation, stale
  Run/Claim versions, absent grant, closed Envelope, Actor mismatch, live ordinary,
  active/released Claim, terminal force, non-Human type, absent capability,
  missing reference/grant reference and incorrect Envelope/revision/ordinal.
  Four referenced-grant negatives retain another valid grant. Every rejection
  checks its error type, all projection tables and journal bytes unchanged.
  Retained tagged null-observation wire cases execute under the full gate.
- The successful authority controls attempt an otherwise valid ObserveRun on
  the predecessor, with the exact executor, current version, valid observation
  time, activated-harness profile and explicit grant. It must return the exact
  RunInactiveError for the interrupted predecessor and leave all facts unchanged.
- Actual CLI TestP22RecoveryUsesClaimScopeAndExactHumanOutsideAdoption rejects
  ordinary forged scope and force-protected adoption with the exact conflict/
  exit 4 and no durable change. Known-live Human force accepts precisely identified
  outside work, preserves outside-scope classification and Human reference, and
  leaves the successor Claim's original base.txt scope intact.
- VAL-018/021: TestP22RecoverySurvivesActualResponseLossAndAppendProjectionFault
  builds a populated workspace through actual fid operations: activation, Issue,
  Envelope/Obligation/policy, Claim/Run, external baseline, actual execution,
  Receipt and linked Evidence. A child process calls actual RecoverRun and exits
  before returning a result. Three modes prove pre-append failure, durable append
  followed by projection failure, and successful commit followed by response loss.
  The test-only SQLite trigger/scalar passes dry-run and fails the second call
  (actual projection); no production fault API or failure framework was added.
- Pre-append failure leaves no Event or partial facts. Post-append failure has
  exactly one durable Event and unchanged projection; fresh fid catch-up applies
  it. Exact-key retry returns the original successor/Claim/link once, with the
  correct replay flag. A different key is rejected by actual CLI as already
  linked; changed same-key request is rejected as typed IdempotencyConflictError
  by the domain. Neither changes any table or Event count. All populated Issue,
  revision, sources, grants, Obligations, policy, Receipt/Evidence and consumed
  budget rows are compared unchanged. Recovery creates no new Attempt or Receipt.

### Files, convergence and mandatory deviations

- Product: internal/workflow/recovery.go and reconciliation.go. New direct tests:
  recovery_scope_test.go and recovery_authority_test.go. New process fixtures:
  p22_recovery_fault_test.go and p22_recovery_scope_test.go. The retained full
  p22_final_journey_test.go adds one actual inspection after its checkpoint.
- Main repair: edb134f9fcd9a9512047b5f3172d6a9d398cac35,
  `fix(recovery): bind adoption and force authority to durable scope`.
  Lint-only followup: 52d93d41ba5cd29231001bcefc5c53d290122645. The current primary binding is the final retained-caller convergence commit.
- Scope/checkpoint red ended `FAIL` (exit_code: 1): ordinary forged scope was
  accepted, force-outside was rejected, and the current checkpoint digest was
  rejected. The first checkpoint setup also had a missing SQL quote; that failed
  construction is retained separately. After the causal repair the byte-identical
  focused command passed at 20.878s.
- Authority red ended `FAIL` (exit_code: 1): referenced wrong-Issue scope and
  missing required force capability were accepted despite another valid grant.
  The exact-grant check repaired these; the same command passed at 86.012s.
  Two fixture attempts used inadmissible state: revision2 is correctly forbidden
  by this one-revision schema, and the accepted release reason is operator-released.
  Corrected fixtures test a closed Envelope and that admitted release reason;
  force wrong-revision reference rejection remains. No schema was loosened.
- Initial interruption proof ended `FAIL` (exit_code: 1). Its replay assertion
  read the adapter's deliberately unadvanced ID after replay, and it assumed a
  CLI idempotency-conflict response that this command does not expose.
  The corrected proof reads the returned successor ID and checks the existing
  domain conflict type for a changed digest, retaining actual public already-linked
  rejection for another key. The same full focused command passed at 10.147s.
- The retained full journey red ended `FAIL` (exit_code: 1) with
  `reconciliation digest is stale or changed`: it reused its pre-checkpoint
  digest. Introducing commit 24e3cbdb22dd9e43b55636da313dcef1ef633b6f,
  Phase22/Task18. One actual post-checkpoint inspection preserves the complete
  original proof/import/parity journey; the identical command passed at 4.850s.
  The Task18 report is unchanged.
- Required lint first failed with `* staticcheck: 1` (exit_code: 1), QF1001 on
  the adoption predicate. An equivalent De Morgan simplification repaired it;
  the byte-identical lint passed with `0 issues.` (exit_code: 0). The entire final
  sequence restarted with format; neither command nor acceptance semantics changed.
- No dependency, architecture, product policy, wire, Event, schema, protected
  control, threshold or validation assertion changed. Independent final-review-2
  remains required after all owners. Human judgments remain blank.

- The first required tagged gate ended `FAIL` (exit_code: 1) at 135.047s:
  the dual adapter, early clean/dirty journey and Human fixture generator still
  reused their pre-checkpoint digest. The shared p22_fake_test.go checkpoint
  adapter now performs an actual inspection before the next request; its
  introducing commit is f3be616f0587881ee3f5c3716ffbf6167e3f2a20, Phase22/Task6.
  The early journey adds the same explicit read; origin
  c69456fc9bf42f811e4999f1968a48c4a82d3cc5, Phase22/Task6. The Human generator
  (24e3cbdb22dd9e43b55636da313dcef1ef633b6f, Task18) consumes that corrected
  adapter without editing its assertions. Commit
  2dd2c8bdca27be464fa66e91755d8d91cb93d2dd records these seven fixture lines.
  No corpus row, expected response, gate, or original owner report changed.
  The unchanged full sequence reruns after this causal fixture convergence;
  the failed receipt remains below.

- The subsequent required tagged gate ended `FAIL` (exit_code: 1) at
  140.856s: the shared checkpoint advance helper also serves the fake adapter,
  which has no executable. Its attempted inspection reported
  `process error = exec: no command, want exit error`. The helper now limits
  the process read to an adapter with its real binary configured; fake response
  validation and all 41 corpus cases remain intact. Commit
  11e8a84a2ec910597e7d152f2105948ecdd2a99f records this test-only correction.
  The focused dual-adapter command passed at 7.762s (exit_code: 0), and the
  entire unchanged final gate sequence runs again below.

### Continuation command receipts

Working directory: `/home/mark/projects/fid`. Root captured the following output
tokens and observed terminations in this recording session. All failed attempts
remain; the five final task commands retain their original order.

- `go test -race ./internal/workflow -run '^TestRecoveryUses(ImmutableClaimScopeForOrdinaryAndForceAdoption|LatestDurableCheckpoint)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 1: Claim scope and checkpoint red.
- `go test -race ./internal/workflow -run '^TestRecoveryUses(ImmutableClaimScopeForOrdinaryAndForceAdoption|LatestDurableCheckpoint)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 2: scope and checkpoint corrected fixture red.
- `go test -race ./internal/workflow -run '^TestRecoveryUses(ImmutableClaimScopeForOrdinaryAndForceAdoption|LatestDurableCheckpoint)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	20.878s` (terminated `exit_code: 0`). Attempt 3: scope and checkpoint green.
- `go test -race ./internal/workflow -run '^TestRecoverRunRequiresExactObservationVersionsAndAuthority$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 4: exact authority and observation matrix red.
- `go test -race ./internal/workflow -run '^TestRecoverRunRequiresExactObservationVersionsAndAuthority$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	86.012s` (terminated `exit_code: 0`). Attempt 5: exact authority matrix green and valid fencing.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestP22RecoverySurvivesActualResponseLossAndAppendProjectionFault|TestP22RecoveryCrashProcess)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 6: real process recovery interruption proof.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestP22RecoverySurvivesActualResponseLossAndAppendProjectionFault|TestP22RecoveryCrashProcess)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.147s` (terminated `exit_code: 0`). Attempt 7: interruption proof corrected replay and domain conflict checks.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22RecoveryUsesClaimScopeAndExactHumanOutsideAdoption$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.241s` (terminated `exit_code: 0`). Attempt 8: actual CLI Claim scope and Human force proof.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 9: retained checkpoint journey convergence check.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.850s` (terminated `exit_code: 0`). Attempt 10: retained checkpoint journey green.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 11: required Task5 post-gate.
- `golangci-lint run ./...` - BLOCKED - `* staticcheck: 1` (terminated `exit_code: 1`). Attempt 12: required Task5 post-gate.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 13: required Task5 post-gate after lint repair.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 14: required Task5 post-gate after lint repair.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	339.754s	coverage: 61.2% of statements` (terminated `exit_code: 0`). Attempt 15: required Task5 full race coverage.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.2%` (terminated `exit_code: 0`). Attempt 16: required Task5 final gates.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 17: required Task5 final gates.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 18: required Task5 gate after retained caller convergence.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 19: required Task5 gate after retained caller convergence.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 61.2% of statements` (terminated `exit_code: 0`). Attempt 20: required Task5 full race after fixture convergence.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.2%` (terminated `exit_code: 0`). Attempt 21: required Task5 final gates after fixture convergence.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 22: required Task5 final gates after fixture convergence.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22DualAdapterRecoveryConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.762s` (terminated `exit_code: 0`). Attempt 23: shared fake adapter causal correction.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 24: final Task5 gates after shared adapter correction.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 25: final Task5 gates after shared adapter correction.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 61.2% of statements` (terminated `exit_code: 0`). Attempt 26: final Task5 gates after shared adapter correction.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.2%` (terminated `exit_code: 0`). Attempt 27: final Task5 gates after shared adapter correction.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	142.025s` (terminated `exit_code: 0`). Attempt 28: final Task5 gates after shared adapter correction.

Final closure: all five unchanged commands passed after the final fixture correction. The previously failed integration gate now passed at 142.025s, terminated exit_code: 0; full race CLI package 426.722s, aggregate coverage 65.2%. No outstanding Task 5 repair remains; independent final-review-2 determines final finding disposition.
