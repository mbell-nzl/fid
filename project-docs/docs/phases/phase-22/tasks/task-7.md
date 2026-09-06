# Phase 22 - Task 7: Durable budget accounting and execution admission

## Objective

Record trusted resource consumption and atomically enforce unknown/exhausted budgets at execution admission and Receipt commit. Retain the same sample chain across linked Runs.

## Validation Contract Coverage

- VAL-025: fid-owned resource units are charged only from observed execution facts.
- VAL-026: Budget samples are idempotent and ordered across recovery.
- VAL-027: Unknown budget measurements remain null and block admission.
- VAL-028: Budget exhaustion commits its charge, Stop, and Attention atomically.
- VAL-029: Consumed budgets survive successor Runs and rebuild.
- VAL-030: Budget warning and exhaustion states have distinct behavior.
- VAL-035: Reopened evidence recovery branches have direct package tests and measured coverage.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- DEC-184/public-interface section 2.1: add exactly run.observe-budget to
  profile/grant acceptance; activation and binary capability IDs are not actions.
- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- contracts.md section 6; schema.md budget Effects.
- CF-034: capture evidence package coverage before edits and after completion, with direct tests of budget and abandoned-execution branches.
- Five fid-owned units cannot be supplied by a caller; external measurements require activated harness provenance and evidence reference.
- Unknown is null, baseline/delta sample order is explicit, and overflow rejects.
- The final admitted process may terminate and write a Receipt after exhaustion; a later process must never spawn.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.
- DEC-192 progressive fence: `internal/journal/` is already in scope. Delete
  exactly the `"run.budget-observed"` and `"run.execution-blocked"` keys from
  `unavailableRecoveryEvents` in `internal/journal/project_effects.go` as part
  of implementing those Events, and prove live-apply/replay parity for each.
  Leave every other key fenced.

## Scope

### Files To Create Or Modify

- `internal/workflow/budget.go`
- `internal/workflow/budget_test.go`
- `internal/workflow/envelope.go`
- `internal/workflow/run_read.go`
- `internal/evidence/execution.go`
- `internal/evidence/execution_test.go`
- `internal/evidence/budget_test.go`
- `internal/workflow/activation.go`
- `internal/journal/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-7-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-6-done.md` - actual predecessor hand-back.

### Out Of Scope

- Telemetry summary capability, provider integration, budget edits/resets, Review-cycle charging, and public budget wiring.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] fid-owned resource units are charged only from observed execution facts. (VAL-025)
- [ ] Budget samples are idempotent and ordered across recovery. (VAL-026)
- [ ] Unknown budget measurements remain null and block admission. (VAL-027)
- [ ] Budget exhaustion commits its charge, Stop, and Attention atomically. (VAL-028)
- [ ] Consumed budgets survive successor Runs and rebuild. (VAL-029)
- [ ] Budget warning and exhaustion states have distinct behavior. (VAL-030)
- [ ] Reopened evidence recovery branches have direct package tests and measured coverage. (VAL-035)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestFinalAllowedAttemptChargesAndRaisesStopAtomically fails before the Attempt and exhaustion effects are coupled.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-7-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

Before any test or production edit, record this baseline command once:

```text
go test -race -cover ./internal/evidence
```

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
go test -race -cover ./internal/evidence
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

- Start after task 6 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 3.
- Commit using Conventional Commits with `Refs: phase-22/task-7`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
