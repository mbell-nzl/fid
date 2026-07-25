# Phase 19 - Task 2: Workflow authoring core (validation + effect builder)

## Objective

Add the pure domain core that turns a workflow definition request into
validated, deterministically ordered journal effects, and the assignment
request validator that Task 4 will use. No CLI, no journal append, no DB
write in this task — this is the de-risking step before the command surface,
exactly as Phase 18's registry loader preceded the consumer swap.

## Validation Contract Coverage

- VAL-006: the effect builder produces a deterministic, ordered effect set
  matching the golden fixture across repeated builds.
- VAL-007: every enumerated invalid definition is rejected with a distinct,
  value-naming error; a valid definition is accepted.
- VAL-008: cancel-edge materialization behaves per DEC-077 for one, zero, and
  multiple `cancelled`-category states.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **Pure core, no side effects.** The new code takes a parsed request plus a
  read-only `*sql.DB` (for duplicate-slug and workflow-existence lookups) and
  returns either a validation error or a `[]journal.Effect`. It must not
  append to the journal, must not open a write lock, and must not execute any
  `INSERT`/`UPDATE`/`DELETE`.
- **Effect order is a contract, not an accident** (VAL-006). Emit in exactly
  this order: the `workflows` insert; `workflow_states` inserts in declared
  order with `ordinal` = 1-based declaration index; `workflow_transitions`
  inserts in declared order, then materialized cancel edges in declared-state
  order; `workflow_transition_requires` inserts in declared `--require`
  order. Two builds of the same request must produce an identical slice — do
  not iterate a Go map to emit effects.
- **Validation rules (DEC-076, DEC-077, DEC-078).** Reject, each with its own
  error naming the offending value:
  - workflow slug not matching `^[a-z0-9][a-z0-9-]{0,63}$`;
  - workflow slug already present in `workflows`;
  - empty `--states`;
  - state name not matching `^[a-z][a-z0-9_]{0,63}$`;
  - duplicate state name;
  - `len(--categories) != len(--states)`;
  - category outside `unstarted|started|completed|cancelled`;
  - more than one `cancelled`-category state;
  - transition whose `from` or `to` is not a declared state;
  - duplicate transition;
  - `--require` naming an edge that is not a declared or materialized
    transition;
  - rule name outside the implemented set;
  - reserved-unimplemented rule name, with a version-scoped message:
    `required_gates_current` → v0.2b, `has_completion` → v0.2c,
    `linked_to_decision` → v0.2d;
  - `--require` params that are not valid JSON, or `all_acs_verified` without
    a non-empty `roles` array.
- **Implemented rule vocabulary is exactly four names**: `all_acs_verified`,
  `has_completion_report`, `concerns_addressed`,
  `all_child_tasks_done_or_cancelled`. Source them from the existing
  `internal/workflow` constants; do not re-declare string literals.
- **Cancel-edge materialization (DEC-077).** Exactly one `cancelled`-category
  state → materialize `<state> → <cancelled-state>` for every non-cancelled
  declared state, skipping any the user already declared (no duplicate
  effect). Zero cancelled states → materialize nothing. More than one → a
  validation error (listed above).
- **Initial-state resolution belongs here too**, as a pure exported helper:
  the first declared state whose category is `unstarted`, else the first
  declared state. Task 4 consumes it; do not duplicate the rule there.
- **Assignment request validation** (used by Task 4): exactly one of kind
  scope or project scope; kind scope requires the workflow's `applies_to` to
  equal the kind or be `any` (NULL); project scope requires `applies_to` to be
  `any`; the workflow must exist. Each violation returns its own error.
- **Unicode arrow normalization is a parsing concern, not a core concern.**
  The core receives already-split `from`/`to` state names. Task 3 owns
  normalization.
- **Respect the file/function size budgets** (`AGENTS.md`): non-test
  functions ≤ 60 lines, non-test files ≤ 400 lines. Split the validator and
  the effect builder into separate files if either approaches the limit.

## Scope

### Files To Create Or Modify

- `internal/workflow/authoring.go` (new) — definition request type,
  validation, cancel-edge materialization, initial-state helper
- `internal/workflow/authoring_effects.go` (new) — effect builder
- `internal/workflow/assignment.go` (new) — assignment request validation
- `internal/workflow/authoring_test.go`,
  `internal/workflow/authoring_effects_test.go`,
  `internal/workflow/assignment_test.go` (new)
- `internal/workflow/rules.go` — export the implemented-rule set if a helper
  is needed; do not change existing rule semantics

### Files To Read For Context

- `internal/workflow/registry.go`, `internal/workflow/workflow.go`,
  `internal/workflow/rules.go`
- `internal/journal/event.go`, `internal/journal/project_effects.go` (effect
  shape: `Op`, `Table`, `PK`, `Values`)
- `migrations/0007_v02a_workflow_tables.up.sql` (the shape the effects must
  reproduce), `migrations/0009_v02a_workflow_assignments.up.sql` (Task 1)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow define/assign`
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-076, DEC-077, DEC-078)

### Out Of Scope

- Any cobra command, flag parsing, or arrow normalization (Task 3).
- Any journal append, write lock, or projection call (Tasks 3, 4).
- Any change to `internal/issues` (Task 4).
- Any change to the four already-implemented rule evaluators.

## Acceptance Criteria

- [ ] The effect builder emits the contracted order and produces an identical
  slice on repeated builds of the same request (VAL-006)
- [ ] Every enumerated rejection case in Binding Constraints returns a
  distinct error naming the offending value, and a valid definition is
  accepted (VAL-007)
- [ ] One-cancelled-state fixtures materialize cancel edges for every
  non-cancelled state, deduplicated against explicitly declared cancel edges
  (VAL-008)
- [ ] A zero-cancelled-state fixture materializes no cancel edge (VAL-008)
- [ ] A two-cancelled-state fixture is rejected (VAL-007/VAL-008)
- [ ] The initial-state helper returns the first `unstarted` state, and the
  first declared state when no `unstarted` state exists
- [ ] Assignment validation rejects unknown workflow, incompatible
  `applies_to` for kind scope, non-`any` workflow for project scope, both
  scopes, and neither scope
- [ ] No production path in this task writes to the DB or the journal
- [ ] Existing behavior is unchanged; the full suite stays green

## First Failing Test

`TestBuildDefinitionEffectsIsOrderedAndDeterministic`: build effects for a
fixture definition twice and deep-equal both against the golden ordered slice
— fails until the effect builder exists.

## Documentation Updates

- None expected in this task; the CLI spec is updated with the command
  surface in Task 3. Record any core-shape decision in the completion report.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. A non-terminating command is
never `PASS`.

Conditional command required by this task:

- None. This task is additive and touches no migration or CLI path.

## Notes For Developer

- Depends on Task 1 (the `workflow_assignments` table must exist for the
  assignment validator's lookups).
- `internal/workflow` is a COV-001 domain package with an 80% coverage
  threshold. This task adds a lot of new code to it; keep test coverage of
  the new files high or the phase gate will regress.
- Table-driven tests are the project norm and fit the rejection matrix
  directly.
- Escalation tripwire: if a validation rule in the Binding Constraints list
  turns out to be unsatisfiable or contradictory against the specs, stop and
  escalate. Do not relax the rule or silently widen the accepted vocabulary.
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-2-done.md`.
