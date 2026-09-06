# Phase 22 - Task 7: Durable budget accounting and execution admission - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Attempt 4 implemented the durable unknown/exhausted admission producer, boundary Stop/Attention effects, successor-chain read, rebuild and real-binary evidence. The required full race gate then failed in three retained tests. Two assertions displaced by Task 7's accepted derived-budget and Stop-producer behavior are outside this recovery's exact allowlist; the remaining journal assertion is also outside scope.
- Recovery attempts: 4
- Next action: Preserve the failed full-gate receipt and obtain a bounded same-task scope amendment for `cmd/fid/generic_run_test.go`, `cmd/fid/contract_surface_test.go`, and the independently diagnosed `internal/journal/recovery_schema_test.go` assertion if its accepted contract/provenance permits repair. Do not alter any protected wire, error, exit, schema, policy, control, or gate definition.
- Human attention required: no

## Summary of what was implemented

- Added a narrow budget-accounting core that derives consumption, remaining capacity, warning state, and exhaustion from append-only observed samples.
- Added direct red/green tests for the final admitted attempt, null-to-baseline recovery, duplicate idempotency, conflicting duplicate keys, stale predecessor rejection, and int64 overflow rejection.
- Converged envelope authoring to the accepted `budget-exhausted` policy, full accepted budget-unit set, zero warning threshold, and the sole `run.observe-budget` profile/grant capability.
- Applied the pre-authorized DEC-192 change by making `run.budget-observed` and `run.execution-blocked` projectable, while retaining every other recovery fence entry; added live-apply/replay parity coverage for the newly available event types.
- Added the first durable execution-accounting slice: `RunGate` reads each envelope budget's ordered cross-Run sample chain, derives fid-only charges from observed Attempt/Receipt facts, and persists the corresponding `budget_samples` plus Run runtime version update in the same Event transaction.
- Added pre-spawn `run.execution-blocked` commits for unknown and exhausted chains. Each appends exact run-scoped Stop/Attention effects and one runtime-version increment without an Attempt or child process.
- Coupled the first exhaustion boundary to the admitted Attempt/Receipt Event effects, while preserving a later admitted Receipt and preventing duplicate `attempts` charges on that Receipt.
- Made Run reads derive each Budget from the envelope-wide sample chain. The real process recovery journey proves a successor Run retains the original `attempts` charge; the Phase 21 journey proves the boundary, later denial, and rebuild journal path.

## Files created or modified

- `internal/workflow/budget.go`
- `internal/workflow/budget_test.go`
- `internal/workflow/envelope.go`
- `internal/workflow/run_read.go`
- `internal/workflow/activation.go`
- `internal/journal/project_effects.go`
- `internal/journal/checkpoint_projection_test.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-7-done.md`

## Documentation updated

- Added this recovery checkpoint report. No contract, schema, public-interface, policy, dependency, task-allocation, or gate-definition change is proposed.

## Validation Contract Coverage

- VAL-025: PARTIAL - `TestFIDBudgetChargesDeriveOnlyObservedAttemptAndReceiptFacts` proves the Attempt alone charges `attempts`, while the Receipt charges only non-success `failures` and measured `milliseconds`; this fixed the observed duplicate-source conflict. The real binary recovery journey executes the automatic start charge.
- VAL-026: PARTIAL - the preserved direct ordering/idempotency tests remain green; `TestP22EarlyRecoveryJourney` proves a successor Run derives the original Envelope sample as `consumed: 1, remaining: 2`.
- VAL-027: PARTIAL - `TestBudgetAdmissionDeniesUnknownAndExhaustedChainsBeforeAttempt` proves the exact unknown/exhausted classification. `appendExecutionBlocked` commits the accepted Event and effects, but the CLI closed-envelope renderer is outside the supplied allowlist.
- VAL-028: PARTIAL - `TestBudgetAdmissionBlockAndEffectsUseExactBudgetStopShape`, `TestBudgetBoundaryBlocksOnlyOnFirstExhaustion`, and `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` prove the boundary Event's sample/Stop/Attention/link effects and later no-Attempt denial. The full gate is red.
- VAL-029: PARTIAL - `TestP22EarlyRecoveryJourney` executes a real process, recovers to a successor, and asserts the original chain's `1/2` totals; the Phase 21 journey runs rebuild after the exhaustion/denial Event sequence. Full-gate completion remains blocked by out-of-allowlist retained assertions.
- VAL-030: PARTIAL - `BudgetAdmissionBlocks` admits warning budgets and rejects only unknown/exhausted values; full query-surface evidence remains unfinished.
- VAL-035: PARTIAL - the pre-edit 25.2% baseline and attempt-4 24.0% measurement are recorded. Direct evidence-package budget tests execute, but the required full gate is red.

