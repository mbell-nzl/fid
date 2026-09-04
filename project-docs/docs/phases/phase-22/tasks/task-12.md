# Phase 22 - Task 12: Recovery freshness and actionable resume

## Objective

Evaluate evidence against the current observed target after recovery and render only valid next actions. Keep proof validity and Artifact availability distinct.

## Validation Contract Coverage

- VAL-048: A changed repository target stales prior repository-bound proof.
- VAL-049: An unchanged target retains proof through recovery.
- VAL-050: Artifact loss preserves Evidence identity while changing only availability and policy outcome.
- VAL-051: Fresh resume reconstructs bounded recovery context without prior chat.
- VAL-052: Resume action eligibility matches actual mutation predicates.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 9; public-interface.md Resume.
- DEC-185/contracts sections 3.1/3.2 bind all comparisons. Reuse task 4's
  observer and independent vectors; do not invent a freshness-specific digest.
  Use a real repository to capture canonical dirty targets through public
  inspection and retained Receipt/Evidence input paths. Prove an unchanged
  dirty target stays fresh and a one-byte same-HEAD change becomes stale.
  Legacy arbitrary digest strings are preserved, never silently refreshed.
- Run/Registration repository matching must allow the reconciled current target within unchanged approved location/path scope; do not silently widen scope.
- Adoption, branch rename, and recovery do not refresh proof.
- An unavailable Findings dimension is explicit; no empty-list claim of completed Review support.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/workflow/resume.go`
- `internal/workflow/resume_test.go`
- `internal/workflow/evaluation.go`
- `internal/workflow/evaluation_test.go`
- `internal/evidence/evidence.go`
- `internal/evidence/execution.go`
- `internal/evidence/registration.go`
- `internal/workflow/run.go`
- `internal/claims/generic.go`
- `internal/validation/`
- `cmd/fid/resume.go`
- `cmd/fid/evidence.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-12-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-11-done.md` - actual predecessor hand-back.

### Out Of Scope

- External-artifact freshness targets, Review freshness, multiple Envelope Revisions, and Human briefs.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] A changed repository target, including one dirty byte at unchanged HEAD,
  stales prior repository-bound proof without changing Evidence identity. (VAL-048)
- [ ] An unchanged clean or canonical dirty target retains proof through
  recovery and branch rename; dirty coverage cannot use synthetic matching
  strings without actual observation. (VAL-049)
- [ ] Artifact loss preserves Evidence identity while changing only availability and policy outcome. (VAL-050)
- [ ] Fresh resume reconstructs bounded recovery context without prior chat. (VAL-051)
- [ ] Resume action eligibility matches actual mutation predicates. (VAL-052)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestRecoveredDirtyTargetStalesOldProofWithoutChangingEvidenceIdentity fails against a stored-target resume projection.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-12-done.md`: implementation,
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

- Start after task 11 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 4 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-12`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
