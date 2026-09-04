# Phase 22 - Task 13: Recovery archive format and validation

## Objective

Create bounded manifest-v2 recovery exports and validate their full identity, references, and effect semantics without mutation. Preserve unavailable Artifact metadata and project-source references.

## Validation Contract Coverage

- VAL-054: Manifest v2 contains all facts needed for recovery without project prose or SQLite.
- VAL-055: Valid hashes cannot conceal invalid recovery history or references.
- VAL-056: Archive extraction enforces size and path boundaries.
- VAL-057: Archive validation has no project or destination effects.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 10; public-interface.md section 5.
- Use stdlib ZIP/JSON and existing Event appliers; no runtime dependency.
- Bound actual expanded bytes and line sizes, reject unknown Effects and special/path-escaping entries.
- Manifest-v1 validation stays read-only; import-v1 is explicitly unsupported.
- Split export.go by archive writing/validation responsibilities while preserving its accepted behavior; it already exceeds the file target.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/export/`
- `internal/journal/event.go`
- `internal/journal/canonical.go`
- `internal/journal/replay.go`
- `internal/config/project_config.go`
- `cmd/fid/export.go`
- `cmd/fid/export_test.go`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-13-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-12-done.md` - actual predecessor hand-back.

### Out Of Scope

- Import publication, old-journal migration, direct database copy, and source-prose export.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Manifest v2 contains all facts needed for recovery without project prose or SQLite. (VAL-054)
- [ ] Valid hashes cannot conceal invalid recovery history or references. (VAL-055)
- [ ] Archive extraction enforces size and path boundaries. (VAL-056)
- [ ] Archive validation has no project or destination effects. (VAL-057)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects fails before semantic validation is added.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-13-done.md`: implementation,
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

- Start after task 12 closes and the independent task-12 recovery checkpoint passes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 5.
- Commit using Conventional Commits with `Refs: phase-22/task-13`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

