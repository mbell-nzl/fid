# Phase 22 - Task 15: Import CLI and three-way recovery parity

## Objective

Expose import through the built CLI and prove live, rebuilt, and imported recovery state equivalence. Re-export the restored project deterministically.

## Validation Contract Coverage

- VAL-063: Live, rebuilt, and imported current state are semantically equivalent.
- VAL-064: Re-export of reconstructed state preserves canonical recovery content.
- VAL-065: The built CLI completes export, validate, import, rebuild, and resume.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- contracts.md section 10; schema.md section 6.
- Use named JSON reads and ordered Events, not an opaque SQLite dump, as parity evidence.
- Add portable recovery cases to the shared semantic runner/fake now; the same
  expectations run against fake and fid before later capability enablement.
- Normalize only archive/destination absolute paths and publication metadata; preserve entity/Event IDs, counters, history, statuses, errors, and null dimensions.
- Exercise active/interrupted/abandoned Runs, deduplicated Attention, exhausted and unknown budgets, partial activation, and missing Artifacts in the corpus.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `cmd/fid/export.go`
- `cmd/fid/export_test.go`
- `cmd/fid/integration/p22_import_recovery_test.go`
- `cmd/fid/integration/testdata/p22_conformance/`
- `internal/export/`
- `internal/workflow/resume.go`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-15-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-14-done.md` - actual predecessor hand-back.

### Out Of Scope

- Changing the acceptance journey to avoid a failure, whole-product qualification, and new semantic operations.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Live, rebuilt, and imported current state are semantically equivalent. (VAL-063)
- [ ] Re-export of reconstructed state preserves canonical recovery content. (VAL-064)
- [ ] The built CLI completes export, validate, import, rebuild, and resume. (VAL-065)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestRecoveryLiveRebuildImportParity compares the same durable journey through all three paths.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-15-done.md`: implementation,
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

- Start after task 14 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 5 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-15`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
