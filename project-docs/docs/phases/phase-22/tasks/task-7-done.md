# Phase 22 - Task 7: Durable budget accounting and execution admission - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The accounting core, accepted authoring policy, and DEC-192 recovery-fence availability are checkpointed. Durable `budget_samples` projection, execution/Receipt coupling, accepted Event producers and effects, recovery/rebuild evidence, and the required quality gate remain unfinished.
- Recovery attempts: 1
- Next action: Continue phase-22/task-7 from this checkpoint. Persist `BudgetSample` chains in `budget_samples`, wire automatic fid observations and unknown/exhausted admission into `internal/evidence.RunGate`, then replace the generic Event fixtures with exact budget/Stop/Attention effects and add recovery, rebuild, and real-entrypoint evidence before running the exact full gate.
- Human attention required: no

## Summary of what was implemented

- Added a narrow budget-accounting core that derives consumption, remaining capacity, warning state, and exhaustion from append-only observed samples.
- Added direct red/green tests for the final admitted attempt, null-to-baseline recovery, duplicate idempotency, conflicting duplicate keys, stale predecessor rejection, and int64 overflow rejection.
- Converged envelope authoring to the accepted `budget-exhausted` policy, full accepted budget-unit set, zero warning threshold, and the sole `run.observe-budget` profile/grant capability.
- Applied the pre-authorized DEC-192 change by making `run.budget-observed` and `run.execution-blocked` projectable, while retaining every other recovery fence entry; added live-apply/replay parity coverage for the newly available event types.

## Files created or modified

- `internal/workflow/budget.go`
- `internal/workflow/budget_test.go`
- `internal/workflow/envelope.go`
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

- VAL-025: UNVERIFIED - `run.observe-budget` authoring acceptance is converged, but automatic fid-owned charging and caller-substitution rejection are not yet connected to execution.
- VAL-026: PARTIAL - `TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines` in `internal/workflow/budget_test.go` covers the in-memory accounting chain; durable recovery and rebuild are unfinished.
- VAL-027: PARTIAL - the same direct test proves a null sample leaves totals null until a baseline; no durable admission denial exists yet.
- VAL-028: PARTIAL - `TestFinalAllowedAttemptChargesAndRaisesStopAtomically` in `internal/workflow/budget_test.go` proves the derived boundary state only. `TestBudgetObservedEventProjectsIdenticallyLiveAndReplay` and `TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay` in `internal/journal/checkpoint_projection_test.go` prove generic projector live/replay parity after DEC-192; accepted budget/Stop/Attention effects and atomicity are unfinished.
- VAL-029: UNVERIFIED - successor Run and rebuild evidence is unfinished.
- VAL-030: PARTIAL - the accounting result distinguishes warning from exhaustion; durable query and execution behavior are unfinished.
- VAL-035: PARTIAL - required pre-edit evidence coverage was captured at 25.2%; direct budget coverage was added, while the abandoned-execution branch and final measured coverage are unfinished.

## Tests run

- Required pre-edit baseline: `go test -race -cover ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.250s\tcoverage: 25.2% of statements` (terminated `exit_code: 0`)
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

## Exact quality-gate results

- `go test -race -cover ./internal/evidence` - BLOCKED - not rerun after implementation because the task is checkpointed before the required evidence wiring is complete; no observed termination result.
- `gofmt -l cmd internal` - BLOCKED - not started because the task implementation is incomplete; no observed termination result.
- `golangci-lint run ./...` - BLOCKED - not started because the task implementation is incomplete; no observed termination result.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not started because the task implementation is incomplete; no observed termination result.
- `go tool cover -func=coverage.out` - BLOCKED - not started because its required coverage profile was not produced; no observed termination result.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not started because the task implementation is incomplete; no observed termination result.

## Deviations from spec

- No accepted contract or schema deviation. This is an incomplete same-task checkpoint, not a completed implementation.
- Mandatory checkpoint-report correction: the first `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-7-done.md project-docs/docs/phases/phase-22/tasks/task-7.md` attempt ended `FAIL report-lint found 3 failure(s)` with `exit_code: 3`. Root cause: the recovery category was not one of the template's typed values and incomplete gates were labelled `UNVERIFIED` instead of `BLOCKED`. Repair: use `in-scope-implementation` and record the unstarted required commands as blocked with no observed termination. The byte-identical later command is recorded after this committed correction.
- Convergence provenance: `cmd/fid/envelope_test.go` was introduced by `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` (`feat(envelope): add one-revision values`, phase-21/task-6); `cmd/fid/integration/p21_conformance_test.go` by `a99b4270f70e98111b5bae43c54ed92aed7078c6` (`test(conformance): add phase 21 wire corpus`, phase-21/task-22); and `cmd/fid/integration/p22_fake_test.go` plus `cmd/fid/integration/p22_recovery_journey_test.go` by `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` (`test(recovery): checkpoint task-6 CLI gap`, phase-22/task-6). Task 7's accepted section-6 authoring change invalidated each retained `"exhaustion_action": "unavailable"` fixture row; this continuation changes that exact row to `"budget-exhausted"`. Evidence command: `go test ./internal/workflow ./internal/journal ./cmd/fid -run '^(TestFinalAllowedAttemptChargesAndRaisesStopAtomically|TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines|TestBudgetSamplesRejectOverflow|TestBudgetObservedEventProjectsIdenticallyLiveAndReplay|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestRemainingRecoveryEventsRemainFenced|TestEnvelope)' -count=1`.
- Convergence provenance: `internal/journal/checkpoint_projection_test.go` was introduced by `59c62b9e1d82c5f698d1e173908446d35b64f539` (`fix(recovery): admit checkpoint projection`, phase-22/task-4). DEC-192 and schema §5 changed the retained `TestRemainingRecoveryEventsRemainFenced` rows for `run.budget-observed` and `run.execution-blocked` from rejection to availability. This continuation removes exactly those two rows, adds their live/replay coverage, and uses the same focused journal command above. The five remaining fence rows are unchanged.
- Convergence provenance: `internal/workflow/envelope.go`, `internal/workflow/run_read.go`, and `internal/evidence/execution.go` were introduced as the old immutable-null-budget surface before the phase planning base (the convergence inventory identifies the generic Run predecessor as `5ca732f` phase-21/task-11 and Attempt/Receipt as `b694991` phase-21/task-15). Their subject-preserving convergence remains assigned to task 7; only the accepted authoring validator is changed in this checkpoint, and durable execution convergence remains unfinished.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a37e7d9a18fcb8d19432e1ce5100a33db839f978`

## Commit subject

- `feat(budget): converge budget event admission policy`
