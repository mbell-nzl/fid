# Phase 23 - Task 5: Finding dispositions and resolution authority

## Objective

Implement the Finding lifecycle: Builder-recorded `addressed` assertions,
resolution and inapplicability only through a later qualifying independent
Review, and a derived Finding state.

## Validation Contract Coverage

- VAL-020: Findings are created only with their Review, atomically, and are
  immutable.
- VAL-021: `finding.address` requires an active mutating Claim and leaves the
  Finding unresolved and still blocking.
- VAL-022: Only a later qualifying independent Review resolves or marks a
  Finding inapplicable, and no Review resolves a Finding it created.
- VAL-023: `Finding.state` is derived from the latest disposition.
- VAL-003: task 1's fence assertion for `finding.addressed` converges here
  without losing its parity subject.


## Binding Constraints

- `contracts.md` section 5 is the exact authority; `schema.md` section 2 fixes
  `finding_dispositions` including its discriminator CHECK and the partial
  unique index that caps terminal dispositions at one.
- `addressed` is an assertion of rework, never a resolution. It must not change
  blocking state, satisfy an Obligation, or close a Stop.
- Resolution arrives only through `resolves_finding_ids` on a later Review
  Record that binds the same Issue and revision, satisfies task 4 independence,
  and has verdict `pass` or `changes-required`. There is no `finding.resolve`
  leaf and none may be added.
- A Review may not resolve a Finding it created, nor one already resolved or
  inapplicable. Both rejections name the exact Finding.
- DEC-192: delete exactly the `"finding.addressed"` key and prove its
  live-apply and replay parity. Converge task 1's fence assertion for that key
  by substituting a still-fenced Event.
- Add the `finding.address` grant action exactly here.
- Do not change Obligation evaluation in this task; the
  `finding-unresolved` reason is delivered with the review Evidence path in
  task 6.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Finding creation effects, dispositions, derived state.
- `internal/workflow/authority.go` - the one new grant action.
- `internal/journal/project_effects.go` - the one fence key and its effects.
- `internal/cli/codes.go` - `finding-not-found`, `finding-already-disposed`.
- `cmd/fid/finding.go` - new leaf wiring for `address` and `list`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-5-done.md`.

### Files To Read For Context

- `contracts.md` sections 4 and 5; `public-interface.md` sections 3, 4, 4.1, 5;
  `schema.md` section 2.
- `task-4-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Obligation evaluation changes, Evidence linking, budgets, and Human text
  beyond the minimum required to return a result.

## Acceptance Criteria

- [ ] Findings have no independent create or update path and appear only inside
  their Review Event. (VAL-020)
- [ ] Address succeeds only under an active mutating Claim, and afterwards the
  Finding is `addressed`, unresolved, and still blocking. (VAL-021)
- [ ] Self-resolution, non-independent resolution, and a second terminal
  disposition all reject; a valid independent resolution records its qualifying
  Review. (VAL-022)
- [ ] A supplied `state` rejects and the derived state follows the dispositions.
  (VAL-023)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestFindingAddressDoesNotResolveOrUnblock` fails against the absent
disposition writer.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-5-done.md` in the full format,
  including the fence-key deletion and its parity proof.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
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

- This task belongs to QA Review window 2.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-5`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
