# Phase 22 - Task 9: Stop, budget, and closure commands

## Objective

Wire the new domain operations, expose committed execution denials unambiguously, and remove run.finish. Apply Stop guards through the same evaluator used by completion and Issue transition.

## Validation Contract Coverage

- VAL-036: Stop, Attention, budget, and closure leaves implement their declared contracts.
- VAL-037: All guarded actions apply the same scoped Stop predicates.
- VAL-038: A committed execution denial names its durable effects without implying execution success.
- VAL-039: Caller-reported run.finish is absent after derived closure lands.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- public-interface.md Stop/Attention/budget/closure leaves and section 4.
- Committed denials name their Event and replay result while retaining a nonzero exit.
- DEC-197: complete the minimum domain/read support required by this task's
  accepted leaves. Add Stop/Attention filtered read APIs and the durable
  run.observe-budget producer using the accepted sample/authority/version/
  idempotency/Event contract. Reuse task-2/task-7 rules rather than implementing
  business rules in CLI adapters. Extend the evidence live-frame builder only
  for the accepted pre-start committed denial; preserve ordinary started-frame
  and duplicate-terminal validation. Record predecessor provenance in the
  report; prior completed tasks remain closed.
- Extend task 6's shared case runner/fake with these implemented operations;
  identical semantic expectations must pass against both fake and fid now.
- DEC-196: strengthen the existing shared runner into semantic conformance
  while extending it. The fake must derive typed results from requests and
  independent in-memory state, never echo fixture-declared expected fields,
  outcomes, or error codes. Shared assertions compare independently declared
  values and identity relationships, including same-key replay without a
  second mutation and rejection without effects. Include the existing recovery
  cases and this task's Stop/Attention/budget/closure operations; record
  task-6 provenance at f3be616f0587881ee3f5c3716ffbf6167e3f2a20.
- Remove finish parsing/outcome requests and current run.finished fixture acceptance only after close/abandon replacements pass.
- Remove run.finish from activated profile and Envelope grant acceptance in
  the same task; both the retired action and retired leaf must reject. Final
  action membership is public-interface section 2.1, never inferred from help.
- Distinguish recovery errors with no durable Stop from explicit stop.raise
  and committed budget denials; Attention requires disposition Event IDs.
- Retained tests converge by invalidating change, including Phase 21 journey closure steps.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `cmd/fid/stop.go`
- `cmd/fid/attention.go`
- `cmd/fid/run.go`
- `cmd/fid/gate.go`
- `cmd/fid/live_output.go`
- `cmd/fid/response_errors.go`
- `internal/cli/`
- `internal/evidence/live.go` - only pre-start committed-denial terminal
  construction in AppendFidError and its minimum call-boundary support.
- `internal/evidence/live_test.go` - corresponding pre-start/started/duplicate
  terminal cases, without changing observed execution outcomes.
- `internal/workflow/stop_read.go` and `internal/workflow/stop_read_test.go` -
  accepted Stop/Attention list filters, limits and read-only results.
- `internal/workflow/run_budget.go` and `internal/workflow/run_budget_test.go` -
  durable run.observe-budget operation, authority, expected versions,
  source/predecessor deduplication, measurement binding and atomic Effects.
- `internal/workflow/budget.go` - only shared helper support needed by that
  producer, retaining the existing accepted admission and accounting rules.
- `internal/workflow/run.go` - only removal of the retired run.finish action,
  request, validation and producer, plus the shared applicable-Stop check on
  run.start and necessary unused-import cleanup.
- `internal/workflow/recovery.go` - only shared applicable-Stop checks before
  ordinary/force recovery mutation, retaining its existing authority rules.
- `internal/evidence/execution.go` - only invoking the shared applicable-Stop
  predicate before new Attempt/spawn; retain committed budget denial and
  already-admitted Receipt behavior and idempotent result replay.
- Matching in-package tests for the three shared Stop call sites above.
- `internal/workflow/evaluation.go`
- `internal/workflow/activation.go` - remove run.finish action acceptance.
- `internal/workflow/envelope.go` - converge grant validation for retired action.
- `internal/journal/` - retire current run.finished Event acceptance.
- `internal/validation/`
- `internal/issues/transition.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `cmd/fid/integration/p22_fake_test.go` - DEC-196 semantic runner/fake repair.
- `cmd/fid/integration/testdata/p22_conformance/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-9-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-8-done.md` - actual predecessor hand-back.

### Out Of Scope

- Automatic retries, Human Outcome resolution, and silent gate/error remapping.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Stop, Attention, budget, and closure leaves implement their declared contracts. (VAL-036)
- [ ] Stop/Attention lists apply accepted filters/limits without mutation, and
  run.observe-budget durably applies a qualifying external sample with version,
  idempotency and live/rebuild parity. Fid-owned unit, missing authority,
  malformed measurement, stale version and conflicting predecessor/source
  reject without effects. A pre-start committed-denial terminal has no Attempt
  or start/output frame; started and duplicate-terminal validation remains.
  (VAL-036/038; task-2/task-7/task-8 provenance under DEC-197)
- [ ] All guarded actions apply the same scoped Stop predicates. (VAL-037)
- [ ] A committed execution denial names its durable effects without implying execution success. (VAL-038)
- [ ] Shared fake/fid cases assert typed recovery lineage and replay identity,
  Stop/Attention disposition relationships, budget counters and denial effects,
  closure outcomes, and rejected-request state preservation. Wrong types or
  values must fail, even when every required result field is present. The fake
  derives outcomes independently of expected results. (VAL-036..038; supports
  task-18 VAL-073, whose final completeness ownership remains there)
- [ ] Caller-reported run.finish leaf, action, and current Event acceptance are
  absent after derived closure lands; old names reject in the corpus. (VAL-039)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestBudgetDenialReportsCommittedStopWithoutStartFrame fails before the new failure/live representation exists.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-9-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 8 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 3 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-9`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
