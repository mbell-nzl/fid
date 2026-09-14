# Phase 23 - Task 8: Generic Concerns and append-only Deviations

## Objective

Replace completion-linked Concerns with generic, subject-addressed Concerns
that carry a single irreversible disposition, and add append-only Deviations.
Neither is Evidence and neither blocks a guarded action.

## Validation Contract Coverage

- VAL-034: Concerns are generic, never completion-linked, and never block a
  guarded action by themselves.
- VAL-035: `concern.dispose` is a single irreversible transition with a
  mandatory reason and a real Finding for `converted-to-finding`.
- VAL-036: Deviations are append-only with explicit supersession.
- VAL-037: Neither a Concern nor a Deviation can be linked as Evidence.
- VAL-003: task 1's fence assertions for the three judgment Events converge
  here without losing their parity subject.


## Binding Constraints

- `contracts.md` section 8 is the exact authority; `schema.md` section 2 fixes
  `judgment_concerns` and `deviations` including their composite CHECKs.
- The new Concern table is created as `judgment_concerns` and is renamed to
  `concerns` only in task 15, after the V1 table is dropped. Effect rules in
  this task name `judgment_concerns`; task 15 owns that single-word change.
- Blocking work is expressed as a blocking Finding or a Stop, never as a
  deferred Concern (PRD-REV-008). An open Concern of any severity must leave
  both guarded actions unaffected.
- A Deviation has no lifecycle and no update path. A changed judgment is a new
  Deviation referencing its predecessor; self-reference is rejected.
- Concerns and Deviations record only the fields in the accepted contract. Do
  not add prose bodies, role columns, triage links, or completion links.
- DEC-192: delete exactly the `"concern.recorded"`, `"concern.disposed"`, and
  `"deviation.recorded"` keys and prove live-apply and replay parity for each.
  Converge task 1's fence assertions for those keys by substituting a
  still-fenced Event; never delete the parity assertion.
- Add the `concern.record`, `concern.dispose`, and `deviation.record` grant
  actions exactly here.
- The V1 `fid concern raise` and `fid concern address` leaves and their
  fail-closed writers stay in place until task 15 removes them. Do not delete
  or alias them in this task.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Concern and Deviation rules and effects.
- `internal/workflow/authority.go` - the three new grant actions.
- `internal/journal/project_effects.go`, `current_event_catalog.go` - the three
  fence keys, effects, and activity event names.
- `internal/activity/` - the added activity event names.
- `internal/cli/codes.go` - `concern-not-found`, `deviation-not-found`,
  `concern-already-disposed`.
- `cmd/fid/concern.go` - add `record`, `dispose`; `cmd/fid/deviation.go` - new.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-8-done.md`.

### Files To Read For Context

- `contracts.md` sections 8 and 13; `public-interface.md` sections 2, 3, 4, 4.1,
  5; `schema.md` sections 2 and 6.
- `task-7-done.md`; `execution-scope.md`; `v1-convergence-inventory.md`.

### Out Of Scope

- Removing the V1 Concern surface, its readers, or `issues_enriched`.
- Decision References, Human Outcomes, and Human text.

## Acceptance Criteria

- [ ] Every closed subject type is accepted, no completion column exists, and an
  open blocking-severity Concern leaves both guarded actions unaffected.
  (VAL-034)
- [ ] The first disposal succeeds; a second disposal, a missing reason, and an
  unknown Finding all reject with their exact codes. (VAL-035)
- [ ] No Deviation update path exists and a superseding Deviation references its
  predecessor with both readable. (VAL-036)
- [ ] `evidence.link` naming a Concern or a Deviation rejects on the closed
  evidence-type domain with no Event. (VAL-037)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestConcernDisposeIsSingleIrreversibleTransition` fails against the absent
generic Concern writer.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3: add `fid concern
  dispose` and `fid deviation list` to the accepted review-and-judgment block,
  per the accepted contract.
- `project-docs/docs/phases/phase-23/tasks/task-8-done.md` in the full format,
  including the three fence-key deletions and their parity proofs.
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

- This task belongs to QA Review window 3.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-8`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
