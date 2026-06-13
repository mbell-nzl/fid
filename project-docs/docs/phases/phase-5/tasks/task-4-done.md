# Phase 5 - Task 4: `internal/workflow` hard-coded graphs + rule functions - Completion Report

## Summary of what was implemented

- Added the new `internal/workflow` domain package with the hard-coded v0.1
  `task-default`, `epic-default`, and `subtask-default` workflow graphs from
  TRD §4.6.
- Implemented exported graph lookups for workflow definitions, state-category
  mapping, allowed-edge checks, and per-edge rule references, including the
  universal non-cancelled `-> cancelled` edge.
- Implemented the pure rule functions `all_acs_verified`,
  `has_completion_report`, `concerns_addressed`, and
  `all_child_tasks_done_or_cancelled` as DB-agnostic Go functions returning
  structured unmet-condition details for Task 5 to render.
- Added exhaustive workflow and rule package tests; `internal/workflow` reached
  95.5% statement coverage.

## Files created or modified

- `internal/workflow/workflow.go`
- `internal/workflow/workflow_test.go`
- `internal/workflow/rules.go`
- `internal/workflow/rules_test.go`
- `project-docs/docs/phases/phase-5/tasks/task-4-done.md`

## Documentation updated

- `project-docs/specs/02_TRD.md` - verified no change. Task 4 implements the
  default workflow tables and rule names from §4.5 and §4.6 exactly as written,
  so no spec text change was needed.
- `project-docs/docs/phases/phase-5/tasks/task-4-done.md` - added this
  completion report with VAL coverage and gate evidence.

## Validation Contract Coverage

- `VAL-012` - PASS. `TestConcernsAddressed` proves
  `concerns_addressed` returns unsatisfied when open concern IDs are present and
  satisfied when no open concerns remain.
- `VAL-006` - ADVANCED. `TestWorkflowDefinitionsMatchTRDExactly`,
  `TestCategoryForStateMatchesDefaultWorkflows`, and `TestIsAllowedEdge` prove
  the task, epic, and subtask workflow graphs expose the exact TRD §4.6 states,
  status categories, allowed edges, disallowed edges, and universal cancel
  semantics.
- `VAL-008` - ADVANCED. `TestRulesForEdgeMatchesTRDRequirements`,
  `TestAllACsVerified`, and `TestHasCompletionReport` prove the
  `in_progress -> dev_complete` task edge carries the TRD rule references and
  the pure rule functions return structured unmet AC IDs and completion-report
  issue IDs.
- `VAL-009` - ADVANCED. `TestRulesForEdgeMatchesTRDRequirements` and
  `TestAllACsVerified` prove the `qa_review -> qa_pass` edge carries the
  `qa_lead` role parameter and the rule function honours that required-role
  set.
- `VAL-010` - ADVANCED. `TestRulesForEdgeMatchesTRDRequirements` and
  `TestAllChildTasksDoneOrCancelled` prove the epic
  `in_progress -> wrapped` edge carries the child-task rule and that unmet
  results name the blocking child IDs.

## Tests run

- `go test ./internal/workflow -run TestIsAllowedEdge -count=1` - FAIL before
  implementation (`undefined: IsAllowedEdge`)
- `go test ./internal/workflow -run TestIsAllowedEdge -count=1` - PASS
- `go test ./internal/workflow -run 'Test(CategoryForStateMatchesDefaultWorkflows|RulesForEdgeMatchesTRDRequirements|AllACsVerified|HasCompletionReport|ConcernsAddressed|AllChildTasksDoneOrCancelled)$' -count=1`
  - FAIL before rule implementation (`undefined: AllACsVerified`,
  `undefined: HasCompletionReport`, `undefined: ConcernStatus`)
- `go test ./internal/workflow -run 'Test(IsAllowedEdge|CategoryForStateMatchesDefaultWorkflows|RulesForEdgeMatchesTRDRequirements|AllACsVerified|HasCompletionReport|ConcernsAddressed|AllChildTasksDoneOrCancelled)$' -count=1`
  - PASS
- `go test ./internal/workflow -count=1` - PASS
- `go test -cover ./internal/workflow` - PASS (`coverage: 95.5% of statements`)
- `go test -race -coverprofile=coverage.out ./...` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `cde39fd`

## Commit subject

- `feat(workflow): add default workflow graphs`
