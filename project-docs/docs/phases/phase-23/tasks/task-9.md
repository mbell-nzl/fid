# Phase 23 - Task 9: Decision References and protected Decisions

## Objective

Record operational Decision References that point at canonical project prose,
and make an unsatisfied protected Decision block exactly its declared actions
until a matching approved reference and Human Outcome both exist.

## Validation Contract Coverage

- VAL-038: `decision.record-reference` stores a reference only and rejects a
  digest mismatch.
- VAL-039: An unsatisfied protected Decision blocks exactly its declared
  affected actions.
- VAL-040: Satisfaction requires the approved reference and the approved Human
  Outcome together.
- VAL-041: Activation raises one Stop per unsatisfied protected Decision and the
  satisfying pair resolves that Stop and its Attention in one Event.
- VAL-042: `linked_to_decision` becomes an implemented transition rule.
- VAL-003: task 1's fence assertion for `decision.referenced` converges here
  without losing its parity subject.


## Binding Constraints

- `contracts.md` section 9 is the exact authority; `schema.md` section 2 fixes
  `decision_refs`, and task 2 already landed `envelope_protected_decisions`.
- fid stores no decision body. Verify the `Source` digest against the observed
  project file under the write lock and reject `source-digest-mismatch` with no
  Event on any mismatch.
- A `rejected` or `deferred` reference never satisfies a protected Decision.
  Authority matching uses the declared required capabilities, never a role
  string or an actor-type heuristic.
- Envelope activation is permitted with an unsatisfied protected Decision;
  execution of the affected actions is not. Raise exactly one Stop per
  unsatisfied protected Decision, subject Envelope, blocked actions exactly the
  declared affected actions, resolver capability `stop.resolve`, resolution
  condition `decision-satisfied` naming that class.
- Satisfying the pair resolves that Stop and its Attention in the same Event,
  leaving every other Stop untouched.
- `linked_to_decision` moves from reserved to implemented and is evaluated
  against `decision_refs` in the Issue's scope. Do not reintroduce the removed
  escalation seed or any role-bound rule parameter.
- Human Outcome recording is task 10. Implement the satisfaction predicate over
  `human_outcomes` here and use a seeded projection fixture until task 10
  supplies the writer; do not add a temporary writer.
- DEC-192: delete exactly the `"decision.referenced"` key and prove its
  live-apply and replay parity; converge task 1's fence assertion for that key.
- Add the `decision.record-reference` grant action exactly here.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Decision Reference rules, effects, satisfaction.
- `internal/workflow/envelope.go`, `evaluation.go`, `resolve.go`, `rules.go`,
  `registry.go`, `authority.go`.
- `internal/issues/transition_rules.go` - the `linked_to_decision` evaluator.
- `internal/journal/project_effects.go`, `current_event_catalog.go`.
- `internal/cli/codes.go` - `protected-decision-required`,
  `source-digest-mismatch`, `decision-not-found`.
- `cmd/fid/decision.go` - new leaf wiring for `record-reference` and `list`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-9-done.md`.

### Files To Read For Context

- `contracts.md` sections 2, 9, 10; `public-interface.md` sections 2, 3, 4, 4.1,
  5; `schema.md` sections 1, 2, 6.2.
- `task-8-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Human Outcome recording, authority, and provenance.
- Attention resolution through a Human Outcome, which is task 11.

## Acceptance Criteria

- [ ] No prose body column exists; a matching digest records and a mismatch
  rejects with no Event. (VAL-038)
- [ ] Only the declared affected action is blocked, with the exact code.
  (VAL-039)
- [ ] Only the both-approved state satisfies; neither fact alone and no
  `rejected` or `deferred` reference does. (VAL-040)
- [ ] Two unsatisfied protected Decisions raise exactly two Stops, and
  satisfying one resolves exactly that Stop and Attention in one Event.
  (VAL-041)
- [ ] `linked_to_decision` is listed as implemented and is satisfied only with a
  qualifying reference. (VAL-042)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestProtectedDecisionBlocksOnlyDeclaredAffectedActions` fails against the
absent protected-Decision evaluation.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3 and section 5: add
  `fid decision list` and the Phase 23 closed codes owned by this task.
- `project-docs/docs/phases/phase-23/tasks/task-9-done.md` in the full format,
  including the fence-key deletion and its parity proof.
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
- Commit using Conventional Commits with `Refs: phase-23/task-9`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
