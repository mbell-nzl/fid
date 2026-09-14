# Phase 23 - Task 17: Final capability matrix and availability audit

## Objective

Advertise `reviews.coverage.v1` and `human.outcomes.v1` and flip the
`attention.v1` Human Outcome limit only after every replacement and removal
proof has passed, and prove every surface that reports availability agrees.

## Validation Contract Coverage

- VAL-047: `attention.v1` reports `human_outcome_resolution:true` only once the
  resolution path exists.
- VAL-066: The final binary capability matrix matches the accepted table exactly
  across every reporting surface.

## Binding Constraints

- `contracts.md` section 15 is the exact authority for the three changed rows.
  `work.claims.v1` stays false with owner phase 24 and
  `telemetry.trajectory.v1` false with owner phase 25.
- Schema presence alone does not enable a row. Change the matrix only after
  tasks 4 through 16 have all closed with their proofs; cite them.
- Capability discovery, `resume`, the export manifest profile, and the shared
  fake and fid fixtures must agree row for row. A disagreement is a defect in
  this task, never an accepted difference.
- Do not obscure a still-unavailable capability by renaming, merging, or
  omitting its row. Negative discovery of `work.claims.v1` remains mandatory.
- Remove the last fail-closed judgment guard and its tests here if any survived
  task 15, and prove absence rather than leaving a dormant code path.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/workflow/capabilities.go`.
- `internal/export/archive_profiles.go`.
- `cmd/fid/capabilities.go` and its fixtures.
- Matching in-package tests and `cmd/fid/integration/` capability fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-17-done.md`.

### Files To Read For Context

- `contracts.md` section 15; `public-interface.md` section 6.
- `task-16-done.md` and the close evidence of tasks 4 through 15.
- `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Any behavior change to the capabilities being advertised.
- The conformance corpus audit, which is task 18.

## Acceptance Criteria

- [ ] All four reporting surfaces produce identical rows equal to the accepted
  table, including both false rows and their owner phases. (VAL-066)
- [ ] `human_outcome_resolution` is true only in this final matrix and the path
  is demonstrated working. (VAL-047)
- [ ] No fail-closed judgment guard or dormant capability-unavailable path
  survives.
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestFinalCapabilityMatrixMatchesAcceptedTable` fails against the retained
Phase-23-unavailable rows.

## Documentation Updates

- `project-docs/specs/02_TRD.md` section 4: mark `reviews.coverage.v1` and
  `human.outcomes.v1` qualified, per the accepted contract.
- `project-docs/docs/phases/phase-23/tasks/task-17-done.md` in the full format,
  citing the closing evidence for every prerequisite task.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task belongs to QA Review window 6.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-17`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