## Tests run

- Required pre-edit baseline: `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.250s\tcoverage: 25.2% of statements` (terminated `exit_code: 0`)
- Required pre-edit baseline for recovery attempt 3: `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.254s\tcoverage: 25.2% of statements` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestFinalAllowedAttemptChargesAndRaisesStopAtomically$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Initial compile receipt: `AdmitBudgetSample` and `BudgetSample` were absent. A minimal accounting API was added; the immediately following semantic red attempt observed null consumption, null remaining, and no exhaustion state.
- `go test ./internal/workflow -run '^TestFinalAllowedAttemptChargesAndRaisesStopAtomically$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Semantic red receipt: `admission = workflow.BudgetAdmission{Admitted:true, ... Consumed:(*int64)(nil), Remaining:(*int64)(nil)}, Exhausted:false}, want admitted exhausted final attempt with zero remaining` (terminated `exit_code: 1`). Root cause: the initial accounting seam accepted the sample without deriving totals. Repair: derive known consumption, remaining, warning, and exhaustion from append-only samples. The byte-identical later command passed.
- `go test ./internal/workflow -run '^TestFinalAllowedAttemptChargesAndRaisesStopAtomically$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.004s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Semantic red receipt: `duplicate admission = workflow.BudgetAdmission{Admitted:false, ...}, budget sample conflict; want idempotent 40` (terminated `exit_code: 1`). Root cause: generated sample identity was compared before normalizing a replayed caller value. Repair: normalize an omitted sample ID to its source sample ID before semantic equality. The byte-identical later command passed.
- `go test ./internal/workflow -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.004s` (terminated `exit_code: 0`)
- Post-format focused verification: `gofmt -w internal/workflow/budget.go internal/workflow/budget_test.go && go test ./internal/workflow -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.005s` (terminated `exit_code: 0`)
- `go test ./internal/journal -run '^TestBudgetObservedEventProjectsIdenticallyLiveAndReplay$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Semantic red receipt: `ProjectEvent() error = Event-v2 "run.budget-observed" does not permit effect application before its recovery domain is available` (terminated `exit_code: 1`). Root cause: both accepted Task 7 event keys remained in `unavailableRecoveryEvents`. Repair: DEC-192 pre-authorized deletion of exactly `run.budget-observed` and `run.execution-blocked`; no other fence entry changed. The later focused command passed.
- `go test ./internal/journal -run '^(TestBudgetObservedEventProjectsIdenticallyLiveAndReplay|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.319s` (terminated `exit_code: 0`)
- `gofmt -w internal/workflow/budget.go internal/workflow/envelope.go internal/workflow/activation.go internal/journal/checkpoint_projection_test.go && go test ./internal/workflow ./internal/journal ./cmd/fid -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow|TestBudgetObservedEventProjectsIdenticallyLiveAndReplay|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestEnvelope)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.749s` (terminated `exit_code: 0`)
  - The preceding package lines were `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.006s` and `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.403s`; the command termination is the observed combined result.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-7-done.md project-docs/docs/phases/phase-22/tasks/task-7.md` - PASS - `PASS report-lint completed with no failures` (terminated `exit_code: 0`)
- `go test ./internal/evidence -run '^TestFIDBudgetChargesDeriveOnlyObservedAttemptAndReceiptFacts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red TDD evidence: `undefined: fidBudgetCharges`. Repair: add the fid-observed Attempt/Receipt charge derivation seam; the byte-identical command later passed.
- `go test ./internal/evidence -run '^TestFIDBudgetChargesDeriveOnlyObservedAttemptAndReceiptFacts$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.004s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestDeriveBudgetPreservesUnknownAndExhaustedChains$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red TDD evidence: `undefined: DeriveBudget`. Repair: expose the existing pure chain derivation for the execution admission reader; the byte-identical command later passed.
- `go test ./internal/workflow -run '^TestDeriveBudgetPreservesUnknownAndExhaustedChains$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.005s` (terminated `exit_code: 0`)
- `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.266s\tcoverage: 24.5% of statements` (terminated `exit_code: 0`)
- `go test ./internal/evidence ./internal/workflow -run '^(TestFIDBudgetChargesDeriveOnlyObservedAttemptAndReceiptFacts|TestDeriveBudgetPreservesUnknownAndExhaustedChains|TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.004s` (terminated `exit_code: 0`)
  - The preceding package line was `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.004s`; the command termination is the observed combined result.
