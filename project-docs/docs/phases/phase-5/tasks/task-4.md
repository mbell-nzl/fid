# Phase 5 - Task 4: `internal/workflow` hard-coded graphs + rule functions

## Objective

Create the `internal/workflow` package with the v0.1 hard-coded default
workflows (`task-default`, `epic-default`, `subtask-default`) from TRD §4.6 —
states, status categories, allowed transition edges, and per-edge rule
references — plus the four named transition-rule functions: `all_acs_verified`
(parameterised by required verifier roles), `has_completion_report`,
`concerns_addressed`, and `all_child_tasks_done_or_cancelled`. This is the pure
domain layer that Task 5 wires into `fid issue transition`. No CLI in this task.

## Requirements

- PRD refs: PRD-WF-005 (default workflows), PRD-WF-006 (state→category),
  PRD-WF-007 (`all_acs_verified` parameterised by roles), PRD-WF-008
  (`has_completion_report`), PRD-WF-009 (`concerns_addressed`), PRD-WF-010
  (`all_child_tasks_done_or_cancelled`), PRD-WF-013 (cancel from any
  non-cancelled state)
- TRD refs: TRD §4.6 (default workflow tables), TRD-ARCH-009 (rules as named Go
  functions in v0.1), TRD §4.5 (app-layer owns transition rules)
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/workflow` ≥ 80%)

## Validation Contract Coverage

- VAL-012: `concerns_addressed` returns satisfied iff no open concerns; both
  states proven by unit test (no default edge wires it).
- VAL-006 (advances): default workflow graphs + category mapping for all three
  kinds.
- VAL-008/009/010 (advances): rule functions and edge legality the transition
  command will exercise in Task 5.

## Scope

### Files To Create Or Modify

- `internal/workflow/workflow.go` (new) - workflow definitions: per-`workflow_id`
  states with categories, allowed edges, per-edge rule references; lookups for
  "is edge allowed", "category of state", "rules for edge", and the universal
  cancel edge (PRD-WF-013).
- `internal/workflow/rules.go` (new) - the four rule functions operating over a
  small data interface (counts/IDs of verifications, completions, open
  concerns, child task states) so they stay DB-agnostic and unit-testable;
  return structured results naming unmet conditions + entity IDs.
- `internal/workflow/workflow_test.go`, `internal/workflow/rules_test.go` (new).

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.6 (the three workflow tables), §4.5
- `internal/issues/issue.go` (`statusCategory` mapping already in use — keep
  consistent or migrate the source of truth here)
- `project-docs/specs/03_TRD_SCHEMA.md` (ac_verifications, task_completions,
  concerns, issues parent/child columns)

### Out Of Scope

- `fid issue transition` CLI and the DB-backed rule-input gathering (Task 5).
- Configurable workflow tables/engine (v0.2, PRD-WF-011/012).
- Escalation workflow (`escalation-default`, v0.2).

## Acceptance Criteria

- [ ] All three default workflows expose their states with the exact status
  categories from TRD §4.6, and every allowed edge from §4.6 is represented;
  disallowed edges are reported as not allowed. (VAL-006)
- [ ] The universal "any non-cancelled state → cancelled" edge is allowed and a
  cancelled→cancelled edge is not. (VAL-006, supports VAL-011)
- [ ] `all_acs_verified` is parameterised by a role set and returns unmet
  details (naming the failing AC IDs) when any required role lacks a passing
  verification. (VAL-008/009)
- [ ] `has_completion_report` returns satisfied iff a completion report exists
  for the issue. (VAL-008)
- [ ] `concerns_addressed` returns satisfied iff there are no open concerns;
  both states unit-tested. (VAL-012)
- [ ] `all_child_tasks_done_or_cancelled` returns satisfied iff every child task
  is done/cancelled and names blocking child IDs otherwise. (VAL-010)
- [ ] COV-001 holds for `internal/workflow` (≥ 80%).

## TDD Plan

- First failing test: `IsAllowedEdge("task-default", "todo", "in_progress")`
  true and a disallowed edge false. Fails (no package yet).
- Additional tests: category mapping per state per workflow; cancel edge from
  several non-cancelled states; each rule function both states with structured
  unmet output (IDs present); parameterised role sets for `all_acs_verified`
  (developer vs qa_lead).
- Regression risks: keep the status-category mapping identical to the one
  `internal/issues` currently uses (avoid a second divergent source of truth —
  prefer making this package the source and having issues consume it, or
  document why duplication is acceptable this phase).

## Documentation Updates

- Completion report: cite VAL-006/008/009/010/012 advancement and confirm the
  rule functions match TRD §4.6 edges verbatim.
- `project-docs/specs/02_TRD.md` - `verified-no-change` expected; the package
  implements §4.6 as written. Record the verdict rationale.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: status-category values must match those used by
  `internal/issues` read model (`statusCategory`).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command as
`PASS` (or `BLOCKED` with reason) in the completion report.

Conditional commands:

- None.

## Notes For Developer

- Keep rule functions DB-agnostic: accept resolved inputs (e.g. a struct of
  verification statuses by role + AC IDs, completion presence, open-concern
  IDs, child task states). Task 5 supplies these from the projection. This keeps
  rules pure and unit-testable per TRD-ARCH-009.
- Rule results must carry enough structure for Task 5 to render PRD-WF-004 rich
  errors (unmet condition name + entity IDs).
- The v0.1 default task workflow does not wire `concerns_addressed` to any edge;
  implement and unit-test it anyway (PRD-WF-009).
