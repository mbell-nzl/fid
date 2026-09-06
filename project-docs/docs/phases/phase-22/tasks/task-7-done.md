# Phase 22 - Task 7: Durable budget accounting and execution admission - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: implementation
- Condition: The tested accounting core is checkpointed, but durable sample projection, admission/Receipt coupling, the two Event producers and their live/replay parity, retained fixture convergence, and the complete required quality gate remain unfinished.
- Recovery attempts: 0
- Next action: Continue phase-22/task-7 from this checkpoint. Persist `BudgetSample` chains in `budget_samples`, wire automatic fid observations and unknown/exhausted admission into `internal/evidence.RunGate`, then add the required journal and entrypoint evidence before running the exact full gate.
- Human attention required: no

## Summary of what was implemented

- Added a narrow budget-accounting core that derives consumption, remaining capacity, warning state, and exhaustion from append-only observed samples.
- Added direct red/green tests for the final admitted attempt, null-to-baseline recovery, duplicate idempotency, conflicting duplicate keys, stale predecessor rejection, and int64 overflow rejection.

## Files created or modified

- `internal/workflow/budget.go`
- `internal/workflow/budget_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-7-done.md`

## Documentation updated

- Added this recovery checkpoint report. No contract, schema, public-interface, policy, dependency, task-allocation, or gate-definition change is proposed.

## Validation Contract Coverage

- VAL-025: UNVERIFIED - automatic fid-owned charging and caller-substitution rejection are not yet connected to execution.
- VAL-026: PARTIAL - `TestBudgetSamplesAreIdempotentOrderedAndRecoverUnknownBaselines` in `internal/workflow/budget_test.go` covers the in-memory accounting chain; durable recovery and rebuild are unfinished.
- VAL-027: PARTIAL - the same direct test proves a null sample leaves totals null until a baseline; no durable admission denial exists yet.
- VAL-028: PARTIAL - `TestFinalAllowedAttemptChargesAndRaisesStopAtomically` in `internal/workflow/budget_test.go` proves the derived boundary state only; Event/Stop/Attention atomicity is unfinished.
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

## Exact quality-gate results

- `go test -race -cover ./internal/evidence` - UNVERIFIED after implementation; reserved for the complete task gate after the evidence changes land.
- `gofmt -l cmd internal` - UNVERIFIED; not run because the task implementation is incomplete.
- `golangci-lint run ./...` - UNVERIFIED; not run because the task implementation is incomplete.
- `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED; not run because the task implementation is incomplete.
- `go tool cover -func=coverage.out` - UNVERIFIED; not run because the task implementation is incomplete.
- `go test -tags integration -race ./cmd/fid/integration/...` - UNVERIFIED; not run because the task implementation is incomplete.

## Deviations from spec

- No accepted contract or schema deviation. This is an incomplete same-task checkpoint, not a completed implementation.
- Convergence provenance: `internal/workflow/envelope.go`, `internal/workflow/run_read.go`, and `internal/evidence/execution.go` were introduced as the old immutable-null-budget surface before the phase planning base (the convergence inventory identifies the generic Run predecessor as `5ca732f` phase-21/task-11 and Attempt/Receipt as `b694991` phase-21/task-15). Their subject-preserving convergence remains assigned to task 7 and has not yet been edited.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