- Required pre-edit baseline for recovery attempt 4: `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t(cached)\tcoverage: 24.5% of statements` (terminated `exit_code: 0`)
- `go test ./internal/workflow ./internal/evidence -run '^(TestBudgetAdmissionBlockAndEffectsUseExactBudgetStopShape|TestBudgetAdmissionDeniesUnknownAndExhaustedChainsBeforeAttempt)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red TDD evidence: `BudgetAdmissionBlocks`, `BudgetStopEffects`, and `budgetAdmissionBlocks` were undefined. Repair: add the accepted budget classification and exact Stop/Attention effect seam; the byte-identical later command passed.
- `go test ./internal/workflow ./internal/evidence -run '^(TestBudgetAdmissionBlockAndEffectsUseExactBudgetStopShape|TestBudgetAdmissionDeniesUnknownAndExhaustedChainsBeforeAttempt)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.004s` (terminated `exit_code: 0`)
- `go test ./internal/evidence -run '^TestBudgetBoundaryBlocksOnlyOnFirstExhaustion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red TDD evidence: `budgetBoundaryBlocks` was undefined. Repair: emit an exhaustion Stop only when a known available budget crosses to zero.
- `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - First real-entrypoint failure: `admit fid budget sample "attempts" with prior samples ... budget sample conflict`. Root cause: the Receipt producer repeated the Attempt's source sample. Repair: Receipt charge derivation excludes `attempts`; its later real-binary command passed.
- `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t1.816s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Red recovery evidence: successor Budget output retained null totals. Root cause: `loadRun` returned immutable binding JSON without deriving the Envelope-wide sample chain. Repair: derive each Run Budget from `budget_samples`; the byte-identical later command passed.
- `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t1.559s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.259s\tcoverage: 24.0% of statements` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Immutable required-gate failure: `TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets` retained null budget counters after Task 7's accepted derived-budget read; `TestPhase21AttentionBoundary` retained the pre-Task-7 exact producer list; `TestRecoverySchemaFreshInitAndReplay` received `NOT NULL constraint failed: run_recoveries.prior_run_id` instead of its accepted closed Event-v2 rejection. The first two assertions are displaced by this continuation, but all three files are outside the exact allowlist.
- `go tool cover -func=coverage.out` - BLOCKED - the failed required coverage command did not produce an accepted quality-gate profile; no observed termination result.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not started after the required repository-wide race gate failed outside the exact allowlist; no observed termination result.

## Deviations from spec

- No accepted contract or schema deviation. This is an incomplete same-task checkpoint, not a completed implementation.
- Checkpoint documentation-commit deviation: the first `git add project-docs/docs/phases/phase-22/tasks/task-7-done.md && git diff --cached --check && git commit -m "docs(task): record phase-22 task-7 report lint" -m "Refs: phase-22/task-7" && project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-7-done.md project-docs/docs/phases/phase-22/tasks/task-7.md && git status --porcelain -- project-docs/docs/phases/phase-22/tasks/task-7-done.md && git status --short` attempt ended `fatal: Unable to create '/home/mark/projects/fid/.git/index.lock': Read-only file system` (terminated `exit_code: 128`). Root cause: the sandbox denied Git's index-lock write. Repair: rerun the byte-identical command with the required repository-write approval after adding this durable receipt; the later result is recorded after its observed termination.
- Checkpoint documentation-commit recovery result: the byte-identical command later passed with final token ` M project-docs/docs/phases/phase-22/run-ledger.md` (terminated `exit_code: 0`), creating report commit `c3e375b` and leaving only the runner-owned ledger dirty.
- Mandatory checkpoint-report correction: the first `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-7-done.md project-docs/docs/phases/phase-22/tasks/task-7.md` attempt ended `FAIL report-lint found 3 failure(s)` with `exit_code: 3`. Root cause: the recovery category was not one of the template's typed values and incomplete gates were labelled `UNVERIFIED` instead of `BLOCKED`. Repair: use `in-scope-implementation` and record the unstarted required commands as blocked with no observed termination. The byte-identical later command is recorded after this committed correction.
- Convergence provenance: `cmd/fid/envelope_test.go` was introduced by `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` (`feat(envelope): add one-revision values`, phase-21/task-6); `cmd/fid/integration/p21_conformance_test.go` by `a99b4270f70e98111b5bae43c54ed92aed7078c6` (`test(conformance): add phase 21 wire corpus`, phase-21/task-22); and `cmd/fid/integration/p22_fake_test.go` plus `cmd/fid/integration/p22_recovery_journey_test.go` by `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` (`test(recovery): checkpoint task-6 CLI gap`, phase-22/task-6). Task 7's accepted section-6 authoring change invalidated each retained `"exhaustion_action": "unavailable"` fixture row; this continuation changes that exact row to `"budget-exhausted"`. Evidence command: `go test ./internal/workflow ./internal/journal ./cmd/fid -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow|TestBudgetObservedEventProjectsIdenticallyLiveAndReplay|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestEnvelope)' -count=1`.
- Convergence provenance: `internal/journal/checkpoint_projection_test.go` was introduced by `59c62b9e1d82c5f698d1e173908446d35b64f539` (`fix(recovery): admit checkpoint projection`, phase-22/task-4). DEC-192 and schema §5 changed the retained `TestRemainingRecoveryEventsRemainFenced` rows for `run.budget-observed` and `run.execution-blocked` from rejection to availability. This continuation removes exactly those two rows, adds their live/replay coverage, and uses the same focused journal command above. The five remaining fence rows are unchanged.
- Convergence provenance: `internal/workflow/envelope.go`, `internal/workflow/run_read.go`, and `internal/evidence/execution.go` were introduced as the old immutable-null-budget surface before the phase planning base (the convergence inventory identifies the generic Run predecessor as `5ca732f` phase-21/task-11 and Attempt/Receipt as `b694991` phase-21/task-15). Their subject-preserving convergence remains assigned to task 7; only the accepted authoring validator is changed in this checkpoint, and durable execution convergence remains unfinished.
- Recovery attempt 2 re-read the Task 7 authority and current producers. `internal/evidence/execution.go` still appends only the existing `execution_attempts` Effect before spawning and existing Receipt/Artifact Effects after termination; `internal/workflow/run_read.go` still reads only `run_budget_bindings`. This confirms the remaining work is the in-scope durable producer/read slice, rather than a report or DEC-192 fence defect.
- Recovery attempt 3 convergence: the Task 7 accepted change now requires `internal/evidence/execution.go` to replace its pre-task-7 Attempt-only and Receipt/Artifact-only Event effects. The implementation loads `budget_samples` using the Envelope/revision/budget key rather than only the current Run, so successor Runs retain the same ordered chain. It appends only fid-derived `attempts`, non-success `failures`, elapsed `milliseconds`, and spawn-failure `tool-failures` values; external units still receive no automatic sample. The next continuation must complete the accepted unknown/exhausted rejection and Stop/Attention Event effects rather than changing this bounded behavior.
- Mandatory attempt-4 gate deviation: the first `go test -race -coverprofile=coverage.out ./...` attempt ended `FAIL` (terminated `exit_code: 1`). Root cause: the accepted Task 7 read now derives known fid-owned zero totals and adds `internal/workflow/budget.go` as an automatic Stop producer, invalidating two retained cmd/fid assertions; a separate journal fixture also fails outside the supplied allowlist. Repair is not authorized in this continuation because those exact test paths are not allowed. The later byte-identical command has not run.
- Attempt-4 convergence provenance: `cmd/fid/generic_run_test.go:45` was introduced before Phase 22 and retains the old immutable-null Budget assertion; Task 7's accepted durable samples/read derivation displaces that assertion. `cmd/fid/contract_surface_test.go:99` retains the old exact automatic Attention producer inventory; Task 7's accepted boundary Stop producer displaces that row. Both need subject-preserving convergence, but neither is in the exact recovery allowlist.

## Concerns, technical debt, or notes for the Tech Lead

- The full-gate scope gap needs a bounded same-task recovery directive before the retained assertions can converge.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d85ab7018903300ccb53199c49fdea04e1a0bdc6`

## Commit subject

- `feat(budget): enforce durable execution admission`
